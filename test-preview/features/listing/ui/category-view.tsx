"use client";

import { useListingFlow } from "../controller/use-listing-flow";
import { useDraftStore } from "../model/draft-store";
import { FlowLayout } from "./flow-layout";

const categories = [
  ["Недвижимость", "⌂", "Квартиры, дома, участки"],
  ["Транспорт", "⌁", "Авто, мотоциклы, запчасти"],
  ["Электроника", "◫", "Телефоны, техника, игры"],
  ["Хобби и отдых", "✦", "Книги, спорт, коллекции"],
  ["Личные вещи", "♙", "Одежда, обувь, аксессуары"],
  ["Для дома", "▤", "Мебель, ремонт, растения"],
] as const;

export function CategoryView() {
  const selected = useDraftStore((state) => state.category);
  const setField = useDraftStore((state) => state.setField);
  const { go } = useListingFlow();

  const selectCategory = (category: string) => {
    setField("category", category);
    setField("subcategory", "");
    go("/add-item/title");
  };

  return (
    <FlowLayout
      current={1}
      title="Что вы продаёте?"
      lead="Выберите раздел — так покупатели быстрее найдут объявление."
      aside={
        <div className="aside-note note-yellow">
          <span>01</span>
          <strong>Категория влияет на поиск</strong>
          <p>Мы покажем подходящие параметры на следующем шаге.</p>
        </div>
      }
    >
      <div className="category-choice-list">
        {categories.map(([name, icon, description]) => (
          <button
            className={selected === name ? "choice-row is-selected" : "choice-row"}
            data-onboarding-id={name === "Хобби и отдых" ? "category-hobby" : undefined}
            data-onboarding-page="/add-item/category"
            data-onboarding-label="Категория «Хобби и отдых»"
            key={name}
            onClick={() => selectCategory(name)}
            type="button"
          >
            <span className="choice-icon">{icon}</span>
            <span>
              <strong>{name}</strong>
              <small>{description}</small>
            </span>
            <b>→</b>
          </button>
        ))}
      </div>
    </FlowLayout>
  );
}
