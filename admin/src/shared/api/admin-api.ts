import type {
  DetailedScenarioAnalytics,
  Element,
  ElementInput,
  ProjectAnalytics,
  ProjectList,
  ProjectWithElements,
  Scenario,
  ScenarioInput,
  ScenarioList,
  ScenarioTestToken,
  ScenarioWithSteps,
  Step,
  StepInput,
} from './types';

const BFF = '/api/backend';

export class ApiError extends Error {
  constructor(
    message: string,
    public readonly status: number,
  ) {
    super(message);
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${BFF}${path}`, {
    ...init,
    headers: {
      ...(init?.body ? { 'content-type': 'application/json' } : {}),
      ...init?.headers,
    },
  });
  if (response.status === 401 && typeof window !== 'undefined') window.location.replace('/login');
  if (!response.ok) {
    const body = (await response.json().catch(() => null)) as { message?: string } | null;
    throw new ApiError(body?.message || `API вернул ${response.status}`, response.status);
  }
  if (response.status === 204) return undefined as T;
  return response.json() as Promise<T>;
}

const json = (value: unknown) => JSON.stringify(value);

export const adminApi = {
  listProjects: () => request<ProjectList>('/projects?limit=100'),
  getProject: (projectId: string) => request<ProjectWithElements>(`/projects/${projectId}`),
  createProject: (name: string) =>
    request<ProjectWithElements>('/projects', {
      method: 'POST',
      body: json({ name, elements: [] }),
    }),

  listElements: (projectId: string) => request<Element[]>(`/projects/${projectId}/elements`),
  createElement: (projectId: string, input: ElementInput) =>
    request<Element>(`/projects/${projectId}/elements`, { method: 'POST', body: json(input) }),
  updateElement: (projectId: string, elementId: string, input: ElementInput) =>
    request<Element>(`/projects/${projectId}/elements/${elementId}`, {
      method: 'PATCH',
      body: json(input),
    }),
  deleteElement: (projectId: string, elementId: string) =>
    request<void>(`/projects/${projectId}/elements/${elementId}`, { method: 'DELETE' }),

  listScenarios: (projectId: string) =>
    request<ScenarioList>(`/projects/${projectId}/scenarios?limit=100`),
  getScenario: (scenarioId: string) => request<ScenarioWithSteps>(`/scenarios/${scenarioId}`),
  createScenario: (projectId: string, input: ScenarioInput) =>
    request<Scenario>(`/projects/${projectId}/scenarios`, { method: 'POST', body: json(input) }),
  updateScenario: (scenarioId: string, input: ScenarioInput) =>
    request<Scenario>(`/scenarios/${scenarioId}`, { method: 'PATCH', body: json(input) }),
  deleteScenario: (scenarioId: string) =>
    request<void>(`/scenarios/${scenarioId}`, { method: 'DELETE' }),
  transitionScenario: (scenarioId: string, action: 'publish' | 'enable' | 'disable') =>
    request<Scenario>(`/scenarios/${scenarioId}/${action}`, { method: 'POST' }),
  createTestToken: (scenarioId: string) =>
    request<ScenarioTestToken>(`/scenarios/${scenarioId}/test-tokens`, { method: 'POST' }),

  createStep: (scenarioId: string, input: StepInput) =>
    request<Step>(`/scenarios/${scenarioId}/steps`, { method: 'POST', body: json(input) }),
  updateStep: (scenarioId: string, stepId: string, input: StepInput) =>
    request<Step>(`/scenarios/${scenarioId}/steps/${stepId}`, {
      method: 'PATCH',
      body: json(input),
    }),
  deleteStep: (scenarioId: string, stepId: string) =>
    request<void>(`/scenarios/${scenarioId}/steps/${stepId}`, { method: 'DELETE' }),
  reorderSteps: (scenarioId: string, orderedStepIds: string[]) =>
    request<Step[]>(`/scenarios/${scenarioId}/steps/order`, {
      method: 'PUT',
      body: json({ ordered_step_ids: orderedStepIds }),
    }),

  getProjectAnalytics: (projectId: string, period = '') =>
    request<ProjectAnalytics>(`/projects/${projectId}/analytics/total${period}`),
  getScenarioAnalytics: (scenarioId: string, period = '') =>
    request<DetailedScenarioAnalytics>(`/scenarios/${scenarioId}/analytics/detailed${period}`),
};

export async function downloadReport(scenarioId: string, period = ''): Promise<void> {
  const response = await fetch(`${BFF}/scenarios/${scenarioId}/report/pdf${period}`);
  if (!response.ok) {
    const body = (await response.json().catch(() => null)) as { message?: string } | null;
    throw new ApiError(body?.message || 'Не удалось сформировать PDF', response.status);
  }
  const link = document.createElement('a');
  link.href = URL.createObjectURL(await response.blob());
  link.download = `scenario-${scenarioId}-report.pdf`;
  link.click();
  URL.revokeObjectURL(link.href);
}
