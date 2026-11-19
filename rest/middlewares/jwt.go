package middleware

import (
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
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return false
	}

	headerArray := strings.Split(header, " ")
	if len(headerArray) != 2 {
		utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return false
	}

	accessToken := headerArray[1] //header.payload.signature

	if utils.VerifyJwt(accessToken) {
		return true
	}
	utils.SendError(w, http.StatusUnauthorized, "Unauthorized")
	return false
}
