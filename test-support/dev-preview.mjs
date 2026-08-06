import { spawn } from "node:child_process";

const npm = process.platform === "win32" ? "npm.cmd" : "npm";

function start(args) {
  return spawn(npm, args, { stdio: "inherit" });
}

function wait(child) {
  return new Promise((resolve) => child.once("exit", (code) => resolve(code ?? 1)));
}

const build = start(["run", "build", "-w", "@interactive-onboarding/sdk"]);
const buildCode = await wait(build);
if (buildCode !== 0) process.exit(buildCode);

const mock = start(["run", "mock-api"]);
const preview = start(["run", "dev", "-w", "@interactive-onboarding/test-preview"]);
let stopping = false;

function stop(code) {
  if (stopping) return;
  stopping = true;
  mock.kill("SIGTERM");
  preview.kill("SIGTERM");
  process.exitCode = code;
}

mock.once("exit", (code) => stop(code ?? 1));
preview.once("exit", (code) => stop(code ?? 0));
process.once("SIGINT", () => stop(130));
process.once("SIGTERM", () => stop(143));
