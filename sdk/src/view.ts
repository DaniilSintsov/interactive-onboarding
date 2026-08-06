type InvitationContent = {
  name: string;
  description: string;
  onAccept: () => void;
  onDecline: () => void;
};

type HintContent = {
  title: string;
  description: string;
  position: number;
  total: number;
  onNext: () => void;
  onClose: () => void;
};

export type OnboardingView = HTMLElement & {
  showInvitation(content: InvitationContent): void;
  showHint(content: HintContent): void;
  positionNear(target: Element): void;
  hideCard(): void;
};

const TAG_NAME = "interactive-onboarding-card";

const TEMPLATE = `
  <style>
    :host {
      --ink: #111;
      --blue: #00aaff;
      --sun: #ffcc00;
      position: fixed;
      inset: 0 auto auto 0;
      width: 0;
      height: 0;
      z-index: 2147483000;
      color: var(--ink);
      font-family: "Avenir Next", "Helvetica Neue", sans-serif;
    }

    [hidden] { display: none !important; }

    .card {
      box-sizing: border-box;
      position: fixed;
      width: min(326px, calc(100vw - 24px));
      padding: 20px;
      overflow: hidden;
      border: 1px solid #ddd;
      border-radius: 14px;
      background: #fff;
      box-shadow: 0 14px 40px rgba(0, 0, 0, 0.16);
      pointer-events: auto;
      animation: arrive 180ms ease-out both;
    }

    .card::before {
      content: "";
      position: absolute;
      inset: 0 0 auto;
      height: 5px;
      background: var(--blue);
    }

    :host([data-mode="invitation"]) .card {
      right: 24px;
      bottom: 24px;
    }

    .eyebrow {
      margin: 2px 34px 8px 0;
      color: var(--blue);
      font-size: 11px;
      font-weight: 800;
      letter-spacing: 0.04em;
      line-height: 1.3;
      text-transform: uppercase;
    }

    h2 {
      margin: 0 30px 8px 0;
      color: var(--ink);
      font-family: "Avenir Next", "Helvetica Neue", sans-serif;
      font-size: 21px;
      font-weight: 700;
      letter-spacing: -0.02em;
      line-height: 1.2;
    }

    p {
      margin: 0;
      color: #666;
      font-size: 14px;
      line-height: 1.55;
    }

    .close {
      position: absolute;
      top: 14px;
      right: 12px;
      display: grid;
      width: 32px;
      height: 32px;
      padding: 0;
      place-items: center;
      border: 0;
      border-radius: 50%;
      background: transparent;
      color: #666;
      cursor: pointer;
      font: 25px/1 Georgia, serif;
    }

    .close:hover { background: #f2f2f2; color: var(--ink); }

    .close:focus-visible,
    button:focus-visible {
      outline: 3px solid var(--sun);
      outline-offset: 2px;
    }

    .actions {
      display: flex;
      gap: 10px;
      margin-top: 18px;
    }

    .actions button {
      min-height: 42px;
      padding: 0 16px;
      border: 1px solid transparent;
      border-radius: 8px;
      cursor: pointer;
      font: 700 14px/1 "Avenir Next", "Helvetica Neue", sans-serif;
    }

    .accept { background: var(--blue); color: white; }
    .accept:hover { background: #008aed; }
    .decline { border-color: #ddd !important; background: #fff; color: #444; }
    .decline:hover { background: #f2f2f2; }

    @keyframes arrive {
      from { opacity: 0; transform: translateY(8px) scale(0.985); }
      to { opacity: 1; }
    }

    @media (max-width: 520px) {
      :host([data-mode="invitation"]) .card { right: 12px; bottom: 12px; }
    }

    @media (prefers-reduced-motion: reduce) {
      .card { animation: none; }
    }
  </style>
  <section class="card" role="dialog" aria-modal="false" aria-live="polite" aria-labelledby="onboarding-title" hidden>
    <button class="close" type="button" aria-label="Закрыть подсказки" title="Закрыть подсказки"><span aria-hidden="true">×</span></button>
    <div class="eyebrow"></div>
    <h2 id="onboarding-title"></h2>
    <p></p>
    <div class="actions">
      <button class="accept" type="button">Начать</button>
      <button class="decline" type="button">Не сейчас</button>
    </div>
  </section>
`;

function registerElement(): void {
  if (customElements.get(TAG_NAME)) return;

  class OnboardingCardElement extends HTMLElement {
    private readonly card: HTMLElement;
    private readonly eyebrow: HTMLElement;
    private readonly cardTitle: HTMLElement;
    private readonly description: HTMLElement;
    private readonly actions: HTMLElement;
    private readonly accept: HTMLButtonElement;
    private readonly decline: HTMLButtonElement;
    private readonly close: HTMLButtonElement;
    private onClose: () => void = () => undefined;

    constructor() {
      super();
      const shadow = this.attachShadow({ mode: "open" });
      shadow.innerHTML = TEMPLATE;

      this.card = this.required(shadow, ".card");
      this.eyebrow = this.required(shadow, ".eyebrow");
      this.cardTitle = this.required(shadow, "h2");
      this.description = this.required(shadow, "p");
      this.actions = this.required(shadow, ".actions");
      this.accept = this.required(shadow, ".accept");
      this.decline = this.required(shadow, ".decline");
      this.close = this.required(shadow, ".close");
      this.close.addEventListener("click", () => this.onClose());
    }

    showInvitation(content: InvitationContent): void {
      this.dataset.mode = "invitation";
      this.card.style.removeProperty("transform");
      this.eyebrow.textContent = "Короткий маршрут";
      this.cardTitle.textContent = content.name;
      this.description.textContent = content.description || "Покажем основные шаги прямо на странице.";
      this.actions.hidden = false;
      this.accept.textContent = "Начать";
      this.accept.onclick = content.onAccept;
      this.decline.hidden = false;
      this.decline.onclick = content.onDecline;
      this.onClose = content.onDecline;
      this.card.hidden = false;
    }

    showHint(content: HintContent): void {
      this.dataset.mode = "hint";
      this.eyebrow.textContent = `Шаг ${content.position} из ${content.total}`;
      this.cardTitle.textContent = content.title;
      this.description.textContent = content.description;
      this.actions.hidden = false;
      this.accept.textContent = "Далее";
      this.accept.onclick = content.onNext;
      this.decline.hidden = true;
      this.onClose = content.onClose;
      this.card.hidden = false;
    }

    positionNear(target: Element): void {
      const targetRect = target.getBoundingClientRect();
      const cardRect = this.card.getBoundingClientRect();
      const gap = 14;
      const padding = 12;
      const left = Math.min(
        Math.max(padding, targetRect.left),
        Math.max(padding, window.innerWidth - cardRect.width - padding),
      );
      const below = targetRect.bottom + gap;
      const top = below + cardRect.height <= window.innerHeight - padding
        ? below
        : Math.max(padding, targetRect.top - cardRect.height - gap);

      this.card.style.transform = `translate(${Math.round(left)}px, ${Math.round(top)}px)`;
    }

    hideCard(): void {
      this.card.hidden = true;
      this.accept.onclick = null;
      this.decline.onclick = null;
      this.onClose = () => undefined;
    }

    private required<T extends Element>(root: ShadowRoot, selector: string): T {
      const element = root.querySelector<T>(selector);
      if (!element) throw new Error(`Onboarding view is missing ${selector}`);
      return element;
    }
  }

  customElements.define(TAG_NAME, OnboardingCardElement);
}

export function createOnboardingView(): OnboardingView {
  registerElement();
  return document.createElement(TAG_NAME) as OnboardingView;
}
