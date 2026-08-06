import { randomUUID } from "node:crypto";
import { createServer } from "node:http";

const port = Number(process.env.MOCK_API_PORT ?? 8080);
const now = "2026-08-06T12:00:00.000Z";
const projectId = "11111111-1111-4111-8111-111111111111";
const scenarioId = "22222222-2222-4222-8222-222222222222";

const elementDefinitions = [
  ["category-hobby", "Категория «Хобби и отдых»"],
  ["listing-title", "Название объявления"],
  ["subcategory-ebooks", "Подкатегория «Электронные книги»"],
  ["listing-photo", "Фотография товара"],
  ["listing-description", "Описание объявления"],
  ["listing-price", "Цена объявления"],
  ["publish-listing", "Кнопка «Разместить»"],
];

const elements = elementDefinitions.map(([key, label], index) => ({
  id: `33333333-3333-4333-8333-33333333333${index}`,
  project_id: projectId,
  key,
  label,
  description: "",
  created_at: now,
  updated_at: now,
}));

const stepDefinitions = [
  ["Выберите категорию", "Начните с «Хобби и отдых».", "/add-item/category", "click"],
  ["Введите название", "Напишите «Электронная книга».", "/add-item/title", "change"],
  ["Уточните подкатегорию", "Выберите «Электронные книги».", "/add-item/title", "click"],
  ["Добавьте фотографию", "Загрузите один снимок товара.", "/add-item/details", "change"],
  ["Опишите товар", "Расскажите о состоянии и комплекте.", "/add-item/details", "change"],
  ["Укажите цену", "Введите цену больше нуля.", "/add-item/details", "change"],
  ["Разместите объявление", "Проверьте поля и опубликуйте объявление.", "/add-item/details", "manual"],
];

const steps = stepDefinitions.map(([title, description, page_path, event], index) => ({
  id: `44444444-4444-4444-8444-44444444444${index}`,
  scenario_id: scenarioId,
  element_id: elements[index].id,
  step_num: index + 1,
  title,
  description,
  frontend_data: {
    page_path,
    advance: event === "manual" ? { mode: "manual" } : { mode: "target_event", event },
  },
  created_at: now,
  updated_at: now,
}));

const project = {
  id: projectId,
  name: "Тестовый классифайд",
  project_key: "pk_demo_avito",
  created_at: now,
  updated_at: now,
};

const scenario = {
  id: scenarioId,
  project_id: projectId,
  name: "Первое объявление",
  description: "Семь подсказок помогут разместить первую вещь.",
  page_pattern: "/",
  status: "enabled",
  published_at: now,
  created_at: now,
  updated_at: now,
};

const sessions = [];
const events = [];

function runtimeScenario() {
  return {
    id: scenario.id,
    name: scenario.name,
    description: scenario.description,
    page_pattern: scenario.page_pattern,
    steps: steps.map((step, index) => ({ ...step, element: elements[index] })),
  };
}

function send(response, status, body, contentType = "application/json") {
  response.writeHead(status, { "content-type": contentType });
  response.end(body === undefined ? undefined : contentType === "application/json" ? JSON.stringify(body) : body);
}

async function readJson(request) {
  const chunks = [];
  for await (const chunk of request) chunks.push(chunk);
  return chunks.length ? JSON.parse(Buffer.concat(chunks).toString("utf8")) : {};
}

createServer(async (request, response) => {
  const url = new URL(request.url ?? "/", `http://${request.headers.host}`);
  const { pathname } = url;

  if (request.method === "GET" && pathname === "/__events") {
    send(response, 200, { sessions, events });
    return;
  }
  if (request.method === "POST" && pathname === "/api/v1/sdk/scenarios/resolve") {
    const body = await readJson(request);
    if (body.user_id === "demo-expert" || body.page !== "/") send(response, 204);
    else send(response, 200, runtimeScenario());
    return;
  }
  if (request.method === "POST" && pathname === "/api/v1/sdk/sessions") {
    const body = await readJson(request);
    const session = {
      id: randomUUID(),
      scenario_id: body.scenario_id,
      user_id: body.user_id,
      status: "active",
      started_at: new Date().toISOString(),
    };
    sessions.push(session);
    send(response, 201, session);
    return;
  }
  if (request.method === "POST" && pathname === "/api/v1/sdk/events") {
    const body = await readJson(request);
    const event = { ...body, received_at: new Date().toISOString() };
    events.push(event);
    send(response, 202, { event, duplicate: false });
    return;
  }

  if (request.method === "GET" && pathname === "/api/v1/projects") {
    send(response, 200, { items: [project], total: 1, limit: 100, offset: 0 });
    return;
  }
  if (request.method === "POST" && pathname === "/api/v1/projects") {
    const body = await readJson(request);
    send(response, 201, { ...project, id: randomUUID(), name: body.name, elements: [] });
    return;
  }
  if (request.method === "GET" && pathname === `/api/v1/projects/${projectId}`) {
    send(response, 200, { ...project, elements });
    return;
  }
  if (request.method === "GET" && pathname === `/api/v1/projects/${projectId}/elements`) {
    send(response, 200, elements);
    return;
  }
  if (request.method === "GET" && pathname === `/api/v1/projects/${projectId}/scenarios`) {
    send(response, 200, {
      items: [{ ...scenario, steps_count: steps.length }],
      total: 1,
      limit: 100,
      offset: 0,
    });
    return;
  }
  if (request.method === "GET" && pathname === `/api/v1/scenarios/${scenarioId}`) {
    send(response, 200, { ...scenario, steps });
    return;
  }
  if (request.method === "GET" && pathname === `/api/v1/scenarios/${scenarioId}/preview`) {
    send(response, 200, runtimeScenario());
    return;
  }
  if (request.method === "GET" && pathname === `/api/v1/projects/${projectId}/analytics/total`) {
    const completed = events.filter((event) => event.type === "onboarding_completed").length;
    send(response, 200, {
      project_id: projectId,
      total_scenarios: 1,
      enabled_scenarios: 1,
      sessions_started: sessions.length,
      sessions_completed: completed,
      sessions_skipped: events.filter((event) => event.type === "onboarding_skipped").length,
      completion_rate: sessions.length ? completed / sessions.length : 0,
      skip_rate: 0,
    });
    return;
  }
  if (request.method === "GET" && pathname === `/api/v1/scenarios/${scenarioId}/analytics/detailed`) {
    send(response, 200, {
      scenario_id: scenarioId,
      started: sessions.length,
      completed: events.filter((event) => event.type === "onboarding_completed").length,
      skipped: events.filter((event) => event.type === "onboarding_skipped").length,
      completion_rate: 0,
      skip_rate: 0,
      average_completion_time_seconds: 0,
      steps: steps.map((step) => ({
        step_id: step.id,
        position: step.step_num,
        title: step.title,
        shown: events.filter((event) => event.type === "step_shown" && event.step_id === step.id).length,
        completed: events.filter((event) => event.type === "step_completed" && event.step_id === step.id).length,
        skipped: 0,
        completion_rate: 0,
        skip_rate: 0,
        drop_off_rate: 0,
      })),
    });
    return;
  }
  if (request.method === "GET" && pathname === `/api/v1/scenarios/${scenarioId}/report/pdf`) {
    send(response, 200, "%PDF-1.4\n%%EOF", "application/pdf");
    return;
  }

  send(response, 404, { code: "not_found", message: `${request.method} ${pathname}` });
}).listen(port, "127.0.0.1", () => {
  console.log(`Mock API listening on http://127.0.0.1:${port}`);
});
