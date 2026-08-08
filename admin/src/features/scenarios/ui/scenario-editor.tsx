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
  initialStepValues,
  toStepInput,
  type StepFormValues,
} from '@/features/scenarios/model/step-form';
import { buildScenarioPreviewUrl } from '@/features/scenarios/model/preview';
import { ScenarioStatus } from './scenario-status';

function StepModal({
  step,
  elements,
  pending,
  onCancel,
  onSave,
}: {
  step: Step | null;
  elements: Element[];
  pending: boolean;
  onCancel: () => void;
  onSave: (input: StepInput) => void;
}) {
  const [form] = Form.useForm<StepFormValues>();
  const mode = Form.useWatch('mode', form);

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
        <Form.Item label="Элемент интерфейса" name="element_id" rules={[{ required: true, message: 'Выберите элемент' }]}>
          <Select
            showSearch
            optionFilterProp="label"
            placeholder="Выберите data-onboarding-id"
            options={elements.map((element) => ({
              value: element.id,
              label: `${element.label} · ${element.key}`,
            }))}
          />
        </Form.Item>
        <Form.Item
          label="Маршрут"
          name="page_path"
          rules={[
            { required: true, whitespace: true, message: 'Укажите маршрут' },
            { pattern: /^\//, message: 'Маршрут начинается с /' },
            { max: 2048 },
          ]}
        >
          <Input placeholder="/add-item/details" />
        </Form.Item>
        <div className="toolbar">
          <Form.Item label="Переход" name="mode" style={{ flex: 1, minWidth: 220 }}>
            <Select
              options={[
                { value: 'target_event', label: 'После действия с элементом' },
                { value: 'manual', label: 'Вручную из приложения' },
              ]}
            />
          </Form.Item>
          {mode !== 'manual' ? (
            <Form.Item label="Событие" name="event" style={{ flex: 1, minWidth: 180 }}>
              <Select
                options={[
                  { value: 'click', label: 'Click' },
                  { value: 'change', label: 'Change с непустым значением' },
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
    () => new Map((elements.data ?? []).map((element) => [element.id, element.label] as const)),
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
            <Form.Item label="Точка входа / page pattern" name="page_pattern" rules={[{ required: true, whitespace: true }, { max: 2048 }]}>
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
              disabled={!editable || elements.data.length === 0}
              onClick={() => setEditingStep(null)}
            >
              + Добавить шаг
            </Button>
          </div>
          {elements.data.length === 0 ? (
            <Alert
              type="warning"
              showIcon
              message="Сначала добавьте элемент интерфейса"
              action={<Button href={`/projects/${projectId}?tab=elements`} size="small">К элементам</Button>}
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
                    {elementNames.get(step.element_id) || 'Неизвестный элемент'} · <code>{step.frontend_data.page_path}</code> ·{' '}
                    {step.frontend_data.advance.mode === 'manual'
                      ? 'ручное завершение'
                      : step.frontend_data.advance.event}
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
          key={editingStep?.id || 'new-step'}
          step={editingStep}
          elements={elements.data}
          pending={saveStep.isPending}
          onCancel={() => setEditingStep(undefined)}
          onSave={(input) => saveStep.mutate({ step: editingStep, input })}
        />
      ) : null}
    </main>
  );
}
