package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// generate random data for line chart
func generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(300)})
	}
	return items
}

func httpserver(w http.ResponseWriter, _ *http.Request) {
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Line example in Westeros theme",
			Subtitle: "Line chart rendered by the http server this time",
		}))

	// Put data into instance
	line.SetXAxis([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
		AddSeries("Category A", generateLineItems()).
		AddSeries("Category B", generateLineItems()).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.Render(w)
}

func main2() {
	port := "10002"
	http.HandleFunc("/", httpserver)
	log.Printf("Running at port: " + port)

	res, err := http.Get("")

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	http.ListenAndServe(":"+port, nil)
}

type TimeSeries struct {
	QC_SUMMARY map[string]string
	STATION    []Station
	//STATION []map[string]string
	SUMMARY map[string]string
	UNITS   map[string]string
}

type Observation struct {
	AirTempSet1 []float32 `json:"air_temp_set_1"`
	DateTime    []string  `json:"date_time"`
}

type Station struct {
	County       string
	Observations Observation `json:"OBSERVATIONS"`
	//Observations []Observation `json:"OBSERVATIONS"`
	//Observations map[string][]string
}

var fiveMinuteTemplate = "https://api.mesowest.net/v2/stations/timeseries?stid=STATION&recent=4320&obtimezone=local&complete=1&hfmetars=1&token=d8c6aee36a994f90857925cea26934be"

func main() {

	var url = strings.ReplaceAll(fiveMinuteTemplate, "STATION", "kpdx")
	fmt.Println(url)

	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	//fmt.Println(body)

	//var foo map[string]string
	var foo = map[string]map[string]string{}
	err = json.Unmarshal(body, &foo)

	fmt.Println("------ Foo ------------")
	fmt.Println(foo)

	for k := range foo {
		fmt.Println("Key: " + k)
		fmt.Println(foo[k])
	}

	fmt.Println("------Foo UNITS ------------")
	fmt.Println(foo["STATION"])

	var timeSeries TimeSeries
	err = json.Unmarshal(body, &timeSeries)

	//fmt.Println("-------- TimeSeries ----------")
	//fmt.Println(timeSeries)

	fmt.Println("------ UNITS ------------")
	fmt.Println(timeSeries.UNITS)

	//fmt.Println("------ Stations ------------")
	//fmt.Println(timeSeries.STATION)

	//for k := range timeSeries.STATION {
	//	fmt.Println(k)
	//	fmt.Println(timeSeries.STATION[k])
	//}

	//fmt.Println("------ Station ------------")
	//fmt.Println(timeSeries.STATION[0])

	fmt.Println("------ Observations ------------")
	fmt.Println(timeSeries.STATION[0].Observations)

	fmt.Println("------ Observations AirTempSet1 ------------")
	fmt.Println(timeSeries.STATION[0].Observations.AirTempSet1)

	//	for k := range timeSeries.STATION[0].Observations {
	//		fmt.Println(k)
	//		//fmt.Println(timeSeries.STATION[0].Observations[k])
	//	}
	//
	//	fmt.Println("------ AirTempSet1 ------------")
	//	fmt.Println(timeSeries.STATION[0].Observations["metar_set_1"])

	str := "2022-07-24T12:15:00-0700"
	//str = strings.ReplaceAll(str, "-0700", "0700")
	fmt.Println("String: " + str)
	layout := time.RFC3339
	layout = "2006-01-02T15:04:05-0700"
	fmt.Println("Layout: " + layout)
	t, err := time.Parse(layout, str)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)

	timesToTemps := getTimesToTemps(timeSeries.STATION[0].Observations)

	for k := range timesToTemps {
		fmt.Print(k)
		fmt.Print(" = ")
		fmt.Println(timesToTemps[k])
	}

}

/**
 * Process the Observations Map provided and return a Map of temperature
 * reading time to temperature reading in degrees fahrenheit.
 */
func getTimesToTemps(observation Observation) map[time.Time]float32 {

	result := make(map[time.Time]float32)

	airTemp := observation.AirTempSet1
	times := observation.DateTime

	for i, t := range times {
		mTime, err := time.Parse("2006-01-02T15:04:05-0700", t)
		if err != nil {
			fmt.Println(err)
		}
		result[mTime] = airTemp[i]
	}

	return result
}
