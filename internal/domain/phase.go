package domain

import (
	"fmt"
	"strconv"
	"strings"
)

type Phase struct {
	Name      string
	Duration  int  // seconds (absolute)
	IsPercent bool
	Percent   int  // percentage value (if IsPercent is true)
}

func (p Phase) GetSeconds(totalMins int) int {
	if p.IsPercent {
		return (p.Percent * totalMins * 60) / 100
	}
	return p.Duration
}

func (p Phase) FormatDuration(totalMins int) string {
	s := p.GetSeconds(totalMins)
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}

func ParsePhases(input string) ([]Phase, error) {
	var phases []Phase
	parts := strings.Split(input, ",")
	for _, part := range parts {
		kv := strings.Split(strings.TrimSpace(part), ":")
		if len(kv) == 2 {
			name := strings.TrimSpace(kv[0])
			mins, err := strconv.Atoi(strings.TrimSpace(kv[1]))
			if err == nil {
				phases = append(phases, Phase{Name: name, Duration: mins * 60})
			}
		}
	}
	if len(phases) == 0 {
		return nil, fmt.Errorf("no valid phases found")
	}
	return phases, nil
}

func ValidatePhases(phases []Phase, totalMinutes int) error {
	sumSeconds := 0
	for _, p := range phases {
		sumSeconds += p.GetSeconds(totalMinutes)
	}
	if sumSeconds != totalMinutes*60 {
		return fmt.Errorf("allocated time (%d mins) does not match total duration (%d mins)", sumSeconds/60, totalMinutes)
	}
	return nil
}
