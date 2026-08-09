import type { Step, StepInput } from '@/shared/api/types';

export type StepFormValues = {
  element_id: string;
  title: string;
  description: string;
  page_path: string;
  mode: 'target_event' | 'manual';
  event?: 'click' | 'change';
};

export function toStepInput(values: StepFormValues): StepInput {
  return {
    element_id: values.element_id,
    title: values.title,
    description: values.description,
    frontend_data: {
      page_path: values.page_path,
      advance:
        values.mode === 'manual'
          ? { mode: 'manual' }
          : { mode: 'target_event', event: values.event || 'click' },
    },
  };
}

export function normalizePagePaths(pagePaths: readonly string[] | undefined): string[] {
  return [...new Set((pagePaths ?? []).map((pagePath) => pagePath.trim()).filter(Boolean))];
}

export function selectPagePath(currentPagePath: string, pagePaths: readonly string[] | undefined): string {
  const normalizedPagePaths = normalizePagePaths(pagePaths);
  const current = currentPagePath.trim();
  if (normalizedPagePaths.length === 0) return current;
  return normalizedPagePaths.includes(current) ? current : normalizedPagePaths[0];
}

export function getStepCatalogWarnings({
  step,
  isAvailable,
  pagePaths,
}: {
  step: Step | null;
  isAvailable: boolean;
  pagePaths: readonly string[] | undefined;
}): string[] {
  if (!step) return [];

  const warnings: string[] = [];
  const normalizedPagePaths = normalizePagePaths(pagePaths);

  if (!isAvailable) {
    warnings.push('Элемент больше не найден в CI-каталоге. Для новых шагов он недоступен.');
  }

  if (normalizedPagePaths.length > 0 && !normalizedPagePaths.includes(step.frontend_data.page_path.trim())) {
    warnings.push(
      `Маршрут шага ${step.frontend_data.page_path || 'пустой'} не совпадает с каталогом. В форме выбран ${normalizedPagePaths[0]}.`,
    );
  }

  return warnings;
}

export function initialStepValues(step: Step | null): StepFormValues {
  const advance = step?.frontend_data.advance;
  return {
    element_id: step?.element_id || '',
    title: step?.title || '',
    description: step?.description || '',
    page_path: step?.frontend_data.page_path || '',
    mode: advance?.mode || 'target_event',
    event: advance?.mode === 'target_event' ? advance.event : 'click',
  };
}
