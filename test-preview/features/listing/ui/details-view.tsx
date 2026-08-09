"use client";

import { useState } from "react";
import { Alert, Button, Input, Radio, Segmented, Select, Switch } from "antd";

import { useListingFlow } from "../controller/use-listing-flow";
import { type ListingDraft, useDraftStore } from "../model/draft-store";
import { type DraftErrors, validateDraft } from "../model/validation";
import { FlowLayout } from "./flow-layout";

const fieldTargets: Partial<Record<keyof DraftErrors, string>> = {
  photoName: "listing-photo",
  description: "listing-description",
  price: "listing-price",
};

export function DetailsView() {
  const draft = useDraftStore();
  const { setField } = draft;
  const [errors, setErrors] = useState<DraftErrors>({});
  const { go } = useListingFlow();

  const clearError = (field: keyof DraftErrors) => {
    setErrors((current) => ({ ...current, [field]: undefined }));
  };

  const publish = () => {
    const nextErrors = validateDraft(draft);
    setErrors(nextErrors);
    const firstField = Object.keys(nextErrors)[0] as keyof DraftErrors | undefined;
    if (firstField) {
      const target = fieldTargets[firstField];
      if (target) {
        document.querySelector<HTMLElement>(`[data-onboarding-id="${target}"]`)?.focus();
      }
      return;
    }
    go("/add-item/success");
  };

  const missingPreviousStep =
    !draft.category || !draft.title.trim() || !draft.subcategory;

  return (
    <FlowLayout
      current={3}
      title="Расскажите о вещи"
      lead="Честные детали и хорошие фотографии экономят время вам и покупателю."
      aside={
        <div className="details-summary">
          <span className="section-kicker">Ваше объявление</span>
          <div className="summary-thumb">{draft.photoName ? "✓" : "+ фото"}</div>
          <strong>{draft.title || "Название появится здесь"}</strong>
          <p>{draft.subcategory || "Категория не выбрана"}</p>
          <b>{draft.price ? `${Number(draft.price).toLocaleString("ru-RU")} ₽` : "Цена не указана"}</b>
          <small>Черновик сохраняется в этой вкладке автоматически</small>
        </div>
      }
    >
      {missingPreviousStep ? (
        <Alert
          className="form-alert"
          type="warning"
          showIcon
          message="Заполните предыдущие шаги"
          description="Для публикации нужны категория, название и подкатегория."
          action={<Button onClick={() => go("/add-item/category")}>Вернуться</Button>}
        />
      ) : null}

      <div className="details-form">
        <section className="form-section">
          <div className="form-section-title">
            <span>01</span>
            <div><h2>Состояние и тип продажи</h2><p>Помогите покупателю понять предложение.</p></div>
          </div>
          <div className="form-control">
            <label>Состояние</label>
            <Segmented<ListingDraft["condition"]>
              block
              value={draft.condition}
              onChange={(value) => setField("condition", value)}
              options={[{ value: "new", label: "Новое" }, { value: "used", label: "Б/у" }]}
            />
          </div>
          <div className="form-control">
            <label>Вид объявления</label>
            <Radio.Group
              value={draft.saleType}
              onChange={(event) => setField("saleType", event.target.value)}
            >
              <Radio value="personal">Продаю своё</Radio>
              <Radio value="resale">Товар приобретён на продажу</Radio>
            </Radio.Group>
          </div>
        </section>

        <section className="form-section">
          <div className="form-section-title">
            <span>02</span>
            <div><h2>Фотографии</h2><p>Можно добавить до 10 снимков.</p></div>
          </div>
          <div className="photo-uploader">
            <label htmlFor="listing-photo">Главная фотография</label>
            <input
              id="listing-photo"
              data-onboarding-id="listing-photo"
              data-onboarding-page="/add-item/details"
              data-onboarding-label="Главная фотография"
              type="file"
              accept="image/png,image/jpeg,image/webp"
              aria-invalid={Boolean(errors.photoName)}
              onChange={(event) => {
                setField("photoName", event.target.files?.[0]?.name ?? "");
                clearError("photoName");
              }}
            />
            <small>{draft.photoName || "JPG, PNG или WebP · до 10 МБ"}</small>
            {errors.photoName ? <span className="field-error">{errors.photoName}</span> : null}
          </div>
        </section>

        <section className="form-section">
          <div className="form-section-title">
            <span>03</span>
            <div><h2>Описание</h2><p>Состояние, комплект и причина продажи.</p></div>
          </div>
          <div className="form-control">
            <label htmlFor="listing-description">Описание объявления</label>
            <Input.TextArea
              id="listing-description"
              data-onboarding-id="listing-description"
              data-onboarding-page="/add-item/details"
              data-onboarding-label="Описание объявления"
              value={draft.description}
              status={errors.description ? "error" : undefined}
              onChange={(event) => {
                setField("description", event.target.value);
                clearError("description");
              }}
              maxLength={1500}
              showCount
              autoSize={{ minRows: 6, maxRows: 12 }}
              placeholder="Например: пользовались год, экран без царапин, зарядка в комплекте…"
            />
            {errors.description ? <span className="field-error">{errors.description}</span> : null}
          </div>
          <label className="switch-row">
            <Switch
              checked={draft.hasMultiple}
              onChange={(checked) => setField("hasMultiple", checked)}
            />
            <span><strong>Есть несколько штук</strong><small>Покупатель сможет выбрать количество</small></span>
          </label>
        </section>

        <section className="form-section">
          <div className="form-section-title">
            <span>04</span>
            <div><h2>Место и цена</h2><p>Покажем объявление людям поблизости.</p></div>
          </div>
          <div className="form-two-columns">
            <div className="form-control">
              <label htmlFor="listing-address">Адрес</label>
              <Input
                id="listing-address"
                value={draft.address}
                onChange={(event) => setField("address", event.target.value)}
              />
            </div>
            <div className="form-control">
              <label htmlFor="listing-price">Цена</label>
              <Input
                id="listing-price"
                data-onboarding-id="listing-price"
                data-onboarding-page="/add-item/details"
                data-onboarding-label="Цена"
                value={draft.price}
                status={errors.price ? "error" : undefined}
                inputMode="numeric"
                suffix="₽"
                onChange={(event) => {
                  setField("price", event.target.value.replace(/\D/g, ""));
                  clearError("price");
                }}
                placeholder="0"
              />
              {errors.price ? <span className="field-error">{errors.price}</span> : null}
            </div>
          </div>
          <div className="map-placeholder" aria-label="Схема района">
            <span className="map-road road-one" /><span className="map-road road-two" />
            <b>●</b><small>Москва, Белорусская</small>
          </div>
        </section>

        <section className="form-section">
          <div className="form-section-title">
            <span>05</span>
            <div><h2>Связь с покупателями</h2><p>Настройте удобный способ ответа.</p></div>
          </div>
          <div className="form-two-columns">
            <div className="form-control">
              <label htmlFor="listing-phone">Телефон</label>
              <Input
                id="listing-phone"
                value={draft.phone}
                onChange={(event) => setField("phone", event.target.value)}
              />
            </div>
            <div className="form-control">
              <label htmlFor="call-device">Принимать звонки</label>
              <Select
                id="call-device"
                value={draft.callDevice}
                onChange={(value) => setField("callDevice", value)}
                options={[
                  { value: "phone", label: "На телефоне" },
                  { value: "browser", label: "В браузере" },
                ]}
              />
            </div>
          </div>
          <div className="form-control">
            <label>Способ связи</label>
            <Radio.Group
              value={draft.contactMethod}
              onChange={(event) => setField("contactMethod", event.target.value)}
            >
              <Radio value="calls-and-messages">Звонки и сообщения</Radio>
              <Radio value="calls">Только звонки</Radio>
              <Radio value="messages">Только сообщения</Radio>
            </Radio.Group>
          </div>
        </section>

        <div className="publish-panel">
          <div><strong>Почти готово</strong><span>Проверьте обязательные поля перед публикацией.</span></div>
          <div>
            <Button size="large" onClick={() => go("/")}>Сохранить и выйти</Button>
            <Button
              type="primary"
              size="large"
              data-onboarding-id="publish-listing"
              data-onboarding-page="/add-item/details"
              data-onboarding-label="Кнопка «Разместить»"
              onClick={publish}
            >
              Разместить
            </Button>
          </div>
        </div>
      </div>
    </FlowLayout>
  );
}
