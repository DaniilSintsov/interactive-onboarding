"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Alert, Input, Skeleton } from "antd";

import { useListingFlow } from "../controller/use-listing-flow";
import { useDraftStore } from "../model/draft-store";
import { FlowLayout } from "./flow-layout";

type SuggestionsResponse = { items: string[] };

async function loadSuggestions(title: string): Promise<SuggestionsResponse> {
  const response = await fetch(`/api/catalog/suggestions?q=${encodeURIComponent(title)}`);
  if (!response.ok) throw new Error("suggestions_failed");
  return response.json();
}

export function TitleView() {
  const title = useDraftStore((state) => state.title);
  const subcategory = useDraftStore((state) => state.subcategory);
  const setField = useDraftStore((state) => state.setField);
  const [error, setError] = useState("");
  const { go } = useListingFlow();
  const cleanTitle = title.trim();
  const suggestions = useQuery({
    queryKey: ["subcategory-suggestions", cleanTitle],
    queryFn: () => loadSuggestions(cleanTitle),
    enabled: cleanTitle.length >= 2,
  });

  const selectSubcategory = (value: string) => {
    if (!cleanTitle) {
      setError("Сначала введите название объявления");
      return;
    }
    setField("subcategory", value);
    go("/add-item/details");
  };

  return (
    <FlowLayout
      current={2}
      title="Назовите вещь"
      lead="Точное название работает лучше описания из общих слов."
      aside={
        <div className="aside-note note-blue">
          <span>Подсказка</span>
          <strong>Пишите как покупатель</strong>
          <p>Марка, модель и тип вещи помогают попасть в точный поиск.</p>
          <div className="example-title">Электронная книга PocketBook 628</div>
        </div>
      }
    >
      <div className="title-field-block">
        <label htmlFor="listing-title">Название объявления</label>
        <Input
          id="listing-title"
          data-onboarding-id="listing-title"
          data-onboarding-page="/add-item/title"
          data-onboarding-label="Название объявления"
          value={title}
          status={error ? "error" : undefined}
          onChange={(event) => {
            setField("title", event.target.value);
            setError("");
          }}
          maxLength={80}
          showCount
          placeholder="Например, электронная книга"
          size="large"
        />
        {error ? <span className="field-error">{error}</span> : null}
      </div>

      <div className="suggestions-block">
        <div className="suggestion-heading">
          <strong>Подходящая категория</strong>
          <small>уточним раздел автоматически</small>
        </div>

        {cleanTitle.length < 2 ? (
          <div className="empty-suggestion">Введите хотя бы два символа</div>
        ) : suggestions.isPending ? (
          <Skeleton active paragraph={{ rows: 2 }} title={false} />
        ) : suggestions.isError ? (
          <Alert type="error" showIcon message="Не удалось загрузить варианты" />
        ) : (
          <div className="subcategory-list">
            {suggestions.data.items.map((item) => (
              <button
                className={subcategory === item ? "subcategory-row is-selected" : "subcategory-row"}
                data-onboarding-id={item === "Электронные книги" ? "subcategory-ebooks" : undefined}
                data-onboarding-page="/add-item/title"
                data-onboarding-label="Подкатегория «Электронные книги»"
                key={item}
                onClick={() => selectSubcategory(item)}
                type="button"
              >
                <span>Хобби и отдых <b>›</b> {item}</span>
                <strong>{subcategory === item ? "Выбрано" : "Выбрать"}</strong>
              </button>
            ))}
          </div>
        )}
      </div>
    </FlowLayout>
  );
}
