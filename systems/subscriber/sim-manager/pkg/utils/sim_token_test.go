package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const textToEncode = "{ 'iccid': '8910300000003540855' }"
const testKey = "the-key-has-to-be-32-bytes-long!"

func Test_encrypt(t *testing.T) {

	r, err := encrypt(textToEncode, testKey)
	if !assert.NoError(t, err) {
		assert.FailNow(t, "encrypt failed")
	}
	res, err := decrypt(r, testKey)
	assert.NoError(t, err)
	assert.Equal(t, textToEncode, res)
}

func Test_encryptErrors(t *testing.T) {
	tests := []struct {
		key  string
		text string
	}{
		{key: "short testKey", text: textToEncode},
		{key: "", text: textToEncode},
	}
	for _, tt := range tests {
		_, err := encrypt(tt.text, tt.key)
		assert.Error(t, err)
	}
}

func Test_decryptErrors(t *testing.T) {
	encrText, err := encrypt(textToEncode, testKey)
	if err != nil {
		assert.FailNow(t, "encrypt failed", err)
	}

	tests := []struct {
		key  string
		text string
	}{
		{key: "short key", text: encrText},
		{key: "", text: textToEncode},
		{key: "", text: "Not encoded text"},
	}
	for _, tt := range tests {
		_, err := decrypt(tt.text, tt.key)
		assert.Error(t, err)
	}
}
