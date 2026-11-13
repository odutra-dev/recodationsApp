package storage

import (
	"context"
	"io"
)

// UploadResult contains metadata about the uploaded object.
type UploadResult struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	URL    string `json:"url"`
	Size   int64  `json:"size"`
}

// StorageService defines the contract for blob/file storage operations.
type StorageService interface {
	UploadFile(ctx context.Context, bucket string, filename string, contentType string, reader io.Reader) (*UploadResult, error)
	DeleteFile(ctx context.Context, bucket string, filename string) error
	GetFileURL(bucket string, filename string) string
	EnsureBucket(ctx context.Context, bucket string) error
}
