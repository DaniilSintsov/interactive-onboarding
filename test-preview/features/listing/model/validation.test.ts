import { describe, expect, it } from "vitest";

import { initialDraft } from "./draft-store";
import { validateDraft } from "./validation";

describe("validateDraft", () => {
  it("проверяет обязательные поля и цену", () => {
    expect(Object.keys(validateDraft(initialDraft))).toEqual([
      "category",
      "title",
      "subcategory",
      "photoName",
      "description",
      "price",
    ]);

    expect(
      validateDraft({
        ...initialDraft,
        category: "Хобби и отдых",
        title: "Электронная книга",
        subcategory: "Электронные книги",
        photoName: "reader.jpg",
        description: "Работает, без царапин",
        price: "3500",
      }),
    ).toEqual({});
  });
});
