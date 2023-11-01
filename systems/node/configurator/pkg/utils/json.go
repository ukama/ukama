/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/wI2L/jsondiff"
)

func JsonDiff(srcFile string, targetFile string) ([]string, bool, error) {
	change := false
	source, err := os.ReadFile(srcFile)
	if err != nil {
		log.Warningf("File may be added in latest commit. Error reading json file %s.  Error %+v", srcFile, err)
		source = []byte("{}")
	}

	target, err := os.ReadFile(targetFile)
	if err != nil {
		log.Warningf("File may be deleted in latest commit. Error reading json file %s: %v", targetFile, err)
		return nil, change, nil
	}

	patch, err := jsondiff.CompareJSON(source, target)
	if err != nil {
		log.Errorf("error comparing json file %s and %s: %v", srcFile, targetFile, err)
		return nil, change, err
	}
	var changedValues []string
	for _, op := range patch {
		fmt.Printf("%s\n", op)
		changedValues = append(changedValues, op.String())
		change = true
	}

	return changedValues, change, nil
}
