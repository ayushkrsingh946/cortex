package skilldest

import (
	"testing"

	"github.com/refactor-ia/cortex/internal/runtimematrix"
)

func TestRootForDirectoryRuntime(t *testing.T) {
	root, ok := rootFor(runtimematrix.RuntimeDirectory)
	if !ok || root != RootKindDirectory {
		t.Fatalf("rootFor(directory) = (%q, %v)", root, ok)
	}
}
