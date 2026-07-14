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
	errPct2 = fmt.Errorf("100%% duplicated") // want `duplicate error message "100%% duplicated" \(first used at c\.go:25\)`
)

var (
	errBare1 = fmt.Errorf("50%")
	errBare2 = fmt.Errorf("50%")
)

func wrapIndexed(err error) error {
	return fmt.Errorf("indexed wrap: %[1]w", err)
}

func wrapIndexedAgain(err error) error {
	return fmt.Errorf("indexed wrap: %[1]w", err) // want `duplicate error message "indexed wrap: %\[1\]w" \(first used at c\.go:35\)`
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
