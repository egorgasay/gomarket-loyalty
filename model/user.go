package model

type User struct {
	Login    string `bson:"_id"`
	Password string `bson:"password"`
	Cookie   string `bson:"cookie"`
	Bonus    int    `bson:"bonus"`
}
