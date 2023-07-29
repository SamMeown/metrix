package retryable

import (
	"context"
	"github.com/SamMeown/metrix/internal/backoff"
	"github.com/SamMeown/metrix/internal/storage"
)

type Storage struct {
	s storage.MetricsStorage
	b backoff.Backoff
}

func NewStorage(s storage.MetricsStorage, retryableErrFunc func(error) bool) *Storage {
	return &Storage{
		s: s,
		b: backoff.NewBackoff([]int{1, 3, 5}, retryableErrFunc),
	}
}

func (s Storage) GetGauge(ctx context.Context, name string) (gauge *float64, err error) {
	err = s.b.RetryContext(ctx, func() (e error) {
		gauge, e = s.s.GetGauge(ctx, name)
		return
	})

	return
}

func (s Storage) GetCounter(ctx context.Context, name string) (counter *int64, err error) {
	err = s.b.RetryContext(ctx, func() (e error) {
		counter, e = s.s.GetCounter(ctx, name)
		return
	})

	return
}

func (s Storage) GetMany(ctx context.Context, names storage.MetricsStorageKeys) (items storage.MetricsStorageItems, err error) {
	err = s.b.RetryContext(ctx, func() (e error) {
		items, e = s.s.GetMany(ctx, names)
		return
	})

	return
}

func (s Storage) GetAll(ctx context.Context) (items storage.MetricsStorageItems, err error) {
	err = s.b.RetryContext(ctx, func() (e error) {
		items, e = s.s.GetAll(ctx)
		return
	})

	return
}

func (s Storage) SetGauge(ctx context.Context, name string, value float64) error {
	return s.b.RetryContext(ctx, func() error {
		return s.s.SetGauge(ctx, name, value)
	})
}

func (s Storage) SetCounter(ctx context.Context, name string, value int64) error {
	return s.b.RetryContext(ctx, func() error {
		return s.s.SetCounter(ctx, name, value)
	})
}

func (s Storage) SetMany(ctx context.Context, items storage.MetricsStorageItems) error {
	return s.b.RetryContext(ctx, func() error {
		return s.s.SetMany(ctx, items)
	})
}
