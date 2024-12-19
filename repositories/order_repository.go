package repositories

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"

	"../models"
)

var (
	orders             []models.Order
	orderConfirmations []models.OrderConfirmation
	mutex              sync.Mutex
)

func loadData() error {
	file, err := os.Open("data.json")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&orders)
	if err != nil {
		return err
	}

	err = decoder.Decode(&orderConfirmations)
	if err != nil {
		return err
	}

	return nil
}

func saveData() error {
	mutex.Lock()
	defer mutex.Unlock()

	data := struct {
		Orders             []models.Order             `json:"orders"`
		OrderConfirmations []models.OrderConfirmation `json:"orderConfirmations"`
	}{
		Orders:             orders,
		OrderConfirmations: orderConfirmations,
	}

	file, err := os.Create("data.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func GetOrder(id string) (*models.Order, error) {
	for _, order := range orders {
		if order.ID == id {
			return &order, nil
		}
	}
	return nil, errors.New("order not found")
}

func CreateOrder(order models.Order) error {
	orders = append(orders, order)
	return saveData()
}

func UpdateOrder(order models.Order) error {
	for i, o := range orders {
		if o.ID == order.ID {
			orders[i] = order
			return saveData()
		}
	}
	return errors.New("order not found")
}

func DeleteOrder(id string) error {
	for i, order := range orders {
		if order.ID == id {
			orders = append(orders[:i], orders[i+1:]...)
			return saveData()
		}
	}
	return errors.New("order not found")
}

func GetOrderConfirmation(id string) (*models.OrderConfirmation, error) {
	for _, oc := range orderConfirmations {
		if oc.ID == id {
			return &oc, nil
		}
	}
	return nil, errors.New("order confirmation not found")
}

func CreateOrderConfirmation(orderConfirmation models.OrderConfirmation) error {
	orderConfirmations = append(orderConfirmations, orderConfirmation)
	return saveData()
}

func UpdateOrderConfirmation(orderConfirmation models.OrderConfirmation) error {
	for i, oc := range orderConfirmations {
		if oc.ID == orderConfirmation.ID {
			orderConfirmations[i] = orderConfirmation
			return saveData()
		}
	}
	return errors.New("order confirmation not found")
}

func DeleteOrderConfirmation(id string) error {
	for i, oc := range orderConfirmations {
		if oc.ID == id {
			orderConfirmations = append(orderConfirmations[:i], orderConfirmations[i+1:]...)
			return saveData()
		}
	}
	return errors.New("order confirmation not found")
}
