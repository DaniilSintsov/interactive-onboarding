'use client';

import { useState } from 'react';
import { App as AntApp, ConfigProvider } from 'antd';
import ruRU from 'antd/locale/ru_RU';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: { staleTime: 15_000, retry: 1, refetchOnWindowFocus: false },
        },
      }),
  );

  return (
    <ConfigProvider
      locale={ruRU}
      theme={{
        token: {
          colorPrimary: '#00aaff',
          colorInfo: '#00aaff',
          colorSuccess: '#64b82a',
          colorWarning: '#f5a623',
          colorError: '#ff4053',
          colorText: '#17181a',
          colorBgLayout: '#f3f4f1',
          borderRadius: 10,
          borderRadiusLG: 16,
          fontFamily: '"Avenir Next", "Helvetica Neue", sans-serif',
        },
        components: {
          Button: { controlHeight: 40, fontWeight: 600 },
          Card: { headerFontSize: 16 },
          Table: { headerBg: '#f3f4f1', headerColor: '#555a60' },
        },
      }}
    >
      <AntApp>
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
      </AntApp>
    </ConfigProvider>
  );
}
