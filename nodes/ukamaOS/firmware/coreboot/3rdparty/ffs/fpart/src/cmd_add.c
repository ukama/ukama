/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fpart/src/cmd_add.c $                                         */
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
 *    File: add.c
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: --add implementation
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
#include <clib/raii.h>

#include "main.h"

int command_add(args_t * args)
{
	assert(args != NULL);

	/* ========================= */

	int add(args_t * args, off_t poffset)
	{
		int rc = 0;

		off_t offset = 0;
		uint32_t size = 0;
		uint32_t flags = 0;

		rc = parse_offset(args->offset, &offset);
		if (rc < 0)
			return rc;
		rc = parse_size(args->size, &size);
		if (rc < 0)
			return rc;
		rc = parse_size(args->flags, &flags);
		if (rc < 0)
			return rc;

		ffs_type_t type = FFS_TYPE_DATA;
		if (args->logical == f_LOGICAL)
			type = FFS_TYPE_LOGICAL;

		const char * target = args->target;
		int debug = args->debug;

		RAII(FILE *, file, fopen_generic(target, "r+", debug), fclose);
		RAII(ffs_t *, ffs,  __ffs_fopen(file, poffset), __ffs_fclose);

		rc = __ffs_entry_add(ffs, args->name, offset, size,
				     type, flags);
		if (rc < 0)
			return rc;

		if (args->verbose == f_VERBOSE)
			printf("%llx: %s: add partition at offset '%llx' size "
			       "'%x' type '%d' flags '%x'\n", (long long)poffset,
				args->name, (long long)offset, size, type, flags);

		return rc;
	}

	/* ========================= */

	return command(args, add);
}
