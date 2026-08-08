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

export function initialStepValues(step: Step | null): StepFormValues {
  const advance = step?.frontend_data.advance;
  return {
    element_id: step?.element_id || '',
    title: step?.title || '',
    description: step?.description || '',
    page_path: step?.frontend_data.page_path || '/',
    mode: advance?.mode || 'target_event',
    event: advance?.mode === 'target_event' ? advance.event : 'click',
  };
}
