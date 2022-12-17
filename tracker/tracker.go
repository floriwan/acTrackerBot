package tracker

import (
	"acTrackerBot/config"
	"acTrackerBot/tracker/acdb"
	"acTrackerBot/tracker/types"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

// channel to send new registration to tracker
var AddRegistrationChannel chan string

// channel to remove registration from tracker
var RemoveRegistrationChannel chan string

var updateChannel chan types.AircraftInformation
var stopChannels map[string]chan int

func StartUp() <-chan types.AircraftInformation {

	if AddRegistrationChannel == nil || RemoveRegistrationChannel == nil {
		log.Fatalf("add and remove channel is not initialized")
	}

	stopChannels = make(map[string]chan int)

	readRegistrationDatabase()
	readRegistrationList()

	go runTracker()

	updateChannel = make(chan types.AircraftInformation)
	return updateChannel
}

func runTracker() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	log.Printf("tracker is running\n")

	for {
		select {
		case reg := <-AddRegistrationChannel:
			addNewReg(reg)
		case reg := <-RemoveRegistrationChannel:
			removeReg(reg)
		case <-stop:
			log.Printf("stopping tracker\n")
			saveRegistrationList()
			stopAllAircraftTracker()
		}
	}
}

func stopAllAircraftTracker() {
	log.Println("stopping all aircraft tracker")
	for _, v := range stopChannels {
		v <- 0
	}
}

func readRegistrationDatabase() {
	acdb.Setup(config.Conf)
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

func removeReg(reg string) error {
	c, ok := stopChannels[reg]
	if !ok {
		return fmt.Errorf("no update process for reg '%v' found", reg)
	}

	// send something to stop the go routine
	log.Printf("stopping update process for '%v'\n", reg)
	c <- 0
	delete(stopChannels, reg)
	close(c)
	return nil
}

func addNewReg(reg string) error {
	if !acdb.IsValidReg(reg) {
		return fmt.Errorf("registration '%v' is not a valid registration", reg)
	}
	go startAircraftTracker(config.Conf.UpdateIntervall, reg)
	return nil
}

func readRegistrationList() error {
	file, err := os.Open(config.Conf.Callsignllistfilename)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// after all registrations are read, add them to the tracker
	for _, line := range lines {
		addNewReg(line)
	}

	return nil
}

func saveRegistrationList() error {
	file, err := os.Create(config.Conf.Callsignllistfilename)
	if err != nil {
		return err
	}
	defer file.Close()

	keys := make([]string, 0, len(stopChannels))
	for k := range stopChannels {
		keys = append(keys, k)
	}

	w := bufio.NewWriter(file)
	for _, line := range keys {
		fmt.Fprintln(w, line)
	}
	return w.Flush()

}

func startAircraftTracker(interval int, reg string) {
	log.Printf("starting new aircraft tracker for '%v'\n", reg)

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
				//log.Printf("'%v' new aircraft state: %v\n", reg, getCurrentAircraftInfo(reg))
				info := getCurrentAircraftInfo(reg)
				aircraftData := acdb.GetAircraftData(reg)
				info.IcaoType = aircraftData.Icaotype
				updateChannel <- info
			}
			//else {
			//	log.Printf("'%v' no status change: %v\n", reg, getCurrentAircraftInfo(reg))
			//}

		case <-sc:
			log.Printf("stop aircraft tracker for registration '%v'\n", reg)
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
