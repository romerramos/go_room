package bills

import "time"

type Model struct {
	ID        int       `json:"id"`
	Amount    int       `json:"amount"`
	DueDate   string    `json:"due_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
