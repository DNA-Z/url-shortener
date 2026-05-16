package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func (h *URLHandler) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only Delete requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var shortURLs []string
	if err := json.Unmarshal(body, &shortURLs); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if len(shortURLs) == 0 {
		http.Error(w, "URL IDs list is empty", http.StatusBadRequest)
		return
	}

	go func(usrID string, urls []string) {
		if err := h.urlService.DeleteUserURLs(userID, shortURLs); err != nil {
			log.Printf("Failed to delete URLs: %v", err)
		}
	}(userID, shortURLs)

	w.WriteHeader(http.StatusAccepted)
}
