package storage

import (
	"context"
	"golang.org/x/exp/maps"
)

const (
	MetricsTypeGauge   = "gauge"
	MetricsTypeCounter = "counter"
)

type MetricsStorageItems struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

type MetricsStorageKeys struct {
	Gauges   []string
	Counters []string
}

type MetricsStorageGetter interface {
	GetGauge(ctx context.Context, name string) (*float64, error)
	GetCounter(ctx context.Context, name string) (*int64, error)
	GetMany(ctx context.Context, names MetricsStorageKeys) (MetricsStorageItems, error)
	GetAll(ctx context.Context) (MetricsStorageItems, error)
}

type MetricsStorageSetter interface {
	SetGauge(ctx context.Context, name string, value float64) error
	SetCounter(ctx context.Context, name string, value int64) error
	SetMany(ctx context.Context, items MetricsStorageItems) error
}

type MetricsStorage interface {
	MetricsStorageGetter
	MetricsStorageSetter
	Ping(ctx context.Context) error
}

func New() MetricsStorage {
	return memStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

type memStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (m memStorage) GetGauge(ctx context.Context, name string) (*float64, error) {
	if val, ok := m.gauges[name]; ok {
		return &val, nil
	}

	return nil, nil
}

func (m memStorage) GetCounter(ctx context.Context, name string) (*int64, error) {
	if val, ok := m.counters[name]; ok {
		return &val, nil
	}

	return nil, nil
}

func (m memStorage) GetMany(ctx context.Context, names MetricsStorageKeys) (MetricsStorageItems, error) {
	rv := MetricsStorageItems{
		Gauges:   make(map[string]float64, len(names.Gauges)),
		Counters: make(map[string]int64, len(names.Counters)),
	}
	for _, v := range names.Gauges {
		gauge, _ := m.GetGauge(ctx, v)
		if gauge == nil {
			continue
		}
		rv.Gauges[v] = *gauge
	}
	for _, v := range names.Counters {
		counter, _ := m.GetCounter(ctx, v)
		if counter == nil {
			continue
		}
		rv.Counters[v] = *counter
	}

	return rv, nil
}

func (m memStorage) GetAll(ctx context.Context) (MetricsStorageItems, error) {
	rv := MetricsStorageItems{
		Gauges:   make(map[string]float64, len(m.gauges)),
		Counters: make(map[string]int64, len(m.counters)),
	}
	for k, v := range m.gauges {
		rv.Gauges[k] = v
	}
	for k, v := range m.counters {
		rv.Counters[k] = v
	}

	return rv, nil
}

func (m memStorage) SetMany(ctx context.Context, items MetricsStorageItems) error {
	maps.Copy(m.gauges, items.Gauges)
	for k, v := range items.Counters {
		m.SetCounter(ctx, k, v)
	}

	return nil
}

func (m memStorage) SetGauge(ctx context.Context, name string, value float64) error {
	m.gauges[name] = value
	return nil
}

func (m memStorage) SetCounter(ctx context.Context, name string, value int64) error {
	m.counters[name] += value
	return nil
}

func (m memStorage) Ping(ctx context.Context) error {
	return nil
}
