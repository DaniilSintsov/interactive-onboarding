"use client";

const categories = [
  { icon: "🚙", name: "Авто" },
  { icon: "🏠", name: "Недвижимость" },
  { icon: "🧳", name: "Жильё для путешествия" },
  { icon: "🪑", name: "Для дома и дачи" },
  { icon: "⚙️", name: "Запчасти" },
  { icon: "🛠️", name: "Услуги" },
  { icon: "📱", name: "Электроника" },
  { icon: "💼", name: "Работа" },
  { icon: "🏗️", name: "Для бизнеса" },
  { icon: "👟", name: "Одежда и обувь" },
];

const listings = [
  { icon: "📖", title: "Электронная книга PocketBook 628", price: "9 500 ₽", place: "Москва, Белорусская", tone: "sand" },
  { icon: "💻", title: "Чехол-магнит для iPad mini", price: "1 615 ₽", place: "Москва", tone: "gray" },
  { icon: "💡", title: "Настольная лампа, винтаж", price: "4 800 ₽", place: "Химки", tone: "green" },
  { icon: "🚲", title: "Городской велосипед", price: "24 000 ₽", place: "Москва, Сокол", tone: "blue" },
  { icon: "🪑", title: "Кресло из массива бука", price: "12 900 ₽", place: "Москва", tone: "coral" },
  { icon: "📷", title: "Плёночный фотоаппарат Zenit", price: "8 500 ₽", place: "Москва, Арбатская", tone: "dark" },
  { icon: "🪴", title: "Монстера в керамическом кашпо", price: "1 300 ₽", place: "Долгопрудный", tone: "mint" },
  { icon: "🎧", title: "Беспроводные наушники", price: "6 900 ₽", place: "Москва", tone: "lilac" },
];

export function HomeView() {
  return (
    <div className="home-page page-enter">
      <section className="market-categories" aria-label="Категории объявлений">
        <div className="category-grid">
          {categories.map((category) => (
            <button className="category-tile" key={category.name} type="button">
              <strong>{category.name}</strong>
              <span aria-hidden="true">{category.icon}</span>
            </button>
          ))}
        </div>
      </section>

      <section className="listing-section" aria-labelledby="listings-title">
        <h1 id="listings-title">Рекомендации для вас</h1>
        <div className="listing-grid">
          {listings.map((listing) => (
            <article className="listing-card" key={listing.title}>
              <div className={`listing-image listing-${listing.tone}`}>
                <span aria-hidden="true">{listing.icon}</span>
                <button type="button" aria-label={`Добавить «${listing.title}» в избранное`}>♡</button>
              </div>
              <h2>{listing.title}</h2>
              <strong>{listing.price}</strong>
              <p>{listing.place}</p>
              <small>Сегодня</small>
            </article>
          ))}
        </div>
      </section>
    </div>
  );
}
