// Package plugin registers duperr as a golangci-lint module plugin.
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/TylerDavidBailey/duperr"
)

//nolint:gochecknoinits // module plugins register themselves in init
func init() {
	register.Plugin("duperr", New)
}

// New returns the duperr linter plugin. Settings are not used.
func New(_ any) (register.LinterPlugin, error) {
	return plugin{}, nil
}

type plugin struct{}

// BuildAnalyzers returns the duperr analyzer.
func (plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{duperr.Analyzer}, nil
}

// GetLoadMode reports that duperr needs type information.
func (plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
