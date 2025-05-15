package handler

import (
	"io"
	"net/http"
)

// Обработчик для работы с инвентарем
func (h *CustomHandler) InventoryHandler(w http.ResponseWriter, r *http.Request) {
	h.LoggerINFO.Printf("InventoryHandler - %s request received", r.Method)

	switch r.Method {
	case http.MethodGet:
		h.getAllInventory(w, r)
	case http.MethodPost:
		h.addInventory(w, r)
	default:
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		h.LoggerERROR.Printf("InventoryHandler - Method %s not allowed", r.Method)
	}
}

// Обработчик для работы с инвентарем по ID
func (h *CustomHandler) InventoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	h.LoggerINFO.Printf("InventoryByIDHandler - %s request received", r.Method)

	switch r.Method {
	case http.MethodGet:
		h.getInventoryByID(w, r)
	case http.MethodPut:
		h.updateInventoryByID(w, r)
	case http.MethodDelete:
		h.deleteInventoryByID(w, r)
	default:
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed") // JSON PASE
		h.LoggerERROR.Printf("InventoryByIDHandler - Method %s not allowed", r.Method)
	}
}

// Получение всего инвентаря
func (h *CustomHandler) getAllInventory(w http.ResponseWriter, r *http.Request) {
	// Проверяем заголовок Content-Type
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	h.LoggerINFO.Println("getAllInventory - Fetching all inventory items")

	// Получаем все элементы через сервис
	data, status, err := h.Service.GetAllInventoryItems()
	if err != nil {
		h.LoggerERROR.Println(err)
		h.respondWithError(w, status, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		h.LoggerERROR.Println(err)
	}
}

// Добавление нового элемента в инвентарь
func (h *CustomHandler) addInventory(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	h.LoggerINFO.Println("addInventory - Adding new inventory item")

	defer r.Body.Close()
	// Читаем тело запроса
	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.LoggerERROR.Println("Error reading request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Добавляем элемент через сервис
	status, err := h.Service.AddInventoryItem(data)
	if err != nil {
		h.LoggerERROR.Println("Service error:", err)
		h.respondWithError(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.LoggerINFO.Println("addInventory - Inventory item added successfully")
}

// Получение элемента инвентаря по ID
func (h *CustomHandler) getInventoryByID(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	id := r.PathValue("id")
	h.LoggerINFO.Printf("getInventoryByID - Fetching inventory item with ID: %s", id)

	// Получаем элемент по ID через сервис
	data, status, err := h.Service.GetInventoryItemByID(id)
	if err != nil {
		h.LoggerERROR.Println(err)
		h.respondWithError(w, status, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		h.LoggerERROR.Println(err)
		h.respondWithError(w, http.StatusInternalServerError, "An error occurred while processing the request")
		return
	}
}

// Обновление элемента инвентаря по ID
func (h *CustomHandler) updateInventoryByID(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	id := r.PathValue("id")
	h.LoggerINFO.Printf("updateInventoryByID - Updating inventory item with ID: %s", id)

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.LoggerERROR.Println("Error reading request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Обновляем элемент через сервис
	if status, err := h.Service.UpdateInventoryItemByID(id, body); err != nil {
		h.LoggerERROR.Println(err)
		h.respondWithError(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Inventory item updated successfully"))
	h.LoggerINFO.Printf("updateInventoryByID - Inventory item with ID %s updated successfully", id)
}

// Удаление элемента инвентаря по ID
func (h *CustomHandler) deleteInventoryByID(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	id := r.PathValue("id")
	h.LoggerINFO.Printf("deleteInventoryByID - Deleting inventory item with ID: %s", id)

	// Удаляем элемент через сервис
	if status, err := h.Service.DeleteInventoryItemByID(id); err != nil {
		h.LoggerERROR.Println(err)
		h.respondWithError(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Inventory item deleted successfully"))
	h.LoggerINFO.Printf("deleteInventoryByID - Inventory item with ID %s deleted successfully", id)
}
