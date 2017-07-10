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
      - sudo apt-get install fluidsynth fluid-soundfont-gs
      - Play MIDI-FILE: fluidsynth -a alsa -m alsa_seq -i /usr/share/sounds/sf2/FluidR3_GM.sf2 MIDI-FILE
      - Start as daemon: fluidsynth -a alsa -m alsa_seq -s -i /usr/share/sounds/sf2/FluidR3_GM.sf2
      - Problems:
        - Sound quality not quite as good as Timidity (modulo timidity's scratchiness)
          - Try other soundfonts?
        - Persistence after stream is closed:
          - If closed and program exits before note offs sent, tone continues.
          - Pitch bends persist to next MIDI to play.

Troubleshooting:
- Does 'aplaymidi --port 128:0 MIDI-FILE' play sound?
- See http://haskell.cs.yale.edu/euterpea/midi-on-linux/
