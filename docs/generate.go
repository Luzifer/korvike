// Generate documentation from real executions
package main

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"

	"github.com/Luzifer/korvike/functions"
)

type (
	docFunc functions.FunctionDocumentation
)

var (
	//go:embed header.md
	docHeader string

	//go:embed func.gotmpl
	docTpl string
)

func main() {
	var err error

	if os.Getenv("KORVIKE") == "" {
		logrus.Fatal("provide path to `korvike` binary in $KORVIKE")
	}

	funcs := functions.GetFunctionDocumentation()

	slices.SortFunc(funcs, func(a, b functions.FunctionDocumentation) int {
		return strings.Compare(a.Title, b.Title)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println(docHeader) //nolint:forbidigo // intended to write docs to stdout

	for _, f := range funcs {
		df := docFunc(f)
		logger := logrus.WithField("func", df.Title)

		if err = df.run(ctx); err != nil {
			logger.WithError(err).Fatal("running function")
		}

		out, err := df.render()
		if err != nil {
			logger.WithError(err).Fatal("rending function")
		}

		fmt.Println(out) //nolint:forbidigo // intended to write docs to stdout
	}
}

func (d docFunc) render() (string, error) {
	buf := new(bytes.Buffer)

	tpl, err := template.New(d.Title).Parse(docTpl)
	if err != nil {
		return "", fmt.Errorf("parsing func.gotmpl: %w", err)
	}

	if err = tpl.Execute(buf, d); err != nil {
		return "", fmt.Errorf("running func.gotmpl: %w", err)
	}

	return strings.TrimSpace(buf.String()), nil
}

func (d *docFunc) run(ctx context.Context) (err error) {
	if d.Output != "" {
		// We do have fake output, don't run.
		return nil
	}

	tmpDir, err := os.MkdirTemp("", "")
	if err != nil {
		return fmt.Errorf("creating tempdir: %w", err)
	}
	defer func() {
		if removeErr := os.RemoveAll(tmpDir); removeErr != nil {
			err = errors.Join(err, removeErr)
		}
	}()

	for i, pc := range d.PreCommands {
		//#nosec:G204 // locally defined commands in tempdir
		cmd := exec.CommandContext(ctx, "/usr/bin/env", "bash", "-euxo", "pipefail", "-c", pc)
		cmd.Dir = tmpDir
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			return fmt.Errorf("running pre-command %d: %w", i, err)
		}
	}

	buf := new(bytes.Buffer)

	var args []string
	for k, v := range d.Vars {
		args = append(args, "-v", fmt.Sprintf("%s=%s", k, v))
	}

	var env []string
	for k, v := range d.Environ {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	//#nosec:G204 // locally defined commands in tempdir
	cmd := exec.CommandContext(ctx, os.Getenv("KORVIKE"), args...)
	cmd.Dir = tmpDir
	cmd.Env = env
	cmd.Stderr = os.Stderr
	cmd.Stdout = buf
	cmd.Stdin = strings.NewReader(d.Template)

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("running korvike: %w", err)
	}

	d.Output = strings.TrimSpace(buf.String())

	return nil
}
