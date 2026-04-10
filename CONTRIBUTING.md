# Contributing to Shapin

Thank you for your interest in contributing!

## Reporting issues

Open an issue on [GitHub](https://github.com/Kirskov/Shapin/issues) with:

- A clear description of the problem
- The command you ran and the full output
- The relevant CI file(s) if applicable

## Development setup

```sh
git clone https://github.com/Kirskov/Shapin.git
cd Shapin
go build -o shapin ./cmd/shapin
go test ./...
```

Requires Go 1.25+.

## Making changes

1. Fork the repository and create a branch from `main`
2. Make your changes
3. Run the tests: `go test ./...`
4. Open a pull request

## Commit style

This project uses [Conventional Commits](https://www.conventionalcommits.org):

```
feat: add support for new provider
fix: skip commented image lines
docs: update README
chore: update dependencies
refactor: extract helper function
test: add cases for prefix stem matching
```

## Adding a provider

1. Create `internal/providers/myprovider.go` implementing the `contract.Provider` interface
2. Register it in `internal/scanner/runner.go`
3. Add tests in `internal/providers/providers_test.go`
4. Document it in `README.md` under `## Providers`

## Adding a built-in stem mapping

Built-in mappings live in `internal/providers/util.go` in the `builtinStemMappings` map. Add the stem (uppercase) and the corresponding Docker Hub image path. Update the table in `README.md`.

## Code of Conduct

Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.
