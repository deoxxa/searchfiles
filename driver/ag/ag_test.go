package ag

import (
	"testing"

	"fknsrs.biz/p/searchfiles/tests"
)

func TestShared(t *testing.T) {
	tests.Test_All(Default, t)
}

func BenchmarkShared(b *testing.B) {
	tests.Benchmark_All(Default, b)
}
