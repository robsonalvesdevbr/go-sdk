# go-sdk

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.26%2B-00ADD8.svg)](go.mod)

> Manage the Golang SDK from your terminal — list, inspect and install Go versions.

`go-sdk` is a small command-line tool (built with [Cobra](https://github.com/spf13/cobra)) that helps
you manage the Go SDK installed on your machine. It can list the stable Go versions available
upstream, show which version is currently installed, and download/install a specific or the latest
version.

## Requirements

- **Go 1.26+** (to build/install the tool).
- **Linux / amd64** — the installer downloads `linux-amd64` archives and extracts them into
  `/usr/local` (both are currently hardcoded).
- **Write access to `/usr/local`** — installing a Go version usually requires running with `sudo`.
- **Internet access** to `go.dev` and `api.github.com`.
- **`GITHUB_TOKEN`** (optional) — increases the GitHub API rate limit used to list versions.

## Installation

Install the latest published version:

```bash
go install github.com/robsonalvesdevbr/go-sdk/cmd/go-sdk@latest
```

Or build it locally from a clone:

```bash
git clone https://github.com/robsonalvesdevbr/go-sdk.git
cd go-sdk

# Build a binary
go build -o go-sdk ./cmd/go-sdk

# ...or run it directly without building
go run ./cmd/go-sdk
```

## Running with sudo

Installing a Go version into `/usr/local` requires root, but `sudo` runs with the root user's
`secure_path` (set in `/etc/sudoers`), which usually does **not** include your `go` toolchain or
`$GOPATH/bin`. That is why commands like `sudo go run ...` or `sudo go-sdk ...` may fail with
`'go': command not found` or `'go-sdk': command not found`.

Use one of the following instead:

```bash
# 1. Build a self-contained binary and run that binary with sudo (recommended)
go build -o go-sdk ./cmd/go-sdk
sudo ./go-sdk install -l

# 2. After `go install`, call the binary by its absolute path
go install github.com/robsonalvesdevbr/go-sdk/cmd/go-sdk@latest
sudo "$(go env GOPATH)/bin/go-sdk" install -l

# 3. Preserve your PATH when elevating (keeps `go` / `go-sdk` resolvable)
sudo env "PATH=$PATH" go run ./cmd/go-sdk install -l
sudo env "PATH=$PATH" go-sdk install -l
```

> Tip: if you'd rather avoid `sudo` entirely, install into a writable directory with
> `--dir` / `GO_SDK_INSTALL_DIR` (see [`install`](#install)).

## Usage

```text
Manage Golang SDK

Usage:
  go-sdk [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  current     Current Go version
  help        Help about any command
  install     Install a specific version of Go
  list        List stable Go versions

Flags:
  -h, --help   help for go-sdk

Use "go-sdk [command] --help" for more information about a command.
```

| Command      | Description                                              |
| ------------ | -------------------------------------------------------- |
| `list`       | List the stable Go versions available upstream           |
| `current`    | Show the Go version currently installed on the system    |
| `install`    | Download and install a specific or the latest Go version |
| `completion` | Generate the shell autocompletion script                 |
| `help`       | Help about any command                                   |

## Commands

### `list`

Lists all stable Go versions (fetched from the `golang/go` tags on the GitHub API). The version
currently installed on the system is marked with `(current)`.

```bash
go-sdk list
```

```text
go1.23.4
go1.24.0
go1.24.1 (current)
go1.24.2
```

### `current`

Shows the Go version reported by `go version` on the host.

```bash
go-sdk current
```

```text
System Go version: go version go1.24.1 linux/amd64
```

Use `-l/--local` to also print the path of the `go` binary in use:

```bash
go-sdk current -l
```

```text
System Go version: go version go1.24.1 linux/amd64
Local Go version: /usr/local/go/bin/go
```

### `install`

Downloads a Go archive from `https://go.dev/dl/` and extracts it into `/usr/local`.

Install the latest version:

```bash
sudo go-sdk install --latest
# or
sudo go-sdk install -l
```

Install a specific version (without the `go` prefix):

```bash
sudo go-sdk install --version-number 1.24.1
# or
sudo go-sdk install -v 1.24.1
```

Install without root into a custom directory (no `sudo` needed):

```bash
go-sdk install --latest --dir "$HOME/.local"
# or via environment variable
GO_SDK_INSTALL_DIR="$HOME/.local" go-sdk install --latest
```

When using a custom directory, add `<dir>/go/bin` to your `PATH` (e.g. `$HOME/.local/go/bin`).

Flags:

```text
  -d, --dir string              Target install directory (default /usr/local or $GO_SDK_INSTALL_DIR)
  -h, --help                    help for install
  -l, --latest                  Install latest version of Go
  -v, --version-number string   Version number to install (e.g. 1.16.3)
```

> `--latest` and `--version-number` are mutually exclusive. If the requested version is not present
> in the available list, the command fails with `version <n> is not available`. Before downloading,
> `go-sdk` checks that the target directory is writable and fails fast with guidance if it is not.

### `completion`

Generates an autocompletion script for your shell (bash, zsh, fish, powershell). For example, with
fish:

```bash
go-sdk completion fish | source
```

## Environment variables

| Variable              | Purpose                                                                          |
| --------------------- | -------------------------------------------------------------------------------- |
| `GITHUB_TOKEN`        | Optional GitHub token used when listing versions, to raise the rate limit.        |
| `GO_SDK_INSTALL_DIR`  | Optional install destination (default `/usr/local`). The `--dir` flag overrides it. |

## How it works

- **Latest version** is resolved from `https://go.dev/VERSION?m=text`.
- **Available versions** come from the `golang/go` repository tags via the GitHub API, sorted by
  semantic version.
- **Installation** verifies the target directory is writable, downloads
  `https://go.dev/dl/<version>.linux-amd64.tar.gz`, extracts it into the install directory
  (`/usr/local` by default), and removes the temporary archive.
- The current version is read by executing `go version` on the host.

## Project structure

```text
.
├── cmd/
│   ├── root.go              # Root Cobra command; registers subcommands
│   └── go-sdk/
│       └── main.go          # Entrypoint
└── internal/
    ├── cli/
    │   ├── functions.go     # Shared RunEFunc type
    │   └── build/
    │       ├── current.go   # `current` command
    │       ├── list.go      # `list` command
    │       └── install.go   # `install` command
    └── sdk/
        ├── info.go          # Version discovery (go version, GitHub tags)
        └── install.go       # Download & extract logic
```

## Limitations & notes

- **Linux/amd64 only** — the download URL and the `/usr/local` install target are hardcoded.
- **Requires privileges** to write into `/usr/local` (run with `sudo`), or use `--dir` /
  `GO_SDK_INSTALL_DIR` to install into a writable directory without root.
- **Network on startup** — every invocation fetches the version list from the GitHub API during
  initialization; set `GITHUB_TOKEN` to avoid rate limiting.

## Contributing

This project is open source and welcomes contributions and suggestions. Please read
[CONTRIBUTING.md](./CONTRIBUTING.md) for how to report issues, propose changes, and follow the
commit conventions.

## License

Released under the [MIT License](./LICENSE).
