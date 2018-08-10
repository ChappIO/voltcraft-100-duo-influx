package main

import (
	"time"
	"fmt"
)

type Measurement struct {
	USBPath              string
	State                ChargerState
	Mode                 ChargerMode
	MilliampHoursCharged int16
	ChargeTime           time.Duration
	Current              struct {
		Millivolts int16
		Milliamps  int16
	}
	CellInfo Cells
}

type Cells struct {
	Count int16
	Cells []Cell
}

type Cell struct {
	Millivolts int16
}

func (measure Measurement) String() string {
	return fmt.Sprintf(
		"%s [%s] %8s - %4dmAh %4dmA %4dmV %d cells%v  ",
		measure.Mode.name(),
		measure.State.name(),
		measure.ChargeTime,
		measure.MilliampHoursCharged,
		measure.Current.Milliamps,
		measure.Current.Millivolts,
		measure.CellInfo.Count,
		measure.CellInfo.Cells,
	)
}
