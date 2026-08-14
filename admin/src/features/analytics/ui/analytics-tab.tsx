'use client';

import { useState } from 'react';
import type { Dayjs } from 'dayjs';
import { useQuery } from '@tanstack/react-query';
import { Alert, App, Button, Card, DatePicker, Empty, Select, Skeleton, Statistic, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { adminApi, downloadReport } from '@/shared/api/admin-api';
import type { ScenarioAnalytics, ScenarioSummary, StepAnalytics } from '@/shared/api/types';
import { formatPercent, periodQuery } from '@/shared/lib/format';
import { ScenarioStatus } from '@/features/scenarios/ui/scenario-status';

type ScenarioAnalyticsRow = {
  scenario: ScenarioSummary;
  analytics: ScenarioAnalytics;
};

export function AnalyticsTab({ projectId }: { projectId: string }) {
  const [scenarioChoice, setScenarioChoice] = useState('');
  const [range, setRange] = useState<[Dayjs, Dayjs] | null>(null);
  const [downloading, setDownloading] = useState(false);
  const { message } = App.useApp();
  const period = periodQuery(range);
  const scenarios = useQuery({
    queryKey: ['scenarios', projectId],
    queryFn: () => adminApi.listScenarios(projectId),
  });
  const scenarioItems = scenarios.data?.items || [];
  const selectedScenarioId = scenarioChoice || scenarioItems[0]?.id || '';
  const projectAnalytics = useQuery({
    queryKey: ['project-analytics', projectId, period],
    queryFn: () => adminApi.getProjectAnalytics(projectId, period),
  });
  const allScenarioAnalytics = useQuery({
    queryKey: ['scenario-analytics-all', projectId, period, scenarioItems.map(({ id }) => id)],
    queryFn: async () =>
      // ponytail: one request per scenario; add bulk project analytics endpoint if scenario count becomes real bottleneck.
      Promise.all(
        scenarioItems.map(async (scenario) => ({
          scenario,
          analytics: await adminApi.getScenarioAnalyticsTotal(scenario.id, period),
        })),
      ),
    enabled: scenarioItems.length > 0,
  });
  const scenarioAnalytics = useQuery({
    queryKey: ['scenario-analytics', selectedScenarioId, period],
    queryFn: () => adminApi.getScenarioAnalytics(selectedScenarioId, period),
    enabled: Boolean(selectedScenarioId),
  });

  async function getPdf() {
    if (!selectedScenarioId) return;
    setDownloading(true);
    try {
      await downloadReport(selectedScenarioId, period);
    } catch (error) {
      message.error(error instanceof Error ? error.message : 'Не удалось скачать PDF');
    } finally {
      setDownloading(false);
    }
  }

  const columns: ColumnsType<StepAnalytics> = [
    { title: '#', dataIndex: 'position', key: 'position', width: 60 },
    { title: 'Шаг', dataIndex: 'title', key: 'title', render: (value) => <b>{value}</b> },
    { title: 'Показы', dataIndex: 'shown', key: 'shown' },
    { title: 'Выполнено', dataIndex: 'completed', key: 'completed' },
    { title: 'Пропущено', dataIndex: 'skipped', key: 'skipped' },
    { title: 'Выполнение', dataIndex: 'completion_rate', key: 'completion', render: formatPercent },
    { title: 'Отвал шага', dataIndex: 'drop_off_rate', key: 'dropoff', render: formatPercent },
  ];
  const scenarioColumns: ColumnsType<ScenarioAnalyticsRow> = [
    {
      title: 'Сценарий',
      key: 'name',
      render: (_, row) => (
        <div>
          <b>{row.scenario.name}</b>
          <div className="muted"><code>{row.scenario.page_pattern}</code></div>
        </div>
      ),
    },
    {
      title: 'Статус',
      key: 'status',
      width: 140,
      render: (_, row) => <ScenarioStatus status={row.scenario.status} />,
    },
    { title: 'Стартов', dataIndex: ['analytics', 'started'], key: 'started', width: 110 },
    { title: 'Завершено', dataIndex: ['analytics', 'completed'], key: 'completed', width: 120 },
    { title: 'Пропущено', dataIndex: ['analytics', 'skipped'], key: 'skipped', width: 120 },
    {
      title: 'Выполнение',
      dataIndex: ['analytics', 'completion_rate'],
      key: 'completion_rate',
      width: 130,
      render: formatPercent,
    },
    {
      title: 'Среднее время',
      dataIndex: ['analytics', 'average_completion_time_seconds'],
      key: 'average_completion_time_seconds',
      width: 140,
      render: (value: number) => `${Math.round(value)} сек.`,
    },
  ];

  return (
    <section>
      <div className="toolbar">
        <DatePicker.RangePicker
          allowClear
          format="DD.MM.YYYY"
          onChange={(value) =>
            setRange(value?.[0] && value[1] ? [value[0], value[1]] : null)
          }
        />
        <span className="muted">to считается исключительно: выбран последний полный день.</span>
      </div>
      {projectAnalytics.isError ? <Alert type="error" showIcon message={projectAnalytics.error.message} /> : null}
      {projectAnalytics.isPending ? <Skeleton active /> : null}
      {projectAnalytics.data ? (
        <div className="metric-grid">
          <Card className="metric-card"><Statistic title="Начато сессий" value={projectAnalytics.data.sessions_started} /></Card>
          <Card className="metric-card"><Statistic title="Завершено" value={projectAnalytics.data.sessions_completed} /></Card>
          <Card className="metric-card"><Statistic title="Пропущено" value={projectAnalytics.data.sessions_skipped} /></Card>
          <Card className="metric-card"><Statistic title="Доля завершений" value={formatPercent(projectAnalytics.data.completion_rate)} /></Card>
        </div>
      ) : null}

      <Card className="content-card" title="Все сценарии проекта" style={{ marginTop: 18 }}>
        {scenarios.isError ? <Alert type="error" showIcon message={scenarios.error.message} /> : null}
        {allScenarioAnalytics.isError ? (
          <Alert type="error" showIcon message={allScenarioAnalytics.error.message} />
        ) : null}
        {scenarios.isPending || allScenarioAnalytics.isPending ? <Skeleton active /> : null}
        {!scenarioItems.length && !scenarios.isPending ? <Empty description="Сценариев пока нет" /> : null}
        {allScenarioAnalytics.data?.length ? (
          <Table<ScenarioAnalyticsRow>
            rowKey={(row) => row.scenario.id}
            columns={scenarioColumns}
            dataSource={allScenarioAnalytics.data}
            pagination={false}
            scroll={{ x: 980 }}
          />
        ) : null}
      </Card>

      <Card className="content-card" title="Детали выбранного сценария" style={{ marginTop: 18 }}>
        <div className="toolbar">
          <Select
            className="analytics-scenario-picker"
            value={selectedScenarioId || undefined}
            placeholder="Выберите сценарий"
            loading={scenarios.isPending}
            options={scenarios.data?.items.map((scenario) => ({ value: scenario.id, label: scenario.name }))}
            onChange={setScenarioChoice}
          />
          <span className="toolbar-spacer" />
          <Button disabled={!selectedScenarioId} loading={downloading} onClick={getPdf}>Скачать PDF</Button>
        </div>
        {scenarioAnalytics.isError ? <Alert type="error" showIcon message={scenarioAnalytics.error.message} /> : null}
        {scenarioAnalytics.data ? (
          <>
            <div className="toolbar">
              <b>{scenarioAnalytics.data.started} стартов</b>
              <span className="muted">Завершено {formatPercent(scenarioAnalytics.data.completion_rate)}</span>
              <span className="muted">Среднее время {Math.round(scenarioAnalytics.data.average_completion_time_seconds)} сек.</span>
            </div>
            <Table<StepAnalytics>
              rowKey="step_id"
              columns={columns}
              dataSource={scenarioAnalytics.data.steps}
              pagination={false}
              scroll={{ x: 840 }}
            />
          </>
        ) : null}
        {scenarioAnalytics.isPending && selectedScenarioId ? <Skeleton active /> : null}
      </Card>
    </section>
  );
}
