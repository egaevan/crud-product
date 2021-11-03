package repository

import (
	"context"
	"crud-product/model"
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Product struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &Product{
		DB: db,
	}
}

func (p *Product) Find(ctx context.Context, productID int) (*model.Product, error) {
	query := `
			SELECT 
				product_id,
				name,
				path,
				price,
				stock,
			    brand_id
			FROM 
				product
			WHERE
				product_id = ? AND flag_active = 1`

	prod := model.Product{}

	err := p.DB.QueryRowContext(ctx, query, productID).Scan(&prod.ID, &prod.Name, &prod.Path, &prod.Price, &prod.Stock, &prod.BrandID)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("data not found")
	}

	if err != nil {
		return nil, err
	}
	return &prod, nil
}

func (p *Product) Fetch(ctx context.Context, brandId int) (result []model.Product, err error) {
	query := `
			SELECT 
				product_id,
				name,
				path,
				price,
				stock 
			FROM 
				product
			WHERE
				flag_active = 1 AND brand_id = ?`

	rows, err := p.DB.QueryContext(ctx, query, brandId)
	if err != nil {
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Error(errRow)
		}
	}()

	result = make([]model.Product, 0)

	for rows.Next() {
		t := model.Product{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Path,
			&t.Price,
			&t.Stock,
		)

		if err != nil {
			log.Error(err)
			return nil, err
		}

		result = append(result, t)
	}
	return result, nil
}

func (p *Product) Update(ctx context.Context, product model.Product, brandId int) error {
		query := `
				UPDATE 
				    product
				SET
					name = ?, 
					path = ?, 
					price = ?, 
					stock = ?
				WHERE
					product_id = ?
			`

	_, err := p.DB.ExecContext(ctx, query,
		product.Name, product.Path, product.Price, product.Stock, product.BrandID)

	if err != nil {
		return err
	}
	return nil
}

func (p *Product) Store(ctx context.Context, product model.Product) error {

	query := `
			INSERT INTO product
				(name, path, price, stock, brand_id)
			VALUES
				(?, ?, ?, ?, ?)`

	_, err := p.DB.ExecContext(ctx, query,
		product.Name, product.Path, product.Price, product.Stock, product.BrandID)

	if err != nil {
		return err
	}
	return nil
}

func (p *Product) Delete(ctx context.Context, productID int) error {
	query := `
				UPDATE 
					product
				SET
					flag_active = 0
				WHERE
					product_id = ?
			`

	_, err := p.DB.ExecContext(ctx, query, productID)

	if err != nil {
		return err
	}

	return nil
}