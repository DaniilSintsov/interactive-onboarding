'use client';

import { useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { adminApi } from '@/shared/api/admin-api';

type InventoryRow = {
  rowKey: string;
  key: string;
  label: string;
  page: string;
};

export function ElementsTab({ projectId }: { projectId: string }) {
  const elements = useQuery({
    queryKey: ['elements', projectId],
    queryFn: () => adminApi.listElements(projectId),
  });

  const rows = useMemo<InventoryRow[]>(() => {
    return (elements.data ?? [])
      .map((element): InventoryRow => {
        return {
          rowKey: element.id,
          key: element.key,
          label: element.label,
          page: element.page.trim(),
        };
      })
      .sort((left, right) => left.key.localeCompare(right.key));
  }, [elements.data]);

  const columns: ColumnsType<InventoryRow> = [
    {
      title: 'Элемент',
      key: 'element',
      render: (_, row) => (
        <div>
          <b>{row.label}</b>
          <div className="muted">
            <code>{row.key}</code>
          </div>
        </div>
      ),
    },
    {
      title: 'Страница',
      key: 'page',
      render: (_, row) =>
        row.page ? (
          <Tag>
            <code>{row.page}</code>
          </Tag>
        ) : (
          <span className="muted">Страница не указана</span>
        ),
    },
  ];

  return (
    <section>
      <div className="toolbar">
        <div>
          <b>Элементы интерфейса проекта</b>
          <div className="muted">Только чтение. Для каждого элемента показана привязанная страница.</div>
        </div>
      </div>
      {elements.isError ? <Alert type="error" showIcon message={elements.error.message} /> : null}
      {!elements.isError && !elements.isPending ? (
        <Alert
          style={{ marginBottom: 16 }}
          type="info"
          showIcon
          message="Вкладка только для чтения"
          description="Название, ключ и страница приходят из backend. Ручное создание и редактирование элементов отключено."
        />
      ) : null}
      <Table<InventoryRow>
        rowKey="rowKey"
        loading={elements.isPending}
        dataSource={rows}
        columns={columns}
        pagination={false}
        scroll={{ x: 760 }}
      />
    </section>
  );
}
