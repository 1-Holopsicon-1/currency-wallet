package models

type Response struct {
	Status  int
	Message interface{} `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
