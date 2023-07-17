package saver

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/SamMeown/metrix/internal/models"
	"github.com/SamMeown/metrix/internal/storage"
)

type MetricsStorageSaver struct {
	storage storage.MetricsStorage
	file    *os.File
	writer  *bufio.Writer
	reader  *bufio.Reader
}

func NewMetricsStorageSaver(storage storage.MetricsStorage, path string) (*MetricsStorageSaver, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &MetricsStorageSaver{
		storage: storage,
		file:    file,
		writer:  bufio.NewWriter(file),
		reader:  bufio.NewReader(file),
	}, nil
}

func (s *MetricsStorageSaver) Load() error {
	for {
		data, err := s.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		var metrics models.Metrics
		err = json.Unmarshal(data, &metrics)
		if err != nil {
			return err
		}

		if metrics.MType == storage.MetricsTypeGauge {
			s.storage.SetGauge(metrics.ID, *metrics.Value)
		} else {
			s.storage.SetCounter(metrics.ID, *metrics.Delta)
		}
	}

	return nil
}

func (s *MetricsStorageSaver) Save() error {
	snapshot, err := s.storage.GetAll()
	if err != nil {
		return err
	}

	metricsList := make([]models.Metrics, 0)
	for name, value := range snapshot.Gauges {
		value := value
		metrics := models.Metrics{
			ID:    name,
			MType: storage.MetricsTypeGauge,
			Value: &value,
		}
		metricsList = append(metricsList, metrics)
	}

	for name, value := range snapshot.Counters {
		value := value
		metrics := models.Metrics{
			ID:    name,
			MType: storage.MetricsTypeCounter,
			Delta: &value,
		}
		metricsList = append(metricsList, metrics)
	}

	s.file.Truncate(0)
	s.file.Seek(0, 0)

	for _, metrics := range metricsList {
		data, err := json.Marshal(&metrics)
		if err != nil {
			return err
		}

		if _, err := s.writer.Write(data); err != nil {
			return err
		}
		if err := s.writer.WriteByte('\n'); err != nil {
			return err
		}
	}

	return s.writer.Flush()
}
