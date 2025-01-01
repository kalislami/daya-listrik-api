package models

import "time"

type EnergyRecord struct {
	ID        int       `json:"id"`
	Date      time.Time `json:"date"`
	Usage     float64   `json:"usage"`
	Device    string    `json:"device"`
}