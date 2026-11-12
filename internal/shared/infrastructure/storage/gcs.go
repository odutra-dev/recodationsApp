package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// GCSStorage implements StorageService using the official Google Cloud Storage Go SDK.
// It supports both local floci-gcp emulator (via STORAGE_EMULATOR_HOST) and real GCP.
type GCSStorage struct {
	client       *storage.Client
	projectID    string
	emulatorHost string
}

// NewGCSStorage creates a new GCS storage service instance.
func NewGCSStorage(ctx context.Context, projectID string) (*GCSStorage, error) {
	emulatorHost := os.Getenv("STORAGE_EMULATOR_HOST")
	var opts []option.ClientOption

	if emulatorHost != "" {
		// When running with floci-gcp, bypass Google auth credentials
		opts = append(opts, option.WithoutAuthentication())
	}

	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GCS client: %w", err)
	}

	return &GCSStorage{
		client:       client,
		projectID:    projectID,
		emulatorHost: emulatorHost,
	}, nil
}

// EnsureBucket creates the bucket if it doesn't already exist.
func (s *GCSStorage) EnsureBucket(ctx context.Context, bucket string) error {
	bkt := s.client.Bucket(bucket)
	_, err := bkt.Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		if err := bkt.Create(ctx, s.projectID, nil); err != nil {
			// Ignore already-owned errors
			if !strings.Contains(err.Error(), "conflict") && !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("failed to create bucket %s: %w", bucket, err)
			}
		}
	}
	return nil
}

// UploadFile uploads an object to the specified bucket in GCS.
func (s *GCSStorage) UploadFile(ctx context.Context, bucket string, filename string, contentType string, reader io.Reader) (*UploadResult, error) {
	if err := s.EnsureBucket(ctx, bucket); err != nil {
		return nil, err
	}

	bkt := s.client.Bucket(bucket)
	obj := bkt.Object(filename)

	writer := obj.NewWriter(ctx)
	if contentType != "" {
		writer.ContentType = contentType
	}

	written, err := io.Copy(writer, reader)
	if err != nil {
		_ = writer.Close()
		return nil, fmt.Errorf("failed copying data to GCS object: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return &UploadResult{
		Bucket: bucket,
		Key:    filename,
		URL:    s.GetFileURL(bucket, filename),
		Size:   written,
	}, nil
}

// DeleteFile deletes an object from GCS.
func (s *GCSStorage) DeleteFile(ctx context.Context, bucket string, filename string) error {
	obj := s.client.Bucket(bucket).Object(filename)
	if err := obj.Delete(ctx); err != nil && err != storage.ErrObjectNotExist {
		return fmt.Errorf("failed to delete object %s from GCS: %w", filename, err)
	}
	return nil
}

// GetFileURL returns the URL for accessing the object.
func (s *GCSStorage) GetFileURL(bucket string, filename string) string {
	if s.emulatorHost != "" {
		host := s.emulatorHost
		if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
			host = "http://" + host
		}
		return fmt.Sprintf("%s/storage/v1/b/%s/o/%s", host, bucket, filename)
	}
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket, filename)
}
