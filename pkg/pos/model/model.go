package model

import (
	"database/sql"
	"log"
	"os"
)

type Models struct {
	Employee EmployeeModel
	Product ProductModule
}



func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Employee: EmployeeModel{
			DB:       	db,
			InfoLog:  	infoLog,
			ErrorLog: 	errorLog,
		},
		Product: ProductModule{
			DB:			db,
			InfoLog:	infoLog,
			ErrorLog: 	errorLog,
		},
	}
}