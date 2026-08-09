"use client";

import { useState } from "react";
import { Button, Tag } from "antd";

import { useListingFlow } from "@/features/listing/controller/use-listing-flow";
import { usePreviewUserStore } from "@/features/onboarding/model/preview-user-store";

export function SiteShell({ children }: { children: React.ReactNode }) {
  const [search, setSearch] = useState("");
  const userId = usePreviewUserStore((state) => state.userId);
  const onboarded = usePreviewUserStore((state) => state.onboarded);
  const regenerateUser = usePreviewUserStore((state) => state.regenerateUser);
  const { go, startOver } = useListingFlow();

  return (
    <div className="site-frame">
      <header className="site-header">
        <div className="header-topline">
          <button className="brand" type="button" onClick={() => go("/")}>
            <span className="brand-mark" aria-hidden="true">
              <i />
              <i />
              <i />
              <i />
            </span>
            <span>рядом</span>
          </button>

          <div className="user-mode">
            <span className="user-id-label">user_id</span>
            <code title={userId}>{userId || "—"}</code>
            <Tag
              color={onboarded ? "green" : "blue"}
              aria-live="polite"
            >
              {onboarded ? "Опытный" : "Новичок"}
            </Tag>
            <Button size="small" onClick={regenerateUser} aria-label="Сгенерировать новый user_id">
              Новый ID
            </Button>
          </div>

          <Button type="primary" onClick={startOver}>
            Разместить объявление
          </Button>
        </div>

        <div className="header-search">
          <button className="catalog-button" type="button" aria-label="Категории">
            <i /><i /><i /><i /><i /><i /><i /><i /><i />
          </button>
          <div className="search-field">
            <input
              value={search}
              onChange={(event) => setSearch(event.target.value)}
              placeholder="Поиск по объявлениям"
              aria-label="Поиск по объявлениям"
            />
            <button type="button">Найти</button>
          </div>
          <span className="header-location">● Москва</span>
        </div>
      </header>
      <main>{children}</main>
      <footer className="site-footer">
        <span>© 2026 Рядом</span>
        <span>Безопасность · Помощь · Для бизнеса</span>
      </footer>
    </div>
  );
}
