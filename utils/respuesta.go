package utils

type Respuesta struct {
	Msg        string `json:"message"`
	Data       string `json:"data"`
	StatusCode int    `json:"code"`
}
