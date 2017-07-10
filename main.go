package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rakyll/portmidi"
)

// Helpful links:
// portmidi lib docs: http://portmedia.sourceforge.net/portmidi/doxygen/
// MIDI general messages: https://www.midi.org/specifications/item/table-1-summary-of-midi-message
// MIDI Control Change messages: http://nickfever.com/music/midi-cc-list
// Concisely on pitch bends: https://www.midikits.net/midi_analyser/pitch_bend.htm
// Verbosely on pitch bends: http://www.infocellar.com/sound/midi/pitch-bends.htm

const Volume = 127
// Latency when opening midi output stream. Greater than 0 so as
// timestamp in events are honored. See portmidi Pm_OpenOutput doc on
// 'latency'.
const Latency = 1

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
	out, err := portmidi.NewOutputStream(2, 1024, Latency)
	if err != nil {
		log.Fatal(err)
	}

	t0 := portmidi.Timestamp(portmidi.Time())
	out.Write([]portmidi.Event{
		// Set up for pitch bends
		portmidi.Event {
			Timestamp: portmidi.Timestamp(t0),
			Status: 0xB0,  // Control Change
			Data1: 0x64,  // controller number for RPN LSB
			Data2: 0x00,  // controller value (0x7F would reset)
		},
		// Set up for pitch bends
		portmidi.Event {
			Timestamp: portmidi.Timestamp(t0),
			Status: 0xB0,  // Control Change
			Data1: 0x65,  // controller number for RPN MSB
			Data2: 0x00,  // controller value (0x7F would reset)
		},
		portmidi.Event {
			Timestamp: portmidi.Timestamp(t0),
			Status: 0xC0,  // Program Change, channel 0
			Data1: 0,  // Acoustic Grand Piano
			Data2: 0,
		},
		portmidi.Event {
			Timestamp: portmidi.Timestamp(t0),
			Status: 0xC1,  // Program Change, channel 1
			Data1: 59,  // Muted trumpet
			Data2: 0,
		},
		portmidi.Event {
			Timestamp: portmidi.Timestamp(t0),
			Status: 0x90,  // Note on, channel 0
			Data1: 60,  // C4
			Data2: Volume,
		},
		portmidi.Event {
			Timestamp: portmidi.Timestamp(t0+1000),
			Status: 0x91,  // Note on, channel 1
			Data1: 64,  // E4
			Data2: Volume,
		},
		portmidi.Event {
			Timestamp: portmidi.Timestamp(t0+1000),
			Status: 0x90,  // Note on, channel 0
			Data1: 67,  // G
			Data2: Volume,
		},
		portmidi.Event {
			Timestamp: portmidi.Timestamp(t0+1100),
			Status: 0x81,  // Note off
			Data1: 64,  // E
			Data2: Volume,
		},
		portmidi.Event {
			Timestamp: portmidi.Timestamp(t0+2000),
			Status: 0x80,  // Note off
			Data1: 60,  // C
			Data2: Volume,
		},
		portmidi.Event {
			Timestamp: portmidi.Timestamp(t0+2000),
			Status: 0xB0,  // Control Change
			Data1: 0x06,  // controller number for Data Entry
			Data2: 24,  // Pitch bend + or - 12 semitones
		},
		portmidi.Event {
			Timestamp: portmidi.Timestamp(t0+2000),
			Status: 0xE0,  // Pitch bend
			Data1: 0x00,  // LSB
			Data2: 0x00,  // MSB
		},
		portmidi.Event {
			Timestamp: portmidi.Timestamp(t0+3000),
			Status: 0x80,  // Note off
			Data1: 67,  // G
			Data2: Volume,
		},
	})
	time.Sleep(4 * time.Second)

	out.Close()

	portmidi.Terminate()
}
