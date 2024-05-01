package store

import (
	"reflect"
	"testing"

	"github.com/mdbottino/log-based-kv-store/mocks"
)

func TestNewStore(t *testing.T) {
	store := NewStore("./data", mocks.MockFileSystem{})
	expectedType := "store.Store"

	if reflect.TypeOf(store).String() != expectedType {
		t.Fatalf("store is of the wrong type, %s", expectedType)
	}
}

func TestStoreSet(t *testing.T) {
	store := NewStore("./data", mocks.MockFileSystem{})
	err := store.Set("key", "value")
	if err != nil {
		t.Fatalf("failed to store a key in the store")
	}
}

func TestStoreGetEmptyStore(t *testing.T) {
	store := NewStore("./data", mocks.MockFileSystem{})

	_, err := store.Get("key")
	if err == nil {
		t.Fatalf("incorrectly retrieved the wrong value")
	}
}

func TestStoreGetAfterSet(t *testing.T) {
	store := NewStore("./data", mocks.MockFileSystem{})
	key := "key"
	value := "value"

	err := store.Set(key, value)
	if err != nil {
		t.Fatalf("failed to store a key in the store")
	}

	err = store.Set("some other key", "some other value")
	if err != nil {
		t.Fatalf("failed to store a key in the store")
	}

	result, err := store.Get(key)
	if err != nil {
		t.Fatalf("failed to retrieve the key from the store")
	}

	if result != value {
		t.Fatalf("failed to retrieve the right value for the provided key")
	}
}