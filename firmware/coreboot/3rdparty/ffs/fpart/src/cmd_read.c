/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fpart/src/cmd_read.c $                                        */
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
 *    File: read.c
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: --read implementation
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

int command_read(args_t * args)
{
	assert(args != NULL);

	if (access(args->path, F_OK) == 0 && args->force != f_FORCE) {
		UNEXPECTED("output file '%s' already exists,  use --force "
		 	   "to overwrite\n", args->path);
		return -1;
	}

	list_t list;
	list_init(&list);

	list_iter_t it;
	ffs_entry_node_t * node;

	/* ========================= */

	int read(args_t * args, off_t poffset)
	{
		const char * target = args->target;
		int debug = args->debug;

		RAII(FILE*, file, fopen_generic(target, "r", debug), fclose);
		if (file == NULL)
			return -1;
		RAII(ffs_t*, ffs, __ffs_fopen(file, poffset), __ffs_fclose);
		if (ffs == NULL)
			return -1;

		size_t block_size = ffs->hdr->block_size;

		if (args->block != NULL) {
			if (parse_size(args->block, &block_size) < 0)
				return -1;
			if (block_size & (ffs->hdr->block_size-1)) {
				UNEXPECTED("'%x' block size must be multiple "
					   "of target block size '%x'",
					   block_size, ffs->hdr->block_size);
				return -1;
			}
		}

		RAII(void*, block, malloc(block_size), free);
		if (block == NULL) {
			ERRNO(errno);
			return -1;
		}

		ffs_entry_t entry;
		if (__ffs_entry_find(ffs, args->name, &entry) == false) {
			UNEXPECTED("'%s' partition not found => %s",
				   args->target, args->name);
			return -1;
		}

		list_iter_init(&it, &list, LI_FLAG_FWD);
		list_for_each(&it, node, node) {
			if (node->entry.base == entry.base &&
			    node->entry.size == entry.size)
				return 0;
		}

		node = (ffs_entry_node_t *)malloc(sizeof(*node));
		assert(node != NULL);

		memcpy(&node->entry, &entry, sizeof(node->entry));
		list_add_tail(&list, &node->node);

		RAII(FILE*, out, fopen(args->path, "w+"), fclose);
		if (out == NULL) {
			ERRNO(errno);
			return -1;
		}

		size_t data_size = entry.actual;
		off_t data_offset = 0;

		while (0 < data_size) {
			ssize_t rc;
			rc = __ffs_entry_read(ffs, args->name,
					      block, data_offset,
					      min(block_size, data_size));

			rc = fwrite(block, 1, rc, out);
			if (rc <= 0 && ferror(out)) {
				ERRNO(errno);
				return -1;
			}

			data_size -= rc;
			data_offset += rc;
		}

		if (args->verbose == f_VERBOSE)
			printf("%llx: %s: wrote '%x' bytes to file '%s'\n",
			       poffset, args->name, entry.actual, args->path);

		return 0;
	}

	/* ========================= */

	int rc = command(args, read);

	while (!list_empty(&list))
		free(container_of(list_remove_head(&list),
				  ffs_entry_node_t, node));

	return rc;
}

