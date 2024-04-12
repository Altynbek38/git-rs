package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Order struct {
	Id          int            `json:"id"`
	EmployeeID  int            `json:"employee_id"`
	TotalPrice  float64        `json:"total_price"`
	TotalPaid   float64        `json:"total_paid"`
	TotalReturn float64        `json:"total_return"`
	ReceiptID   string         `json:"receipt_id"`
	Products    []OrderProduct `json:"products"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type OrderModule struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (o OrderModule) Create(order *Order) error {
	query := `
			INSERT INTO order (employee_id, total_price, total_paid, total_return,receipt_id, products)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING id
			`
	args := []interface{}{order.EmployeeID, order.TotalPrice, order.TotalPaid, order.TotalReturn, order.ReceiptID, order.Products}
	fmt.Println(args...)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return o.DB.QueryRowContext(ctx, query, args...).Scan(&order.Id)
}

func (o OrderModule) Get(id int) (*Order, error) {
	query := `
			SELECT * FROM order 
			WHERE id = $1
			`
	var order Order
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := o.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&order.Id, &order.EmployeeID, &order.TotalPrice, &order.TotalPaid,
		&order.TotalReturn, &order.ReceiptID, &order.Products, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o OrderModule) GetAll() (*[]Order, error) {
	query := `SELECT * from order`

	var orders []Order
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := o.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ord Order
		err := rows.Scan(&ord.Id, &ord.EmployeeID, &ord.TotalPrice, &ord.TotalPaid,
			&ord.TotalReturn, &ord.ReceiptID, &ord.Products, &ord.CreatedAt, &ord.UpdatedAt)

		if err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &orders, nil
}

func (o OrderModule) Update(id int, order *Order) error {
	query := `
			UPDATE order 
			SET employee_id = $1, total_price = $2, total_paid = $3, total_return = $4, receipt_id = $5, products = $6, updated_at = $7
			WHERE id = $8
			RETURNING updated_at
			`
	args := []interface{}{order.EmployeeID, order.TotalPrice, order.TotalPaid, order.TotalReturn, order.ReceiptID, order.Products, time.Now()}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return o.DB.QueryRowContext(ctx, query, args...).Scan(&order.UpdatedAt)
}

func (o OrderModule) Delete(id int) error {
	query := `
			DELETE FROM order 
			WHERE id = $1
			`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := o.DB.ExecContext(ctx, query, id)
	return err
}
