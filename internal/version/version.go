// Package version resolves a descriptive version string for local builds
// that weren't stamped via ldflags.
package version

import (
	"os/exec"
	"strings"
)

// ResolveDev builds a descriptive version string for a local ("dev") build by
// querying the git working tree: the most recent tag, the short commit hash,
// and a "-dev" suffix, e.g. "v1.4.0-abc1234-dev". Falls back to just the
// commit hash, or the literal "dev", if git/tags aren't available.
func ResolveDev() string {
	tag := gitTag()
	commit := gitShortCommit()
	if tag == "" && commit == "" {
		return "dev"
	}
	parts := []string{tag, commit, "dev"}
	return strings.Join(nonEmpty(parts), "-")
}

func gitTag() string {
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func gitShortCommit() string {
	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func nonEmpty(parts []string) []string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
