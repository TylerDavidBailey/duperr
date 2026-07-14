// Package plugin registers duperr as a golangci-lint module plugin.
//
// It lives outside the root package so that importing duperr pulls in no
// init function or plugin machinery; golangci-lint custom reaches it via
// the import field in .custom-gcl.yml.
package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/TylerDavidBailey/duperr"
)

//nolint:gochecknoinits // module plugins register themselves in init
func init() {
	register.Plugin("duperr", newPlugin)
}

func newPlugin(_ any) (register.LinterPlugin, error) {
	return plugin{}, nil
}

type plugin struct{}

func (plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{duperr.Analyzer}, nil
}

func (plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
