package handler

import (
	"encoding/json"
	"net/http"
)

// Обработчик запроса на получение общей суммы продаж
func (h *CustomHandler) GetTotalSalesHandler(w http.ResponseWriter, r *http.Request) {
	h.LoggerINFO.Println("GetTotalSalesHandler - Received request to get total sales.")

	// Получаем общую сумму продаж через сервис
	totalSales, err := h.Service.GetTotalSales()
	if err != nil {
		// Обрабатываем ошибку и возвращаем клиенту ошибку сервера
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		h.LoggerERROR.Printf("Error getting total sales: %v", err)
		return
	}

	// Формируем ответ в формате JSON
	response := map[string]float64{"total_sales": totalSales}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Обрабатываем ошибку при кодировании ответа
		h.LoggerERROR.Printf("GetTotalSalesHandler - Error encoding response: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
	h.LoggerINFO.Println("GetTotalSalesHandler - Successfully responded with total sales.")
}

// Обработчик запроса на получение популярных товаров
func (h *CustomHandler) GetPopularItemsHandler(w http.ResponseWriter, r *http.Request) {
	h.LoggerINFO.Println("GetPopularItemsHandler - Received request to get popular items.")
	// Получаем популярные товары через сервис
	popularItems, err := h.Service.GetPopularItems()
	if err != nil {
		// Обрабатываем ошибку и возвращаем клиенту ошибку сервера
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		h.LoggerERROR.Printf("Error getting popular items: %v", err)
		return
	}

	// Формируем ответ в формате JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(popularItems)
	if err != nil {
		// Обрабатываем ошибку при кодировании ответа
		h.LoggerERROR.Printf("GetPopularItemsHandler - Error encoding popular items response: %v", err)
		h.respondWithError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
	h.LoggerINFO.Println("GetPopularItemsHandler - Successfully responded with popular items.")
}
