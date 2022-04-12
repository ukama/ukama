/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fpart/src/cmd_user.c $                                        */
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
 *    File: user.c
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: --user implementation
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

int command_user(args_t * args)
{
	assert(args != NULL);

	char name[strlen(args->name) + 2];
	strcpy(name, args->name);
	if (strrchr(name, '$') == NULL)
		strcat(name, "$");

	char full_name[page_size];
	regex_t rx;

	off_t __poffset;
	ffs_t * __ffs;

	uint32_t user = 0;
	if (args->user != NULL)
		if (parse_size(args->user, &user) < 0)
			return -1;

	if (FFS_USER_WORDS <= user) {
		UNEXPECTED("invalid user word '%d', valid range [0..%d]",
			   user, FFS_USER_WORDS);
		return -1;
	}

	int user_entry(ffs_entry_t * entry)
	{
		assert(entry != NULL);

		if (__ffs_entry_name(__ffs, entry, full_name,
			 	     sizeof full_name) < 0)
			return -1;

		if (regexec(&rx, full_name, 0, NULL, 0) == REG_NOMATCH)
			return 0;

		if (entry->flags & FFS_FLAGS_PROTECTED &&
		    args->force != f_FORCE) {
			printf("%llx: %s: protected, skipping user[%d]\n",
			       (long long)__poffset, full_name, user);
			return 0;
		}

		uint32_t value = 0;
		if (args->value != NULL) {
			if (args->value != NULL)
				if (parse_size(args->value, &value) < 0)
					return -1;

			if (__ffs_entry_user_put(__ffs, full_name,
						 user, value) < 0)
				return -1;

			if (args->verbose == f_VERBOSE)
				printf("%llx: %s: user[%d] = '%x'\n", (long long)__poffset,
				       args->name, user, value);
		} else {
			if (__ffs_entry_user_get(__ffs, args->name,
						 user, &value) < 0)
				return -1;

			printf("%llx: %s: user[%d] = '%x'\n", (long long)__poffset,
			       args->name, user, value);
		}

		return 0;
	}

	int __user(args_t * args, off_t poffset)
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

		int rc = __ffs_iterate_entries(ffs, user_entry);
		if (rc == 1)
			return -1;

		return rc;
	}

	/* ========================= */

	if (regcomp(&rx, name, REG_ICASE | REG_NOSUB) != 0) {
		ERRNO(errno);
		return-1;
	}

	int rc = command(args, __user);

	regfree(&rx);

	return rc;
}
