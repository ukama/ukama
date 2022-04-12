/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fcp/src/misc.c $                                              */
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
 *    File: fcp_misc.c
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: misc helpers
 *    Date: 01/30/2013
 */

#include <sys/types.h>
#include <sys/stat.h>

#include <fcntl.h>
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include <unistd.h>
#include <getopt.h>
#include <errno.h>
#include <ctype.h>
#include <regex.h>

#include <clib/attribute.h>
#include <clib/version.h>
#include <clib/list.h>
#include <clib/list_iter.h>
#include <clib/misc.h>
#include <clib/min.h>
#include <clib/err.h>
#include <clib/raii.h>

#include "misc.h"
#include "main.h"

#define COMPARE_SIZE	256UL

static regex_t * regex_create(const char * str)
{
	assert(str != NULL);

	regex_t * self = (regex_t *)malloc(sizeof(*self));
	if (self == NULL) {
		ERRNO(errno);
		return self;
	}

	if (regcomp(self, str, REG_ICASE | REG_NOSUB) != 0) {
		free(self);
		ERRNO(errno);
		return NULL;
	}

	return self;
}

static int regex_free(regex_t * self)
{
	if (self == NULL)
		return 0;

	regfree(self);
	free(self);

	return 0;
}

int parse_offset(const char *str, off_t *offset)
{
	assert(offset != NULL);

	if (str == NULL) {
		*offset = 0;
		return 0;
	}

	char *end = NULL;

	errno = 0;
	*offset = strtoull(str, &end, 0);
	if (errno != 0) {
		ERRNO(errno);
		return -1;
	}

	if (*end != '\0') {
		if (!strcmp(end, "KiB")    ||
		    !strcasecmp(end, "KB") ||
		    !strcasecmp(end, "K"))
			*offset <<= 10;
		else if (!strcmp(end, "MiB")    ||
			 !strcasecmp(end, "MB") ||
			 !strcasecmp(end, "M"))
			*offset <<= 20;
		else if (!strcmp(end, "GiB")    ||
		    	 !strcasecmp(end, "GB") ||
			 !strcasecmp(end, "G"))
			*offset <<= 30;
		else {
			UNEXPECTED("invalid offset specified '%s'", end);
			return -1;
		}
	}

	return 0;
}

int parse_size(const char *str, uint32_t *size)
{
	assert(size != NULL);

	if (str == NULL) {
		*size = 0;
		return 0;
	}

	char *end = NULL;

	errno = 0;
	*size = strtoul(str, &end, 0);
	if (errno != 0) {
		ERRNO(errno);
		return -1;
	}

	if (*end != '\0') {
		if (!strcmp(end, "KiB") || !strcasecmp(end, "K") ||
		    !strcasecmp(end, "KB"))
			*size <<= 10;
		else if (!strcmp(end, "MiB") || !strcasecmp(end, "M") ||
			 !strcasecmp(end, "MB"))
			*size <<= 20;
		else if (!strcmp(end, "GiB") || !strcasecmp(end, "G") ||
			 !strcasecmp(end, "GB"))
			*size <<= 30;
		else {
			UNEXPECTED("invalid size specified '%s'", end);
			return -1;
		}
	}

	return 0;
}

int parse_number(const char *str, uint32_t *num)
{
	assert(num != NULL);

	if (str == NULL) {
		*num = 0;
		return 0;
	}

	char *end = NULL;

	errno = 0;
	*num = strtoul(str, &end, 0);
	if (errno != 0) {
		ERRNO(errno);
		return -1;
	}

	if (*end != '\0') {
		UNEXPECTED("invalid number specified '%s'", end);
		return -1;
	}

	return 0;
}

int parse_path(const char * path, char ** type, char ** target, char ** name)
{
	assert(path != NULL);

	*type = *target = *name = NULL;

	char * delim1 = strchr(path, ':');
	char * delim2 = strrchr(path, ':');

	if (delim1 == NULL && delim2 == NULL) {	// <target>
		if (asprintf(target, "%s", path) < 0) {
			ERRNO(errno);
			return -1;
		}
	} else if (delim1 == delim2) {		// <target>:<name>
		if (asprintf(target, "%.*s", (uint32_t)(delim1 - path), path) < 0) {
			ERRNO(errno);
			return -1;
		}
		delim1++;
		if (asprintf(name, "%s", delim1) < 0) {
			ERRNO(errno);
			return -1;
		}

		if (valid_type(*target) == true) {
			*type = *target;
			*target = *name;
			*name = NULL;
		}
	} else if (delim1 != delim2) {		// <type>:<target>:<name>
		if (asprintf(type, "%.*s", (uint32_t)(delim1 - path), path) < 0) {
			ERRNO(errno);
			return -1;
		}
		delim1++;
		if (asprintf(target, "%.*s", (uint32_t)(delim2 - delim1), delim1) < 0) {
			ERRNO(errno);
			return -1;
		}
		delim2++;
		if (asprintf(name, "%s", delim2) < 0) {
			ERRNO(errno);
			return -1;
		}
	}

	return 0;
}

int dump_errors(const char * name, FILE * out)
{
	assert(name != NULL);

	if (out == NULL)
		out = stderr;

	err_t * err = NULL;

	while ((err = err_get()) != NULL) {
		switch (err_type(err)) {
		case ERR_VERSION:
			fprintf(out, "%s: %s : %s(%d) : v%d.%02d.%04d %.*s\n",
				basename((char *)name), err_type_name(err),
				basename(err_file(err)), err_line(err),
				VER_TO_MAJOR(err_code(err)),
				VER_TO_MINOR(err_code(err)),
				VER_TO_PATCH(err_code(err)),
				err_size(err), (char *)err_data(err));
			break;
		default:
			fprintf(out, "%s: %s : %s(%d) : (code=%d) %.*s\n",
				basename((char *)name), err_type_name(err),
				basename(err_file(err)), err_line(err),
				err_code(err), err_size(err),
				(char *)err_data(err));
		}
	}

	return 0;
}

int check_file(const char * path, FILE * file, off_t offset) {
	assert(file != NULL);

	switch (__ffs_fcheck(file, offset)) {
	case 0:
		return 0;
	case FFS_CHECK_HEADER_MAGIC:
		UNEXPECTED("'%s' no partition table found at offset '%llx'\n",
			   path, (long long)offset);
		return -1;
	case FFS_CHECK_HEADER_CHECKSUM:
		UNEXPECTED("'%s' partition table at offset '%llx', is "
			   "corrupted\n", path, (long long)offset);
		return -1;
	case FFS_CHECK_ENTRY_CHECKSUM:
		UNEXPECTED("'%s' partition table at offset '%llx', has "
			   "corrupted entries\n", path, (long long)offset);
		return -1;
	default:
		return -1;
	}
}

entry_list_t * entry_list_create(ffs_t * ffs)
{
	entry_list_t * self = (entry_list_t *)malloc(sizeof(*self));
	if (self == NULL) {
		ERRNO(errno);
		return NULL;
	}

	list_init(&self->list);

	self->ffs = ffs;

	return self;
}

entry_list_t * entry_list_create_by_regex(ffs_t * ffs, const char * name)
{
	assert(ffs != NULL);

	if (name == NULL)
		name = ".*";

	RAII(regex_t*, rx, regex_create(name), regex_free);
	if (rx == NULL)
		return NULL;

	entry_list_t * self = entry_list_create(ffs);
	if (self == NULL)
		return NULL;

	int entry_list(ffs_entry_t * entry)
	{
		assert(entry != NULL);

		char full_name[page_size];
		if (__ffs_entry_name(ffs, entry, full_name,
				     sizeof full_name) < 0)
			return -1;
		if (regexec(rx, full_name, 0, NULL, 0) == REG_NOMATCH)
			return 0;

		entry_node_t * entry_node;
		entry_node = (entry_node_t *)malloc(sizeof(*entry_node));
		if (entry_node == NULL) {
			ERRNO(errno);
			return -1;
		}

		memcpy(&entry_node->entry, entry, sizeof(*entry));
		list_add_tail(&self->list, &entry_node->node);

		return 0;
	}

	if (__ffs_iterate_entries(ffs, entry_list) < 0) {
		entry_list_delete(self);
		return NULL;
	}

	return self;
}

int entry_list_add(entry_list_t * self, ffs_entry_t * entry)
{
	assert(self != NULL);
	assert(entry != NULL);

	entry_node_t * entry_node;
	entry_node = (entry_node_t *)malloc(sizeof(*entry_node));
	if (entry_node == NULL) {
		ERRNO(errno);
		return -1;
	}

	memcpy(&entry_node->entry, entry, sizeof(entry_node->entry));
	list_add_tail(&self->list, &entry_node->node);

	return 0;
}

int entry_list_add_child(entry_list_t * self, ffs_entry_t * parent)
{
	assert(self != NULL);
	assert(parent != NULL);

	int child_entry_list(ffs_entry_t * child)
	{
		assert(child != NULL);

		if (child->pid != parent->id)
			return 0;

		entry_node_t * entry_node;
		entry_node = (entry_node_t *)malloc(sizeof(*entry_node));
		if (entry_node == NULL) {
			ERRNO(errno);
			return -1;
		}

		memcpy(&entry_node->entry, child, sizeof(*child));
		list_add_tail(&self->list, &entry_node->node);

		return 0;
	}

	if (__ffs_iterate_entries(self->ffs, child_entry_list) < 0)
		return -1;

	return 0;
}

int entry_list_remove(entry_list_t * self, entry_node_t * node)
{
	assert(self != NULL);
	assert(node != NULL);

	list_remove_node(&self->list, &node->node);

	return 0;
}

int entry_list_delete(entry_list_t * self)
{
	if (self == NULL)
		return 0;

	while (!list_empty(&self->list)) {
		free(container_of(list_remove_head(&self->list), entry_node_t,
				  node));
	}

	return 0;
}

int entry_list_exists(entry_list_t * self, ffs_entry_t * entry)
{
	assert(self != NULL);
	assert(entry != NULL);

	list_iter_t it;
	entry_node_t * entry_node;

	list_iter_init(&it, &self->list, LI_FLAG_FWD);

	list_for_each(&it, entry_node, node) {
		ffs_entry_t * __entry = &entry_node->entry;

		if (__entry->base == entry->base &&
		    __entry->size == entry->size)
			return 1;
	}

	return 0;
}

ffs_entry_t * entry_list_find(entry_list_t * self, const char * name)
{
	assert(self != NULL);
	assert(name != NULL);

	list_iter_t it;
	entry_node_t * entry_node = NULL;

	list_iter_init(&it, &self->list, LI_FLAG_FWD);

	list_for_each(&it, entry_node, node) {
		ffs_entry_t * entry = &entry_node->entry;

		if (strcmp(entry->name, name) == 0)
			return entry;
	}

	return NULL;
}

int entry_list_dump(entry_list_t * self, FILE * out)
{
	assert(self != NULL);
	if (out == NULL)
		out = stdout;

	list_iter_t it;
	entry_node_t * entry_node = NULL;

	list_iter_init(&it, &self->list, LI_FLAG_FWD);

	list_for_each(&it, entry_node, node) {
		ffs_entry_t * entry = &entry_node->entry;

		fprintf(stderr, "id[%d] pid[%d] name[%s]\n",
			entry->id, entry->pid, entry->name);
	}

	return 0;
}

FILE *__fopen(const char * type, const char * target, const char * mode,
	      int debug)
{
	assert(target != NULL);
	assert(mode != NULL);

	FILE *file = NULL;
	uint32_t port = 0;

	if (type == NULL)
		type = TYPE_FILE;

	if (strcasecmp(type, TYPE_AA) == 0) {
		if (parse_number(target, &port) < 0)
			return NULL;
                UNEXPECTED("Removed support");
		//file = fopen_aaflash(port, mode, debug);
	} else if (strcasecmp(type, TYPE_RW) == 0) {
                UNEXPECTED("Removed support");
		//file = fopen_rwflash(target, mode, debug);
	} else if (strcasecmp(type, TYPE_SFC) == 0) {
		UNEXPECTED("FIX ME");
		return NULL;
	} else if (strcasecmp(type, TYPE_FILE) == 0) {
		file = fopen(target, mode);
		if (file == NULL)
			ERRNO(errno);
	} else {
		errno = EINVAL;
		ERRNO(errno);
	}

	return file;
}

int is_file(const char * type, const char * target, const char * name)
{
	return type == NULL && target != NULL && name == NULL;
}

int valid_type(const char * type)
{
	return type == NULL ? 0 :
	       strcasecmp(type, TYPE_FILE) == 0 ||
	       strcasecmp(type, TYPE_RW) == 0   ||
	       strcasecmp(type, TYPE_AA) == 0   ||
	       strcasecmp(type, TYPE_SFC) == 0;
}

int fcp_read_entry(ffs_t * src, const char * name, FILE * out)
{
	assert(src != NULL);
	assert(name != NULL);

	uint32_t block_size;
	if (__ffs_info(src, FFS_INFO_BLOCK_SIZE, &block_size) < 0)
		return -1;

	uint32_t block_count;
	if (__ffs_info(src, FFS_INFO_BLOCK_COUNT, &block_count) < 0)
		return -1;

	size_t buffer_size = block_size * block_count;
	RAII(void*, buffer, malloc(buffer_size), free);
	if (buffer == NULL) {
		ERRNO(errno);
		return -1;
	}

	ffs_entry_t entry;
	if (__ffs_entry_find(src, name, &entry) == false) {
		UNEXPECTED("'%s' partition not found => %s",
			   src->path, name);
		return -1;
	}

	uint32_t poffset;
	if (__ffs_info(src, FFS_INFO_OFFSET, &poffset) < 0)
		return -1;

	uint32_t total = 0;
	uint32_t size = entry.actual;
	off_t offset = 0;

	if (isatty(fileno(stderr))) {
		fprintf(stderr, "%8x: %s: read partition %8x/%8x",
			poffset, name, entry.actual, total);
	}

	while (0 < size) {
		size_t count = min(buffer_size, size);

		ssize_t rc;
		rc = __ffs_entry_read(src, name, buffer, offset, count);


		rc = fwrite(buffer, 1, rc, out);
		if (rc <= 0 && ferror(out)) {
			ERRNO(errno);
			return -1;
		}

		size -= rc;
		total += rc;
		offset += rc;

		if (isatty(fileno(stderr))) {
			fprintf(stderr, "\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b");
			fprintf(stderr, "%8x/%8x", entry.actual, total);
		}
	}

	if (isatty(fileno(stderr))) {
		fprintf(stderr, "\n");
	}

	return total;
}

int fcp_write_entry(ffs_t * dst, const char * name, FILE * in)
{
	assert(dst != NULL);
	assert(name != NULL);

	uint32_t block_size;
	if (__ffs_info(dst, FFS_INFO_BLOCK_SIZE, &block_size) < 0)
		return -1;

	uint32_t block_count;
	if (__ffs_info(dst, FFS_INFO_BLOCK_COUNT, &block_count) < 0)
		return -1;

	size_t buffer_size = block_size * block_count;
	RAII(void*, buffer, malloc(buffer_size), free);
	if (buffer == NULL) {
		ERRNO(errno);
		return -1;
	}

	ffs_entry_t entry;
	if (__ffs_entry_find(dst, name, &entry) == false) {
		UNEXPECTED("'%s' partition not found => %s",
			   dst->path, name);
		return -1;
	}

	uint32_t poffset;
	if (__ffs_info(dst, FFS_INFO_OFFSET, &poffset) < 0)
		return -1;

	uint32_t total = 0;
	uint32_t size = entry.actual;
	off_t offset = 0;

	if (isatty(fileno(stderr))) {
		fprintf(stderr, "%8x: %s: write partition %8x/%8x",
			poffset, name, entry.actual, total);
	}

	while (0 < size) {
		size_t count = min(buffer_size, size);

		ssize_t rc;
		rc = fread(buffer, 1, count, in);
		if (rc <= 0) {
			if (feof(in)) {
				break;
			} else if (ferror(in)) {
				ERRNO(errno);
				return -1;
			}
		}

		rc = __ffs_entry_write(dst, name, buffer, offset, rc);

		if (__ffs_fsync(dst) < 0)
			return -1;

		size -= rc;
		total += rc;
		offset += rc;

		if (isatty(fileno(stderr))) {
			fprintf(stderr, "\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b");
			fprintf(stderr, "%8x/%8x", entry.actual, total);
		}
	}

	if (isatty(fileno(stderr))) {
		fprintf(stderr, "\n");
	}

	return total;
}

int fcp_erase_entry(ffs_t * dst, const char * name, char fill)
{
	assert(dst != NULL);
	assert(name != NULL);

	uint32_t block_size;
	if (__ffs_info(dst, FFS_INFO_BLOCK_SIZE, &block_size) < 0)
		return -1;

	uint32_t block_count;
	if (__ffs_info(dst, FFS_INFO_BLOCK_COUNT, &block_count) < 0)
		return -1;

	size_t buffer_size = block_size * block_count;
	RAII(void*, buffer, malloc(buffer_size), free);
	if (buffer == NULL) {
		ERRNO(errno);
		return -1;
	}

	memset(buffer, fill, buffer_size);

	ffs_entry_t entry;
	if (__ffs_entry_find(dst, name, &entry) == false) {
		UNEXPECTED("'%s' partition not found => %s",
			   dst->path, name);
		return -1;
	}

	uint32_t poffset;
	if (__ffs_info(dst, FFS_INFO_OFFSET, &poffset) < 0)
		return -1;

	uint32_t total = 0;
	uint32_t size = entry.size * block_size;
	off_t offset = 0;

	if (isatty(fileno(stderr))) {
		fprintf(stderr, "%8x: %s: erase partition %8x/%8x",
			poffset, name, entry.actual, total);
	}

	while (0 < size) {
		size_t count = min(buffer_size, size);

		ssize_t rc;
		rc = __ffs_entry_write(dst, name, buffer, offset, count);

		if (__ffs_fsync(dst) < 0)
			return -1;

		size -= rc;
		total += rc;
		offset += rc;

		if (isatty(fileno(stderr))) {
			fprintf(stderr, "\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b");
			fprintf(stderr, "%8x/%8x", entry.size * block_size,
				total);
		}
	}

	if (__ffs_entry_truncate(dst, name, 0ULL) < 0) {
		ERRNO(errno);
		return -1;
	}

	if (isatty(fileno(stderr))) {
		fprintf(stderr, "\n");
	}

	return total;
}

int fcp_copy_entry(ffs_t * src, const char * src_name,
		   ffs_t * dst, const char * dst_name)
{
	assert(src != NULL);
	assert(src_name != NULL);
	assert(dst != NULL);
	assert(dst_name != NULL);

	uint32_t block_size;
	if (__ffs_info(src, FFS_INFO_BLOCK_SIZE, &block_size) < 0)
		return -1;

	uint32_t block_count;
	if (__ffs_info(dst, FFS_INFO_BLOCK_COUNT, &block_count) < 0)
		return -1;

	size_t buffer_size = block_size * block_count;
	RAII(void*, buffer, malloc(buffer_size), free);
	if (buffer == NULL) {
		ERRNO(errno);
		return -1;
	}

	ffs_entry_t src_entry;
	if (__ffs_entry_find(src, src_name, &src_entry) == false) {
		UNEXPECTED("'%s' partition not found => %s",
			   src->path, src_name);
		return -1;
	}

	ffs_entry_t dst_entry;
	if (__ffs_entry_find(dst, dst_name, &dst_entry) == false) {
		UNEXPECTED("'%s' partition not found => %s",
			   dst->path, dst_name);
		return -1;
	}

	uint32_t total = 0;
	uint32_t size = src_entry.actual;
	off_t offset = 0;

	if (isatty(fileno(stderr))) {
		fprintf(stderr, "%8llx: %s: copy partition %8x/%8x",
			(long long)src->offset, dst_name, src_entry.actual, total);
	}

	while (0 < size) {
		size_t count = min(buffer_size, size);

		ssize_t rc;
		rc = __ffs_entry_read(src, src_name, buffer, offset, count);
		if (rc < 0)
			return -1;
		rc = __ffs_entry_write(dst, dst_name, buffer, offset, rc);
		if (rc < 0)
			return -1;

		if (__ffs_fsync(dst) < 0)
			return -1;

		size -= rc;
		total += rc;
		offset += rc;

		if (isatty(fileno(stderr))) {
			fprintf(stderr, "\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b");
			fprintf(stderr, "%8x/%8x", (uint32_t)src_entry.actual, total);
		}
	}

	if (isatty(fileno(stderr))) {
		fprintf(stderr, "\n");
	}

	return total;
}

int fcp_compare_entry(ffs_t * src, const char * src_name,
		      ffs_t * dst, const char * dst_name)
{
	assert(src != NULL);
	assert(src_name != NULL);
	assert(dst != NULL);
	assert(dst_name != NULL);

	uint32_t block_size;
	if (__ffs_info(src, FFS_INFO_BLOCK_SIZE, &block_size) < 0)
		return -1;

	uint32_t block_count;
	if (__ffs_info(dst, FFS_INFO_BLOCK_COUNT, &block_count) < 0)
		return -1;

	size_t buffer_size = block_size * block_count;

	RAII(void*, src_buffer, malloc(buffer_size), free);
	if (src_buffer == NULL) {
		ERRNO(errno);
		return -1;
	}
	RAII(void*, dst_buffer, malloc(buffer_size), free);
	if (dst_buffer == NULL) {
		ERRNO(errno);
		return -1;
	}

	ffs_entry_t src_entry;
	if (__ffs_entry_find(src, src_name, &src_entry) == false) {
		UNEXPECTED("'%s' partition not found => %s",
			   src->path, src_name);
		return -1;
	}

	ffs_entry_t dst_entry;
	if (__ffs_entry_find(dst, dst_name, &dst_entry) == false) {
		UNEXPECTED("'%s' partition not found => %s",
			   dst->path, dst_name);
		return -1;
	}

	uint32_t total = 0;
	uint32_t size = src_entry.actual;
	off_t offset = 0;

	if (isatty(fileno(stderr))) {
		fprintf(stderr, "%8llx: %s: compare partition %8x/%8x",
			(long long)src->offset, dst_name, src_entry.actual, total);
	}

	while (0 < size) {
		size_t count = min(buffer_size, size);

		ssize_t rc;
		rc = __ffs_entry_read(src, src_name, src_buffer, offset, count);
		if (rc < 0)
			return -1;
		rc = __ffs_entry_read(dst, dst_name, dst_buffer, offset, rc);
		if (rc < 0)
			return -1;

		char * src_ptr = src_buffer;
		char * dst_ptr = dst_buffer;
		size_t cnt = 0;

		while (cnt < count) {
			size_t cmp_sz = min(count - cnt, COMPARE_SIZE);

			if (memcmp(src_ptr, dst_ptr, cmp_sz) != 0) {
				UNEXPECTED("MISCOMPARE! '%s' != '%s' at "
					   "offset '%llx'\n", src_name,
					   dst_name, (long long)offset + cnt);

				if (isatty(fileno(stderr)))
					fprintf(stderr, " <== [ERROR]\n");

				return -1;
			}

			src_ptr += cmp_sz;
			dst_ptr += cmp_sz;
			cnt += cmp_sz;
		}

		size -= rc;
		total += rc;
		offset += rc;

		if (isatty(fileno(stderr))) {
			fprintf(stderr, "\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b");
			fprintf(stderr, "%8x/%8x", src_entry.actual, total);
		}
	}

	if (isatty(fileno(stderr))) {
		fprintf(stderr, "\n");
	}

	return total;
}
