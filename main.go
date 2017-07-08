// Requires installation of timidity

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rakyll/portmidi"
)

func main() {
	portmidi.Initialize()
	fmt.Printf("CountDevices: %v\n", portmidi.CountDevices())
	fmt.Printf("DefaultInputDevice: %v\n", portmidi.DefaultInputDeviceID())
	fmt.Printf("DefaultOutputDevice: %v\n", portmidi.DefaultOutputDeviceID())
	for device := 0; device < portmidi.CountDevices(); device++ {
		fmt.Printf("Info: %v %+v\n", device, portmidi.Info(portmidi.DeviceID(device)))
	}
	out, err := portmidi.NewOutputStream(2, 1024, 0)
	if err != nil {
		log.Fatal(err)
	}

	// note on events to play C major chord
	// out.WriteShort(0x90, 60, 100)
	out.WriteShort(0x90, 64, 100)
	// out.WriteShort(0x90, 67, 100)

	// notes will be sustained for 2 seconds
	time.Sleep(2 * time.Second)

	// note off events
	// out.WriteShort(0x80, 60, 100)
	out.WriteShort(0x80, 64, 100)
	// out.WriteShort(0x80, 67, 100)

	out.Close()

	portmidi.Terminate()
}
