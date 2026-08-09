import {
  hasMeaningfulValue,
  parseRuntimeScenario,
  parseRuntimeScenarioResolveResponse,
} from "./logic.js";
import type {
  CreateOnboardingOptions,
  EventType,
  Onboarding,
  RuntimeScenario,
  RuntimeStep,
  StartOptions,
} from "./types.js";
import { createOnboardingView, type OnboardingView } from "./view.js";

export type {
  CreateOnboardingOptions,
  FrontendStepData,
  Onboarding,
  OnboardingUserState,
  StartOptions,
} from "./types.js";

type Progress = {
  version: 2;
  projectKey: string;
  userId: string;
  usesTestToken: boolean;
  isTest: boolean;
  sessionId: string | null;
  scenario: RuntimeScenario;
  stepIndex: number;
  accepted: boolean;
  shownStepIds: string[];
};

type Highlight = {
  element: HTMLElement | SVGElement;
  outline: string;
  outlinePriority: string;
  outlineOffset: string;
  outlineOffsetPriority: string;
};

type EventData = Record<string, string>;

function assertOption(value: string, name: string): string {
  const trimmed = value.trim();
  if (!trimmed) throw new TypeError(`${name} must be a non-empty string`);
  return trimmed;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function parseProgress(value: unknown, projectKey: string): Progress | null {
  if (!isRecord(value)) return null;

  try {
    const scenario = parseRuntimeScenario(value.scenario);
    if (
      value.version !== 2 ||
      value.projectKey !== projectKey ||
      typeof value.userId !== "string" ||
      typeof value.usesTestToken !== "boolean" ||
      typeof value.isTest !== "boolean" ||
      (value.sessionId !== null && typeof value.sessionId !== "string") ||
      !Number.isInteger(value.stepIndex) ||
      Number(value.stepIndex) < 0 ||
      Number(value.stepIndex) >= scenario.steps.length ||
      typeof value.accepted !== "boolean" ||
      !Array.isArray(value.shownStepIds) ||
      !value.shownStepIds.every((id) => typeof id === "string")
    ) {
      return null;
    }

    if (
      (value.isTest && value.sessionId !== null) ||
      (!value.isTest && value.accepted && !value.sessionId) ||
      (!value.accepted && value.sessionId !== null)
    ) {
      return null;
    }

    return {
      version: 2,
      projectKey,
      userId: value.userId,
      usesTestToken: value.usesTestToken,
      isTest: value.isTest,
      sessionId: value.sessionId,
      scenario,
      stepIndex: Number(value.stepIndex),
      accepted: value.accepted,
      shownStepIds: [...value.shownStepIds],
    };
  } catch {
    return null;
  }
}

export function createOnboarding(options: CreateOnboardingOptions): Onboarding {
  const projectKey = assertOption(options.projectKey, "projectKey");
  const runtimeUrl = assertOption(options.runtimeUrl, "runtimeUrl").replace(/\/+$/, "");
  const onUserStateChange = options.onUserStateChange;
  const storageKey = `@interactive-onboarding/sdk:v2:${projectKey}`;

  let progress: Progress | null = null;
  let view: OnboardingView | null = null;
  let highlighted: Highlight | null = null;
  let target: Element | null = null;
  let targetEvents: { name: string; listener: EventListener }[] = [];
  let observer: MutationObserver | null = null;
  let observerTimer: number | null = null;
  let routeTimer: number | null = null;
  let positionFrame: number | null = null;
  let listenersAttached = false;
  let completingStepId: string | null = null;
  let accepting = false;
  let startController: AbortController | null = null;
  let sessionController: AbortController | null = null;
  let startPromise: Promise<void> | null = null;
  let pendingDismiss: (() => void) | null = null;
  let generation = 0;
  let eventQueue: Promise<void> = Promise.resolve();

  function requireBrowser(): void {
    if (
      typeof window === "undefined" ||
      typeof document === "undefined" ||
      typeof customElements === "undefined"
    ) {
      throw new Error("Onboarding SDK can only start in a browser");
    }
  }

  function readProgress(): Progress | null {
    try {
      const raw = window.sessionStorage.getItem(storageKey);
      if (!raw) return null;
      const stored = parseProgress(JSON.parse(raw) as unknown, projectKey);
      if (!stored) window.sessionStorage.removeItem(storageKey);
      return stored;
    } catch {
      return null;
    }
  }

  function writeProgress(): void {
    if (!progress) return;
    try {
      window.sessionStorage.setItem(storageKey, JSON.stringify(progress));
    } catch {
      // Storage may be unavailable in privacy-restricted browser contexts.
    }
  }

  function removeProgress(): void {
    try {
      window.sessionStorage.removeItem(storageKey);
    } catch {
      // Storage may be unavailable in privacy-restricted browser contexts.
    }
  }

  async function post(
    path: string,
    body: unknown,
    signal?: AbortSignal,
    keepalive = false,
    testToken?: string,
  ): Promise<unknown> {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      "X-Project-Key": projectKey,
    };
    if (testToken) headers["X-Scenario-Test-Token"] = testToken;

    const request: RequestInit = {
      method: "POST",
      credentials: "same-origin",
      keepalive,
      headers,
      body: JSON.stringify(body),
    };
    if (signal) request.signal = signal;

    const response = await fetch(`${runtimeUrl}${path}`, request);
    if (response.status === 204) return null;
    if (!response.ok) {
      const detail = (await response.text()).slice(0, 300);
      throw new Error(`Onboarding API ${path} failed (${response.status})${detail ? `: ${detail}` : ""}`);
    }
    return response.json() as Promise<unknown>;
  }

  function recordEvent(
    state: Progress,
    type: EventType,
    data: EventData = {},
    stepId?: string,
  ): Promise<void> {
    if (state.isTest || !state.sessionId) return Promise.resolve();

    const payload: Record<string, unknown> = {
      id: crypto.randomUUID(),
      session_id: state.sessionId,
      type,
      data,
      occurred_at: new Date().toISOString(),
    };
    if (stepId) payload.step_id = stepId;

    eventQueue = eventQueue
      .then(() => post("/events", payload, undefined, true))
      .then(() => undefined)
      .catch((error: unknown) => {
        console.warn("Onboarding tracking event was not delivered", error);
      });
    return eventQueue;
  }

  function restoreHighlight(): void {
    if (!highlighted) return;
    const { element } = highlighted;

    if (highlighted.outline) {
      element.style.setProperty("outline", highlighted.outline, highlighted.outlinePriority);
    } else {
      element.style.removeProperty("outline");
    }
    if (highlighted.outlineOffset) {
      element.style.setProperty(
        "outline-offset",
        highlighted.outlineOffset,
        highlighted.outlineOffsetPriority,
      );
    } else {
      element.style.removeProperty("outline-offset");
    }
    highlighted = null;
  }

  function highlight(element: Element): void {
    restoreHighlight();
    if (!(element instanceof HTMLElement || element instanceof SVGElement)) return;

    highlighted = {
      element,
      outline: element.style.getPropertyValue("outline"),
      outlinePriority: element.style.getPropertyPriority("outline"),
      outlineOffset: element.style.getPropertyValue("outline-offset"),
      outlineOffsetPriority: element.style.getPropertyPriority("outline-offset"),
    };
    element.style.setProperty("outline", "3px solid #00aaff", "important");
    element.style.setProperty("outline-offset", "4px", "important");
  }

  function stopWaitingForTarget(): void {
    observer?.disconnect();
    observer = null;
    if (observerTimer !== null) window.clearTimeout(observerTimer);
    observerTimer = null;
  }

  function clearTarget(): void {
    stopWaitingForTarget();
    if (target) {
      for (const event of targetEvents) target.removeEventListener(event.name, event.listener);
    }
    targetEvents = [];
    target = null;
    pendingDismiss = null;
    restoreHighlight();
    view?.hideCard();
  }

  function positionView(): void {
    if (!view || !target || positionFrame !== null) return;
    positionFrame = window.requestAnimationFrame(() => {
      positionFrame = null;
      if (view && target) view.positionNear(target);
    });
  }

  const onViewportChange = (): void => positionView();

  const onEscape = (event: KeyboardEvent): void => {
    if (event.key !== "Escape" || accepting) return;
    if (progress) void skipOnboarding("user_closed");
    else pendingDismiss?.();
  };

  function ensureView(): OnboardingView {
    if (!view) {
      view = createOnboardingView();
      document.body.append(view);
    }
    if (!listenersAttached) {
      document.addEventListener("keydown", onEscape);
      window.addEventListener("resize", onViewportChange);
      window.addEventListener("scroll", onViewportChange, { capture: true, passive: true });
      listenersAttached = true;
    }
    return view;
  }

  function currentStep(state = progress): RuntimeStep | null {
    return state?.scenario.steps[state.stepIndex] ?? null;
  }

  function findTarget(step: RuntimeStep): Element | null {
    const selector = `[data-onboarding-id="${CSS.escape(step.element.key)}"]`;
    return document.querySelector(selector);
  }

  function showStep(step: RuntimeStep, element: Element): void {
    if (!progress || currentStep()?.id !== step.id) return;
    stopWaitingForTarget();
    clearTarget();
    target = element;
    highlight(element);

    const state = progress.shownStepIds.includes(step.id)
      ? progress
      : { ...progress, shownStepIds: [...progress.shownStepIds, step.id] };
    if (state !== progress) {
      progress = state;
      writeProgress();
      void recordEvent(state, "step_shown", {}, step.id);
    }

    ensureView().showHint({
      title: step.title,
      description: step.description,
      position: step.step_num,
      total: state.scenario.steps.length,
      onNext: () => void completeCurrentStep(),
      onClose: () => void skipOnboarding("user_closed"),
    });
    ensureView().positionNear(element);

    if (step.frontend_data.advance.mode === "target_event") {
      const eventName = step.frontend_data.advance.event;
      const listener: EventListener = (event) => {
        if (currentStep()?.id !== step.id) return;
        if (eventName === "change" && !hasMeaningfulValue(event.target)) return;
        void completeCurrentStep();
      };
      const eventNames = eventName === "change" ? ["input", "change"] : [eventName];
      targetEvents = eventNames.map((name) => ({ name, listener }));
      for (const event of targetEvents) element.addEventListener(event.name, event.listener);
    }
  }

  function advanceAfterMissingTarget(step: RuntimeStep): void {
    if (!progress || currentStep()?.id !== step.id) return;
    const state = progress;
    void recordEvent(state, "step_skipped", { reason: "target_not_found" }, step.id);
    clearTarget();

    if (state.stepIndex === state.scenario.steps.length - 1) {
      void recordEvent(state, "onboarding_skipped", { reason: "target_not_found" });
      progress = null;
      removeProgress();
      return;
    }

    progress = { ...state, stepIndex: state.stepIndex + 1 };
    writeProgress();
    renderForCurrentPage();
  }

  function waitForTarget(step: RuntimeStep): void {
    clearTarget();
    const found = findTarget(step);
    if (found) {
      showStep(step, found);
      return;
    }

    observer = new MutationObserver(() => {
      const added = findTarget(step);
      if (added) showStep(step, added);
    });
    observer.observe(document.documentElement, { childList: true, subtree: true });
    observerTimer = window.setTimeout(() => advanceAfterMissingTarget(step), 5_000);
  }

  function expectedPath(state: Progress, path: string): boolean {
    return state.scenario.page_pattern === path || state.scenario.steps.some(
      (step) => step.frontend_data.page_path === path,
    );
  }

  function renderForCurrentPage(): void {
    if (!progress) return;
    ensureView();
    if (routeTimer !== null) window.clearTimeout(routeTimer);
    routeTimer = null;

    const state = progress;
    const step = currentStep(state);
    if (!step) return;
    const path = window.location.pathname;

    if (!state.accepted) {
      if (!expectedPath(state, path)) {
        void skipOnboarding("left_expected_path");
        return;
      }
      clearTarget();
      ensureView().showInvitation({
        name: state.scenario.name,
        description: state.scenario.description,
        onAccept: () => void acceptInvitation(),
        onDecline: () => void skipOnboarding("invitation_declined"),
      });
      return;
    }

    if (step.frontend_data.page_path === path) {
      waitForTarget(step);
      return;
    }

    clearTarget();
    if (expectedPath(state, path)) return;

    if (
      step.frontend_data.advance.mode === "manual" &&
      state.shownStepIds.includes(step.id)
    ) {
      // Gives success-page code one turn to call completeCurrentStep().
      routeTimer = window.setTimeout(() => void skipOnboarding("left_expected_path"), 0);
      return;
    }
    void skipOnboarding("left_expected_path");
  }

  async function acceptInvitation(): Promise<void> {
    const state = progress;
    if (!state || state.accepted || accepting) return;

    accepting = true;
    view?.setError(null);
    view?.setBusy(true);
    try {
      let sessionId: string | null = null;
      if (!state.isTest) {
        const controller = new AbortController();
        sessionController = controller;
        sessionId = await createSession(state.scenario.id, state.userId, controller.signal);
        if (progress !== state) return;
      }

      progress = { ...state, accepted: true, sessionId };
      writeProgress();
      const step = currentStep();
      if (!step) return;

      if (step.frontend_data.page_path === window.location.pathname) {
        renderForCurrentPage();
      } else {
        clearTarget();
        const destination = new URL(step.frontend_data.page_path, window.location.href);
        window.location.assign(`${destination.pathname}${destination.search}${destination.hash}`);
      }
    } catch (error) {
      if ((error as { name?: unknown }).name !== "AbortError") {
        console.error("Onboarding session could not be started", error);
        view?.setError("Не удалось начать маршрут. Попробуйте ещё раз.");
      }
    } finally {
      accepting = false;
      sessionController = null;
      view?.setBusy(false);
    }
  }

  async function skipOnboarding(reason: string): Promise<void> {
    if (!progress) return;
    const state = progress;
    progress = null;
    removeProgress();
    if (routeTimer !== null) window.clearTimeout(routeTimer);
    routeTimer = null;
    clearTarget();
    await recordEvent(state, "onboarding_skipped", { reason });
  }

  async function completeCurrentStep(): Promise<void> {
    const state = progress;
    const step = currentStep(state);
    if (!state || !step || !state.accepted || completingStepId === step.id) return;

    completingStepId = step.id;
    if (routeTimer !== null) window.clearTimeout(routeTimer);
    routeTimer = null;
    clearTarget();

    const completion = recordEvent(state, "step_completed", {}, step.id);
    if (state.stepIndex === state.scenario.steps.length - 1) {
      const finished = recordEvent(state, "onboarding_completed");
      progress = null;
      removeProgress();
      await finished;
      if (!state.isTest) onUserStateChange?.({ userId: state.userId, onboarded: true });
    } else {
      progress = { ...state, stepIndex: state.stepIndex + 1 };
      writeProgress();
      renderForCurrentPage();
      await completion;
    }
    completingStepId = null;
  }

  async function resolveScenarios(
    userId: string,
    testToken: string | undefined,
    signal: AbortSignal,
  ) {
    const response = await post(
      "/scenarios/resolve",
      { page: window.location.pathname, user_id: userId },
      signal,
      false,
      testToken,
    );
    return parseRuntimeScenarioResolveResponse(response);
  }

  async function createSession(
    scenarioId: string,
    userId: string,
    signal: AbortSignal,
  ): Promise<string> {
    const response = await post(
      "/sessions",
      { scenario_id: scenarioId, user_id: userId },
      signal,
    );
    if (!isRecord(response) || typeof response.id !== "string") {
      throw new TypeError("Onboarding session response is invalid");
    }
    return response.id;
  }

  async function runStart(options: StartOptions, runGeneration: number): Promise<void> {
    requireBrowser();
    const userId = assertOption(options.userId, "userId");
    const testToken = options.testToken === undefined
      ? undefined
      : assertOption(options.testToken, "testToken");
    const usesTestToken = testToken !== undefined;

    if (
      progress &&
      (progress.userId !== userId || (!progress.isTest && progress.usesTestToken !== usesTestToken))
    ) {
      throw new Error("Call destroy() before changing onboarding user or mode");
    }

    ensureView();
    const stored = progress ?? readProgress();
    if (
      stored &&
      stored.userId === userId &&
      (stored.isTest || stored.usesTestToken === usesTestToken)
    ) {
      progress = stored;
      renderForCurrentPage();
      return;
    }
    if (stored) removeProgress();

    const controller = new AbortController();
    startController = controller;
    const response = await resolveScenarios(userId, testToken, controller.signal);
    if (runGeneration !== generation) return;

    const createProgress = (scenario: RuntimeScenario): Progress => ({
      version: 2,
      projectKey,
      userId,
      usesTestToken,
      isTest: response.is_test,
      sessionId: null,
      scenario,
      stepIndex: 0,
      accepted: false,
      shownStepIds: [],
    });

    if (response.scenarios.length === 0) {
      clearTarget();
      return;
    }

    if (response.scenarios.length === 1) {
      progress = createProgress(response.scenarios[0]!);
      writeProgress();
      renderForCurrentPage();
      return;
    }

    clearTarget();
    const dismiss = () => {
      clearTarget();
    };
    pendingDismiss = dismiss;
    ensureView().showScenarioSelector({
      scenarios: response.scenarios.map(({ id, name, description }) => ({ id, name, description })),
      onSelect: (scenarioId) => {
        const scenario = response.scenarios.find(({ id }) => id === scenarioId);
        if (!scenario) return;
        progress = createProgress(scenario);
        writeProgress();
        void acceptInvitation();
      },
      onDecline: dismiss,
    });
  }

  async function start(options: StartOptions): Promise<void> {
    if (startPromise) return startPromise;
    const runGeneration = generation;
    startPromise = runStart(options, runGeneration);
    try {
      await startPromise;
    } finally {
      startPromise = null;
      startController = null;
    }
  }

  function destroy(): void {
    generation += 1;
    startController?.abort();
    startController = null;
    sessionController?.abort();
    sessionController = null;
    if (routeTimer !== null) window.clearTimeout(routeTimer);
    routeTimer = null;
    if (positionFrame !== null) window.cancelAnimationFrame(positionFrame);
    positionFrame = null;
    clearTarget();

    if (listenersAttached) {
      document.removeEventListener("keydown", onEscape);
      window.removeEventListener("resize", onViewportChange);
      window.removeEventListener("scroll", onViewportChange, true);
      listenersAttached = false;
    }
    view?.remove();
    view = null;
    progress = null;
    completingStepId = null;
    accepting = false;
  }

  return { start, completeCurrentStep, destroy };
}
