//go:build !go1.21

package sdk

func init() {
	panic("This feature requires Go 1.21+.")
}
