/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

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

type Codec interface {
	GetIccidFromToken(string) (string, error)
	GenerateTokenFromIccid(string) (string, error)
}

type tokenCodec struct {
	Key string
}

func NewTokenCodec(key string) Codec {
	return &tokenCodec{
		Key: key,
	}
}

func (c *tokenCodec) GetIccidFromToken(simToken string) (string, error) {
	str, err := decrypt(simToken, c.Key)
	if err != nil {
		return "", err
	}

	var iccidEnvelope struct {
		ICCID string `json:"iccid,string"`
	}

	err = json.Unmarshal([]byte(str), &iccidEnvelope)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal sim token")
	}

	return iccidEnvelope.ICCID, nil
}

func (c *tokenCodec) GenerateTokenFromIccid(iccid string) (string, error) {
	iccidEnvelope := struct {
		ICCID string `json:"iccid,string"`
	}{
		ICCID: iccid,
	}

	tokenJson, err := json.Marshal(iccidEnvelope)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal sim token")
	}

	return encrypt(string(tokenJson), c.Key)
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
