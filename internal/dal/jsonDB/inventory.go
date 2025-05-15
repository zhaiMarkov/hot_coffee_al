package jsondb

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"hot-coffee/internal/config"
	"hot-coffee/internal/domain"
)

// Получение данных об инвентаре из файла inventory.json
func (j *JsonDB) GetInventoryItems() ([]byte, error) {
	path := filepath.Join(config.Dir, "inventory.json")
	file, err := os.OpenFile(path, os.O_RDONLY, 0o666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Десериализация данных об инвентаре из JSON
func (j *JsonDB) UnmarshalInventoryItems(data []byte) ([]*domain.InventoryItem, error) {
	var inventoryItem []*domain.InventoryItem

	err := json.Unmarshal(data, &inventoryItem)
	if err != nil {
		return nil, err
	}
	return inventoryItem, nil
}

// Сериализация данных об инвентаре в JSON
func (j *JsonDB) MarshalInventoryItems(inventoryItems []*domain.InventoryItem) ([]byte, error) {
	return json.Marshal(inventoryItems)
}

// Десериализация одного элемента инвентаря из JSON
func (j *JsonDB) UnmarshalJsonInventory(data []byte) (*domain.InventoryItem, error) {
	var inventoryItem *domain.InventoryItem

	err := json.Unmarshal(data, &inventoryItem)
	if err != nil {
		return nil, err
	}
	return inventoryItem, nil
}

// Сериализация одного элемента инвентаря в JSON
func (j *JsonDB) MarshalJsonInventory(inventoryItem *domain.InventoryItem) ([]byte, error) {
	return json.Marshal(inventoryItem)
}

// Сохранение данных об инвентаре в файл inventory.json
func (j *JsonDB) SaveInventoryItems(data []byte) error {
	path := filepath.Join(config.Dir, "inventory.json")

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
