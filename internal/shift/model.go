package shift

import "time"

type ShiftType struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	StartTime string    `json:"startTime"`
	EndTime   string    `json:"endTime"`
	Color     string    `json:"color,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateRequest struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Color     string `json:"color"`
}

type UpdateRequest struct {
	Name      *string `json:"name"`
	StartTime *string `json:"startTime"`
	EndTime   *string `json:"endTime"`
	Color     *string `json:"color"`
}
