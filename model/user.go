package model

type User struct {
	Email     string   `json:"email" bson:"email"`
	Password  string   `json:"password" bson:"password"`
	Name      string   `json:"name" bson:"name"`
	Birthdate string   `json:"birthdate" bson:"birthdate"`
	Gender    string   `json:"gender" bson:"gender"`
	ID        string   `json:"id" bson:"_id"`
	Friends   []string `json:"friends" bson:"friends"`
}
