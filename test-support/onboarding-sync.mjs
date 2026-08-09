import assert from "node:assert/strict";
import fs from "node:fs";
import path from "node:path";
import ts from "typescript";

const root = process.cwd();
const defaultSource = path.join(root, "test-preview");
const keyPattern = /^[a-z0-9]+(?:-[a-z0-9]+)*$/;
const catalogMetadataPrefix = "onboarding-catalog:";

function catalogDescription(pagePaths) {
  return `${catalogMetadataPrefix}${JSON.stringify({ page_paths: pagePaths })}`;
}

function parseArgs(argv) {
  const options = {};
  for (let index = 0; index < argv.length; index += 1) {
    const argument = argv[index];
    if (argument === "--self-test") options.selfTest = true;
    else if (argument === "--skip-api") options.skipApi = true;
    else if (argument.startsWith("--backend-url=")) options.backendUrl = argument.slice(14);
    else if (argument === "--backend-url") options.backendUrl = argv[++index];
    else if (argument.startsWith("--project-key=")) options.projectKey = argument.slice(14);
    else if (argument === "--project-key") options.projectKey = argv[++index];
    else if (argument.startsWith("--source=")) options.source = argument.slice(9);
    else if (argument === "--source") options.source = argv[++index];
    else if (argument.startsWith("--manifest=")) options.manifest = argument.slice(11);
    else if (argument === "--manifest") options.manifest = argv[++index];
    else if (argument.startsWith("--revision=")) options.revision = argument.slice(11);
    else if (argument === "--revision") options.revision = argv[++index];
    else throw new Error(`unknown argument: ${argument}`);
  }
  return options;
}

function sourceFiles(directory) {
  const result = [];
  for (const entry of fs.readdirSync(directory, { withFileTypes: true })) {
    if (entry.name === "node_modules" || entry.name === ".next") continue;
    const entryPath = path.join(directory, entry.name);
    if (entry.isDirectory()) result.push(...sourceFiles(entryPath));
    else if (/\.[jt]sx$/.test(entry.name)) result.push(entryPath);
  }
  return result.sort();
}

function literalValues(expression) {
  if (!expression) return [];
  if (ts.isStringLiteralLike(expression)) return [expression.text];
  if (ts.isConditionalExpression(expression)) {
    return [...literalValues(expression.whenTrue), ...literalValues(expression.whenFalse)];
  }
  if (
    ts.isParenthesizedExpression(expression) ||
    ts.isAsExpression(expression) ||
    ts.isSatisfiesExpression(expression) ||
    ts.isNonNullExpression(expression)
  ) {
    return literalValues(expression.expression);
  }
  if (
    ts.isBinaryExpression(expression) &&
    [ts.SyntaxKind.AmpersandAmpersandToken, ts.SyntaxKind.QuestionQuestionToken].includes(
      expression.operatorToken.kind,
    )
  ) {
    return literalValues(expression.right);
  }
  if (
    expression.kind === ts.SyntaxKind.NullKeyword ||
    (ts.isIdentifier(expression) && expression.text === "undefined")
  ) {
    return [];
  }
  throw new Error(`unsupported onboarding attribute expression: ${expression.getText()}`);
}

function attributeValues(opening, name) {
  const attribute = opening.attributes.properties.find(
    (property) => ts.isJsxAttribute(property) && property.name.getText() === name,
  );
  if (!attribute || !ts.isJsxAttribute(attribute) || !attribute.initializer) return [];
  if (ts.isStringLiteral(attribute.initializer)) return [attribute.initializer.text];
  return literalValues(attribute.initializer.expression);
}

function requiredSingleAttribute(opening, name, filePath) {
  const values = [...new Set(attributeValues(opening, name).map((value) => value.trim()).filter(Boolean))];
  if (values.length !== 1) {
    throw new Error(`${filePath}: ${name} must contain one literal value`);
  }
  return values[0];
}

function collectSource(filePath, content, catalog, occurrences) {
  const source = ts.createSourceFile(filePath, content, ts.ScriptTarget.Latest, true, ts.ScriptKind.TSX);
  const visit = (node) => {
    if (ts.isJsxOpeningElement(node) || ts.isJsxSelfClosingElement(node)) {
      const keys = [...new Set(attributeValues(node, "data-onboarding-id").map((key) => key.trim()).filter(Boolean))];
      if (keys.length > 0) {
        const page = requiredSingleAttribute(node, "data-onboarding-page", filePath);
        const label = requiredSingleAttribute(node, "data-onboarding-label", filePath);
        if (!page.startsWith("/")) throw new Error(`${filePath}: onboarding page must start with /`);
        if (page.length > 2048) throw new Error(`${filePath}: onboarding page is longer than 2048 characters`);
        if (label.length > 255) throw new Error(`${filePath}: onboarding label is longer than 255 characters`);

        for (const key of keys) {
          if (!keyPattern.test(key)) throw new Error(`${filePath}: invalid onboarding key ${JSON.stringify(key)}`);
          if (key.length > 255) throw new Error(`${filePath}: onboarding key is longer than 255 characters`);
          const occurrence = `${key}\0${page}`;
          if (occurrences.has(occurrence)) {
            throw new Error(`duplicate data-onboarding-id ${JSON.stringify(key)} on ${page}`);
          }
          occurrences.add(occurrence);

          const current = catalog.get(key);
          if (current && current.label !== label) {
            throw new Error(`conflicting labels for data-onboarding-id ${JSON.stringify(key)}`);
          }
          if (current) current.pagePaths.add(page);
          else catalog.set(key, { key, label, pagePaths: new Set([page]) });
        }
      }
    }
    ts.forEachChild(node, visit);
  };
  visit(source);
}

function collectManifest(sourceDirectory, revision) {
  const catalog = new Map();
  const occurrences = new Set();
  for (const filePath of sourceFiles(sourceDirectory)) {
    collectSource(filePath, fs.readFileSync(filePath, "utf8"), catalog, occurrences);
  }
  if (catalog.size === 0) throw new Error(`no data-onboarding-id found in ${sourceDirectory}`);

  return {
    revision,
    elements: [...catalog.values()]
      .map(({ key, label, pagePaths }) => {
        const page_paths = [...pagePaths].sort();
        const description = catalogDescription(page_paths);
        if (description.length > 2000) throw new Error(`onboarding metadata for ${JSON.stringify(key)} is too long`);
        return {
          key,
          label,
          description,
          page_paths,
        };
      })
      .sort((left, right) => left.key.localeCompare(right.key)),
  };
}

function normalizeBackendUrl(value) {
  const normalized = String(value ?? "").trim().replace(/\/+$/, "");
  if (!normalized) throw new Error("--backend-url is required");
  const url = new URL(normalized);
  if (url.protocol !== "http:" && url.protocol !== "https:") {
    throw new Error("--backend-url must use http or https");
  }
  return url.toString().replace(/\/+$/, "");
}

function elementPayload(element) {
  return { key: element.key, label: element.label, description: element.description };
}

async function request(backendUrl, requestPath, init) {
  const response = await fetch(`${backendUrl}/api/v1${requestPath}`, {
    ...init,
    headers: {
      ...(init?.body ? { "content-type": "application/json" } : {}),
      ...init?.headers,
    },
  });
  if (!response.ok) {
    const body = await response.text();
    throw new Error(`${init?.method ?? "GET"} ${requestPath} failed: ${response.status} ${body}`);
  }
  return response.status === 204 ? undefined : response.json();
}

async function findProject(backendUrl, projectKey) {
  let offset = 0;
  while (true) {
    const page = await request(backendUrl, `/projects?limit=100&offset=${offset}`);
    const project = page.items.find((item) => item.project_key === projectKey);
    if (project) return project;
    offset += page.items.length;
    if (offset >= page.total || page.items.length === 0) break;
  }
  throw new Error(`project not found for project_key ${JSON.stringify(projectKey)}`);
}

function planSync(existing, desired) {
  const desiredByKey = new Map(desired.map((element) => [element.key, element]));
  const existingByKey = new Map(existing.map((element) => [element.key, element]));
  return {
    create: desired.filter((element) => !existingByKey.has(element.key)),
    update: desired.filter((element) => {
      const current = existingByKey.get(element.key);
      return current && (current.label !== element.label || current.description !== element.description);
    }),
    stale: existing.filter((element) => !desiredByKey.has(element.key)),
  };
}

async function syncElements(backendUrl, projectKey, manifest) {
  const project = await findProject(backendUrl, projectKey);
  const existing = await request(backendUrl, `/projects/${project.id}/elements`);
  const plan = planSync(existing, manifest.elements);

  for (const element of plan.create) {
    await request(backendUrl, `/projects/${project.id}/elements`, {
      method: "POST",
      body: JSON.stringify(elementPayload(element)),
    });
  }
  const existingByKey = new Map(existing.map((element) => [element.key, element]));
  for (const element of plan.update) {
    await request(backendUrl, `/projects/${project.id}/elements/${existingByKey.get(element.key).id}`, {
      method: "PATCH",
      body: JSON.stringify(elementPayload(element)),
    });
  }
  for (const element of plan.stale) {
    const stale = { ...element, description: catalogDescription([]) };
    if (element.description === stale.description) continue;
    await request(backendUrl, `/projects/${project.id}/elements/${element.id}`, {
      method: "PATCH",
      body: JSON.stringify(elementPayload(stale)),
    });
  }
  return { projectId: project.id, created: plan.create.length, updated: plan.update.length, stale: plan.stale.length };
}

function writeManifest(filePath, manifest) {
  fs.mkdirSync(path.dirname(filePath), { recursive: true });
  fs.writeFileSync(filePath, `${JSON.stringify(manifest, null, 2)}\n`);
}

function selfTest() {
  const source = `
    <button
      data-onboarding-id={active ? "save-item" : undefined}
      data-onboarding-page="/items/new"
      data-onboarding-label="Сохранить"
    />`;
  const catalog = new Map();
  collectSource("fixture.tsx", source, catalog, new Set());
  assert.deepEqual([...catalog.keys()], ["save-item"]);
  assert.deepEqual([...catalog.get("save-item").pagePaths], ["/items/new"]);
  assert.equal(catalogDescription(["/items/new"]), 'onboarding-catalog:{"page_paths":["/items/new"]}');
  assert.deepEqual(planSync([{ id: "1", key: "old", label: "Old", description: "" }], [
    { key: "new", label: "New", description: "", page_paths: ["/"] },
  ]), {
    create: [{ key: "new", label: "New", description: "", page_paths: ["/"] }],
    update: [],
    stale: [{ id: "1", key: "old", label: "Old", description: "" }],
  });
  assert.equal(normalizeBackendUrl("http://localhost:8080/"), "http://localhost:8080");
  console.log("onboarding-sync self-test passed");
}

async function main() {
  const options = parseArgs(process.argv.slice(2));
  if (options.selfTest) return selfTest();

  const source = path.resolve(root, options.source ?? defaultSource);
  const revision = String(options.revision ?? process.env.GITHUB_SHA ?? process.env.CI_COMMIT_SHA ?? "local").trim();
  const manifest = collectManifest(source, revision);
  console.log(`onboarding catalog: ${manifest.elements.length} elements, revision=${revision}`);
  if (options.manifest) {
    const manifestPath = path.resolve(root, options.manifest);
    writeManifest(manifestPath, manifest);
    console.log(`onboarding manifest: ${manifestPath}`);
  }

  if (options.skipApi) return;
  const projectKey = String(options.projectKey ?? "").trim();
  if (!projectKey) throw new Error("--project-key is required");
  const result = await syncElements(normalizeBackendUrl(options.backendUrl), projectKey, manifest);
  console.log(
    `onboarding sync: project=${result.projectId} created=${result.created} updated=${result.updated} stale=${result.stale}`,
  );
}

main().catch((error) => {
  console.error(error instanceof Error ? error.message : error);
  process.exit(1);
});
