package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInstallDirectoryTargetRendersQAActorsAsSkills pins the #344 decision:
// the QA actors project into a directory target as self-contained SKILL.md
// files with valid frontmatter, and no QA backend binding is created.
func TestInstallDirectoryTargetRendersQAActorsAsSkills(t *testing.T) {
	home := canonicalHome(t)
	target := filepath.Join(home, "agent-skills")
	if err := os.MkdirAll(target, 0o700); err != nil {
		t.Fatal(err)
	}
	directoryLifecycleInstall(t, home, target)
	for _, actor := range []string{"requirements-analyst", "adversarial-tester", "evidence-auditor", "test-designer", "test-runner", "exploratory-tester"} {
		path := filepath.Join(target, "skills", "cortex-"+actor, "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("QA actor %s missing: %v", actor, err)
		}
		content := string(data)
		if !strings.HasPrefix(content, "---\n") || !strings.Contains(content, "name: \"cortex-"+actor+"\"") || !strings.Contains(content, "description:") || !strings.Contains(content, "license:") {
			t.Fatalf("QA actor %s frontmatter incomplete:\n%.200s", actor, content)
		}
	}
}
