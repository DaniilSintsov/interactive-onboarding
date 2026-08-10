"use client";

export function FlowLayout({
  current,
  title,
  lead,
  children,
  aside,
}: {
  current: 1 | 2 | 3;
  title: string;
  lead: string;
  children: React.ReactNode;
  aside: React.ReactNode;
}) {
  return (
    <div className="flow-page page-enter" aria-label={`Шаг ${current} из 3`}>
      <div className="flow-grid">
        <section className="flow-main">
          <header className="flow-heading">
            <span>Новое объявление</span>
            <h1>{title}</h1>
            <p>{lead}</p>
          </header>
          {children}
        </section>
        <aside className="flow-aside">{aside}</aside>
      </div>
    </div>
  );
}
