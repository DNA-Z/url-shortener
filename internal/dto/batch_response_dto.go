package dto

type BatchResponseDto struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}
