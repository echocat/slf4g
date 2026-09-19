//go:build !go1.22

package log

import (
	"fmt"
	"strings"
)

func formatStrSlice(in []string) string {
	var result strings.Builder
	for i, v := range in {
		if i > 0 {
			result.WriteByte(' ')
		}
		result.WriteString(v)
	}
	return result.String()
}

func formatAnySlice(in []any) string {
	var result strings.Builder
	for i, v := range in {
		if i > 0 {
			result.WriteByte(' ')
		}
		result.WriteString(fmt.Sprint(v))
	}
	return result.String()
}
