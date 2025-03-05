/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"github.com/jszwec/csvutil"
	"github.com/ukama/ukama/testing/services/dummy/dsimfactory/pkg/db"
)

type RawSim struct {
	Iccid          string `csv:"ICCID"`
	Msisdn         string `csv:"MSISDN"`
	SmDpAddress    string `csv:"SmDpAddress"`
	ActivationCode string `csv:"ActivationCode"`
	QrCode         string `csv:"QrCode"`
	IsPhysical     string `csv:"IsPhysical"`
	Imsi           string `csv:"IMSI"`
}

func ParseBytesToRawSim(b []byte) ([]RawSim, error) {
	var r []RawSim
	err := csvutil.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func RawSimToPb(r []RawSim) []db.Sim {
	var s []db.Sim
	for _, value := range r {
		s = append(s, db.Sim{
			Iccid:          value.Iccid,
			Msisdn:         value.Msisdn,
			SmDpAddress:    value.SmDpAddress,
			ActivationCode: value.ActivationCode,
			IsPhysical:     value.IsPhysical == "TRUE",
			QrCode:         value.QrCode,
			Imsi:           value.Imsi,
		})
	}
	return s
}
