package native

import (
	"testing"

	"fknsrs.biz/p/searchfiles/tests"
)

func TestShared(t *testing.T) {
	tests.All(Default, t)
}
