package aeroDataBox

var FlightStatusResult []FlightStatus

type FlightStatus struct {
	FlightDeparture Departure     `json:"departure"`
	FlightArrival   Arrival       `json:"arrival"`
	FlightStatus    string        `json:"status"`
	FlightAircraft  AircraftShort `json:"aircraft"`
}

type AircraftShort struct {
	Reg   string `json:"reg"`
	ModeS string `json:"modeS"`
	Model string `json:"model"`
}

type Departure struct {
	DepartureAirport Airport `json:"airport"`
	ScheduledTimeUtc string  `json:"scheduledTimeUtc"`
}

type Arrival struct {
	ArrivalAirport   Airport `json:"airport"`
	ScheduledTimeUtc string  `json:"scheduledTimeUtc"`
}

type Airport struct {
	Icao string `json:"icao"`
	Name string `json:"name"`
}

type Aircraft struct {
	Id               int           `json:"id"`
	Reg              string        `json:"reg"`
	Active           bool          `json:"active"`
	Serial           string        `json:"serial"`
	HexIcao          string        `json:"hexIcao"`
	AirlineName      string        `json:"airlineName"`
	IataCodeShort    string        `json:"iataCodeShort"`
	IcaoCode         string        `json:"icaoCode"`
	Model            string        `json:"model"`
	ModelCode        string        `json:"modelCode"`
	NumSeats         int           `json:"numSeats"`
	RolloutDate      string        `json:"rolloutDate"`
	FirstFlightDate  string        `json:"firstFlightDate"`
	DeliveryData     string        `json:"deliveryDate"`
	RegistrationDate string        `json:"registrationDate"`
	TypeName         string        `json:"typeName"`
	NumEngines       int           `json:"numEngines"`
	EngineType       string        `json:"engineType"`
	IsFreighter      bool          `json:"isFreighter"`
	ProductionLine   string        `json:"productionLine"`
	AgeYears         float32       `json:"ageYears"`
	Verified         bool          `json:"verified"`
	NumRegistration  int           `json:"numRegistrations"`
	Image            AircraftImage `json:"image"`
}

type AircraftImage struct {
	Url         string `json:"url"`
	WebUrl      string `json:"webUrl"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	License     string `json:"license"`
}

/*
{
  "id": 23431,
  "reg": "D-AIUD",
  "active": true,
  "serial": "6033",
  "hexIcao": "3C66A4",
  "airlineName": "Lufthansa",
  "iataCodeShort": "32A",
  "icaoCode": "A320",
  "model": "A320",
  "modelCode": "320-214",
  "numSeats": 168,
  "rolloutDate": "2014-03-13",
  "firstFlightDate": "2014-03-13",
  "deliveryDate": "2014-03-24",
  "registrationDate": "2014-03-24",
  "typeName": "Airbus A320 (Sharklets)",
  "numEngines": 2,
  "engineType": "Jet",
  "isFreighter": false,
  "productionLine": "Airbus A320",
  "ageYears": 8.8,
  "verified": true,
  "image": {
    "url": "https://farm6.staticflickr.com/5543/9486512167_c4406dc98a_z.jpg",
    "webUrl": "https://www.flickr.com/photos/41153475@N04/9486512167/",
    "author": "markyharky",
    "title": "D-AIZP A320 Lufthansa",
    "description": "D-AIZP A320 Lufthansa",
    "license": "AttributionCC",
    "htmlAttributions": [
      "Original of \"<span property='dc:title' itemprop='name'>D-AIZP A320 Lufthansa</span>\" by  <a rel='dc:creator nofollow' property='cc:attributionName' href='https://www.flickr.com/photos/41153475@N04/9486512167/' target='_blank'><span itemprop='producer'>markyharky</span></a>.",
      "Shared and hosted by <span itemprop='publisher'>Flickr</span> under <a target=\"_blank\" rel=\"license cc:license nofollow\" href=\"https://creativecommons.org/licenses/by/2.0/\">CC BY</a>"
    ]
  },
  "numRegistrations": 1
}*/
