package main

import (
	"math"
	"math/rand"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// SoundSystem handles retro-synthesized audio effects using pure mathematics
// (sine, square, triangle, and noise waves with customized ADSR envelopes).
type SoundSystem struct {
	context       *audio.Context
	arrowShoot    []byte
	explosion     []byte
	iceFreeze     []byte
	goldMineClick []byte
	unitHurt      []byte
	enemyDeath    []byte
	baseHurt      []byte
	win           []byte
	gameOver      []byte
}

// SynthParams contains configuration options for mathematical waveform rendering.
type SynthParams struct {
	WaveType     string  // "sine", "square", "triangle", "sawtooth"
	Duration     float64 // sound length in seconds
	StartFreq    float64 // starting frequency in Hz
	EndFreq      float64 // ending frequency in Hz (enables pitch sweeping)
	VibratoFreq  float64 // vibrato frequency in Hz
	VibratoDepth float64 // vibrato depth in Hz
	NoiseFilter  float64 // 0.0 to 1.0 (coefficient for lowpass noise filtering)
	NoiseBlend   float64 // 0.0 to 1.0 (blend factor between pure tone and noise)
	DutyCycle    float64 // 0.0 to 1.0 (for square wave pulse width)
	Volume       float64 // 0.0 to 1.0
	Attack       float64 // attack duration in seconds
	Decay        float64 // decay duration in seconds
	Sustain      float64 // sustain amplitude level (0.0 to 1.0)
	Release      float64 // release duration in seconds
}

// NewSoundSystem creates a SoundSystem instance and pre-renders all audio assets in memory.
func NewSoundSystem(ctx *audio.Context) *SoundSystem {
	sys := &SoundSystem{context: ctx}
	if ctx == nil {
		return sys
	}

	// 1. PlayArrowShoot (crispy 8-bit bow-shoot or laser shot)
	sys.arrowShoot = generateSound(SynthParams{
		WaveType:    "square",
		Duration:    0.15,
		StartFreq:   800,
		EndFreq:     150,
		DutyCycle:   0.5,
		Volume:      0.15,
		Attack:      0.01,
		Decay:       0.05,
		Sustain:     0.2,
		Release:     0.09,
	})

	// 2. PlayExplosion (rumbling low-pass filtered noise)
	sys.explosion = generateSound(SynthParams{
		WaveType:    "noise",
		Duration:    0.8,
		NoiseFilter: 0.12,
		NoiseBlend:  1.0,
		Volume:      0.4,
		Attack:      0.02,
		Decay:       0.3,
		Sustain:     0.2,
		Release:     0.48,
	})

	// 3. PlayIceFreeze (crystal-like chime with shimmering vibrato)
	sys.iceFreeze = generateSound(SynthParams{
		WaveType:     "sine",
		Duration:     0.4,
		StartFreq:    1200,
		EndFreq:      2200,
		VibratoFreq:  18,
		VibratoDepth: 150,
		Volume:       0.25,
		Attack:       0.05,
		Decay:        0.1,
		Sustain:      0.4,
		Release:      0.25,
	})

	// 4. PlayGoldMineClick (iconic 8-bit double coin chime)
	sys.goldMineClick = generateSequence([]SynthParams{
		{
			WaveType:  "square",
			Duration:  0.07,
			StartFreq: 988, // B5
			EndFreq:   988,
			DutyCycle: 0.5,
			Volume:    0.15,
			Attack:    0.01,
			Decay:     0.02,
			Sustain:   0.5,
			Release:   0.04,
		},
		{
			WaveType:  "square",
			Duration:  0.22,
			StartFreq: 1319, // E6
			EndFreq:   1319,
			DutyCycle: 0.5,
			Volume:    0.15,
			Attack:    0.01,
			Decay:     0.05,
			Sustain:   0.4,
			Release:   0.16,
		},
	})

	// 5. PlayUnitHurt (sharp, medium-low pitch impact)
	sys.unitHurt = generateSound(SynthParams{
		WaveType:    "sawtooth",
		Duration:    0.15,
		StartFreq:   320,
		EndFreq:     90,
		Volume:      0.25,
		Attack:      0.01,
		Decay:       0.04,
		Sustain:     0.3,
		Release:     0.1,
	})

	// 6. PlayEnemyDeath (crunchy crumbling noise + tone decay)
	sys.enemyDeath = generateSound(SynthParams{
		WaveType:    "sawtooth",
		Duration:    0.3,
		StartFreq:   200,
		EndFreq:     50,
		NoiseBlend:  0.6,
		NoiseFilter: 0.25,
		Volume:      0.3,
		Attack:      0.01,
		Decay:       0.08,
		Sustain:     0.3,
		Release:     0.21,
	})

	// 7. PlayBaseHurt (heavy metallic alarm clonk)
	sys.baseHurt = generateSound(SynthParams{
		WaveType:    "square",
		Duration:    0.4,
		StartFreq:   130,
		EndFreq:     50,
		NoiseBlend:  0.2,
		NoiseFilter: 0.5,
		Volume:      0.4,
		Attack:      0.02,
		Decay:       0.12,
		Sustain:     0.4,
		Release:     0.26,
	})

	// 8. PlayWin (cheerful major arpeggio fanfare)
	sys.win = generateSequence([]SynthParams{
		{
			WaveType:  "triangle",
			Duration:  0.08,
			StartFreq: 523.25, // C5
			EndFreq:   523.25,
			Volume:    0.25,
			Attack:    0.01,
			Decay:     0.02,
			Sustain:   0.7,
			Release:   0.05,
		},
		{
			WaveType:  "triangle",
			Duration:  0.08,
			StartFreq: 659.25, // E5
			EndFreq:   659.25,
			Volume:    0.25,
			Attack:    0.01,
			Decay:     0.02,
			Sustain:   0.7,
			Release:   0.05,
		},
		{
			WaveType:  "triangle",
			Duration:  0.08,
			StartFreq: 783.99, // G5
			EndFreq:   783.99,
			Volume:    0.25,
			Attack:    0.01,
			Decay:     0.02,
			Sustain:   0.7,
			Release:   0.05,
		},
		{
			WaveType:  "triangle",
			Duration:  0.4,
			StartFreq: 1046.50, // C6
			EndFreq:   1046.50,
			Volume:    0.25,
			Attack:    0.02,
			Decay:     0.1,
			Sustain:   0.6,
			Release:   0.28,
		},
	})

	// 9. PlayGameOver (sad descending chromatic/slide)
	sys.gameOver = generateSequence([]SynthParams{
		{
			WaveType:  "square",
			Duration:  0.15,
			StartFreq: 392.00, // G4
			EndFreq:   392.00,
			Volume:    0.25,
			Attack:    0.01,
			Decay:     0.03,
			Sustain:   0.5,
			Release:   0.11,
		},
		{
			WaveType:  "square",
			Duration:  0.15,
			StartFreq: 349.23, // F4
			EndFreq:   349.23,
			Volume:    0.25,
			Attack:    0.01,
			Decay:     0.03,
			Sustain:   0.5,
			Release:   0.11,
		},
		{
			WaveType:  "square",
			Duration:  0.20,
			StartFreq: 311.13, // Eb4
			EndFreq:   311.13,
			Volume:    0.25,
			Attack:    0.01,
			Decay:     0.04,
			Sustain:   0.5,
			Release:   0.15,
		},
		{
			WaveType:  "square",
			Duration:  0.65,
			StartFreq: 261.63, // C4
			EndFreq:   130.81, // descending pitch slide to C3
			Volume:    0.25,
			Attack:    0.02,
			Decay:     0.1,
			Sustain:   0.4,
			Release:   0.53,
		},
	})

	return sys
}

// PlayArrowShoot triggers the arrow shot sound.
func (s *SoundSystem) PlayArrowShoot() {
	s.play(s.arrowShoot)
}

// PlayExplosion triggers the explosion sound.
func (s *SoundSystem) PlayExplosion() {
	s.play(s.explosion)
}

// PlayIceFreeze triggers the freeze sound.
func (s *SoundSystem) PlayIceFreeze() {
	s.play(s.iceFreeze)
}

// PlayGoldMineClick triggers the gold mine click coin sound.
func (s *SoundSystem) PlayGoldMineClick() {
	s.play(s.goldMineClick)
}

// PlayUnitHurt triggers the unit hurt sound.
func (s *SoundSystem) PlayUnitHurt() {
	s.play(s.unitHurt)
}

// PlayEnemyDeath triggers the enemy death sound.
func (s *SoundSystem) PlayEnemyDeath() {
	s.play(s.enemyDeath)
}

// PlayBaseHurt triggers the base damaged sound.
func (s *SoundSystem) PlayBaseHurt() {
	s.play(s.baseHurt)
}

// PlayWin triggers the win victory fanfare.
func (s *SoundSystem) PlayWin() {
	s.play(s.win)
}

// PlayGameOver triggers the game over sad sound.
func (s *SoundSystem) PlayGameOver() {
	s.play(s.gameOver)
}

// play plays the pre-rendered PCM buffer.
func (s *SoundSystem) play(pcm []byte) {
	if s.context == nil || len(pcm) == 0 {
		return
	}
	player := s.context.NewPlayerFromBytes(pcm)
	player.Play()
}

// generateSound creates a raw stereo, 16-bit, 44100Hz PCM byte buffer.
func generateSound(params SynthParams) []byte {
	const sampleRate = 44100
	const numChannels = 2
	const bytesPerSample = 2

	numSamples := int(params.Duration * sampleRate)
	if numSamples <= 0 {
		return nil
	}

	buf := make([]byte, numSamples*numChannels*bytesPerSample)

	var phase float64
	var filterVal float64

	for i := 0; i < numSamples; i++ {
		t := float64(i) / sampleRate
		progress := t / params.Duration

		// Pitch sliding
		var currentFreq float64
		if params.StartFreq == params.EndFreq {
			currentFreq = params.StartFreq
		} else if params.StartFreq > 0 && params.EndFreq > 0 {
			currentFreq = params.StartFreq * math.Pow(params.EndFreq/params.StartFreq, progress)
		} else {
			currentFreq = params.StartFreq + progress*(params.EndFreq-params.StartFreq)
		}

		// Apply Vibrato modulation
		if params.VibratoFreq > 0 && params.VibratoDepth > 0 {
			vibratoOffset := math.Sin(2.0*math.Pi*params.VibratoFreq*t) * params.VibratoDepth
			currentFreq += vibratoOffset
		}

		// Phase accumulator
		phase += currentFreq / sampleRate
		for phase >= 1.0 {
			phase -= 1.0
		}
		for phase < 0.0 {
			phase += 1.0
		}

		// Oscillator waveform rendering
		var oscAmp float64
		switch params.WaveType {
		case "sine":
			oscAmp = math.Sin(2.0 * math.Pi * phase)
		case "square":
			duty := params.DutyCycle
			if duty <= 0 || duty >= 1.0 {
				duty = 0.5
			}
			if phase < duty {
				oscAmp = 1.0
			} else {
				oscAmp = -1.0
			}
		case "triangle":
			if phase < 0.5 {
				oscAmp = 4.0*phase - 1.0
			} else {
				oscAmp = 3.0 - 4.0*phase
			}
		case "sawtooth":
			oscAmp = 2.0*phase - 1.0
		default:
			oscAmp = math.Sin(2.0 * math.Pi * phase)
		}

		// Noise generator
		rawNoise := rand.Float64()*2.0 - 1.0
		filterCoeff := params.NoiseFilter
		if filterCoeff <= 0 {
			filterCoeff = 1.0
		}
		filterVal = filterVal + filterCoeff*(rawNoise-filterVal)
		noiseAmp := filterVal

		// Blend oscillator with noise
		blend := params.NoiseBlend
		if blend < 0 {
			blend = 0
		}
		if blend > 1 {
			blend = 1
		}
		amp := (1.0-blend)*oscAmp + blend*noiseAmp

		// Apply ADSR envelope
		env := getEnvelope(t, params.Duration, params.Attack, params.Decay, params.Sustain, params.Release)

		finalAmp := amp * env * params.Volume
		if finalAmp > 1.0 {
			finalAmp = 1.0
		}
		if finalAmp < -1.0 {
			finalAmp = -1.0
		}

		// Convert float float volume sample to signed 16-bit int PCM
		val := int16(finalAmp * 32760)
		left := uint16(val)
		right := uint16(val)

		idx := i * 4
		buf[idx] = byte(left)
		buf[idx+1] = byte(left >> 8)
		buf[idx+2] = byte(right)
		buf[idx+3] = byte(right >> 8)
	}

	return buf
}

// generateSequence chains multiple sound configurations into a continuous buffer.
func generateSequence(notes []SynthParams) []byte {
	var totalBuf []byte
	for _, note := range notes {
		noteBuf := generateSound(note)
		totalBuf = append(totalBuf, noteBuf...)
	}
	return totalBuf
}

// getEnvelope computes the ADSR multiplier for the given time.
func getEnvelope(t, duration, attack, decay, sustain, release float64) float64 {
	if t < 0 || t > duration {
		return 0.0
	}
	// Fallback to auto-scaling envelope if parameters do not fit
	if attack+decay+release > duration {
		rampUp := duration * 0.1
		if t < rampUp {
			return t / rampUp
		}
		return 1.0 - (t-rampUp)/(duration-rampUp)
	}

	if t < attack {
		if attack > 0 {
			return t / attack
		}
		return 1.0
	}
	if t < attack+decay {
		if decay > 0 {
			progress := (t - attack) / decay
			return 1.0 - (1.0-sustain)*progress
		}
		return sustain
	}
	if t < duration-release {
		return sustain
	}
	// Release phase
	timeInRelease := t - (duration - release)
	if release > 0 {
		progress := timeInRelease / release
		if progress > 1.0 {
			progress = 1.0
		}
		return sustain * (1.0 - progress)
	}
	return 0.0
}
