package main

import (
	"encoding/json"
	"net/http"
	"pos-rs/pkg/pos/model"
	"fmt"
	"github.com/gorilla/mux"
)

func (app *Application) createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct model.Product

	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	err = app.Models.Product.Create(&newProduct)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusCreated, newProduct)
}

func (app *Application) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["productId"]

	Product, err := app.Models.Product.Get(param)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusFound, Product)
}

func (app *Application) getAllProduct(w http.ResponseWriter, r *http.Request) {
	products, err := app.Models.Product.GetAll()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("Hello from getAllProduct")
	app.respondWithJSON(w, http.StatusFound, products)
}

func (app *Application) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["productId"]

	var updatedProduct model.Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	err = app.Models.Product.Update(id, &updatedProduct)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusOK, updatedProduct)
}

func (app *Application) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["productId"]

	err := app.Models.Product.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

