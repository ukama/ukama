/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/test/err.c $                                             */
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

#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <errno.h>
#include <string.h>
#include <libgen.h>

#include <clib/version.h>
#include <clib/err.h>

extern char * program_invocation_short_name;

#define FOO_MAJOR	0x01
#define FOO_MINOR	0x00
#define FOO_PATCH	0x00
#define FOO_VER		VER(FOO_MAJOR, FOO_MINOR, FOO_PATCH)

#define FOO_UNEXPECTED(f, ...)		({		\
	UNEXPECTED(f, ##__VA_ARGS__);			\
	VERSION(FOO_VER, "%s", program_invocation_short_name);	\
					})


int main (int argc, char * argv[]) {
	ERRNO(EINVAL);
	FOO_UNEXPECTED("cannot frob the ka-knob");

	goto error;

	if (false) {
		err_t * err = NULL;
error:
		while ((err = err_get()) != NULL) {
			switch (err_type(err)) {
			case ERR_VERSION:
				fprintf(stderr, "%s: %s : %s(%d) : v%d.%02d.%04d %.*s\n",
					basename((char *)argv[0]),
					err_type_name(err), basename(err_file(err)), err_line(err),
					VER_TO_MAJOR(err_code(err)), VER_TO_MINOR(err_code(err)),
					VER_TO_PATCH(err_code(err)),
					err_size(err), (char *)err_data(err));
				break;
			default:
				fprintf(stderr, "%s: %s : %s(%d) : (code=%d) %.*s\n",
					basename((char *)argv[0]),
					err_type_name(err), basename(err_file(err)), err_line(err),
					err_code(err), err_size(err), (char *)err_data(err));
			}
		}
	}

	return 0;
}

