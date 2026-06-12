# Repository Guidelines

## Project Structure & Module Organization

This repository is a compact Go module for `github.com/beevik/etree`, a pure-Go XML element tree library. Source files live at the repository root:

- `etree.go`: core document, element, token, read/write, and indentation APIs.
- `path.go`: XPath-like path parsing and query support.
- `helpers.go`: shared helper logic.
- `*_test.go`: unit and example tests, including `etree_test.go`, `path_test.go`, and `example_test.go`.
- `README.md` and `RELEASE_NOTES.md`: user documentation and release history.
- `.github/workflows/go.yml`: CI build, test, and CodeQL analysis.

## Build, Test, and Development Commands

- `go build -v ./...`: compile all packages; mirrors the CI build step.
- `go test -v ./...`: run the full test suite with verbose output; mirrors CI.
- `go test -run TestPath ./...`: run a focused test by name while iterating.
- `gofmt -w *.go`: format root Go files before committing.

The module targets Go `1.23.0` and CI currently checks Go `1.23` and `1.25.x`.

## Coding Style & Naming Conventions

Use standard Go formatting and idioms. Keep indentation and alignment to `gofmt`; do not hand-align with spaces. Public exported APIs should use clear Go names such as `Document`, `Element`, or `ReadSettings`, with comments that begin with the exported identifier. Prefer small helpers near related tests or implementation code instead of adding new packages unless the module structure truly needs it.

## Testing Guidelines

Tests use Go’s standard `testing` package. Add or update `TestXxx` functions in the relevant `*_test.go` file and keep examples as `ExampleXxx` functions when they document public behavior. For parser, XML serialization, path query, or compatibility changes, include focused regression coverage and run `go test -v ./...` before opening a PR.

## Commit & Pull Request Guidelines

Recent commits use short, imperative or descriptive messages, for example `Update README`, `Add iterator versions of element queries`, and `bug fix: InsertChildAt wrong index handling`. Keep commits focused on one logical change.

Pull requests should include a concise summary, the motivation or linked issue, and the tests run. For public API changes, update `README.md` or examples when needed and mention compatibility considerations.

## Security & Configuration Tips

This library has no runtime service configuration or external dependencies. Avoid adding new dependencies without a clear need, and be careful with XML parsing behavior because it can affect validation, entity handling, and compatibility with Go’s `encoding/xml`.
