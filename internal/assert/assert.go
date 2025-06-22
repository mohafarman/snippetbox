package assert

import (
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	// INFO: Indicated to the go test runner that this is a helper function
	// so that eg. t.Errorf reports the filename and line of the calling func
	t.Helper()

	if actual != expected {
		t.Errorf("got %v, expected %v", actual, expected)
	}
}
