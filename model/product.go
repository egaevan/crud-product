package model

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Price   int `json:"price"`
	Stock   int `json:"stock"`
	BrandID int `json:"brand_id"`
}