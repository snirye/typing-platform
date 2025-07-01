package core

import (
	"log"
	"os"
)

// Logger wraps a standard logger for file logging
// This abstraction allows us to swap out logging implementations if needed.
type Logger struct {
	*log.Logger
}

// NewLogger creates a new logger that writes to the specified file path.
func NewLogger(filePath string) (*Logger, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	logger := log.New(file, "[ascii-type] ", log.LstdFlags|log.Lshortfile)
	return &Logger{logger}, nil
}
