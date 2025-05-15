package handler

import (
	"encoding/json"
	"net/http"

	"hot-coffee/internal/domain"
)

// Метод для отправки ответа с ошибкой
func (h *CustomHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	// Устанавливаем заголовок ответа в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// Устанавливаем код состояния HTTP
	w.WriteHeader(code)

	// Создаем структуру ошибки
	errMsg := domain.Error{
		Code:    code,
		Message: message,
	}

	// Кодируем и отправляем сообщение об ошибке в формате JSON
	if err := json.NewEncoder(w).Encode(errMsg); err != nil {
		// Логируем ошибку, если не удалось закодировать сообщение об ошибке
		h.LoggerERROR.Println("Failed to encode error message:", err)
	}
}
