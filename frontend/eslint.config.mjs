import { defineConfig, globalIgnores } from "eslint/config";
import nextPlugin from "@next/eslint-plugin-next";

// eslint-config-next / eslint-plugin-react are not ESLint 10 + TypeScript 7 ready yet.
// Lint with the Next plugin; type-check covers TypeScript via tsc.
const eslintConfig = defineConfig([
  nextPlugin.configs["core-web-vitals"],
  {
    rules: {
      "@next/next/no-img-element": "warn",
    },
  },
  globalIgnores([
    ".next/**",
    "out/**",
    "build/**",
    "next-env.d.ts",
    "node_modules/**",
  ]),
]);

export default eslintConfig;
