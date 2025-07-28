package api

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/object-storage-service/domain"
	"github.com/yourusername/object-storage-service/persistence"
)

func setupTestServer() (*Server, domain.Storage) {
	storage := persistence.NewInMemoryStorage()
	server := NewServer(storage, "8080")
	return server, storage
}

func TestPutObject(t *testing.T) {
	server, storage := setupTestServer()
	ts := httptest.NewServer(server.router)
	defer ts.Close()

	url := ts.URL + "/objects/testbucket/testobject"
	payload := []byte("test object data")

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created; got %d", resp.StatusCode)
	}

	storedData, err := storage.Get("testbucket", "testobject")
	if err != nil {
		t.Fatalf("expected object to be stored, got error: %v", err)
	}
	if string(storedData) != string(payload) {
		t.Errorf("stored data mismatch; expected %q got %q", payload, storedData)
	}
}

func TestGetObject(t *testing.T) {
	server, storage := setupTestServer()
	ts := httptest.NewServer(server.router)
	defer ts.Close()

	// Pre-store an object
	bucket, objectID := "testbucket", "testobject"
	content := []byte("hello world")
	if _, err := storage.Put(bucket, objectID, content); err != nil {
		t.Fatalf("failed to store object: %v", err)
	}

	url := ts.URL + "/objects/" + bucket + "/" + objectID

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("could not send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK; got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	if !bytes.Equal(body, content) {
		t.Errorf("expected body %q; got %q", content, body)
	}

	// Test object not found case
	notFoundURL := ts.URL + "/objects/" + bucket + "/nonexistent"
	respNF, err := http.Get(notFoundURL)
	if err != nil {
		t.Fatalf("could not send GET request for not found: %v", err)
	}
	defer respNF.Body.Close()

	if respNF.StatusCode != http.StatusNotFound && respNF.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 or 404 for not found; got %d", respNF.StatusCode)
	}
}

func TestDeleteObject(t *testing.T) {
	server, storage := setupTestServer()
	ts := httptest.NewServer(server.router)
	defer ts.Close()

	// Pre-store an object
	bucket, objectID := "testbucket", "testobject"
	content := []byte("to be deleted")
	if _, err := storage.Put(bucket, objectID, content); err != nil {
		t.Fatalf("failed to store object: %v", err)
	}

	url := ts.URL + "/objects/" + bucket + "/" + objectID

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatalf("could not create DELETE request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("could not send DELETE request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK on delete; got %d", resp.StatusCode)
	}

	// Verify deletion
	_, err = storage.Get(bucket, objectID)
	if err == nil {
		t.Errorf("expected object to be deleted, but still found")
	}

	// Delete nonexistent object
	reqNF, err := http.NewRequest(http.MethodDelete, ts.URL+"/objects/"+bucket+"/nonexistent", nil)
	if err != nil {
		t.Fatalf("could not create DELETE request for nonexistent: %v", err)
	}

	respNF, err := http.DefaultClient.Do(reqNF)
	if err != nil {
		t.Fatalf("could not send DELETE request for nonexistent: %v", err)
	}
	defer respNF.Body.Close()

	if respNF.StatusCode != http.StatusNotFound && respNF.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 or 404 for deleting nonexistent; got %d", respNF.StatusCode)
	}
}

func TestPut_DeDuplicateAndOverwrite(t *testing.T) {
	server, storage := setupTestServer()
	ts := httptest.NewServer(server.router)
	defer ts.Close()

	// First insert — should succeed
	ok, err := storage.Put("bucket1", "obj1", []byte("original"))
	if err != nil || !ok {
		t.Fatalf("expected first insert to succeed, got err: %v", err)
	}

	// Duplicate insert with same content — should deduplicate silently (no error, ok = false)
	ok, err = storage.Put("bucket1", "obj1", []byte("original"))
	if err != nil {
		t.Fatalf("expected no error on duplicate data, got: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false on duplicate data, got true")
	}

	// Insert with same ID but different content — should overwrite successfully
	ok, err = storage.Put("bucket1", "obj1", []byte("updated"))
	if err != nil || !ok {
		t.Fatalf("expected overwrite to succeed, got err: %v", err)
	}

	// Validate content was updated
	data, err := storage.Get("bucket1", "obj1")
	if err != nil {
		t.Fatalf("expected get to succeed, got err: %v", err)
	}
	if !bytes.Equal(data, []byte("updated")) {
		t.Fatalf("expected updated content, got: %s", data)
	}
}
