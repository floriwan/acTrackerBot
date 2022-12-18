package tracker

import (
	"acTrackerBot/tracker/types"
	"log"
)

type aircraftData struct {
	LastAdsbExchData    types.AdsbExchData
	CurrentAdsbExchData types.AdsbExchData
	LastAircraftInfo    types.AircraftInformation
	CurrentAircraftInfo types.AircraftInformation
}

var lastAircraftInfo map[string]*aircraftData

func init() {
	lastAircraftInfo = make(map[string]*aircraftData)
}

func processData(reg string, data types.AdsbExchData) {

	_, ok := lastAircraftInfo[reg]
	if !ok {
		addAircraft(reg, data)
	}

	updateAircraft(reg, data)
}

func getCurrentAircraftInfo(reg string) types.AircraftInformation {
	return lastAircraftInfo[reg].CurrentAircraftInfo
}

func newStatus(reg string) bool {
	data, ok := lastAircraftInfo[reg]
	if !ok {
		log.Printf("aircraft ref '%v' not found in history data", reg)
		return false
	}

	if data.CurrentAircraftInfo.Status != data.LastAircraftInfo.Status {
		return true
	}
	return false
}

func updateAircraft(reg string, data types.AdsbExchData) {
	acData := lastAircraftInfo[reg]

	acData.LastAircraftInfo = acData.CurrentAircraftInfo
	acData.LastAdsbExchData = acData.CurrentAdsbExchData

	acData.CurrentAdsbExchData = data

	if len(data.Ac) == 0 {
		//log.Printf("no data available for %v", reg)
		acData.CurrentAircraftInfo.Status = types.Parking
		return
	}

	//log.Printf("new data gs:%v altGeom:%v, altBaro:%v", data.Ac[0].Gs, data.Ac[0].Alt_geom, data.Ac[0].Alt_baro)

	acData.CurrentAircraftInfo.Speed = data.Ac[0].Gs
	acData.CurrentAircraftInfo.AltGeom = data.Ac[0].Alt_geom
	acData.CurrentAircraftInfo.Lat = data.Ac[0].Lat
	acData.CurrentAircraftInfo.Lon = data.Ac[0].Lon
	acData.CurrentAircraftInfo.Callsign = data.Ac[0].Flight
	acData.CurrentAircraftInfo.Squawk = data.Ac[0].Squawk

	if data.Ac[0].Gs > 0 {
		acData.CurrentAircraftInfo.Status = types.Taxing
	}

	if data.Ac[0].Gs > 0 && data.Ac[0].Alt_geom > 0 {
		acData.CurrentAircraftInfo.Status = types.Airborn
	}

	//log.Printf("prev status:%v, gs:%v -> current status:%v, gs:%v",
	//	acData.LastAircraftInfo.Status, acData.LastAdsbExchData.Ac[0].Gs,
	//	acData.CurrentAircraftInfo.Status, acData.CurrentAdsbExchData.Ac[0].Gs)

}

func addAircraft(reg string, data types.AdsbExchData) {
	acData := aircraftData{CurrentAircraftInfo: types.AircraftInformation{Reg: reg, Status: types.Parking}}
	acData.CurrentAdsbExchData = data
	lastAircraftInfo[reg] = &acData
}
