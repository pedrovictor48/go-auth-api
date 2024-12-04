package repository

import (
	"auth_api/model"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoAuthor = errors.New("autor não encontrado")
)

type PostRepository struct {
	client *mongo.Client
}

func NewPostRepository(client *mongo.Client) PostRepository {
	return PostRepository{client}
}

func (p *PostRepository) CreatePost(post model.Post) error {
	var err error
	collection := p.client.Database("testdb").Collection("posts")
	_, err = collection.InsertOne(context.TODO(), post)
	if err != nil {
		return err
	}
	return nil
}

type PostResponse struct {
	Content    string `json:"content" bson:"content"`
	AuthorName string `json:"author" bson:"author"`
	Date       string `json:"date" bson:"date"`
}

func (p *PostRepository) GetPostsById(authorId string) ([]PostResponse, error) {
	var err error
	collection := p.client.Database("testdb").Collection("posts")
	cursor, err := collection.Find(context.Background(), bson.M{"author": authorId})
	if err != nil {
		return nil, err
	}
	var posts []PostResponse

	objectID, err := primitive.ObjectIDFromHex(authorId)
	if err != nil {
		return nil, ErrNoAuthor
	}
	userCollection := p.client.Database("testdb").Collection("users")
	var userAuthor model.User
	err = userCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&userAuthor)
	if err != nil {
		return nil, ErrNoAuthor
	}
	for cursor.Next(context.Background()) {
		var postDB model.Post
		err = cursor.Decode(&postDB)
		if err != nil {
			return nil, err
		}
		post := PostResponse{
			Content:    postDB.Content,
			AuthorName: userAuthor.Name,
			Date:       postDB.Date,
		}
		posts = append(posts, post)
	}
	if err != nil {
		return nil, err
	}
	return posts, nil
}
