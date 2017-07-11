Install prerequisites:
- libportmidi
- github.com/rakyll/portmidi
- Install a synthesizer
  - Ubuntu: https://help.ubuntu.com/community/Midi/SoftwareSynthesisHowTo
    - Timidity++
      - sudo apt-get install timidity
      - Play MIDI-FILE: timidity MIDI-FILE
      - Start as daemon: See synth.sh
      - Problems:
        - Scratchy sound
    - FluidSynth
      - sudo apt-get install fluidsynth fluid-soundfont-gs
      - Play MIDI-FILE: fluidsynth -a alsa -m alsa_seq -g2 -i /usr/share/sounds/sf2/FluidR3_GM.sf2 MIDI-FILE
      - Start as daemon: See synth.sh
      - Problems:
        - Low values of -g are low volume, high values are scratchy.
        - Inexplicably gets persistently mangled notes when experimenting with -g.
        - Sound quality slightly less than Timidity, but acceptable.
        - Persistence after stream is closed:
          - If closed and program exits before note offs sent, tone continues.
          - Pitch bends persist to next MIDI to play.

Troubleshooting:
- Does 'aplaymidi --port 128:0 MIDI-FILE' play sound?
- See http://haskell.cs.yale.edu/euterpea/midi-on-linux/
