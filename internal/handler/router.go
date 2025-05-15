package handler

import "net/http"

func (h *CustomHandler) Routing() *http.ServeMux {
	router := http.NewServeMux()

	// Order
	router.HandleFunc("/order", h.OrderHandler)
	router.HandleFunc("/order/{id}", h.OrderByIDHandler)
	router.HandleFunc("/order/{id}/close", h.CloseOrderHandler)

	// Menu
	router.HandleFunc("/menu", h.MenuHandler)
	router.HandleFunc("/menu/{id}", h.MenuByIDHandler)

	// Inventory
	router.HandleFunc("/inventory", h.InventoryHandler)
	router.HandleFunc("/inventory/{id}", h.InventoryByIDHandler)

	// aggregation
	router.HandleFunc("/reports/total-sales", h.GetTotalSalesHandler)
	router.HandleFunc("/reports/popular-items", h.GetPopularItemsHandler)

	router.HandleFunc("/", h.RootHandler)
	return router
}
