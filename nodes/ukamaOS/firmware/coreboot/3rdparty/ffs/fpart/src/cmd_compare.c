/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fpart/src/cmd_compare.c $                                     */
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
 *    File: compare.c
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: --compare implementation
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

int command_compare(args_t * args)
{
	assert(args != NULL);

	char name[strlen(args->name) + 2];
	strcpy(name, args->name);
	if (strrchr(name, '$') == NULL)
		strcat(name, "$");

	char full_name[page_size];
	regex_t rx;

	off_t __poffset;
	ffs_t * __in, * __out;

	list_t list;
	list_init(&list);

	list_iter_t it;
	ffs_entry_node_t * node;

	/* ========================= */

	int compare_entry(ffs_entry_t * src)
	{
                if (__ffs_entry_name(__in, src, full_name,
				     sizeof full_name) < 0)
			return -1;
		if (regexec(&rx, full_name, 0, NULL, 0) == REG_NOMATCH)
			return 0;

		if (src->flags & FFS_FLAGS_PROTECTED &&
		    args->force != f_FORCE) {
			if (args->verbose == f_VERBOSE)
				printf("%llx: %s: protected, skipping "
				       "compare\n", __poffset, full_name);
			return 0;
		}

		list_iter_init(&it, &list, LI_FLAG_FWD);
		list_for_each(&it, node, node) {
			if (node->entry.base == src->base &&
			    node->entry.size == src->size)
				return 0;
		}

		node = (ffs_entry_node_t *)malloc(sizeof(*node));
		assert(node != NULL);

		memcpy(&node->entry, src, sizeof(node->entry));
		list_add_tail(&list, &node->node);

		ffs_entry_t dst;
		if (__ffs_entry_find(__out, full_name, &dst) == false) {
			UNEXPECTED("'%s' entry not found '%s'",
				   args->target, args->name);
			return -1;
		}

		if (verify_operation(full_name, __in, src, __out, &dst) < 0)
			return -1;

		size_t block_size = __in->hdr->block_size;

		if (args->block != NULL) {
			if (parse_size(args->block, &block_size) < 0)
				return -1;
			if (block_size & (__in->hdr->block_size-1)) {
				UNEXPECTED("'%x' block size must be multiple "
					   "of target block size '%x'",
					   block_size, __in->hdr->block_size);
				return -1;
			}
		}

		RAII(void*, src_block, malloc(block_size), free);
		if (src_block == NULL) {
			ERRNO(errno);
			return -1;
		}
		RAII(void*, dst_block, malloc(block_size), free);
		if (dst_block == NULL) {
			ERRNO(errno);
			return -1;
		}

		size_t data_size = src->actual;
		off_t data_offset = 0;

		while (0 < data_size) {
			ssize_t rc;

			rc = __ffs_entry_read(__in, full_name, src_block,
					      data_offset,
					      min(block_size, data_size));
			if (rc < 0)
				return -1;

			rc = __ffs_entry_read(__out, full_name, dst_block,
					     data_offset, rc);
			if (rc < 0)
				return -1;

			if (memcmp(src_block, dst_block, block_size) != 0) {
				if (args->verbose == f_VERBOSE)
					printf("%llx: %s: miscompare at offset"
					       "'%llx'\n", __poffset, full_name,
					       data_offset);
				UNEXPECTED("data miscompare at offset range "
					   "[%llx..%llx]", data_offset,
					   data_offset + block_size - 1);
				return -1;
			}


			data_offset += rc;
			data_size -= rc;
		}

		if (args->verbose == f_VERBOSE)
			printf("%llx: %s: compare to '%s'\n", __poffset,
		       	       full_name, args->path);

		return 0;
	}

	int compare(args_t * args, off_t poffset)
	{
		__poffset = poffset;

		const char * path = args->path;
		const char * target = args->target;
		int debug = args->debug;

		RAII(FILE*, inf, fopen_generic(path, "r", debug), fclose);
		if (inf == NULL)
			return -1;
		RAII(FILE*, outf, fopen_generic(target, "r", debug), fclose);
		if (outf == NULL)
			return -1;

		RAII(ffs_t*, in, __ffs_fopen(inf, poffset), __ffs_fclose);
		if (in == NULL)
			return -1;
		RAII(ffs_t*, out, __ffs_fopen(outf, poffset), __ffs_fclose);
		if (out == NULL)
			return -1;

		__in = in;
		__out = out;

		int rc = __ffs_iterate_entries(in, compare_entry);
		if (rc == 1)
			rc = -1;

		return rc;
	}

	/* ========================= */

	if (regcomp(&rx, name, REG_ICASE | REG_NOSUB) != 0) {
		ERRNO(errno);
		return -1;
	}

	int rc = command(args, compare);

	regfree(&rx);

	while (!list_empty(&list))
		free(container_of(list_remove_head(&list),
				  ffs_entry_node_t, node));

	return rc;
}
