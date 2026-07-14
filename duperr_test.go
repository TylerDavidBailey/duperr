package duperr_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/TylerDavidBailey/duperr"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()
	analysistest.Run(t, analysistest.TestData(), duperr.Analyzer, "a", "b")
}
