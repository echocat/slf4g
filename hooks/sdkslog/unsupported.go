//go:build !go1.21
// +build !go1.21

package hook_sdkslog

func init() {
	panic("This feature requires Go 1.21+.")
}
