/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/src/tree_iter.c $                                        */
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
 *   File: tree_iter.c
 * Author: Shaun Wetzstein <shaun@us.ibm.com>
 *  Descr:
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
#include "tree_iter.h"

/* ======================================================================= */

int tree_iter_init(tree_iter_t * self, tree_t * tree, uint32_t flags)
{
	assert(self != NULL);
	assert(tree != NULL);

	self->node = NULL;
	self->safe = NULL;

	self->tree = tree;
	self->flags = flags;

	return tree_iter_clear(self);
}

int tree_iter_clear(tree_iter_t * self)
{
	assert(self != NULL);

	if (self->flags & TI_FLAG_FWD) {
		self->node = tree_min(self->tree);
		if (self->node != NULL)
			self->safe = tree_node_next(self->node);
	} else if (self->flags & TI_FLAG_BWD) {
		self->node = tree_max(self->tree);
		if (self->node != NULL)
			self->safe = tree_node_prev(self->node);
	} else {
		UNEXPECTED("invalid tree_iter flags");
		return -1;
	}

	return 0;
}

tree_node_t *tree_iter_elem(tree_iter_t * self)
{
	assert(self != NULL);
	return self->node;
}

tree_node_t *tree_iter_inc1(tree_iter_t * self)
{
	return tree_iter_inc2(self, 1);
}

tree_node_t *tree_iter_inc2(tree_iter_t * self, size_t count)
{
	assert(self != NULL);

	for (size_t i = 0; i < count && self->node != NULL; i++) {
		if (self->flags & TI_FLAG_FWD) {
			self->node = self->safe;
			if (self->node != NULL)
				self->safe = tree_node_next(self->node);
		} else if (self->flags & TI_FLAG_BWD) {
			self->node = self->safe;
			if (self->node != NULL)
				self->safe = tree_node_prev(self->node);
		} else {
			UNEXPECTED("invalid tree_iter flags");
			return NULL;
		}
	}

	return self->node;
}

tree_node_t *tree_iter_dec1(tree_iter_t * self)
{
	return tree_iter_dec2(self, 1);
}

tree_node_t *tree_iter_dec2(tree_iter_t * self, size_t count)
{
	assert(self != NULL);

	for (size_t i = 0; i < count && self->node != NULL; i++) {
		if (self->flags & TI_FLAG_FWD) {
			self->node = self->safe;
			if (self->node != NULL)
				self->safe = tree_node_prev(self->node);
		} else if (self->flags & TI_FLAG_BWD) {
			self->node = self->safe;
			if (self->node != NULL)
				self->safe = tree_node_next(self->node);
		} else {
			UNEXPECTED("invalid tree_iter flags");
			return NULL;
		}
	}

	return self->node;
}

/* ======================================================================= */
