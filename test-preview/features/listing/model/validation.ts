import type { ListingDraft } from "./draft-store";

export type RequiredDraftField =
  | "category"
  | "title"
  | "subcategory"
  | "photoName"
  | "description"
  | "price";

export type DraftErrors = Partial<Record<RequiredDraftField, string>>;

export function validateDraft(draft: ListingDraft): DraftErrors {
  const errors: DraftErrors = {};

  if (!draft.category) errors.category = "Выберите категорию";
  if (!draft.title.trim()) errors.title = "Введите название";
  if (!draft.subcategory) errors.subcategory = "Выберите подкатегорию";
  if (!draft.photoName) errors.photoName = "Добавьте хотя бы одну фотографию";
  if (!draft.description.trim()) errors.description = "Добавьте описание";
  if (!draft.price.trim() || Number(draft.price) <= 0) {
    errors.price = "Укажите цену больше нуля";
  }

  return errors;
}
