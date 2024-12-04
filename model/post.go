package model

type Post struct {
	Content string `json:"content" bson:"content"`
	Author  string `json:"author" bson:"author"`
	Date    string `json:"date" bson:"date"`
}
