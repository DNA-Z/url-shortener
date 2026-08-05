package handler

import (
	"encoding/json"
	"net"
	"net/http"
)

// StatsGet обрабатывает запрос на получение статистики сервиса
//
// Пример запроса:
//
//	GET /api/internal/stats HTTP/1.1
//
// Пример ответа:
//
//	{
//	  "urls": 123,
//	  "users": 45
//	}
//
// Возможные статусы:
// - 200 OK - успешный ответ
// - 403 Forbidden - IP не в доверенной подсети
// - 400 Bad Request - неверный метод запроса
// - 500 Internal Server Error - внутренняя ошибка сервера
func (h *URLHandler) StatsGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}

	if h.trustedSubnet == "" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP == "" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	clientIP := net.ParseIP(realIP)
	if clientIP == nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_, trustedNet, err := net.ParseCIDR(h.trustedSubnet)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if !trustedNet.Contains(clientIP) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	stats, err := h.urlService.GetStats()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]int{
		"urls":  stats.URLs,
		"users": stats.Users,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
