package models

import "time"

type Project struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	OwnerID   int       `db:"owner_id" json:"owner_id"`
	State     string    `db:"state" json:"state"`
	Private   bool      `db:"private" json:"private"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
