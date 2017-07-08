package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rakyll/portmidi"
)

// TODO: Constants, such as for volume

func main() {
	portmidi.Initialize()
	fmt.Printf("CountDevices: %v\n", portmidi.CountDevices())
	fmt.Printf("DefaultInputDevice: %v\n", portmidi.DefaultInputDeviceID())
	fmt.Printf("DefaultOutputDevice: %v\n", portmidi.DefaultOutputDeviceID())
	for device := 0; device < portmidi.CountDevices(); device++ {
		fmt.Printf("Info: %v %+v\n", device, portmidi.Info(portmidi.DeviceID(device)))
	}
	// TODO: Instead of hardcoded 2, search the portmidi.Info for the
	// first port which is not Midi Through Port-0 and
	// IsOutputAvailable.
	out, err := portmidi.NewOutputStream(2, 1024, 0)
	if err != nil {
		log.Fatal(err)
	}

	out.WriteShort(
		0x90, // Note on
		60,   // Middle C
		100)  // Volume
	time.Sleep(1 * time.Second)
	out.WriteShort(0x90, 64, 100)
	time.Sleep(1 * time.Second)
	out.WriteShort(
		0x80, // Note off
		60,
		100)
	// Note off, E, 100 volume
	out.WriteShort(0x80, 64, 100)

	out.Close()

	portmidi.Terminate()
}
