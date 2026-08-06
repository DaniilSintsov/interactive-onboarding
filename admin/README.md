# Admin

Защищённая Next.js-панель управления онбордингами. Браузер обращается только к
`/api/backend/*`; BFF добавляет Bearer и проксирует запрос в Go API.

```bash
cp .env.example .env.local
npm install
npm run dev
```

Переменные: `ADMIN_API_URL` (по умолчанию `http://localhost:8080`), общий
`ADMIN_PASSWORD`, `SESSION_SECRET` длиной от 32 символов и публичный
`NEXT_PUBLIC_PREVIEW_URL` тестового классифайда.
