package compare

import (
	"testing"
)

func BenchmarkNoPipe(b *testing.B) {
	NoPipe()
}

func BenchmarkWithPipe(b *testing.B) {
	WithPipe()
}

// go test -bench=. -run=none
