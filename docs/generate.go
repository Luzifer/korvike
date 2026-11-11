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
	"go.yaml.in/yaml/v4"
)

type (
	docFile struct {
		Header    string    `yaml:"header"`
		Functions []docFunc `yaml:"functions"`
	}

	docFunc struct {
		Title       string            `yaml:"title"`
		Description string            `yaml:"description"`
		Environ     map[string]string `yaml:"environ"`
		PreCommands []string          `yaml:"preCommands"`
		Vars        map[string]string `yaml:"vars"`
		Template    string            `yaml:"template"`
		Output      string            `yaml:"output"`
	}
)

var (
	//go:embed docs.yaml
	docSrc string

	//go:embed func.gotmpl
	docTpl string
)

func main() {
	var (
		doc docFile
		err error
	)

	if os.Getenv("KORVIKE") == "" {
		logrus.Fatal("provide path to `korvike` binary in $KORVIKE")
	}

	if err = yaml.NewDecoder(strings.NewReader(docSrc)).Decode(&doc); err != nil {
		logrus.WithError(err).Fatal("parsing config")
	}

	slices.SortFunc(doc.Functions, func(a, b docFunc) int {
		return strings.Compare(a.Title, b.Title)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println(doc.Header) //nolint:forbidigo // intended to write docs to stdout

	for _, df := range doc.Functions {
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
