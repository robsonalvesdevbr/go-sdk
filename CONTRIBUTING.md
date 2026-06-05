# Contributing to go-sdk

Thanks for your interest in improving **go-sdk**! This project is open source and free to modify,
and contributions and suggestions are very welcome.

## Ways to contribute

- **Report a bug** — open an issue describing what happened, what you expected, and the steps to
  reproduce it (include your OS, architecture and `go version` output).
- **Suggest a feature** — open an issue explaining the use case and the behavior you'd like.
- **Send a pull request** — fix a bug, add a feature, or improve the docs.

## Development setup

Requirements: **Go 1.26+** on **Linux/amd64** (see the [README](./README.md#requirements) for the
full list).

```bash
git clone https://github.com/robsonalvesdevbr/go-sdk.git
cd go-sdk

# Run the CLI without building
go run ./cmd/go-sdk --help

# Build a binary
go build -o go-sdk ./cmd/go-sdk
```

> Tip: `install`-related code writes to `/usr/local`. Be careful when testing it manually — prefer a
> disposable VM or container.

## Before opening a pull request

Run the standard Go checks and make sure they pass:

```bash
go build ./...
go vet ./...
go test ./...
gofmt -l .   # should print nothing; run `gofmt -w .` to fix formatting
```

Keep changes focused and small when possible, and update the documentation (README) when you change
user-facing behavior.

## Commit messages

This project follows the [Conventional Commits](https://www.conventionalcommits.org) specification.
Use a `type(scope): description` format, for example:

```text
feat(install): support installing a specific patch version
fix(list): mark the current version correctly
docs(readme): document the GITHUB_TOKEN environment variable
```

Common types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`.

## Pull request checklist

- [ ] The branch builds (`go build ./...`) and is formatted (`gofmt -l .` is clean).
- [ ] `go vet ./...` and `go test ./...` pass.
- [ ] Commits follow Conventional Commits.
- [ ] Documentation is updated when behavior changes.
- [ ] The PR description explains the motivation and the change.

## Code of conduct

Please be respectful and constructive in all interactions. We want this to be a welcoming project
for contributors of all experience levels.
