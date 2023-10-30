/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package archiver

import (
	"compress/flate"
	"fmt"
	"os"

	arc "github.com/mholt/archiver/v3"
	log "github.com/sirupsen/logrus"
)

/* Set default tar options */
func defaultTarOptions() *arc.Tar {
	return &arc.Tar{
		OverwriteExisting:      true,
		MkdirAll:               true,
		ImplicitTopLevelFolder: true,
		StripComponents:        0,
		ContinueOnError:        true,
	}
}

/* Check if a file exist */
func validateSource(fname string) bool {
	fileinfo, err := os.Stat(fname)
	if os.IsNotExist(err) {
		log.Errorf("Error validating source for tar %s.", err.Error())

		return false
	}

	return !fileinfo.IsDir()
}

/* Check if file a directory */
func validateDir(fname string) bool {
	fileinfo, err := os.Stat(fname)
	if os.IsNotExist(err) {
		log.Errorf("Error validating directory %s.", err.Error())

		return false
	}

	return fileinfo.IsDir()
}

/* Unarchive a tar.gz capp file
   lib does support multiple format extension but our capp is only having targ.gz */
func Unarchive(fname string, dest string) error {
	log.Debugf("Unarchiving contents of %s", fname)

	/* Validate source */
	if !validateSource(fname) {
		return fmt.Errorf("invalid source '%s' to untar", fname)
	}

	/* Validate destination */
	if !validateDir(dest) {
		return fmt.Errorf("invalid destination '%s' to untar", dest)
	}

	/* archive interface */
	iface, err := arc.ByExtension(fname)
	if err != nil {
		log.Errorf("Error reading tar type for %s source %s.", fname, err.Error())

		return err
	}

	/* Setting archive options */
	switch v := iface.(type) {
	case *arc.TarGz:
		v.Tar = defaultTarOptions()
		v.CompressionLevel = flate.DefaultCompression
	default:
		log.Errorf("Error unkown tar format for source %s", fname)

		return fmt.Errorf("unkown format for tar source")
	}

	u, ok := iface.(arc.Unarchiver)
	if !ok {
		return fmt.Errorf("no matching unarchive fromatter found")
	}

	/* Unarchiving */
	err = u.Unarchive(fname, dest)
	if err != nil {
		log.Debugf("Failed to unarchive %s to %s: %s", fname, dest, err.Error())

		return err
	} else {
		log.Debugf("Unarchive %s to %s", fname, dest)
	}

	return err
}
