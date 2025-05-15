package service

import (
	"hot-coffee/internal/domain"
)

type ServiceModule interface {
	OrderService
	MenuService
	InventoryService
	AggregationsService
}

type OrderService interface {
	AddOrder(data []byte) (int, error)
	GetAllOrders() ([]byte, int, error)
	GetOrderByID(id string) ([]byte, int, error)
	UpdateOrderByID(id string, data []byte) (int, error)
	DeleteOrderByID(id string) (int, error)
	CloseOrderByID(id string) (int, error)
}

type MenuService interface {
	AddMenu([]byte) (int, error)
	GetAllMenuItems() ([]byte, int, error)
	GetMenuItemByID(id string) ([]byte, int, error)
	UpdateMenuItemByID(id string, data []byte) (int, error)
	DeleteMenuItemByID(id string) (int, error)
}
type InventoryService interface {
	AddInventoryItem(data []byte) (int, error)
	GetAllInventoryItems() ([]byte, int, error)
	GetInventoryItemByID(id string) ([]byte, int, error)
	UpdateInventoryItemByID(id string, data []byte) (int, error)
	DeleteInventoryItemByID(id string) (int, error)
}

type AggregationsService interface {
	GetTotalSales() (float64, error)
	GetPopularItems() ([]domain.ProductSales, error)
}
