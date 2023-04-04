package runctx_test

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"fknsrs.biz/p/searchfiles/internal/runctx"
)

func TestRun(t *testing.T) {
	t.Parallel()

	var testErr = fmt.Errorf("test")

	tests := []struct {
		name       string
		ctx        func(ctx context.Context) context.Context
		command    []string
		checkError runctx.CheckErrorFunc
		expected   []string
		err        error
		rootErr    error
	}{
		{
			name:     "success zero lines",
			command:  []string{"true"},
			expected: []string{},
		},
		{
			name:     "success one line",
			command:  []string{"echo", "hello", "world"},
			expected: []string{"hello world"},
		},
		{
			name:     "success two lines",
			command:  []string{"sh", "-c", "echo a; echo b"},
			expected: []string{"a", "b"},
		},
		{
			name:    "failure via exit status",
			command: []string{"sh", "-c", "echo 'hello' && echo 'world' >&2 && exit 1"},
			err:     fmt.Errorf("runctx.Run: command failed with exit code 1: exit status 1"),
		},
		{
			name:    "failure via checkError",
			command: []string{"sh", "-c", "echo test_stdout && echo test_stderr >&2 && exit 1"},
			checkError: func(cmd *exec.Cmd, err error, stdout, stderr *bytes.Buffer) error {
				return testErr
			},
			err:     fmt.Errorf("runctx.Run: command failed: test"),
			rootErr: testErr,
		},
		{
			name:    "success via checkError",
			command: []string{"sh", "-c", "echo test_stdout && echo test_stderr >&2 && exit 1"},
			checkError: func(cmd *exec.Cmd, err error, stdout, stderr *bytes.Buffer) error {
				return nil
			},
			expected: []string{"test_stdout"},
		},
		{
			name: "context canceled",
			ctx: func(ctx context.Context) context.Context {
				ctx2, cancel := context.WithCancel(ctx)
				defer cancel()
				return ctx2
			},
			command: []string{"sleep", "1"},
			err:     fmt.Errorf("runctx.Run: context canceled"),
			rootErr: context.Canceled,
		},
		{
			name: "context timed out",
			ctx: func(ctx context.Context) context.Context {
				ctx2, _ := context.WithDeadline(ctx, time.Now().Add(time.Millisecond*100))
				return ctx2
			},
			command: []string{"sleep", "1"},
			err:     fmt.Errorf("runctx.Run: context deadline exceeded"),
			rootErr: context.DeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			ctx := context.Background()
			if tt.ctx != nil {
				ctx = tt.ctx(ctx)
			}

			var arguments []string
			if len(tt.command) > 1 {
				arguments = tt.command[1:]
			}

			lines, err := runctx.Run(ctx, tt.command[0], arguments, tt.checkError)

			if tt.err != nil {
				a.EqualError(err, tt.err.Error())

				if tt.rootErr != nil {
					a.ErrorIs(err, tt.rootErr)
				}

				return
			}

			a.NoError(err)
			a.Equal(tt.expected, lines)
		})
	}
}
