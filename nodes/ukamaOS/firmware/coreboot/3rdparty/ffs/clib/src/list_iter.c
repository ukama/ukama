/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/src/list_iter.c $                                        */
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
 *   File: list_iter.c
 * Author: Shaun Wetzstein <shaun@us.ibm.com>
 *  Descr: List iterator
 *   Note:
 *   Date: 10/22/10
 */

#include <unistd.h>
#include <stdarg.h>
#include <stdlib.h>
#include <malloc.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <limits.h>

#include "libclib.h"
#include "list_iter.h"

/* ======================================================================= */

static inline list_node_t *__list_iter_bwd(list_iter_t * self)
{
	assert(self != NULL);

	assert(self->list != NULL);
	assert(self->node != NULL);

	if (self->node == &self->list->node)
		self->node = NULL;
	else
		self->node = list_node_prev(self->node);

	return self->node;
}

static inline list_node_t *__list_iter_fwd(list_iter_t * self)
{
	assert(self != NULL);
	assert(self->list != NULL);
	assert(self->node != NULL);

	if (self->node == &self->list->node)
		self->node = NULL;
	else
		self->node = list_node_next(self->node);

	return self->node;
}

void list_iter_init(list_iter_t * self, list_t * list, uint32_t flags)
{
	assert(self != NULL);
	assert(list != NULL);

	self->flags = flags;
	self->list = list;

	if (self->flags & LI_FLAG_BWD)
		self->node = list_tail(self->list);
	else
		self->node = list_head(self->list);
}

void list_iter_clear(list_iter_t * self)
{
	assert(self != NULL);

	if (self->flags & LI_FLAG_BWD)
		self->node = list_tail(self->list);
	else
		self->node = list_head(self->list);
}

list_node_t *list_iter_elem(list_iter_t * self)
{
	assert(self != NULL);
	return self->node;
}

list_node_t *list_iter_inc1(list_iter_t * self)
{
	return list_iter_inc2(self, 1);
}

list_node_t *list_iter_inc2(list_iter_t * self, size_t count)
{
	assert(self != NULL);

	for (size_t i = 0; i < count && self->node != NULL; i++) {
		if (self->flags & LI_FLAG_BWD)
			__list_iter_bwd(self);
		else
			__list_iter_fwd(self);
	}

	return self->node;
}

list_node_t *list_iter_dec1(list_iter_t * self)
{
	return list_iter_dec2(self, 1);
}

list_node_t *list_iter_dec2(list_iter_t * self, size_t count)
{
	assert(self != NULL);

	for (size_t i = 0; i < count && self->node != NULL; i++) {
		if (self->flags & LI_FLAG_BWD)
			__list_iter_fwd(self);
		else
			__list_iter_bwd(self);
	}

	return self->node;
}

/* ======================================================================= */
