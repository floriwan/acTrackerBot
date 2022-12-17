package types

type AircraftInformation struct {
	Reg      string
	IcaoType string
	Status   AircraftStatus
	Lat      float32
	Lon      float32
	AltGeom  int
	Speed    float32
}
