import type { Metadata, Viewport } from "next";
import { AntdRegistry } from "@ant-design/nextjs-registry";

import { Providers } from "@/shared/ui/providers";
import { SiteShell } from "@/shared/ui/site-shell";

import "./globals.css";

export const metadata: Metadata = {
  title: "Рядом — объявления по соседству",
  description: "Тестовый классифайд для интерактивного онбординга",
};

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  themeColor: "#f7f8fb",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ru" data-scroll-behavior="smooth">
      <body>
        <AntdRegistry>
          <Providers>
            <SiteShell>{children}</SiteShell>
          </Providers>
        </AntdRegistry>
      </body>
    </html>
  );
}
