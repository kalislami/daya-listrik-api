package models

import (
	"time"
)

type EnergyRecord struct {
	ID       int      `gorm:"primaryKey;autoIncrement"`
	Date     time.Time `gorm:"autoCreateTime"`
	Usage    float64
	Device   string
	Duration float64 `gorm:"default:1"`
}
