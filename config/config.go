package config

import (
	"encoding/json"
	"io"
	"os"
)

type Configuration struct {
	Discordbottoken        string `json:"discordbottoken"`
	Discrodwebhookid       string `json:"discrodwebhookid"`
	Discrodwebhooktoken    string `json:"discrodwebhooktoken"`
	Adsbrapidapikey        string `json:"adsbrapidapikey"`
	Adsbrapidapihost       string `json:"adsbrapidapihost"`
	Flightawareapikey      string `json:"flightawareapikey"`
	Acdburl                string `json:"acdburl"`
	Acdbfilename           string `json:"acdbfilename"`
	UpdateIntervall        int    `json:"updateIntervall"`
	Callsignllistfilename  string `json:"callsignllistfilename"`
	Aerodataboxrapidapikey string `json:"aerodataboxrapidapikey"`
	Aerodataboxrapihost    string `json:"aerodataboxrapihost"`
}

var Conf = Configuration{}

func ReadConfig() {

	cfile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer cfile.Close()

	conf, err := io.ReadAll(cfile)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(conf, &Conf); err != nil {
		panic(err)
	}
}
