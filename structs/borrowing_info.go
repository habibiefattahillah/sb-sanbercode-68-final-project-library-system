package structs

import "time"

type BorrowingInfo struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	BorrowedAt time.Time  `json:"borrowed_at"`
	ReturnedAt *time.Time `json:"returned_at,omitempty"`
}
