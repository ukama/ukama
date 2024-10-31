/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ukama/ukama/utils/msgcli/internal/scan"
	"github.com/ukama/ukama/utils/msgcli/util"
)

var (
	outputFile   = os.Stdout
	outputFormat = util.EnumParam{
		Values: []string{"json", "yaml", "toml"},
	}
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans for events",
	Long: `The scan command scans for events registered, listened and sent
by the given Ukama service.`,

	Aliases:      []string{"s"},
	SilenceUsage: true,

	RunE: func(cmd *cobra.Command, args []string) error {
		src, err := cmd.Flags().GetString("source")
		if err != nil {
			return err
		}

		outputFilename, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		if outputFilename != "" {
			outputFile, err = os.Create(outputFilename)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer outputFile.Close()
		}

		if outputFormat.String() == "" {
			err = outputFormat.Set(defaultOutputFormat)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		cfg := &util.Config{
			OutputFormat: outputFormat.String(),
		}

		return scan.Run(src, outputFile, cfg)
	},
}

func init() {
	eventsCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP("source", "s", ".", "source directory to start from")
	scanCmd.Flags().StringP("output", "o", "", "name of the file to write to (default \"Stdout\")")
	scanCmd.Flags().VarP(&outputFormat, "format", "f",
		fmt.Sprintf("output format. Must match one of the following: %q (default \"json\" )",
			outputFormat.Values))
}
