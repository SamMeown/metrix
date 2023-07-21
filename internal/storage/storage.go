package storage

import "golang.org/x/exp/maps"

const (
	MetricsTypeGauge   = "gauge"
	MetricsTypeCounter = "counter"
)

type MetricsStorageItems struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

type MetricsStorageGetter interface {
	GetGauge(name string) (*float64, error)
	GetCounter(name string) (*int64, error)
	GetAll() (MetricsStorageItems, error)
}

type MetricsStorage interface {
	MetricsStorageGetter
	SetGauge(name string, value float64) error
	SetCounter(name string, value int64) error
	SetMany(items MetricsStorageItems) error
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

func (m memStorage) GetGauge(name string) (*float64, error) {
	if val, ok := m.gauges[name]; ok {
		return &val, nil
	}

	return nil, nil
}

func (m memStorage) GetCounter(name string) (*int64, error) {
	if val, ok := m.counters[name]; ok {
		return &val, nil
	}

	return nil, nil
}

func (m memStorage) GetAll() (MetricsStorageItems, error) {
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

func (m memStorage) SetMany(items MetricsStorageItems) error {
	maps.Copy(m.gauges, items.Gauges)
	maps.Copy(m.counters, items.Counters)

	return nil
}

func (m memStorage) SetGauge(name string, value float64) error {
	m.gauges[name] = value
	return nil
}

func (m memStorage) SetCounter(name string, value int64) error {
	m.counters[name] += value
	return nil
}
