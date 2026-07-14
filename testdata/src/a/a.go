package a

import (
	"errors"
	"fmt"
)

const dbMsg = "constant db message"

var (
	errFirst  = errors.New("connecting to db")
	errSecond = errors.New("connecting to db") // want `duplicate error message "connecting to db" \(first used at a\.go:11\)`
	errThird  = fmt.Errorf("connecting to db") // want `duplicate error message "connecting to db" \(first used at a\.go:11\)`
)

var (
	errConst1 = errors.New(dbMsg)
	errConst2 = errors.New(dbMsg) // want `duplicate error message "constant db message" \(first used at a\.go:17\)`
)

func closeOnce(err error) error {
	return fmt.Errorf("closing db: %w", err)
}

func closeTwice(err error) error {
	return fmt.Errorf("closing db: %w", err) // want `duplicate error message "closing db: %w" \(first used at a\.go:22\)`
}

func readFirst(path string, err error) error {
	return fmt.Errorf("reading %s: %w", path, err)
}

func readSecond(path string, err error) error {
	return fmt.Errorf("reading %s: %w", path, err)
}

func dynamicFirst(msg string) error {
	return errors.New(msg)
}

func dynamicSecond(msg string) error {
	return errors.New(msg)
}
