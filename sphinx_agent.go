package main

import (
	"github.com/yunge/sphinx"
        "fmt"
"time"
"strconv"
	"github.com/yvasiyarov/newrelic_platform_go"
)

const (
    MIN_PAUSE_TIME = 30
)

type SphinxStatusData map[string]string
type MetricsDataSource struct {
	SphinxHost        string
	Port              int
	ConnectionTimeout int

        PreviousData SphinxStatusData
        LastData     SphinxStatusData
        LastUpdateTime time.Time
}

func NewMetricsDataSource(sphinxHost string, port int, connectionTimeout int) *MetricsDataSource {
	if port == 0 {
		port = 9312
	}
	ds := &MetricsDataSource{
		SphinxHost:        sphinxHost,
		Port:              port,
		ConnectionTimeout: connectionTimeout,
	}
	return ds
}

func (ds *MetricsDataSource) GetData(key string) (float64, error) {
    if err := ds.CheckAndUpdateData(); err != nil {
        return 0, err
    }

    prev, last, err := ds.GetOriginalData(key)

    if err != nil {
        return 0, err
    }
    return last - prev, nil
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
        if startTime.Sub(ds.LastUpdateTime) > time.Second * MIN_PAUSE_TIME {
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

type Metrica struct {
	Name       string
	Units      string
	DataKey    string
	DataSource *MetricsDataSource
}

func (metrica *Metrica) GetName() string {
	return metrica.Name
}
func (metrica *Metrica) GetUnits() string {
	return metrica.Units
}
func (metrica *Metrica) GetValue() (float64, error) {
    return metrica.DataSource.GetData(metrica.DataKey)
}

func AddMetrcas(component newrelic_platform_go.IComponent, dataSource *MetricsDataSource) {
    metricas := []*Metrica{
        &Metrica{
            DataKey: "queries",
            Name: "Queries",
            Units: "Queries/second",
        },
        &Metrica{
            DataKey: "connections",
            Name: "Connections",
            Units: "connections/second",
        },
        &Metrica{
            DataKey: "maxed_out",
            Name: "Maxed out",
            Units: "connections/second",
        },
        &Metrica{
            DataKey: "command_search",
            Name: "Command search",
            Units: "command/second",
        },
        &Metrica{
            DataKey: "command_excerpt",
            Name: "Command excerpt",
            Units: "command/second",
        },
        &Metrica{
            DataKey: "command_update",
            Name: "Command update",
            Units: "command/second",
        },
        &Metrica{
            DataKey: "command_keywords",
            Name: "Command keywords",
            Units: "command/second",
        },
        &Metrica{
            DataKey: "command_persist",
            Name: "Command persist",
            Units: "command/second",
        },
        &Metrica{
            DataKey: "command_flushattrs",
            Name: "Command flushattrs",
            Units: "command/second",
        },
    }
    for _, m := range metricas {
        m.DataSource = dataSource
        component.AddMetrica(m)
    }
}

func main() {
        plugin := newrelic_platform_go.NewNewrelicPlugin("0.0.1", "7bceac019c7dcafae1ef95be3e3a3ff8866de246", 60)
        component := newrelic_platform_go.NewPluginComponent("Sphinx component", "com.github.yvasiyarov.Sphinx")
        plugin.AddComponent(component)

	ds := NewMetricsDataSource("web-d5.butik.ru", 0, 0)
        AddMetrcas(component, ds)

        plugin.Verbose = true
  
        plugin.Run()
        //plugin.Harvest()
}
