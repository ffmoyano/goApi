package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type Payload struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
}

func base64encode(word []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(word), "=")
}

func base64decode(word string) string {
	length := len(word) % 4
	if length%4 > 0 {
		word += strings.Repeat("=", 4-length)
	}
	decoded, err := base64.URLEncoding.DecodeString(word)
	if err != nil {
		log.Fatalf("Decoding Error %s", err)
	}
	return string(decoded)
}

func Hash(src string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(src))
	return base64encode(h.Sum(nil))
}

func isValidHash(value string, hash string, secret string) bool {
	return hash == Hash(value, secret)
}

func AssembleJWT(payload Payload, secret string) (string, error) {
	jsonHeader, err := json.Marshal(Header{"HS256", "JWT"})
	if err != nil {
		return fmt.Sprintf("Error while encoding to json the gopherjwt header: %s", err), err
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Sprintf("Error while encoding to json the gopherjwt payload: %s", err), err
	}
	signatureValue := base64encode(jsonHeader) + "." + base64encode(jsonPayload)
	return signatureValue + "." + Hash(signatureValue, secret), nil
}

func Decode(jwt string, secret string) (interface{}, error) {
	token := strings.Split(jwt, ".")
	// check if the jwt token contains
	// header, payload and token
	if len(token) != 3 {
		splitErr := errors.New(" Invalid token: token should contain header, payload and secret")
		return nil, splitErr
	}
	// base64decode payload
	decodedPayload := base64decode(token[1])
	payload := Payload{}
	// parses payload from string to a struct
	ParseErr := json.Unmarshal([]byte(decodedPayload), &payload)
	if ParseErr != nil {
		return nil, fmt.Errorf(" Invalid payload: %s", ParseErr.Error())
	}
	// checks if the token has expired.
	if payload.Exp != 0 && time.Now().Unix() > payload.Exp {
		return nil, errors.New(" Expired token: token has expired")
	}
	signatureValue := token[0] + "." + token[1]
	// verifies if the header and signature is exactly whats in
	// the signature
	if isValidHash(signatureValue, token[2], secret) == false {
		return nil, errors.New(" Invalid token")
	}
	return payload, nil
}
