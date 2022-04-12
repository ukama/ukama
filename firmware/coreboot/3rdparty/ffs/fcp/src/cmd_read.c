/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fcp/src/cmd_read.c $                                          */
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
 *    File: cmd_read.c
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: read implementation
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

#include <clib/attribute.h>
#include <clib/list.h>
#include <clib/list_iter.h>
#include <clib/misc.h>
#include <clib/min.h>
#include <clib/err.h>
#include <clib/raii.h>

#include "misc.h"
#include "main.h"

static int __read(args_t * args, off_t offset, entry_list_t * done_list)
{
	assert(args != NULL);

	char * type = args->src_type;
	char * target = args->src_target;
	char * name = args->src_name;

	char * out_path = args->dst_target;

	RAII(FILE*, file, __fopen(type, target, "r", debug), fclose);
	if (file == NULL)
		return -1;
	if (check_file(target, file, offset) < 0)
		return -1;
	RAII(ffs_t*, ffs, __ffs_fopen(file, offset), __ffs_fclose);
	if (ffs == NULL)
		return -1;

	done_list->ffs = ffs;

	if (ffs->count <= 0)
		return 0;

	uint32_t block_size;
	if (__ffs_info(ffs, FFS_INFO_BLOCK_SIZE, &block_size) < 0)
		return -1;

	ffs_entry_t entry;
	if (__ffs_entry_find(ffs, name, &entry) == false) {
		UNEXPECTED("partition entry '%s' not found\n", name);
		return -1;
	}

	if (entry.actual == 0) {
		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: %s: is empty (skip)\n",
		       		(long long)offset, name);
		return 0;
	}

	char full_name[page_size];
	if (__ffs_entry_name(ffs, &entry, full_name,
			     sizeof full_name) < 0)
		return -1;

	if (entry.type == FFS_TYPE_LOGICAL) {
		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: %s: logical (skip)\n",
		       		(long long)offset, name);
		return 0;
	}

	if (entry_list_exists(done_list, &entry) == 1) {
		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: %s: write to '%s' (skip)\n",
			        (long long)offset, full_name, out_path);
		return 0;
	}

	if (entry_list_add(done_list, &entry) < 0)
		return -1;

	if (strcmp(out_path, "-") == 0) {
		if (fcp_read_entry(ffs, full_name, stdout) < 0)
			return -1;
	} else {
		if (access(out_path, F_OK) == 0 && args->force != f_FORCE) {
			UNEXPECTED("output file '%s' already exists, "
				   "use --force to overwrite\n", out_path);
			return -1;
		}

		RAII(FILE*, out, fopen(out_path, "w+"), fclose);
		if (out == NULL) {
			ERRNO(errno);
			return -1;
		}

		if (fcp_read_entry(ffs, full_name, out) < 0)
			return -1;

		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: %s: write to '%s' (done)\n",
		       		(long long)offset, full_name, out_path);
	}

	return 0;
}

int command_read(args_t * args)
{
	assert(args != NULL);

	int rc = 0;

	RAII(entry_list_t*, done_list, entry_list_create(NULL),
	     entry_list_delete);
	if (done_list == NULL)
		return -1;

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

		rc = __read(args, offset, done_list);
		if (rc < 0)
			break;

		if (*end == '\0')
			break;
		end++;
	}

	return rc;
}
