package functions

import (
	"strings"

	"github.com/echocat/slf4g/native/execution"
	"github.com/echocat/slf4g/native/formatter/encoding"
)

// EncodeMultilineWithIndent will indent the given message per line with the given indents
// and writes the result to the given `to` encoding.TextEncoder. The first line can have a
// different indent than the following lines.
func EncodeMultilineWithIndent(firstLineIndent, otherLinesIndent string, to encoding.TextEncoder, message string) error {
	message = strings.ReplaceAll(message, "\r", "")
	for i, line := range strings.Split(message, "\n") {
		ident := &firstLineIndent
		var prefixExecution execution.Execution
		if i > 0 {
			ident = &otherLinesIndent
			prefixExecution = to.WriteByteChecked('\n')
		}

		if err := execution.Execute(
			prefixExecution,
			to.WriteStringPChecked(ident),
			to.WriteStringChecked(line),
		); err != nil {
			return err
		}
	}
	return nil
}

// IndentMultiline will indent the given message per line with the given indents.
// The first line can have a different indent than the following lines.
func IndentMultiline(firstLineIndent, otherLinesIndent string, message string) string {
	encoder := encoding.NewBufferedTextEncoder()

	if err := EncodeMultilineWithIndent(firstLineIndent, otherLinesIndent, encoder, message); err != nil {
		panic(err)
	}

	return encoder.String()
}
