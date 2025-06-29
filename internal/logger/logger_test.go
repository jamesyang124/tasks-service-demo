package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet_Bootstrap(t *testing.T) {
	// Test that logger can be retrieved
	logger := Get()
	assert.NotNil(t, logger)

	// Test that basic logging works without panicking
	assert.NotPanics(t, func() {
		logger.Info("test message")
		logger.Error("test error")
	})
}
