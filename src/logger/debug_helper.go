package logger

import "fmt"

// LogTokenProcessing logs token processing information for debugging
func LogTokenProcessing(format string, args ...interface{}) {
	if CurrentLevel <= LevelDebug {
		msg := fmt.Sprintf(format, args...)
		fmt.Printf("[TOKEN-DEBUG] %s\n", msg)
	}
}
