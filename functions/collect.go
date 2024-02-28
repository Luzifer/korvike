// Package functions contains custom functions specific to korvike and
// returns the sprig functions with added korvike specific functions
package functions

import (
	"log"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

var (
	templateFunctions     = sprig.FuncMap()
	templateFunctionsLock sync.Mutex
)

func registerFunction(name string, f interface{}) {
	templateFunctionsLock.Lock()
	defer templateFunctionsLock.Unlock()

	if _, ok := templateFunctions[name]; ok {
		log.Printf("overwriting existing function %q", name)
	}

	templateFunctions[name] = f
}

// GetFunctionMap exports all functions used in korvike to be used in own projects
// Example:
//
//	import korvike "github.com/Luzifer/korvike"
//	tpl := template.New("mytemplate").Funcs(korvike.GetFunctionMap())
func GetFunctionMap() template.FuncMap {
	return templateFunctions
}
