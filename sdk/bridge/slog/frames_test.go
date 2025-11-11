//go:build go1.21

package sdk_test

import (
	"context"
	sdk "log/slog"
	"testing"

	"github.com/echocat/slf4g/internal/test/assert"
	bridge "github.com/echocat/slf4g/sdk/bridge/slog"
)

func Test_detectSkipFramesFromSdk(t *testing.T) {
	handler := &dummyHandler{detectSkipFramesUsing: bridge.DefaultDetectSkipFrames}

	sl := sdk.New(handler)
	sl.Log(context.Background(), bridge.LevelInfo, "info")
	assert.ToBeEqual(t, uint16(3), handler.detectedSkipFrames)
}

type dummyHandler struct {
	detectedSkipFrames    uint16
	detectSkipFramesUsing bridge.DetectSkipFrames
}

func (instance *dummyHandler) Enabled(_ context.Context, _ sdk.Level) bool {
	return true
}

func (instance *dummyHandler) Handle(_ context.Context, _ sdk.Record) error {
	instance.detectedSkipFrames = instance.detectSkipFramesUsing(1)
	return nil
}

func (instance *dummyHandler) WithAttrs(_ []sdk.Attr) sdk.Handler {
	panic("should never be called")
}

func (instance *dummyHandler) WithGroup(_ string) sdk.Handler {
	panic("should never be called")
}
