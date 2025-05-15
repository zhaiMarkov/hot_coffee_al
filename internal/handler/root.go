package handler

import "net/http"

func (h *CustomHandler) RootHandler(w http.ResponseWriter, r *http.Request) {
	h.respondWithError(w, http.StatusNotFound, "Not found")
}
