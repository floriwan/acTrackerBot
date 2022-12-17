package types

//var FlightawareFlights []Flight

type FlightawareFlights struct {
	Flights []Flight `json:"flights"`
}
type Airport struct {
	Code           string `json:"code"`
	Code_iata      string `json:"code_iata"`
	Timezone       string `json:"timezone"`
	Name           string `json:"name"`
	City           string `json:"city"`
	AirportInfoUrl string `json:"airport_info_url"`
}

type Flight struct {
	Ident       string  `json:"ident"`
	Origin      Airport `json:"origin"`
	Destination Airport `json:"destination"`
}

/*
{
	"flights": [
	  {
		"ident": "DLH690",
		"ident_icao": "DLH690",
		"ident_iata": "LH690",
		"fa_flight_id": "DLH690-1671123360-schedule-0465",
		"actual_off": "2022-12-17T18:02:12Z",
		"actual_on": null,
		"foresight_predictions_available": true,
		"predicted_out": null,
		"predicted_off": null,
		"predicted_on": null,
		"predicted_in": null,
		"predicted_out_source": null,
		"predicted_off_source": null,
		"predicted_on_source": null,
		"predicted_in_source": null,
		"origin": {
		  "code": "EDDF",
		  "code_icao": "EDDF",
		  "code_iata": "FRA",
		  "code_lid": null,
		  "timezone": "Europe/Berlin",
		  "name": "Frankfurt Int'l",
		  "city": "Frankfurt am Main",
		  "airport_info_url": "/airports/EDDF"
		},
		"destination": {
		  "code": "LLBG",
		  "code_icao": "LLBG",
		  "code_iata": "TLV",
		  "code_lid": null,
		  "timezone": "Asia/Jerusalem",
		  "name": "Ben Gurion Int'l",
		  "city": "Tel Aviv",
		  "airport_info_url": "/airports/LLBG"
		},
		"waypoints": [],
		"first_position_time": "2022-12-17T17:46:24Z",
		"last_position": {
		  "fa_flight_id": "DLH690-1671123360-schedule-0465",
		  "altitude": 350,
		  "altitude_change": "-",
		  "groundspeed": 466,
		  "heading": 124,
		  "latitude": 44.71124,
		  "longitude": 18.86666,
		  "timestamp": "2022-12-17T19:11:10Z",
		  "update_type": "A"
		},
		"bounding_box": [
		  50.05,
		  8.52568,
		  44.71124,
		  18.86666
		],
		"ident_prefix": null,
		"aircraft_type": "A21N"
	  }
	],
	"links": null,
	"num_pages": 1
  }
*/
