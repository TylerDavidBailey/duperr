// Duperr reports duplicate error messages within a package: two errors
// built from the same message are indistinguishable when debugging, since
// a log line or test failure carrying the message cannot be traced back to
// a single call site.
//
// Usage:
//
//	go vet -vettool=$(which duperr) ./...
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/TylerDavidBailey/duperr"
)

func main() {
	singlechecker.Main(duperr.Analyzer)
}
