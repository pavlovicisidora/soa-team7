package model

// User predstavlja korisnika u našem sistemu.
// Za sada sadrži samo ID, dodaj ostatak kad se poveze gateway za stakeholdera
type Follow struct {
	UserID string `json:"userId"`
}
