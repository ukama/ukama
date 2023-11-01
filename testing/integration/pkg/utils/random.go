/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	b64 "encoding/base64"
	"math/rand"
	"strings"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/goombaio/namegenerator"
)

var node_type = []string{"hnode", "tnode", "anode"}
var country_code = []string{"us", "uk", "eu", "pk", "rc"}

func RandomPort() int {
	n, _ := faker.RandomInt(1024, 65336)
	return n[0]
}

func RandomInt(m int) int {
	n, _ := faker.RandomInt(0, m, 1)
	return n[0]
}

func RandomGetNodeId(typ string) string {
	nIndex := rand.Int() % len(node_type)
	cIndex := rand.Int() % len(country_code)
	if typ == "" {
		typ = node_type[nIndex]
	}

	return strings.ToLower(country_code[cIndex] + "-" + RandomString(6) + "-" + typ + "-" + RandomString(2) + "-" + RandomString(4))
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

func RandomPastDate(year int) string {
	t := time.Date(RandomIntInRange(1900, year), time.Month(RandomInt(12)), RandomInt(28), RandomInt(24), RandomInt(59), 16, 0, time.UTC)
	tmp := t.Format(time.RFC3339)
	return tmp
}

func GenerateFutureDate(a time.Duration) string {
	t := time.Now()
	f := t.Add(a)
	tmp := f.Format(time.RFC3339)
	return tmp
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

func GenerateRandomUTCPastDate(year int) string {
	t := time.Date(RandomIntInRange(1900, year), time.Month(RandomInt(12)), RandomInt(28), RandomInt(24), RandomInt(59), 16, 0, time.UTC).UTC()
	tmp := t.Format(time.RFC3339)
	return tmp
}

func GenerateUTCFutureDate(a time.Duration) string {
	t := time.Now().UTC()
	f := t.Add(a)
	tmp := f.Format(time.RFC3339)
	return tmp
}

func GenerateUTCDate() string {
	t := time.Now().UTC()
	tmp := t.Format(time.RFC3339)
	return tmp
}


func Contains(elems []string, v string) bool {
    for _, s := range elems {
        if v == s {
            return true
        }
    }
    return false
}
