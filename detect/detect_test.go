package detect

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"fknsrs.biz/p/searchfiles"
	"fknsrs.biz/p/searchfiles/driver/ag"
	"fknsrs.biz/p/searchfiles/driver/grep"
	"fknsrs.biz/p/searchfiles/driver/pt"
	"fknsrs.biz/p/searchfiles/driver/rg"
)

func TestNames(t *testing.T) {
	a := assert.New(t)
	a.ElementsMatch([]string{"ag", "grep", "native", "pt", "rg"}, searchfiles.DriverNames())
}

func TestDetectDefault(t *testing.T) {
	a := assert.New(t)

	driverName, err := Detect(context.Background(), nil)
	a.NoError(err)
	a.Equal("ag", driverName)
}

func TestDetectNativeOnly(t *testing.T) {
	a := assert.New(t)

	driverName, err := Detect(context.Background(), []string{"native"})
	a.NoError(err)
	a.Equal("native", driverName)
}

func TestDetectNone(t *testing.T) {
	a := assert.New(t)

	driverName, err := Detect(context.Background(), []string{})
	a.ErrorIs(err, ErrNoWorkingDriver)
	a.Equal("", driverName)
}

func TestDetectBrokenExceptGrep(t *testing.T) {
	a := assert.New(t)

	agProgram := ag.Default.Program
	ptProgram := pt.Default.Program
	rgProgram := rg.Default.Program
	defer func() {
		ag.Default.Program = agProgram
		pt.Default.Program = ptProgram
		rg.Default.Program = rgProgram
	}()
	ag.Default.Program = "xxx-does-not-exist"
	pt.Default.Program = "xxx-does-not-exist"
	rg.Default.Program = "xxx-does-not-exist"

	driverName, err := Detect(context.Background(), nil)
	a.NoError(err)
	a.Equal("grep", driverName)
}

func TestDetectBrokenExceptNative(t *testing.T) {
	a := assert.New(t)

	agProgram := ag.Default.Program
	grepProgram := grep.Default.Program
	ptProgram := pt.Default.Program
	rgProgram := rg.Default.Program
	defer func() {
		ag.Default.Program = agProgram
		grep.Default.Program = grepProgram
		pt.Default.Program = ptProgram
		rg.Default.Program = rgProgram
	}()
	ag.Default.Program = "xxx-does-not-exist"
	grep.Default.Program = "xxx-does-not-exist"
	pt.Default.Program = "xxx-does-not-exist"
	rg.Default.Program = "xxx-does-not-exist"

	driverName, err := Detect(context.Background(), nil)
	a.NoError(err)
	a.Equal("native", driverName)
}
