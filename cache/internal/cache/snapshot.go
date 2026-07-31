package cache

import "time"

type Snapshot struct {
	Version int             `json:"version"`
	SavedAt time.Time       `json:"saved_at"`
	Data    map[string]Item `json:"data"`
}
