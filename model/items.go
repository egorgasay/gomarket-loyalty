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
	User  string `bson:"user"`
	Order string `bson:"_id" json:"number"`
	Bonus int    `bson:"bonus" json:"accrual"`
	Time  string `bson:"time" json:"upload_time"`
}

type RequestNameItems struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Query  Query  `json:"query"`
	Sort   string `json:"sort"`
}

type Query struct {
	Ids   []int `json:"ids"`
	Price Price `json:"price"`
}
type Price struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type Names struct {
	Full   []string `json:"full"`
	Partly []string `json:"partly"`
}

type ResponseNameItems struct {
	Offset int       `json:"offset"`
	Total  int       `json:"total"`
	Items  []ItemRes `json:"items"`
}

type ItemRes struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Sizes       []Size `json:"sizes"`
}

type Size struct {
	Size  string `json:"size"`
	Count int    `json:"count"`
}
