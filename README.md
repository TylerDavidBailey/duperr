# duperr

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

## What it checks

- Constant message strings passed to `errors.New` and `fmt.Errorf`, compared
  per package. Named constants and constant concatenation count; dynamic
  messages are ignored.
- `fmt.Errorf` format strings are only compared when they contain no verbs
  besides `%w` — dynamic verbs (`%s`, `%d`, …) already make the resulting
  messages distinct at runtime.
- Files ending in `_test.go` are skipped.

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
    version: latest
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
