import { describe, expect, it } from 'vitest';
import { getStepCatalogWarnings, initialStepValues, selectPagePath, toStepInput } from './step-form';

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

describe('selectPagePath', () => {
  it('keeps a matching route from catalog', () => {
    expect(selectPagePath('/add-item/details', ['/add-item/details', '/add-item/title'])).toBe('/add-item/details');
  });

  it('falls back to first catalog route when current path mismatches', () => {
    expect(selectPagePath('/legacy', ['/add-item/details', '/add-item/title'])).toBe('/add-item/details');
  });

  it('keeps manual value when catalog has no routes', () => {
    expect(selectPagePath('/legacy', [])).toBe('/legacy');
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

describe('getStepCatalogWarnings', () => {
  const step = {
    id: 'step-id',
    scenario_id: 'scenario-id',
    element_id: 'element-id',
    step_num: 2,
    title: 'Шаг',
    description: 'Описание',
    frontend_data: { page_path: '/legacy', advance: { mode: 'manual' as const } },
    created_at: '',
    updated_at: '',
  };

  it('warns when element is no longer in catalog and has no routes', () => {
    expect(getStepCatalogWarnings({ step, isAvailable: false, pagePaths: [] })).toEqual([
      'Элемент больше не найден в CI-каталоге. Для новых шагов он недоступен.',
    ]);
  });

  it('warns when step route mismatches catalog', () => {
    expect(getStepCatalogWarnings({ step, isAvailable: true, pagePaths: ['/add-item/details'] })).toEqual([
      'Маршрут шага /legacy не совпадает с каталогом. В форме выбран /add-item/details.',
    ]);
  });
});
