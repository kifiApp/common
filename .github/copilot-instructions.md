# Kifi Common

This module (`github.com/kifiApp/common`) is the shared Go contract and utility library for Kifi services. It requires Go 1.23.0 (toolchain `go1.23.1`).

## Commands

```sh
# Compile and run all package tests (there are currently no committed *_test.go files).
go test ./...

# Run one test once tests exist.
go test ./common/jwt -run '^TestName$'

# Compile all packages without running tests.
go build ./...

# Run the configured linter.
golangci-lint run

# Regenerate Go protobuf and gRPC bindings after a .proto change.
make generate
```

`make generate` requires `protoc`, `protoc-gen-go`, and `protoc-gen-go-grpc`, and uses the well-known protobuf types from `/usr/local/include`. The generated files record `protoc` 3.14.0, `protoc-gen-go` 1.35.1, and `protoc-gen-go-grpc` 1.2.0. `make modsync` runs generation, upgrades all Go dependencies, and tidies the module; do not use it as a routine regeneration command.

## Architecture

- `proto/` is the source of truth for the shared API. Domain schemas (`accounts.proto`, `goals.proto`, `family.proto`, `devices.proto`) define messages, while `account-service.proto`, `goals-service.proto`, and `users.proto` define the gRPC services and import those domain schemas.
- Every schema uses `go_package = "github.com/kifiApp/common/gen;gen"`. Generation writes the checked-in Go bindings to `gen/gen/`; downstream services import them as `github.com/kifiApp/common/gen` (package `gen`).
- `common/jwt` provides the JWT interface and HMAC implementation used by services. `common/logger` owns process-wide Zap/slog and OpenTelemetry logging/tracing setup; its constructors return shutdown functions for exporters. `common/tools` exposes the ISO-3166 alpha-2 country-name map.

## Repository conventions

- Make protobuf changes in `proto/`, never in `gen/gen/*.pb.go` or `*_grpc.pb.go`. Regenerate all bindings with `make generate` and commit the matching generated changes with the schema change.
- Treat protobuf field numbers, field names, package names, and RPC method names as wire contracts. Preserve existing field numbers; add fields with new numbers rather than repurposing or renumbering fields.
- Keep all protobuf schemas on the shared `go_package` above so the generator continues to produce one Go package and imports resolve consistently.
- Follow `.golangci.yaml` for lint behavior. It deliberately excludes generated code and enables strict checks including `gofmt`, `errcheck`, `staticcheck`, `gosec`, and `govet`.

## Coding rules

- When a function can fail or an expected result may be absent, return the result with an `error` or `bool` rather than panicking or using a sentinel value. Do not add an unused `error` or `bool` return when failure or absence is not possible.
- Use TDD for behavior changes: write a failing test first, then implement the smallest change that makes it pass.
- Prefer the simplest implementation that satisfies the requirement; avoid speculative generalization or cleverness.
- Keep each function and type focused on one responsibility. Separate unrelated concerns such as database access, business logic, and transport mapping.
- Prefer concrete or generic types over `interface{}`/`any` and runtime type assertions. Introduce an interface only when a concrete need exists, such as test mocking or multiple implementations.
