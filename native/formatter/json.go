package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/native/hints"
)

const (
	// DefaultKeyLevel is the default key to write the level of log entries to
	// the output with. See Json.KeyLevel for more information.
	DefaultKeyLevel = "level"
)

// Json is an implementation of Formatter which formats given log entries in a
// JSON format (https://en.wikipedia.org/wiki/JSON) where every log.Entry is one
// line in the output.
type Json struct {
	// KeyLevel is the key to write the level of log entries to the output with.
	// If not set DefaultKeyLevel is used.
	KeyLevel string

	// PrintRootLogger will (if set to true) also print the field logger for the
	// root logger. If set to false the logger field will be only printed for
	// every logger but not for the root one. If not set set
	// DefaultPrintRootLogger will be used.
	PrintRootLogger *bool
}

// NewJson creates a new instance of Text which is ready to use.
func NewJson(customizer ...func(*Json)) *Json {
	result := &Json{}
	for _, c := range customizer {
		c(result)
	}
	return result
}

// Format implements Formatter.Format()
func (instance *Json) Format(event log.Event, using log.Provider, _ hints.Hints) ([]byte, error) {
	fail := func(err error) ([]byte, error) {
		return nil, fmt.Errorf("cannot format event (%v): %w", event, err)
	}

	to := new(bytes.Buffer)
	enc := json.NewEncoder(to)

	if err := to.WriteByte('{'); err != nil {
		return fail(err)
	}

	if err := instance.encode(to, enc, instance.getLevelKey(), event.GetLevel()); err != nil {
		return fail(err)
	}

	printRootLogger := DefaultPrintRootLogger
	if v := instance.PrintRootLogger; v != nil {
		printRootLogger = *v
	}

	loggerKey := using.GetFieldKeysSpec().GetLogger()
	if err := event.ForEach(func(k string, v interface{}) error {
		if vl, ok := v.(fields.Lazy); ok {
			v = vl.Get()
		}

		if !printRootLogger && k == loggerKey && v == "ROOT" {
			return nil
		}
		if _, err := to.Write([]byte(",")); err != nil {
			return err
		}
		return instance.encode(to, enc, k, v)
	}); err != nil {
		return fail(err)
	}

	if _, err := to.WriteString("}\n"); err != nil {
		return fail(err)
	}

	return to.Bytes(), nil
}

func (instance *Json) encode(buf *bytes.Buffer, enc *json.Encoder, k string, v interface{}) error {
	if err := enc.Encode(k); err != nil {
		return err
	}
	buf.Truncate(buf.Len() - 1) // Because someone believe it is a great idea to but a \n always everywhere ...

	if _, err := buf.Write([]byte(":")); err != nil {
		return err
	}

	if ve, ok := v.(error); ok {
		v = ve.Error()
	}
	if vs, ok := v.(string); ok {
		v = strings.TrimRightFunc(vs, unicode.IsSpace)
	}
	if vs, ok := v.(*string); ok {
		v = strings.TrimRightFunc(*vs, unicode.IsSpace)
	}
	if err := enc.Encode(v); err != nil {
		return err
	}
	buf.Truncate(buf.Len() - 1) // Because someone believe it is a great idea to but a \n always everywhere ...

	return nil
}

func (instance *Json) getLevelKey() string {
	if v := instance.KeyLevel; v != "" {
		return v
	}
	return DefaultKeyLevel
}
