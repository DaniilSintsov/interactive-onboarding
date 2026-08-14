import { describe, expect, it } from 'vitest';
import { initialStepValues, toStepInput } from './step-form';

const base = {
  element_id: 'element-id',
  title: 'Добавьте фото',
  description: 'Загрузите хотя бы один снимок',
  page_path: '/add-item/details',
} as const;

describe('toStepInput', () => {
  it('builds target-event frontend_data', () => {
    expect(toStepInput({ ...base, mode: 'target_event', event: 'change' })).toEqual({
      element_id: 'element-id',
      title: 'Добавьте фото',
      description: 'Загрузите хотя бы один снимок',
      frontend_data: {
        page_path: '/add-item/details',
        advance: { mode: 'target_event', event: 'change' },
      },
    });
  });

  it('drops a stale event in manual mode', () => {
    expect(toStepInput({ ...base, mode: 'manual', event: 'click' }).frontend_data.advance).toEqual({
      mode: 'manual',
    });
  });
});

describe('initialStepValues', () => {
  it('keeps an empty legacy page_path for editing', () => {
    expect(
      initialStepValues({
        id: 'step-id',
        scenario_id: 'scenario-id',
        element_id: 'element-id',
        step_num: 1,
        title: 'Шаг',
        description: 'Описание',
        frontend_data: { page_path: '', advance: { mode: 'manual' } },
        created_at: '',
        updated_at: '',
      }).page_path,
    ).toBe('');
  });
});
