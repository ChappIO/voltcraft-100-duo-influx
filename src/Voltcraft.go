package main

import (
	"sync"
	"github.com/karalabe/hid"
	"time"
)

func NewVoltcraft() *Voltcraft {
	return &Voltcraft{
		ProductName: "C8051F3xx Development Board",
		lock:        sync.Mutex{},
		foundPaths:  make(map[string]bool),
	}
}

type Voltcraft struct {
	ProductName string
	OnData      func(measurement Measurement)
	foundPaths  map[string]bool
	lock        sync.Mutex
}

func (volt *Voltcraft) StartScanning() {
	for {
		volt.Scan()
		time.Sleep(10 * time.Second)
	}
}

func (volt *Voltcraft) Scan() {
	for _, info := range hid.Enumerate(0, 0) {
		Debug.Printf("Found: %s %s", info.Path, info.Product)
		// Is this a charger?
		if info.Product != volt.ProductName {
			continue
		}

		volt.lock.Lock()
		if volt.foundPaths[info.Path] {
			Debug.Printf("Already monitoring %s", info.Path)
			continue
		} else {
			Info.Printf("Start monitoring %s", info.Path)
			dev, err := info.Open()

			if err != nil {
				panic(err)
			}
			go func() {
				defer func() {
					if r := recover(); r != nil {
						err := r.(error)
						Warn.Printf("Stop monitoring %s: %s", info.Path, err)

						volt.lock.Lock()
						volt.foundPaths[info.Path] = false
						volt.lock.Unlock()
					}
				}()

				volt.foundPaths[dev.Path] = true
				volt.lock.Unlock()

				charger := Charger{
					Device: &*dev,
					OnData: volt.OnData,
				}
				charger.Monitor()
				Info.Print("Done")
			}()
		}
	}
}
