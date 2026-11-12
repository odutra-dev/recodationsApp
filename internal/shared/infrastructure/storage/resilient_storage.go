package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"time"

	"github.com/sony/gobreaker"
)

// ResilientStorage wraps any StorageService with a Circuit Breaker and Exponential Backoff Retry.
type ResilientStorage struct {
	next           StorageService
	circuitBreaker *gobreaker.CircuitBreaker
	maxRetries     int
	baseBackoff    time.Duration
	maxBackoff     time.Duration
}

// NewResilientStorage wraps the underlying StorageService with resilience mechanisms.
func NewResilientStorage(next StorageService) *ResilientStorage {
	settings := gobreaker.Settings{
		Name:        "GCSStorageCircuitBreaker",
		MaxRequests: 3,
		Interval:    10 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			slog.Warn("Circuit Breaker state changed", "name", name, "from", from.String(), "to", to.String())
		},
	}

	return &ResilientStorage{
		next:           next,
		circuitBreaker: gobreaker.NewCircuitBreaker(settings),
		maxRetries:     3,
		baseBackoff:    100 * time.Millisecond,
		maxBackoff:     2 * time.Second,
	}
}

func (r *ResilientStorage) EnsureBucket(ctx context.Context, bucket string) error {
	_, err := r.circuitBreaker.Execute(func() (any, error) {
		return nil, r.retry(ctx, func() error {
			return r.next.EnsureBucket(ctx, bucket)
		})
	})
	return err
}

func (r *ResilientStorage) UploadFile(ctx context.Context, bucket string, filename string, contentType string, reader io.Reader) (*UploadResult, error) {
	result, err := r.circuitBreaker.Execute(func() (any, error) {
		var res *UploadResult
		err := r.retry(ctx, func() error {
			var rErr error
			res, rErr = r.next.UploadFile(ctx, bucket, filename, contentType, reader)
			return rErr
		})
		return res, err
	})

	if err != nil {
		return nil, fmt.Errorf("resilient upload failed: %w", err)
	}

	return result.(*UploadResult), nil
}

func (r *ResilientStorage) DeleteFile(ctx context.Context, bucket string, filename string) error {
	_, err := r.circuitBreaker.Execute(func() (any, error) {
		return nil, r.retry(ctx, func() error {
			return r.next.DeleteFile(ctx, bucket, filename)
		})
	})
	return err
}

func (r *ResilientStorage) GetFileURL(bucket string, filename string) string {
	return r.next.GetFileURL(bucket, filename)
}

// retry executes the operation with exponential backoff and jitter.
func (r *ResilientStorage) retry(ctx context.Context, op func() error) error {
	var err error
	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = op()
		if err == nil {
			return nil
		}

		if attempt == r.maxRetries {
			break
		}

		// Calculate exponential backoff with full jitter
		backoff := r.baseBackoff * time.Duration(1<<attempt)
		if backoff > r.maxBackoff {
			backoff = r.maxBackoff
		}
		jitter := time.Duration(rand.Int63n(int64(backoff / 2)))
		sleepDuration := backoff/2 + jitter

		slog.Warn("Storage operation failed, retrying...", "attempt", attempt+1, "error", err, "backoff", sleepDuration)
		select {
		case <-time.After(sleepDuration):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return err
}
