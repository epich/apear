Install prerequisites:
- libportmidi
- github.com/rakyll/portmidi
- Install a synthesizer
  - Ubuntu: https://help.ubuntu.com/community/Midi/SoftwareSynthesisHowTo
    - Timidity++
      - sudo apt-get install timidity
      - Play MIDI-FILE: timidity MIDI-FILE
      - Start as daemon: timidity -iA -Os
      - Problems:
        - Scratchy sound
          - Install from source to see if problem still exists?
    - FluidSynth
      - sudo apt-get install fluidsynth qsynth fluid-soundfont-gs
      - Play MIDI-FILE: fluidsynth -a alsa -m alsa_seq -i /usr/share/sounds/sf2/FluidR3_GM.sf2 MIDI-FILE
      - Start as daemon: fluidsynth -a alsa -m alsa_seq -s -i /usr/share/sounds/sf2/FluidR3_GM.sf2
      - Problems:
        - Instrument 0 (Acoustic Grand Piano) sounds more like a harpsichord.
          - Possibly try other soundfonts.
        - Pitch Bends persist across MIDI stream sessions.
          - Possible workaround (untested): Reset pitch bends either at program start or all program exits

Troubleshooting:
- Does 'aplaymidi --port 128:0 MIDI-FILE' play sound?
- See http://haskell.cs.yale.edu/euterpea/midi-on-linux/
