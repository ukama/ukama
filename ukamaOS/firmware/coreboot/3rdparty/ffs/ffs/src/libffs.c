/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: ffs/src/libffs.c $                                            */
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
 *   File: libffs.c
 * Author: Shaun Wetzstein <shaun@us.ibm.com>
 *  Descr: FFS IO interface
 *   Note:
 *   Date: 05/07/12
 */

#include <sys/types.h>
#include <sys/stat.h>

#include <stdlib.h>
#include <stdarg.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <unistd.h>
#include <fcntl.h>
#include <endian.h>
#include <libgen.h>
#include <regex.h>

#include "libffs.h"

#include <clib/builtin.h>
#include <clib/checksum.h>
#include <clib/misc.h>
#include <clib/err.h>
#include <clib/raii.h>

#ifndef be32toh
#include <byteswap.h>
#if __BYTE_ORDER == __LITTLE_ENDIAN
#define be32toh(x) __bswap_32(x)
#define htobe32(x) __bswap_32(x)
#else
#define be32toh(x) (x)
#define htobe32(x) (x)
#endif
#endif

#define FFS_ENTRY_EXTENT	10UL

/* ============================================================ */

static void __hdr_be32toh(ffs_hdr_t * hdr)
{
	assert(hdr != NULL);

	hdr->magic = be32toh(hdr->magic);
	hdr->version = be32toh(hdr->version);
	hdr->size = be32toh(hdr->size);
	hdr->entry_size = be32toh(hdr->entry_size);
	hdr->entry_count = be32toh(hdr->entry_count);
	hdr->block_size = be32toh(hdr->block_size);
	hdr->block_count = be32toh(hdr->block_count);
	hdr->checksum = be32toh(hdr->checksum);
}

static void __hdr_htobe32(ffs_hdr_t * hdr)
{
	assert(hdr != NULL);

	hdr->magic = htobe32(hdr->magic);
	hdr->version = htobe32(hdr->version);
	hdr->size = htobe32(hdr->size);
	hdr->entry_size = htobe32(hdr->entry_size);
	hdr->entry_count = htobe32(hdr->entry_count);
	hdr->block_size = htobe32(hdr->block_size);
	hdr->block_count = htobe32(hdr->block_count);
	hdr->checksum = htobe32(hdr->checksum);
}

static void __entry_be32toh(ffs_entry_t * entry)
{
	assert(entry != NULL);

	entry->base = be32toh(entry->base);
	entry->size = be32toh(entry->size);
	entry->pid = be32toh(entry->pid);
	entry->id = be32toh(entry->id);
	entry->type = be32toh(entry->type);
	entry->flags = be32toh(entry->flags);
	entry->actual = be32toh(entry->actual);
	entry->checksum = be32toh(entry->checksum);

	entry->resvd[0] = be32toh(entry->resvd[0]);
	entry->resvd[1] = be32toh(entry->resvd[1]);
	entry->resvd[2] = be32toh(entry->resvd[2]);
	entry->resvd[3] = be32toh(entry->resvd[3]);

	for (int j = 0; j < FFS_USER_WORDS; j++)
		entry->user.data[j] = be32toh(entry->user.data[j]);
}

static void __entry_htobe32(ffs_entry_t * entry)
{
	assert(entry != NULL);

	entry->base = htobe32(entry->base);
	entry->size = htobe32(entry->size);
	entry->pid = htobe32(entry->pid);
	entry->id = htobe32(entry->id);
	entry->type = htobe32(entry->type);
	entry->flags = htobe32(entry->flags);
	entry->actual = htobe32(entry->actual);
	entry->checksum = htobe32(entry->checksum);

	entry->resvd[0] = htobe32(entry->resvd[0]);
	entry->resvd[1] = htobe32(entry->resvd[1]);
	entry->resvd[2] = htobe32(entry->resvd[2]);
	entry->resvd[3] = htobe32(entry->resvd[3]);

	for (int j = 0; j < FFS_USER_WORDS; j++)
		entry->user.data[j] = htobe32(entry->user.data[j]);
}

static int __hdr_read(ffs_hdr_t * hdr, FILE * file, off_t offset)
{
	assert(hdr != NULL);

	if (fseeko(file, offset, SEEK_SET) != 0) {
		ERRNO(errno);
		return -1;
	}

	size_t rc = fread(hdr, 1, sizeof(*hdr), file);
	if (rc <= 0 && ferror(file)) {
		ERRNO(errno);
		return -1;
	}

	uint32_t ck = memcpy_checksum(NULL, (void *)hdr,
				      offsetof(ffs_hdr_t, checksum));

	__hdr_be32toh(hdr);

	if (hdr->magic != FFS_MAGIC) {
		ERROR(ERR_UNEXPECTED, FFS_CHECK_HEADER_MAGIC,
		      "magic number mismatch '%x' != '%x'",
		      hdr->magic, FFS_MAGIC);
		return -1;
	}

	if (hdr->checksum != ck) {
		ERROR(ERR_UNEXPECTED, FFS_CHECK_HEADER_CHECKSUM,
		      "header checksum mismatch '%x' != '%x'",
		      hdr->checksum, ck);
		return -1;
	}

	return 0;
}

static int __hdr_write(ffs_hdr_t * hdr, FILE * file, off_t offset)
{
	assert(hdr != NULL);
	assert(hdr->magic == FFS_MAGIC);

	hdr->checksum = memcpy_checksum(NULL, (void *)hdr,
					offsetof(ffs_hdr_t, checksum));
	hdr->checksum = htobe32(hdr->checksum);

	if (fseeko(file, offset, SEEK_SET) != 0) {
		ERRNO(errno);
		return -1;
	}

	size_t size = sizeof(*hdr);

	if (0 < hdr->entry_count) {
		size += hdr->entry_count * hdr->entry_size;

		for (size_t i=0; i<hdr->entry_count; i++) {
			ffs_entry_t *e = hdr->entries + i;

			__entry_htobe32(e);

			e->checksum = memcpy_checksum(NULL, (void *)e,
						      offsetof(ffs_entry_t,
							       checksum));
			e->checksum = htobe32(e->checksum);
		}
	}

	__hdr_htobe32(hdr);

	size_t rc = fwrite(hdr, 1, size, file);
	if (rc <= 0 && ferror(file)) {
		ERRNO(errno);
		return -1;
	}

	__hdr_be32toh(hdr);
	if (0 < hdr->entry_count)
		for (size_t i=0; i<hdr->entry_count; i++)
			__entry_be32toh(hdr->entries + i);

	return 0;
}

static int __entries_read(ffs_hdr_t * hdr, FILE * file, off_t offset)
{
	assert(hdr != NULL);
	assert(hdr->magic == FFS_MAGIC);

	if (0 < hdr->entry_count) {
		if (fseeko(file, offset, SEEK_SET) != 0) {
			ERRNO(errno);
			return -1;
		}

		size_t size = hdr->entry_count * hdr->entry_size;

		size_t rc = fread(hdr->entries, 1, size, file);
		if (rc <= 0 && ferror(file)) {
			ERRNO(errno);
			return -1;
		}

		for (size_t i=0; i<hdr->entry_count; i++) {
			ffs_entry_t *e = hdr->entries + i;

			uint32_t ck = memcpy_checksum(NULL, (void *)e,
                                                      offsetof(ffs_entry_t,
							      checksum));
			__entry_be32toh(e);

			if (e->checksum != ck) {
				ERROR(ERR_UNEXPECTED, FFS_CHECK_ENTRY_CHECKSUM,
				      "'%s' entry checksum mismatch '%x' != "
				      "'%x'", e->name, e->checksum, ck);
				return -1;
			}
		}
	}

	return 0;
}

#if 0
static void __entries_write(ffs_hdr_t * hdr, FILE * file, off_t offset)
{
	if (hdr == NULL)
		ffs_throw(UNEX, 10400, "NULL hdr pointer");
	if (hdr->magic != FFS_MAGIC)
		ffs_throw(UNEX, 10401, "magic number mismatch '%x' != "
			  "'%x'", hdr->magic, FFS_MAGIC);

	if (0 < hdr->entry_count) {
		size_t size = hdr->entry_count * hdr->entry_size;

		for (size_t i = 0; i < hdr->entry_count; i++) {
			ffs_entry_t *e = hdr->entries + i;
			__entry_htobe32(e);
			e->checksum = memcpy_checksum(NULL, (void *)e,
						      offsetof(ffs_entry_t,
							       checksum));
			e->checksum = htobe32(e->checksum);
		}

		if (fseeko(file, offset, SEEK_SET) != 0)
			ffs_throw(ERR, 10402, "%s (errno=%d)",
				  strerror(errno), errno);

		size_t rc = fwrite(hdr->entries, 1, size, file);
		if (rc <= 0 && ferror(file))
			ffs_throw(ERR, 10403, "%s (errno=%d)",
			  	  strerror(errno), errno);

		fflush(file);

		for (size_t i = 0; i < hdr->entry_count; i++)
			__entry_be32toh(hdr->entries + i);
	}
}
#endif

static ffs_entry_t *__iterate_entries(ffs_hdr_t * self,
				      int (*func) (ffs_entry_t *))
{
	assert(self != NULL);
	assert(func != NULL);

	for (uint32_t i = 0; i < self->entry_count; i++) {
		if (func(self->entries + i) != 0)
			return self->entries + i;
	}

	return NULL;
}

static ffs_entry_t *__find_entry(ffs_hdr_t * self, const char *path)
{
	assert(self != NULL);

	if (path == NULL || *path == '\0')
		return NULL;

	if (*path == '/')
		path++;

	char __path[strlen(path) + 1];
	strcpy(__path, path);
	path = __path;

	ffs_entry_t root = {.id = FFS_PID_TOPLEVEL }, *parent = &root;

	char *name;
	while (parent != NULL && (name = strtok((char *)path, "/")) != NULL) {
		path = NULL;

		int find_entry(ffs_entry_t * child) {
			return (parent->id == child->pid &&
				strncmp(name, child->name,
					sizeof(child->name)) == 0);
		}

		parent = __iterate_entries(self, find_entry);
	}

	return parent;
}

/* ============================================================ */

int __ffs_fcheck(FILE *file, off_t offset)
{
	assert(file != NULL);

	RAII(ffs_hdr_t*, hdr, malloc(sizeof(*hdr)), free);
	if (hdr == NULL) {
		ERRNO(errno);
		return -1;
	}
	memset(hdr, 0, sizeof(*hdr));

	if (fseeko(file, offset, SEEK_SET) != 0) {
		ERRNO(errno);
		return -1;
	}

	if (fread(hdr, 1, sizeof(*hdr), file) <= 0 && ferror(file)) {
		ERRNO(errno);
		return -1;
	}

	uint32_t ck = memcpy_checksum(NULL, (void *)hdr,
				      offsetof(ffs_hdr_t, checksum));
	__hdr_be32toh(hdr);

	if (hdr->magic != FFS_MAGIC) {
		ERROR(ERR_UNEXPECTED, FFS_CHECK_HEADER_MAGIC,
		      "header magic mismatch '%x' != '%x'",
		      hdr->magic, FFS_MAGIC);
		return FFS_CHECK_HEADER_MAGIC;
	}
	if (hdr->checksum != ck) {
		ERROR(ERR_UNEXPECTED, FFS_CHECK_HEADER_CHECKSUM,
		      "header checksum mismatch '%x' != '%x'",
		      hdr->checksum, ck);
	 	return FFS_CHECK_HEADER_CHECKSUM;
	}

	size_t size = hdr->entry_count * hdr->entry_size;

	hdr = (ffs_hdr_t *)realloc(hdr, sizeof(*hdr) + size);
	if (hdr == NULL) {
		ERRNO(errno);
		return -1;
	}
	memset(hdr->entries, 0, size);

	if (0 < hdr->entry_count) {
		if (fread(hdr->entries, 1, size, file) <= 0 && ferror(file)) {
			ERRNO(errno);
			return -1;
		}

		for (size_t i = 0; i < hdr->entry_count; i++) {
			ffs_entry_t *e = hdr->entries + i;

			uint32_t ck = memcpy_checksum(NULL, (void *)e,
						      offsetof(ffs_entry_t,
							      checksum));

			__entry_be32toh(e);

			if (e->checksum != ck) {
				ERROR(ERR_UNEXPECTED, FFS_CHECK_ENTRY_CHECKSUM,
				      "'%s' entry checksum mismatch '%x' != "
				      "'%x'", e->name, e->checksum, ck);
				return FFS_CHECK_ENTRY_CHECKSUM;
			}
		}
	}

	return 0;
}

int __ffs_check(const char *path, off_t offset)
{
	if (path == NULL || *path == '\0') {
		UNEXPECTED("invalid path '%s'\n", path);
		return -1;
	}

	RAII(FILE*, file, fopen(path, "r"), fclose);
	if (file == NULL) {
		ERRNO(errno);
		return -1;
	}

	return __ffs_fcheck(file, offset);
}

ffs_t *__ffs_fcreate(FILE *file, off_t offset, uint32_t block_size,
		     uint32_t block_count)
{
	assert(file != NULL);

	if (!is_pow2(block_size)) {
		UNEXPECTED("'%d' invalid block size (must be non-0 and a "
			   "power of 2)", block_size);
		return NULL;
	}
	if (!is_pow2(block_count)) {
		UNEXPECTED("'%d' invalid block count (must be non-0 and a "
			   "power of 2)", block_count);
		return NULL;
	}
	if (offset & (block_size - 1)) {
		UNEXPECTED("'%lld' invalid offset (must be 'block size' "
			   "aligned)", (long long)offset);
		return NULL;
	}

	ffs_t *self = (ffs_t *) malloc(sizeof(*self));
	if (self == NULL) {
		ERRNO(errno);
		goto error;
	}

	memset(self, 0, sizeof(*self));
	self->file = file;
	self->offset = offset;
	self->count = FFS_ENTRY_EXTENT;
	self->dirty = true;

	self->hdr = (ffs_hdr_t *) malloc(sizeof(*self->hdr));
	if (self->hdr == NULL) {
		ERRNO(errno);
		goto error;
	}

	self->hdr->magic = FFS_MAGIC;
	self->hdr->version = FFS_VERSION_1;
	self->hdr->size = 1;
	self->hdr->entry_size = sizeof(*self->hdr->entries);
	self->hdr->entry_count = 0;
	self->hdr->block_size = block_size;
	self->hdr->block_count = block_count;
	self->hdr->checksum = 0;
	self->hdr->resvd[0] = 0;
	self->hdr->resvd[1] = 0;
	self->hdr->resvd[2] = 0;
	self->hdr->resvd[3] = 0;

	size_t size = self->count * self->hdr->entry_size;

	self->hdr = (ffs_hdr_t *) realloc(self->hdr, sizeof(*self->hdr) + size);
	if (self->hdr == NULL) {
		ERRNO(errno);
		goto error;
	}
	memset(self->hdr->entries, 0, size);

	if (__ffs_entry_add(self, FFS_PARTITION_NAME, offset, block_size,
			    FFS_TYPE_PARTITION, FFS_FLAGS_PROTECTED) < 0)
		goto error;
	if (__ffs_entry_truncate(self, FFS_PARTITION_NAME, block_size) < 0)
		goto error;

	if (false) {
 error:
		if (self != NULL) {
			if (self->file != NULL)
				fclose(self->file), self->file = NULL;
			if (self->path != NULL)
				free(self->path), self->path = NULL;
			if (self->hdr != NULL)
				free(self->hdr), self->hdr = NULL;
			free(self), self = NULL;
		}
	}

	return self;
}

ffs_t *__ffs_create(const char *path, off_t offset, uint32_t block_size,
		    uint32_t block_count)
{
	assert(path != NULL);

	FILE * file = fopen(path, "r+");
	if (file == NULL) {
		ERRNO(errno);
		return NULL;
	}

	ffs_t * self = __ffs_fcreate(file, offset, block_size, block_count);
	if (self != NULL)
		self->path = strdup(path);

	return self;
}

ffs_t *__ffs_fopen(FILE * file, off_t offset)
{
	assert(file != NULL);

	ffs_t *self = (ffs_t *) malloc(sizeof(*self));
	if (self == NULL) {
		ERRNO(errno);
		goto error;
	}

	memset(self, 0, sizeof(*self));
	self->file = file;
	self->count = 0;
	self->offset = offset;
	self->dirty = false;

	self->hdr = (ffs_hdr_t *) malloc(sizeof(*self->hdr));
	if (self->hdr == NULL) {
		ERRNO(errno);
		goto error;
	}
	memset(self->hdr, 0, sizeof(*self->hdr));

	if (__hdr_read(self->hdr, self->file, self->offset) < 0)
		goto error;

	self->count = max(self->hdr->entry_count, FFS_ENTRY_EXTENT);
	size_t size = self->count * self->hdr->entry_size;

	self->hdr = (ffs_hdr_t *)realloc(self->hdr, sizeof(*self->hdr) + size);
	if (self->hdr == NULL) {
		ERRNO(errno);
		goto error;
	}
	memset(self->hdr->entries, 0, size);

	if (0 < self->hdr->entry_count) {
		if (__entries_read(self->hdr, self->file,
	 		           self->offset + sizeof(*self->hdr)) < 0)
			goto error;
	}

	if (false) {
 error:
		if (self != NULL) {
			if (self->hdr != NULL)
				free(self->hdr), self->hdr = NULL;

			free(self), self = NULL;
		}
	}

	return self;
}

ffs_t *__ffs_open(const char *path, off_t offset)
{
	assert(path != NULL);

	FILE *file = fopen(path, "r+");
	if (file == NULL) {
		ERRNO(errno);
		return NULL;
	}

	ffs_t *self = __ffs_fopen(file, offset);
	if (self != NULL)
		self->path = strdup(path);

	return self;
}

static int ffs_flush(ffs_t * self)
{
	assert(self != NULL);

	if (__hdr_write(self->hdr, self->file, self->offset) < 0)
		return -1;

	if (fflush(self->file) != 0) {
		ERRNO(errno);
		return -1;
	}

	self->dirty = false;

	return 0;
}

int __ffs_info(ffs_t * self, int name, uint32_t *value)
{
	assert(self != NULL);
	assert(value != NULL);

	switch (name) {
	case FFS_INFO_MAGIC:
		*value = self->hdr->magic;
		break;
	case FFS_INFO_VERSION:
		*value = self->hdr->version;
		break;
	case FFS_INFO_ENTRY_SIZE:
		*value = self->hdr->entry_size;
		break;
	case FFS_INFO_ENTRY_COUNT:
		*value = self->hdr->entry_count;
		break;
	case FFS_INFO_BLOCK_SIZE:
		*value = self->hdr->block_size;
		break;
	case FFS_INFO_BLOCK_COUNT:
		*value = self->hdr->block_count;
		break;
	case FFS_INFO_OFFSET:
		*value = self->offset;
		break;
	default:
		UNEXPECTED("'%d' invalid info field", name);
		return -1;
	}

	return 0;
}


int __ffs_fclose(ffs_t * self)
{
	if (self == NULL)
		return 0;

	if (self->dirty == true)
		if (ffs_flush(self) < 0)
			return -1;

	if (self->hdr != NULL)
		free(self->hdr), self->hdr = NULL;

	memset(self, 0, sizeof(*self));
	free(self);

	return 0;
}

int __ffs_close(ffs_t * self)
{
	if (self == NULL)
		return 0;

	if (self->dirty == true)
		if (ffs_flush(self) < 0)
			return -1;

	if (self->path != NULL)
		free(self->path), self->path = NULL;
	if (self->file != NULL)
		fclose(self->file), self->file = NULL;

	return __ffs_fclose(self);
}

int __ffs_fsync(ffs_t * self)
{
	assert(self != NULL);

	if (fflush(self->file) < 0) {
		ERRNO(errno);
		return -1;
	}

	return 0;
}

static ffs_entry_t *__add_entry_check(ffs_hdr_t * self, off_t offset,
				      size_t size)
{
	assert(self != NULL);

	int find_overlap(ffs_entry_t * entry) {
		if (entry->type == FFS_TYPE_LOGICAL)
			return 0;

		off_t entry_start = entry->base;
		off_t entry_end = entry_start + entry->size - 1;

		off_t new_start = offset / self->block_size;
		off_t new_end = new_start + (size / self->block_size) - 1;

		return !(new_start < entry_start && new_end < entry_start) &&
  		       !(entry_end < new_start && entry_end < new_end);
	}

	return __iterate_entries(self, find_overlap);
}

int __ffs_iterate_entries(ffs_t * self, int (*func) (ffs_entry_t *))
{
	return __iterate_entries(self->hdr, func) != NULL;
}

int __ffs_list_entries(ffs_t * self, const char * name, bool user, FILE * out)
{
	assert(self != NULL);

	if (out == NULL)
		out = stdout;

	char full_name[4096];
	regex_t rx;

	int print_entry(ffs_entry_t * entry)
	{
		uint32_t offset = entry->base * self->hdr->block_size;
		uint32_t size = entry->size * self->hdr->block_size;

                if (__ffs_entry_name(self, entry, full_name,
				     sizeof full_name) < 0)
			return -1;

		if (regexec(&rx, full_name, 0, NULL, 0) == REG_NOMATCH)
			return 0;

		fprintf(stdout, "%3d [%08x-%08x:%8x] "
			"[%c%c%c%c%c%c%c%c%c%c] %s\n",
			entry->id, offset, offset+size-1, entry->actual,
			entry->type == FFS_TYPE_LOGICAL ? 'l' : 'd',
	/* reserved */	'-', '-', '-', '-', '-', '-', '-',
			entry->flags & FFS_FLAGS_U_BOOT_ENV ? 'b' : '-',
			entry->flags & FFS_FLAGS_PROTECTED ? 'p' : '-',
			full_name);

		if (user == true) {
			for (int i=0; i<FFS_USER_WORDS; i++) {
				fprintf(stdout, "[%2d] %8x ", i,
					entry->user.data[i]);
				if ((i+1) % 4 == 0)
					fprintf(stdout, "\n");
			}
		}

		return 0;
	}

	if (0 < self->count) {
		if (regcomp(&rx, name, REG_ICASE | REG_NOSUB) != 0) {
			ERRNO(errno);
			return-1;
		}

		fprintf(out, "========================[ PARTITION TABLE 0x%llx "
			"]=======================\n", (long long)self->offset);
		fprintf(out, "vers:%04x size:%04x * blk:%06x blk(s):%06x * "
			"entsz:%06x ent(s):%06x\n",
			self->hdr->version, self->hdr->size,
			self->hdr->block_size, self->hdr->block_count,
			self->hdr->entry_size, self->hdr->entry_count);
		fprintf(out, "------------------------------------------------"
			"---------------------------\n");

		(void)__iterate_entries(self->hdr, print_entry);

		fprintf(stdout, "\n");

		regfree(&rx);
	}

	return 0;
}

int __ffs_entry_find(ffs_t *self, const char *path, ffs_entry_t *entry)
{
	assert(self != NULL);
	assert(path != NULL);

	ffs_entry_t *__entry = __find_entry(self->hdr, path);
	if (__entry != NULL && entry != NULL)
		*entry = *__entry;

	return __entry != NULL;
}

int __ffs_entry_find_parent(ffs_t *self, const char *path, ffs_entry_t *entry)
{
	assert(self != NULL);
	assert(path != NULL);

	if (*path == '/')
		path++;

	char __path[strlen(path) + 1];
	strcpy(__path, path);
	char *parent_path = dirname(__path);

	int found = 0;

	if (strcmp(parent_path, ".") != 0) {
		ffs_entry_t parent;

		found = __ffs_entry_find(self, parent_path, &parent);

		if (found && entry != NULL)
			*entry = parent;
	}

	return found;
}

int __ffs_entry_name(ffs_t *self, ffs_entry_t *entry, char *name, size_t size)
{
	assert(self != NULL);
	assert(entry != NULL);

	ffs_hdr_t *hdr = self->hdr;

	int __entry_name(ffs_entry_t *parent, char *name, size_t size) {
		assert(parent != NULL);
		assert(name != NULL);

		if (parent->pid != FFS_PID_TOPLEVEL) {
			for (uint32_t i = 0; i < hdr->entry_count; i++) {
				if (hdr->entries[i].id == parent->pid) {
					__entry_name(hdr->entries + i, name,
						     size);
					break;
				}
			}
		}

		if (strlen(name) + strlen(parent->name) < size)
			strcat(name, parent->name);

		if (parent->id != entry->id) {
			if (strlen(name) + strlen("/") < size)
				strcat(name, "/");
		}

		return 0;
	}

	memset(name, 0, size);

	return __entry_name(entry, name, size);
}

int __ffs_entry_add(ffs_t * self, const char *path, off_t offset, uint32_t size,
		    ffs_type_t type, uint32_t flags)
{
	assert(self != NULL);
	assert(path != NULL);

	if (__ffs_entry_find(self, path, NULL) == true) {
		UNEXPECTED("'%s' entry already exists", path);
		return -1;
	}

	ffs_entry_t parent = {.id = FFS_PID_TOPLEVEL };
	(void)__ffs_entry_find_parent(self, path, &parent);

	ffs_hdr_t *hdr = self->hdr;

	if (type != FFS_TYPE_LOGICAL) {
		ffs_entry_t *overlap = __add_entry_check(hdr, offset, size);
		if (overlap != NULL) {
			UNEXPECTED("'%s' at offset %lld and size %d overlaps "
				   "'%s' at offset %d and size %d",
				   path, (long long)offset, size, overlap->name,
				   overlap->base * hdr->block_size,
				   overlap->size * hdr->block_size);
			return -1;
		}
	}

	int find_empty(ffs_entry_t * empty) {
		return empty->type == 0;
	}

	ffs_entry_t *entry = __iterate_entries(hdr, find_empty);
	if (entry == NULL) {
		if (self->count <= hdr->entry_count) {
			size_t new_size;
			new_size = hdr->entry_size *
					(self->count + FFS_ENTRY_EXTENT);

			self->hdr = (ffs_hdr_t *) realloc(self->hdr,
							  sizeof(*self->hdr) +
							  new_size);
			assert(self->hdr != NULL);

			if (hdr != self->hdr)
				hdr = self->hdr;

			memset(hdr->entries + self->count, 0,
			       FFS_ENTRY_EXTENT * hdr->entry_size);

			self->count += FFS_ENTRY_EXTENT;
		}

		entry = hdr->entries + hdr->entry_count;
	}

	uint32_t max_id = 0;

	int find_max_id(ffs_entry_t * max) {
		if (max_id < max->id)
			max_id = max->id;
		return 0;
	}

	(void)__iterate_entries(hdr, find_max_id);

	char name[strlen(path) + 1];
	strcpy(name, path);
	strncpy(entry->name, basename(name), sizeof(entry->name));
	entry->id = max_id + 1;
	entry->pid = parent.id;
	entry->base = offset / hdr->block_size;
	entry->size = size / hdr->block_size;
	entry->type = type;
	entry->flags = flags;
	entry->checksum = 0;

	hdr->entry_count++;

    // Need to update 'part' entry as well as ffs hdr
    // if the required number of blocks changes
    uint32_t blocksNeeded = (hdr->entry_count * hdr->entry_size + FFS_HDR_SIZE_NO_ENTRY) / hdr->block_size;
    if(hdr->entry_count * hdr->entry_size + FFS_HDR_SIZE_NO_ENTRY % hdr->block_size)
    {
        blocksNeeded++;
    }

    if(hdr->size != blocksNeeded)
    {
        ffs_entry_t entry;
        if (__ffs_entry_find(self, "part", &entry) == false) {
            UNEXPECTED("entry '%s' not found in table at offset '%llx'",
                       "part", (long long)self->offset);
                       return -1;
        }

        int find_entry_id(ffs_entry_t * __entry) {
              return entry.id == __entry->id;
        }

        ffs_entry_t *entry_p = __iterate_entries(hdr, find_entry_id);
        assert(entry_p != NULL);

        hdr->size = blocksNeeded;
        entry_p->size = blocksNeeded;
        entry_p->actual = blocksNeeded * hdr->block_size;
    }
	self->dirty = true;

	return 0;
}

int __ffs_entry_delete(ffs_t * self, const char *path)
{
	assert(self != NULL);
	assert(path != NULL);

	ffs_entry_t entry;
	if (__ffs_entry_find(self, path, &entry) == false) {
		UNEXPECTED("entry '%s' not found in table at offset '%llx'",
			  path, (long long)self->offset);
		return -1;
	}

	if (entry.type == FFS_TYPE_PARTITION) {
		UNEXPECTED("'%s' cannot --delete partition type entries", path);
		return -1;
	}

	uint32_t children = 0;

	int find_children(ffs_entry_t * child) {
		if (entry.id == child->pid)
			children++;
		return 0;
	}

	ffs_hdr_t *hdr = self->hdr;

	(void)__iterate_entries(hdr, find_children);

	if (0 < children) {
		UNEXPECTED("'%s' has '%d' children, --delete those first",
			   path, children);
		return -1;
	}

	int find_entry_id(ffs_entry_t * __entry) {
		return entry.id == __entry->id;
	}

	ffs_entry_t *entry_p = __iterate_entries(hdr, find_entry_id);
	assert(entry_p != NULL);

	int start = entry_p - hdr->entries;
	int count = hdr->entry_count - start;

	memmove(entry_p, entry_p + 1, hdr->entry_size * count);

	hdr->entry_count = max(0UL, hdr->entry_count - 1);
	memset(hdr->entries + hdr->entry_count, 0, hdr->entry_size);

	self->dirty = true;

	return 0;
}

int __ffs_entry_user_get(ffs_t *self, const char *path, uint32_t word,
			 uint32_t *value)
{
	assert(self != NULL);
	assert(path != NULL);
	assert(value != NULL);

	if (FFS_USER_WORDS <= word) {
		UNEXPECTED("word '%d' outside range [0..%d]",
			   word, FFS_USER_WORDS - 1);
		return -1;
	}

	ffs_entry_t *entry = __find_entry(self->hdr, path);
	if (entry == NULL) {
		UNEXPECTED("entry '%s' not found in partition table at "
			   "offset '%llx'", path, (long long)self->offset);
		return -1;
	}

	*value = entry->user.data[word];

	return 0;
}

int __ffs_entry_user_put(ffs_t *self, const char *path, uint32_t word,
			 uint32_t value)
{
	assert(self != NULL);
	assert(path != NULL);

	if (FFS_USER_WORDS <= word) {
		UNEXPECTED("word '%d' outside range [0..%d]",
			   word, FFS_USER_WORDS - 1);
		return -1;
	}

	ffs_entry_t *entry = __find_entry(self->hdr, path);
	if (entry == NULL) {
		UNEXPECTED("entry '%s' not found in partition table at "
			   "offset '%llx'", path, (long long)self->offset);
		return -1;
	}

	entry->user.data[word] = value;
	self->dirty = true;

	return 0;
}

ssize_t __ffs_entry_hexdump(ffs_t * self, const char *path, FILE * out)
{
	assert(self != NULL);
	assert(path != NULL);

	ffs_entry_t entry;
	if (__ffs_entry_find(self, path, &entry) == false) {
		UNEXPECTED("entry '%s' not found in table at offset '%llx'",
			   path, (long long)self->offset);
		return -1;
	}

	size_t size = entry.size * self->hdr->block_size;
	if (entry.actual < size)
		size = entry.actual;

	off_t offset = entry.base * self->hdr->block_size;

	if (fseeko(self->file, offset, SEEK_SET) != 0) {
		ERRNO(errno);
		return -1;
	}

	ssize_t total = 0;

	size_t block_size = self->hdr->block_size;
	char block[block_size];
	while (0 < size) {
		size_t rc = fread(block, 1, min(block_size, size),
				  self->file);
		if (rc <= 0) {
			if (ferror(self->file)) {
				ERRNO(errno);
				return -1;
			}
			break;
		}

		dump_memory(out, offset + total, block, rc);

		total += rc;
		size -= rc;
	}

	return total;
}

ssize_t __ffs_entry_truncate(ffs_t * self, const char *path, size_t size)
{
	assert(self != NULL);
	assert(path != NULL);

	ffs_entry_t * entry = __find_entry(self->hdr, path);
	if (entry == NULL) {
		UNEXPECTED("entry '%s' not found in partition table at "
			   "offset '%llx'", path, (long long)self->offset);
		return -1;
	}

	if ((entry->size * self->hdr->block_size) < size) {
		errno = EFBIG;
		ERRNO(errno);
		return -1;
	} else {
		entry->actual = size;
		self->dirty = true;
	}

	return 0;
}

ssize_t __ffs_entry_read(ffs_t * self, const char *path, void *buf,
			 off_t offset, size_t count)
{
	assert(self != NULL);
	assert(path != NULL);
	assert(buf != NULL);

	if (count == 0)
		return 0;

	ffs_entry_t entry;
	if (__ffs_entry_find(self, path, &entry) == false) {
		UNEXPECTED("entry '%s' not found in partition table at "
			   "offset '%llx'", path, (long long)self->offset);
		return -1;
	}

	size_t entry_size = entry.size * self->hdr->block_size;
	if (entry.actual < entry_size)
		entry_size = entry.actual;
	off_t entry_offset = entry.base * self->hdr->block_size;

	if (entry_size <= offset)
		return 0;
	else
		count = min(count, (entry_offset + entry_size) - offset);

	ssize_t total = 0;

	if (fseeko(self->file, entry_offset, SEEK_SET) != 0) {
		ERRNO(errno);
		return -1;
	}

	if (fseeko(self->file, offset, SEEK_CUR) != 0) {
		ERRNO(errno);
		return -1;
	}

	while (0 < count) {
		size_t rc = fread(buf + total, 1, count, self->file);
		if (rc <= 0) {
			if (ferror(self->file)) {
				ERRNO(errno);
				return -1;
			}
			break;
		}

		total += rc;
		count -= rc;
	}

	return total;
}

ssize_t __ffs_entry_write(ffs_t * self, const char *path, const void *buf,
			  off_t offset, size_t count)
{
	assert(self != NULL);
	assert(path != NULL);
	assert(buf != NULL);

	if (count == 0)
		return 0;

	ffs_entry_t *entry = __find_entry(self->hdr, path);
	if (entry == NULL) {
		UNEXPECTED("entry '%s' not found in partition table at "
			   "offset '%llx'", path, (long long)self->offset);
		return -1;
	}

	size_t entry_size = entry->size * self->hdr->block_size;
	off_t entry_offset = entry->base * self->hdr->block_size;

	if (entry_size <= offset)
		return 0;
	else
		count = min(count, (entry_offset + entry_size) - offset);

	ssize_t total = 0;

	if (fseeko(self->file, entry_offset, SEEK_SET) != 0) {
		ERRNO(errno);
		return -1;
	}

	if (fseeko(self->file, offset, SEEK_CUR) != 0) {
		ERRNO(errno);
		return -1;
	}

	while (0 < count) {
		size_t rc = fwrite(buf + total, 1, count, self->file);
		if (rc <= 0) {
			if (ferror(self->file)) {
				ERRNO(errno);
				return -1;
			}
			break;
		}
		total += rc;
		count -= rc;
	}

	fflush(self->file);

	if (entry->actual < (uint32_t) total) {
		entry->actual = (uint32_t) total;
		self->dirty = true;
	}

	return total;
}

#if 0
ssize_t __ffs_entry_copy(ffs_t *self, ffs_t *in, const char *path)
{
	assert(self != NULL);
	assert(in != NULL);
	assert(path != NULL);

	if (*path == '\0')
		return 0;

	ffs_entry_t *src = __find_entry(in->hdr, path);
	if (src == NULL) {
		UNEXPECTED("entry '%s' not found in table at offset '%llx'",
			   path, in->offset);
		return -1;
	}

	ffs_entry_t *dest = __find_entry(self->hdr, path);
	if (dest == NULL) {
		UNEXPECTED("entry '%s' not found in table at offset '%llx'",
			   path, self->offset);
		return -1;
	}

	if (src->base != dest->base) {
		UNEXPECTED("partition '%s' offsets differ '%x' != '%x'",
			   path, src->base, dest->base);
		return -1;
	}

	if (src->size != dest->size) {
		UNEXPECTED("partition '%s' sizes differ '%x' != '%x'",
			   path, src->size, dest->size);
		return -1;
	}

	size_t block_size = self->hdr->block_size;
	off_t src_offset = src->base * in->hdr->block_size;
	size_t src_actual = src->actual;
	off_t dest_offset = dest->base * self->hdr->block_size;

	if (fseeko(in->file, src_offset, SEEK_SET) != 0) {
		ERRNO(errno);
		return -1;
	}

	if (fseeko(self->file, dest_offset, SEEK_SET) != 0) {
		ERRNO(errno);
		return -1;
	}

	ssize_t total = 0;
	char block[block_size];

	while (0 < src_actual) {
		size_t count =  min(src_actual, block_size);

		size_t rc = fread(block, 1, count, in->file);
		if (rc <= 0 && ferror(in->file)) {
			ERRNO(errno);
			return -1;
		}

		rc = fwrite(block, 1, rc, self->file);
		if (rc <= 0 && ferror(self->file)) {
			ERRNO(errno);
			return -1;
		}

		total += rc;
		src_actual -= rc;
	}

	if (dest->actual != (uint32_t)total) {
		dest->actual = (uint32_t) total;
		self->dirty = true;
	}

	return total;
}

ssize_t __ffs_entry_compare(ffs_t * self, ffs_t * in, const char *path)
{
	assert(self != NULL);
	assert(in != NULL);
	assert(path != NULL);

	if (*path == '\0')
		return 0;

	ffs_entry_t *src = __find_entry(in->hdr, path);
	if (src == NULL) {
		UNEXPECTED("entry '%s' not found in table at offset '%llx'",
			   path, in->offset);
		return -1;
	}

	ffs_entry_t *dest = __find_entry(self->hdr, path);
	if (dest == NULL) {
		UNEXPECTED("entry '%s' not found in table at offset '%llx'",
			   path, self->offset);
		return -1;
	}

	if (src->base != dest->base) {
		UNEXPECTED("partition '%s' offsets differ '%x' != '%x'",
			   path, src->base, dest->base);
		return -1;
	}

	if (src->size != dest->size) {
		UNEXPECTED("partition '%s' sizes differ '%x' != '%x'",
			   path, src->size, dest->size);
		return -1;
	}

	if (src->actual != dest->actual) {
		UNEXPECTED("partition '%s' actual sizes differ '%x' != '%x'",
			   path, src->actual, dest->actual);
		return -1;
	}

	off_t offset = src->base * self->hdr->block_size;
	size_t actual = src->actual;

	if (fseeko(in->file, offset, SEEK_SET) != 0) {
		ERRNO(errno);
		return -1;
	}

	if (fseeko(self->file, offset, SEEK_SET) != 0) {
		ERRNO(errno);
		return -1;
	}

	ssize_t total = 0;
	size_t size = 256;
	char __src[size], __dest[size];

	while (0 < actual) {
		size = min(actual, size);

		size_t rc = fread(__src, 1, size, in->file);
		if (rc <= 0 && ferror(in->file)) {
			ERRNO(errno);
			return -1;
		}

		rc = fread(__dest, 1, size, self->file);
		if (rc <= 0 && ferror(self->file)) {
			ERRNO(errno);
			return -1;
		}

		if (memcmp(__src, __dest, size) != 0) {
			printf("==========> '%s'\n", self->path);
			dump_memory(stdout, offset + total, __dest, size);

			printf("==========> '%s'\n", in->path);
			dump_memory(stdout, offset + total, __src, size);

			break;
		}

		actual -= size;
		total += size;
	}

	return total;
}
#endif

int __ffs_entry_list(ffs_t * self, ffs_entry_t ** list)
{
	assert(self != NULL);
	assert(list != NULL);

	size_t size = 0, count = 0;
	*list = NULL;

	int name_list(ffs_entry_t * entry) {
		if (size <= count) {
			size += 10;

			*list = realloc(*list, size * sizeof(**list));
			if (*list == NULL) {
				ERRNO(errno);
				return -1;
			}
		}

		(*list)[count++] = *entry;

		return 0;
	}

	if (__ffs_iterate_entries(self, name_list) != 0) {
		if (*list != NULL)
			free(*list), *list = NULL;
		count = 0;
		return -1;
	}

	return count;
}

/* ============================================================ */
