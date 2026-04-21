//go:build linux

package internal

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// List returns all stored key aliases (Linux implementation).
// Uses secret-tool search to query Secret Service for entries with service="mk".
func List() ([]string, error) {
	out, err := exec.Command("secret-tool", "search", "--all", "service", service).CombinedOutput()
	if err != nil {
		output := string(out)
		// secret-tool returns exit code 1 when no matching entries found.
		if strings.Contains(output, "no matching") || strings.Contains(output, "No matching") || len(out) == 0 {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to list keyring items: %w\n"+
			"Make sure gnome-keyring is running and libsecret-tools is installed:\n"+
			"  sudo apt install libsecret-tools", err)
	}

	var aliases []string
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "label = ") {
			label := strings.TrimPrefix(line, "label = ")
			alias := extractAliasFromLabel(label)
			if alias != "" {
				aliases = append(aliases, alias)
			}
		}
	}

	return aliases, nil
}

// extractAliasFromLabel extracts the alias from a label like "Password for 'openai' on 'mk'".
func extractAliasFromLabel(label string) string {
	prefix := "Password for '"
	if !strings.HasPrefix(label, prefix) {
		return ""
	}
	rest := label[len(prefix):]
	end := strings.Index(rest, "'")
	if end == -1 {
		return ""
	}
	return rest[:end]
}
