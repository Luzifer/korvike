// Package functions contains custom functions specific to korvike and
// returns the sprig functions with added korvike specific functions
package functions

import (
	"maps"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type (
	// FuncMapOpt adds a group of functions to a function map.
	FuncMapOpt func(template.FuncMap)
)

var (
	templateFunctions     = make(template.FuncMap)
	templateFunctionsSafe = make(template.FuncMap)
	templateFunctionsLock sync.Mutex
)

//revive:disable-next-line:flag-parameter // Safety is registry metadata, not control flow chosen by callers.
func registerFunction(name string, f any, isSafe bool) {
	templateFunctionsLock.Lock()
	defer templateFunctionsLock.Unlock()

	templateFunctions[name] = f

	if isSafe {
		templateFunctionsSafe[name] = f
	}
}

// GetFunctionMap exports all functions used in korvike to be used in own projects
// Example:
//
//	import korvike "github.com/Luzifer/korvike"
//	tpl := template.New("mytemplate").Funcs(korvike.GetFunctionMap())
func GetFunctionMap(opts ...FuncMapOpt) template.FuncMap {
	if len(opts) == 0 {
		// Backwards compatibility: If no options are given, load all Sprig
		// and all Korvike functions
		opts = []FuncMapOpt{WithAll}
	}

	funcs := make(template.FuncMap)

	for _, opt := range opts {
		opt(funcs)
	}

	return funcs
}

// WithAll includes all available Sprig and Korvike functions in the function map.
func WithAll(fm template.FuncMap) {
	maps.Insert(fm, maps.All(sprig.FuncMap()))
	maps.Insert(fm, maps.All(templateFunctions))
}

// WithKorvikeAll includes all Korvike functions in the function map.
func WithKorvikeAll(fm template.FuncMap) {
	maps.Insert(fm, maps.All(templateFunctions))
}

// WithKorvikeSafe includes only Korvike functions safe for untrusted templates.
//
// Safe functions must not allow template authors to read host files, access
// secrets, execute nested templates with broader privileges, perform network
// access, mutate process state, or otherwise gain capabilities beyond pure
// data formatting and transformation.
func WithKorvikeSafe(fm template.FuncMap) {
	maps.Insert(fm, maps.All(templateFunctionsSafe))
}

// WithSafe includes Sprig's hermetic functions and Korvike functions safe for untrusted templates.
func WithSafe(fm template.FuncMap) {
	maps.Insert(fm, maps.All(sprig.HermeticTxtFuncMap()))
	maps.Insert(fm, maps.All(templateFunctionsSafe))
}
