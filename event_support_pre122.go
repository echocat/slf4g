//go:build !go1.22

package log

import "fmt"

func formatStrSlice(in []string) string {
	var result string
	for i, v := range in {
		if i > 0 {
			result += " "
		}
		result += fmt.Sprint(v)
	}
	return result
}

func formatAnySlice(in []interface{}) string {
	var result string
	for i, v := range in {
		if i > 0 {
			result += " "
		}
		result += fmt.Sprint(v)
	}
	return result
}
