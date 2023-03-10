package tracker

import (
	"acTrackerBot/config"
	"acTrackerBot/tracker/acdb"
	"acTrackerBot/tracker/aeroDataBox"
	"acTrackerBot/tracker/types"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

// send aircraft update information back through this channel
var updateChannel chan types.AircraftInformation

// all aircraft trackers have a channel in this map to stop the goroutine
var stopChannels map[string]chan int

func StartUp() <-chan types.AircraftInformation {

	// this channels must be initialized from the sender
	if AddRegistrationChannel == nil || RemoveRegistrationChannel == nil {
		log.Fatalf("add and remove channel is not initialized")
	}

	stopChannels = make(map[string]chan int)

	readRegistrationDatabase()
	if err := readRegistrationList(); err != nil {
		log.Printf("unable to import aircraft list %v\n", err)
	}

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
			if reg == "" {
				log.Println("error: recieve empty add message")
				continue
			}
			if err := addNewReg(reg); err != nil {
				log.Printf("%v\n", err.Error())
			}
		case reg := <-RemoveRegistrationChannel:
			if reg == "" {
				log.Println("error: recieve empty remove message")
				continue
			}
			if err := removeReg(reg); err != nil {
				log.Printf("%v\n", err.Error())
			}
		case <-stop:
			if err := saveRegistrationList(); err != nil {
				log.Printf("can not save registraion list %v\n", err.Error())
			}
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

	if reg == "" {
		return fmt.Errorf("no registration set to remove update process")
	}

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
	log.Printf("import last registration list to %v\n", config.Conf.Callsignllistfilename)
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
		if err := addNewReg(line); err != nil {
			log.Printf("can not add registration %v %v", line, err)
		}
	}

	return nil
}

func saveRegistrationList() error {
	log.Printf("save active registration list to %v\n", config.Conf.Callsignllistfilename)
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
			//log.Printf("update aircraft '%v'\n", reg)
			data := requestAdsbExchangeData(reg)
			processData(reg, data)

			if newStatus(reg) {
				sendUpdate(reg)
			}

		case <-sc:
			log.Printf("stop aircraft tracker for registration '%v'\n", reg)
			return
		}
	}

}

func sendUpdate(reg string) {
	info := getCurrentAircraftInfo(reg)
	info.Origin = "N/A"
	info.Destination = "N/A"
	info.FlightStatus = "N/A"

	flightStatus, err := aeroDataBox.GetFlightInfo(reg)
	if err != nil {
		log.Printf("error: %v\n", err)
		updateChannel <- info
		return
	}

	addFlightStatus(flightStatus, &info)

	log.Printf("<- send flight status: %v\n", info)
	updateChannel <- info
}

func addFlightStatus(info *aeroDataBox.FlightStatus, aicraftInfo *types.AircraftInformation) {

	aicraftInfo.Origin = "N/A"
	aicraftInfo.Destination = "N/A"
	aicraftInfo.FlightStatus = "N/A"

	if info.FlightDeparture.DepartureAirport.Icao != "" {
		aicraftInfo.Origin = info.FlightDeparture.DepartureAirport.Icao
	}

	if info.FlightArrival.ArrivalAirport.Icao != "" {
		aicraftInfo.Destination = info.FlightArrival.ArrivalAirport.Icao
	}

	if info.FlightStatus != "" {
		aicraftInfo.FlightStatus = info.FlightStatus
	}

}

/*
	func addFlightawareData(callsign string, aircraftInfo *types.AircraftInformation) {
		url := fmt.Sprintf("https://aeroapi.flightaware.com/aeroapi/flights/search?query=-idents+%v", callsign)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("x-apikey", config.Conf.Flightawareapikey)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("unable to request flightaware data: error %v\n", err)
			return
		}

		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		flightawareData := types.FlightawareFlights{}
		err = json.Unmarshal(body, &flightawareData)
		if err != nil {
			log.Printf("can not unmarshal %v\n%v\n", string(body), err)
			return
		}

		log.Printf("flightaware data len: %v\n", len(flightawareData.Flights))
		if len(flightawareData.Flights) < 1 {
			log.Printf("no flightaware data available\n")
			return
		}

		log.Printf("%v : %v -> %v\n", callsign, flightawareData.Flights[0].Origin.Code, flightawareData.Flights[0].Destination.Code)

		if flightawareData.Flights[0].Origin.Code != "" {
			aircraftInfo.Origin = flightawareData.Flights[0].Origin.Code
		}

		if flightawareData.Flights[0].Destination.Code != "" {
			aircraftInfo.Destination = flightawareData.Flights[0].Destination.Code
		}

}
*/
func requestAdsbExchangeData(reg string) (data types.AdsbExchData) {
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
	body, _ := io.ReadAll(res.Body)

	data = types.AdsbExchData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("can not unmarshal %v\n%v\n", string(body), err)
		return
	}

	return data
}
