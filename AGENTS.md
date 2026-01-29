# Repository Guidelines

## Project Structure & Module Organization
- `backend/` contains the Go CLI/server, including `cmd/`, `internal/` (services, translator, server, config), and `pkg/` shared utilities.
- `web/` is the Vite + Vue 3 frontend (`src/` for app code, `src/views/`, `src/components/`, `src/locales/`).
- `backend/webui/` embeds the built frontend (`dist/`) for the backend binary.
- `tests/` holds CLI sample files and a shell test harness.
- Config lives in `backend/config.yaml`, `.env.sample`, and `CONFIGURATION.md`.

## Build, Test, and Development Commands
Run these from the repo root unless noted:
- `make -C backend binary` builds the Go binary into `backend/xcstrings-translator`.
- `make -C backend ui` installs web deps and builds the frontend into `backend/webui/dist`.
- `make -C backend test` runs Go tests in the backend module.
- `npm --prefix web run dev` starts the Vite dev server for the UI.
- `npm --prefix web run build` produces a production UI build.
- `bash tests/test.sh` runs the CLI smoke test script (requires `jq`).

## Coding Style & Naming Conventions
- Go: keep files gofmt’d; prefer `camelCase` for locals and `PascalCase` for exported symbols.
- Vue/TS: 2-space indentation, `PascalCase` for components, `camelCase` for composables and variables.
- Follow existing module boundaries (`internal/` for app-only packages; `pkg/` for reusable helpers).

## Testing Guidelines
- Backend uses Go’s `testing` package; place tests alongside source with `_test.go` suffix.
- Frontend linting is via `vue-tsc --noEmit` (`npm --prefix web run lint`).
- CLI smoke tests live in `tests/` and can be extended with additional `.xcstrings` fixtures.

## Commit & Pull Request Guidelines
- Commit messages follow Conventional Commits style seen in history: `feat: ...`, `fix: ...`, `refactor(scope): ...`.
- PRs should include a short summary, testing notes (commands run), and screenshots for UI changes.

## Configuration & Secrets
- Copy `.env.sample` to `.env` for local development; never commit secrets.
- Update `backend/config.yaml` for provider credentials and service settings.
