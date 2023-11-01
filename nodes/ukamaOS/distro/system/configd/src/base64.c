/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <string.h>
#include "base64.h"

/* aaaack but it's fast and const should make it shared text page. */
static const unsigned char base64_char_set[256] =
{
		/* ASCII table */
		64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
		64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
		64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 62, 64, 64, 64, 63,
		52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 64, 64, 64, 64, 64, 64,
		64,  0,  1,  2,  3,  4,  5,  6,  7,  8,  9, 10, 11, 12, 13, 14,
		15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 64, 64, 64, 64, 64,
		64, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
		41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 64, 64, 64, 64, 64,
		64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
		64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
		64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
		64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
		64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
		64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
		64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
		64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64
};

int base64_decode_len(const char *encodedData)
{
	int outputBytes;
	const unsigned char *input;
	int inputBytes;

	input = (const unsigned char *) encodedData;
	while (base64_char_set[*(input++)] <= 63);

	inputBytes = (input - (const unsigned char *) encodedData) - 1;
	outputBytes = ((inputBytes + 3) / 4) * 3;

	return outputBytes + 1;
}

int base64_decode(char *decodedData, const char *encodedData)
{
	int outputBytes;
	const unsigned char *input;
	unsigned char *output;
	int inputBytes;

	input = (const unsigned char *) encodedData;
	while (base64_char_set[*(input++)] <= 63);
	inputBytes = (input - (const unsigned char *) encodedData) - 1;
	outputBytes = ((inputBytes + 3) / 4) * 3;

	output = (unsigned char *) decodedData;
	input = (const unsigned char *) encodedData;

	while (inputBytes > 4) {
		*(output++) =
				(unsigned char) (base64_char_set[*input] << 2 | base64_char_set[input[1]] >> 4);
		*(output++) =
				(unsigned char) (base64_char_set[input[1]] << 4 | base64_char_set[input[2]] >> 2);
		*(output++) =
				(unsigned char) (base64_char_set[input[2]] << 6 | base64_char_set[input[3]]);
		input += 4;
		inputBytes -= 4;
	}

	if (inputBytes > 1) {
		*(output++) =
				(unsigned char) (base64_char_set[*input] << 2 | base64_char_set[input[1]] >> 4);
	}
	if (inputBytes > 2) {
		*(output++) =
				(unsigned char) (base64_char_set[input[1]] << 4 | base64_char_set[input[2]] >> 2);
	}
	if (inputBytes > 3) {
		*(output++) =
				(unsigned char) (base64_char_set[input[2]] << 6 | base64_char_set[input[3]]);
	}

	*(output++) = '\0';
	outputBytes -= (4 - inputBytes) & 3;
	return outputBytes;
}

