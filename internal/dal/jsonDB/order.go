package jsondb

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"hot-coffee/internal/config"
	"hot-coffee/internal/domain"
)

// Получение заказов из файла order.json
func (j *JsonDB) GetOrders() ([]byte, error) {
	path := filepath.Join(config.Dir, "order.json")

	file, err := os.OpenFile(path, os.O_RDONLY, 0o755)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Сохранение заказов в файл order.json
func (j *JsonDB) SaveOrders(data []byte) error {
	path := filepath.Join(config.Dir, "order.json")

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// Десериализация массива заказов из JSON
func (j *JsonDB) UnmarshalJsonOrders(data []byte) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := json.Unmarshal(data, &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// Сериализация массива заказов в JSON
func (j *JsonDB) MarshalJsonOrders(orders []*domain.Order) ([]byte, error) {
	ordersJson, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}

	return ordersJson, nil
}

// Десериализация одного заказа из JSON
func (j *JsonDB) UnmarshalJsonOrderItem(data []byte) (*domain.Order, error) {
	var order domain.Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// Сериализация одного заказа в JSON
func (j *JsonDB) MarshalJsonOrderItem(order *domain.Order) ([]byte, error) {
	orderJson, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	return orderJson, nil
}
