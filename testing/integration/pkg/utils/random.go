package utils

import (
	b64 "encoding/base64"
	"strings"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/goombaio/namegenerator"
)

func RandomPort() int {
	n, _ := faker.RandomInt(1024, 65336)
	return n[0]
}

func RandomInt(m int) int {
	n, _ := faker.RandomInt(0, m)
	return n[0]
}

func RandomIntInRange(min int, max int) int {
	n, _ := faker.RandomInt(min, max, 1)
	return n[0]
}

func RandomName() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	name := nameGenerator.Generate()

	return name
}

func RandomBytes(n int) []byte {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[RandomInt(len(letterBytes)-1)]
	}

	return b
}

func RandomIPv4() string {
	return faker.IPv4()
}

func IPv4CIDRToStringNotation(i string) string {
	s := strings.Split(i, "/")
	return s[0]

}

func RandomString(n int) string {
	b := RandomBytes(n)
	return string(b)
}

func RandomBase64String(n int) string {

	return b64.StdEncoding.EncodeToString(RandomBytes(n))
}
