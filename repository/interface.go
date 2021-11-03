package repository

import (
	"context"
	"crud-product/model"
)

type ProductRepository interface {
	Find(context.Context, int) (*model.Product, error)
	Fetch(context.Context, int) ([]model.Product, error)
	Store(context.Context, model.Product) error
	Update(context.Context, model.Product, int) error
	Delete(context.Context, int) error
}

type UserRepository interface {
	FindOne(context.Context, string, string) (model.User, error)
	Store(context.Context, model.User) error
}