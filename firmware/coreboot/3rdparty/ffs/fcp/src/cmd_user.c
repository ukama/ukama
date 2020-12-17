/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fcp/src/cmd_user.c $                                          */
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
 *    File: cmd_user.c
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: user implementation
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

#define DELIM	'='

static int __user(args_t * args, off_t offset)
{
	assert(args != NULL);

	char * type = args->dst_type;
	char * target = args->dst_target;
	char * name = args->dst_name;

	RAII(FILE*, file, __fopen(type, target, "r+", debug), fclose);
	if (file == NULL)
		return -1;
	if (check_file(target, file, offset) < 0)
		return -1;
	RAII(ffs_t*, ffs, __ffs_fopen(file, offset), __ffs_fclose);
	if (ffs == NULL)
		return -1;

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

	char full_name[page_size];
	if (__ffs_entry_name(ffs, &entry, full_name,
			     sizeof full_name) < 0)
		return -1;

	if (args->protected != f_PROTECTED &&
	    entry.flags && FFS_FLAGS_PROTECTED) {
		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: %s: protected (skip)\n",
				(long long)offset, full_name);
		return 0;
	}

	// parse <word>[=<value>]

	if (args->opt_nr <= 1) {
		for (uint32_t word=0; word<FFS_USER_WORDS; word++) {
			uint32_t value = 0;

			if (__ffs_entry_user_get(ffs, name, word, &value) < 0)
				return -1;

			fprintf(stdout, "%8llx: %s: [%02d] = %08x\n",
				(long long)offset, full_name, word, value);
		}
		fprintf(stdout, "\n");
	} else {
		for (int i=1; i<args->opt_nr; i++) {
			if (args->opt[i] == '\0') {
				UNEXPECTED("invalid user '%s', use form "
					   "'<word>[=<value>]'\n",
					   args->opt[i]);
				return -1;
			}

			char * __value = strrchr(args->opt[i], DELIM);
			if (__value != NULL)
				*__value = '\0';

			uint32_t word = 0;
			if (parse_number(args->opt[i], &word) < 0)
				return -1;

			uint32_t value = 0;
			if (__value != NULL) {	// write
				if (parse_number(__value+1, &value) < 0)
					return -1;

				*__value = DELIM;

				if (__ffs_entry_user_put(ffs, full_name,
							 word, value) < 0)
					return -1;

				if (args->verbose == f_VERBOSE)
					fprintf(stderr, "%8llx: %s: [%02d] = "
						"%08x\n", (long long)offset, full_name,
						word, value);
			} else {		// read
				if (__ffs_entry_user_get(ffs, full_name,
							 word, &value) < 0)
					return -1;

				if (isatty(fileno(stdout)))
					fprintf(stdout, "%8llx: %s: [%02d] = "
						"%08x\n", (long long)offset, full_name,
						word, value);
				else
					fprintf(stdout, "%x", value);
			}
		}
	}

	return 0;
}

int command_user(args_t * args)
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

		rc = __user(args, offset);
		if (rc < 0)
			break;

		if (*end == '\0')
			break;
		end++;
	}

	return rc;
}
