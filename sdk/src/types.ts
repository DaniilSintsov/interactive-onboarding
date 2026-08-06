export type FrontendStepData = {
  page_path: string;
  advance:
    | { mode: "target_event"; event: "click" | "change" }
    | { mode: "manual" };
};

export type RuntimeElement = {
  id: string;
  key: string;
  label: string;
  description: string;
};

export type RuntimeStep = {
  id: string;
  step_num: number;
  title: string;
  description: string;
  frontend_data: FrontendStepData;
  element: RuntimeElement;
};

export type RuntimeScenario = {
  id: string;
  name: string;
  description: string;
  page_pattern: string;
  steps: RuntimeStep[];
};

export type EventType =
  | "step_shown"
  | "step_completed"
  | "step_skipped"
  | "onboarding_completed"
  | "onboarding_skipped";

export type CreateOnboardingOptions = {
  projectKey: string;
  runtimeUrl: string;
};

export type StartOptions = {
  userId: string;
  preview?: boolean;
};

export type Onboarding = {
  start(options: StartOptions): Promise<void>;
  completeCurrentStep(): Promise<void>;
  destroy(): void;
};
