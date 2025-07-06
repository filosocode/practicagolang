package controllers

import (
	"encoding/json"
	"net/http"
)

type Saludo struct {
	Msg        string `json:"message"`
	StatusCode int    `json:"code"`
}

func GetInitRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	saludo := Saludo{
		Msg:        "API funcionando correctamente",
		StatusCode: http.StatusOK,
	}

	json.NewEncoder(w).Encode(saludo)
}
