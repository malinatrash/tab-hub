package models

type Project struct {
	ID      int    `db:"id"`
	Name    string `db:"name"`
	OwnerID int    `db:"owner_id"`
	State   []byte `db:"state"`
	Private bool   `db:"private"`
}
