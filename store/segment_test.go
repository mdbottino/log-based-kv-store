package store

import (
	"reflect"
	"testing"

	"github.com/mdbottino/log-based-kv-store/mocks"
)

func TestNewSegment(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()
	defer mockFs.Clear()

	segment := NewSegment("./data", mockFs)
	expectedType := "store.Segment"

	if reflect.TypeOf(segment).String() != expectedType {
		t.Fatalf("segment is of the wrong type, %s", expectedType)
	}
}

func TestEmptySegmentSize(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()
	defer mockFs.Clear()

	segment := NewSegment("./data", mockFs)
	size, err := segment.Size()
	if err != nil || size > 0 {
		t.Fatalf("segment retrieved an invalid size")
	}
}

func TestSegmentSize(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()
	defer mockFs.Clear()

	segment := NewSegment("./data", mockFs)
	segment.Append("key", "value")

	// 3 chars for the key, 1 for the ':', 5 for the value and 1 for the newline
	expected := 10

	size, _ := segment.Size()
	if size != expected {
		t.Fatalf("segment retrieved the wrong size")
	}
}
