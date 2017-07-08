Install prerequisites:
- libportmidi
- github.com/rakyll/portmidi
- Install Timidity++ per https://help.ubuntu.com/community/Midi/SoftwareSynthesisHowTo

Troubleshooting:
- Does 'timidity MIDI-FILE' play sound?
- Can you start 'timidity -iA and -Os' and hear sound?
- Does 'aplaymidi --port 128:0 MIDI-FILE' play sound?
- See http://haskell.cs.yale.edu/euterpea/midi-on-linux/
