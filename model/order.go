package model

import "time"

type Order struct {
	Kind   string
	Price  float64
	Amount float64
	Time   time.Time
}
