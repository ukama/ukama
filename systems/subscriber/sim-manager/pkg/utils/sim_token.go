package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ukama/ukama/systems/common/errors"
)

func GetIccidFromToken(simToken string, key string) (string, error) {
	str, err := decrypt(simToken, key)
	if err != nil {
		return "", err
	}

	var iccidEnvelope struct {
		ICCID string `json:"iccid"`
	}

	err = json.Unmarshal([]byte(str), &iccidEnvelope)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal sim token")
	}

	return iccidEnvelope.ICCID, nil
}

func GenerateTokenFromIccid(iccid string, key string) (string, error) {
	iccidEnvelope := struct {
		ICCID string `json:"iccid"`
	}{
		ICCID: iccid,
	}

	tokenJson, err := json.Marshal(iccidEnvelope)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal sim token")
	}

	return encrypt(string(tokenJson), key)
}

func encrypt(plaintext string, key string) (string, error) {
	if len(key) != 32 {
		return "", fmt.Errorf("key must be 32 bytes")
	}

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	b := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(b), nil
}

func decrypt(base64Str string, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}
	c, err := aes.NewCipher([]byte(key)[:32])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	o, err := gcm.Open(nil, nonce, ciphertext, nil)
	return string(o), err
}
