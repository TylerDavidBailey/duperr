package a

import "errors"

var (
	errTest1 = errors.New("duplicate in tests")
	errTest2 = errors.New("duplicate in tests")
)
