import { describe, expect, it } from 'vitest';
import {
  firstElementIdForPage,
  initialStepValues,
  isElementMissingFromPage,
  toStepInput,
} from './step-form';

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

describe('firstElementIdForPage', () => {
  it('selects the first element of the new page', () => {
    expect(
      firstElementIdForPage(
        [
          { id: 'other-page', page: '/other' },
          { id: 'first-match', page: ' /add-item/category ' },
          { id: 'second-match', page: '/add-item/category' },
        ],
        '/add-item/category',
      ),
    ).toBe('first-match');
  });

  it('clears the selection when the page has no elements', () => {
    expect(firstElementIdForPage([{ id: 'other-page', page: '/other' }], '/empty')).toBeUndefined();
  });

  it('treats the newly selected element as available', () => {
    const pageElements = [{ id: 'first-match', page: '/add-item/category' }];
    const selectedElementId = firstElementIdForPage(pageElements, '/add-item/category');

    expect(isElementMissingFromPage(pageElements, selectedElementId)).toBe(false);
    expect(isElementMissingFromPage(pageElements, 'saved-on-other-page')).toBe(true);
  });
});
