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
	s.b.Retry(func() error {
		gauge, err = s.s.GetGauge(ctx, name)
		return err
	})

	return
}

func (s Storage) GetCounter(ctx context.Context, name string) (counter *int64, err error) {
	s.b.Retry(func() error {
		counter, err = s.s.GetCounter(ctx, name)
		return err
	})

	return
}

func (s Storage) GetMany(ctx context.Context, names storage.MetricsStorageKeys) (items storage.MetricsStorageItems, err error) {
	s.b.Retry(func() error {
		items, err = s.s.GetMany(ctx, names)
		return err
	})

	return
}

func (s Storage) GetAll(ctx context.Context) (items storage.MetricsStorageItems, err error) {
	s.b.Retry(func() error {
		items, err = s.s.GetAll(ctx)
		return err
	})

	return
}

func (s Storage) SetGauge(ctx context.Context, name string, value float64) error {
	return s.b.Retry(func() error {
		return s.s.SetGauge(ctx, name, value)
	})
}

func (s Storage) SetCounter(ctx context.Context, name string, value int64) error {
	return s.b.Retry(func() error {
		return s.s.SetCounter(ctx, name, value)
	})
}

func (s Storage) SetMany(ctx context.Context, items storage.MetricsStorageItems) error {
	return s.b.Retry(func() error {
		return s.s.SetMany(ctx, items)
	})
}
