'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Alert, Button, Form, Input } from 'antd';

export function LoginForm({ returnTo }: { returnTo: string }) {
  const router = useRouter();
  const [error, setError] = useState('');
  const [pending, setPending] = useState(false);

  async function submit({ password }: { password: string }) {
    setPending(true);
    setError('');
    try {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'content-type': 'application/json' },
        body: JSON.stringify({ password }),
      });
      if (!response.ok) {
        const body = (await response.json().catch(() => null)) as { message?: string } | null;
        throw new Error(body?.message || 'Не удалось войти');
      }
      router.replace(returnTo);
      router.refresh();
    } catch (submitError) {
      setError(submitError instanceof Error ? submitError.message : 'Не удалось войти');
    } finally {
      setPending(false);
    }
  }

  return (
    <Form layout="vertical" requiredMark={false} onFinish={submit} size="large">
      {error ? <Alert className="login-alert" type="error" message={error} showIcon /> : null}
      <Form.Item
        label="Общий пароль команды"
        name="password"
        rules={[{ required: true, message: 'Введите пароль' }]}
      >
        <Input.Password autoComplete="current-password" placeholder="••••••••••••" autoFocus />
      </Form.Item>
      <Button type="primary" htmlType="submit" loading={pending} block>
        Войти в панель
      </Button>
    </Form>
  );
}
