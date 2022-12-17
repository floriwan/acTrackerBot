package types

//go:generate go run golang.org/x/tools/cmd/stringer -type=AircraftStatus
type AircraftStatus int

/*
run go generate tracker/types/aircraftStatus.go to generate stringer function
*/
const (
	Standing AircraftStatus = iota
	Moving
	Airborn
)
