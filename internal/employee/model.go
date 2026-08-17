package employee

import "time"

type Employee struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Position  string    `json:"position"`
	Active    bool      `json:"active"`
	UserID    *int64    `json:"userId,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Position string `json:"position"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateRequest struct {
	Name     *string `json:"name"`
	Phone    *string `json:"phone"`
	Position *string `json:"position"`
	Active   *bool   `json:"active"`
}
