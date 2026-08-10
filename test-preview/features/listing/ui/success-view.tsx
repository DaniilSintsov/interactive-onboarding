"use client";

import { useEffect, useRef, useState } from "react";
import { Button } from "antd";

import { useOnboarding } from "@/features/onboarding/ui/onboarding-provider";

import { useListingFlow } from "../controller/use-listing-flow";
import { useDraftStore } from "../model/draft-store";

export function SuccessView() {
  const resetDraft = useDraftStore((state) => state.resetDraft);
  const { completeCurrentStep } = useOnboarding();
  const { go, startOver } = useListingFlow();
  const [savedTitle] = useState(() => useDraftStore.getState().title);
  const [trackingDone, setTrackingDone] = useState(false);
  const started = useRef(false);

  useEffect(() => {
    if (started.current) return;
    started.current = true;
    void completeCurrentStep()
      .catch(() => undefined)
      .finally(() => {
        resetDraft();
        setTrackingDone(true);
      });
  }, [completeCurrentStep, resetDraft]);

  return (
    <div className="success-page page-enter">
      <div className="success-confetti" aria-hidden="true">
        <i /><i /><i /><i /><i />
      </div>
      <section className="success-card">
        <div className="success-check">✓</div>
        <span className="section-kicker">Объявление опубликовано</span>
        <h1>Готово — теперь его увидят покупатели</h1>
        <p>
          «{savedTitle || "Ваше объявление"}» прошло проверку и уже появилось в
          поиске. Первые просмотры обычно приходят в течение часа.
        </p>
        <div className="success-stats">
          <div><strong>0</strong><span>просмотров</span></div>
          <div><strong>30 дней</strong><span>срок размещения</span></div>
          <div><strong>{trackingDone ? "готово" : "…"}</strong><span>статус онбординга</span></div>
        </div>
        <div className="success-actions">
          <Button type="primary" size="large" onClick={() => go("/")}>На главную</Button>
          <Button size="large" onClick={startOver}>Разместить ещё одно</Button>
        </div>
      </section>
      <aside className="success-tip">
        <span>Совет</span>
        Отвечайте на сообщения быстро — активные продавцы получают больше откликов.
      </aside>
    </div>
  );
}
