package eventbus

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const stateFile = ".adxctl/eventbus.state"

type State struct {
	RunnerType string `json:"runner_type"`
	Identifier string `json:"identifier"`
}

func GetStateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, stateFile), nil
}

func SaveState(state State) error {
	path, err := GetStateFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadState() (*State, error) {
	path, err := GetStateFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func ClearState() error {
	path, err := GetStateFilePath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}
