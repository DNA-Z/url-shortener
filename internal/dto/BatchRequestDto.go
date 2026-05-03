package dto

type BatchRequestDto struct {
	ID          string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}
