package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"testapp/models"
	"testapp/repositories"

	"github.com/gorilla/mux"
)

func GetOrderConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orderConfirmationChan := make(chan *models.OrderConfirmation)
	orderChan := make(chan *models.Order)

	go func() {
		for _, oc := range repositories.OrderConfirmations {
			if oc.ID == id {
				orderConfirmationChan <- &oc
				return
			}
		}
		orderConfirmationChan <- nil
	}()

	go func() {
		for _, order := range repositories.Orders {
			if order.ID == id {
				orderChan <- &order
				return
			}
		}
		orderChan <- nil
	}()

	var orderConfirmation *models.OrderConfirmation
	var order *models.Order

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
		OrderConfirmation *models.OrderConfirmation `json:"orderConfirmation"`
		Order             *models.Order             `json:"order"`
	}{
		OrderConfirmation: orderConfirmation,
		Order:             order,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
