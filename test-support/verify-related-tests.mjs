import assert from "node:assert/strict";
import path from "node:path";
import { spawnSync } from "node:child_process";

const root = process.cwd();

function run(command, args, capture = false) {
  const result = spawnSync(command, args, {
    cwd: root,
    encoding: capture ? "utf8" : undefined,
    env: process.env,
    stdio: capture ? ["ignore", "pipe", "inherit"] : "inherit",
  });
  if (result.error) throw result.error;
  if (result.status !== 0) process.exit(result.status ?? 1);
  return result.stdout ?? "";
}

function changedFiles() {
  const value = process.env.HARNEST_CHANGED_FILES;
  if (!value) throw new Error("HARNEST_CHANGED_FILES is required; run harnest verify --changed");
  return value.split(/\r?\n/).filter(Boolean).map((file) => file.replaceAll("\\", "/"));
}

function isWorkspaceConfig(file) {
  return /\/(package\.json|tsconfig(?:\.[^/]+)?\.json|vitest\.config\.[^/]+)$/.test(file);
}

function runWorkspaceTest(workspace) {
  run("npm", ["test", "-w", workspace]);
}

function runFrontend(files) {
  if (files.some((file) => file === "package.json" || file === "package-lock.json")) {
    run("npm", ["test"]);
    return;
  }

  const full = new Set();
  if (files.some((file) => file.startsWith("sdk/"))) {
    runWorkspaceTest("@interactive-onboarding/sdk");
    runWorkspaceTest("@interactive-onboarding/test-preview");
    full.add("test-preview");
  }

  for (const [directory, workspace] of [
    ["admin", "@interactive-onboarding/admin"],
    ["test-preview", "@interactive-onboarding/test-preview"],
  ]) {
    const related = files.filter((file) => file.startsWith(`${directory}/`));
    if (related.length === 0 || full.has(directory)) continue;
    if (related.some(isWorkspaceConfig)) {
      runWorkspaceTest(workspace);
      continue;
    }
    run("npm", [
      "exec", "--workspace", workspace, "--", "vitest", "related",
      "--run", "--passWithNoTests", ...related.map((file) => file.slice(directory.length + 1)),
    ]);
  }
}

function parsePackages(output) {
  return output.trim().split("\n").filter(Boolean).map((line) => {
    const [importPath, directory, deps = "", tests = "", externalTests = ""] = line.split("\t");
    return {
      importPath,
      directory: path.resolve(directory),
      references: new Set([...deps.split(","), ...tests.split(","), ...externalTests.split(",")].filter(Boolean)),
    };
  });
}

function selectImpactedPackages(packages, changedImports) {
  return packages
    .filter((pkg) => changedImports.has(pkg.importPath) || [...changedImports].some((name) => pkg.references.has(name)))
    .map((pkg) => pkg.importPath)
    .sort();
}

function runBackend(files) {
  const backendFiles = files.filter((file) => file.startsWith("backend/"));
  if (backendFiles.some((file) => file === "backend/go.mod" || file === "backend/go.sum")) {
    run("go", ["-C", "backend", "test", "-mod=readonly", "./..."]);
    return;
  }

  const goFiles = backendFiles.filter((file) => file.endsWith(".go"));
  if (goFiles.length === 0) return;

  const template = '{{.ImportPath}}\t{{.Dir}}\t{{join .Deps ","}}\t{{join .TestImports ","}}\t{{join .XTestImports ","}}';
  const packages = parsePackages(run("go", ["-C", "backend", "list", "-mod=readonly", "-f", template, "./..."], true));
  const changedDirectories = new Set(goFiles.map((file) => path.dirname(path.resolve(root, file))));
  const changedImports = new Set(packages.filter((pkg) => changedDirectories.has(pkg.directory)).map((pkg) => pkg.importPath));
  const targets = new Set(selectImpactedPackages(packages, changedImports));

  for (const directory of changedDirectories) {
    if ([...packages].some((pkg) => pkg.directory === directory)) continue;
    targets.add(`./${path.relative(path.join(root, "backend"), directory).replaceAll("\\", "/")}`);
  }
  if (targets.size > 0) run("go", ["-C", "backend", "test", "-mod=readonly", ...[...targets].sort()]);
}

function selfTest() {
  const packages = [
    { importPath: "example/a", references: new Set() },
    { importPath: "example/b", references: new Set(["example/a"]) },
    { importPath: "example/c", references: new Set(["other"]) },
  ];
  assert.deepEqual(selectImpactedPackages(packages, new Set(["example/a"])), ["example/a", "example/b"]);
  assert.equal(isWorkspaceConfig("admin/vitest.config.ts"), true);
  assert.equal(isWorkspaceConfig("admin/src/page.tsx"), false);
  console.log("related-test runner self-test passed");
}

const mode = process.argv[2];
if (mode === "--self-test") selfTest();
else if (mode === "frontend" || mode === "backend") {
  const files = changedFiles();
  if (files.includes("test-support/verify-related-tests.mjs")) selfTest();
  if (mode === "frontend") runFrontend(files);
  else runBackend(files);
} else throw new Error("usage: verify-related-tests.mjs frontend|backend|--self-test");
