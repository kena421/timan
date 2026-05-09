package domain

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Blueprint represents a full event profile template.
type Blueprint struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Total            int     `json:"total"`             // Total duration in minutes
	Phases           []Phase `json:"phases"`            // Breakdown of the event
	WarningValue     int     `json:"warning_value"`     // Alert threshold value
	WarningIsPercent bool    `json:"warning_is_percent"` // If true, warning is calculated as % of total
}

// AppState persists small pieces of application configuration across restarts.
type AppState struct {
	LastEventID string `json:"last_event_id"` // Tracks the last profile the user selected
}

// EventStore handles the loading and saving of both event templates and app state.
type EventStore struct {
	path      string // Path to events.json
	statePath string // Path to state.json
}

// NewEventStore initializes the local file store in the user's config directory.
func NewEventStore() *EventStore {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".config", "timan")
	os.MkdirAll(dir, 0755)
	
	return &EventStore{
		path:      filepath.Join(dir, "events.json"),
		statePath: filepath.Join(dir, "state.json"),
	}
}

// LoadAll retrieves all saved event templates from disk.
func (s *EventStore) LoadAll() ([]Blueprint, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return []Blueprint{}, nil
	}

	var events []Blueprint
	err = json.Unmarshal(data, &events)
	return events, err
}

// Save persists the provided list of event templates to events.json.
func (s *EventStore) Save(events []Blueprint) error {
	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// LoadState retrieves persistent app configuration (like LastEventID).
func (s *EventStore) LoadState() (AppState, error) {
	data, err := os.ReadFile(s.statePath)
	if err != nil {
		return AppState{}, nil
	}

	var state AppState
	err = json.Unmarshal(data, &state)
	return state, err
}

// SaveState persists small app-level config like the last used profile.
func (s *EventStore) SaveState(state AppState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.statePath, data, 0644)
}
