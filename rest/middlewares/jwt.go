package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ShajibEwuCse19/gopay-core/utils"
)

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if !checkAccessToken(w, r) {
				return 
			}
			next.ServeHTTP(w, r)
		},
	)
}

func checkAccessToken(w http.ResponseWriter, r *http.Request) bool {
	header := r.Header.Get("Authorization")
	if header == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	headerArray := strings.Split(header, " ")
	if len(headerArray) != 2 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	accessToken := headerArray[1] //header.payload.signature

	if utils.VerifyJwt(accessToken) {
		fmt.Println("Token is valid")
		fmt.Println("Token: ", accessToken)
		return true
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
} 