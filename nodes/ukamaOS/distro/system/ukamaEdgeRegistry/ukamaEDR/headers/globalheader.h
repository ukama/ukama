/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef HEADERS_GLOBALHEADERS_H_
#define HEADERS_GLOBALHEADERS_H_

#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>

#define TRUE 				true
#define FALSE				false

#define UKAMA_FREE(mem)			if(mem) { \
									free(mem);\
									mem = NULL; \
								}

#define MAX(a,b)			(((a) < (b))?(b):(a))
#define MIN(a,b)			(((a) < (b))?(a):(b))

typedef enum ReturnValue {
	RET_GEN_ERROR = -1, /* Anything less than zero is error.*/
	RET_SUCCESS	  =  0,
} ReturnValue;

#endif /* HEADERS_GLOBALHEADERS_H_*/
