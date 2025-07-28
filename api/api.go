package api

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yourusername/object-storage-service/domain"
)

// putObjectHandler uploads an object to a bucket.
// @Summary Upload an object
// @Description Upload an object to the specified bucket with objectID.
// @Tags objects
// @Accept application/octet-stream
// @Produce application/json
// @Param bucket path string true "Bucket name"
// @Param objectID path string true "Object ID"
// @Param data body string true "Object data"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /objects/{bucket}/{objectID} [put]
func putObjectHandler(storage domain.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bucket := vars["bucket"]
		objectID := vars["objectID"]

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Request error:", err)
			http.Error(w, "unable to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		_, err = storage.Put(bucket, objectID, data)
		if err != nil {
			log.Println("Request error:", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"` + objectID + `"}`))
	}
}

// getObjectHandler downloads an object from a bucket.
// @Summary Download an object
// @Description Download an object by bucket and objectID.
// @Tags objects
// @Produce application/octet-stream
// @Param bucket path string true "Bucket name"
// @Param objectID path string true "Object ID"
// @Success 200 {string} string "Object data"
// @Failure 404 {string} string "Not Found"
// @Router /objects/{bucket}/{objectID} [get]
func getObjectHandler(storage domain.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bucket := vars["bucket"]
		objectID := vars["objectID"]

		data, err := storage.Get(bucket, objectID)
		if err != nil {
			log.Println("Request error:", err)
			http.Error(w, "object not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

// deleteObjectHandler deletes an object from a bucket.
// @Summary Delete an object
// @Description Delete an object by bucket and objectID.
// @Tags objects
// @Param bucket path string true "Bucket name"
// @Param objectID path string true "Object ID"
// @Success 200 {string} string "Deleted"
// @Failure 404 {string} string "Not Found"
// @Router /objects/{bucket}/{objectID} [delete]
func deleteObjectHandler(storage domain.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bucket := vars["bucket"]
		objectID := vars["objectID"]

		err := storage.Delete(bucket, objectID)
		if err != nil {
			log.Println("Request error:", err)
			http.Error(w, "object not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
