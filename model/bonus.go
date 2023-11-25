package model

type Mechanic struct {
	Match      string `json:"match" bson:"_id"`
	Reward     int    `json:"reward" bson:"reward"`
	RewardType string `json:"reward_type" bson:"reward_type"`
}
