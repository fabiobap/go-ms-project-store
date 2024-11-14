package dto

import "time"

type NewTokenDTO struct {
	UserID    uint64
	Name      string
	ExpiresAt time.Time
	Abilities []string
}
