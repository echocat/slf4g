//go:build go1.21
// +build go1.21

package sdk

import (
	"context"
	"fmt"
	sdk "log/slog"

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

	parent         *Handler
	fieldKeyPrefix string
	attrs          attrs
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

	fieldsMap, err := fields.AsMap(fds)
	if err != nil {
		return nil, err
	}

	return logger.NewEvent(l, fieldsMap), nil
}

func (instance *Handler) fieldsOfRecord(logger log.CoreLogger, record sdk.Record) fields.Fields {
	fdsSpec := logger.GetProvider().GetFieldKeysSpec()

	vs := make(attrs, 2+record.NumAttrs())
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
			Key:   instance.fieldKeyPrefix + v.Key,
			Value: v.Value,
		}
		i++
		return true
	})

	return fields.NewLineage(vs, instance.fields())
}

func (instance *Handler) fields() fields.Fields {
	if parent := instance.parent; parent != nil {
		return fields.NewLineage(instance.attrs, parent.fields())
	}
	return instance.attrs
}

func (instance *Handler) levelOfRecord(record sdk.Record) (level.Level, error) {
	return instance.getLevelMapper().FromSdk(record.Level)
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
	nvs := instance.attrs.clone()
	nvs.add(instance.fieldKeyPrefix, vs...)
	return &Handler{
		instance.Delegate,
		instance.LevelMapper,
		instance.DetectSkipFrames,
		instance,
		instance.fieldKeyPrefix,
		nvs,
	}
}

// WithGroup implements [sdk.Handler.WithGroup]
func (instance *Handler) WithGroup(key string) sdk.Handler {
	return &Handler{
		instance.Delegate,
		instance.LevelMapper,
		instance.DetectSkipFrames,
		instance,
		instance.fieldKeyPrefix + key + ".",
		nil,
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
