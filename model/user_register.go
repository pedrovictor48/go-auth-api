package model

type UserRegister struct {
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	Name      string `json:"name" bson:"name"`
	Birthdate string `json:"birthdate" bson:"birthdate"`
	Gender    string `json:"gender" bson:"gender"`
}
