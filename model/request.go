package model

type (
	NewHeapRequest struct {
		Host   string `json:"host" validate:"required"`
		Port   string `json:"port" validate:"required"`
		Binary string `json:"binary" validate:"required"`
		Git    string `json:"git"`
	}
)
