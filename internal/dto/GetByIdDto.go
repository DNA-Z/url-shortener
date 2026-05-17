package dto

type GetByIdDto struct {
	OriginalUrl string `json:"original-url"`
	IsDeleted   bool   `json:"is_deleted"`
}
