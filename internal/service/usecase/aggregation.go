package usecase

import (
	"fmt"

	"hot-coffee/internal/domain"
)

// Пробует вызвать и при выявлении ошибки выводит то что есть какая то проблема
func (a *Application) GetTotalSales() (float64, error) {
	totalSales, err := a.Repository.GetTotalSales()
	if err != nil {
		return 0, fmt.Errorf("error fetching total sales: %w", err)
	}
	return totalSales, nil
}

// Старается взять самые знаменитые позиции
func (a *Application) GetPopularItems() ([]domain.ProductSales, error) {
	popularItems, err := a.Repository.GetPopularItems()
	if err != nil {
		return nil, fmt.Errorf("error fetching popular items: %w", err)
	}
	return popularItems, nil
}
