//go:build darwin

package internal

import (
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"
)

// List returns all stored key aliases (macOS implementation).
// Uses "security dump-keychain" to find entries where svce == "mk",
// then extracts the corresponding acct (alias) value.
func List() ([]string, error) {
	out, err := exec.Command("/usr/bin/security", "dump-keychain").Output()
	if err != nil {
		return nil, fmt.Errorf("failed to dump keychain: %w", err)
	}

	var aliases []string
	lines := strings.Split(string(out), "\n")

	// Parse the dump-keychain output.
	// Each generic password entry looks like:
	//
	//   keychain: "/Users/.../login.keychain-db"
	//   class: "genp"
	//   attributes:
	//       0x00000007 <blob>="mk"
	//       0x00000008 <blob>="nvidia"
	//       "acct"<blob>="nvidia"
	//       "svce"<blob>="mk"
	//       ...
	//   password: <base64-encoded>
	//
	// Both hex codes and text labels appear. We match on both patterns.

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Look for the start of a new entry: "class:" line
		if !strings.Contains(line, "class:") {
			i++
			continue
		}

		// Scan forward within this entry's attributes block
		var svce, acct string
		for j := i + 1; j < len(lines); j++ {
			attrLine := lines[j]

			// Stop at the next entry or end of attributes
			if strings.Contains(attrLine, "class:") || strings.Contains(attrLine, "password:") {
				break
			}

			trimmed := strings.TrimSpace(attrLine)

			// Try to extract service value from this line
			if svce == "" {
				if val := extractAttrValue(trimmed, "svce"); val != "" {
					svce = val
				}
			}

			// Try to extract account value from this line
			if acct == "" {
				if val := extractAttrValue(trimmed, "acct"); val != "" {
					acct = val
				}
			}
		}

		if svce == service && acct != "" {
			aliases = append(aliases, acct)
		}

		i++
	}

	return aliases, nil
}

// extractAttrValue tries to extract a value from a keychain attribute line.
// It handles multiple formats found in security dump-keychain output:
//
//   "svce"<blob>="mk"          (text label + blob)
//   svce<blob>="mk"            (text label without quotes)
//   0x00000007 <blob>="mk"     (hex code for svce)
//   0x00000008 <blob>="nvidia" (hex code for acct)
func extractAttrValue(line, attrName string) string {
	// Pattern 1: text label with quotes: "svce"<blob>="value"
	if idx := strings.Index(line, `"`+attrName+`"<blob>="`); idx != -1 {
		start := idx + len(`"`+attrName+`"<blob>="`)
		rest := line[start:]
		if end := strings.Index(rest, `"`); end != -1 {
			return rest[:end]
		}
	}
	// Pattern 1b: text label with quotes + hex blob: "acct"<blob>=0xE58898...
	if idx := strings.Index(line, `"`+attrName+`"<blob>=0x`); idx != -1 {
		start := idx + len(`"`+attrName+`"<blob>=0x`)
		hexStr := strings.TrimRight(line[start:], " \t")
		if decoded, err := hex.DecodeString(hexStr); err == nil {
			return string(decoded)
		}
	}

	// Pattern 2: text label without quotes: svce<blob>="value"
	if idx := strings.Index(line, attrName+`<blob>="`); idx != -1 {
		start := idx + len(attrName+`<blob>="`)
		rest := line[start:]
		if end := strings.Index(rest, `"`); end != -1 {
			return rest[:end]
		}
	}
	// Pattern 2b: text label without quotes + hex blob: acct<blob>=0xE58898...
	if idx := strings.Index(line, attrName+`<blob>=0x`); idx != -1 {
		start := idx + len(attrName+`<blob>=0x`)
		hexStr := strings.TrimRight(line[start:], " \t")
		if decoded, err := hex.DecodeString(hexStr); err == nil {
			return string(decoded)
		}
	}

	// Pattern 3: hex code for known attributes
	//   0x00000007 = svce (service)
	//   0x00000008 = acct (account)
	var hexCode string
	switch attrName {
	case "svce":
		hexCode = "0x00000007"
	case "acct":
		hexCode = "0x00000008"
	}

	if hexCode != "" {
		// Pattern 3a: hex code with quoted blob value
		//   0x00000008 <blob>="nvidia"
		if idx := strings.Index(line, hexCode+` <blob>="`); idx != -1 {
			start := idx + len(hexCode+` <blob>="`)
			rest := line[start:]
			if end := strings.Index(rest, `"`); end != -1 {
				val := rest[:end]
				if val == "<NULL>" {
					return ""
				}
				return val
			}
		}

		// Pattern 3b: hex code with hex-encoded blob value (non-ASCII, e.g. Chinese)
		//   0x00000008 <blob>=0xE58898E58898...
		hexPrefix := hexCode + ` <blob>=0x`
		if idx := strings.Index(line, hexPrefix); idx != -1 {
			start := idx + len(hexPrefix)
			hexStr := strings.TrimRight(line[start:], " \t")
			decoded, err := hex.DecodeString(hexStr)
			if err == nil {
				return string(decoded)
			}
		}
	}

	return ""
}
