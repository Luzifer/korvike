package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

func init() {
	if err := rconfig.Parse(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("korvike %s\n", version)
		os.Exit(0)
	}
}

func openInput(f string) (io.Reader, error) {
	if f == "-" {
		return os.Stdin, nil
	}

	if _, err := os.Stat(f); err != nil {
		return nil, errors.New("Could not find file " + f)
	}

	return os.Open(f)
}

func openOutput(f string) (io.Writer, error) {
	if f == "-" {
		return os.Stdout, nil
	}

	return os.Create(f)
}

func main() {
	in, err := openInput(cfg.Input)
	if err != nil {
		log.Fatalf("Unable to open input: %s", err)
	}

	out, err := openOutput(cfg.Output)
	if err != nil {
		log.Fatalf("Unable to open output: %s", err)
	}

	rawTpl, err := ioutil.ReadAll(in)
	if err != nil {
		log.Fatalf("Unable to read from input: %s", err)
	}

	tpl, err := template.New("in").Funcs(korvike.GetFunctionMap()).Parse(string(rawTpl))
	if err != nil {
		log.Fatalf("Unable to parse template: %s", err)
	}

	vars := map[string]interface{}{}
	for k, v := range env.ListToMap(cfg.KeyPairs) {
		vars[k] = v
	}

	korvike.SetSubTemplateVariables(vars)
	if err := tpl.Execute(out, vars); err != nil {
		log.Fatalf("Unable to execute template: %s", err)
	}
}
