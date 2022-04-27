/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fcp/src/cmd_list.c $                                          */
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
 *    File: cmd_list.c
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: list implementation
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

#include "misc.h"
#include "main.h"

static int list(args_t * args, off_t offset)
{
	assert(args != NULL);

	char * type = args->dst_type;
	char * target = args->dst_target;
	char * name = args->dst_name;

	RAII(FILE*, file, __fopen(type, target, "r", debug), fclose);
	if (file == NULL)
		return -1;
	if (check_file(target, file, offset) < 0)
		return -1;
	RAII(ffs_t*, ffs, __ffs_fopen(file, offset), __ffs_fclose);
	if (ffs == NULL)
		return -1;

	if (ffs->count <= 0)
		return 0;

	fprintf(stdout, "========================[ PARTITION TABLE"
		" 0x%llx ]=======================\n", (long long)ffs->offset);
	fprintf(stdout, "vers:%04x size:%04x * blk:%06x blk(s):"
		"%06x * entsz:%06x ent(s):%06x\n",
		ffs->hdr->version, ffs->hdr->size, ffs->hdr->block_size,
		ffs->hdr->block_count, ffs->hdr->entry_size,
		ffs->hdr->entry_count);
	fprintf(stdout, "----------------------------------"
	       "----------------------------------"
	       "-------\n");

	RAII(entry_list_t*, entry_list, entry_list_create_by_regex(ffs, name),
	     entry_list_delete);
	if (entry_list == NULL)
		return -1;

	list_iter_t it;
	entry_node_t * entry_node;

	list_iter_init(&it, &entry_list->list, LI_FLAG_FWD);
	list_for_each(&it, entry_node, node) {
		ffs_entry_t * entry = &entry_node->entry;
		char full_name[page_size];

		if (__ffs_entry_name(ffs, entry, full_name,
				     sizeof full_name) < 0)
			return -1;

		uint32_t offset = entry->base * ffs->hdr->block_size;
		uint32_t size = entry->size * ffs->hdr->block_size;

		char type;
		if (entry->type == FFS_TYPE_LOGICAL)
			type ='l';
		else if (entry->type == FFS_TYPE_DATA)
			type ='d';
		else if (entry->type == FFS_TYPE_PARTITION)
			type ='p';
		fprintf(stdout, "%3d [%08x-%08x] [%8x:%8x] ",
			entry->id, offset, offset+size-1, size, entry->actual);

		fprintf(stdout, "[%c%c%c%c%c%c] %s\n",
			type, '-', '-', '-',
			entry->flags & FFS_FLAGS_U_BOOT_ENV ? 'b' : '-',
			entry->flags & FFS_FLAGS_PROTECTED ? 'p' : '-',
			full_name);

		if (args->verbose == f_VERBOSE) {
			for (int i=0; i<FFS_USER_WORDS; i++) {
				fprintf(stdout, "[%2d] %8x ", i,
					entry->user.data[i]);
				if ((i+1) % 4 == 0)
					fprintf(stdout, "\n");
			}
		}
	}
	fprintf(stdout, "\n");

	return 0;
}

int command_list(args_t * args)
{
	assert(args != NULL);

	int rc = 0;

	char * end = (char *)args->offset;
	while (rc == 0 && end != NULL && *end != '\0') {
		errno = 0;
		off_t offset = strtoull(end, &end, 0);
		if (end == NULL || errno != 0) {
			UNEXPECTED("invalid --offset specified '%s'",
				   args->offset);
			return -1;
		}

		if (*end != ',' && *end != ':' && *end != '\0') {
			UNEXPECTED("invalid --offset separator "
				   "character '%c'", *end);
			return -1;
		}

		rc = list(args, offset);
		if (rc < 0)
			break;

		if (*end == '\0')
			break;
		end++;
	}

	return rc;
}
