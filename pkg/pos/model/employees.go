package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Employee struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Password    string `json:"password"`
	IsAdmin     bool   `json:"isAdmin"`
	PhoneNumber string `json:"phoneNumber"`
	Enrolled    string `json:"enrolled"`
}

type EmployeeModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (e EmployeeModel) Register(emp *Employee) error {
	query := `
			INSERT INTO employee (id, name, surname, password, is_admin, phone_number, enrolled) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, password
			`
	args := []interface{}{emp.Id, emp.Name, emp.Surname, emp.Password, emp.IsAdmin, emp.PhoneNumber, emp.Enrolled}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return e.DB.QueryRowContext(ctx, query, args...).Scan(&emp.Id, &emp.Password)
}

func (e EmployeeModel) Get(id string) (*Employee, error) {
	query := `SELECT * FROM employee where id = $1`

	var emp Employee
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := e.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&emp.Id, &emp.Name, &emp.Surname, &emp.Password, &emp.IsAdmin, &emp.PhoneNumber, &emp.Enrolled)

	if err != nil {
		return nil, err
	}

	return &emp, nil
}

func (e EmployeeModel) GetAll() (*[]Employee, error) {
	query := `SELECT * from employee`

	var emp []Employee
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := e.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee Employee
		err := rows.Scan(&employee.Id, &employee.Name, &employee.Surname, &employee.Password, &employee.IsAdmin, &employee.PhoneNumber, &employee.Enrolled)
		if err != nil {
			return nil, err
		}
		emp = append(emp, employee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &emp, nil
}

func (e EmployeeModel) Update(id string, emp *Employee) error {
	query := `
			UPDATE employee 
			SET name = $1, surname = $2, password = $3, is_admin = $4, phone_number = $5
			WHERE id = $6
			RETURNING id, password
	`
	args := []interface{}{emp.Name, emp.Surname, emp.Password, emp.IsAdmin, emp.PhoneNumber, id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return e.DB.QueryRowContext(ctx, query, args...).Scan(&emp.Id, &emp.Password)
}

func (e EmployeeModel) Delete(id string) error {
	query := `DELETE FROM employee WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := e.DB.ExecContext(ctx, query, id)
	return err
}
