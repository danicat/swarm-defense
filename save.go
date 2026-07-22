package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ScoreData represents the JSON structure for high score persistence.
type ScoreData struct {
	HighScore int `json:"high_score"`
}

const (
	primaryPath  = "/Users/petruzalek/.gemini/antigravity-cli/swarm_defense_score.json"
	fallbackPath = "highscores.json"
)

// LoadHighScore loads the high score from the primary file path, falling back
// to the local fallback file path if needed. Returns 0 if no score can be loaded.
func LoadHighScore() int {
	// Try loading from primary path first
	score, err := loadFromFile(primaryPath)
	if err == nil {
		return score
	}

	// Try loading from fallback path
	score, err = loadFromFile(fallbackPath)
	if err == nil {
		return score
	}

	return 0
}

// SaveHighScore saves the high score to both the primary and fallback paths
// to ensure persistence and redundancy. It handles errors gracefully without panicking.
func SaveHighScore(score int) {
	data := ScoreData{HighScore: score}

	// Try to save to primary path first
	err := saveToFile(primaryPath, data)
	if err != nil {
		// If primary fails, make sure fallback is attempted
		_ = saveToFile(fallbackPath, data)
	} else {
		// If primary succeeds, also update fallback for sync/redundancy
		_ = saveToFile(fallbackPath, data)
	}
}

// loadFromFile reads and decodes the ScoreData from a given path.
func loadFromFile(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data ScoreData
	if err := decoder.Decode(&data); err != nil {
		return 0, err
	}

	return data.HighScore, nil
}

// saveToFile writes the ScoreData to the given path using atomic write where possible,
// ensuring the directory exists and data is correctly flushed to disk.
func saveToFile(path string, data ScoreData) error {
	// Ensure the parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Marshal data with indentation for readability
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Attempt atomic write: write to temp file first, sync, close, rename
	tempFile, err := os.CreateTemp(dir, "swarm_defense_score_temp_*.json")
	if err == nil {
		tempName := tempFile.Name()
		defer os.Remove(tempName) // Cleanup if we return early with error
		defer tempFile.Close()

		if _, err := tempFile.Write(bytes); err == nil {
			if err := tempFile.Sync(); err == nil {
				if err := tempFile.Close(); err == nil {
					if err := os.Rename(tempName, path); err == nil {
						return nil // Atomic write succeeded!
					}
				}
			}
		}
	}

	// Fallback to direct file write if atomic write fails (e.g. due to filesystem limitations)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(bytes); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}

	return nil
}
