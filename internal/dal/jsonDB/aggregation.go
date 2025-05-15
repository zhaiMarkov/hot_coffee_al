package jsondb

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"hot-coffee/internal/config"
	"hot-coffee/internal/domain"
)

// GetTotalSales вычисляет общую сумму продаж на основе завершенных заказов
func (j *JsonDB) GetTotalSales() (float64, error) {
	// Собираем путь к файлу order.json
	path := filepath.Join(config.Dir, "order.json")

	// Открываем файл
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Декодируем содержимое файла в срез структур domain.Order
	var orders []domain.Order
	if err := json.NewDecoder(file).Decode(&orders); err != nil {
		return 0, err
	}

	// Вычисляем общую сумму продаж
	totalSales := 0.0
	for _, order := range orders {
		if order.Status == domain.StatusCompleted {
			for _, item := range order.Items {
				// Получаем цену товара
				price, err := j.getItemPrice(item.ProductID)
				if err != nil {
					return 0, fmt.Errorf("could not get price for item %s: %w", item.ProductID, err)
				}
				totalSales += float64(item.Quantity) * price
			}
		}
	}
	return totalSales, nil
}

// GetPopularItems находит самые популярные товары (по количеству проданных) из завершенных заказов
func (j *JsonDB) GetPopularItems() ([]domain.ProductSales, error) {
	// Собираем путь к файлу order.json
	path := filepath.Join(config.Dir, "order.json")
	// Открываем файл
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Декодируем содержимое файла в срез структур domain.Order
	var orders []domain.Order
	if err := json.NewDecoder(file).Decode(&orders); err != nil {
		return nil, err
	}

	// Создаем словарь для хранения количества продаж каждого товара
	itemSales := make(map[string]int)
	for _, order := range orders {
		if order.Status == domain.StatusCompleted {
			for _, item := range order.Items {
				itemSales[item.ProductID] += item.Quantity
			}
		}
	}

	// Создаем срез структур domain.ProductSales с популярными товарами
	popularItems := make([]domain.ProductSales, 0)
	for itemID, salesCount := range itemSales {
		popularItems = append(popularItems, domain.ProductSales{
			ProductID: itemID,
			Quantity:  salesCount,
		})
	}
	return popularItems, nil
}

// getItemPrice получает цену определенного товара
func (j *JsonDB) getItemPrice(productID string) (float64, error) {
	// Собираем путь к файлу menu.json
	path := filepath.Join(config.Dir, "menu.json")

	// Открываем файл
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Декодируем содержимое файла в срез структур domain.MenuItem
	var menuItems []domain.MenuItem
	if err := json.NewDecoder(file).Decode(&menuItems); err != nil {
		return 0, err
	}

	// Находим и возвращаем цену товара
	for _, item := range menuItems {
		if item.ID == productID {
			return item.Price, nil
		}
	}

	return 0, errors.New("item not found")
}
