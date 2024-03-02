package main

import (
	"encoding/json"
	"net/http"
	"pos-rs/pkg/pos/model"
	"time"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) registerEmployee(w http.ResponseWriter, r *http.Request) {
	var newEmployee model.Employee

	err := json.NewDecoder(r.Body).Decode(&newEmployee)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newEmployee.Password), bcrypt.DefaultCost)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	newEmployee.Password = string(hashedPassword) 

	newEmployee.Enrolled = time.Now().Format(time.RFC3339) 

	err = app.Models.Employee.Register(&newEmployee)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, newEmployee)
}


func (app *Application) logInEmployee(w http.ResponseWriter, r *http.Request) {
	var logInRequest struct {
		Id       string `json:"id"`
		Password string `string:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&logInRequest)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	Employee, err := app.Models.Employee.Get(logInRequest.Id)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(logInRequest.Password), 14)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(Employee.Password), hashedPassword)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	app.respondWithJSON(w, http.StatusOK, logInRequest)
}

func (app *Application) updateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedEmployee model.Employee
	err := json.NewDecoder(r.Body).Decode(&updatedEmployee)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	err = app.Models.Employee.Update(id, &updatedEmployee)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, updatedEmployee)
}

func (app *Application) getAllEmployee(w http.ResponseWriter, r *http.Request) {
	employees, err := app.Models.Employee.GetAll()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusFound, employees)
}

func (app *Application) getEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	Employee, err := app.Models.Employee.Get(param)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Employee Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusFound, Employee)
}

func (app *Application) deleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := app.Models.Employee.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
