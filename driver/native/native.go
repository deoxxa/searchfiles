package native

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"fknsrs.biz/p/searchfiles"
)

var (
	Default *Driver
)

func init() {
	Default = &Driver{}
	searchfiles.Register("native", Default)
}

type Driver struct{}

func (d *Driver) SelfTest(ctx context.Context) error {
	return nil
}

func (d *Driver) SearchLiteral(ctx context.Context, directory, query string) ([]string, error) {
	a, err := d.search(ctx, directory, regexp.QuoteMeta(query))
	if err != nil {
		return nil, fmt.Errorf("native.Driver.SearchLiteral: %w", err)
	}

	return a, nil
}

func (d *Driver) SearchRegexp(ctx context.Context, directory, query string) ([]string, error) {
	a, err := d.search(ctx, directory, query)
	if err != nil {
		return nil, fmt.Errorf("native.Driver.SearchLiteral: %w", err)
	}

	return a, nil
}

func (d *Driver) search(ctx context.Context, directory, query string) ([]string, error) {
	re, err := regexp.Compile(query)
	if err != nil {
		return nil, fmt.Errorf("native.Driver.search: could not compile query: %w", err)
	}

	collector := &matchCollector{ctx: ctx, directory: directory, regexp: re}

	if err := filepath.Walk(directory, collector.walk); err != nil {
		return nil, fmt.Errorf("native.Driver.search: could not walk directory: %w", err)
	}

	return []string(collector.files), nil
}

type matchCollector struct {
	ctx       context.Context
	directory string
	regexp    *regexp.Regexp
	files     []string
}

func (c *matchCollector) walk(path string, info fs.FileInfo, pathErr error) error {
	if pathErr != nil {
		return pathErr
	}

	if !info.Mode().IsRegular() {
		return nil
	}

	if err := c.ctx.Err(); err != nil {
		return fmt.Errorf("native.matchCollector.walk: %w", err)
	}

	fd, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("native.matchCollector.walk: could not open file %q: %w", path, err)
	}
	defer fd.Close()

	matched, err := matchReader(c.ctx, c.regexp, fd)
	if err != nil {
		return fmt.Errorf("native.matchCollector.walk: could not search file %q: %w", path, err)
	}

	if matched {
		c.files = append(c.files, strings.TrimPrefix(path, c.directory))
	}

	if err := fd.Close(); err != nil {
		return fmt.Errorf("native.matchCollector.walk: could not close file %q: %w", path, err)
	}

	return nil
}

func matchReader(ctx context.Context, re *regexp.Regexp, fd *os.File) (bool, error) {
	ch := make(chan bool, 1)
	go func() {
		ch <- re.MatchReader(bufio.NewReader(fd))
	}()

	select {
	case <-ctx.Done():
		return false, fmt.Errorf("native.matchReader: %w", ctx.Err())
	case r := <-ch:
		return r, nil
	}
}
