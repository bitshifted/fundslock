# Agent Guidelines & Runbook (AGENTS.md)

Welcome, AI Agent! This file serves as your comprehensive instruction manual and onboarding guide for working in the **backend** directory. It contains our architectural layout, service components, operational commands, coding guidelines, and testing requirements to ensure your contributions align perfectly with our standards.

---

## 1. Service Overview

**fundslock-be** is the REST API backend for the FundsLock decentralized escrow system. It acts as an off-chain layer providing authentication services and facilitating communication with the Graph protocol for tracking FundsLock smart contract events.

### Key Technology Stack
- **Language:** Go (`1.25.x`)
- **CLI parser:** Kong (`github.com/alecthomas/kong`)
- **HTTP Routing:** Chi (`github.com/go-chi/chi/v5`)
- **Logging:** Zerolog (`github.com/rs/zerolog`)
- **Linting:** Golangci-lint (`v2.12.2`)
- **License Compliance:** Addlicense (`v1.1.1`)
- **Task Orchestration:** `Makefile`

---

## 2. Directory Architecture

```
/
├── cli/                      # CLI Commands defined using kong parser
│   ├── cli.go                # Command structures, start server & version runners
│   └── version.go            # Version structures and constants
├── log/                      # Logging packages
│   └── logger.go             # Logging configuration with zerolog (Debug/Info modes)
├── srv/                      # HTTP REST API server
│   └── server.go             # Router (chi), middleware, server config & starting logic
├── target/                   # Build artifacts, test results & coverage (generated on build)
│   ├── coverage.html         # HTML coverage report
│   ├── coverage.out          # Raw coverage output
│   └── ...                   # Platform-specific build binaries (linux/windows/macos)
├── .golangci.yml             # golangci-lint configuration file
├── Dockerfile                # Ubuntu 24.04-based container image definition
├── go.mod                    # Go module dependencies
├── go.sum                    # Go module checksums
├── main.go                   # Main entry point (parses CLI options and boots logger/commands)
└── Makefile                  # Task automation commands
```

---

## 3. Architecture & Execution Flow

The backend service operates as a CLI tool built with Kong.

### Execution Flow:
1. `main.go` parses arguments and flags using `kong`.
2. It detects the presence of `--enable-debug` to initialize the logger package (`log/logger.go`) with `zerolog.DebugLevel` and custom console output. Otherwise, it defaults to `zerolog.InfoLevel`.
3. It dispatches commands to the `cli/` command runner methods (`cli.VersionCmd.Run` or `cli.StartCmd.Run`).
4. Running the `start` command invokes `srv.Start()` to bind and run the HTTP server.

### HTTP Server Setup (`srv/server.go`):
- Implemented with the lightweight `go-chi/chi/v5` router.
- Employs `chi/middleware.Logger`.
- **Listen Address:** Defaulting to `:3000`.
- **Timeout Constraints:** Read timeout `10s`, Write timeout `10s`, Idle timeout `120s`.
- Root endpoint (`GET /`) responds with a simple `"Hello, world!"` message.

---

## 4. Development & Operational Runbook

Our `Makefile` orchestrates development tasks. You are expected to use these targets to validate any code modifications.

### Main Verification Pipeline
Before submitting any changes, you must ensure they pass our build pipeline:
```bash
# Performs a clean, initializes targets, checks licenses, lints, tests, and builds binaries for all platforms
make build
```

### Individual Verification Steps
```bash
# Clean build target directory
make clean

# Initialize target directories
make init

# Check code correctness, quality, and style with golangci-lint
make lint

# Run all unit and integration tests and generate coverage report
make test

# Add license headers to all Go source files, Dockerfile, and Makefile
make add-license-headers

# Check that license headers are present across source files
make check-license-headers
```

---

## 5. Coding Standards & Conventions

AI agents must strictly adhere to these standards when writing or modifying code in this workspace:

### Strict File Header Rule
Every single `.go` file **must** begin with the exact Copyright and SPDX license headers. Failure to include these will cause `make check-license-headers` to fail.
```go
// Copyright 2026 Bitshift ED
// SPDX-License-Identifier: MPL-2.0
```
Use `make add-license-headers` to apply this header automatically to new files.

### Code Formatting
Code must be formatted using `gofmt` and `goimports`. The `golangci-lint` tool in `make lint` enforces formatting guidelines automatically.

### Linting & Quality Rules
We use a strict `golangci-lint` setup (`.golangci.yml`). Ensure that you:
- Avoid unused imports or variables.
- Always check returned errors.
- Do not bypass linter warnings using generic suppresses unless absolutely necessary and documented.
- Check that your functions conform to complexity limits (e.g., `gocyclo` max complexity `15`, `funlen` max lines `100` / statements `60`).

---

## 6. Testing & Validation Strategy

We aim for comprehensive coverage on any added features or bug fixes.
1. **Unit Testing:** Write unit tests using the standard Go `testing` package (and optionally asserting libraries already established in `go.mod`/`go.sum`). Keep tests alongside the implementation with the suffix `_test.go`.
2. **Mocking & Isolation:** When implementing network or external dependency-reliant code, write mocks or use interfaces to facilitate unit testing without external resources.
3. **Execution:** Always verify your changes pass by running `make test`. Ensure coverage does not regress.

---

## 7. Running & Deployment

### Local Running
To start the server locally on port 3000:
```bash
go run main.go start
```
To enable debug logging:
```bash
go run main.go  --enable-debug start
```

### Docker Packaging
The service can be containerized using the provided `Dockerfile`.
The Dockerfile expects a binary to exist at the build argument path `binary_location` (defaults to `target/linux-amd64/fundslock-be`).

To build the container image:
1. Compile the target binaries first:
   ```bash
   make build
   ```
2. Build the Docker image:
   ```bash
   docker build -t bitshifted/fundslock-be:latest .
   ```
