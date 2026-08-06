'use client';

import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Skeleton } from 'antd';
import { adminApi } from '@/shared/api/admin-api';
import { formatDate } from '@/shared/lib/format';

export function ProjectsPage() {
  const projects = useQuery({ queryKey: ['projects'], queryFn: adminApi.listProjects });

  return (
    <main>
      <header className="page-heading">
        <div>
          <p className="eyebrow">01 / Управление</p>
          <h1>Проекты онбординга</h1>
          <p>Каждый проект связывает элементы продукта, сценарии и аналитику в одном контуре.</p>
        </div>
        <Button type="primary" size="large" href="/projects/create">+ Новый проект</Button>
      </header>

      {projects.isPending ? <Skeleton active paragraph={{ rows: 6 }} /> : null}
      {projects.isError ? <Alert type="error" showIcon message={projects.error.message} /> : null}
      {projects.data?.items.length === 0 ? (
        <Card className="content-card empty-state">
          <Empty description="Проектов пока нет">
            <Button type="primary" href="/projects/create">Создать первый проект</Button>
          </Empty>
        </Card>
      ) : null}
      {projects.data?.items.length ? (
        <div className="project-grid">
          {projects.data.items.map((project, index) => (
            <Link key={project.id} href={`/projects/${project.id}`}>
              <Card className="project-tile" styles={{ body: { padding: 24 } }}>
                <span className="project-index">PROJECT / {String(index + 1).padStart(2, '0')}</span>
                <h3>{project.name}</h3>
                <div className="project-tile-footer">
                  <span>Обновлён {formatDate(project.updated_at)}</span>
                  <b>Открыть →</b>
                </div>
              </Card>
            </Link>
          ))}
        </div>
      ) : null}
    </main>
  );
}
