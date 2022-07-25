package main

var fiveMinuteTemplate = "https://api.mesowest.net/v2/stations/timeseries?stid=STATION&recent=4320&obtimezone=local&complete=1&hfmetars=1&token=d8c6aee36a994f90857925cea26934be"

type TimeSeries struct {
	//QC_SUMMARY map[string]string
	STATION []Station
	//SUMMARY    map[string]string
	//UNITS      map[string]string
}

type Station struct {
	County       string
	Observations Observation `json:"OBSERVATIONS"`
}

type Observation struct {
	AirTempSet1 []float32 `json:"air_temp_set_1"`
	DateTime    []string  `json:"date_time"`
}
