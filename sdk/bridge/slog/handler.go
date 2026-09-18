//go:build go1.21

package sdk

import (
	"context"
	"fmt"
	sdk "log/slog"
	"sync"

	log "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/fields"
	"github.com/echocat/slf4g/level"
)

func NewHandler(target log.CoreLogger, customizer ...func(*Handler)) *Handler {
	result := &Handler{
		Delegate: target,
	}

	for _, c := range customizer {
		c(result)
	}

	return result
}

type Handler struct {
	// Delegate is the [log.CoreLogger] of the slf4g framework where to forward all logged
	// events of this implementation to.
	//
	// If empty the result of [log.GetRootLogger] will be used.
	Delegate log.CoreLogger

	// LevelMapper holds the mapper which is used to transform the levels
	// between [sdk.Level] and [github.com/echocat/slf4g/level.Level].
	//
	// If empty [DefaultLevelMapper] will be used.
	LevelMapper LevelMapper

	// DetectSkipFrames defines for how much frames an element [log.Event]
	// should be skipped while reporting to [Handler.Delegate].
	//
	// If empty [DefaultDetectSkipFrames] will be used.
	DetectSkipFrames DetectSkipFrames

	parent       *Handler
	fieldKeyPath *fieldKeyPath
	attrs        attrs
}

// fieldKeyPath keeps group construction linear and resolves only paths used by attributes.
type fieldKeyPath struct {
	parent   *fieldKeyPath
	key      string
	length   int
	once     sync.Once
	resolved string
}

func newFieldKeyPath(parent *fieldKeyPath, key string) *fieldKeyPath {
	length := len(key) + 1
	if parent != nil {
		length += parent.length
	}
	return &fieldKeyPath{parent: parent, key: key, length: length}
}

func (instance *fieldKeyPath) prefix() string {
	if instance == nil {
		return ""
	}
	instance.once.Do(func() {
		result := make([]byte, instance.length)
		offset := len(result)
		for current := instance; current != nil; current = current.parent {
			offset--
			result[offset] = '.'
			offset -= len(current.key)
			copy(result[offset:], current.key)
		}
		instance.resolved = string(result)
	})
	return instance.resolved
}

// Enabled implements [sdk.Handler.Enabled]
func (instance *Handler) Enabled(_ context.Context, sl sdk.Level) bool {
	l, err := instance.mapFromSdkLevel(sl)
	if err != nil {
		return false
	}

	return instance.getDelegate().IsLevelEnabled(l)
}

// Handle implements [sdk.Handler.Handle]
func (instance *Handler) Handle(_ context.Context, record sdk.Record) error {
	delegate := instance.getDelegate()
	helperOf(delegate)()
	e, err := instance.eventOfRecord(delegate, record)
	if err != nil {
		return err
	}
	if record.PC != 0 && instance.DetectSkipFrames == nil {
		candidate := eventWithProgramCounter{Event: e, programCounter: record.PC}
		if delegate.Accepts(candidate) {
			e = candidate
		}
	}

	skipFrames := instance.getDetectSkipFrames()(1)
	delegate.Log(e, skipFrames)
	return nil
}

func (instance *Handler) eventOfRecord(logger log.CoreLogger, record sdk.Record) (log.Event, error) {
	l, err := instance.levelOfRecord(record)
	if err != nil {
		return nil, err
	}

	fds := instance.fieldsOfRecord(logger, record)

	return log.NewEventWithFields(logger, l, fds), nil
}

func (instance *Handler) fieldsOfRecord(logger log.CoreLogger, record sdk.Record) fields.Fields {
	fdsSpec := logger.GetProvider().GetFieldKeysSpec()
	numAttrs := record.NumAttrs()
	keyPrefix := ""
	if numAttrs > 0 {
		keyPrefix = instance.fieldKeyPath.prefix()
	}

	vs := make(attrs, 2+numAttrs)
	var i int
	vs[i] = sdk.Attr{
		Key:   fdsSpec.GetMessage(),
		Value: sdk.StringValue(record.Message),
	}
	i++

	vs[i] = sdk.Attr{
		Key:   fdsSpec.GetTimestamp(),
		Value: sdk.TimeValue(record.Time),
	}
	i++

	record.Attrs(func(v sdk.Attr) bool {
		vs[i] = sdk.Attr{
			Key:   keyPrefix + v.Key,
			Value: v.Value,
		}
		i++
		return true
	})

	return fields.NewLineage(vs, instance.fields())
}

func (instance *Handler) fields() fields.Fields {
	var lineage []*Handler
	totalAttrs := 0
	for current := instance; current != nil; current = current.parent {
		lineage = append(lineage, current)
		totalAttrs += len(current.attrs)
	}

	result := make(attrs, 0, totalAttrs)
	handledKeys := make(map[string]struct{}, totalAttrs)
	for start := 0; start < len(lineage); {
		end := start + 1
		for end < len(lineage) && lineage[end].fieldKeyPath == lineage[start].fieldKeyPath {
			end++
		}

		var segment attrs
		indexes := map[string]int{}
		for i := end - 1; i >= start; i-- {
			for _, attr := range lineage[i].attrs {
				if index, exists := indexes[attr.Key]; exists {
					segment[index] = attr
				} else {
					indexes[attr.Key] = len(segment)
					segment = append(segment, attr)
				}
			}
		}
		for _, attr := range segment {
			if _, handled := handledKeys[attr.Key]; !handled {
				handledKeys[attr.Key] = struct{}{}
				result = append(result, attr)
			}
		}
		start = end
	}
	return result
}

func (instance *Handler) levelOfRecord(record sdk.Record) (level.Level, error) {
	return instance.mapFromSdkLevel(record.Level)
}

func (instance *Handler) mapFromSdkLevel(sl sdk.Level) (level.Level, error) {
	l, err := instance.getLevelMapper().FromSdk(sl)
	if err != nil {
		return 0, fmt.Errorf("cannot map SDK's level %d to slf4g's level: %w", sl, err)
	}
	return l, nil
}

// WithAttrs implements [sdk.Handler.WithAttrs]
func (instance *Handler) WithAttrs(vs []sdk.Attr) sdk.Handler {
	nvs := make(attrs, 0, len(vs))
	keyPrefix := ""
	if len(vs) > 0 {
		keyPrefix = instance.fieldKeyPath.prefix()
	}
	nvs.add(keyPrefix, vs...)
	return &Handler{
		Delegate:         instance.Delegate,
		LevelMapper:      instance.LevelMapper,
		DetectSkipFrames: instance.DetectSkipFrames,
		parent:           instance,
		fieldKeyPath:     instance.fieldKeyPath,
		attrs:            nvs,
	}
}

// WithGroup implements [sdk.Handler.WithGroup]
func (instance *Handler) WithGroup(key string) sdk.Handler {
	if key == "" {
		return instance
	}
	return &Handler{
		Delegate:         instance.Delegate,
		LevelMapper:      instance.LevelMapper,
		DetectSkipFrames: instance.DetectSkipFrames,
		parent:           instance,
		fieldKeyPath:     newFieldKeyPath(instance.fieldKeyPath, key),
	}
}

func (instance *Handler) getDelegate() log.CoreLogger {
	if v := instance.Delegate; v != nil {
		return v
	}
	return log.GetRootLogger()
}

func (instance *Handler) getLevelMapper() LevelMapper {
	if v := instance.LevelMapper; v != nil {
		return v
	}
	return DefaultLevelMapper
}

func (instance *Handler) getDetectSkipFrames() DetectSkipFrames {
	if v := instance.DetectSkipFrames; v != nil {
		return v
	}
	return DefaultDetectSkipFrames
}

func helperOf(instance log.CoreLogger) func() {
	if wh, ok := instance.(interface {
		Helper() func()
	}); ok {
		return wh.Helper()
	}
	return func() {}
}
