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

type aircraft struct {
	Icao         string
	Reg          string
	Icaotype     string
	Year         string
	Manufacturer string
	Model        string
	Ownop        string
	Faa_pia      bool
	Faa_ladd     bool
	Short_type   string
	Mil          bool
}

var aircrafts map[string]aircraft

func Setup(conf config.Configuration) {
	aircrafts = make(map[string]aircraft)

	if err := downloadAircraftData(conf.Acdburl, conf.Acdbfilename); err != nil {
		log.Fatalf("unable to download aircraft data %v", err)
	}
	if err := readAircraftData(conf.Acdbfilename); err != nil {
		panic(err)
	}
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
		a := aircraft{}
		err = json.Unmarshal([]byte(v), &a)
		if err != nil {
			fmt.Printf("read line: %v\n", v)
			return err
		}
		aircrafts[a.Reg] = a
	}

	fmt.Printf("aircraft database size: %v\n", len(aircrafts))
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

	fmt.Printf("download aircraft database: %v\n", url)
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

	fmt.Printf("is file %v older than 1 day ", filename)
	info, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("does not exists -> false\n")
		return true
	}

	today := time.Now()
	yesterday := today.Add(-24 * time.Hour)

	fmt.Printf("file exists -> %v\n", yesterday.After(info.ModTime()))
	return yesterday.After(info.ModTime())

}
