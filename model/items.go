package model

type Items struct {
	Items []Item `json:"items"`
}
type Item struct {
	Id    int `json:"id"`
	Price int `json:"price"`
	Count int `json:"count"`
}

type Order struct {
	Order string `bson:"order"`
	Bonus int    `bson:"bonus"`
}
