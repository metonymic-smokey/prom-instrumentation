package main

import (
	"encoding/json"
	"errors"
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

func getTempData() (Response, error) {
	url := fmt.Sprintf("https://api.tomorrow.io/v4/timelines?location=%f,%f&fields=temperature&timesteps=%s&units=%s", 73.98529171943665, 40.75872069597532, "1h", "metric")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Response{}, errors.New("error in GET request")
	}
	req.Header.Add("apikey", "APIKEY")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Response{}, errors.New("error reading response")
	}

	var dat Response
	if err := json.Unmarshal(body, &dat); err != nil {
		return Response{}, errors.New("error unmarshalling JSON")
	}

	return dat, nil
}
