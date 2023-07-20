package models

type Spot struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Website     NullString `json:"website"`
	Coordinates string     `json:"coordinates"`
	Description NullString `json:"description"`
	Rating      float64    `json:"rating"`
}
