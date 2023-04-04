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

func All(driver searchfiles.Driver, t *testing.T) {
	for _, fn := range []func(driver searchfiles.Driver, t *testing.T){
		SearchLiteral_PositiveCases,
		SearchLiteral_QueryNotFound,
		SearchLiteral_RootDirNotFound,
		SearchLiteral_QueryNotLiteralMatch,
		SearchRegexp_PositiveCaseSingleFile,
		SearchRegexp_PositiveCaseMultipleFiles,
		SearchRegexp_QueryNotFound,
		SearchRegexp_InvalidRegex,
		SearchRegexp_RootDirNotFound,
	} {
		pc := reflect.ValueOf(fn).Pointer()
		f := runtime.FuncForPC(pc)
		t.Run(path.Base(f.Name()), func(t *testing.T) { fn(driver, t) })
	}
}

func SearchLiteral_PositiveCases(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchLiteral(context.Background(), getRoot(), "test")
	a.NoError(err)
	a.ElementsMatch([]string{"/file1.txt", "/file2.txt", "/file4.txt", "/subdir/file3.txt"}, results)
}

func SearchLiteral_QueryNotFound(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchLiteral(context.Background(), getRoot(), "notfound")
	a.NoError(err)
	a.Empty(results)
}

func SearchLiteral_RootDirNotFound(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchLiteral(context.Background(), "/directory-does-not-exist", "test")
	a.Error(err)
	a.Empty(results)
}

func SearchLiteral_QueryNotLiteralMatch(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchLiteral(context.Background(), getRoot(), "Test")
	a.NoError(err)
	a.Empty(results)
}

func SearchRegexp_PositiveCaseSingleFile(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchRegexp(context.Background(), getRoot(), `\d{3}-\d{3}-\d{4}`)
	a.NoError(err)
	a.ElementsMatch([]string{"/file4.txt"}, results)
}

func SearchRegexp_PositiveCaseMultipleFiles(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchRegexp(context.Background(), getRoot(), `test`)
	a.NoError(err)
	a.ElementsMatch([]string{"/file1.txt", "/file2.txt", "/file4.txt", "/subdir/file3.txt"}, results)
}

func SearchRegexp_QueryNotFound(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchRegexp(context.Background(), getRoot(), `notfound`)
	a.NoError(err)
	a.Empty(results)
}

func SearchRegexp_InvalidRegex(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchRegexp(context.Background(), getRoot(), `[`)
	a.Error(err)
	a.Empty(results)
}

func SearchRegexp_RootDirNotFound(driver searchfiles.Driver, t *testing.T) {
	a := assert.New(t)
	results, err := driver.SearchRegexp(context.Background(), "/directory-does-not-exist", `test`)
	a.Error(err)
	a.Empty(results)
}
