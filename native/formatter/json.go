package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/native/formatter/hints"
)

const (
	DefaultLevelKey = "level"
)

type Json struct {
	LevelKey          string
	PrintGlobalLogger bool
}

func NewJson() *Json {
	return &Json{}
}

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

	loggerKey := using.GetFieldKeysSpec().GetLogger()
	if err := event.ForEach(func(k string, v interface{}) error {
		if vl, ok := v.(fields.Lazy); ok {
			v = vl.Get()
		}

		if !instance.PrintGlobalLogger && k == loggerKey && v == "ROOT" {
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
	if v := instance.LevelKey; v != "" {
		return v
	}
	return DefaultLevelKey
}
