package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rakyll/portmidi"
)

// Helpful links:
// portmidi lib docs: http://portmedia.sourceforge.net/portmidi/doxygen/
// MIDI general messages: https://www.midi.org/specifications/item/table-1-summary-of-midi-message
// MIDI Control Change messages: http://nickfever.com/music/midi-cc-list
// Concisely on pitch bends: https://www.midikits.net/midi_analyser/pitch_bend.htm
// Verbosely on pitch bends: http://www.infocellar.com/sound/midi/pitch-bends.htm

// Chroma is a musical note independent of octave.
type Chroma int

// Values are the same as the MIDI notes modulo 12
const (
	C  Chroma = iota // C natural
	Cs               // C sharp
	D                // ...
	Ds
	E
	F
	Fs
	G
	Gs
	A
	As
	B
)

func (chroma Chroma) String() string {
	switch chroma {
	case C:
		return "C"
	case Cs:
		return "C#"
	case D:
		return "D"
	case Ds:
		return "D#"
	case E:
		return "E"
	case F:
		return "F"
	case Fs:
		return "F#"
	case G:
		return "G"
	case Gs:
		return "G#"
	case A:
		return "A"
	case As:
		return "A#"
	case B:
		return "B"
	}
	return ""
}

// Convert user inputted character to its intended chroma.
func InputToChroma(bytes []byte) Chroma {
	switch instring := string(bytes); instring {
	case "c":
		return C
	case "C":
		return Cs
	case "d":
		return D
	case "D":
		return Ds
	case "e":
		return E
	case "f":
		return F
	case "F":
		return Fs
	case "g":
		return G
	case "G":
		return Gs
	case "a":
		return A
	case "A":
		return As
	case "b":
		return B
	default:
		// TODO: Return error code and have program continue
		log.Fatal("Unrecognized input: %s", bytes)
	}
	return C
}

// Lowest musical note eligible to play
const NOTE_LOWER = int64(E + 2*12) // E1

// Highest musical note eligible to play
const NOTE_UPPER = int64(A + 7*12) // A6

// Volume (ie how loud) as passed to portmidi.
//
// NB: a flag to the MIDI synthesizer also affects volume.
const VOLUME = 127

func getRandInclusive(lo int64, hi int64) int64 {
	return lo + rand.Int63n(hi-lo+1)
}

func execCmd(cmd *exec.Cmd) string {
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	portmidi.Initialize()
	log.Printf(
		"portmidi CountDevices: %v DefaultInputDevice: %v DefaultOutputDevice: %v\n",
		portmidi.CountDevices(),
		portmidi.DefaultInputDeviceID(),
		portmidi.DefaultOutputDeviceID())
	for device := 0; device < portmidi.CountDevices(); device++ {
		log.Printf("portmidi DeviceID: %v %+v\n", device, portmidi.Info(portmidi.DeviceID(device)))
	}
	// TODO: Instead of hardcoded 2, search the portmidi.Info for the
	// first port which is not Midi Through Port-0 and
	// IsOutputAvailable.
	out, err := portmidi.NewOutputStream(
		2,
		1024,
		// Latency when opening midi output stream.
		//
		// Using 0 means MIDI events are sent right away, but timestamps
		// are not honored.
		//
		// Using 1 means 1ms delay before MIDI events are sent, but
		// timestamps are honored. Care is necessary to send all note
		// offs, pitch bend resets, etc, because fluidsynth will persist
		// those across MIDI stream sessions (but not across restarts of
		// fluidsynth).
		//
		// The portmidi Pm_OpenOutput doc on has more on this 'latency'
		// field.
		0)
	if err != nil {
		log.Fatal(err)
	}

	// Put tty into mode making single byte input on stdin promptly available.
	//
	// TODO: Maybe use github.com/pkg/term to do this?
	// TODO: Would like to change how input is shown. eg 'eF' -> 'E F#'
	ttyOrig := execCmd(exec.Command("stty", "-F", "/dev/tty", "-g"))
	log.Printf("Original tty: %s", ttyOrig)
	exec.Command("stty", "-F", "/dev/tty", "-icanon", "min", "1").Run()
	defer exec.Command("stty", "-F", "/dev/tty", strings.TrimSpace(ttyOrig)).Run()
	var b []byte = make([]byte, 1)

	correct_queries := 0
	total_queries := 0
	for {
		note0 := getRandInclusive(NOTE_LOWER+12, NOTE_UPPER-12)
		note1 := getRandInclusive(note0-12, note0+12)
		input_str := " "
		for input_str == " " {
			go func(note0, note1 int64) {
				out.WriteShort(0x90, note0, VOLUME)
				time.Sleep(500 * time.Millisecond)
				out.WriteShort(0x80, note0, VOLUME)
				out.WriteShort(0x90, note1, VOLUME)
				time.Sleep(500 * time.Millisecond)
				out.WriteShort(0x80, note1, VOLUME)
			}(note0, note1)
			os.Stdin.Read(b)
			input_str = string(b)
		}
		fmt.Printf("\n")
		inputted_chroma := InputToChroma(b)
		actual_chroma := Chroma(note0 % 12)
		if inputted_chroma == actual_chroma {
			correct_queries++
		}
		total_queries++
		log.Printf(
			"Correct/total: %v/%v. Inputted, actual are %v, %v%v %v%v.\n",
			correct_queries, total_queries,
			inputted_chroma, actual_chroma, note0/12-1,
			Chroma(note1%12), note1/12-1)
	}

	out.WriteShort(0xC0, 1, 0) // instrument 1
	out.WriteShort(0x90, 60, VOLUME)
	time.Sleep(1 * time.Second)
	out.WriteShort(0x80, 60, VOLUME)
	time.Sleep(1 * time.Second)
	out.WriteShort(0x90, 64, VOLUME)
	time.Sleep(1 * time.Second)
	out.WriteShort(0x80, 64, VOLUME)

	// t0 := portmidi.Timestamp(portmidi.Time())
	// out.Write([]portmidi.Event{
	// 	// Set up for pitch bends
	// 	portmidi.Event {
	// 		Timestamp: portmidi.Timestamp(t0),
	// 		Status: 0xB0,  // Control Change
	// 		Data1: 0x64,  // controller number for RPN LSB
	// 		Data2: 0x00,  // controller value (0x7F would reset)
	// 	},
	// 	// Set up for pitch bends
	// 	portmidi.Event {
	// 		Timestamp: portmidi.Timestamp(t0),
	// 		Status: 0xB0,  // Control Change
	// 		Data1: 0x65,  // controller number for RPN MSB
	// 		Data2: 0x00,  // controller value (0x7F would reset)
	// 	},
	// 	portmidi.Event {
	// 		Timestamp: portmidi.Timestamp(t0),
	// 		Status: 0xC0,  // Program Change, channel 0
	// 		Data1: 4,
	// 		Data2: 0,
	// 	},
	// 	portmidi.Event {
	// 		Timestamp: portmidi.Timestamp(t0),
	// 		Status: 0xC1,  // Program Change, channel 1
	// 		Data1: 59,  // Muted trumpet
	// 		Data2: 0,
	// 	},
	// 	portmidi.Event {
	// 		Timestamp: portmidi.Timestamp(t0),
	// 		Status: 0x90,  // Note on, channel 0
	// 		Data1: 60,  // C4
	// 		Data2: VOLUME,
	// 	},
	// 	portmidi.Event {
	// 		Timestamp: portmidi.Timestamp(t0+1000),
	// 		Status: 0x91,  // Note on, channel 1
	// 		Data1: 64,  // E4
	// 		Data2: VOLUME,
	// 	},
	// 	portmidi.Event {
	// 		Timestamp: portmidi.Timestamp(t0+1000),
	// 		Status: 0x90,  // Note on, channel 0
	// 		Data1: 67,  // G
	// 		Data2: VOLUME,
	// 	},
	// 	portmidi.Event {
	// 		Timestamp: portmidi.Timestamp(t0+1100),
	// 		Status: 0x81,  // Note off
	// 		Data1: 64,  // E
	// 		Data2: VOLUME,
	// 	},
	// 	portmidi.Event {
	// 		Timestamp: portmidi.Timestamp(t0+2000),
	// 		Status: 0x80,  // Note off
	// 		Data1: 60,  // C
	// 		Data2: VOLUME,
	// 	},
	// 	portmidi.Event {
	// 		Timestamp: portmidi.Timestamp(t0+2000),
	// 		Status: 0xB0,  // Control Change
	// 		Data1: 0x06,  // controller number for Data Entry
	// 		Data2: 24,  // Pitch bend + or - 12 semitones
	// 	},
	// 	portmidi.Event {
	// 		Timestamp: portmidi.Timestamp(t0+2000),
	// 		Status: 0xE0,  // Pitch bend
	// 		Data1: 0x00,  // LSB
	// 		Data2: 0x00,  // MSB
	// 	},
	// 	portmidi.Event {
	// 		Timestamp: portmidi.Timestamp(t0+3000),
	// 		Status: 0x80,  // Note off
	// 		Data1: 67,  // G
	// 		Data2: VOLUME,
	// 	},
	// })
	// time.Sleep(3 * time.Second)

	out.Close()

	portmidi.Terminate()
}
