package model

import "mime/multipart"

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name" form:"name" `
	Path string `json:"-"`
	UrlImage *multipart.FileHeader `json:"url_image" form:"url_image"`
	Price   int `json:"price" form:"price"`
	Stock   int `json:"stock" form:"stock"`
	BrandID int `json:"brand_id" form:"brand_id"`
}