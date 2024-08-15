//go:build slf4gcompat || go1.23

package testlog

func (instance *coreLogger) logLogDepth(str string, _ uint16) {
	instance.tb.Log(str)
}
