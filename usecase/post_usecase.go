package usecase

import (
	"auth_api/model"
	"auth_api/repository"
)

type PostUsecase struct {
	repository repository.PostRepository
}

func NewPostUsecase(repo repository.PostRepository) PostUsecase {
	return PostUsecase{
		repository: repo,
	}
}

func (u *PostUsecase) CreatePost(post model.Post) error {
	var err error
	postRepository := u.repository
	err = postRepository.CreatePost(post)
	if err != nil {
		return err
	}
	return nil
}

func (u *PostUsecase) GetPostsById(authorId string) ([]repository.PostResponse, error) {
	var err error
	postRepository := u.repository
	posts, err := postRepository.GetPostsById(authorId)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
