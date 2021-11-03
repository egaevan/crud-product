package usecase

import (
	"context"

	"crud-product/model"
	"crud-product/repository"
	log "github.com/sirupsen/logrus"
)

type User struct {
	UserRepo repository.UserRepository
}

func NewUser(userRepo repository.UserRepository) UserUsecase {
	return &User{
		UserRepo: userRepo,
	}
}

func (u *User) Login(ctx context.Context, user model.User) (model.User, error) {

	user, err := u.UserRepo.FindOne(ctx, user.Email, user.Password)
	if err != nil {
		log.Error(err)
		return user, err
	}

	return user, nil
}

func (u *User) CreateUser(ctx context.Context, user model.User) error {
	err := u.UserRepo.Store(ctx, user)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}