package a

import "errors"

var errCrossFile = errors.New("connecting to db") // want `duplicate error message "connecting to db" \(first used at a\.go:11\)`
