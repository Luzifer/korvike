package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/Luzifer/go_helpers/v2/env"
	korvike "github.com/Luzifer/korvike/functions"
	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		Input          string   `flag:"in,i" default:"-" description:"Source to read the template from ('-' = stdin)"`
		KeyPairs       []string `flag:"key-value,v" default:"" description:"Key-value pairs to use in templating (-v key=value)"`
		Output         string   `flag:"out,o" default:"-" description:"Destination to write the output to ('-' = stdout)"`
		VersionAndExit bool     `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func initApp() (err error) {
	if err := rconfig.Parse(&cfg); err != nil {
		return fmt.Errorf("parsing CLI options: %w", err)
	}

	return nil
}

func openInput(f string) (io.Reader, error) {
	if f == "-" {
		return os.Stdin, nil
	}

	if _, err := os.Stat(f); err != nil {
		return nil, fmt.Errorf("finding file  %q: %s", f, err)
	}

	r, err := os.Open(f) //#nosec G304 // Intended to use user-given file
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	return r, nil
}

func openOutput(f string) (io.Writer, error) {
	if f == "-" {
		return os.Stdout, nil
	}

	w, err := os.Create(f) //#nosec G304 // Intended to use user-given file
	if err != nil {
		return nil, fmt.Errorf("creating file: %w", err)
	}

	return w, nil
}

func main() {
	var err error
	if err = initApp(); err != nil {
		log.Fatalf("initializing app: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("korvike %s\n", version) //nolint:forbidigo
		os.Exit(0)
	}

	in, err := openInput(cfg.Input)
	if err != nil {
		log.Fatalf("opening input: %s", err)
	}

	out, err := openOutput(cfg.Output)
	if err != nil {
		log.Fatalf("opening output: %s", err)
	}

	rawTpl, err := io.ReadAll(in)
	if err != nil {
		log.Fatalf("reading from input: %s", err)
	}

	tpl, err := template.New("in").Funcs(korvike.GetFunctionMap()).Parse(string(rawTpl))
	if err != nil {
		log.Fatalf("parsing template: %s", err)
	}

	vars := map[string]interface{}{}
	for k, v := range env.ListToMap(cfg.KeyPairs) {
		vars[k] = v
	}

	korvike.SetSubTemplateVariables(vars)
	if err := tpl.Execute(out, vars); err != nil {
		log.Fatalf("executing template: %s", err)
	}
}
