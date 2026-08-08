'use client';

import { useRouter } from 'next/navigation';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { App, Button, Card, Form, Input } from 'antd';
import { adminApi } from '@/shared/api/admin-api';

export function CreateProjectPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { message } = App.useApp();
  const createProject = useMutation({
    mutationFn: (name: string) => adminApi.createProject(name),
    onSuccess: (project) => {
      void queryClient.invalidateQueries({ queryKey: ['projects'] });
      message.success('Проект создан');
      router.push(`/projects/${project.id}`);
    },
    onError: (error) => message.error(error.message),
  });

  return (
    <main>
      <header className="page-heading">
        <div>
          <p className="eyebrow">Новый контур</p>
          <h1>Создать проект</h1>
          <p>Сейчас нужно только название. Элементы интерфейса добавляются после создания.</p>
        </div>
        <Button href="/projects">← К проектам</Button>
      </header>
      <Card className="content-card" styles={{ body: { padding: 30 } }}>
        <Form
          layout="vertical"
          requiredMark={false}
          onFinish={({ name }: { name: string }) => createProject.mutate(name.trim())}
        >
          <Form.Item
            label="Название проекта"
            name="name"
            rules={[{ required: true, whitespace: true, message: 'Укажите название' }, { max: 255 }]}
          >
            <Input size="large" placeholder="Например, Размещение объявления" autoFocus />
          </Form.Item>
          <div className="technical-note">
            API получит <code>elements: []</code>. Публичный project_key появится автоматически.
          </div>
          <div className="form-actions">
            <Button href="/projects">Отмена</Button>
            <Button type="primary" htmlType="submit" loading={createProject.isPending}>Создать проект</Button>
          </div>
        </Form>
      </Card>
    </main>
  );
}
