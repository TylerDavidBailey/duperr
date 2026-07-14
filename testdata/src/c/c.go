package c

import (
	"errors"
	"fmt"
	stderrors "errors"
)

var (
	errAlias1 = errors.New("alias duplicate")
	errAlias2 = stderrors.New("alias duplicate") // want `duplicate error message "alias duplicate" \(first used at c\.go:10\)`
)

var (
	errConcat1 = errors.New("split " + "message")
	errConcat2 = errors.New(`split message`) // want `duplicate error message "split message" \(first used at c\.go:15\)`
)

var (
	errEmpty1 = errors.New("")
	errEmpty2 = errors.New("") // want `duplicate error message "" \(first used at c\.go:20\)`
)

var (
	errPct1 = fmt.Errorf("100%% duplicated")
	errPct2 = fmt.Errorf("100%% duplicated") // want `duplicate error message "100% duplicated" \(first used at c\.go:25\)`
)

// Messages compare by their runtime text: %% in a fmt.Errorf format renders
// as a single percent sign, so it matches the unescaped errors.New literal
// but not the escaped one.
var (
	errUnescaped = errors.New("75% literal")
	errEscaped   = fmt.Errorf("75%% literal") // want `duplicate error message "75% literal" \(first used at c\.go:33\)`
	errTwoPct    = errors.New("75%% literal")
)

// A literal %w in errors.New stays literal at runtime, so it never matches
// the fmt.Errorf wrap template sharing its source text.
var errLiteralVerb = errors.New("literal wrap: %w")

func wrapLiteralVerb(err error) error {
	return fmt.Errorf("literal wrap: %w", err)
}

var (
	errBare1 = fmt.Errorf("50%")
	errBare2 = fmt.Errorf("50%")
)

func wrapIndexed(err error) error {
	return fmt.Errorf("indexed wrap: %[1]w", err)
}

func wrapIndexedAgain(err error) error {
	return fmt.Errorf("indexed wrap: %[1]w", err) // want `duplicate error message "indexed wrap: %\[1\]w" \(first used at c\.go:52\)`
}

type fakeErrors struct{}

func (fakeErrors) New(_ string) error { return nil }

func shadowedPackageName() {
	errors := fakeErrors{}
	_ = errors.New("shadowed duplicate")
	_ = errors.New("shadowed duplicate")
}

func indirectCall() {
	newErr := errors.New
	_ = newErr("indirect duplicate")
	_ = newErr("indirect duplicate")
}
