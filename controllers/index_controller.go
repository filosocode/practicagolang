package controllers

import (
	"encoding/json"
	"net/http"
)

// Saludo define la estructura de la respuesta JSON
type Saludo struct {
	Msg        string `json:"message"`
	StatusCode int    `json:"code"`
}

// GetInitRoute maneja la ruta GET /api y responde con un mensaje de bienvenida
func GetInitRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	saludo := Saludo{
		Msg:        "API funcionando correctamente",
		StatusCode: http.StatusOK,
	}

	json.NewEncoder(w).Encode(saludo)
}
