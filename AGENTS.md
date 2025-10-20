# Repository Guidelines

## Project Structure & Module Organization
The project targets Go 1.23+ with a standard layout. `cmd/downloader` hosts the CLI entrypoint, while `cmd/gui` wraps the GUI built with Fyne; both depend on reusable packages under `pkg/`. Key subpackages include `pkg/core` for orchestration, `pkg/sites` for site-specific scrapers, and `pkg/util` for image and filesystem helpers. Supporting flag parsing and other internals live in `internal/`. Docs for contributors live in `docs/`, assets in `img/`, and release artifacts appear under `build/` when generated.

## Build, Test, and Development Commands
Run `go build -o ./bin/comics-downloader ./cmd/downloader` for a local CLI binary. The GUI can be built with `go build -o ./bin/comics-downloader-gui ./cmd/gui` once Fyne prerequisites are met. Cross-platform artifacts are automated via Make targets such as `make linux-x86-64-build` or the aggregate `make builds`. Use `fyne-cross windows -output comics-downloader-gui-windows.exe ./cmd/gui` when Docker-based cross-compilation is preferred. Static checks run via `make lint` (wraps `golangci-lint run` and matches CI), and `go test -v ./...` should pass before pushing; CI mirrors this command set and reports coverage with Coveralls.

## Coding Style & Naming Conventions
Format all Go files with `gofmt` (invoked via `go fmt ./...`) to preserve canonical tab-indented style and import ordering. Follow idiomatic Go naming: exported identifiers use CamelCase, unexported use camelCase, constants are ALL_CAPS only when enumerations require clarity. Keep packages small and cohesive; prefer descriptive filenames like `mangadex.go` aligning with the site handled. Log through the existing Logrus helpers to stay consistent with current output.

## Testing Guidelines
Place tests alongside implementation in `*_test.go` files, using table-driven cases where practical. Run targeted suites with `go test ./pkg/sites` during feature work, and use `go test -cover ./...` to ensure coverage does not regress—Coveralls tracks the master branch. When scraping new sources, add integration-style tests that stub HTTP calls where possible to keep the suite deterministic.

## Commit & Pull Request Guidelines
Commit subjects should be concise and imperative (`add`, `fix`, `update`), optionally prefixed with scopes like `fix:` as seen in history. Commit early and keep diffs focused; avoid mixing GUI and CLI changes in one commit. Pull requests must target `master`, include a brief summary of the change, reference related issues, and note how to reproduce or test the update (logs, commands, screenshots for GUI tweaks). Ensure tests pass locally before requesting review.
