import { LoginForm } from '@/features/auth/ui/login-form';

export default async function LoginPage({
  searchParams,
}: {
  searchParams: Promise<{ from?: string }>;
}) {
  const { from } = await searchParams;
  const returnTo = from?.startsWith('/') && !from.startsWith('//') ? from : '/projects';

  return (
    <main className="login-page">
      <section className="login-story" aria-label="О продукте">
        <div className="brand-lockup brand-lockup-light">
          <span className="brand-mark" aria-hidden="true">
            <i />
            <i />
            <i />
            <i />
          </span>
          <span>ONBOARD CONTROL</span>
        </div>
        <div className="login-story-copy">
          <p className="eyebrow">Продуктовая диспетчерская</p>
          <h1>Помощь приходит ровно в нужный момент.</h1>
          <p>
            Собирайте многостраничные сценарии, проверяйте их на тестовом классифайде и
            следите, где пользователям нужна поддержка.
          </p>
        </div>
        <div className="signal-line" aria-hidden="true">
          <span>ENTRY</span>
          <b />
          <span>GUIDE</span>
          <b />
          <span>RESULT</span>
        </div>
      </section>
      <section className="login-panel">
        <div className="login-card">
          <p className="eyebrow">Защищённая зона</p>
          <h2>Вход администратора</h2>
          <p className="muted">Пароль хранится только на сервере. Сессия действует 7 дней.</p>
          <LoginForm returnTo={returnTo} />
        </div>
      </section>
    </main>
  );
}
