package novel

import (
	"strings"
)

func wordsInOrder(s, query string) bool {
	sLower := strings.ToLower(s)
	words := strings.Fields(strings.ToLower(query))

	pos := 0
	for _, w := range words {
		idx := strings.Index(sLower[pos:], w)
		if idx == -1 {
			return false
		}
		pos += idx + len(w)
	}
	return true
}
