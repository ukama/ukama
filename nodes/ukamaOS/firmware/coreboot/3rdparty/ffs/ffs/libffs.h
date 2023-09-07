/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: ffs/libffs.h $                                                */
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

#ifndef __LIBFFS_H__
#define __LIBFFS_H__

#include <stdbool.h>
#include <stdarg.h>


#include "ffs.h"

typedef struct ffs_entry ffs_entry_t;
typedef struct ffs_hdr ffs_hdr_t;
typedef enum type ffs_type_t;

#define FFS_EXCEPTION_DATA	1024

/* ============================================================ */

/*!
 * @brief ffs I/O interface
 */
struct ffs {
    ffs_hdr_t * hdr;

    char * path;
    off_t offset;

    FILE * file;
    uint32_t count;

    bool dirty;
};

typedef struct ffs ffs_t;

struct ffs_exception {
    int rc;
    char data[FFS_EXCEPTION_DATA];
};

typedef struct ffs_exception ffs_exception_t;

#define FFS_PARTITION_NAME		"part"

#define FFS_INFO_ERROR			0
#define FFS_INFO_MAGIC			1
#define FFS_INFO_VERSION		2
#define FFS_INFO_ENTRY_SIZE		3
#define FFS_INFO_ENTRY_COUNT		4
#define FFS_INFO_BLOCK_SIZE		5
#define FFS_INFO_BLOCK_COUNT		6
#define FFS_INFO_OFFSET			8

#define FFS_CHECK_PATH			-3
#define FFS_CHECK_HEADER_MAGIC		-4
#define FFS_CHECK_HEADER_CHECKSUM	-5
#define FFS_CHECK_ENTRY_CHECKSUM	-6

#ifdef __cplusplus
extern "C" {
#endif

extern int __ffs_fcheck(FILE *, off_t)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

extern int __ffs_check(const char *, off_t)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

extern ffs_t * __ffs_fcreate(FILE *, off_t, uint32_t, uint32_t)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

extern ffs_t * __ffs_create(const char *, off_t, uint32_t, uint32_t)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

extern ffs_t * __ffs_fopen(FILE *, off_t)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

extern ffs_t * __ffs_open(const char *, off_t)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

extern int __ffs_info(ffs_t *, int, uint32_t *)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

extern int __ffs_close(ffs_t *)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

extern int __ffs_fclose(ffs_t *)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

extern int __ffs_fsync(ffs_t *)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

extern int __ffs_list_entries(ffs_t *, const char *, bool, FILE *)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

extern int __ffs_iterate_entries(ffs_t *, int (*)(ffs_entry_t*))
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

extern int __ffs_entry_find(ffs_t *, const char *, ffs_entry_t *)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

extern int __ffs_entry_find_parent(ffs_t *, const char *, ffs_entry_t *)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

extern int __ffs_entry_name(ffs_t *, ffs_entry_t *, char *, size_t)
/*! @cond */ __nonnull ((1,2,3)) /*! @endcond */ ;

extern int __ffs_entry_add(ffs_t *, const char *, off_t,
			    uint32_t, ffs_type_t, uint32_t)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

extern int __ffs_entry_delete(ffs_t *, const char *)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

extern int __ffs_entry_user_get(ffs_t *, const char *, uint32_t, uint32_t *)
/*! @cond */ __nonnull ((1,2,4)) /*! @endcond */ ;

extern int __ffs_entry_user_put(ffs_t *, const char *, uint32_t, uint32_t)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

extern ssize_t __ffs_entry_hexdump(ffs_t *, const char *, FILE *)
/*! @cond */ __nonnull ((1,2,3)) /*! @endcond */ ;

extern ssize_t __ffs_entry_truncate(ffs_t *, const char *, size_t)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

extern ssize_t __ffs_entry_read(ffs_t *, const char *, void *, off_t, size_t)
/*! @cond */ __nonnull ((1,2,3)) /*! @endcond */ ;

extern ssize_t __ffs_entry_write(ffs_t *, const char *, const void *,
				 off_t, size_t)
/*! @cond */ __nonnull ((1,2,3)) /*! @endcond */ ;

extern ssize_t __ffs_entry_copy(ffs_t *, ffs_t *, const char *)
/*! @cond */ __nonnull ((1,2,3)) /*! @endcond */ ;

extern ssize_t __ffs_entry_compare(ffs_t *, ffs_t *, const char *)
/*! @cond */ __nonnull ((1,2,3)) /*! @endcond */ ;

extern int __ffs_entry_list(ffs_t *, ffs_entry_t ** list)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

#ifdef __cplusplus
}
#endif

/* ============================================================ */

#endif /* __LIBFFS_H__ */
