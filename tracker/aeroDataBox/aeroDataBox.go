package aeroDataBox

import (
	"acTrackerBot/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func GetFlightInfo(reg string) (status *FlightStatus, err error) {

	url := fmt.Sprintf("https://aerodatabox.p.rapidapi.com/flights/reg/%v?withAircraftImage=false&withLocation=false", reg)
	data, err := sendRequest(url)
	if err != nil {
		return nil, fmt.Errorf("can not send http request to aerodatabox %v", err)
	}

	flightStatusList := FlightStatusResult
	if err := json.Unmarshal(data, &flightStatusList); err != nil {
		return nil, fmt.Errorf("can not unmarshall aircraft data %v %v", string(data), err)
	}

	for _, v := range flightStatusList {
		log.Printf("flight status: %+v\n", v)
		departureTime := v.FlightDeparture.ScheduledTimeUtc
		arrivalTime := v.FlightArrival.ScheduledTimeUtc

		if departureTime == "" || arrivalTime == "" {
			log.Printf("departure %v or arrival %v time not set", departureTime, arrivalTime)
			continue
		}

		depTime, err := time.Parse("2006-01-02 15:04Z", departureTime)
		if err != nil {
			log.Printf("can not convert departure time string %v %v", departureTime, err)
			continue
		}

		arrTime, err := time.Parse("2006-01-02 15:04Z", arrivalTime)
		if err != nil {
			log.Printf("can not convert arrival time string %v %v", arrivalTime, err)
			continue
		}

		utcTime := time.Now().UTC()

		if depTime.Before(utcTime) && arrTime.After(utcTime) {
			//fmt.Printf("-> found flight : %+v\n", v)
			status = &v
			return status, nil
		}

	}
	return nil, fmt.Errorf("no flight status found for registration %v", reg)
}

func GetAircraftInfo(reg string) (data *Aircraft, err error) {

	if config.Conf.Aerodataboxrapidapikey == "" ||
		config.Conf.Aerodataboxrapihost == "" {
		return nil, fmt.Errorf("no api key %v or no api host %v set", config.Conf.Aerodataboxrapidapikey, config.Conf.Aerodataboxrapihost)
	}

	url := fmt.Sprintf("https://aerodatabox.p.rapidapi.com/aircrafts/reg/%v?withImage=true", reg)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", config.Conf.Aerodataboxrapidapikey)
	req.Header.Add("X-RapidAPI-Host", config.Conf.Aerodataboxrapihost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can not send http request to aerodatabox %v", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("can not unmarshall aircraft data %v %v", string(body), err)
	}

	return data, nil

}

func sendRequest(url string) (data []byte, err error) {
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", config.Conf.Aerodataboxrapidapikey)
	req.Header.Add("X-RapidAPI-Host", config.Conf.Aerodataboxrapihost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can not send http request to aerodatabox %v", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return body, nil

}
