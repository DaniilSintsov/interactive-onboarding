'use client';

import { useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  Alert,
  App,
  Button,
  Card,
  Form,
  Input,
  Modal,
  Popconfirm,
  Select,
  Skeleton,
  Space,
} from 'antd';
import { adminApi } from '@/shared/api/admin-api';
import type { Element, ScenarioInput, Step, StepInput } from '@/shared/api/types';
import {
  firstElementIdForPage,
  initialStepValues,
  isElementMissingFromPage,
  toStepInput,
  type StepFormValues,
} from '@/features/scenarios/model/step-form';
import { buildScenarioPreviewUrl, isScenarioPathname } from '@/features/scenarios/model/preview';
import { ScenarioStatus } from './scenario-status';

const scenarioPathnameRules = [
  { required: true, whitespace: true },
  { max: 2048 },
  {
    validator: async (_: unknown, value?: string) => {
      if (!value) return;
      if (!isScenarioPathname(value)) {
        throw new Error('Укажите путь вида /catalog без домена, параметров после ? и #.');
      }
    },
  },
];

function formatStepAdvance(step: Step): string {
  if (step.frontend_data.advance.mode === 'manual') return 'Завершение: по команде из приложения';
  return step.frontend_data.advance.event === 'click'
    ? 'Завершение: после клика по элементу'
    : 'Завершение: после ввода или выбора значения';
}

function StepModal({
  projectId,
  step,
  allElements,
  pending,
  onCancel,
  onSave,
}: {
  projectId: string;
  step: Step | null;
  allElements: Element[];
  pending: boolean;
  onCancel: () => void;
  onSave: (input: StepInput) => void;
}) {
  const [form] = Form.useForm<StepFormValues>();
  const mode = Form.useWatch('mode', form);
  const event = Form.useWatch('event', form);
  const elementId = Form.useWatch('element_id', form) ?? '';
  const pagePath = (Form.useWatch('page_path', form) ?? '').trim();
  const initialPagePath = step?.frontend_data.page_path ?? '';
  const pages = useQuery({
    queryKey: ['pages', projectId],
    queryFn: () => adminApi.listPages(projectId),
  });
  const pageElements = useQuery({
    queryKey: ['elements', projectId, pagePath],
    queryFn: () => adminApi.listElements(projectId, pagePath),
    enabled: Boolean(pagePath),
  });
  const projectPagePaths = useMemo(
    () => [...new Set((pages.data ?? []).map((page) => page.trim()).filter(Boolean))],
    [pages.data],
  );
  const pageOptions = useMemo(() => {
    const nextPagePaths = [...projectPagePaths];
    if (initialPagePath.trim() && !nextPagePaths.includes(initialPagePath.trim())) nextPagePaths.unshift(initialPagePath.trim());
    return nextPagePaths.map((pathname) => ({
      value: pathname,
      label: pathname === initialPagePath.trim() && !projectPagePaths.includes(pathname) ? `${pathname} · вне текущего списка` : pathname,
    }));
  }, [initialPagePath, projectPagePaths]);
  const selectedElement = useMemo(
    () =>
      (pageElements.data ?? []).find((element) => element.id === elementId) ??
      allElements.find((element) => element.id === elementId) ??
      null,
    [allElements, elementId, pageElements.data],
  );
  const elementOptions = useMemo(() => {
    const options = (pageElements.data ?? []).map((element) => ({
      value: element.id,
      label: `${element.label} · ${element.key}`,
    }));
    if (selectedElement && !options.some((option) => option.value === selectedElement.id)) {
      options.unshift({
        value: selectedElement.id,
        label: `${selectedElement.label} · ${selectedElement.key} · вне выбранной страницы`,
      });
    }
    return options;
  }, [pageElements.data, selectedElement]);
  const pageMissing =
    pages.isSuccess && Boolean(initialPagePath.trim()) && !projectPagePaths.includes(initialPagePath.trim());
  const elementMissing =
    Boolean(pagePath) &&
    pageElements.isSuccess &&
    isElementMissingFromPage(pageElements.data, elementId);

  return (
    <Modal
      title={step ? `Шаг ${step.step_num}: ${step.title}` : 'Новый шаг'}
      open
      width={680}
      onCancel={onCancel}
      onOk={() => form.submit()}
      okText="Сохранить шаг"
      cancelText="Отмена"
      confirmLoading={pending}
      destroyOnHidden
    >
      <Form
        form={form}
        layout="vertical"
        requiredMark={false}
        initialValues={initialStepValues(step)}
        onFinish={(values) => onSave(toStepInput(values))}
      >
        {pages.isError ? (
          <Alert
            style={{ marginBottom: 16 }}
            type="error"
            showIcon
            message="Не удалось загрузить страницы проекта"
            description={pages.error.message}
          />
        ) : null}
        {pageMissing ? (
          <Alert
            style={{ marginBottom: 16 }}
            type="warning"
            showIcon
            message="Маршрут шага больше не найден в проекте"
            description={`Сохранённый маршрут ${initialPagePath} оставлен в форме, но отсутствует в текущем списке страниц.`}
          />
        ) : null}
        {elementMissing ? (
          <Alert
            style={{ marginBottom: 16 }}
            type="warning"
            showIcon
            message="Элемент шага больше не найден на выбранной странице"
            description="Выберите доступный элемент или смените страницу."
          />
        ) : null}
        <Form.Item
          label="Страница шага"
          extra="Подсказка появится на этой странице рядом с выбранным элементом. При смене страницы выбирается первый доступный элемент."
          name="page_path"
          rules={[
            { required: true, whitespace: true, message: 'Выберите страницу' },
            { pattern: /^\//, message: 'Маршрут начинается с /' },
            { max: 2048 },
          ]}
        >
          <Select
            showSearch
            optionFilterProp="label"
            loading={pages.isPending}
            placeholder="Выберите страницу проекта"
            options={pageOptions}
            notFoundContent={pages.isPending ? 'Загрузка...' : 'Страницы не найдены'}
            onChange={(nextPagePath) =>
              form.setFieldValue('element_id', firstElementIdForPage(allElements, nextPagePath))
            }
          />
        </Form.Item>
        {pageElements.isError ? (
          <Alert
            style={{ marginBottom: 16 }}
            type="error"
            showIcon
            message="Не удалось загрузить элементы страницы"
            description={pageElements.error.message}
          />
        ) : null}
        <Form.Item
          label="Целевой элемент"
          extra="SDK ищет на странице элемент с этим data-onboarding-id и показывает подсказку рядом с ним."
          name="element_id"
          rules={[{ required: true, message: 'Выберите элемент' }]}
        >
          <Select
            showSearch
            optionFilterProp="label"
            loading={pageElements.isPending}
            disabled={!pagePath || pageElements.isError}
            placeholder={pagePath ? 'Выберите data-onboarding-id' : 'Сначала выберите страницу'}
            options={elementOptions}
            notFoundContent={pagePath ? 'На странице нет элементов' : 'Сначала выберите страницу'}
          />
        </Form.Item>
        <div className="toolbar">
          <Form.Item
            label="Как завершается шаг"
            extra={
              mode === 'manual'
                ? 'Шаг завершится только когда приложение вызовет completeCurrentStep().'
                : 'Шаг завершится автоматически после выбранного действия пользователя.'
            }
            name="mode"
            style={{ flex: 1, minWidth: 220 }}
          >
            <Select
              options={[
                { value: 'target_event', label: 'После действия пользователя' },
                { value: 'manual', label: 'По команде из приложения' },
              ]}
            />
          </Form.Item>
          {mode !== 'manual' ? (
            <Form.Item
              label="Действие, завершающее шаг"
              extra={
                event === 'change'
                  ? 'Ждём непустое значение или выбранный файл: пустое изменение шаг не завершает.'
                  : 'Первый клик по целевому элементу завершит шаг.'
              }
              name="event"
              style={{ flex: 1, minWidth: 180 }}
            >
              <Select
                options={[
                  { value: 'click', label: 'Клик по элементу' },
                  { value: 'change', label: 'Ввод или выбор значения' },
                ]}
              />
            </Form.Item>
          ) : null}
        </div>
        <Form.Item label="Заголовок подсказки" name="title" rules={[{ required: true, whitespace: true }, { max: 255 }]}>
          <Input placeholder="Добавьте фотографию" />
        </Form.Item>
        <Form.Item label="Текст подсказки" name="description" rules={[{ required: true, whitespace: true }, { max: 2000 }]}>
          <Input.TextArea rows={4} placeholder="Покажите товар с лучшей стороны — достаточно одного снимка." />
        </Form.Item>
      </Form>
    </Modal>
  );
}

export function ScenarioEditor({ projectId, scenarioId }: { projectId: string; scenarioId: string }) {
  const [editingStep, setEditingStep] = useState<Step | null | undefined>(undefined);
  const queryClient = useQueryClient();
  const { message } = App.useApp();
  const scenario = useQuery({
    queryKey: ['scenario', scenarioId],
    queryFn: () => adminApi.getScenario(scenarioId),
  });
  const elements = useQuery({
    queryKey: ['elements', projectId],
    queryFn: () => adminApi.listElements(projectId),
  });
  const elementNames = useMemo(
    () =>
      new Map(
        (elements.data ?? []).map((element) => [
          element.id,
          element.label,
        ] as const),
      ),
    [elements.data],
  );
  const availableElementCount = useMemo(
    () => (elements.data ?? []).filter((element) => element.page.trim()).length,
    [elements.data],
  );

  const refresh = () =>
    Promise.all([
      queryClient.invalidateQueries({ queryKey: ['scenario', scenarioId] }),
      queryClient.invalidateQueries({ queryKey: ['scenarios', projectId] }),
    ]);

  const updateGeneral = useMutation({
    mutationFn: (input: ScenarioInput) => adminApi.updateScenario(scenarioId, input),
    onSuccess: () => {
      void refresh();
      message.success('Настройки сохранены');
    },
    onError: (error) => message.error(error.message),
  });
  const transition = useMutation({
    mutationFn: (action: 'publish' | 'enable' | 'disable') =>
      adminApi.transitionScenario(scenarioId, action),
    onSuccess: (_, action) => {
      void refresh();
      message.success(
        action === 'disable' ? 'Сценарий выключен' : action === 'enable' ? 'Сценарий включён' : 'Сценарий опубликован',
      );
    },
    onError: (error) => message.error(error.message),
  });
  const saveStep = useMutation({
    mutationFn: ({ step, input }: { step: Step | null; input: StepInput }) =>
      step
        ? adminApi.updateStep(scenarioId, step.id, input)
        : adminApi.createStep(scenarioId, input),
    onSuccess: () => {
      void refresh();
      setEditingStep(undefined);
      message.success('Шаг сохранён');
    },
    onError: (error) => message.error(error.message),
  });
  const removeStep = useMutation({
    mutationFn: (stepId: string) => adminApi.deleteStep(scenarioId, stepId),
    onSuccess: () => {
      void refresh();
      message.success('Шаг удалён');
    },
    onError: (error) => message.error(error.message),
  });
  const reorder = useMutation({
    mutationFn: (ids: string[]) => adminApi.reorderSteps(scenarioId, ids),
    onSuccess: () => void refresh(),
    onError: (error) => message.error(error.message),
  });
  const preview = useMutation({
    mutationFn: async ({ pagePattern, popup }: { pagePattern: string; popup: Window }) => {
      try {
        const { token } = await adminApi.createTestToken(scenarioId);
        const url = buildScenarioPreviewUrl(pagePattern, token);
        popup.location.replace(url);
      } catch (error) {
        popup.close();
        throw error;
      }
    },
    onError: (error) => message.error(error.message),
  });

  if (scenario.isPending || elements.isPending) return <Skeleton active paragraph={{ rows: 14 }} />;
  if (scenario.isError) return <Alert type="error" showIcon message={scenario.error.message} />;
  if (elements.isError) return <Alert type="error" showIcon message={elements.error.message} />;

  const scenarioData = scenario.data;
  const editable = scenarioData.status !== 'enabled';
  const orderedSteps = [...scenarioData.steps].sort((left, right) => left.step_num - right.step_num);

  function moveStep(index: number, direction: -1 | 1) {
    const target = index + direction;
    if (target < 0 || target >= orderedSteps.length) return;
    const ids = orderedSteps.map((step) => step.id);
    [ids[index], ids[target]] = [ids[target], ids[index]];
    reorder.mutate(ids);
  }

  function openPreview() {
    const popup = window.open('about:blank', '_blank');
    if (!popup) {
      message.error('Разрешите всплывающие окна для проверки сценария');
      return;
    }
    popup.opener = null;
    preview.mutate({ pagePattern: scenarioData.page_pattern, popup });
  }

  return (
    <main>
      <header className="page-heading">
        <div>
          <p className="eyebrow">Редактор сценария</p>
          <h1>{scenarioData.name}</h1>
          <p><ScenarioStatus status={scenarioData.status} /> · {orderedSteps.length} шагов</p>
        </div>
        <Space wrap>
          <Button href={`/projects/${projectId}`}>← К проекту</Button>
          {orderedSteps.length > 0 ? (
            <Button loading={preview.isPending} disabled={preview.isPending} onClick={openPreview}>
              Проверить ↗
            </Button>
          ) : null}
          {scenarioData.status === 'in_development' ? (
            <Button
              type="primary"
              disabled={orderedSteps.length === 0}
              loading={transition.isPending}
              onClick={() => transition.mutate('publish')}
            >
              Опубликовать
            </Button>
          ) : null}
          {scenarioData.status === 'disabled' ? (
            <Button type="primary" loading={transition.isPending} onClick={() => transition.mutate('enable')}>
              Включить
            </Button>
          ) : null}
          {scenarioData.status === 'enabled' ? (
            <Popconfirm
              title="Выключить сценарий?"
              description="Runtime перестанет выдавать его новым пользователям."
              okText="Выключить"
              cancelText="Отмена"
              onConfirm={() => transition.mutate('disable')}
            >
              <Button danger loading={transition.isPending}>Выключить</Button>
            </Popconfirm>
          ) : null}
        </Space>
      </header>

      {scenarioData.status === 'enabled' ? (
        <Alert
          style={{ marginBottom: 18 }}
          type="info"
          showIcon
          message="Активный сценарий защищён от изменений"
          description="Выключите сценарий, чтобы изменить поля или шаги."
        />
      ) : null}

      <div className="editor-grid">
        <Card className="content-card sticky-panel" title="Общие поля">
          <Form<ScenarioInput>
            key={scenarioData.updated_at}
            layout="vertical"
            requiredMark={false}
            disabled={!editable}
            initialValues={{
              name: scenarioData.name,
              description: scenarioData.description,
              page_pattern: scenarioData.page_pattern,
            }}
            onFinish={(values) => updateGeneral.mutate(values)}
          >
            <Form.Item label="Название" name="name" rules={[{ required: true, whitespace: true }, { max: 255 }]}>
              <Input />
            </Form.Item>
            <Form.Item
              label="Стартовая страница сценария"
              extra="Сценарий предложится, когда пользователь откроет этот адрес. Укажите только путь, например /add-item/title: без домена, параметров после ? и #."
              name="page_pattern"
              rules={scenarioPathnameRules}
            >
              <Input placeholder="/" />
            </Form.Item>
            <Form.Item label="Описание" name="description" rules={[{ max: 2000 }]}>
              <Input.TextArea rows={5} />
            </Form.Item>
            {editable ? (
              <Button type="primary" htmlType="submit" loading={updateGeneral.isPending} block>
                Сохранить общие поля
              </Button>
            ) : null}
          </Form>
        </Card>

        <section>
          <div className="toolbar">
            <div><b>Упорядоченные шаги</b><div className="muted">Порядок меняется кнопками — без drag-and-drop.</div></div>
            <span className="toolbar-spacer" />
            <Button
              type="primary"
              disabled={!editable || availableElementCount === 0}
              onClick={() => setEditingStep(null)}
            >
              + Добавить шаг
            </Button>
          </div>
          {elements.data.length === 0 ? (
            <Alert
              type="warning"
              showIcon
              message="В проекте пока нет элементов"
              description="Сначала добавьте элементы интерфейса, потом сможете собирать шаги сценария."
              action={<Button href={`/projects/${projectId}?tab=elements`} size="small">К inventory</Button>}
            />
          ) : null}
          {elements.data.length > 0 && availableElementCount === 0 ? (
            <Alert
              style={{ marginBottom: 16 }}
              type="warning"
              showIcon
              message="Для новых шагов нет элементов с привязанной страницей"
              description="Проверьте элементы проекта: у каждого шага должен быть элемент с заполненным page."
            />
          ) : null}
          {orderedSteps.length === 0 && elements.data.length > 0 ? (
            <Card className="empty-state"><h3>Шагов пока нет</h3><p className="muted">Добавьте первый целевой элемент сценария.</p></Card>
          ) : null}
          {orderedSteps.map((step, index) => (
            <Card className="step-card" key={step.id}>
              <div className="step-card-head">
                <span className="step-number">{index + 1}</span>
                <div className="step-summary">
                  <b>{step.title}</b>
                  <span>
                    {elementNames.get(step.element_id) || 'Неизвестный элемент'} · <code>{step.frontend_data.page_path}</code> · {formatStepAdvance(step)}
                  </span>
                  <p>{step.description}</p>
                </div>
                {editable ? (
                  <div className="step-actions">
                    <Button
                      size="small"
                      aria-label={`Поднять шаг ${index + 1}`}
                      disabled={index === 0 || reorder.isPending}
                      onClick={() => moveStep(index, -1)}
                    >↑</Button>
                    <Button
                      size="small"
                      aria-label={`Опустить шаг ${index + 1}`}
                      disabled={index === orderedSteps.length - 1 || reorder.isPending}
                      onClick={() => moveStep(index, 1)}
                    >↓</Button>
                    <Button size="small" onClick={() => setEditingStep(step)}>Изменить</Button>
                    <Popconfirm
                      title="Удалить шаг?"
                      okText="Удалить"
                      cancelText="Отмена"
                      onConfirm={() => removeStep.mutate(step.id)}
                    >
                      <Button size="small" danger>×</Button>
                    </Popconfirm>
                  </div>
                ) : null}
              </div>
            </Card>
          ))}
        </section>
      </div>

      {editingStep !== undefined ? (
        <StepModal
          projectId={projectId}
          key={editingStep?.id || 'new-step'}
          step={editingStep}
          allElements={elements.data}
          pending={saveStep.isPending}
          onCancel={() => setEditingStep(undefined)}
          onSave={(input) => saveStep.mutate({ step: editingStep, input })}
        />
      ) : null}
    </main>
  );
}
