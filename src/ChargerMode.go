package main

import "fmt"

type ChargerMode int16

const (
	Charge     ChargerMode = 60
	Storage    ChargerMode = 572
	FastCharge ChargerMode = 828
	Discharge  ChargerMode = 316
	Balance    ChargerMode = 1084
)

func (mode ChargerMode) name() string {
	switch mode {
	case Charge:
		return "Charge"
	case Discharge:
		return "Discharge"
	case Storage:
		return "Storage"
	case FastCharge:
		return "Fast Charge"
	case Balance:
		return "Balance"

	default:
		return fmt.Sprintf("Mode(%d)", mode)
	}
}
