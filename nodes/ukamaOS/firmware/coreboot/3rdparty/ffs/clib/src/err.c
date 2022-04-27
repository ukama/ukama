/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/src/err.c $                                              */
/*                                                                        */
/* OpenPOWER FFS Project                                                  */
/*                                                                        */
/* Contributors Listed Below - COPYRIGHT 2014,2015                        */
/* [+] International Business Machines Corp.                              */
/*                                                                        */
/*                                                                        */
/* Licensed under the Apache License, Version 2.0 (the "License");        */
/* you may not use this file except in compliance with the License.       */
/* You may obtain a copy of the License at                                */
/*                                                                        */
/*     http://www.apache.org/licenses/LICENSE-2.0                         */
/*                                                                        */
/* Unless required by applicable law or agreed to in writing, software    */
/* distributed under the License is distributed on an "AS IS" BASIS,      */
/* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or        */
/* implied. See the License for the specific language governing           */
/* permissions and limitations under the License.                         */
/*                                                                        */
/* IBM_PROLOG_END_TAG                                                     */

/*
 *   File: err.c
 * Author: Shaun Wetzstein <shaun@us.ibm.com>
 *  Descr: Error container
 *   Date: 04/04/12
 */

#include <errno.h>
#include <assert.h>
#include <pthread.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdarg.h>
#include <execinfo.h>

#include "attribute.h"
#include "min.h"

#include "err.h"

static list_t *__err_key = 0;

static const char *__err_type_name[] = {
	[ERR_NONE] = "none",
	[ERR_ERRNO] = "errno",
	[ERR_UNEXPECTED] = "unexpected",
	[ERR_VERSION] = "version",
};

/* =======================================================================*/

void err_delete(err_t * self)
{
	assert(self != NULL);
	assert(self->magic == ERR_MAGIC);

	memset(self, 0, sizeof(*self));
	free(self);

	return;
}

err_t *err_get(void)
{
	list_t *list = __err_key;

	err_t *self = NULL;

	if (list != NULL) {
		self = container_of(list_remove_tail(list), err_t, node);

		if (list_empty(list)) {
			free(list), list = NULL;
			__err_key = list;
		}
	}

	return self;
}

void err_put(err_t * self)
{
	assert(self != NULL);

	list_t *list = __err_key;
	if (list == NULL) {
		list = (list_t *) malloc(sizeof(*list));
		assert(list != NULL);

		list_init(list);
		__err_key = list;
	}

	list_add_head(list, &self->node);

	return;
}

int err_type1(err_t * self)
{
	assert(self != NULL);
	assert(self->magic == ERR_MAGIC);

	return self->type;
}

void err_type2(err_t * self, int type)
{
	assert(self != NULL);
	assert(self->magic == ERR_MAGIC);

	self->type = type;
}

int err_code1(err_t * self)
{
	assert(self != NULL);
	assert(self->magic == ERR_MAGIC);

	return self->code;
}

void err_code2(err_t * self, int code)
{
	assert(self != NULL);
	assert(self->magic == ERR_MAGIC);

	self->code = code;

	return;
}

const char *err_file1(err_t * self)
{
	assert(self != NULL);
	assert(self->magic == ERR_MAGIC);

	return self->file;
}

void err_file2(err_t * self, const char *file)
{
	assert(self != NULL);
	assert(self->magic == ERR_MAGIC);

	self->file = file;

	return;
}

int err_line1(err_t * self)
{
	assert(self != NULL);
	assert(self->magic == ERR_MAGIC);

	return self->line;
}

void err_line2(err_t * self, int line)
{
	assert(self != NULL);
	assert(self->magic == ERR_MAGIC);

	self->line = line;

	return;
}

const void *err_data1(err_t * self)
{
	assert(self != NULL);
	assert(self->magic == ERR_MAGIC);

	return self->data;
}

void err_data2(err_t * self, int size, const void *data)
{
	assert(self != NULL);
	assert(self->magic == ERR_MAGIC);

	self->size = min(size, ERR_DATA_SIZE);
	memcpy(self->data, data, self->size);

	return;
}

int err_size(err_t * self)
{
	assert(self != NULL);
	assert(self->magic == ERR_MAGIC);

	return self->size;
}

const char *err_type_name(err_t * self)
{
	assert(self != NULL);
	assert(self->magic == ERR_MAGIC);

	return __err_type_name[self->type];
}

