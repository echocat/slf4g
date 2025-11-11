//go:build go1.22

package log

import "fmt"

func formatStrSlice(in []string) string {
	var result []byte
	for i, v := range in {
		if i > 0 {
			result = append(result, ' ')
		}
		result = fmt.Append(result, v)
	}
	return string(result)
}

func formatAnySlice(in []interface{}) string {
	var result []byte
	for i, v := range in {
		if i > 0 {
			result = append(result, ' ')
		}
		result = fmt.Append(result, v)
	}
	return string(result)
}
