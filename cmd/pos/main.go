package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"fmt"
	_ "github.com/lib/pq"
	"pos-rs/pkg/pos/model"
	"github.com/gorilla/mux"
	// "github.com/gorilla/sessions"
)

type Config struct {
	Port string
	Env  string
	DB   struct {
		DSN string
	}
}

type Application struct {
	Config Config
	Models model.Models
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.Port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.DB.DSN, "db-dsn", "postgres://postgres:postgres@localhost/pos_rs?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	db, err := OpenDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &Application{
		Config: cfg,
		Models: model.NewModels(db),
	}

	app.run()
}

func (app *Application) run() {
    r := mux.NewRouter()
    v1 := r.PathPrefix("/api/v1").Subrouter()
	fmt.Println("Running")

    v1.HandleFunc("/employees", app.getAllEmployee).Methods("GET")
    v1.HandleFunc("/employees/{id}", app.getEmployee).Methods("GET")
    v1.HandleFunc("/employees", app.registerEmployee).Methods("POST")
    v1.HandleFunc("/employees/{id}", app.updateEmployee).Methods("PUT")
    v1.HandleFunc("/employees/{id}", app.deleteEmployee).Methods("DELETE")

    v1.HandleFunc("/products", app.getAllProduct).Methods("GET")
    v1.HandleFunc("/products/{productId}", app.getProduct).Methods("GET")
    v1.HandleFunc("/products", app.createProduct).Methods("POST")
    v1.HandleFunc("/products/{productId}", app.updateProduct).Methods("PUT")
    v1.HandleFunc("/products/{productId}", app.deleteProduct).Methods("DELETE")

    log.Printf("Starting server on %s\n", app.Config.Port)
    err := http.ListenAndServe(app.Config.Port, r)
    log.Fatal(err)
}

func OpenDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.DSN)

	if err != nil {
		return nil, err
	}

	return db, nil
}
