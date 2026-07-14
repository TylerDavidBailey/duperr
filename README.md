# duperr

[![CI](https://github.com/TylerDavidBailey/duperr/actions/workflows/ci.yml/badge.svg)](https://github.com/TylerDavidBailey/duperr/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/TylerDavidBailey/duperr.svg)](https://pkg.go.dev/github.com/TylerDavidBailey/duperr)
[![Go Report Card](https://goreportcard.com/badge/github.com/TylerDavidBailey/duperr)](https://goreportcard.com/report/github.com/TylerDavidBailey/duperr)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A Go linter that reports duplicate error messages within a package.

```go
var errDial = errors.New("connecting to db")
...
var errPing = errors.New("connecting to db") // duplicate error message "connecting to db" (first used at dial.go:11)
```

Two errors built from the same message are indistinguishable when debugging:
a log line or test failure carrying the message cannot be traced back to a
single call site. Existing linters don't catch this — `dupl` works on code
fragments, and `goconst` suggests the opposite fix (sharing the string).

Mature codebases accumulate these: duperr reports 20–139 duplicates each in
Caddy, Hugo, and Prometheus. A typical find, from Hugo's config loader:

```go
configs, err = fromLoadConfigResult(d.Fs, d.Logger, res)
if err != nil {
    return nil, fmt.Errorf("failed to create config from modules config: %w", err)
}
if err := configs.transientErr(); err != nil {
    return nil, fmt.Errorf("failed to create config from modules config: %w", err)
}
```

When that message turns up in a log, there is no way to tell which of the
two calls failed.

## What it checks

- Constant message strings passed to `errors.New` and `fmt.Errorf`, compared
  per package. Named constants and constant concatenation count; dynamic
  messages are ignored.
- `fmt.Errorf` format strings are only compared when they contain no verbs
  besides `%w` — dynamic verbs (`%s`, `%d`, …) already make the resulting
  messages distinct at runtime.
- Files ending in `_test.go` and generated files are skipped.

Every occurrence after the first is reported, pointing back to the first.

## Usage

### With `go vet`

```sh
go install github.com/TylerDavidBailey/duperr/cmd/duperr@latest
go vet -vettool=$(which duperr) ./...
```

### With golangci-lint (module plugin)

duperr registers itself as a [module plugin](https://golangci-lint.run/plugins/module-plugins/).
Create `.custom-gcl.yml` next to your `.golangci.yml`:

```yaml
version: v2.12.2
plugins:
  - module: github.com/TylerDavidBailey/duperr
    version: v0.1.0
```

Build the custom binary once with `golangci-lint custom`, then enable the
linter in `.golangci.yml`:

```yaml
linters:
  enable:
    - duperr
  settings:
    custom:
      duperr:
        type: module
        description: reports duplicate error messages within a package
```

Run the built binary (`./custom-gcl` by default) instead of `golangci-lint`.

## See also

[golangci-lint-config](https://github.com/TylerDavidBailey/golangci-lint-config) —
the opinionated golangci-lint config this repo dogfoods.

## License

[MIT](LICENSE)
