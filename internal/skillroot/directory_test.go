package skillroot

import (
	"testing"

	"github.com/refactor-ia/cortex/internal/skilldest"
)

// A caller-owned root is refused before any destination is derived. With an
// empty symbolic plan every input is invalid, so these cases pin only that a
// bad root can never reach derivation; the accepted path is covered end to end
// by internal/cli.
func TestResolveIntoDirectoryRefusesBadRoots(t *testing.T) {
	for name, root := range map[string]string{
		"empty":    "",
		"relative": "relative/path",
		"bare":     "/",
		"dirty":    "/tmp/./x",
	} {
		if _, err := ResolveIntoDirectory(skilldest.Plan{}, root); err == nil {
			t.Fatalf("%s: root %q must be refused", name, root)
		}
	}
}
