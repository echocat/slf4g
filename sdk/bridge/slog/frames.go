package sdk

type DetectSkipFrames func() uint16

var DefaultDetectSkipFrames DetectSkipFrames = func() uint16 {
	return 0
}
