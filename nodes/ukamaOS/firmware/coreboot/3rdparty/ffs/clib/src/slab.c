/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/src/slab.c $                                             */
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
 *   File: slab.c
 * Author: Shaun Wetzstein <shaun@us.ibm.com>
 *  Descr: Slab based allocator
 *   Note:
 *   Date: 06/05/10
 */

#include <unistd.h>
#include <stdlib.h>
#include <malloc.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <limits.h>

#include "libclib.h"
#include "tree_iter.h"

#include "slab.h"

/* ======================================================================= */

/*! @cond */
#define SLAB_NODE_MAGIC             "SLND"

#define SLAB_NODE_MAGIC_CHECK(m) ({                       \
    bool rc = (((m)[0] != SLAB_NODE_MAGIC[0]) ||          \
               ((m)[1] != SLAB_NODE_MAGIC[1]) ||          \
               ((m)[2] != SLAB_NODE_MAGIC[2]) ||          \
               ((m)[3] != SLAB_NODE_MAGIC[3]));           \
    rc;                                                   \
})

typedef struct slab_node slab_node_t;

struct slab_node {
	uint8_t magic[4];

	uint32_t free;
	tree_node_t node;

	uint32_t bitmap[];
};

#define SLAB_ALLOC_COUNT(s)	((s)->hdr.data_size / (s)->hdr.alloc_size)

#define SLAB_PAGE_MAX		UINT16_MAX
#define SLAB_PAGE_DIVISOR	32
/*! @endcond */

/* ======================================================================= */

static slab_node_t *__slab_grow(slab_t * self)
{
	assert(self != NULL);
	assert(!MAGIC_CHECK(self->hdr.id, SLAB_MAGIC));

	slab_node_t *node = NULL;
	if (posix_memalign((void **)&node, self->hdr.align_size,
			   self->hdr.page_size) < 0) {
		ERRNO(errno);
		return NULL;
	}

	memset(node, 0, self->hdr.page_size);

	node->magic[0] = SLAB_NODE_MAGIC[0];
	node->magic[1] = SLAB_NODE_MAGIC[1];
	node->magic[2] = SLAB_NODE_MAGIC[2];
	node->magic[3] = SLAB_NODE_MAGIC[3];

	node->free = SLAB_ALLOC_COUNT(self);

	tree_node_init(&node->node, (const void *)int64_hash1((int64_t) node));
	splay_insert(&self->tree, &node->node);

	self->hdr.page_count++;

	return node;
}

/* ======================================================================= */

int slab_init3(slab_t * self, const char *name, uint32_t alloc_size)
{
	size_t page_size = max(sysconf(_SC_PAGESIZE),
			       __round_pow2(alloc_size * SLAB_PAGE_DIVISOR));
	return slab_init5(self, name, alloc_size, page_size, page_size);
}

int slab_init4(slab_t * self, const char *name, uint32_t alloc_size,
	       size_t page_size)
{
	size_t align_size = (size_t) sysconf(_SC_PAGESIZE);
	return slab_init5(self, name, alloc_size, page_size, align_size);
}

int slab_init5(slab_t * self, const char *name, uint32_t alloc_size,
	       size_t page_size, size_t align_size)
{
	assert(self != NULL);

	if (MAGIC_CHECK(self->hdr.id, SLAB_MAGIC) == false)
		slab_delete(self);

	alloc_size = align(alloc_size, sizeof(void *));
	if (alloc_size < SLAB_ALLOC_MIN || SLAB_ALLOC_MAX < alloc_size) {
		UNEXPECTED("alloc_size out of range [%d..%d]",
			   SLAB_ALLOC_MIN, SLAB_ALLOC_MAX);
		return -1;
	}

	page_size = __round_pow2(page_size);
	if (page_size / alloc_size < SLAB_PAGE_DIVISOR) {
		UNEXPECTED("page_size out of range [%d..%d]",
			   alloc_size * SLAB_PAGE_DIVISOR, SLAB_PAGE_MAX);
		return -1;
	}

	uint32_t __page_size = (uint32_t) sysconf(_SC_PAGESIZE);

	align_size = __round_pow2(align_size);
	if (align_size % __page_size) {
		UNEXPECTED("align_size must be 0x%x aligned", __page_size);
		return -1;
	}

	memset(self, 0, sizeof *self);

	self->hdr.id[IDENT_MAGIC_0] = SLAB_MAGIC[IDENT_MAGIC_0];
	self->hdr.id[IDENT_MAGIC_1] = SLAB_MAGIC[IDENT_MAGIC_1];
	self->hdr.id[IDENT_MAGIC_2] = SLAB_MAGIC[IDENT_MAGIC_2];
	self->hdr.id[IDENT_MAGIC_3] = SLAB_MAGIC[IDENT_MAGIC_3];

	self->hdr.id[IDENT_MAJOR] = CLIB_MAJOR;
	self->hdr.id[IDENT_MINOR] = CLIB_MINOR;
	self->hdr.id[IDENT_PATCH] = CLIB_PATCH;

	if (__BYTE_ORDER == __LITTLE_ENDIAN)
		self->hdr.id[IDENT_FLAGS] |= SLAB_FLAG_LSB;
	if (__BYTE_ORDER == __BIG_ENDIAN)
		self->hdr.id[IDENT_FLAGS] |= SLAB_FLAG_MSB;

	if (name != NULL && *name != '\0')
		strncpy(self->hdr.name, name, sizeof(self->hdr.name));

	self->hdr.page_size = page_size;
	self->hdr.align_size = align_size;
	self->hdr.alloc_size = alloc_size;

	self->hdr.bitmap_size =
	    align(page_size / alloc_size, CHAR_BIT * sizeof(uint32_t));
	self->hdr.bitmap_size /= CHAR_BIT;

	self->hdr.data_size =
	    self->hdr.page_size - sizeof(slab_node_t) - self->hdr.bitmap_size;

	tree_init(&self->tree, default_compare);

	return 0;
}

int slab_delete(slab_t * self)
{
	if (self == NULL)
		return 0;

	if (MAGIC_CHECK(self->hdr.id, SLAB_MAGIC)) {
		UNEXPECTED("'%s' invalid or corrupt slab object",
			   self->hdr.name);
		return -1;
	}

	tree_iter_t it;
	tree_iter_init(&it, &self->tree, TI_FLAG_FWD);

	slab_node_t *node;
	tree_for_each(&it, node, node) {
		if (SLAB_NODE_MAGIC_CHECK(node->magic)) {
			UNEXPECTED("'%s' invalid or corrupt slab_node object "
				   "=> '%.4s'", self->hdr.name, node->magic);
			return -1;
		}

		if (splay_remove(&self->tree, &node->node) < 0)
			return -1;
		memset(node, 0, sizeof(*node));

		free(node);
	}

	memset(self, 0, sizeof(*self));

	return 0;
}

void *slab_alloc(slab_t * self)
{
	assert(self != NULL);
	assert(!MAGIC_CHECK(self->hdr.id, SLAB_MAGIC));

	tree_iter_t it;
	tree_iter_init(&it, &self->tree, TI_FLAG_FWD);

	slab_node_t *node;
	tree_for_each(&it, node, node)
	    if (0 < node->free)
		break;

	if (tree_iter_elem(&it) == NULL)
		node = __slab_grow(self);

	assert(node != NULL);
	assert(0 < node->free);

	void *ptr = NULL;

	uint32_t map_pos;
	for (map_pos = 0; map_pos < self->hdr.bitmap_size; map_pos++)
		if (node->bitmap[map_pos] != UINT32_MAX)
			break;

	if (node->bitmap[map_pos] == UINT32_MAX) {
		UNEXPECTED("'%s' cache is corrupted", self->hdr.name);
		return NULL;
	}
	if (self->hdr.bitmap_size <= map_pos) {
		UNEXPECTED("'%s' cache is corrupted", self->hdr.name);
		return NULL;
	}

	uint32_t bit = clzl(~node->bitmap[map_pos]);
	uint32_t mask = 0x80000000 >> bit;

	if ((node->bitmap[map_pos] & mask) == mask) {
		UNEXPECTED("'%s' cache is corrupted", self->hdr.name);
		return NULL;
	}

	node->bitmap[map_pos] |= mask;
	node->free--;

	ptr = (void *)node->bitmap + self->hdr.bitmap_size +
	    (map_pos * INT32_BIT + bit) * self->hdr.alloc_size;

	return ptr;
}

int slab_free(slab_t * self, void *ptr)
{
	assert(self != NULL);
	assert(!MAGIC_CHECK(self->hdr.id, SLAB_MAGIC));

	if (ptr == NULL)
		return 0;

	slab_node_t *node = (slab_node_t *) ((uintptr_t) ptr &
					     ~(self->hdr.page_size - 1));
	assert(node != NULL);

	if (SLAB_NODE_MAGIC_CHECK(node->magic)) {
		int64_t hash = int64_hash1((int64_t) node);
		if (splay_find(&self->tree, (const void *)hash) == NULL) {
			UNEXPECTED("'%s' invalid slab_node pointer, %p",
				   self->hdr.name, ptr);
			return -1;
		}
	}

	void *data = (void *)node->bitmap + self->hdr.bitmap_size;
	assert(data != NULL);

	if (ptr < data) {
		UNEXPECTED("'%s' pointer out-of-range, %p",
			   self->hdr.name, ptr);
		return -1;
	}

	size_t slot = (ptr - data) / self->hdr.alloc_size;
	uint32_t mask = 0x80000000 >> slot;
	size_t map_pos = slot / INT32_BIT;

	if ((node->bitmap[map_pos] & mask) != mask) {
		UNEXPECTED("'%s' double free detected, %p",
			   self->hdr.name, ptr);
		return -1;
	}

	node->bitmap[map_pos] &= ~mask;
	node->free++;

	if (SLAB_ALLOC_COUNT(self) - node->free <= 0) {
		splay_remove(&self->tree, &node->node);
		self->hdr.page_count = max(0UL, self->hdr.page_count - 1);
		memset(node, 0, sizeof(*node));
		free(node);
	}

	return 0;
}

size_t slab_alloc_size(slab_t * self)
{
	assert(self != NULL);
	return self->hdr.alloc_size;
}

size_t slab_page_size(slab_t * self)
{
	assert(self != NULL);
	return self->hdr.page_size;
}

size_t slab_data_size(slab_t * self)
{
	assert(self != NULL);
	return self->hdr.data_size;
}

size_t slab_bitmap_size(slab_t * self)
{
	assert(self != NULL);
	return self->hdr.bitmap_size;
}

size_t slab_align_size(slab_t * self)
{
	assert(self != NULL);
	return self->hdr.align_size;
}

void slab_dump(slab_t * self, FILE * out)
{
	if (out == NULL)
		out = stdout;

	if (self != NULL) {
		if (MAGIC_CHECK(self->hdr.id, SLAB_MAGIC)) {
			UNEXPECTED("'%s' invalid or corrupt slab object",
				   self->hdr.name);
			return;
		}

		fprintf(out, "%s: page_size: %d bitmap_size: %d data_size: "
			"%d alloc_size: %d -- page_count: %d\n",
			self->hdr.name, self->hdr.page_size,
			self->hdr.bitmap_size, self->hdr.data_size,
			self->hdr.alloc_size, self->hdr.page_count);

		tree_iter_t it;
		tree_iter_init(&it, &self->tree, TI_FLAG_FWD);

		slab_node_t *node;
		tree_for_each(&it, node, node) {
			fprintf(out, "magic[%.4s] node: %p bitmap: %p data: "
				"%p -- alloc: %d free: %d\n",
				node->magic, &node->node, node->bitmap,
				(void *)node->bitmap + self->hdr.bitmap_size,
				self->hdr.data_size / self->hdr.alloc_size -
				node->free, node->free);

			dump_memory(out, (unsigned long)node, node,
				    self->hdr.page_size);
		}
	}
}

/* ======================================================================= */
