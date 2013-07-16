package main

import (
	"github.com/yvasiyarov/newrelic_platform_go"
)

type MetricsDataSource struct {
}

type Metrica struct {
    Name string
    Units string
    DataKey string
    Datasource *MetricsDataSource
}

func (metrica *Metrica) GetName() string {
	return metrica.Name
}
func (metrica *WaveMetrica) GetUnits() string {
	return metrica.Units
}
func (metrica *WaveMetrica) GetValue() (float64, error) {
	metrica.sawtoothCounter++
	if metrica.sawtoothCounter > metrica.sawtoothMax {
		metrica.sawtoothCounter = 0
	}
	return float64(metrica.sawtoothCounter), nil
}


func main() {
	plugin := newrelic_platform_go.NewNewrelicPlugin("0.0.1", "7bceac019c7dcafae1ef95be3e3a3ff8866de246", 60)
	component := newrelic_platform_go.NewPluginComponent("Wave component", "com.exmaple.plugin.gowave")
	plugin.AddComponent(component)

	m := &WaveMetrica{
		sawtoothMax:     10,
		sawtoothCounter: 5,
	}

	component.AddMetrica(m)
	plugin.Verbose = true
	plugin.Run()
}
