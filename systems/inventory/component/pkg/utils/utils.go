/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package utils

import (
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/gitClient"
)

type Component struct {
	Category      string `json:"category" yaml:"category"`
	Type          string `json:"type" yaml:"type"`
	Description   string `json:"description" yaml:"description"`
	UserId        string `json:"ownerId" yaml:"ownerId"`
	ImagesURL     string `json:"imagesURL" yaml:"imagesURL"`
	DatasheetURL  string `json:"datasheetURL" yaml:"datasheetURL"`
	InventoryID   string `json:"inventoryID" yaml:"inventoryID"`
	PartNumber    string `json:"partNumber" yaml:"partNumber"`
	Manufacturer  string `json:"manufacturer" yaml:"manufacturer"`
	Managed       string `json:"managed" yaml:"managed"`
	Warranty      uint32 `json:"warranty" yaml:"warranty"`
	Specification string `json:"specification" yaml:"specification"`
}

func GetEnvironmentField(env gitClient.Environment, envName string) []gitClient.Company {
	envName = strings.TrimSpace(envName)
	if envName == "" {
		envName = "production" // default
	}

	envNameRunes := []rune(strings.ToLower(envName))
	if len(envNameRunes) > 0 {
		envNameRunes[0] = []rune(strings.ToUpper(string(envNameRunes[0])))[0]
		envName = string(envNameRunes)
	}

	envValue := reflect.ValueOf(env)
	field := envValue.FieldByName(envName)

	if !field.IsValid() || field.Kind() != reflect.Slice {
		log.Warnf("Environment field '%s' not found, defaulting to Production", envName)
		return env.Production
	}

	if field.CanInterface() {
		if companies, ok := field.Interface().([]gitClient.Company); ok {
			return companies
		}
	}

	log.Warnf("Failed to convert environment field '%s', defaulting to Production", envName)
	return env.Production
}
