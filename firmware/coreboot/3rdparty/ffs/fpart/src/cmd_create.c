/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fpart/src/cmd_create.c $                                      */
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
 *    File: create.c
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: --create implementation
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

int command_create(args_t * args)
{
	assert(args != NULL);

	uint32_t block = 0;
	uint32_t size = 0;
	uint32_t pad = 0xff;

	if (parse_size(args->block, &block) < 0)
		return -1;
	if (parse_size(args->size, &size) < 0)
		return -1;
	if (args->pad != NULL)
		if (parse_size(args->pad, &pad) < 0)
			return -1;

	struct stat st;
	if (stat(args->target, &st) < 0) {
		if (errno == ENOENT) {
			create_regular_file(args->target, size, (uint8_t)pad);
		} else {
			ERRNO(errno);
			return -1;
		}
	} else {
		if (st.st_size != size) {
			create_regular_file(args->target, size, (uint8_t)pad);
		} else {
			if (args->force != f_FORCE && st.st_size != size) {
				UNEXPECTED("--size '%d' differs from actual "
					   "size '%lld', use --force to "
					   "override", size, (long long)st.st_size);
				return -1;
			}
		}
	}

	/* ========================= */

	int create(args_t * args, off_t poffset)
	{
		if (args->verbose == f_VERBOSE)
			printf("%llx: create partition table\n", (long long)poffset);

		const char * target = args->target;
		int debug = args->debug;

		RAII(FILE*, file, fopen_generic(target, "r+", debug), fclose);
		if (file == NULL)
			return -1;
		RAII(ffs_t*, ffs, __ffs_fcreate(file, poffset, block,
		     size / block), __ffs_fclose);
		if (ffs == NULL)
			return -1;

		return 0;
	}

	/* ========================= */

	return command(args, create);
}
