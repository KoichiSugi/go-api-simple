package service

import "net/http"

type ServiceImpl interface {
	GetAllEmployees(w http.ResponseWriter, r *http.Request)
}
