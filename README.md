Personal workspace for reusable Go components and service scaffolding templates. Keeps common building blocks in one place so new services start quickly with familiar defaults.

## Goals

- Store frequently used packages (logging, response helpers, redis, http, jwt, oauth2, etc.) inside one repository.
- Provide a CLI that can scaffold a new service from curated templates instead of repeating boilerplate.

## Contents

- **CLI** (`cli/`): Cobra-based tool with template loader/renderer (`tmpl/`), bundled templates like `go-service`, `go-clean-service`.
- **Templates**: described via `template.yaml`, rendered using Go templates (`*.tmpl`); include build scripts, docker compose, env templates, etc.
- **Packages**: directories such as `logx`, `response`, `redis`, `http`, `jwt`, `password`, `helper`, `oauth2`, … hold reusable Go code.

## Direction

- Keep expanding templates and CLI utilities.
- Future ideas include template registry, interactive prompts, dry-run mode, etc.

