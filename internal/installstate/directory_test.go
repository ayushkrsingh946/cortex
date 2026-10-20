package installstate

import (
	"testing"

	"github.com/refactor-ia/cortex/internal/runtimematrix"
	"github.com/refactor-ia/cortex/internal/skilldest"
)

func TestValidPairAcceptsDirectory(t *testing.T) {
	if !validPair(runtimematrix.RuntimeDirectory, skilldest.RootKindDirectory) {
		t.Fatal("the (directory, directory) pair must be valid")
	}
	if validPair(runtimematrix.RuntimeDirectory, skilldest.RootKindPiUserAgent) || validPair(runtimematrix.RuntimeDirectory, "") {
		t.Fatal("directory must only pair with the directory root kind")
	}
	// The directory runtime is skill-only: it never carries actor-aware state.
	if validActorState(Manifest{owner: "cortex", scope: "user", runtimeID: runtimematrix.RuntimeDirectory, rootKind: skilldest.RootKindDirectory}) {
		t.Fatal("directory runtime must not satisfy actor-aware state")
	}
}
