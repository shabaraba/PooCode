package logger

import "fmt"

// TokenDebug logs token processing information for debugging
// This is simplified to not depend on global state
func TokenDebug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("[TOKEN-DEBUG] %s\n", msg)
}
