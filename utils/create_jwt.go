package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type PayLoad struct {
	Sub         string `json:"sub"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email,omitempty"`
	FullName    string `json:"full_name"`
	KycStatus   string `json:"kyc_status"`
	Status      string `json:"status"`
	Iat         int64  `json:"iat"`
	Exp         int64  `json:"exp"`
}

func Base64ByteArrayEncode(data []byte) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

func CreateJwt(secret string, data PayLoad) (string, error) {
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	byteArrayHeader, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	byteArrayPayLoad, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	headerBase64 := Base64ByteArrayEncode(byteArrayHeader)
	payLoadBase64 := Base64ByteArrayEncode(byteArrayPayLoad)

	unsignedMessage := headerBase64 + "." + payLoadBase64

	hmacSha256 := hmac.New(sha256.New, []byte(secret))
	hmacSha256.Write([]byte(unsignedMessage))
	signature := hmacSha256.Sum(nil)

	signatureBase64 := Base64ByteArrayEncode(signature)

	jwt := unsignedMessage + "." + signatureBase64

	return jwt, nil
}
