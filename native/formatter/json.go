package formatter

import (
	"fmt"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/native/hints"
	nlevel "github.com/echocat/slf4g/native/level"
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

	// LevelFormatter is used to format the level.Level of a given log.Entry.
	// into the field with key of KeyLevel.
	LevelFormatter Level

	// PrintRootLogger will (if set to true) also print the field logger for the
	// root logger. If set to false the logger field will be only printed for
	// every logger but not for the root one. If not set set
	// DefaultPrintRootLogger will be used.
	PrintRootLogger *bool

	// KeySorter will force the printed fields to be sorted using this sorter.
	// The fields which contains the level.Level will be always the first,
	// regardless of the result of the KeySorter. If this field is empty the
	// fields are not sorted and the order is not deterministic and reliable.
	KeySorter fields.KeySorter
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
	if event == nil {
		return []byte{}, nil
	}

	to := newJsonEncoderBuffered()

	if err := executeChecked(
		to.WriteByteChecked('{'),
		instance.encodeLevelChecked(event, using, &to.jsonEncoder),
		instance.encodeValuesChecked(event, using, &to.jsonEncoder),
		to.WriteByteChecked('}'),
		to.WriteByteChecked('\n'),
	); err != nil {
		return nil, fmt.Errorf("cannot format event (%v): %w", event, err)
	}

	return to.Bytes(), nil
}

func (instance *Json) getLevelKey() string {
	if v := instance.KeyLevel; v != "" {
		return v
	}
	return DefaultKeyLevel
}

func (instance *Json) encodeLevelChecked(of log.Event, using log.Provider, to *jsonEncoder) checkedExecution {
	return func() error {
		lvl, err := instance.formatLevel(of, using)
		if err != nil {
			return err
		}
		return to.WriteKeyValue(instance.getLevelKey(), lvl)
	}
}

func (instance *Json) formatLevel(of log.Event, using log.Provider) (interface{}, error) {
	return instance.getLevelFormatter(using).FormatLevel(of.GetLevel(), using)
}

func (instance *Json) getLevelFormatter(using log.Provider) Level {
	if v := instance.LevelFormatter; v != nil {
		return v
	}
	if v, ok := using.(nlevel.NamesAware); ok {
		return NewNamesBasedLevel(v.GetLevelNames())
	}
	return DefaultLevel
}

func (instance *Json) encodeValuesChecked(of log.Event, using log.Provider, to *jsonEncoder) checkedExecution {
	return func() error {
		printRootLogger := instance.getPrintRootLogger()
		loggerKey := using.GetFieldKeysSpec().GetLogger()
		consumer := func(k string, v interface{}) error {
			if vl, ok := v.(fields.Lazy); ok {
				v = vl.Get()
			}
			if !printRootLogger && k == loggerKey && v == "ROOT" {
				return nil
			}
			return executeChecked(
				to.WriteByteChecked(','),
				to.WriteKeyValueChecked(k, v),
			)
		}
		if sorter := instance.KeySorter; sorter != nil {
			return fields.SortedForEach(of, sorter, consumer)
		}
		return of.ForEach(consumer)
	}
}

func (instance *Json) getPrintRootLogger() bool {
	if v := instance.PrintRootLogger; v != nil {
		return *v
	}
	return DefaultPrintRootLogger
}
