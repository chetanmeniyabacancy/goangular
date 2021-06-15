package controllers

type CommonError struct {
	Status int64 `json:"status"`
    Message string `json:"message"`
}

type CommonSuccess struct {
	Status int64 `json:"status"`
    Message string `json:"message"`
}