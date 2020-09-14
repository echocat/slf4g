package formatter

import (
	"encoding/json"
	"fmt"
	log "github.com/echocat/slf4g"
)

type SimpleValueFormatter struct {
	QuoteType QuoteType
}

func (instance *SimpleValueFormatter) FormatValue(v interface{}, _ log.Provider) ([]byte, error) {
	if instance.QuoteType == QuoteTypeMinimal {
		switch str := v.(type) {
		case string:
			if stringNeedsQuoting(str) {
				return json.Marshal(str)
			}
			return []byte(str), nil
		case *string:
			if stringNeedsQuoting(*str) {
				return json.Marshal(*str)
			}
			return []byte(*str), nil
		}
		return json.Marshal(v)
	}

	if instance.QuoteType == QuoteTypeEverything {
		return json.Marshal(fmt.Sprint(v))
	}

	return json.Marshal(v)
}

type QuoteType uint8

const (
	QuoteTypeMinimal    QuoteType = 0
	QuoteTypeNormal     QuoteType = 1
	QuoteTypeEverything QuoteType = 2
)

func stringNeedsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == ':' ||
			ch == '/' || ch == '\\' ||
			ch == '@' || ch == '^' || ch == '+' || ch == '#' ||
			ch == '(' || ch == ')' ||
			ch == '[' || ch == ']' ||
			ch == '{' || ch == '}') {
			return true
		}
	}
	return false
}
