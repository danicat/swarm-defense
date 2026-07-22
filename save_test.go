package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadHighScore(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "swarm_defense_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testPath := filepath.Join(tempDir, "highscores.json")

	// Test case 1: File does not exist initially
	score, err := loadFromFile(testPath)
	if err == nil {
		t.Errorf("expected error when loading non-existent file, got none")
	}
	if score != 0 {
		t.Errorf("expected score 0, got %d", score)
	}

	// Test case 2: Save and load score
	expectedScore := 4250
	data := ScoreData{HighScore: expectedScore}
	err = saveToFile(testPath, data)
	if err != nil {
		t.Fatalf("failed to save score: %v", err)
	}

	loadedScore, err := loadFromFile(testPath)
	if err != nil {
		t.Fatalf("failed to load score: %v", err)
	}

	if loadedScore != expectedScore {
		t.Errorf("expected loaded score %d, got %d", expectedScore, loadedScore)
	}

	// Test case 3: Overwrite with higher score
	newHighScore := 9999
	newData := ScoreData{HighScore: newHighScore}
	err = saveToFile(testPath, newData)
	if err != nil {
		t.Fatalf("failed to save new score: %v", err)
	}

	loadedScore, err = loadFromFile(testPath)
	if err != nil {
		t.Fatalf("failed to load new score: %v", err)
	}

	if loadedScore != newHighScore {
		t.Errorf("expected loaded new score %d, got %d", newHighScore, loadedScore)
	}
}
