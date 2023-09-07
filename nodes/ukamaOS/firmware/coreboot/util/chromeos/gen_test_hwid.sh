#!/bin/sh
#
# This file is part of the coreboot project.
#
# Copyright 2019 Google Inc.
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; version 2 of the License.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

main() {
  if [ "$#" != 1 ]; then
    echo "Usage: $0 MAINBOARD_PARTNUMBER" >&2
    exit 1
  fi

  # Generate a test-only Chrome OS HWID v2 string
  local board="$1"
  local prefix="$(echo "${board}" | tr a-z A-Z) TEST"
  # gzip has second-to-last 4 bytes in CRC32.
  local crc32="$(printf "${prefix}" | gzip -1 | tail -c 8 | head -c 4 | \
		 hexdump -e '1/4 "%04u" ""' | tail -c 4)"

  echo "${prefix}" "${crc32}"
}
main "$@"
