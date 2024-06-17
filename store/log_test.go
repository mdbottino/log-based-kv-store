package store

import (
	"reflect"
	"testing"

	"github.com/mdbottino/log-based-kv-store/mocks"
)

func TestNewLog(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()
	defer mockFs.Clear()

	log := NewLog("./data", mockFs)
	expectedType := "store.Log"

	if reflect.TypeOf(log).String() != expectedType {
		t.Fatalf("log is of the wrong type, %s", expectedType)
	}
}

func TestLogAppendEmpty(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()
	defer mockFs.Clear()

	log := NewLog("./data", mockFs)
	log.Append("key", "value")

	// 3 chars for the key, 1 for the ':', 5 for the value and 1 for the newline
	expected := 10
	if mockFs.GetHandle(0).GetSize() != expected {
		t.Fatalf("the log attempted to write unexpected data")
	}
}

func TestLogAppend(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()
	defer mockFs.Clear()

	log := NewLog("./data", mockFs)
	log.Append("key", "value")
	log.Append("more", "data")

	// 10 chars for the first key-value pair and 10 for the other one
	expected := 20
	if mockFs.GetHandle(0).GetSize() != expected {
		t.Fatalf("the log attempted to write unexpected data")
	}
}

func TestLogFindKeyNotPresent(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()
	defer mockFs.Clear()

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
	defer mockFs.Clear()

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

func TestLogFindExistingKeyWithMultipleEntries(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()
	defer mockFs.Clear()

	log := NewLog("./data", mockFs)
	log.Append("key", "value")
	log.Append("more", "data")
	log.Append("key", "updated value")

	val, err := log.Find("key")
	if err != nil {
		t.Fatalf("failed to retrieve a key from the log")
	}

	if val != "updated value" {
		t.Fatalf("retrieved the wrong value from the log")
	}
}

func TestAddSegment(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()
	defer mockFs.Clear()

	log := NewLog("./data", mockFs)
	log.AddSegment()

	if len(log.segments) != 2 {
		t.Fatalf("failed to create another segment")
	}
}

func TestFindKeyOnTheNewestSegment(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()
	defer mockFs.Clear()

	log := NewLog("./data", mockFs)
	log.Append("key", "value")
	log.Append("more", "data")
	log.AddSegment()
	log.Append("key", "new data")

	if len(log.segments) != 2 {
		t.Fatalf("failed to create another segment")
	}

	val, err := log.Find("key")
	if err != nil {
		t.Fatalf("an error occurred while finding the key")
	}

	if val != "new data" {
		t.Fatalf("retrieved the wrong value")
	}
}

func TestFindKeyOnTheOldestSegment(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()
	defer mockFs.Clear()

	log := NewLog("./data", mockFs)
	log.Append("key", "value")
	log.Append("more", "data")
	log.AddSegment()
	log.Append("other", "stuff")

	if len(log.segments) != 2 {
		t.Fatalf("failed to create another segment")
	}

	val, err := log.Find("key")
	if err != nil {
		t.Fatalf("an error occurred while finding the key")
	}

	if val != "value" {
		t.Fatalf("retrieved the wrong value")
	}
}

func TestFindKeyNotPresentWithMoreThanOnceSegment(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()
	defer mockFs.Clear()

	log := NewLog("./data", mockFs)
	log.Append("key", "value")
	log.AddSegment()
	log.Append("other", "stuff")

	if len(log.segments) != 2 {
		t.Fatalf("failed to create another segment")
	}

	_, err := log.Find("not.here")
	if err == nil {
		t.Fatalf("retrieved wrong value from the log")
	}
}

func TestGetLastSegmentOneSegment(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()
	defer mockFs.Clear()

	log := NewLog("./data", mockFs)
	log.Append("key", "value")

	segment := log.GetLatestSegment()

	val, _ := segment.Find("key")
	if val != "value" {
		t.Fatalf("retrieved wrong value from the log")
	}

	_, err := segment.Find("other")
	if err == nil {
		t.Fatalf("retrieved wrong value from the log")
	}
}

func TestGetLastSegmentMoreThanOneSegment(t *testing.T) {
	mockFs := mocks.NewMockFileSystem()
	defer mockFs.Clear()

	log := NewLog("./data", mockFs)
	log.Append("key", "value")
	log.AddSegment()
	log.Append("other", "stuff")

	segment := log.GetLatestSegment()

	val, _ := segment.Find("other")
	if val != "stuff" {
		t.Fatalf("retrieved wrong value from the log")
	}

	_, err := segment.Find("key")
	if err == nil {
		t.Fatalf("retrieved wrong value from the log")
	}
}
