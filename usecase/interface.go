package usecase

import (
	"context"
	"crud-product/model"
)

type ProductUsecase interface {
	GetProduct(context.Context, int) (*model.Product, error)
	GetProductAll(context.Context, int) ([]model.Product, error)
	SendProduct(context.Context, model.Product) (*model.Product, error)
	UpdateProduct(context.Context, model.Product, int) (*model.Product, error)
	DeleteProduct(context.Context, int) error
}

type UserUsecase interface {
	Login(context.Context, model.User) (model.User, error)
	CreateUser(context.Context, model.User) error
}