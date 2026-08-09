# interactive-onboarding

## CI-синхронизация элементов

Элементы test-preview участвуют в онбординге, если у них заданы статические атрибуты:

```tsx
<button
  data-onboarding-id="publish-listing"
  data-onboarding-page="/add-item/details"
  data-onboarding-label="Кнопка «Разместить»"
/>
```

После установки зависимостей deploy job запускает:

```bash
npm run onboarding:sync -- \
  --backend-url="$ONBOARDING_BACKEND_URL" \
  --project-key="$ONBOARDING_PROJECT_KEY" \
  --revision="${CI_COMMIT_SHA:-local}"
```

Токен не нужен: команда рассчитана на внутренний backend URL, доступный только из доверенного test-preview CI/CD. Скрипт создаёт и обновляет элементы проекта, сохраняет привязку к страницам и помечает исчезнувшие ключи недоступными, не удаляя используемые сценариями записи. `npm run onboarding:sync:test` проверяет парсер без обращения к API.

`docker compose up --build` запускает эту команду автоматически в one-shot сервисе `onboarding-sync`. Backend URL, project key и `${CI_COMMIT_SHA:-local}` передаются из Compose; `test-preview` стартует только после успешной синхронизации. Проект с указанным `NEXT_PUBLIC_PROJECT_KEY` должен уже существовать; Admin стартует независимо, чтобы его можно было создать на чистой БД.
