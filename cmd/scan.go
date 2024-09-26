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

	"github.com/ukama/msgcli/internal/scan"
)

const (
	defaultOutputFormat = "json"
)

var (
	oFile        = os.Stdout
	outputFormat = scan.Param{
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
		src, err := cmd.Flags().GetString("source-file")
		if err != nil {
			return err
		}

		ofilename, err := cmd.Flags().GetString("output-file")
		if err != nil {
			return err
		}

		if ofilename != "" {
			oFile, err = os.Create(ofilename)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer oFile.Close()
		}

		if outputFormat.String() == "" {
			err = outputFormat.Set(defaultOutputFormat)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		cfg := &scan.Config{
			OutputFormat: outputFormat.String(),
		}

		return scan.Run(src, oFile, cfg)
	},
}

func init() {
	eventsCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP("source-file", "s", ".", "Source directory to start from")
	scanCmd.Flags().StringP("output-file", "o", "", "The name of the file to write to (default \"Stdout\")")
	scanCmd.Flags().VarP(&outputFormat, "output-format", "f",
		fmt.Sprintf("Output format. Must match one of the following: %q (default \"json\" )",
			outputFormat.Values))
}
