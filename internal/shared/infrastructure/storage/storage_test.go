package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockStorage implements StorageService for unit testing.
type mockStorage struct {
	uploadFunc func(ctx context.Context, bucket string, filename string, contentType string, reader io.Reader) (*UploadResult, error)
	deleteFunc func(ctx context.Context, bucket string, filename string) error
	calls      int
}

func (m *mockStorage) EnsureBucket(ctx context.Context, bucket string) error {
	return nil
}

func (m *mockStorage) UploadFile(ctx context.Context, bucket string, filename string, contentType string, reader io.Reader) (*UploadResult, error) {
	m.calls++
	if m.uploadFunc != nil {
		return m.uploadFunc(ctx, bucket, filename, contentType, reader)
	}
	return &UploadResult{Bucket: bucket, Key: filename, URL: "http://mock/" + filename, Size: 10}, nil
}

func (m *mockStorage) DeleteFile(ctx context.Context, bucket string, filename string) error {
	m.calls++
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, bucket, filename)
	}
	return nil
}

func (m *mockStorage) GetFileURL(bucket string, filename string) string {
	return "http://mock/" + filename
}

func TestResilientStorageSuccess(t *testing.T) {
	mock := &mockStorage{}
	resilient := NewResilientStorage(mock)

	ctx := context.Background()
	reader := bytes.NewReader([]byte("test data"))

	res, err := resilient.UploadFile(ctx, "photos", "test.jpg", "image/jpeg", reader)
	require.NoError(t, err)
	assert.Equal(t, "http://mock/test.jpg", res.URL)
	assert.Equal(t, 1, mock.calls)
}

func TestResilientStorageRetryOnFailure(t *testing.T) {
	attempts := 0
	mock := &mockStorage{
		uploadFunc: func(ctx context.Context, bucket string, filename string, contentType string, reader io.Reader) (*UploadResult, error) {
			attempts++
			if attempts < 3 {
				return nil, errors.New("temporary network glitch")
			}
			return &UploadResult{Bucket: bucket, Key: filename, URL: "http://mock/" + filename, Size: 10}, nil
		},
	}

	resilient := NewResilientStorage(mock)
	resilient.baseBackoff = 5 * time.Millisecond
	resilient.maxBackoff = 20 * time.Millisecond

	ctx := context.Background()
	reader := bytes.NewReader([]byte("test data"))

	res, err := resilient.UploadFile(ctx, "photos", "test.jpg", "image/jpeg", reader)
	require.NoError(t, err)
	assert.Equal(t, "http://mock/test.jpg", res.URL)
	assert.Equal(t, 3, attempts, "Should have retried twice and succeeded on 3rd attempt")
}

func TestResilientStorageCircuitBreakerTrips(t *testing.T) {
	mock := &mockStorage{
		uploadFunc: func(ctx context.Context, bucket string, filename string, contentType string, reader io.Reader) (*UploadResult, error) {
			return nil, errors.New("hard backend outage")
		},
	}

	resilient := NewResilientStorage(mock)
	resilient.baseBackoff = 1 * time.Millisecond
	resilient.maxBackoff = 2 * time.Millisecond

	ctx := context.Background()

	// Cause failures to trip breaker
	for i := 0; i < 5; i++ {
		reader := bytes.NewReader([]byte("data"))
		_, _ = resilient.UploadFile(ctx, "photos", "test.jpg", "image/jpeg", reader)
	}

	// Now breaker should be open
	reader := bytes.NewReader([]byte("data"))
	_, err := resilient.UploadFile(ctx, "photos", "test.jpg", "image/jpeg", reader)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "circuit breaker is open")
}
