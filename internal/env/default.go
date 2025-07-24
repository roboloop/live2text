//go:build !debug

package env

func IsDebugMode() bool {
	return false
}
