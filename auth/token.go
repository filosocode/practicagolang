package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerarToken(correo string) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = correo
	claims["autorized"] = true
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	// CORREGIDO: usar método de firma compatible con clave secreta
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return jwtToken.SignedString([]byte(os.Getenv("API_SECRET")))
}

func ValidarToken(r *http.Request) (string, error) {
	jwtToken := ExtraerToken(r)

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// CORREGIDO: typo y verificación correcta del método
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metodo de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)

		// CORREGIDO: devolver el ID del claim
		if id, ok := claims["id"].(string); ok {
			return id, nil
		}
		return "", fmt.Errorf("token sin claim 'id' válido")
	}

	return "", fmt.Errorf("token inválido")
}

func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(b))
}

func ExtraerToken(r *http.Request) string {
	parametros := r.URL.Query()
	token := parametros.Get("token")
	if token != "" {
		return token
	}

	// CORREGIDO: typo en el header
	tokenString := r.Header.Get("Authorization")
	if len(strings.Split(tokenString, " ")) == 2 {
		return strings.Split(tokenString, " ")[1]
	}
	return ""
}
