package dto

type GetByIdDto struct {
	ShortURL  string `json:"short_url"`
	IsDeleted bool   `json:"is_deleted"`
}
