//go:build slf4gcompat

package testlog

func (instance *coreLogger) logLogDepth(str string, _ uint16) {
	instance.tb.Log(str)
}
