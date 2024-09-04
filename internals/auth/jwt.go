package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Payload struct {
	Sub  int32  `json:"sub"`
	Name string `json:"name"`
	Exp  int64  `json:"exp"`
	Iat  int64  `json:"iat"`
}

func GenerateJWTToken(id int32, name string) (string, error) {
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}
	payload := Payload{
		Sub:  id,
		Name: name,
		Exp:  time.Now().Add(24 * time.Hour).Unix(),
		Iat:  time.Now().Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

	signingInput := headerEncoded + "." + payloadEncoded
	signature, err := generateSignature(signingInput, os.Getenv("JWT_SECRET"))
	if err != nil {
		return "", err
	}

	token := signingInput + "." + signature
	return token, nil
}

func generateSignature(signingInput, secretKey string) (string, error) {
	h := hmac.New(sha256.New, []byte(secretKey))
	_, err := h.Write([]byte(signingInput))
	if err != nil {
		return "", err
	}
	signature := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(signature), nil
}

func VerifyToken(token string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("invalid token format: token should have 3 parts")
	}

	headerEncoded := parts[0]
	payloadEncoded := parts[1]
	providedSignature := parts[2]

	signingInput := headerEncoded + "." + payloadEncoded
	expectedSignature, err := generateSignature(signingInput, os.Getenv("JWT_SECRET"))
	if err != nil {
		return err
	}

	if providedSignature != expectedSignature {
		return errors.New("invalid token signature: token verification failed")
	}

	return CheckTokenExpiration(payloadEncoded)
}

func CheckTokenExpiration(payloadEncoded string) error {
	var payload Payload
	decodedPayload, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return err
	}
	err = json.Unmarshal(decodedPayload, &payload)
	if err != nil {
		return err
	}

	if time.Now().Unix() > payload.Exp {
		return errors.New("token has expired")
	}

	return nil
}

func ExtractSubFromToken(token string) (int32, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, errors.New("invalid token format: token should have 3 parts")
	}

	payloadEncoded := parts[1]
	var payload Payload
	decodedPayload, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(decodedPayload, &payload)
	if err != nil {
		return 0, err
	}

	if payload.Sub == 0 {
		return 0, errors.New("token does not contain a subject claim")
	}

	if err != nil {
		return 0, errors.New("failed to parse subject as int32: " + err.Error())
	}

	return payload.Sub, nil
}
