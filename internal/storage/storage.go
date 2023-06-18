package storage

const (
	MetricsTypeGauge   = "gauge"
	MetricsTypeCounter = "counter"
)

type MetricsStorage interface {
	SetGauge(name string, value float64) error
	SetCounter(name string, value int64) error
	Value(name string) (any, error)
	Values() (map[string]any, error)
}

func New() MetricsStorage {
	return memStorage{make(map[string]any)}
}

type memStorage struct {
	values map[string]any
}

func (m memStorage) Value(name string) (any, error) {
	return m.values[name], nil
}

func (m memStorage) Values() (map[string]any, error) {
	rv := make(map[string]any, len(m.values))
	for k, v := range m.values {
		rv[k] = v
	}

	return rv, nil
}

func (m memStorage) SetGauge(name string, value float64) error {
	m.values[name] = value
	return nil
}

func (m memStorage) SetCounter(name string, value int64) error {
	if _, ok := m.values[name]; !ok {
		m.values[name] = int64(0)
	}
	m.values[name] = m.values[name].(int64) + value

	return nil
}
