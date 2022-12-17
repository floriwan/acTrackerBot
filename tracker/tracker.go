package tracker

import (
	"acTrackerBot/config"
	"acTrackerBot/tracker/acdb"
	"acTrackerBot/tracker/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var updateChannel chan types.AircraftInformation
var stopChannels map[string]chan int

func StartUp() <-chan types.AircraftInformation {
	stopChannels = make(map[string]chan int)
	go readRegistrationDatabase()
	updateChannel = make(chan types.AircraftInformation)
	return updateChannel
}

func readRegistrationDatabase() {
	acdb.Setup(config.Conf)
	//updateChannel <- "db ready"
}

func GetRegListSize() int {
	return len(stopChannels)
}

func GetRegList() string {
	keys := make([]string, 0, len(stopChannels))
	for k := range stopChannels {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func RemoveReg(reg string) error {
	c, ok := stopChannels[reg]
	if !ok {
		return fmt.Errorf("no update process for reg '%v' found", reg)
	}

	// send something to stop the go routine
	log.Printf("stopping update process for '%v'\n", reg)
	c <- 0
	delete(stopChannels, reg)
	return nil
}

func AddNewReg(reg string) error {
	if !acdb.IsValidReg(reg) {
		return fmt.Errorf("registration '%v' is not a valid registration", reg)
	}
	go startTracker(config.Conf.UpdateIntervall, reg)
	return nil
}

func startTracker(interval int, reg string) {
	log.Printf("starting new tracker for '%v'\n", reg)

	sc := make(chan int)
	stopChannels[reg] = sc

	ticker := time.NewTicker(time.Duration(interval * int(time.Minute)))

	for {
		select {
		case <-ticker.C:
			log.Printf("update reg '%v'\n", reg)
			data := requestData(reg)
			processData(reg, data)

			if newStatus(reg) {
				log.Printf("'%v' new aircraft state: %v\n", reg, getCurrentAircraftInfo(reg))
				info := getCurrentAircraftInfo(reg)
				aircraftData := acdb.GetAircraftData(reg)
				info.IcaoType = aircraftData.Icaotype
				updateChannel <- info
			} else {
				log.Printf("'%v' no status change: %v\n", reg, getCurrentAircraftInfo(reg))
			}

		case <-sc:
			log.Printf("stop update for registration '%v'\n", reg)
			return
		}
	}

}

func requestData(reg string) (data types.AdsbExchData) {
	url := fmt.Sprintf("https://adsbexchange-com1.p.rapidapi.com/v2/registration/%v/", reg)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Key", config.Conf.Adsbrapidapikey)
	req.Header.Add("X-RapidAPI-Host", config.Conf.Adsbrapidapihost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error %v\n", err)
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	data = types.AdsbExchData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("can not unmarshal %v\n%v\n", string(body), err)
		return
	}

	return data
}
