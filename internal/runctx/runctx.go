package runctx

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type CheckErrorFunc func(cmd *exec.Cmd, err error, stdout, stderr *bytes.Buffer) error

func Run(ctx context.Context, program string, arguments []string, checkError CheckErrorFunc) ([]string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, program, arguments...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	ch := make(chan error, 1)
	go func() {
		ch <- cmd.Run()
	}()

	select {
	case <-ctx.Done():
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return nil, fmt.Errorf("runctx.Run: %w", ctx.Err())
	case err := <-ch:
		if err != nil && checkError != nil {
			err = checkError(cmd, err, &stdout, &stderr)
		}

		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				return nil, fmt.Errorf("runctx.Run: command failed with exit code %d: %w", exitErr.ExitCode(), err)
			}

			return nil, fmt.Errorf("runctx.Run: command failed: %w", err)
		}
	}

	lines := make([]string, 0)

	for _, e := range strings.Split(stdout.String(), "\n") {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}

		lines = append(lines, e)
	}

	return lines, nil
}
