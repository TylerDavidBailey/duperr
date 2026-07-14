package b

import "errors"

// Same message as package a: not flagged, duplicates are per-package.
var errOnly = errors.New("connecting to db")
