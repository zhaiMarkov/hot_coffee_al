package usecase

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"hot-coffee/internal/domain"
)

func (a *Application) AddOrder(data []byte) (int, error) {
	order, err := a.Repository.UnmarshalJsonOrderItem(data)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid order data")
	}

	order.ID = generateOrderID()
	order.Status = domain.StatusPending
	order.CreatedAt = time.Now()
	if err := CheckOrderFields(order); err != nil {
		return http.StatusBadRequest, err
	}

	// Get menu items
	menuData, err := a.Repository.GetMenuItems()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error getting menu items")
	}

	menuItems, err := a.Repository.UnmarshalJsonMenuItems(menuData)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error unmarshalling menu items")
	}

	// Get inventory items
	inventoryData, err := a.Repository.GetInventoryItems()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error getting inventory items")
	}

	inventoryItems, err := a.Repository.UnmarshalInventoryItems(inventoryData)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error unmarshalling inventory items")
	}

	for _, item := range order.Items {
		menuItem := findMenuItem(item.ProductID, menuItems)
		if menuItem == nil {
			return http.StatusBadRequest, fmt.Errorf("menu item %s not found", item.ProductID)
		}
		if !hasIngredient(menuItem.Ingredients, inventoryItems) {
			return http.StatusConflict, fmt.Errorf("ingredient for menu item %s not found in inventory", menuItem.ID)
		}
		if !checkIngredientsAvailability(item.Quantity, menuItem.Ingredients, inventoryItems) {
			return http.StatusConflict, fmt.Errorf("insufficient ingredients for menu item %s", menuItem.ID)
		}
	}

	// Save the order
	allData, err := a.Repository.GetOrders()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error getting orders")
	}

	orderList, err := a.Repository.UnmarshalJsonOrders(allData)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error unmarshalling orders")
	}

	orderList = append(orderList, order)

	ordersJson, err := a.Repository.MarshalJsonOrders(orderList)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error marshalling orders")
	}

	err = a.Repository.SaveOrders(ordersJson)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error saving orders")
	}

	return http.StatusCreated, nil
}

func (a *Application) GetAllOrders() ([]byte, int, error) {
	// Get all orders
	allData, err := a.Repository.GetOrders()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Check if there are no orders
	if len(allData) <= 2 {
		return nil, http.StatusNotFound, fmt.Errorf("no orders found")
	}

	return allData, http.StatusOK, nil
}

func (a *Application) GetOrderByID(id string) ([]byte, int, error) {
	// Get all orders
	allData, err := a.Repository.GetOrders()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Unmarshal the JSON orders
	orders, err := a.Repository.UnmarshalJsonOrders(allData)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Check if the order exists
	for _, item := range orders {
		if item.ID == id {
			orderJson, err := a.Repository.MarshalJsonOrderItem(item)
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}
			return orderJson, http.StatusOK, nil
		}
	}

	return nil, http.StatusNotFound, fmt.Errorf("order with ID %s not found", id)
}

func (a *Application) UpdateOrderByID(id string, data []byte) (int, error) {
	// Unmarshal the JSON order
	newOrder, err := a.Repository.UnmarshalJsonOrderItem(data)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid order data")
	}

	newOrder.ID = id
	newOrder.Status = domain.StatusPending
	newOrder.CreatedAt = time.Now()
	// Check if all fields are set
	if err := CheckOrderFields(newOrder); err != nil {
		return http.StatusBadRequest, err
	}

	// Get menu items
	menuData, err := a.Repository.GetMenuItems()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error getting menu items")
	}

	menuItems, err := a.Repository.UnmarshalJsonMenuItems(menuData)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error unmarshalling menu items")
	}

	// Get inventory items
	inventoryData, err := a.Repository.GetInventoryItems()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error getting inventory items")
	}

	inventoryItems, err := a.Repository.UnmarshalInventoryItems(inventoryData)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error unmarshalling inventory items")
	}

	for _, item := range newOrder.Items {
		menuItem := findMenuItem(item.ProductID, menuItems)
		if menuItem == nil {
			return http.StatusBadRequest, fmt.Errorf("menu item %s not found", item.ProductID)
		}

		if !hasIngredient(menuItem.Ingredients, inventoryItems) {
			return http.StatusConflict, fmt.Errorf("ingredient for menu item %s not found in inventory", menuItem.ID)
		}

		if !checkIngredientsAvailability(item.Quantity, menuItem.Ingredients, inventoryItems) {
			return http.StatusConflict, fmt.Errorf("insufficient ingredients for menu item %s", menuItem.ID)
		}
	}

	// Get all orders
	allData, err := a.Repository.GetOrders()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Unmarshal the JSON orders
	orders, err := a.Repository.UnmarshalJsonOrders(allData)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Update the order
	for i, item := range orders {
		if item.ID == id {
			if item.Status != domain.StatusPending {
				return http.StatusConflict, fmt.Errorf("order %s is already completed", id)
			}
			orders[i] = newOrder
			break
		}
	}

	// Marshal the JSON orders
	ordersJson, err := a.Repository.MarshalJsonOrders(orders)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Save the orders
	err = a.Repository.SaveOrders(ordersJson)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (a *Application) DeleteOrderByID(id string) (int, error) {
	// Get all orders
	allData, err := a.Repository.GetOrders()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Unmarshal the JSON orders
	orders, err := a.Repository.UnmarshalJsonOrders(allData)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Delete the order
	for i, item := range orders {
		if item.ID == id {
			orders = append(orders[:i], orders[i+1:]...)
			break
		}
	}

	// Marshal the JSON orders
	ordersJson, err := a.Repository.MarshalJsonOrders(orders)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Save the orders
	err = a.Repository.SaveOrders(ordersJson)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusNoContent, nil
}

func (a *Application) CloseOrderByID(id string) (int, error) {
	// Get all orders
	allData, err := a.Repository.GetOrders()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Unmarshal the JSON orders
	orders, err := a.Repository.UnmarshalJsonOrders(allData)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Find the order and update its status
	var targetOrder *domain.Order
	for i, order := range orders {
		if order.ID == id && order.Status == domain.StatusPending {
			orders[i].Status = domain.StatusCompleted
			targetOrder = orders[i]
			break
		}
	}

	// If the order is not found or already completed
	if targetOrder == nil {
		return http.StatusNotFound, fmt.Errorf("order %s not found or already completed", id)
	}

	// Get menu items
	menuData, err := a.Repository.GetMenuItems()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	menuItems, err := a.Repository.UnmarshalJsonMenuItems(menuData)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Get inventory items
	inventoryData, err := a.Repository.GetInventoryItems()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	inventoryItems, err := a.Repository.UnmarshalInventoryItems(inventoryData)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Update inventory
	for _, orderItem := range targetOrder.Items {
		menuItem := findMenuItem(orderItem.ProductID, menuItems)
		if menuItem == nil {
			return http.StatusBadRequest, fmt.Errorf("menu item %s not found", orderItem.ProductID)
		}

		// decrement inventory
		for _, ingredient := range menuItem.Ingredients {
			if err := decrementInventory(ingredient, orderItem.Quantity, inventoryItems); err != nil {
				return http.StatusConflict, fmt.Errorf("insufficient ingredients for menu item %s", menuItem.ID)
			}
		}
	}

	// Save the updated inventory
	updatedInventoryData, err := a.Repository.MarshalInventoryItems(inventoryItems)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if err := a.Repository.SaveInventoryItems(updatedInventoryData); err != nil {
		return http.StatusInternalServerError, err
	}

	// Save the updated orders
	ordersJson, err := a.Repository.MarshalJsonOrders(orders)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if err := a.Repository.SaveOrders(ordersJson); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// Additional functions
func findMenuItem(id string, menuItems []*domain.MenuItem) *domain.MenuItem {
	for _, item := range menuItems {
		if item.ID == id {
			return item
		}
	}
	return nil
}

func checkIngredientsAvailability(quantity int, ingredients []domain.MenuItemIngredient, inventoryItems []*domain.InventoryItem) bool {
	for _, ingredient := range ingredients {
		for _, inventoryItem := range inventoryItems {
			if ingredient.IngredientID == inventoryItem.IngredientID && inventoryItem.Quantity < float64(quantity) {
				return false
			}
		}
	}
	return true
}

func hasIngredient(menuItemIngredients []domain.MenuItemIngredient, inventoryItems []*domain.InventoryItem) bool {
	for _, ingredient := range menuItemIngredients {
		for _, inventoryItem := range inventoryItems {
			if inventoryItem.IngredientID == ingredient.IngredientID {
				return true
			}
		}
	}
	return false
}

func decrementInventory(ingredient domain.MenuItemIngredient, orderQuantity int, inventoryItems []*domain.InventoryItem) error {
	requiredQuantity := ingredient.Quantity * float64(orderQuantity)
	for _, inventoryItem := range inventoryItems {
		if inventoryItem.IngredientID == ingredient.IngredientID {
			if inventoryItem.Quantity < requiredQuantity {
				return fmt.Errorf("insufficient quantity for ingredient %s", ingredient.IngredientID)
			}
			inventoryItem.Quantity -= requiredQuantity
			return nil
		}
	}
	return fmt.Errorf("ingredient %s not found in inventory", ingredient.IngredientID)
}

func generateOrderID() string {
	timestamp := time.Now().UnixNano()
	randomNumber := rand.Intn(10000)
	return strings.ReplaceAll(fmt.Sprintf("ORD-%d-%04d", timestamp, randomNumber), "/", "")
}

func CheckOrderFields(order *domain.Order) error {
	if order.ID == "" {
		return errors.New("order ID is required")
	}

	if order.CustomerName == "" {
		return errors.New("customer name is required")
	}

	if order.Status != domain.StatusPending {
		return fmt.Errorf("invalid order status: %s", order.Status)
	}

	if len(order.Items) == 0 {
		return errors.New("order must contain at least one item")
	}

	for _, item := range order.Items {
		if item.ProductID == "" {
			return errors.New("product ID is required for each item")
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("quantity for product %s must be greater than zero", item.ProductID)
		}
	}

	return nil
}
