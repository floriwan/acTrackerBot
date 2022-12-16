package tracker

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type TrackerError struct {
	desc string
}

func (m TrackerError) Error() string {
	return m.desc
}

var aircraftRegistrations []string

func Add(reg string) error {
	if slices.Contains(aircraftRegistrations, reg) {
		return TrackerError{desc: fmt.Sprintf("registration %v is already in tracked aircraft list", reg)}
	}
	aircraftRegistrations = append(aircraftRegistrations, reg)
	return nil
}

func Remove(reg string) error {
	/*fmt.Printf("remove reg %v\n", reg)
	if slices.Contains(aircraftRegistrations, reg) {
		return TrackerError{desc: fmt.Sprintf("registration %v not found in tracked aircraft list", reg)}
	}*/

	aircraftRegistrations = findAndDelete(aircraftRegistrations, reg)
	fmt.Printf("%v\n", aircraftRegistrations)
	return nil
}

func findAndDelete(l []string, reg string) []string {
	index := 0
	for _, i := range l {
		if i != reg {
			l[index] = i
			index++
		}
	}
	return l[:index]
}

func Size() int {
	return len(aircraftRegistrations)
}

func List() string {
	return strings.Join(aircraftRegistrations, ", ")
}
