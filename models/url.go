package models

type URL struct {
	ID       int64  `json:"-" bson:"_id"`
	Code     string `json:"code" bson:"-"`
	URL      string `json:"url" bson:"url"`
	ExpireAt int64  `json:"expireAt" bson:"expireAt"`
}
