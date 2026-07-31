package cache

import "time"

type Item struct {
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
}
