package main

import (
	"encoding/json"
	"net/http"
	"pos-rs/pkg/pos/model"
	"strconv"
	"github.com/gorilla/mux"
)

func (app *Application) createOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder model.Order

	err := json.NewDecoder((r.Body)).Decode(&newOrder)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalies Request Payload")
		return 
	}

	err = app.Models.Order.Create(&newOrder)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	app.respondWithJSON(w, http.StatusCreated, newOrder)
}

func (app *Application) getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["orderId"]

    orderId, err := strconv.Atoi(param)
    if err != nil {
        app.respondWithError(w, http.StatusBadRequest, "Invalid Order ID")
        return
    }

	Order, err := app.Models.Order.Get(orderId)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Not Found")
	}

	app.respondWithJSON(w, http.StatusFound, Order)
}

func (app *Application) getAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := app.Models.Order.GetAll()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	app.respondWithJSON(w, http.StatusFound, orders)
}

func (app *Application) addProductToOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["orderId"]

    orderId, err := strconv.Atoi(param)
    if err != nil {
        app.respondWithError(w, http.StatusBadRequest, "Invalid Order ID")
        return
    }

	var product model.OrderProduct
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	existingOrder, err := app.Models.Order.Get(orderId)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Order Not Found")
		return
	}


	existingOrder.Products = append(existingOrder.Products, product)
	existingOrder.TotalPrice += float64(product.Price) * float64(product.Qty)

	err = app.Models.Order.Update(orderId, existingOrder)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusOK, existingOrder)
}

func (app *Application) removeProductFromOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["orderId"]

    orderId, err := strconv.Atoi(param)
    if err != nil {
        app.respondWithError(w, http.StatusBadRequest, "Invalid Order ID")
        return
    }

	var product model.OrderProduct
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	existingOrder, err := app.Models.Order.Get(orderId)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Order Not Found")
		return
	}

	var updatedProducts []model.OrderProduct
	var updatedTotalPrice float64
	for _, p := range existingOrder.Products {
		if p.Id != product.Id {
			updatedProducts = append(updatedProducts, p)
			updatedTotalPrice += float64(p.Price)* float64(p.Qty)
		}
	}

	existingOrder.Products = updatedProducts
	existingOrder.TotalPrice = updatedTotalPrice

	err = app.Models.Order.Update(orderId, existingOrder)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusOK, existingOrder)
}

func (app *Application) deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["orderId"]

    orderId, err := strconv.Atoi(param)
    if err != nil {
        app.respondWithError(w, http.StatusBadRequest, "Invalid Order ID")
        return
    }

	err = app.Models.Order.Delete(orderId)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return 
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
