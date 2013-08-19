package main

import (
	"fmt"
	"github.com/yunge/sphinx"
	"strconv"
	"time"
)

type SphinxStatusData map[string]string
type MetricsDataSource struct {
	SphinxHost        string
	Port              int
	ConnectionTimeout int

	PreviousData   SphinxStatusData
	LastData       SphinxStatusData
	LastUpdateTime time.Time
}

func NewMetricsDataSource(sphinxHost string, port int, connectionTimeout int) *MetricsDataSource {
	ds := &MetricsDataSource{
		SphinxHost:        sphinxHost,
		Port:              port,
		ConnectionTimeout: connectionTimeout,
	}
	return ds
}

func (ds *MetricsDataSource) CheckAndGetData(key string) (float64, error) {
	if err := ds.CheckAndUpdateData(); err != nil {
		return 0, err
	}

	prev, last, err := ds.GetOriginalData(key)

	if err != nil {
		return 0, err
	}
	return last - prev, nil
}
func (ds *MetricsDataSource) CheckAndGetLastData(key string) (float64, error) {
	if err := ds.CheckAndUpdateData(); err != nil {
		return 0, err
	}

	_, last, err := ds.GetOriginalData(key)

	if err != nil {
		return 0, err
	}
	return last, nil
}

func (ds *MetricsDataSource) GetOriginalData(key string) (float64, float64, error) {
	previousValue, ok := ds.PreviousData[key]
	if !ok {
		return 0, 0, fmt.Errorf("Can not get data from source \n")
	}
	currentValue, ok := ds.LastData[key]
	if !ok {
		return 0, 0, fmt.Errorf("Can not get data from source \n")
	}

	//some metric calculation can be turned off by sphinx settings
	if previousValue == "OFF" || currentValue == "OFF" {
		return 0, 0, nil
	}

	previousValueConverted, err := strconv.ParseFloat(previousValue, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Can not convert previous value of %s to int \n", key)
	}
	currentValueConverted, err := strconv.ParseFloat(currentValue, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Can not convert current value of %s to int \n", key)
	}

	return previousValueConverted, currentValueConverted, nil
}

func (ds *MetricsDataSource) CheckAndUpdateData() error {
	startTime := time.Now()
	if startTime.Sub(ds.LastUpdateTime) > time.Second*MIN_PAUSE_TIME {
		newData, err := ds.QueryData()
		if err != nil {
			return err
		}

		if ds.PreviousData == nil {
			ds.PreviousData = newData
		} else {
			ds.PreviousData = ds.LastData
		}
		ds.LastData = newData
		ds.LastUpdateTime = startTime
	}

	// check uptime
	//If uptime is less then in previous run - then server were restarted
	if prev, last, err := ds.GetOriginalData("uptime"); err != nil {
		return err
	} else {
		if last < prev {
			ds.PreviousData = ds.LastData
		}
	}
	return nil
}

func (ds *MetricsDataSource) QueryData() (SphinxStatusData, error) {
	client := sphinx.NewClient().SetServer(ds.SphinxHost, ds.Port)

	if ds.ConnectionTimeout != 0 {
		client.SetConnectTimeout(ds.ConnectionTimeout)
	}
	if err := client.Error(); err != nil {
		return nil, err
	}

	defer client.Close()
	status, err := client.Status()
	if err != nil {
		return nil, err
	}

	data := make(SphinxStatusData, len(status))
	for _, row := range status {
		data[row[0]] = row[1]
	}

	return data, nil
}
