package domain

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Event represents a timed session structure with multiple phases.
// It serves as a template or "blueprint" for a specific type of session.
type Event struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Total            int     `json:"total"` // total minutes
	WarningValue     int     `json:"warning_value"`
	WarningIsPercent bool    `json:"warning_is_percent"`
	Phases           []Phase `json:"phases"`
}

// EventStore handles the persistence of Event templates and application state.
type EventStore struct {
	path      string
	statePath string
}

// AppState holds persistent application-wide settings.
type AppState struct {
	LastEventID string `json:"last_event_id"`
}

// NewEventStore initializes a new EventStore, creating the configuration directory if it doesn't exist.
func NewEventStore() *EventStore {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".config", "timan")
	os.MkdirAll(configDir, 0755)
	return &EventStore{
		path:      filepath.Join(configDir, "events.json"),
		statePath: filepath.Join(configDir, "state.json"),
	}
}

// SaveState persists the current application state.
func (s *EventStore) SaveState(state AppState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile(s.statePath, data, 0644)
}

// LoadState retrieves the persisted application state.
func (s *EventStore) LoadState() AppState {
	data, err := os.ReadFile(s.statePath)
	if err != nil {
		return AppState{}
	}
	var state AppState
	json.Unmarshal(data, &state)
	return state
}

// SaveAll persists a slice of Events to the store.
func (s *EventStore) SaveAll(events []Event) error {
	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// LoadAll retrieves all persisted Events from the store.
// If the store file doesn't exist, it returns default event templates.
func (s *EventStore) LoadAll() ([]Event, error) {
	// Check for old blueprints.json for backward compatibility
	oldPath := filepath.Join(filepath.Dir(s.path), "blueprints.json")
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		if _, errOld := os.Stat(oldPath); errOld == nil {
			// Migrate from old path
			data, err := os.ReadFile(oldPath)
			if err == nil {
				var events []Event
				if err := json.Unmarshal(data, &events); err == nil {
					s.SaveAll(events)
					// Optional: os.Remove(oldPath) - keeping it for safety for now
					return events, nil
				}
			}
		}
		return s.getDefaults(), nil
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	var events []Event
	err = json.Unmarshal(data, &events)
	return events, err
}

func (s *EventStore) getDefaults() []Event {
	return []Event{
		{
			ID:               "default-session",
			Name:             "Standard Session",
			Total:            60,
			WarningValue:     5,
			WarningIsPercent: false,
			Phases: []Phase{
				{Name: "Intro", Duration: 5 * 60},
				{Name: "Main Content", Duration: 45 * 60},
				{Name: "Wrap-up", Duration: 10 * 60},
			},
		},
		{
			ID:               "standard-interview",
			Name:             "Tech Interview",
			Total:            60,
			WarningValue:     5,
			WarningIsPercent: false,
			Phases: []Phase{
				{Name: "Intro", Duration: 5 * 60},
				{Name: "Understand", Duration: 10 * 60},
				{Name: "Design", Duration: 20 * 60},
				{Name: "Discussion", Duration: 25 * 60},
			},
		},
	}
}
