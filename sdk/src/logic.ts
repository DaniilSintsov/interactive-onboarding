import type {
  FrontendStepData,
  RuntimeScenario,
  RuntimeScenarioResolveResponse,
  RuntimeStep,
} from "./types.js";

type UnknownRecord = Record<string, unknown>;

function record(value: unknown, path: string): UnknownRecord {
  if (typeof value !== "object" || value === null || Array.isArray(value)) {
    throw new TypeError(`${path} must be an object`);
  }
  return value as UnknownRecord;
}

function string(value: unknown, path: string): string {
  if (typeof value !== "string") {
    throw new TypeError(`${path} must be a string`);
  }
  return value;
}

function frontendData(value: unknown, path: string): FrontendStepData {
  const data = record(value, path);
  const pagePath = string(data.page_path, `${path}.page_path`);
  const advance = record(data.advance, `${path}.advance`);

  if (advance.mode === "manual") {
    return { page_path: pagePath, advance: { mode: "manual" } };
  }
  if (
    advance.mode === "target_event" &&
    (advance.event === "click" || advance.event === "change")
  ) {
    return {
      page_path: pagePath,
      advance: { mode: "target_event", event: advance.event },
    };
  }
  throw new TypeError(`${path}.advance is invalid`);
}

function runtimeStep(value: unknown, index: number): RuntimeStep {
  const path = `scenario.steps[${index}]`;
  const step = record(value, path);
  const element = record(step.element, `${path}.element`);

  if (!Number.isInteger(step.step_num) || Number(step.step_num) < 1) {
    throw new TypeError(`${path}.step_num must be a positive integer`);
  }

  return {
    id: string(step.id, `${path}.id`),
    step_num: Number(step.step_num),
    title: string(step.title, `${path}.title`),
    description: string(step.description, `${path}.description`),
    frontend_data: frontendData(step.frontend_data, `${path}.frontend_data`),
    element: {
      id: string(element.id, `${path}.element.id`),
      key: string(element.key, `${path}.element.key`),
      label: string(element.label, `${path}.element.label`),
      description: string(element.description, `${path}.element.description`),
    },
  };
}

export function parseRuntimeScenario(value: unknown): RuntimeScenario {
  const scenario = record(value, "scenario");
  if (!Array.isArray(scenario.steps) || scenario.steps.length === 0) {
    throw new TypeError("scenario.steps must be a non-empty array");
  }

  return {
    id: string(scenario.id, "scenario.id"),
    name: string(scenario.name, "scenario.name"),
    description: string(scenario.description, "scenario.description"),
    page_pattern: string(scenario.page_pattern, "scenario.page_pattern"),
    steps: scenario.steps.map(runtimeStep).sort((a, b) => a.step_num - b.step_num),
  };
}

export function parseRuntimeScenarioResolveResponse(
  value: unknown,
): RuntimeScenarioResolveResponse {
  const response = record(value, "response");
  if (typeof response.is_test !== "boolean") {
    throw new TypeError("response.is_test must be a boolean");
  }
  if (!Array.isArray(response.scenarios)) {
    throw new TypeError("response.scenarios must be an array");
  }

  return {
    is_test: response.is_test,
    scenarios: response.scenarios.map(parseRuntimeScenario),
  };
}

export function hasMeaningfulValue(target: EventTarget | null): boolean {
  if (typeof target !== "object" || target === null) return false;

  const field = target as {
    type?: unknown;
    files?: { length: number } | null;
    value?: unknown;
  };
  if (field.type === "file") return Boolean(field.files?.length);
  return field.value !== undefined && String(field.value).trim().length > 0;
}
