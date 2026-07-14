// Package duperr defines an Analyzer that reports duplicate error messages
// within a package.
//
// Two errors constructed from the same message string are indistinguishable
// when debugging: a log line or test failure carrying the message cannot be
// traced back to a single call site. The analyzer flags every occurrence
// after the first of a constant message passed to errors.New or fmt.Errorf,
// pointing back to the first.
//
// fmt.Errorf format strings containing verbs other than %w are skipped:
// their dynamic arguments already make the resulting messages distinct.
// Files ending in _test.go and generated files are ignored.
package duperr

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

const doc = `reports duplicate error messages within a package

Two errors constructed from the same message string are indistinguishable
when debugging: a log line or test failure carrying the message cannot be
traced back to a single call site. The analyzer flags every occurrence
after the first of a constant message passed to errors.New or fmt.Errorf,
pointing back to the first.

fmt.Errorf format strings containing verbs other than %w are skipped:
their dynamic arguments already make the resulting messages distinct.
Files ending in _test.go and generated files are ignored.`

// Analyzer reports duplicate error messages within a package. See the
// package documentation for details.
//
//nolint:gochecknoglobals // analyzers are exported as package-level vars by convention
var Analyzer = &analysis.Analyzer{
	Name:     "duperr",
	Doc:      doc,
	URL:      "https://github.com/TylerDavidBailey/duperr",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

//nolint:nilnil // analyzers without a result type return nil, nil by contract
func run(pass *analysis.Pass) (any, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("getting inspector result for %s", pass.Pkg.Path())
	}

	occurrences := map[message][]token.Pos{}

	nodeFilter := []ast.Node{(*ast.File)(nil), (*ast.CallExpr)(nil)}
	insp.Nodes(nodeFilter, func(n ast.Node, push bool) bool {
		if !push {
			return false
		}
		switch n := n.(type) {
		case *ast.File:
			filename := pass.Fset.Position(n.Pos()).Filename
			return !strings.HasSuffix(filename, "_test.go") && !ast.IsGenerated(n)
		case *ast.CallExpr:
			if msg, ok := constantMessage(pass, n); ok {
				occurrences[msg] = append(occurrences[msg], n.Pos())
			}
		}
		return true
	})

	for msg, positions := range occurrences {
		if len(positions) < 2 {
			continue
		}
		// Raw token.Pos order follows file registration order, which is
		// not deterministic across loads; sort by resolved position.
		slices.SortFunc(positions, func(a, b token.Pos) int {
			pa, pb := pass.Fset.Position(a), pass.Fset.Position(b)
			if c := strings.Compare(pa.Filename, pb.Filename); c != 0 {
				return c
			}
			return pa.Offset - pb.Offset
		})
		first := pass.Fset.Position(positions[0])
		firstLoc := fmt.Sprintf("%s:%d", filepath.Base(first.Filename), first.Line)
		for _, pos := range positions[1:] {
			pass.Reportf(pos, "duplicate error message %q (first used at %s)", msg.text, firstLoc)
		}
	}

	return nil, nil
}

// message identifies an error message by what it produces at runtime: the
// message text with fmt.Errorf's %% escapes resolved, plus whether it is a
// wrap template. errors.New("x: %w") and fmt.Errorf("x: %w", err) share
// source text but can never produce the same runtime message.
type message struct {
	wraps bool // fmt.Errorf format containing %w
	text  string
}

// constantMessage returns the message of an errors.New or fmt.Errorf call
// whose runtime message is fully determined by a constant argument.
func constantMessage(pass *analysis.Pass, call *ast.CallExpr) (message, bool) {
	if len(call.Args) == 0 {
		return message{}, false
	}

	fn, ok := typeutil.Callee(pass.TypesInfo, call).(*types.Func)
	if !ok {
		return message{}, false
	}
	isErrorf := false
	switch fn.FullName() {
	case "errors.New":
	case "fmt.Errorf":
		isErrorf = true
	default:
		return message{}, false
	}

	tv := pass.TypesInfo.Types[call.Args[0]]
	if tv.Value == nil || tv.Value.Kind() != constant.String {
		return message{}, false
	}
	msg := constant.StringVal(tv.Value)
	if !isErrorf {
		return message{text: msg}, true
	}
	dynamic, wraps := scanVerbs(msg)
	if dynamic {
		return message{}, false
	}

	return message{wraps: wraps, text: strings.ReplaceAll(msg, "%%", "%")}, true
}

// scanVerbs reports whether format contains a formatting verb other than %w
// (such verbs consume dynamic arguments, so call sites sharing the format
// string still produce distinct messages at runtime) and whether it contains
// %w itself.
func scanVerbs(format string) (dynamic, wraps bool) {
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			continue
		}
		i++
		for i < len(format) && strings.ContainsRune("+-# 0123456789.[]", rune(format[i])) {
			i++
		}
		if i >= len(format) {
			return true, wraps // trailing bare % is malformed; play it safe
		}
		switch format[i] {
		case '%':
		case 'w':
			wraps = true
		default:
			return true, wraps
		}
	}

	return false, wraps
}
