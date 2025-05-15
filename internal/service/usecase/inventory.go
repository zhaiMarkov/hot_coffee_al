package usecase

import (
	"fmt"
	"log"
	"net/http"

	"hot-coffee/internal/domain"
)

func (a *Application) AddInventoryItem(data []byte) (int, error) {
	item, err := a.Repository.UnmarshalJsonInventory(data)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid inventory item data")
	}

	if err := validateInventoryItem(item); err != nil {
		return http.StatusBadRequest, err
	}

	allData, err := a.Repository.GetInventoryItems()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	inventoryItems, err := a.Repository.UnmarshalInventoryItems(allData)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for _, items := range inventoryItems {
		if items.IngredientID == item.IngredientID {
			return http.StatusConflict, fmt.Errorf("inventory item with Ingredient ID %s already exists", item.IngredientID)
		}

		if items.Name == item.Name {
			return http.StatusConflict, fmt.Errorf("inventory item with name %s already exists", item.Name)
		}
	}

	inventoryItems = append(inventoryItems, item)
	log.Printf("Validated and ready for storage: %+v", item)

	updatedData, err := a.Repository.MarshalInventoryItems(inventoryItems)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := a.Repository.SaveInventoryItems(updatedData); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (a *Application) GetAllInventoryItems() ([]byte, int, error) {
	data, err := a.Repository.GetInventoryItems()
	if err != nil {
		return nil, http.StatusNoContent, err
	}

	if len(data) <= 2 {
		return nil, http.StatusNotFound, fmt.Errorf("no inventory items found")
	}
	return data, http.StatusOK, nil
}

func (a *Application) GetInventoryItemByID(id string) ([]byte, int, error) {
	data, err := a.Repository.GetInventoryItems()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	inventryItems, err := a.Repository.UnmarshalInventoryItems(data)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	for _, item := range inventryItems {
		if item.IngredientID == id {
			inventoryItemJson, err := a.Repository.MarshalJsonInventory(item)
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}
			return inventoryItemJson, http.StatusOK, nil

		}
	}
	return nil, http.StatusNotFound, fmt.Errorf("inventory item with Ingredient ID %s not found", id)
}

func (a *Application) UpdateInventoryItemByID(id string, data []byte) (int, error) {
	inventory, err := a.Repository.UnmarshalJsonInventory(data)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid inventory item data")
	}

	if err = validateInventoryItem(inventory); err != nil {
		return http.StatusBadRequest, err
	}

	allData, err := a.Repository.GetInventoryItems()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	inventoryItems, err := a.Repository.UnmarshalInventoryItems(allData)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for i, item := range inventoryItems {
		if item.IngredientID == id {
			inventoryItems[i] = inventory
			break
		}
	}

	inventoryItemsJson, err := a.Repository.MarshalInventoryItems(inventoryItems)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = a.Repository.SaveInventoryItems(inventoryItemsJson)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (a *Application) DeleteInventoryItemByID(id string) (int, error) {
	allData, err := a.Repository.GetInventoryItems()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if len(allData) <= 2 {
		return http.StatusBadRequest, fmt.Errorf("Data is empty")
	}

	inventoryItems, err := a.Repository.UnmarshalInventoryItems(allData)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for i, item := range inventoryItems {
		if item.IngredientID == id {
			inventoryItems = append(inventoryItems[:i], inventoryItems[i+1:]...)
			break
		}
	}

	inventoryItemsJson, err := a.Repository.MarshalInventoryItems(inventoryItems)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = a.Repository.SaveInventoryItems(inventoryItemsJson)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusNoContent, nil
}

func validateInventoryItem(item *domain.InventoryItem) error {
	if item.IngredientID == "" {
		return fmt.Errorf("ingredient ID is required")
	}
	if item.Name == "" {
		return fmt.Errorf("name is required")
	}
	if item.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}
	if item.Unit == "" {
		return fmt.Errorf("unit is required")
	}
	return nil
}
