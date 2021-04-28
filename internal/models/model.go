package models

import "time"

type UserAction struct {
	Metrics   string    `json:"metrics"`
	TargetId  int64     `json:"targetId"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	SubType   string    `json:"subType"`
	UserId    string    `json:"userId"`
}
