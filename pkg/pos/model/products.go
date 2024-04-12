package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Product struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	CategoryId  int    `json:"categoryId"`
	Price       int       `json:"price"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}

type ProductModule struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (p ProductModule) Create(product *Product) error {
	fmt.Println("Hello From Product Module")
	query := `
			INSERT INTO products (name, category_id, price, description, amount)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
			`
	args := []interface{}{product.Name, product.CategoryId, product.Price, product.Description, product.Amount}
	fmt.Println(args...)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	fmt.Println("Buy From Product Module")
	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.Id)
}

func (p ProductModule) Get(id int) (*Product, error) {
	query := `
			SELECT * FROM products 
			WHERE id = $1
			`

	var product Product
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := p.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&product.Id, &product.Name, &product.CategoryId,
		&product.Price, &product.Description, &product.Amount, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (p ProductModule) GetAll() (*[]Product, error) {
	query := `SELECT * from products`

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
		err := rows.Scan(&prd.Id, &prd.Name, &prd.CategoryId, &prd.Price, &prd.Description, &prd.Amount, &prd.CreatedAt, &prd.UpdatedAt)
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

func (p ProductModule) Update(id int, product *Product) error {
	query := `
			UPDATE products
			SET name = $1, category_id = $2, price = $3, description = $4, amount = $5
			WHERE id = $6
			RETURNING updated_at
			`
	args := []interface{}{product.Name, product.CategoryId, product.Price, product.Description, product.Amount, id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.DB.QueryRowContext(ctx, query, args...).Scan(&product.UpdatedAt)
}

func (p ProductModule) Delete(id int) error {
	query := `
			DELETE FROM products
			WHERE id = $1
			`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, id)
	return err
}
