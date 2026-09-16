package tracing

import (
	"context"
	"testing"
)

// Start must work before Initialize: callers' unit tests never set up an
// exporter, and the span they get back must still be usable.
func TestStartBeforeInitialize(t *testing.T) {
	span := Start(context.Background(), "test")
	if span == nil {
		t.Fatal("Start returned a nil span")
	}
	span.End()
}
