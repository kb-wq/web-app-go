package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"../repositories"
)

func GetOrderConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orderConfirmationChan := make(chan *repositories.OrderConfirmation)
	orderChan := make(chan *repositories.Order)

	go func() {
		orderConfirmation, err := repositories.GetOrderConfirmation(id)
		if err != nil {
			orderConfirmationChan <- nil
			return
		}
		orderConfirmationChan <- orderConfirmation
	}()

	go func() {
		order, err := repositories.GetOrder(id)
		if err != nil {
			orderChan <- nil
			return
		}
		orderChan <- order
	}()

	var orderConfirmation *repositories.OrderConfirmation
	var order *repositories.Order

	select {
	case orderConfirmation = <-orderConfirmationChan:
	case <-ctx.Done():
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
		return
	}

	select {
	case order = <-orderChan:
	case <-ctx.Done():
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
		return
	}

	if orderConfirmation == nil {
		http.Error(w, "Order confirmation not found", http.StatusNotFound)
		return
	}

	response := struct {
		OrderConfirmation *repositories.OrderConfirmation `json:"orderConfirmation"`
		Order             *repositories.Order             `json:"order"`
	}{
		OrderConfirmation: orderConfirmation,
		Order:             order,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
