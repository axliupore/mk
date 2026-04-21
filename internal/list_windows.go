//go:build windows

package internal

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// List returns all stored key aliases (Windows implementation).
// Uses cmdkey /list to query Windows Credential Manager for entries whose target starts with "mk:".
func List() ([]string, error) {
	out, err := exec.Command("cmdkey", "/list").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list credentials: %w", err)
	}

	var aliases []string
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Match lines containing "LegacyGeneric:target=mk:"
		// This part is always in English regardless of OS locale.
		if strings.Contains(line, "LegacyGeneric:target="+service+":") {
			alias := extractAliasFromTarget(line)
			if alias != "" {
				aliases = append(aliases, alias)
			}
		}
	}

	return aliases, nil
}

// extractAliasFromTarget extracts the alias from a line like "Target: LegacyGeneric:target=mk:openai".
func extractAliasFromTarget(line string) string {
	idx := strings.Index(line, service+":")
	if idx == -1 {
		return ""
	}
	alias := line[idx+len(service)+1:]
	return strings.TrimSpace(alias)
}
