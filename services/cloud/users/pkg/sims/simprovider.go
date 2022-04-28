package sims

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/ukama/ukama/services/cloud/users/pkg/db"
)

type SimProvider interface {
	GetICCIDWithCode(simCode string) (string, error)
	GetICCIDFromPool() (string, error)
	GetSimToken(iccid string) (string, error)
}

type simProvider struct {
	key     string
	simPool db.SimPoolRepo
}

type SimToken struct {
	ICCID string `json:"iccid"`
}

func NewSimProvider(key string, simPool db.SimPoolRepo) *simProvider {
	return &simProvider{key: key, simPool: simPool}
}

func (i simProvider) GetICCIDWithCode(simCode string) (string, error) {

	str, err := decrypt(simCode, i.key)
	if err != nil {
		return "", err
	}

	var simToken SimToken
	err = json.Unmarshal([]byte(str), &simToken)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal sim token")
	}

	return simToken.ICCID, nil
}

func (i simProvider) GetICCIDFromPool() (string, error) {
	iccid, err := i.simPool.Pop()
	if err != nil {
		return "", err
	}
	return iccid, nil
}

func (i simProvider) GetSimToken(iccid string) (string, error) {
	token := SimToken{
		ICCID: iccid,
	}

	tokenJson, err := json.Marshal(token)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal sim token")
	}

	return encrypt(string(tokenJson), i.key)
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
