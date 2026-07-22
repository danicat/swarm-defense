package main

import (
	"testing"
)

func TestEnvelope(t *testing.T) {
	// Test basic envelope transitions
	// Duration: 1.0s, Attack: 0.1s, Decay: 0.2s, Sustain: 0.5, Release: 0.2s
	// Stage durations are valid: 0.1 + 0.2 + 0.2 = 0.5s <= 1.0s

	// 1. Initial / start of attack
	e := getEnvelope(0.0, 1.0, 0.1, 0.2, 0.5, 0.2)
	if e != 0.0 {
		t.Errorf("Expected 0.0 at t=0, got %f", e)
	}

	// 2. Middle of attack
	e = getEnvelope(0.05, 1.0, 0.1, 0.2, 0.5, 0.2)
	if e != 0.5 {
		t.Errorf("Expected 0.5 at t=0.05, got %f", e)
	}

	// 3. Peak / end of attack
	e = getEnvelope(0.1, 1.0, 0.1, 0.2, 0.5, 0.2)
	if e != 1.0 {
		t.Errorf("Expected 1.0 at t=0.1, got %f", e)
	}

	// 4. Middle of decay (decay from 1.0 to 0.5 over 0.2s; at 0.2s, we are 0.1s into decay, so 50% decay)
	e = getEnvelope(0.2, 1.0, 0.1, 0.2, 0.5, 0.2)
	expectedDecayVal := 1.0 - (1.0-0.5)*0.5
	if mathAbs(e-expectedDecayVal) > 1e-9 {
		t.Errorf("Expected %f at t=0.2, got %f", expectedDecayVal, e)
	}

	// 5. Sustain phase
	e = getEnvelope(0.5, 1.0, 0.1, 0.2, 0.5, 0.2)
	if e != 0.5 {
		t.Errorf("Expected 0.5 in sustain, got %f", e)
	}

	// 6. Release phase (sustain 0.5 to 0.0 over 0.2s starting at t=0.8; at t=0.9, we are 50% into release)
	e = getEnvelope(0.9, 1.0, 0.1, 0.2, 0.5, 0.2)
	if mathAbs(e-0.25) > 1e-9 {
		t.Errorf("Expected 0.25 at t=0.9, got %f", e)
	}

	// 7. Post sound
	e = getEnvelope(1.1, 1.0, 0.1, 0.2, 0.5, 0.2)
	if e != 0.0 {
		t.Errorf("Expected 0.0 post sound, got %f", e)
	}
}

func TestGenerateSound(t *testing.T) {
	params := SynthParams{
		WaveType:    "sine",
		Duration:    0.05,
		StartFreq:   440,
		EndFreq:     440,
		Volume:      0.5,
		Attack:      0.01,
		Decay:       0.01,
		Sustain:     0.8,
		Release:     0.01,
	}

	pcm := generateSound(params)
	expectedLength := int(0.05 * 44100 * 2 * 2) // 44100 samples/sec * 0.05s * 2 channels * 2 bytes/sample
	// Allow a tiny margin due to float rounding of duration * sampleRate
	if len(pcm) != expectedLength {
		t.Errorf("Expected PCM buffer length of %d, got %d", expectedLength, len(pcm))
	}

	// Test that the noise wave generates different bytes
	noiseParams := params
	noiseParams.WaveType = "noise"
	noiseParams.NoiseBlend = 1.0
	noisePcm := generateSound(noiseParams)
	if len(noisePcm) == 0 {
		t.Errorf("Expected noise PCM buffer to be populated, got empty")
	}
}

func TestGenerateSequence(t *testing.T) {
	notes := []SynthParams{
		{
			WaveType:  "square",
			Duration:  0.02,
			StartFreq: 200,
			EndFreq:   200,
			Volume:    0.1,
		},
		{
			WaveType:  "triangle",
			Duration:  0.03,
			StartFreq: 400,
			EndFreq:   400,
			Volume:    0.1,
		},
	}

	seqPcm := generateSequence(notes)
	expectedLen := int((0.02 + 0.03) * 44100 * 4)
	if len(seqPcm) != expectedLen {
		t.Errorf("Expected sequence PCM length of %d, got %d", expectedLen, len(seqPcm))
	}
}

func mathAbs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
