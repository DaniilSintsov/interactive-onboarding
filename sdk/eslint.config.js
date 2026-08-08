import js from "@eslint/js";
import tseslint from "typescript-eslint";

export default [
  { ignores: ["dist"] },
  { ...js.configs.recommended, files: ["src/**/*.ts"] },
  ...tseslint.configs.recommended.map((config) => ({
    ...config,
    files: ["src/**/*.ts"],
  })),
];
