package store

import (
	"reflect"
	"testing"

	"github.com/mdbottino/log-based-kv-store/mocks"
)

func TestNewLog(t *testing.T) {
	log := NewLog("./data", mocks.NewMockFileSystem())
	expectedType := "store.Log"

	if reflect.TypeOf(log).String() != expectedType {
		t.Fatalf("log is of the wrong type, %s", expectedType)
	}
}

func TestLogAppendEmpty(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()

	log := NewLog("./data", mockFs)
	log.Append("key", "value")

	// 3 chars for the key, 1 for the ':', 5 for the value and 1 for the newline
	expected := 10
	if mockFs.Handle.GetSize() != expected {
		t.Fatalf("the log attempted to write unexpected data")
	}
}

func TestLogAppend(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()

	log := NewLog("./data", mockFs)
	log.Append("key", "value")
	log.Append("more", "data")

	// 10 chars for the first key-value pair and 10 for the other one
	expected := 20
	if mockFs.Handle.GetSize() != expected {
		t.Fatalf("the log attempted to write unexpected data")
	}
}

func TestLogFindKeyNotPresent(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()

	log := NewLog("./data", mockFs)
	log.Append("key", "value")
	log.Append("more", "data")

	_, err := log.Find("banana")
	if err == nil {
		t.Fatalf("retrieved the wrong value from the log")
	}

}

func TestLogFindExistingKey(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()

	log := NewLog("./data", mockFs)
	log.Append("key", "value")
	log.Append("more", "data")

	val, err := log.Find("key")
	if err != nil {
		t.Fatalf("failed to retrieve a key from the log")
	}

	if val != "value" {
		t.Fatalf("retrieved the wrong value from the log")
	}
}
