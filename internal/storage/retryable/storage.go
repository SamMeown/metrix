package retryable

import (
	"errors"
	"github.com/SamMeown/metrix/internal/backoff"
	"github.com/SamMeown/metrix/internal/storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"net"
)

type Storage struct {
	s storage.MetricsStorage
	b backoff.Backoff
}

func NewStorage(s storage.MetricsStorage) *Storage {
	return &Storage{
		s: s,
		b: backoff.NewBackoff([]int{1, 3, 5}, isRetryableError),
	}
}

func isRetryableError(err error) bool {
	var pgErr *pgconn.PgError
	var netErr *net.OpError
	if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
		return true
	} else if errors.As(err, &netErr) {
		return true
	}

	return false
}

func (s Storage) GetGauge(name string) (gauge *float64, err error) {
	s.b.Retry(func() error {
		gauge, err = s.s.GetGauge(name)
		return err
	})

	return
}

func (s Storage) GetCounter(name string) (counter *int64, err error) {
	s.b.Retry(func() error {
		counter, err = s.s.GetCounter(name)
		return err
	})

	return
}

func (s Storage) GetMany(names storage.MetricsStorageKeys) (items storage.MetricsStorageItems, err error) {
	s.b.Retry(func() error {
		items, err = s.s.GetMany(names)
		return err
	})

	return
}

func (s Storage) GetAll() (items storage.MetricsStorageItems, err error) {
	s.b.Retry(func() error {
		items, err = s.s.GetAll()
		return err
	})

	return
}

func (s Storage) SetGauge(name string, value float64) error {
	return s.b.Retry(func() error {
		return s.s.SetGauge(name, value)
	})
}

func (s Storage) SetCounter(name string, value int64) error {
	return s.b.Retry(func() error {
		return s.s.SetCounter(name, value)
	})
}

func (s Storage) SetMany(items storage.MetricsStorageItems) error {
	return s.b.Retry(func() error {
		return s.s.SetMany(items)
	})
}
