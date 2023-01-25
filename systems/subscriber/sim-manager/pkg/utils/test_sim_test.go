package utils

import (
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_TestSimUUID(t *testing.T) {
	for i := 0; i < 5; i++ {
		tt := func() string {
			s := uuid.NewV4().String()
			return strings.Replace(s, s[:8], testUUUIDPrefix, 1)
		}()

		t.Run("testUuudIsValid: "+tt, func(t *testing.T) {
			err := testSimUUID(tt)
			assert.NoError(t, err)
		})
	}

	tests := []string{
		"testxuid-89a6-42f6-9f54-46a85fcbe539",
		"4e601397-e639-4857-9b39-7824856ab112",
		"4e601397-xxxuuid-428f-46b5-bc73-a39d821b61eb",
		"testuuid84894fc2bfb58849d15b88a9",
		"testuuid07671441c29a31f9a2f64024651d",
	}

	for i := 0; i < 5; i++ {
		tt := tests[i]

		t.Run("testUuudIsNotValid: "+tt, func(t *testing.T) {
			err := testSimUUID(tt)
			assert.Error(t, err)
		})
	}
}

func Test_GetIccidFromTestSimUUID(t *testing.T) {
	tests := []struct {
		testUUID string
		iccid    string
	}{

		{testUUID: "testxuid-89a6-42f6-9f54-46a85fcbe539",
			iccid: ""},

		{testUUID: "testuuid-fd1b-4163-be40-67fdffe98101",
			iccid: "fd1b4163be4067fdffe98101"},

		{testUUID: "4e601397-e639-4857-9b39-7824856ab112",
			iccid: ""},

		{testUUID: "testuuid-382e-4be8-9f41-db480d48677b",
			iccid: "382e4be89f41db480d48677b"},

		{testUUID: "4e601397-xxxuuid-428f-46b5-bc73-a39d821b61eb",
			iccid: ""},

		{testUUID: "testuuid-14f0-4592-bd94-c20daae240d7",
			iccid: "14f04592bd94c20daae240d7"},

		{testUUID: "testuuid84894fc2bfb58849d15b88a9",
			iccid: ""},

		{testUUID: "testuuid-4bfb-4f4c-b03e-b88294816041",
			iccid: "4bfb4f4cb03eb88294816041"},

		{testUUID: "testuuid07671441c29a31f9a2f64024651d",
			iccid: ""},

		{testUUID: "testuuid-2719-4056-b437-604568556f6b",
			iccid: "27194056b437604568556f6b"},
	}

	for _, tt := range tests {
		t.Run(tt.testUUID, func(t *testing.T) {
			iccid, _ := GetIccidFromTestSimUUID(tt.testUUID)
			assert.Equal(t, iccid, tt.iccid)
		})
	}
}
