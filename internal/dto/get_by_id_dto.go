package dto

type GetByIDDto struct {
	OriginalUrl string `json:"original-url"`
	IsDeleted   bool   `json:"is_deleted"`
}
