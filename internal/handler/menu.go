package handler

import (
	"io"
	"net/http"
)

// MenuHandler обрабатывает запросы для работы с меню (создание, получение всех элементов)
func (h *CustomHandler) MenuHandler(w http.ResponseWriter, r *http.Request) {
	h.LoggerINFO.Printf("MenuHandler - %s request received", r.Method)

	switch r.Method {
	case http.MethodGet:
		h.getAllMenu(w, r)
	case http.MethodPost:
		h.addMenu(w, r)
	default:
		h.LoggerERROR.Printf("MenuHandler - Method %s not allowed", r.Method)
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// MenuByIDHandler обрабатывает запросы для работы с меню по ID (получение, обновление, удаление)
func (h *CustomHandler) MenuByIDHandler(w http.ResponseWriter, r *http.Request) {
	h.LoggerINFO.Printf("MenuByIDHandler - %s request received", r.Method)

	switch r.Method {
	case http.MethodGet:
		h.getMenuByID(w, r)
	case http.MethodPut:
		h.updateMenuByID(w, r)
	case http.MethodDelete:
		h.deleteMenuByID(w, r)
	default:
		h.LoggerERROR.Printf("MenuByIDHandler - Method %s not allowed", r.Method)
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// getAllMenu получает все элементы меню из базы данных
func (h *CustomHandler) getAllMenu(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	h.LoggerINFO.Println("getAllMenu - Fetching all menu items")

	// Вызов сервиса для получения всех элементов меню
	data, status, err := h.Service.GetAllMenuItems()
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

// addMenu добавляет новый элемент в меню
func (h *CustomHandler) addMenu(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	h.LoggerINFO.Println("addMenu - Adding new menu item")

	defer r.Body.Close()

	// Чтение данных из тела запроса
	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.LoggerERROR.Println("Error reading request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Вызов сервиса для добавления нового элемента
	status, err := h.Service.AddMenu(data)
	if err != nil {
		h.LoggerERROR.Println("Service error:", err)
		h.respondWithError(w, status, err.Error())
		return
	}

	// Ответ с кодом 201 Created при успешном добавлении
	w.WriteHeader(http.StatusCreated)
	h.LoggerINFO.Println("addMenu - Menu item added successfully")
}

// getMenuByID получает элемент меню по его ID
func (h *CustomHandler) getMenuByID(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	id := r.PathValue("id")
	h.LoggerINFO.Printf("getMenuByID - Fetching menu item with ID: %s", id)

	// Вызов сервиса для получения элемента по ID
	data, status, err := h.Service.GetMenuItemByID(id)
	if err != nil {
		h.LoggerERROR.Println(err)
		h.respondWithError(w, status, err.Error())
		return
	}

	// Отправляем данные элемента в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		h.LoggerERROR.Println(err)
		h.respondWithError(w, http.StatusInternalServerError, "An error occurred while processing the request")
		return
	}
}

// updateMenuByID обновляет элемент меню по его ID
func (h *CustomHandler) updateMenuByID(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	id := r.PathValue("id")
	h.LoggerINFO.Printf("updateMenuByID - Updating menu item with ID: %s", id)

	defer r.Body.Close()

	// Чтение данных из тела запроса
	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.LoggerERROR.Println("Error reading request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Вызов сервиса для обновления элемента по ID
	status, err := h.Service.UpdateMenuItemByID(id, data)
	if err != nil {
		h.LoggerERROR.Println(err)
		h.respondWithError(w, status, err.Error())
		return
	}

	// Ответ с кодом 200 OK при успешном обновлении
	w.WriteHeader(status)
	h.LoggerINFO.Printf("updateMenuByID - Menu item with ID %s updated successfully", id)
}

// deleteMenuByID удаляет элемент меню по его ID
func (h *CustomHandler) deleteMenuByID(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	id := r.PathValue("id")
	h.LoggerINFO.Printf("deleteMenuByID - Deleting menu item with ID: %s", id)

	// Вызов сервиса для удаления элемента по ID
	status, err := h.Service.DeleteMenuItemByID(id)
	if err != nil {
		h.LoggerERROR.Println(err)
		h.respondWithError(w, status, err.Error())
		return
	}

	// Ответ с кодом 200 OK при успешном удалении
	w.WriteHeader(status)
	h.LoggerINFO.Printf("deleteMenuByID - Menu item with ID %s deleted successfully", id)
}
