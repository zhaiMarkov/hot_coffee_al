package jsondb

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"hot-coffee/internal/config"
	"hot-coffee/internal/domain"
)

// Получение меню из файла menu.json
func (j *JsonDB) GetMenuItems() ([]byte, error) {
	path := filepath.Join(config.Dir, "menu.json")

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

// Сохранение меню в файл menu.json
func (j *JsonDB) SaveMenuItems(data []byte) error {
	path := filepath.Join(config.Dir, "menu.json")

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

// Десериализация одного элемента меню из JSON
func (j *JsonDB) UnmarshalJsonMenu(data []byte) (*domain.MenuItem, error) {
	var menuItem domain.MenuItem
	err := json.Unmarshal(data, &menuItem)
	if err != nil {
		return nil, err
	}

	return &menuItem, nil
}

// Сериализация одного элемента меню в JSON
func (j *JsonDB) MarshalJsonMenu(menuItem *domain.MenuItem) ([]byte, error) {
	menuItemJson, err := json.Marshal(menuItem)
	if err != nil {
		return nil, err
	}

	return menuItemJson, nil
}

// Десериализация массива элементов меню из JSON
func (j *JsonDB) UnmarshalJsonMenuItems(data []byte) ([]*domain.MenuItem, error) {
	var menuItems []*domain.MenuItem
	err := json.Unmarshal(data, &menuItems)
	if err != nil {
		return nil, err
	}

	return menuItems, nil
}

// Сериализация массива элементов меню в JSON
func (j *JsonDB) MarshalJsonMenuItems(menuItems []*domain.MenuItem) ([]byte, error) {
	menuItemsJson, err := json.Marshal(menuItems)
	if err != nil {
		return nil, err
	}

	return menuItemsJson, nil
}
