/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fpart/src/cmd_write.c $                                       */
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
 *    File: write.c
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: --write implementation
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
#include <clib/list.h>
#include <clib/list_iter.h>
#include <clib/misc.h>
#include <clib/min.h>
#include <clib/err.h>
#include <clib/raii.h>

#include "main.h"

int command_write(args_t * args)
{
	assert(args != NULL);

	char name[strlen(args->name) + 2];
	strcpy(name, args->name);
	if (strrchr(name, '$') == NULL)
		strcat(name, "$");

	char full_name[page_size];
	struct stat st;
	regex_t rx;

  	off_t __poffset;
	ffs_t *__ffs;

	RAII(FILE*, in, fopen(args->path, "r"), fclose);

	list_t list;
	list_init(&list);

	list_iter_t it;
	ffs_entry_node_t * node;

	/* ========================= */

	int write_entry(ffs_entry_t * entry)
	{
		assert(entry != NULL);

                if (__ffs_entry_name(__ffs, entry, full_name,
				     sizeof full_name) < 0)
			return -1;

		if (regexec(&rx, full_name, 0, NULL, 0) == REG_NOMATCH)
			return 0;

		if (entry->flags & FFS_FLAGS_PROTECTED &&
		    args->force != f_FORCE) {
			if (args->verbose == f_VERBOSE)
				printf("%llx: %s: protected, skipping write\n",
			       	       __poffset, full_name);
			return 0;
		}

		if (entry->actual < st.st_size) {
			if (__ffs_entry_truncate(__ffs, full_name,
						 st.st_size) < 0) {
				ERRNO(errno);
				return -1;
			}

			if (args->verbose == f_VERBOSE)
				printf("%llx: %s: truncate size '%llx'\n",
				       __poffset, full_name, st.st_size);
		}

		list_iter_init(&it, &list, LI_FLAG_FWD);
		list_for_each(&it, node, node) {
			if (node->entry.base == entry->base &&
			    node->entry.size == entry->size)
				return 0;
		}

		node = (ffs_entry_node_t *)malloc(sizeof(*node));
		assert(node != NULL);

		memcpy(&node->entry, entry, sizeof(node->entry));
		list_add_tail(&list, &node->node);

		size_t block_size = __ffs->hdr->block_size;
		size_t entry_size = entry->size * block_size;
		size_t data_size = st.st_size;
		off_t data_offset = 0;

		if (entry_size < data_size) {
			UNEXPECTED("'%s' of size '%x' too big for partition "
				  "'%s' of size '%x'", args->path, data_size,
				  args->name, entry_size);
			return -1;
		}

		if (args->block != NULL) {
			if (parse_size(args->block, &block_size) < 0)
				return -1;
			if (block_size & (__ffs->hdr->block_size-1)) {
				UNEXPECTED("'%x' block size must be multiple "
					   "of target block size '%x'",
					   block_size, __ffs->hdr->block_size);
				return -1;
			}
		}

		RAII(void*, block, malloc(block_size), free);
		if (block == NULL) {
			ERRNO(errno);
			return -1;
		}

		if (fseeko(in, 0, SEEK_SET) != 0) {
			ERRNO(errno);
			return -1;
		}

		while (0 < data_size) {
			ssize_t rc = fread(block, 1,
					   min(block_size, data_size), in);
			if (rc <= 0 && ferror(in)) {
				ERRNO(errno);
				return -1;
			}

			rc = __ffs_entry_write(__ffs, full_name, block,
					       data_offset, rc);
			if (rc < 0)
				return -1;

			__ffs_fsync(__ffs);

			data_offset += rc;
			data_size -= rc;
		}

		if (args->verbose == f_VERBOSE)
			printf("%llx: %s: read '%llx' bytes from file '%s'\n",
		       		__poffset, full_name, st.st_size, args->path);

		return 0;
	}

	int write(args_t * args, off_t poffset)
	{
		__poffset = poffset;

		const char * target = args->target;
		int debug = args->debug;

		RAII(FILE*, file, fopen_generic(target, "r+", debug), fclose);
		if (file == NULL)
			return -1;
		RAII(ffs_t*, ffs, __ffs_fopen(file, poffset), __ffs_fclose);
		if (ffs == NULL)
			return -1;

		__ffs = ffs;

		int rc = __ffs_iterate_entries(ffs, write_entry);
		if (rc == 1)
			rc = -1;

		return rc;
	}

	/* ========================= */

	if (stat(args->path, &st) < 0) {
		ERRNO(errno);
		return -1;
	}

	if (regcomp(&rx, name, REG_ICASE | REG_NOSUB) != 0) {
		ERRNO(errno);
		return -1;
	}

	int rc = command(args, write);

	regfree(&rx);

	while (!list_empty(&list))
		free(container_of(list_remove_head(&list),
				  ffs_entry_node_t, node));

	return rc;
}
