package businesshours

import "time"

type BusinessHours struct {
	ID        int64     `json:"id"`
	DayOfWeek int       `json:"dayOfWeek"`
	OpenTime  string    `json:"openTime"`
	CloseTime string    `json:"closeTime"`
	IsClosed  bool      `json:"isClosed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpsertItem struct {
	DayOfWeek int    `json:"dayOfWeek"`
	OpenTime  string `json:"openTime"`
	CloseTime string `json:"closeTime"`
	IsClosed  bool   `json:"isClosed"`
}
