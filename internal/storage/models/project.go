package models

type Project struct {
	Name    string
	OwnerID int
	State   []byte
	Private bool
}
