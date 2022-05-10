package log

import "fmt"

/// Simple logger that can be enabled/disabled based on verbose flag
type Logger struct {
	VerboseEnabled bool
}

func CreateLogger() *Logger {
	return &Logger{
		VerboseEnabled: false,
	}
}

func (self *Logger) Verbose(message string) {
	if self.VerboseEnabled {
		fmt.Println(message)
	}
}
