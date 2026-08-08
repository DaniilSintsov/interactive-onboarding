'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Alert, App, Button, Form, Input, Modal, Popconfirm, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { adminApi } from '@/shared/api/admin-api';
import type { ScenarioInput, ScenarioSummary } from '@/shared/api/types';
import { formatDate } from '@/shared/lib/format';
import { ScenarioStatus } from './scenario-status';

function openPreview(scenario: ScenarioSummary) {
  const url = new URL(scenario.page_pattern, process.env.NEXT_PUBLIC_PREVIEW_URL || 'http://localhost:3000');
  url.searchParams.set('preview', '1');
  window.open(url, '_blank', 'noopener,noreferrer');
}

export function ScenariosTab({ projectId }: { projectId: string }) {
  const [creating, setCreating] = useState(false);
  const [form] = Form.useForm<ScenarioInput>();
  const router = useRouter();
  const queryClient = useQueryClient();
  const { message } = App.useApp();
  const scenarios = useQuery({
    queryKey: ['scenarios', projectId],
    queryFn: () => adminApi.listScenarios(projectId),
  });
  const create = useMutation({
    mutationFn: (values: ScenarioInput) => adminApi.createScenario(projectId, values),
    onSuccess: (scenario) => {
      void queryClient.invalidateQueries({ queryKey: ['scenarios', projectId] });
      message.success('Сценарий создан');
      router.push(`/projects/${projectId}/scenarios/${scenario.id}`);
    },
    onError: (error) => message.error(error.message),
  });
  const remove = useMutation({
    mutationFn: adminApi.deleteScenario,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ['scenarios', projectId] });
      message.success('Сценарий удалён');
    },
    onError: (error) => message.error(error.message),
  });

  const columns: ColumnsType<ScenarioSummary> = [
    {
      title: 'Сценарий',
      key: 'name',
      render: (_, scenario) => (
        <div><b>{scenario.name}</b><div className="muted"><code>{scenario.page_pattern}</code></div></div>
      ),
    },
    { title: 'Статус', dataIndex: 'status', key: 'status', render: (status) => <ScenarioStatus status={status} /> },
    { title: 'Шаги', dataIndex: 'steps_count', key: 'steps', width: 90 },
    { title: 'Обновлён', dataIndex: 'updated_at', key: 'updated', render: formatDate, width: 150 },
    {
      title: '',
      key: 'actions',
      width: 260,
      render: (_, scenario) => (
        <Space wrap>
          {scenario.published_at ? <Button type="link" onClick={() => openPreview(scenario)}>Проверить</Button> : null}
          <Button type="link" onClick={() => router.push(`/projects/${projectId}/scenarios/${scenario.id}`)}>Открыть</Button>
          <Popconfirm
            title="Удалить сценарий?"
            description="Историческая аналитика сохранится."
            okText="Удалить"
            cancelText="Отмена"
            onConfirm={() => remove.mutate(scenario.id)}
          >
            <Button type="link" danger>Удалить</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <section>
      <div className="toolbar">
        <div><b>Сценарии помощи</b><div className="muted">Настройка, публикация и проверка.</div></div>
        <span className="toolbar-spacer" />
        <Button type="primary" onClick={() => setCreating(true)}>+ Новый сценарий</Button>
      </div>
      {scenarios.isError ? <Alert type="error" showIcon message={scenarios.error.message} /> : null}
      <Table<ScenarioSummary>
        rowKey="id"
        loading={scenarios.isPending}
        dataSource={scenarios.data?.items || []}
        columns={columns}
        pagination={false}
        scroll={{ x: 900 }}
      />
      <Modal
        title="Новый сценарий"
        open={creating}
        onCancel={() => setCreating(false)}
        onOk={() => form.submit()}
        okText="Создать"
        cancelText="Отмена"
        confirmLoading={create.isPending}
        destroyOnHidden
      >
        <Form
          form={form}
          layout="vertical"
          requiredMark={false}
          initialValues={{ page_pattern: '/', description: '' }}
          onFinish={(values) => create.mutate(values)}
        >
          <Form.Item label="Название" name="name" rules={[{ required: true, whitespace: true }, { max: 255 }]}>
            <Input placeholder="Разместить первое объявление" autoFocus />
          </Form.Item>
          <Form.Item label="Точка входа / page pattern" name="page_pattern" rules={[{ required: true, whitespace: true }, { max: 2048 }]}>
            <Input placeholder="/" />
          </Form.Item>
          <Form.Item label="Описание" name="description" rules={[{ max: 2000 }]}>
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>
    </section>
  );
}
