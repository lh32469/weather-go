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
	"sort"
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
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  types.ThemeMacarons,
			Width:  "1200px",
			Height: "800px"}),
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

func httpserver2(w http.ResponseWriter, _ *http.Request) {

	timesToTemps := getTimesToTemps(getObservation("e5093"))

	keys := make([]time.Time, 0, len(timesToTemps))

	for fullTime, _ := range timesToTemps {
		keys = append(keys, fullTime)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	// create a new line instance
	chart := charts.NewLine()

	chart.Y = "Temperature"
	// set some global options like Title/Legend/ToolTip or anything else
	chart.SetGlobalOptions(
		charts.WithXAxisOpts(
			opts.XAxis{
				AxisLabel: getXAxisLabel2(),
			},
		),
		charts.WithInitializationOpts(
			opts.Initialization{
				Theme:     types.ThemeWesteros,
				PageTitle: "WeatherGraphApp",
				Width:     "auto",
				Height:    "800px"},
		),
		charts.WithTitleOpts(
			opts.Title{
				Title:    "WeatherGraphApp",
				Subtitle: "Current and 24 hr previous temperature graphs",
			},
		),
	)

	// Put data into instance
	chart.SetXAxis(keys).
		AddSeries("Current", generateTemperatureLineItems(timesToTemps)).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	err := chart.Render(w)
	if err != nil {
		panic(err)
	}
}

func getXAxisLabel() *opts.AxisLabel {
	return &opts.AxisLabel{
		ShowMinLabel: true,
		Formatter:    "{value} kg",
	}
}

func getXAxisLabel2() *opts.AxisLabel {
	return &opts.AxisLabel{
		ShowMinLabel: true,
		Formatter:    opts.FuncOpts(dateFormatter),
	}
}

var dateFormatter = `
function (value, index) {
    var parsed = Date.parse(value);
	console.log(index + value + parsed);
    var date = new Date(parsed);
	console.log('parsed date ' + date);
	var text = (date.getMonth() + 1) + '/' + (date.getDate()) 
		+ ' ' + (date.getHours() + 1) + ':' + (date.getMinutes() + 1);
	return text;
}
`

func main3() {
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

func main() {
	port := "10002"
	http.HandleFunc("/", httpserver2)
	log.Printf("Running at port: " + port)

	http.ListenAndServe(":"+port, nil)
}

func main1() {

	timesToTemps := getTimesToTemps(getObservation("e5093"))

	// Sort the Time keys
	times := make([]time.Time, 0)
	for k, _ := range timesToTemps {
		times = append(times, k)
	}

	sort.Slice(times, func(i, j int) bool {
		return times[i].Before(times[j])
	})

	for _, k := range times {
		fmt.Print(k)
		fmt.Print(" = ")
		fmt.Println(timesToTemps[k])
	}

}

func generateTemperatureLineItems(data map[time.Time]float32) []opts.LineData {

	items := make([]opts.LineData, 0)

	// Sort the Time keys
	times := make([]time.Time, 0)
	for k, _ := range data {
		times = append(times, k)
	}

	sort.Slice(times, func(i, j int) bool {
		return times[i].Before(times[j])
	})

	for _, k := range times {
		fmt.Print(k)
		fmt.Print(" = ")
		fmt.Println(data[k])
		items = append(items, opts.LineData{Value: data[k]})
	}

	return items
}

func getObservation(station string) Observation {

	var url = strings.ReplaceAll(fiveMinuteTemplate, "STATION", station)
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

	var timeSeries TimeSeries
	err = json.Unmarshal(body, &timeSeries)

	return timeSeries.STATION[0].Observations
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
		result[mTime] = airTemp[i]*9/5 + 32
	}

	return result
}
