/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
