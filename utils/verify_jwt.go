package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"strings"
	"github.com/ShajibEwuCse19/gopay-core/config"
)

func VerifyJwt(accessToken string) bool {
	accessToken = strings.Trim(accessToken, "\" \t\r\n")

	parts := strings.Split(accessToken, ".")
	if len(parts) != 3 {
		return false
	}

	header := parts[0]
	payload := parts[1]
	signature := parts[2]

	unsignedMessageBase64 := header + "." + payload

	hmacSha256 := hmac.New(sha256.New, []byte(config.GetConfig().JWTSecretKey))
	hmacSha256.Write([]byte(unsignedMessageBase64))
	newSignature := hmacSha256.Sum(nil)
	newSignatureBase64 := Base64ByteArrayEncode(newSignature)

	return hmac.Equal([]byte(newSignatureBase64), []byte(signature))

}