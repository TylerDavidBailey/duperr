// Duperr reports duplicate error messages within a package.
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
