package domain

import "time"

type Balance struct {
	UserID        int64
	Amount        float64
	LastUpdatedAt time.Time
}
