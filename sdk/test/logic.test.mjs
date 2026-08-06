import assert from "node:assert/strict";
import test from "node:test";

import { hasMeaningfulValue, parseRuntimeScenario } from "../dist/logic.js";

test("runtime data is validated, ordered, and change values are meaningful", () => {
  const step = (id, stepNum, advance) => ({
    id,
    step_num: stepNum,
    title: `Шаг ${stepNum}`,
    description: "Описание",
    frontend_data: { page_path: "/add-item/details", advance },
    element: { id: `element-${id}`, key: id, label: id, description: "" },
  });
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
