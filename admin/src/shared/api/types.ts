export type Project = {
  id: string;
  name: string;
  project_key: string;
  created_at: string;
  updated_at: string;
};

export type ProjectList = {
  items: Project[];
  total: number;
  limit: number;
  offset: number;
};

export type Element = {
  id: string;
  project_id: string;
  key: string;
  label: string;
  description: string;
  created_at: string;
  updated_at: string;
};

export type ProjectWithElements = Project & { elements: Element[] };

export type ScenarioStatus = 'in_development' | 'enabled' | 'disabled';

export type ScenarioSummary = {
  id: string;
  project_id: string;
  name: string;
  page_pattern: string;
  status: ScenarioStatus;
  steps_count: number;
  published_at: string | null;
  created_at: string;
  updated_at: string;
};

export type ScenarioList = {
  items: ScenarioSummary[];
  total: number;
  limit: number;
  offset: number;
};

export type Advance =
  | { mode: 'target_event'; event: 'click' | 'change' }
  | { mode: 'manual' };

export type FrontendStepData = {
  page_path: string;
  advance: Advance;
};

export type Step = {
  id: string;
  scenario_id: string;
  element_id: string;
  step_num: number;
  title: string;
  description: string;
  frontend_data: FrontendStepData;
  created_at: string;
  updated_at: string;
};

export type Scenario = {
  id: string;
  project_id: string;
  name: string;
  description: string;
  page_pattern: string;
  status: ScenarioStatus;
  published_at: string | null;
  created_at: string;
  updated_at: string;
};

export type ScenarioWithSteps = Scenario & { steps: Step[] };

export type ProjectAnalytics = {
  project_id: string;
  total_scenarios: number;
  enabled_scenarios: number;
  sessions_started: number;
  sessions_completed: number;
  sessions_skipped: number;
  completion_rate: number;
  skip_rate: number;
};

export type StepAnalytics = {
  step_id: string;
  position: number;
  title: string;
  shown: number;
  completed: number;
  skipped: number;
  completion_rate: number;
  skip_rate: number;
  drop_off_rate: number;
};

export type DetailedScenarioAnalytics = {
  scenario_id: string;
  started: number;
  completed: number;
  skipped: number;
  completion_rate: number;
  skip_rate: number;
  average_completion_time_seconds: number;
  steps: StepAnalytics[];
};

export type ElementInput = { key: string; label: string; description?: string };
export type ScenarioInput = { name: string; description?: string; page_pattern: string };
export type StepInput = {
  element_id: string;
  title: string;
  description: string;
  frontend_data: FrontendStepData;
};
