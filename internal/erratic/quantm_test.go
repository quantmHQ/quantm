package erratic_test

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.breu.io/quantm/internal/erratic"
)

func TestQuantmError_Hints(t *testing.T) {
	// Test adding various types of hints
	err := erratic.NewBadRequestError(erratic.CommonModule, "test error").
		AddHint("string_key", "string_value").
		AddHint("int_key", 123).
		AddHint("bool_key", true).
		AddHint("map_key", map[string]string{"foo": "bar"})

	// Verify LogValue

	lv := err.LogValue()
	assert.Equal(t, slog.KindGroup, lv.Kind())

	// Verify ToProto (should stringify hints)
	protoStatus := err.ToProto()
	require.NotNil(t, protoStatus)

	// Verify ToConnectError (should stringify hints)
	connectErr := err.ToConnectError()
	require.NotNil(t, connectErr)
	assert.Equal(t, "string_value", connectErr.Meta().Get("string_key"))
	assert.Equal(t, "123", connectErr.Meta().Get("int_key"))
	assert.Equal(t, "true", connectErr.Meta().Get("bool_key"))
}

func TestQuantmError_StackHint(t *testing.T) {
	err := erratic.New(erratic.CommonModule, erratic.CodeSystem, "system error")

	// Add manual stack hint

	_ = err.WithStack("fake stack trace")

	// Should not panic in ToProto due to type assertion

	protoStatus := err.ToProto()
	require.NotNil(t, protoStatus)
}

func TestQuantmError_New(t *testing.T) {
	// Just verify New creates an error without panic
	err := erratic.New(erratic.CommonModule, erratic.CodeUnknown, "unknown")
	assert.NotNil(t, err)
	assert.Equal(t, "unknown", err.Message)
}
