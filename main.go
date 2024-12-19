package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Order struct {
	ID           string            `json:"id"`
	HeaderFields map[string]string `json:"headerFields"`
}

type OrderConfirmation struct {
	ID           string            `json:"id"`
	OrderID      string            `json:"orderId"`
	HeaderFields map[string]string `json:"headerFields"`
}

var orders []Order
var orderConfirmations []OrderConfirmation

func loadData() {
	file, err := os.Open("data.json")
	if err != nil {
		log.Println("Error opening data file:", err)
		return
	}
	defer file.Close()

	var data struct {
		Orders             []Order             `json:"orders"`
		OrderConfirmations []OrderConfirmation `json:"orderConfirmations"`
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		log.Println("Error decoding data:", err)
		return
	}

	orders = data.Orders
	orderConfirmations = data.OrderConfirmations
}

func getOrderConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orderConfirmationChan := make(chan *OrderConfirmation)
	orderChan := make(chan *Order)

	go func() {
		for _, oc := range orderConfirmations {
			if oc.ID == id {
				orderConfirmationChan <- &oc
				return
			}
		}
		orderConfirmationChan <- nil
	}()

	go func() {
		for _, order := range orders {
			if order.ID == id {
				orderChan <- &order
				return
			}
		}
		orderChan <- nil
	}()

	var orderConfirmation *OrderConfirmation
	var order *Order

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
		OrderConfirmation *OrderConfirmation `json:"orderConfirmation"`
		Order             *Order             `json:"order"`
	}{
		OrderConfirmation: orderConfirmation,
		Order:             order,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func main() {
	loadData()

	r := mux.NewRouter()
	r.HandleFunc("/api/order-confirmation/{id}", getOrderConfirmationHandler).Methods("GET")

	http.Handle("/", r)
	fmt.Println("Server started at :8051")
	if err := http.ListenAndServe(":8051", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
