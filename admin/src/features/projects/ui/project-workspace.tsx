'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { Alert, App, Button, Card, Skeleton, Tabs } from 'antd';
import { adminApi } from '@/shared/api/admin-api';
import { ElementsTab } from '@/features/elements/ui/elements-tab';
import { ScenariosTab } from '@/features/scenarios/ui/scenarios-tab';
import { AnalyticsTab } from '@/features/analytics/ui/analytics-tab';

const tabKeys = new Set(['scenarios', 'elements', 'analytics']);

export function ProjectWorkspace({ projectId }: { projectId: string }) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { message } = App.useApp();
  const project = useQuery({
    queryKey: ['project', projectId],
    queryFn: () => adminApi.getProject(projectId),
  });
  const requestedTab = searchParams.get('tab') || 'scenarios';
  const activeTab = tabKeys.has(requestedTab) ? requestedTab : 'scenarios';

  if (project.isPending) return <Skeleton active paragraph={{ rows: 12 }} />;
  if (project.isError) return <Alert type="error" showIcon message={project.error.message} />;
  const projectData = project.data;

  const sdkSample = `const onboarding = createOnboarding({\n  projectKey: "${projectData.project_key}",\n  runtimeUrl: "/api/runtime",\n});\n\nawait onboarding.start({ userId: currentUser.id });`;

  async function copyKey() {
    try {
      await navigator.clipboard.writeText(projectData.project_key);
      message.success('project_key скопирован');
    } catch {
      message.error('Не удалось скопировать ключ');
    }
  }

  return (
    <main>
      <header className="page-heading">
        <div>
          <p className="eyebrow">Проект онбординга</p>
          <h1>{projectData.name}</h1>
          <p>Сценарии, карта элементов и результат прохождений.</p>
        </div>
        <Button href="/projects">← Все проекты</Button>
      </header>

      <div className="project-key-panel">
        <Card className="content-card" title="Публичный project_key">
          <code className="project-key">{projectData.project_key}</code>
          <Button onClick={copyKey}>Копировать ключ</Button>
        </Card>
        <Card className="content-card" title="Минимальное подключение SDK">
          <pre className="sdk-sample"><code>{sdkSample}</code></pre>
        </Card>
      </div>

      <Card className="content-card workspace-tabs">
        <Tabs
          activeKey={activeTab}
          onChange={(key) => router.replace(`/projects/${projectId}?tab=${key}`, { scroll: false })}
          items={[
            { key: 'scenarios', label: 'Сценарии', children: <ScenariosTab projectId={projectId} /> },
            { key: 'elements', label: 'Элементы', children: <ElementsTab projectId={projectId} /> },
            { key: 'analytics', label: 'Аналитика', children: <AnalyticsTab projectId={projectId} /> }
          ]}
        />
      </Card>
    </main>
  );
}
