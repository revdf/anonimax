package config

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/parrotsec/packages/anonsurf/internal/fileops"
)

// UpdateKeys rewrites key=value pairs in the Anonimax config file.
// Missing keys are appended inside the AnonimaxConfig section.
func UpdateKeys(path string, values map[string]string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	pending := make(map[string]string, len(values))
	for key, value := range values {
		pending[key] = value
	}

	inSection := false
	sectionEnd := -1
	for i, line := range lines {
		trimmed := trimComment(line)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			name := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(trimmed, "["), "]"))
			if inSection && sectionEnd < 0 {
				sectionEnd = i
			}
			inSection = name == sectionName
			continue
		}
		if !inSection {
			continue
		}
		sectionEnd = i + 1
		key, _, ok := strings.Cut(trimmed, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if value, found := pending[key]; found {
			lines[i] = key + "=" + formatConfigValue(value)
			delete(pending, key)
		}
	}

	if len(pending) > 0 {
		insertAt := sectionEnd
		if insertAt < 0 {
			return fmt.Errorf("%s section not found in %s", sectionName, path)
		}
		extra := make([]string, 0, len(pending))
		for key, value := range pending {
			extra = append(extra, key+"="+formatConfigValue(value))
		}
		updated := make([]string, 0, len(lines)+len(extra))
		updated = append(updated, lines[:insertAt]...)
		updated = append(updated, extra...)
		updated = append(updated, lines[insertAt:]...)
		lines = updated
	}

	content := strings.Join(lines, "\n")
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return fileops.WriteAtomic(path, []byte(content), 0644)
}

func formatConfigValue(value string) string {
	if value == "" {
		return `""`
	}
	if strings.ContainsAny(value, " \t#;\"'") || strings.Contains(value, "=") {
		return strconvQuote(value)
	}
	return value
}

func strconvQuote(value string) string {
	var builder strings.Builder
	builder.WriteByte('"')
	for _, r := range value {
		switch r {
		case '\\', '"':
			builder.WriteByte('\\')
			builder.WriteRune(r)
		default:
			builder.WriteRune(r)
		}
	}
	builder.WriteByte('"')
	return builder.String()
}
