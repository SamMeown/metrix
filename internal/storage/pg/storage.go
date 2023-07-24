package pg

import (
	"context"
	"database/sql"
	"errors"
	"github.com/SamMeown/metrix/internal/storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"net"
	"time"
)

type Storage struct {
	conn *sql.DB
}

func NewStorage(conn *sql.DB) *Storage {
	return &Storage{conn: conn}
}

type requestExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (s Storage) Bootstrap(ctx context.Context) error {
	tr, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tr.Rollback()

	_, err = tr.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS gauges (
		    name VARCHAR(255) PRIMARY KEY,
		    value DOUBLE PRECISION,
		    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	_, err = tr.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS counters (
		    name VARCHAR(255) PRIMARY KEY,
		    value BIGINT,
		    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	return tr.Commit()
}

func (s Storage) getGauge(re requestExecutor, name string) (*float64, error) {
	row := re.QueryRowContext(
		context.TODO(),
		"SELECT value FROM gauges WHERE name = $1;",
		name,
	)

	var value float64
	err := row.Scan(&value)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &value, nil
}

func (s Storage) GetGauge(name string) (*float64, error) {
	return s.getGauge(s.conn, name)
}

func (s Storage) getCounter(re requestExecutor, name string) (*int64, error) {
	row := re.QueryRowContext(
		context.TODO(),
		"SELECT value FROM counters WHERE name = $1;",
		name,
	)

	var value int64
	err := row.Scan(&value)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &value, nil
}

func (s Storage) GetCounter(name string) (*int64, error) {
	return s.getCounter(s.conn, name)
}

func (s Storage) GetMany(names storage.MetricsStorageKeys) (storage.MetricsStorageItems, error) {
	tr, err := s.conn.BeginTx(
		context.TODO(),
		&sql.TxOptions{Isolation: sql.LevelRepeatableRead},
	)
	if err != nil {
		return storage.MetricsStorageItems{}, err
	}
	defer tr.Rollback()

	rv := storage.MetricsStorageItems{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
	for _, name := range names.Gauges {
		gauge, err := s.getGauge(tr, name)
		if err != nil {
			return storage.MetricsStorageItems{}, err
		}
		if gauge == nil {
			continue
		}
		rv.Gauges[name] = *gauge
	}
	for _, name := range names.Counters {
		counter, err := s.getCounter(tr, name)
		if err != nil {
			return storage.MetricsStorageItems{}, err
		}
		if counter == nil {
			continue
		}
		rv.Counters[name] = *counter
	}

	if err := tr.Commit(); err != nil {
		return storage.MetricsStorageItems{}, err
	}

	return rv, nil
}

func (s Storage) GetAll() (storage.MetricsStorageItems, error) {
	rv := storage.MetricsStorageItems{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}

	rows, err := s.conn.QueryContext(context.TODO(), `
		SELECT name, value as gauge, NULL as counter 
		FROM gauges 
		UNION ALL 
		SELECT name, NULL as gauge, value as counter 
		FROM counters;
	`)
	if err != nil {
		return storage.MetricsStorageItems{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name    string
			gauge   sql.NullFloat64
			counter sql.NullInt64
		)
		err = rows.Scan(&name, &gauge, &counter)
		if err != nil {
			return storage.MetricsStorageItems{}, err
		}

		if gauge.Valid {
			rv.Gauges[name] = gauge.Float64
		} else if counter.Valid {
			rv.Counters[name] = counter.Int64
		}
	}

	if err = rows.Err(); err != nil {
		return storage.MetricsStorageItems{}, err
	}

	return rv, nil
}

func (s Storage) setGauge(re requestExecutor, name string, value float64) error {
	_, err := re.ExecContext(
		context.TODO(),
		`
			INSERT INTO gauges (name, value, updated_at) 
			VALUES ($1, $2, $3) 
			ON CONFLICT(name) DO UPDATE SET value = $2, updated_at = $3;
		`,
		name,
		value,
		time.Now(),
	)

	return err
}

func (s Storage) setCounter(re requestExecutor, name string, value int64) error {
	_, err := re.ExecContext(
		context.TODO(),
		`
			INSERT INTO counters (name, value, updated_at) 
			VALUES ($1, $2, $3) 
			ON CONFLICT(name) DO UPDATE SET value = counters.value + $2, updated_at = $3;
		`,
		name,
		value,
		time.Now(),
	)

	return err
}

func (s Storage) SetGauge(name string, value float64) error {
	return s.setGauge(s.conn, name, value)
}

func (s Storage) SetCounter(name string, value int64) error {
	return s.setCounter(s.conn, name, value)
}

func (s Storage) SetMany(items storage.MetricsStorageItems) error {
	tr, err := s.conn.BeginTx(context.TODO(), nil)
	if err != nil {
		return err
	}
	defer tr.Rollback()

	for name, value := range items.Gauges {
		err := s.setGauge(tr, name, value)
		if err != nil {
			return err
		}
	}

	for name, value := range items.Counters {
		err := s.setCounter(tr, name, value)
		if err != nil {
			return err
		}
	}

	return tr.Commit()
}

func IsRetryableError(err error) bool {
	var pgErr *pgconn.PgError
	var netErr *net.OpError
	if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
		return true
	} else if errors.As(err, &netErr) {
		return true
	}

	return false
}
