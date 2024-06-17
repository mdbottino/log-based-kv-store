package store

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig([]byte{})
	expectedType := "store.Config"

	if reflect.TypeOf(config).String() != expectedType {
		t.Fatalf("config is of the wrong type, %s", expectedType)
	}

	if config.Segments.MaxSize != 0 {
		t.Fatalf("config does not have the correct default value for segments.maxSize")
	}
}

func TestNewConfigWithValues(t *testing.T) {
	config := NewConfig([]byte("segments:\n  maxSize: 2000\n"))

	expected := 2000

	if config.Segments.MaxSize != expected {
		t.Fatalf("config did not load properly the segments.maxSize value")
	}
}
