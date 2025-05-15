package usecase

import (
	"errors"
	"fmt"
	"net/http"

	"hot-coffee/internal/domain"
)

func (a *Application) AddMenu(data []byte) (int, error) {
	// Unmarshal the JSON menu
	menu, err := a.Repository.UnmarshalJsonMenu(data)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid menu data")
	}

	// Check if all fields are set
	if err = CheckMenuItemFields(menu); err != nil {
		return http.StatusBadRequest, err
	}

	// Get all menu items
	allData, err := a.Repository.GetMenuItems()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Unmarshal the JSON menu items
	menuItems, err := a.Repository.UnmarshalJsonMenuItems(allData)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Check if the menu item already exists
	for _, item := range menuItems {
		if item.ID == menu.ID {
			return http.StatusBadRequest, fmt.Errorf("menu item with ID %s already exists", menu.ID)
		}

		if item.Name == menu.Name {
			return http.StatusBadRequest, fmt.Errorf("menu item with name %s already exists", menu.Name)
		}
	}
	// Append the new menu item
	menuItems = append(menuItems, menu)

	// Marshal the JSON menu items
	menuItemsJson, err := a.Repository.MarshalJsonMenuItems(menuItems)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Save the menu items
	err = a.Repository.SaveMenuItems(menuItemsJson)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (a *Application) GetAllMenuItems() ([]byte, int, error) {
	data, err := a.Repository.GetMenuItems()
	if err != nil {
		return nil, http.StatusNoContent, err
	}

	if len(data) <= 2 {
		return nil, http.StatusNotFound, fmt.Errorf("no menu items found")
	}
	return data, http.StatusOK, nil
}

func (a *Application) GetMenuItemByID(id string) ([]byte, int, error) {
	// Get all menu items
	data, err := a.Repository.GetMenuItems()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Unmarshal the JSON menu items
	menuItems, err := a.Repository.UnmarshalJsonMenuItems(data)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Find the menu item by ID
	for _, item := range menuItems {
		if item.ID == id {
			// Marshal the menu item
			menuItemJson, err := a.Repository.MarshalJsonMenu(item)
			if err != nil {
				return nil, http.StatusInternalServerError, err
			}
			return menuItemJson, http.StatusOK, nil
		}
	}

	return nil, http.StatusNotFound, fmt.Errorf("menu item with ID %s not found", id)
}

func (a *Application) UpdateMenuItemByID(id string, data []byte) (int, error) {
	// Get the menu item by ID
	menu, err := a.Repository.UnmarshalJsonMenu(data)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Check if all fields are set
	if err = CheckMenuItemFields(menu); err != nil {
		return http.StatusBadRequest, err
	}

	// Get all menu items
	allData, err := a.Repository.GetMenuItems()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Unmarshal the JSON menu items
	menuItems, err := a.Repository.UnmarshalJsonMenuItems(allData)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Update the menu item
	for i, item := range menuItems {
		if item.ID == id {
			menuItems[i] = menu
			break
		}
	}

	// Marshal the JSON menu items
	menuItemsJson, err := a.Repository.MarshalJsonMenuItems(menuItems)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Save the menu items
	err = a.Repository.SaveMenuItems(menuItemsJson)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (a *Application) DeleteMenuItemByID(id string) (int, error) {
	// Get all menu items
	allData, err := a.Repository.GetMenuItems()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if len(allData) <= 2 {
		return http.StatusBadRequest, fmt.Errorf("Data is empty")
	}

	// Unmarshal the JSON menu items
	menuItems, err := a.Repository.UnmarshalJsonMenuItems(allData)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Find the menu item by ID
	for i, item := range menuItems {
		if item.ID == id {
			// Remove the menu item
			menuItems = append(menuItems[:i], menuItems[i+1:]...)
			break
		}
	}

	// Marshal the JSON menu items
	menuItemsJson, err := a.Repository.MarshalJsonMenuItems(menuItems)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Save the menu items
	err = a.Repository.SaveMenuItems(menuItemsJson)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusNoContent, nil
}

func CheckMenuItemFields(menuItem *domain.MenuItem) error {
	if menuItem.ID == "" {
		return errors.New("menu item ID is required")
	}

	if menuItem.Name == "" {
		return errors.New("menu item name is required")
	}

	if menuItem.Price <= 0 {
		return fmt.Errorf("menu item %s must have a positive price", menuItem.ID)
	}

	for _, ingredient := range menuItem.Ingredients {

		if ingredient.IngredientID == "" {
			return fmt.Errorf("ingredient ID is required for menu item %s", menuItem.ID)
		}

		if ingredient.Quantity <= 0 {
			return fmt.Errorf("ingredient %s in menu item %s must have a positive quantity", ingredient.IngredientID, menuItem.ID)
		}
	}

	return nil
}
