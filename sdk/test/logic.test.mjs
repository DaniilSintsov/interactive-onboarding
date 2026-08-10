import assert from "node:assert/strict";
import test from "node:test";

import {
  hasMeaningfulValue,
  parseRuntimeScenario,
  parseRuntimeScenarioResolveResponse,
} from "../dist/logic.js";

const step = (id, stepNum, advance = { mode: "manual" }) => ({
  id,
  step_num: stepNum,
  title: `Шаг ${stepNum}`,
  description: "Описание",
  frontend_data: { page_path: "/add-item/details", advance },
  element: { id: `element-${id}`, key: id, label: id, description: "" },
});

const runtimeScenario = (id) => ({
  id,
  name: `Маршрут ${id}`,
  description: "Помощь по шагам",
  page_pattern: "/",
  steps: [step(`${id}-step`, 1)],
});

test("runtime data is validated, ordered, and change values are meaningful", () => {
  const scenario = parseRuntimeScenario({
    id: "scenario-1",
    name: "Разместить объявление",
    description: "Помощь по шагам",
    page_pattern: "/",
    steps: [
      step("second", 2, { mode: "manual" }),
      step("first", 1, { mode: "target_event", event: "change" }),
    ],
  });

  assert.deepEqual(scenario.steps.map(({ id }) => id), ["first", "second"]);
  assert.equal(hasMeaningfulValue({ value: "  " }), false);
  assert.equal(hasMeaningfulValue({ value: "Книга" }), true);
  assert.equal(hasMeaningfulValue({ type: "file", files: { length: 1 } }), true);
  assert.throws(
    () => parseRuntimeScenario({ ...scenario, steps: [step("bad", 1, { mode: "target_event", event: "blur" })] }),
    /advance is invalid/,
  );
});

test("resolve contract accepts zero, one, or many scenarios and validates the wrapper", () => {
  assert.deepEqual(parseRuntimeScenarioResolveResponse({ is_test: false, scenarios: [] }), {
    is_test: false,
    scenarios: [],
  });

  const one = parseRuntimeScenarioResolveResponse({
    is_test: true,
    scenarios: [runtimeScenario("one")],
  });
  assert.equal(one.is_test, true);
  assert.equal(one.scenarios[0].id, "one");

  const many = parseRuntimeScenarioResolveResponse({
    is_test: false,
    scenarios: [runtimeScenario("one"), runtimeScenario("two")],
  });
  assert.deepEqual(many.scenarios.map(({ id }) => id), ["one", "two"]);

  assert.throws(
    () => parseRuntimeScenarioResolveResponse({ is_test: "false", scenarios: [] }),
    /is_test must be a boolean/,
  );
  assert.throws(
    () => parseRuntimeScenarioResolveResponse({ is_test: false, scenarios: [{}] }),
    /scenario.steps must be a non-empty array/,
  );
});
