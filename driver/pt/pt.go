package pt

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"fknsrs.biz/p/searchfiles"
	"fknsrs.biz/p/searchfiles/internal/runctx"
)

var (
	Default *Driver
)

func init() {
	Default = &Driver{}
	searchfiles.Register("pt", Default)
}

const (
	DefaultProgram = "pt"
)

type Driver struct {
	Program string
}

func (d *Driver) program() string {
	if d.Program == "" {
		return DefaultProgram
	}
	return d.Program
}

func (d *Driver) SelfTest(ctx context.Context) error {
	if _, err := runctx.Run(ctx, d.program(), []string{"--version"}, nil); err != nil {
		return err
	}

	return nil
}

func (d *Driver) SearchLiteral(ctx context.Context, directory, query string) ([]string, error) {
	if st, err := os.Stat(directory); err != nil {
		return nil, fmt.Errorf("pt.Driver.SearchLiteral: could not stat %q: %w", directory, err)
	} else if !st.IsDir() {
		return nil, fmt.Errorf("pt.Driver.SearchLiteral: %q is not a directory", directory)
	}

	files, err := d.search(ctx, query, directory)
	if err != nil {
		return nil, fmt.Errorf("pt.Driver.SearchLiteral: %w", err)
	}

	return cleanResults(directory, files), nil
}

func (d *Driver) SearchRegexp(ctx context.Context, directory, query string) ([]string, error) {
	if st, err := os.Stat(directory); err != nil {
		return nil, fmt.Errorf("pt.Driver.SearchRegexp: could not stat %q: %w", directory, err)
	} else if !st.IsDir() {
		return nil, fmt.Errorf("pt.Driver.SearchRegexp: %q is not a directory", directory)
	}

	files, err := d.search(ctx, "-e", query, directory)
	if err != nil {
		return nil, fmt.Errorf("pt.Driver.SearchRegexp: %w", err)
	}

	return cleanResults(directory, files), nil
}

func (d *Driver) search(ctx context.Context, args ...string) ([]string, error) {
	files, err := runctx.Run(ctx, d.program(), append([]string{"-l"}, args...), func(cmd *exec.Cmd, err error, stdout, stderr *bytes.Buffer) error {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 0 {
					return nil
				}

				if status.ExitStatus() == 1 && stdout.Len() == 0 && stderr.Len() == 0 {
					return nil
				}
			}
		}

		return err
	})
	if err != nil {
		return nil, fmt.Errorf("pt.Driver.search: could not run command: %w", err)
	}

	return files, nil
}

func cleanResults(directory string, input []string) []string {
	var files []string

	for _, e := range input {
		e = strings.TrimSpace(e)
		e = strings.TrimPrefix(e, directory)
		if e == "" {
			continue
		}

		files = append(files, e)
	}

	return files
}
