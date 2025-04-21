package models

import "time"

type EnergyRecord struct {
	ID        int       `json:"id"`
	Date      time.Time `json:"date"`
	Usage     float64   `json:"usage"`
	Duration  float64   `json:"duration"`
	Device    string    `json:"device"`
}