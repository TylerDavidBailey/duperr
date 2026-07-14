// Package duperr defines an Analyzer that reports duplicate error messages
// within a package.
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

// Analyzer reports duplicate error messages within a package.
//
// Two errors constructed from the same message string are indistinguishable
// when debugging: a log line or test failure carrying the message cannot be
// traced back to a single call site. The analyzer flags every occurrence
// after the first of a constant message passed to errors.New or fmt.Errorf.
//
// fmt.Errorf format strings containing verbs other than %w are skipped:
// their dynamic arguments already make the resulting messages distinct.
// Files ending in _test.go are ignored.
//
//nolint:gochecknoglobals // analyzers are exported as package-level vars by convention
var Analyzer = &analysis.Analyzer{
	Name:     "duperr",
	Doc:      "reports duplicate error messages within a package",
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

	occurrences := map[string][]token.Pos{}

	insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}
		if msg, ok := constantMessage(pass, call); ok {
			occurrences[msg] = append(occurrences[msg], call.Pos())
		}
	})

	for msg, positions := range occurrences {
		if len(positions) < 2 {
			continue
		}
		slices.Sort(positions)
		first := pass.Fset.Position(positions[0])
		firstLoc := fmt.Sprintf("%s:%d", filepath.Base(first.Filename), first.Line)
		for _, pos := range positions[1:] {
			pass.Reportf(pos, "duplicate error message %q (first used at %s)", msg, firstLoc)
		}
	}

	return nil, nil
}

// constantMessage returns the message string of an errors.New or fmt.Errorf
// call outside _test.go files whose runtime message is fully determined by a
// constant argument.
func constantMessage(pass *analysis.Pass, call *ast.CallExpr) (string, bool) {
	if len(call.Args) == 0 {
		return "", false
	}

	fn, ok := typeutil.Callee(pass.TypesInfo, call).(*types.Func)
	if !ok {
		return "", false
	}
	isErrorf := false
	switch fn.FullName() {
	case "errors.New":
	case "fmt.Errorf":
		isErrorf = true
	default:
		return "", false
	}

	if strings.HasSuffix(pass.Fset.Position(call.Pos()).Filename, "_test.go") {
		return "", false
	}

	tv := pass.TypesInfo.Types[call.Args[0]]
	if tv.Value == nil || tv.Value.Kind() != constant.String {
		return "", false
	}
	msg := constant.StringVal(tv.Value)
	if isErrorf && hasDynamicVerbs(msg) {
		return "", false
	}

	return msg, true
}

// hasDynamicVerbs reports whether format contains a formatting verb other
// than %w. Such verbs consume dynamic arguments, so call sites sharing the
// format string still produce distinct messages at runtime.
func hasDynamicVerbs(format string) bool {
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			continue
		}
		i++
		for i < len(format) && strings.ContainsRune("+-# 0123456789.[]", rune(format[i])) {
			i++
		}
		if i >= len(format) {
			return true // trailing bare % is malformed; play it safe
		}
		switch format[i] {
		case '%', 'w':
		default:
			return true
		}
	}
	return false
}
