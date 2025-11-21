package dto

// APIResponse is a shared response wrapper used by swagger examples.
type APIResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// Pagination describes paged list metadata.
type Pagination struct {
	Page     int   `json:"page" example:"1"`
	PageSize int   `json:"page_size" example:"20"`
	Total    int64 `json:"total" example:"100"`
}
