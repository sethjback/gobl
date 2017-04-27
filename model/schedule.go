package model

import "fmt"

// Schedule defines when a backup should run
type Schedule struct {
	ID      string `json:"id,omitempty"`
	JobID   string `json:"jobId"`
	Seconds string `json:"seconds"`
	Minutes string `json:"minutes"`
	Hour    string `json:"hour"`
	DOM     string `json:"dom"`
	MON     string `json:"mon"`
	DOW     string `json:"dow"`
}

func (s *Schedule) String() string {
	return fmt.Sprintf("%s %s %s %s %s %s", s.Seconds, s.Minutes, s.Hour, s.DOM, s.MON, s.DOW)
}
