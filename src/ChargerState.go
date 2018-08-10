package main

import "fmt"

type ChargerState int16

const (
	_ ChargerState = iota
	Active
	Configuration
	Finished
	Error
)

func (state ChargerState) name() string {
	switch state {
	case Active:
		return "Active"
	case Configuration:
		return "Configuration"
	case Finished:
		return "Finished"
	case Error:
		return "Error"
	default:
		return fmt.Sprintf("State(%d)", state)
	}
}
