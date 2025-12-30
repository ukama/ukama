#!/bin/bash
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2024-present, Ukama Inc.

set -euo pipefail

gps_data_count=10
checksum_ok=-1
sentce_syntax_ok=-1
Latitude=""
Longitude=""

# --- Ukama virtual-node support ---
# UKAMA_GPS_MODE=mock|real  (default: real)
# UKAMA_GPS_COORDS="lat,lon" (default in mock: South Pole)
UKAMA_GPS_MODE="${UKAMA_GPS_MODE:-real}"
UKAMA_GPS_COORDS_DEFAULT="-90.000000,0.000000"

# Align with gps.d expectations
GPS_LOC_FILE="/tmp/gps_loc.log"
GPS_RAW_FILE="/tmp/gps_raw.txt"

# Function to convert NMEA format to decimal degrees for Google Maps
google_api_convert() {
	tmp=$(echo $1 | awk -F "." '{print $1}')
	degree=$(($tmp / 100))
	min=$(echo "$tmp % 100" | bc)
	tmp=$(echo $1 | awk -F "." '{print $2}' | awk -F "," '{print $1}')
	point=0.$tmp
	min=$(echo "$min + $point" | bc)
	tmp1=$(echo "$min / 60" | bc)
	google_format=$(echo "$degree + $tmp1" | bc)
	dir=$(echo $1 | awk -F "," '{print $2}')
	if [[ $dir == 'S' || $dir == 'W' ]]; then
		dir="-"
	else
		dir=""
	fi
	echo "${dir}${google_format}"
}

# Function to check checksum of GPS NMEA string
checksum_check() {
	i=1
	nmea_str=$(echo $1 | awk -F "*" '{print $1}' | awk -F "$" '{print $2}')
	ori_sum=$(echo $1 | awk -F "*" '{print $2}')

	nmea_str_len=${#nmea_str}
	checksum=$(printf %d "'${nmea_str:0:1}")

	while true; do
		more_char=$(printf %d "'${nmea_str:$i:1}")
		checksum=$((checksum ^ more_char))
		i=$((i + 1))
		if [ $i -eq $nmea_str_len ]; then break; fi
	done

	checksum=$(printf %02X $checksum)

	if [[ $ori_sum == $checksum ]]; then
		checksum_ok=0
	else
		checksum_ok=-1
	fi
}

# --- Mock mode helpers ---
mock_write_coords() {
	local coords="${UKAMA_GPS_COORDS:-$UKAMA_GPS_COORDS_DEFAULT}"

	# Basic validation: must contain exactly one comma
	if [[ "$coords" != *,* ]]; then
		echo "[$0] Invalid UKAMA_GPS_COORDS='$coords' (expected 'lat,lon')" >&2
		exit 1
	fi

	echo "$coords" > "$GPS_LOC_FILE"
	exit 0
}

mock_fix() {
	# In mock mode, assume fix is OK (locked)
	exit 0
}

mock_get_data() {
	# No-op in mock mode; create raw file for compatibility if needed
	: > "$GPS_RAW_FILE"
	exit 0
}

# Function to gather GPS data from remote host
gather_gps_data() {

	trx_host=$1
	missing_counter=0
	rm -f "$GPS_RAW_FILE"

	while [ $missing_counter -lt 10 ]; do
		echo "[$0] Gathering GPS Data"

        rsync -avz "$trx_host:$GPS_RAW_FILE" "$GPS_RAW_FILE" &
        pid=$!
		i=0
		while true; do
			if [ $i -eq $gps_data_count ]; then
				break
			fi
			i=$((i + 1))
			echo -n "."
			sleep 1
		done

		kill -9 $pid || true
		echo

		gpgsv_present=$(grep -i -c gpgsv "$GPS_RAW_FILE" || true)
		gprmc_present=$(grep -i -c g*rmc "$GPS_RAW_FILE" || true)
		gpgsa_present=$(grep -i -c gpgsa "$GPS_RAW_FILE" || true)

		if [ $gpgsv_present -eq 0 ]; then
			echo "[$0] Incomplete GPS data: GPGSV"
			missing_counter=$((missing_counter + 1))
			continue
		elif [ $gprmc_present -eq 0 ]; then
			echo "[$0] Incomplete GPS data: G*RMC"
			missing_counter=$((missing_counter + 1))
			continue
		elif [ $gpgsa_present -eq 0 ]; then
			echo "[$0] Incomplete GPS data: GPGSA"
			missing_counter=$((missing_counter + 1))
			continue
		else
			break
		fi
	done

	if [ $missing_counter -eq 10 ]; then
		echo "GPS data gathering failed - Exiting"
		exit 1
	fi
}

# Function to check GPS fix
gps_fix() {

	gps_fix_str=$(grep -a -m 1 "GPGSA" "$GPS_RAW_FILE")
	checksum_check "$gps_fix_str"

	if [ $checksum_ok -eq -1 ]; then
		exit 1
	else
		exit 0
	fi
}

# Function to get coordinates from GPS data
get_coordinates() {

	rmc_str=$(grep -a -m 1 "G*RMC" "$GPS_RAW_FILE")
	checksum_check "$rmc_str"

	if [ $checksum_ok -eq -1 ]; then
		exit 1
	else
		Latitude=$(echo "$rmc_str" | awk -F "," '{print $4","$5}')
		Longitude=$(echo "$rmc_str" | awk -F "," '{print $6","$7}')
	fi

	google_lat=$(google_api_convert "$Latitude")
	google_long=$(google_api_convert "$Longitude")

	echo "$google_lat,$google_long" > "$GPS_LOC_FILE"
}

# Main
case "${1:-}" in
	"get_gps_data")
		if [[ "$UKAMA_GPS_MODE" == "mock" ]]; then
			mock_get_data
		fi
		gather_gps_data "${2:-}"
		;;
	"gps_fix")
		if [[ "$UKAMA_GPS_MODE" == "mock" ]]; then
			mock_fix
		fi
		gps_fix
		;;
	"get_coordinates")
		if [[ "$UKAMA_GPS_MODE" == "mock" ]]; then
			mock_write_coords
		fi
		get_coordinates
		;;
	*)
		echo "Invalid argument. Use get_gps_data, gps_fix, or get_coordinates." >&2
		exit 1
		;;
esac
