package acdb

import (
	"acTrackerBot/config"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Aircraft struct {
	Icao         string `json:"icao"`
	Reg          string `json:"reg"`
	Icaotype     string `json:"icaotype"`
	Year         string `json:"year"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Ownop        string `json:"ownop"`
	Faa_pia      bool   `json:"faa_pia"`
	Faa_ladd     bool   `json:"faa_ladd"`
	Short_type   string `json:"short_type"`
	Mil          bool   `json:"mil"`
}

var aircrafts map[string]Aircraft

func Setup(conf config.Configuration) {
	aircrafts = make(map[string]Aircraft)

	if err := downloadAircraftData(conf.Acdburl, conf.Acdbfilename); err != nil {
		log.Fatalf("unable to download aircraft data %v", err)
	}
	if err := readAircraftData(conf.Acdbfilename); err != nil {
		panic(err)
	}
}

func GetAircraftData(reg string) Aircraft {
	return aircrafts[reg]
}

func IsValidReg(reg string) bool {
	_, ok := aircrafts[reg]
	return ok
}

func readAircraftData(filename string) error {

	fi, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return err
	}
	defer fz.Close()

	s, err := ioutil.ReadAll(fz)
	if err != nil {
		return err
	}

	a := strings.Split(string(s), "}")

	for _, v := range a {
		if len(v) == 1 {
			continue
		}
		v = v + "}"
		a := Aircraft{Year: "unknown", Manufacturer: "unknown", Model: "unknown", Ownop: "unknown"}
		err = json.Unmarshal([]byte(v), &a)
		if err != nil {
			log.Printf("read line: %v\n", v)
			return err
		}

		if a.Year == "" {
			a.Year = "unknown"
		}
		if a.Manufacturer == "" {
			a.Manufacturer = "unknown"
		}
		if a.Model == "" {
			a.Model = "unknown"
		}
		if a.Ownop == "" {
			a.Ownop = "unknown"
		}

		aircrafts[a.Reg] = a
	}

	log.Printf("aircraft database size: %v\n", len(aircrafts))
	return nil
}

func downloadAircraftData(url string, filename string) error {
	if !isFileOlderThan1Day(filename) {
		return nil
	}

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	log.Printf("download aircraft database: %v\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func isFileOlderThan1Day(filename string) bool {

	info, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		return true
	}

	today := time.Now()
	yesterday := today.Add(-24 * time.Hour)
	return yesterday.After(info.ModTime())

}
