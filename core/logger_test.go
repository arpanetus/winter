package core

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger("test")
	logger.Info("Works gud info")
	logger.Infof("Works gud %s", "some value")
	logger.Warn("Another warn log and color")
	logger.Err("And some red for errors")
	logger.Note("Noting...")
}
