# Admin

Защищённая Next.js-панель управления онбордингами. Браузер обращается только к
`/api/backend/*`; BFF проверяет admin-сессию и проксирует запрос в Go API по внутренней сети.

```bash
cp .env.example .env.local
npm install
npm run credentials
npm run dev
```

Сохраните `ADMIN_LOGIN_PASSWORD` из вывода генератора в менеджере паролей, а
`ADMIN_PASSWORD_HASH` и `SESSION_SECRET` — в `.env.local` или secret settings платформы.
Не добавляйте `ADMIN_LOGIN_PASSWORD` в env.

Переменные: `ADMIN_API_URL` (по умолчанию `http://localhost:8080`), versioned scrypt hash
`ADMIN_PASSWORD_HASH`, `SESSION_SECRET` и публичный `NEXT_PUBLIC_PREVIEW_URL` тестового
классифайда. Сессия действует 24 часа.
Повторный запуск генератора ротирует credentials и завершает ранее выданные сессии.
При rolling deployment сначала добавьте новые env, затем разверните новый image; удаляйте
старый `ADMIN_PASSWORD` только после остановки старых replicas.

В Docker Compose Go API доступен контейнерам по внутренней сети, а host-доступ для
локальной разработки привязан к `127.0.0.1:8080`. При deployment admin отдельно от
Go, включая Vercel, `ADMIN_API_URL` должен быть доступен из Node runtime безопасным
приватным каналом.
