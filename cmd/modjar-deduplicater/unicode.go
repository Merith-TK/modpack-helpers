package main

import (
	"fmt"
	"strings"
)

func cleanUnicode(s string) string {
	for _, r := range s {
		if r > 127 {
			s = strings.ReplaceAll(s, string(r), fmt.Sprintf("&#%d;", r))
		}
	}
	return s
}
