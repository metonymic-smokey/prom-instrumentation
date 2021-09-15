package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Temperature struct {
	Temp float64 `json:"temperature"`
}

type Final struct {
	StartTime string      `json:"startTime"`
	Values    Temperature `json:"values"`
}

type Interval struct {
	Timestep  string  `json:"timestep"`
	StartTime string  `json:"startTime"`
	EndTime   string  `json:"endTime"`
	TempVal   []Final `json:"intervals"`
}

type Timelines struct {
	Timestep []Interval `json:"timelines"`
}

type Response struct {
	Data Timelines `json:"data"`
}

func getTempData() Response {
	url := fmt.Sprintf("https://api.tomorrow.io/v4/timelines?location=%f,%f&fields=temperature&timesteps=%s&units=%s", 73.98529171943665, 40.75872069597532, "1h", "metric")

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("apikey", "APIKEY")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var dat Response
	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}

	return dat
}
