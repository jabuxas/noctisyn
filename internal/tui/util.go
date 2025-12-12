package tui

import (
	"regexp"
	"strings"
)

var unsafeChars = regexp.MustCompile(`[\/\\:*?"<>|]`)

func safeFilename(title string) string {
	safe := unsafeChars.ReplaceAllString(title, "")
	safe = strings.ReplaceAll(safe, " ", "_")
	safe = strings.TrimSpace(safe)
	if len(safe) > 200 {
		safe = safe[:200]
	}
	if safe == "" {
		safe = "untitled"
	}
	return safe
}
