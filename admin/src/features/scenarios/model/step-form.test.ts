import { describe, expect, it } from 'vitest';
import { toStepInput } from './step-form';

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
