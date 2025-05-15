package handler

import (
	"io"
	"net/http"
)

// OrderHandler обрабатывает запросы для работы с заказами (получение всех, добавление нового)
func (h *CustomHandler) OrderHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAllOrders(w, r)
	case http.MethodPost:
		h.AddOrder(w, r)
	default:
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// OrderByIDHandler обрабатывает запросы для работы с заказами по ID (получение, обновление, удаление)
func (h *CustomHandler) OrderByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetOrderByID(w, r)
	case http.MethodPut:
		h.UpdateOrderByID(w, r)
	case http.MethodDelete:
		h.DeleteOrderByID(w, r)
	default:
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// CloseOrderHandler обрабатывает запрос для закрытия заказа по ID
func (h *CustomHandler) CloseOrderHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CloseOrderByID(w, r)
	default:
		h.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// GetAllOrders получает все заказы
func (h *CustomHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	// Получаем все заказы через сервис
	data, status, err := h.Service.GetAllOrders()
	if err != nil {
		h.respondWithError(w, status, err.Error())
		return
	}
	h.respondWithJSON(w, status, data)
}

// AddOrder добавляет новый заказ
func (h *CustomHandler) AddOrder(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	// Чтение данных из тела запроса
	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request")
	}

	status, err := h.Service.AddOrder(data)
	if err != nil {
		h.respondWithError(w, status, err.Error())
		return
	}
	h.respondWithJSON(w, status, nil)
}

// GetOrderByID получает заказ по его ID
func (h *CustomHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	id := r.PathValue("id")

	// Получаем заказ через сервис
	data, status, err := h.Service.GetOrderByID(id)
	if err != nil {
		h.respondWithError(w, status, err.Error())
		return
	}
	h.respondWithJSON(w, status, data)
}

// UpdateOrderByID обновляет заказ по его ID
func (h *CustomHandler) UpdateOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	id := r.PathValue("id")

	// Чтение данных из тела запроса
	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request")
	}

	// Обновляем заказ через сервис
	status, err := h.Service.UpdateOrderByID(id, data)
	if err != nil {
		h.respondWithError(w, status, err.Error())
		return
	}
	h.respondWithJSON(w, status, nil)
}

// DeleteOrderByID удаляет заказ по его ID
func (h *CustomHandler) DeleteOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	id := r.PathValue("id")

	// Удаляем заказ через сервис
	status, err := h.Service.DeleteOrderByID(id)
	if err != nil {
		h.respondWithError(w, status, err.Error())
		return
	}

	h.respondWithJSON(w, status, nil)
}

// CloseOrderByID закрывает заказ по его ID
func (h *CustomHandler) CloseOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		h.respondWithError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return
	}

	id := r.PathValue("id")

	// Закрываем заказ через сервис
	status, err := h.Service.CloseOrderByID(id)
	if err != nil {
		h.respondWithError(w, status, err.Error())
		return
	}

	h.respondWithJSON(w, status, nil)
}

// respondWithJSON отправляет ответ в формате JSON
func (h *CustomHandler) respondWithJSON(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
	}
}
