package persistence

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/DanielePalaia/object-storage-service/domain"
)

func TestInMemoryStorage_PutGetDelete(t *testing.T) {
	storage := NewInMemoryStorage()

	bucket := "bucket1"
	objectID := "obj1"
	data := []byte("hello world")

	// Test Put
	if _, err := storage.Put(bucket, objectID, data); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Test Get - should return data
	got, err := storage.Get(bucket, objectID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("Get returned wrong data: got %q want %q", got, data)
	}

	// Test Delete
	if err := storage.Delete(bucket, objectID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Test Get after Delete - should return error
	_, err = storage.Get(bucket, objectID)
	if err == nil {
		t.Errorf("expected error getting deleted object, got nil")
	}
}

func TestInMemoryStorage_Deduplication(t *testing.T) {
	storage := NewInMemoryStorage()

	bucket := "bucket1"
	objectID := "obj1"
	data1 := []byte("data1")
	data2 := []byte("data1") // same content, should deduplicate (no error)
	data3 := []byte("data2") // different content, should overwrite

	// Put first object
	if _, err := storage.Put(bucket, objectID, data1); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Put same object again with identical data, should succeed silently
	if _, err := storage.Put(bucket, objectID, data2); err != nil {
		t.Fatalf("Put failed on duplicate data: %v", err)
	}

	// Verify stored data unchanged
	got, err := storage.Get(bucket, objectID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !bytes.Equal(got, data1) {
		t.Errorf("Get returned wrong data after deduplication: got %q want %q", got, data1)
	}

	// Put same object with different data, should overwrite
	if _, err := storage.Put(bucket, objectID, data3); err != nil {
		t.Fatalf("Put failed on overwrite: %v", err)
	}

	// Verify data updated
	got, err = storage.Get(bucket, objectID)
	if err != nil {
		t.Fatalf("Get failed after overwrite: %v", err)
	}
	if !bytes.Equal(got, data3) {
		t.Errorf("Get returned wrong data after overwrite: got %q want %q", got, data3)
	}
}

func TestInMemoryStorage_ConcurrentAccess(t *testing.T) {
	storage := NewInMemoryStorage()
	bucket := "concurrent-bucket"

	var wg sync.WaitGroup
	numOps := 100

	for i := 0; i < numOps; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			objectID := fmt.Sprintf("obj-%d", i)
			data := []byte(fmt.Sprintf("data-%d", i))
			_, err := storage.Put(bucket, objectID, data)
			if err != nil && err != domain.ErrAlreadyExist {
				t.Errorf("unexpected error during concurrent Put: %v", err)
			}
		}(i)
	}

	wg.Wait()

	// Validate all objects exist
	for i := 0; i < numOps; i++ {
		objectID := fmt.Sprintf("obj-%d", i)
		_, err := storage.Get(bucket, objectID)
		if err != nil {
			t.Errorf("missing object after concurrent writes: %s", objectID)
		}
	}
}
