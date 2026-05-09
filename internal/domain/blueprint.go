package domain

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Blueprint struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Total            int     `json:"total"` // total minutes
	WarningValue     int     `json:"warning_value"`
	WarningIsPercent bool    `json:"warning_is_percent"`
	Phases           []Phase `json:"phases"`
}

type BlueprintStore struct {
	path string
}

func NewBlueprintStore() *BlueprintStore {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".config", "timan")
	os.MkdirAll(configDir, 0755)
	return &BlueprintStore{
		path: filepath.Join(configDir, "blueprints.json"),
	}
}

func (s *BlueprintStore) SaveAll(blueprints []Blueprint) error {
	data, err := json.MarshalIndent(blueprints, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *BlueprintStore) LoadAll() ([]Blueprint, error) {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		return s.getDefaults(), nil
	}
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	var blueprints []Blueprint
	err = json.Unmarshal(data, &blueprints)
	return blueprints, err
}

func (s *BlueprintStore) getDefaults() []Blueprint {
	return []Blueprint{
		{
			ID:               "default-interview",
			Name:             "Standard Interview",
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
