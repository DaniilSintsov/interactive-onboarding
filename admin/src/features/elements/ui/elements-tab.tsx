'use client';

import { useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Space, Table, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { parseOnboardingCatalogMetadata } from '@/features/elements/model/onboarding-catalog';
import { adminApi } from '@/shared/api/admin-api';

type InventoryRow = {
  rowKey: string;
  key: string;
  label: string;
  page_paths: string[];
  status: 'available' | 'stale';
};

export function ElementsTab({ projectId }: { projectId: string }) {
  const elements = useQuery({
    queryKey: ['elements', projectId],
    queryFn: () => adminApi.listElements(projectId),
  });

  const rows = useMemo<InventoryRow[]>(() => {
    return (elements.data ?? [])
      .map((element): InventoryRow => {
        const metadata = parseOnboardingCatalogMetadata(element.description);
        return {
          rowKey: element.id,
          key: element.key,
          label: element.label,
          page_paths: metadata?.page_paths ?? [],
          status: metadata && metadata.page_paths.length > 0 ? 'available' : 'stale',
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
      title: 'Маршруты',
      key: 'page_paths',
      render: (_, row) =>
        row.page_paths.length > 0 ? (
          <Space size={[6, 6]} wrap>
            {row.page_paths.map((pagePath) => (
              <Tag key={pagePath}>
                <code>{pagePath}</code>
              </Tag>
            ))}
          </Space>
        ) : (
          <span className="muted">Нет маршрутов в каталоге</span>
        ),
    },
    {
      title: 'Статус',
      key: 'status',
      render: (_, row) =>
        row.status === 'available' ? (
          <Tag color="green">В каталоге и в БД</Tag>
        ) : (
          <Tag color="red">Нет в текущем CI-каталоге</Tag>
        ),
    },
  ];

  const staleCount = rows.filter((row) => row.status === 'stale').length;

  return (
    <section>
      <div className="toolbar">
        <div>
          <b>CI-инвентарь интерфейса</b>
          <div className="muted">Источник: элементы проекта, синхронизированные при деплое test-preview.</div>
        </div>
      </div>
      {elements.isError ? <Alert type="error" showIcon message={elements.error.message} /> : null}
      {!elements.isError && !elements.isPending ? (
        <Alert
          style={{ marginBottom: 16 }}
          type={staleCount > 0 ? 'warning' : 'info'}
          showIcon
          message="Вкладка только для чтения"
          description={
            staleCount > 0
              ? `В текущем CI-каталоге отсутствуют элементы: ${staleCount}. Для новых шагов они недоступны.`
              : 'Маршруты и названия берутся из CI-каталога. Ручное создание и редактирование элементов отключено.'
          }
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
