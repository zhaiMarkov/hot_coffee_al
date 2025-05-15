package dal

import "hot-coffee/internal/domain"

// Общий интерфейс хранилища данных
type DataRepository interface {
	OrderRepository
	MenuRepository
	InventoryRepository
	AgreggationRepository
}

// Интерфейс хранилища заказов
type OrderRepository interface {
	// GetOrders получает все заказы
	GetOrders() ([]byte, error)

	// SaveOrders сохраняет заказы
	SaveOrders([]byte) error

	// UnmarshalJsonOrders десериализует заказы из JSON
	UnmarshalJsonOrders(data []byte) ([]*domain.Order, error)

	// MarshalJsonOrders сериализует заказы в JSON
	MarshalJsonOrders(orders []*domain.Order) ([]byte, error)

	// UnmarshalJsonOrderItem десериализует один заказ из JSON
	UnmarshalJsonOrderItem(data []byte) (*domain.Order, error)

	// MarshalJsonOrderItem сериализует один заказ в JSON
	MarshalJsonOrderItem(order *domain.Order) ([]byte, error)
}

// Интерфейс хранилища меню
type MenuRepository interface {
	// GetMenuItems получает все элементы меню
	GetMenuItems() ([]byte, error)

	// SaveMenuItems сохраняет элементы меню
	SaveMenuItems([]byte) error

	// UnmarshalJsonMenuItems десериализует элементы меню из JSON
	UnmarshalJsonMenuItems(data []byte) ([]*domain.MenuItem, error)

	// MarshalJsonMenuItems сериализует элементы меню в JSON
	MarshalJsonMenuItems(menuItems []*domain.MenuItem) ([]byte, error)

	// UnmarshalJsonMenu десериализует один элемент меню из JSON
	UnmarshalJsonMenu(data []byte) (*domain.MenuItem, error)

	// MarshalJsonMenu сериализует один элемент меню в JSON
	MarshalJsonMenu(menuItem *domain.MenuItem) ([]byte, error)
}

// Интерфейс хранилища инвентаря
type InventoryRepository interface {
	GetInventoryItems() ([]byte, error)
	UnmarshalInventoryItems(data []byte) ([]*domain.InventoryItem, error)

	// MarshalInventoryItems сериализует элементы инвентаря в JSON
	MarshalInventoryItems(inventoryItems []*domain.InventoryItem) ([]byte, error)

	// UnmarshalJsonInventory десериализует один элемент инвентаря из JSON
	UnmarshalJsonInventory(data []byte) (*domain.InventoryItem, error)

	// MarshalJsonInventory сериализует один элемент инвентаря в JSON
	MarshalJsonInventory(inventoryItem *domain.InventoryItem) ([]byte, error)
	SaveInventoryItems(data []byte) error
}

// Интерфейс агрегированных данных
type AgreggationRepository interface {
	GetTotalSales() (float64, error)
	GetPopularItems() ([]domain.ProductSales, error)
}
