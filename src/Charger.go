package main

import (
	"bytes"
	"encoding/binary"
	"github.com/karalabe/hid"
	"time"
)

var ping = []byte{
	15,
	3,
	85,
	0,
	85,
	255,
	255,
}

type Charger struct {
	Device *hid.Device
	OnData func(measurement Measurement)
}

func (charger *Charger) Monitor() {
	for {
		charger.Ping()

		time.Sleep(time.Second)
	}
}

func (charger *Charger) Ping() {
	// Request metrics
	c, err := charger.Device.Write(ping)
	if err != nil {
		panic(err)
	}
	Debug.Printf("Sent ping %s: %d bytes", charger.Device.Product, c)

	// Parse metrics
	result := make([]byte, 40)
	c, err = charger.Device.Read(result)
	if err != nil {
		panic(err)
	}
	Debug.Printf("Received %d bytes", c)

	charger.Process(result[:c])
}

func (charger *Charger) Process(data []byte) {
	buf := bytes.NewBuffer(data)

	// First three bytes are markers
	buf.ReadByte()
	buf.ReadByte()
	buf.ReadByte()

	// Extract all data
	measure := Measurement{
		USBPath: charger.Device.Path,
	}

	readInt16(buf, &measure.State)
	readInt16(buf, &measure.MilliampHoursCharged)

	var chargeTimeSeconds int16;
	readInt16(buf, &chargeTimeSeconds)
	measure.ChargeTime = time.Duration(chargeTimeSeconds) * time.Second

	readInt16(buf, &measure.Current.Millivolts)
	readInt16(buf, &measure.Current.Milliamps)

	var unknown1 int16
	readInt16(buf, &unknown1)
	var unknown2 int16
	readInt16(buf, &unknown2)

	measure.CellInfo.Cells = make([]Cell, 6)

	for i := 0; i < len(measure.CellInfo.Cells); i++ {
		readInt16(buf, &measure.CellInfo.Cells[i])
	}

	var unknown3 int16
	readInt16(buf, &unknown3)
	var unknown4 int16
	readInt16(buf, &unknown4)
	var unknown5 int16
	readInt16(buf, &unknown5)

	readInt16(buf, &measure.CellInfo.Count)
	readInt16(buf, &measure.Mode)

	Debug.Printf(
		"%v",
		measure,
		unknown1,
		unknown2,
		unknown3,
		unknown4,
		unknown5,
		buf.Bytes(),
	)

	if charger.OnData != nil {
		charger.OnData(measure)
	}
}

func readInt16(buf *bytes.Buffer, out interface{}) {
	err := binary.Read(buf, binary.BigEndian, out)
	if err != nil {
		panic(err)
	}
}
