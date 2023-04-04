package tests

import (
	"context"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"fknsrs.biz/p/searchfiles"
)

func getRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "data")
}

func Test_All(driver searchfiles.Driver, t *testing.T) {
	for _, fn := range []func(driver searchfiles.Driver, t *testing.T){
		Test_SearchLiteral_PositiveCases,
		Test_SearchLiteral_QueryNotFound,
		Test_SearchLiteral_RootDirNotFound,
		Test_SearchLiteral_QueryNotLiteralMatch,
		Test_SearchRegexp_PositiveCaseSingleFile,
		Test_SearchRegexp_PositiveCaseMultipleFiles,
		Test_SearchRegexp_QueryNotFound,
		Test_SearchRegexp_InvalidRegex,
		Test_SearchRegexp_RootDirNotFound,
	} {
		pc := reflect.ValueOf(fn).Pointer()
		f := runtime.FuncForPC(pc)
		t.Run(path.Base(f.Name()), func(t *testing.T) { fn(driver, t) })
	}
}

func Test_SearchLiteral_PositiveCases(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchLiteral(context.Background(), getRoot(), "test")
	a.NoError(err)
	a.ElementsMatch([]string{"/file1.txt", "/file2.txt", "/file4.txt", "/subdir/file3.txt"}, results)
}

func Test_SearchLiteral_QueryNotFound(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchLiteral(context.Background(), getRoot(), "notfound")
	a.NoError(err)
	a.Empty(results)
}

func Test_SearchLiteral_RootDirNotFound(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchLiteral(context.Background(), "/directory-does-not-exist", "test")
	a.Error(err)
	a.Empty(results)
}

func Test_SearchLiteral_QueryNotLiteralMatch(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchLiteral(context.Background(), getRoot(), "Test")
	a.NoError(err)
	a.Empty(results)
}

func Test_SearchRegexp_PositiveCaseSingleFile(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchRegexp(context.Background(), getRoot(), `\d{3}-\d{3}-\d{4}`)
	a.NoError(err)
	a.ElementsMatch([]string{"/file4.txt"}, results)
}

func Test_SearchRegexp_PositiveCaseMultipleFiles(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchRegexp(context.Background(), getRoot(), `test`)
	a.NoError(err)
	a.ElementsMatch([]string{"/file1.txt", "/file2.txt", "/file4.txt", "/subdir/file3.txt"}, results)
}

func Test_SearchRegexp_QueryNotFound(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchRegexp(context.Background(), getRoot(), `notfound`)
	a.NoError(err)
	a.Empty(results)
}

func Test_SearchRegexp_InvalidRegex(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchRegexp(context.Background(), getRoot(), `[`)
	a.Error(err)
	a.Empty(results)
}

func Test_SearchRegexp_RootDirNotFound(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchRegexp(context.Background(), "/directory-does-not-exist", `test`)
	a.Error(err)
	a.Empty(results)
}

func Benchmark_All(driver searchfiles.Driver, b *testing.B) {
	for _, fn := range []func(driver searchfiles.Driver, b *testing.B){
		Benchmark_SearchLiteralWithMatches,
		Benchmark_SearchLiteralWithoutMatches,
		Benchmark_SearchRegexpWithMatches,
		Benchmark_SearchRegexpWithoutMatches,
	} {
		pc := reflect.ValueOf(fn).Pointer()
		f := runtime.FuncForPC(pc)
		b.Run(path.Base(f.Name()), func(b *testing.B) { fn(driver, b) })
	}
}

func Benchmark_SearchLiteralWithMatches(driver searchfiles.Driver, b *testing.B) {
	ctx := context.Background()
	root := getRoot()

	for i := 0; i < b.N; i++ {
		_, _ = driver.SearchLiteral(ctx, root, "test")
	}
}

func Benchmark_SearchLiteralWithoutMatches(driver searchfiles.Driver, b *testing.B) {
	ctx := context.Background()
	root := getRoot()

	for i := 0; i < b.N; i++ {
		_, _ = driver.SearchLiteral(ctx, root, "test-xxx-not-found")
	}
}

func Benchmark_SearchRegexpWithMatches(driver searchfiles.Driver, b *testing.B) {
	ctx := context.Background()
	root := getRoot()

	for i := 0; i < b.N; i++ {
		_, _ = driver.SearchRegexp(ctx, root, "test")
	}
}

func Benchmark_SearchRegexpWithoutMatches(driver searchfiles.Driver, b *testing.B) {
	ctx := context.Background()
	root := getRoot()

	for i := 0; i < b.N; i++ {
		_, _ = driver.SearchRegexp(ctx, root, "test-xxx-not-found")
	}
}
