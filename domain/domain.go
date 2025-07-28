package domain

import "errors"

type Object struct {
	ID     string
	Bucket string
	Data   []byte
}

type Storage interface {
	Put(bucket, objectID string, data []byte) (bool, error) // true = created, false = already exists
	Get(bucket, objectID string) ([]byte, error)
	Delete(bucket, objectID string) error
}

var (
	ErrNotFound     = errors.New("object not found")
	ErrAlreadyExist = errors.New("object already exists in bucket")
)
