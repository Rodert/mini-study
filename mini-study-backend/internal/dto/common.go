package dto

// APIResponse is a shared response wrapper used by swagger examples.
type APIResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}
