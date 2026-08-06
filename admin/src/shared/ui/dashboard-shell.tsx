'use client';

import { useRouter } from 'next/navigation';
import { Button, Drawer, Layout, Menu } from 'antd';
import { useUiStore } from '@/shared/model/ui-store';

const navigation = [{ key: '/projects', label: 'Проекты онбординга' }];

function Brand() {
  return (
    <div className="brand-lockup">
      <span className="brand-mark" aria-hidden="true">
        <i />
        <i />
        <i />
        <i />
      </span>
      <span>ONBOARD CONTROL</span>
    </div>
  );
}

function Navigation({ close }: { close?: () => void }) {
  const router = useRouter();
  return (
    <>
      <Brand />
      <p className="nav-caption">Workspace / MVP</p>
      <Menu
        mode="inline"
        theme="dark"
        selectedKeys={['/projects']}
        items={navigation}
        onClick={({ key }) => {
          router.push(key);
          close?.();
        }}
      />
      <div className="nav-footer">
        <span className="system-dot" /> API через защищённый BFF
      </div>
    </>
  );
}

export function DashboardShell({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const navigationOpen = useUiStore((state) => state.navigationOpen);
  const openNavigation = useUiStore((state) => state.openNavigation);
  const closeNavigation = useUiStore((state) => state.closeNavigation);

  async function logout() {
    await fetch('/api/auth/logout', { method: 'POST' });
    router.replace('/login');
    router.refresh();
  }

  return (
    <Layout className="dashboard-layout">
      <Layout.Sider width={252} className="dashboard-sider" breakpoint="lg" collapsedWidth={0}>
        <Navigation />
      </Layout.Sider>
      <Layout>
        <header className="dashboard-topbar">
          <Button className="mobile-menu" onClick={openNavigation} aria-label="Открыть навигацию">
            ☰
          </Button>
          <div className="topbar-context">
            <span>Контур администратора</span>
            <b>Production-ready workspace</b>
          </div>
          <Button onClick={logout}>Выйти</Button>
        </header>
        <Layout.Content className="dashboard-content">{children}</Layout.Content>
      </Layout>
      <Drawer
        placement="left"
        width={272}
        open={navigationOpen}
        onClose={closeNavigation}
        className="mobile-navigation"
      >
        <Navigation close={closeNavigation} />
      </Drawer>
    </Layout>
  );
}
