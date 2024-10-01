/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package scan

import (
	"fmt"
	"io"

	"github.com/ukama/msgcli/util"
)

func Run(dir string, out io.Writer, cfg *util.Config) error {
	data := &util.ResultSet{}

	err := WalkAndParse(dir, data)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	outputBuf, err := util.Serialize(data, cfg.OutputFormat)
	if err != nil {
		return fmt.Errorf("error while serializing output: %w", err)
	}

	_, err = fmt.Fprint(out, outputBuf)
	if err != nil {
		return fmt.Errorf("error while writting output: %w", err)
	}

	return nil
}
