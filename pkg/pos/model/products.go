package model

import (
	"context"
	"database/sql"
	"log"
	"time"
	"fmt"
)

type Product struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	CategoryId  string `json:"categoryId"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"UpdatedAt"`
}

type ProductModule struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (p ProductModule) Create(product *Product) error {
	fmt.Println("Hello From Product Module")
	query := `
			INSERT INTO product (id, title, category_id, price, description, amount, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
			`
	args := []interface{}{product.Id, product.Title, product.CategoryId, product.Price, product.Description, product.Amount, time.Now(), time.Now()}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	fmt.Println("Buy From Product Module")
	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.Id)
}

func (p ProductModule) Get(id string) (*Product, error) {
	query := `
			SELECT * FROM product 
			WHERE id = $1
			`

	var product Product
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := p.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&product.Id, &product.Title, &product.CategoryId,
		&product.Price, &product.Description, &product.Amount, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (p ProductModule) GetAll() (*[]Product, error) {
    query := `SELECT * from product`

    var products []Product
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    rows, err := p.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var prd Product
        err := rows.Scan(&prd.Id, &prd.Title, &prd.CategoryId, &prd.Price, &prd.Description, &prd.Amount, &prd.CreatedAt, &prd.UpdatedAt)
        if err != nil {
            return nil, err
        }
        products = append(products, prd)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return &products, nil
}


func (p ProductModule) Update(id string, product *Product) error {
	query := `
			UPDATE product menu 
			SET title = $1, category_id = $2, price = $3, description = $4, amount = $5
			WHERE id = $6
			RETURNING updated_at
			`
	args := []interface{}{product.Title, product.CategoryId, product.Price, product.Description, product.Amount, id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.UpdatedAt)
}

func (p ProductModule) Delete(id string) error {
	query := `
			DELETE FROM product 
			WHERE id = $1
			`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, id)
	return err
}
