package model

type User struct {
	Login      string `bson:"_id"`
	SpentBonus int    `bson:"spentBonus"`
	Bonus      int    `bson:"bonus"`
}
