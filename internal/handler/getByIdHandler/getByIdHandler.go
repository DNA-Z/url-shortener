package getbyidhandler

import "net/http"

type GetByIdHandler struct {
	Message string `json:"message"`
}

func GetByIdGet(res http.ResponseWriter, req *http.Request) {
	return
}
