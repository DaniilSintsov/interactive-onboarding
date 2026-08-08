"use client";

import { useState } from "react";
import { Button, Segmented } from "antd";

import { useListingFlow } from "@/features/listing/controller/use-listing-flow";
import { type DemoUser, useDraftStore } from "@/features/listing/model/draft-store";

export function SiteShell({ children }: { children: React.ReactNode }) {
  const [search, setSearch] = useState("");
  const demoUser = useDraftStore((state) => state.demoUser);
  const setDemoUser = useDraftStore((state) => state.setDemoUser);
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
            <span>Режим</span>
            <Segmented<DemoUser>
              value={demoUser}
              onChange={setDemoUser}
              options={[
                { value: "demo-novice", label: "Новичок" },
                { value: "demo-expert", label: "Опытный" },
              ]}
              aria-label="Тип демо-пользователя"
            />
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
