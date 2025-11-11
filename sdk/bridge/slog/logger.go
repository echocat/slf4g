//go:build go1.21

package sdk

import (
	sdk "log/slog"

	log "github.com/echocat/slf4g"
)

// New creates an instance SDK's [log/slog.Logger] using the given target logger.
// If the target logger is nil, the result of [log.GetRootLogger] will be used.
func New(target log.CoreLogger, customizer ...func(*Handler)) *sdk.Logger {
	h := NewHandler(target, customizer...)
	return sdk.New(h)
}

// Configure configures the SDK's [log/slog] framework to use slf4g with the result
// of [log.GetRootLogger].
func Configure(customizer ...func(*Handler)) {
	ConfigureWith(log.GetRootLogger(), customizer...)
}

// ConfigureWith configures the SDK's [log/slog] framework to use slf4g with the
// given target Logger.
// If the target logger is nil, the result of [log.GetRootLogger] will be used.
func ConfigureWith(target log.CoreLogger, customizer ...func(*Handler)) {
	logger := New(target, customizer...)
	sdk.SetDefault(logger)
}
