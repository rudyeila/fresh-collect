package model

type Ingredient struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Amount  *float64 `json:"amount"`
	Unit    string   `json:"unit"`
	Shipped bool     `json:"shipped"`
}

type Recipe struct {
	ID          string
	Country     string
	Ingredients []Ingredient
}
