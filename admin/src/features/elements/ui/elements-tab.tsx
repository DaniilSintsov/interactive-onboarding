'use client';

import { useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Alert, App, Button, Form, Input, Modal, Popconfirm, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { adminApi } from '@/shared/api/admin-api';
import type { Element, ElementInput } from '@/shared/api/types';

type ElementDraft = ElementInput & { id?: string };

export function ElementsTab({ projectId }: { projectId: string }) {
  const [editing, setEditing] = useState<Element | null | undefined>(undefined);
  const [form] = Form.useForm<ElementInput>();
  const queryClient = useQueryClient();
  const { message } = App.useApp();
  const elements = useQuery({
    queryKey: ['elements', projectId],
    queryFn: () => adminApi.listElements(projectId),
  });

  const save = useMutation({
    mutationFn: ({ id, ...values }: ElementDraft) =>
      id
        ? adminApi.updateElement(projectId, id, values)
        : adminApi.createElement(projectId, values),
    onSuccess: () => {
      void Promise.all([
        queryClient.invalidateQueries({ queryKey: ['elements', projectId] }),
        queryClient.invalidateQueries({ queryKey: ['project', projectId] }),
      ]);
      setEditing(undefined);
      form.resetFields();
      message.success('Элемент сохранён');
    },
    onError: (error) => message.error(error.message),
  });
  const remove = useMutation({
    mutationFn: (elementId: string) => adminApi.deleteElement(projectId, elementId),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ['elements', projectId] });
      message.success('Элемент удалён');
    },
    onError: (error) => message.error(error.message),
  });

  function open(element: Element | null) {
    setEditing(element);
    form.setFieldsValue(
      element
        ? { key: element.key, label: element.label, description: element.description }
        : { key: '', label: '', description: '' },
    );
  }

  const columns: ColumnsType<Element> = [
    { title: 'Название', dataIndex: 'label', key: 'label', render: (value) => <b>{value}</b> },
    {
      title: 'Технический key',
      dataIndex: 'key',
      key: 'key',
      render: (value) => <code>{value}</code>,
    },
    { title: 'Описание', dataIndex: 'description', key: 'description', render: (value) => value || '—' },
    {
      title: '',
      key: 'actions',
      width: 180,
      render: (_, element) => (
        <Space>
          <Button type="link" onClick={() => open(element)}>Изменить</Button>
          <Popconfirm
            title="Удалить элемент?"
            description="Шаги, которые ссылаются на него, заблокируют удаление."
            okText="Удалить"
            cancelText="Отмена"
            onConfirm={() => remove.mutate(element.id)}
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
        <div>
          <b>Карта интерфейса</b>
          <div className="muted">key совпадает с data-onboarding-id в продукте.</div>
        </div>
        <span className="toolbar-spacer" />
        <Button type="primary" onClick={() => open(null)}>+ Добавить элемент</Button>
      </div>
      {elements.isError ? <Alert type="error" showIcon message={elements.error.message} /> : null}
      <Table<Element>
        rowKey="id"
        loading={elements.isPending}
        dataSource={elements.data || []}
        columns={columns}
        pagination={false}
        scroll={{ x: 760 }}
      />
      <Modal
        title={editing ? 'Изменить элемент' : 'Новый элемент'}
        open={editing !== undefined}
        onCancel={() => setEditing(undefined)}
        okText="Сохранить"
        cancelText="Отмена"
        confirmLoading={save.isPending}
        onOk={() => form.submit()}
        destroyOnHidden
      >
        <Form
          form={form}
          layout="vertical"
          requiredMark={false}
          onFinish={(values) => save.mutate({ ...values, id: editing?.id })}
        >
          <Form.Item
            label="Технический key"
            name="key"
            extra="Тестовый путь: category-hobby, listing-title, subcategory-ebooks, listing-photo, listing-description, listing-price, publish-listing"
            rules={[
              { required: true, whitespace: true, message: 'Укажите key' },
              { max: 255 },
              { pattern: /^[a-z0-9]+(?:-[a-z0-9]+)*$/, message: 'Используйте lowercase kebab-case' },
            ]}
          >
            <Input placeholder="listing-title" />
          </Form.Item>
          <Form.Item label="Название" name="label" rules={[{ required: true, whitespace: true }, { max: 255 }]}>
            <Input placeholder="Поле названия объявления" />
          </Form.Item>
          <Form.Item label="Описание" name="description" rules={[{ max: 2000 }]}>
            <Input.TextArea rows={3} placeholder="Где находится и зачем используется" />
          </Form.Item>
        </Form>
      </Modal>
    </section>
  );
}
