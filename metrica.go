package main

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
	return metrica.DataSource.CheckAndGetLastData(metrica.DataKey)
}

type IncrementalMetrica struct {
	Metrica
}

func (metrica *IncrementalMetrica) GetValue() (float64, error) {
	return metrica.DataSource.CheckAndGetData(metrica.DataKey)
}

var plainMetrics = []*Metrica{
	&Metrica{
		DataKey: "avg_query_wall",
		Name:    "avg/Avg Query Wall Time",
		Units:   "milisecond",
	},
}
var incrementalMetrics = []*Metrica{
	&Metrica{
		DataKey: "queries",
		Name:    "general/Queries",
		Units:   "Queries/second",
	},
	&Metrica{
		DataKey: "connections",
		Name:    "general/Connections",
		Units:   "connections/second",
	},
	&Metrica{
		DataKey: "maxed_out",
		Name:    "error/Maxed out connections",
		Units:   "connections/second",
	},
	&Metrica{
		DataKey: "command_search",
		Name:    "commands/Command search",
		Units:   "command/second",
	},
	&Metrica{
		DataKey: "command_excerpt",
		Name:    "commands/Command excerpt",
		Units:   "command/second",
	},
	&Metrica{
		DataKey: "command_update",
		Name:    "commands/Command update",
		Units:   "command/second",
	},
	&Metrica{
		DataKey: "command_keywords",
		Name:    "commands/Command keywords",
		Units:   "command/second",
	},
	&Metrica{
		DataKey: "command_persist",
		Name:    "commands/Command persist",
		Units:   "command/second",
	},
	&Metrica{
		DataKey: "command_flushattrs",
		Name:    "commands/Command flushattrs",
		Units:   "command/second",
	},
}
