package runtimematrix

import "testing"

// The directory projection target (issue #341) is explicit, never probed: its
// decision comes from DirectoryTargetMatrix and Decide must keep rejecting it.
func TestDirectoryTargetIsExplicitAndUnprobed(t *testing.T) {
	matrix := DirectoryTargetMatrix()
	if len(matrix.Decisions) != 1 || matrix.Decisions[0].ID != RuntimeDirectory ||
		matrix.Decisions[0].Outcome != OutcomePresentCompatible || matrix.Decisions[0].Action != Configure ||
		!matrix.Decisions[0].IncludeInTransaction || !matrix.Decisions[0].TouchAllowed || !matrix.HasCompatible {
		t.Fatalf("directory matrix = %#v", matrix)
	}
	if _, err := Decide([]Observation{{ID: RuntimeDirectory, Present: true, Version: "1", Compatibility: Compatible}}); err == nil {
		t.Fatal("Decide must reject a directory observation")
	}
	for _, id := range runtimeOrder {
		if id == RuntimeDirectory {
			t.Fatal("directory runtime must stay outside runtimeOrder")
		}
	}
}
