/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: ffs/src/libffs2.c $                                           */
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
 *   File: libffs2.c
 * Author: Shaun Wetzstein <shaun@us.ibm.com>
 *  Descr: FFS IO interface
 *   Note:
 *   Date: 05/07/12
 */

#include <sys/types.h>
#include <sys/stat.h>

#include <stdlib.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <unistd.h>
#include <fcntl.h>
#include <endian.h>
#include <libgen.h>

#include "libffs2.h"

#include <clib/builtin.h>
#include <clib/checksum.h>
#include <clib/misc.h>
#include <clib/err.h>
#include <clib/raii.h>

#define FFS_ERRSTR_MAX		1024

struct ffs_error {
	int errnum;
	char errstr[FFS_ERRSTR_MAX];
};

typedef struct ffs_error ffs_error_t;

static ffs_error_t __error;

/* ============================================================ */

extern void ffs_errclr(void)
{
	__error.errnum = __error.errstr[0] = 0;
}

extern int ffs_errnum(void)
{
	return __error.errnum;
}

extern const char *ffs_errstr(void)
{
	return __error.errnum ? __error.errstr : NULL;
}

int ffs_check(const char *path, off_t offset)
{
	RAII(FILE*, file, fopen(path, "r"), fclose);
	if (file == NULL) {
		__error.errnum = -1;
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : %s (errno=%d)\n",
			 program_invocation_short_name, "errno",
			 __FILE__, __LINE__, strerror(errno), errno);

		return -1;
	}

	int rc = __ffs_fcheck(file, offset);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		// __ffs_check will return FFS_CHECK_* const's
	}

	return rc;
}

ffs_t *ffs_create(const char *path, off_t offset, uint32_t block_size,
		  uint32_t block_count)
{
	ffs_t *self = __ffs_create(path, offset, block_size, block_count);
	if (self == NULL) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));
	}

	return self;
}

ffs_t *ffs_open(const char *path, off_t offset)
{
	FILE * file = fopen(path, "r+");
	if (file == NULL) {
		__error.errnum = -1;
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : %s (errno=%d)\n",
			 program_invocation_short_name, "errno",
			 __FILE__, __LINE__, strerror(errno), errno);

		return NULL;
	}

	ffs_t *self = __ffs_fopen(file, offset);
	if (self == NULL) {
		fclose(file);

		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));
	}

	return self;
}

int ffs_info(ffs_t *self, int name, uint32_t *value)
{
	int rc = __ffs_info(self, name, value);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

int ffs_close(ffs_t * self)
{
	int rc = __ffs_close(self);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

int ffs_fsync(ffs_t * self)
{
	int rc = __ffs_fsync(self);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

int ffs_list_entries(ffs_t * self, FILE * out)
{
	int rc = __ffs_list_entries(self, ".*", true, out);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

int ffs_iterate_entries(ffs_t * self, int (*func) (ffs_entry_t *))
{
	int rc = __ffs_iterate_entries(self, func);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

int ffs_entry_find(ffs_t * self, const char *path, ffs_entry_t * entry)
{
	int rc = __ffs_entry_find(self, path, entry);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

int ffs_entry_find_parent(ffs_t * self, const char *path, ffs_entry_t * entry)
{
	int rc = __ffs_entry_find_parent(self, path, entry);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

int ffs_entry_add(ffs_t * self, const char *path, off_t offset, size_t size,
		  ffs_type_t type, uint32_t flags)
{
	int rc = __ffs_entry_add(self, path, offset, size, type, flags);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

int ffs_entry_delete(ffs_t * self, const char *path)
{
	int rc = __ffs_entry_delete(self, path);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

int ffs_entry_user_get(ffs_t * self, const char *path, uint32_t word,
		       uint32_t * value)
{
	int rc = __ffs_entry_user_get(self, path, word, value);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

int ffs_entry_user_put(ffs_t * self, const char *path, uint32_t word,
		       uint32_t value)
{
	int rc = __ffs_entry_user_put(self, path, word, value);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

ssize_t ffs_entry_hexdump(ffs_t * self, const char *path, FILE * out)
{
	ssize_t rc = __ffs_entry_hexdump(self, path, out);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

ssize_t ffs_entry_truncate(ffs_t * self, const char *path, off_t offset,
		           uint8_t pad __unused__)
{
        return ffs_entry_truncate_no_pad(self, path, offset);
}

ssize_t ffs_entry_truncate_no_pad(ffs_t * self, const char *path, off_t offset)
{
	ssize_t rc = __ffs_entry_truncate(self, path, offset);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

ssize_t ffs_entry_read(ffs_t * self, const char *path, void *buf, off_t offset,
		       size_t count)
{
	ssize_t rc = __ffs_entry_read(self, path, buf, offset, count);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

ssize_t ffs_entry_write(ffs_t * self, const char *path, const void *buf,
			off_t offset, size_t count)
{
	ssize_t rc = __ffs_entry_write(self, path, buf, offset, count);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

ssize_t ffs_entry_list(ffs_t * self, ffs_entry_t ** list)
{
	ssize_t rc = __ffs_entry_list(self, list);
	if (rc < 0) {
		err_t *err = err_get();
		assert(err != NULL);

		__error.errnum = err_code(err);
		snprintf(__error.errstr, sizeof __error.errstr,
			 "%s: %s : %s(%d) : (code=%d) %.*s\n",
			 program_invocation_short_name,
			 err_type_name(err), err_file(err), err_line(err),
			 err_code(err), err_size(err), (char *)err_data(err));

		rc = -1;
	}

	return rc;
}

/* ============================================================ */
