/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fpart/src/cmd_list.c $                                        */
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
 *    File: list.c
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: --list implementation
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

int command_list(args_t * args)
{
	assert(args != NULL);

	char name[strlen(args->name) + 2];
	strcpy(name, args->name);
	if (strrchr(name, '$') == NULL)
		strcat(name, "$");

	ffs_t *__ffs;

	char full_name[page_size];
	regex_t rx;

	/* ========================= */

	int __list_entry(ffs_entry_t * entry)
	{

		uint32_t offset = entry->base * __ffs->hdr->block_size;
		uint32_t size = entry->size * __ffs->hdr->block_size;

                if (__ffs_entry_name(__ffs, entry, full_name,
				     sizeof full_name) < 0)
			return -1;

		if (regexec(&rx, full_name, 0, NULL, 0) == REG_NOMATCH)
			return 0;

		char type;
		if (entry->type == FFS_TYPE_LOGICAL) {
			type ='l';
		} else if (entry->type == FFS_TYPE_DATA) {
			type ='d';
		} else if (entry->type == FFS_TYPE_PARTITION) {
			type ='p';
		}
		fprintf(stdout, "%3d [%08x-%08x] [%8x:%8x] ",
			entry->id, offset, offset+size-1,
			size, entry->actual);

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

		return 0;
	}

	int list(args_t * args, off_t poffset)
	{
		const char * target = args->target;
		int debug = args->debug;
		int rc = 0;

		RAII(FILE*, file, fopen_generic(target, "r", debug), fclose);
		if (file == NULL)
			return -1;

		if (__ffs_fcheck(file, poffset) < 0)
			return -1;

		RAII(ffs_t*, ffs, __ffs_fopen(file, poffset), __ffs_fclose);
		if (ffs == NULL)
			return -1;

		if (0 < ffs->count) {
			printf("========================[ PARTITION TABLE"
				" 0x%llx ]=======================\n",
				(long long)ffs->offset);
			printf("vers:%04x size:%04x * blk:%06x blk(s):"
				"%06x * entsz:%06x ent(s):%06x\n",
				ffs->hdr->version,
				ffs->hdr->size,
				ffs->hdr->block_size,
				ffs->hdr->block_count,
				ffs->hdr->entry_size,
				ffs->hdr->entry_count);
			printf("----------------------------------"
			       "----------------------------------"
			       "-------\n");

			__ffs = ffs;

			rc = __ffs_iterate_entries(ffs, __list_entry);
			if (rc == 1)
				rc = -1;

			printf("\n");
		}

		return rc;
	}

	/* ========================= */

	if (regcomp(&rx, name, REG_ICASE | REG_NOSUB) != 0) {
		ERRNO(errno);
		return -1;
	}

	int rc = command(args, list);

	regfree(&rx);

	return rc;
}
