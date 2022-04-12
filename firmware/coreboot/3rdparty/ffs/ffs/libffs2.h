/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: ffs/libffs2.h $                                               */
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

/*!
 * @file libffs2.h
 * @brief FSP Flash File Structure API
 * @details This library provides an API to the FFS partitioning scheme.
 * @author Shaun Wetzstein <shaun@us.ibm.com>
 * @date 2010-2012
 */

#ifndef __LIBFFS2_H__
#define __LIBFFS2_H__

#ifdef __cplusplus
extern "C" {
#endif

#include "libffs.h"

/* FFS_OPEN_* are deprecated, use FFS_CHECK_* instead */
#define FFS_OPEN_INVALID_PATH		FFS_CHECK_PATH
#define FFS_OPEN_FOPEN_FAILURE		-1
#define FFS_OPEN_MAGIC_CHECK		FFS_CHECK_HEADER_MAGIC
#define FFS_OPEN_CRC_CHECK		FFS_CHECK_HEADER_CHECKSUM

/*!
 * @brief Clears the (global) FFS error indicator
 * @memberof ffs
 * @return None
 */
extern void ffs_errclr(void);

/*!
 * @brief Return the error number of an FFS error
 * @memberof ffs
 * @return non-0 on success, '0' if no pending error
 */
extern int ffs_errnum(void);

/*!
 * @brief Return the error string of an FFS error
 * @memberof ffs
 * @return non-NULL on success, 'NULL' if no pending error
 */
extern const char * ffs_errstr(void);

/*!
 * @brief Open the file named 'path' and check the @em FFS partition table
 *        at 'offset' bytes from the beginning of the file (or device).
 * @memberof ffs
 * @param path [in] Path of target file or device
 * @param offset [in] Byte offset, from the beginning of the file (or device),
 *        of the ffs_hdr_t structure
 * @return 0 on success, non-0 otherwise
 *         FFS_CHECK_PATH if path is NULL or ""
 *         FFS_CHECK_HEADER_MAGIC for header magic corruption
 *         FFS_CHECK_HEADER_CHECKSUM for header data corruption
 *         FFS_CHECK_ENTRY_CHECKSUM for [an] entry data corruption
 */
extern int ffs_check(const char *, off_t)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

/*!
 * @brief Create an empty @em FFS partition table and store it at 'offset'
 *        bytes from the beginning of the file (or device) named 'path'
 * @memberof ffs
 * @param path [in] Path of target file or device
 * @param offset [in] Byte offset, from beginning of file (or device),
 *        of the ffs_hdr_t structure
 * @param block_size [in] Block size in bytes
 * @param block_count [in] Number of blocks in the device size
 * @return Pointer to ffs_t (allocated on the heap) on success,
 *         NULL otherwise
 */
extern ffs_t * ffs_create(const char *, off_t, uint32_t, uint32_t)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

/*!
 * @brief Open the file name 'path' and read the @em FFS partition table
 *        at 'offset' bytes from the beggining of the file (or device).
 * @memberof ffs
 * @param path [in] Path of target file or device
 * @param offset [in] Byte offset, from beginning of file (or device),
 *        of the ffs_hdr_t structure
 * @return Pointer to ffs_t (allocated on the heap) on success,
 *         NULL otherwise
 */
extern ffs_t * ffs_open(const char *, off_t)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

/*!
 * @brief Query a @em FFS object for header metadata.
 * @memberof ffs
 * @param self [in] Pointer to ffs object
 * @param name [in] Field number, see ffs.h for details
 * 		FFS_INFO_MAGIC - ffs_hdr::magic
 *		FFS_INFO_VERSION - ffs_hdr::version
 * 		FFS_INFO_ENTRY_SIZE - ffs_hdr::entry_size
 * 		FFS_INFO_ENTRY_COUNT - ffs_hdr::entry_count
 * 		FFS_INFO_BLOCK_SIZE - ffs_hdr::block_size
 * 		FFS_INFO_BLOCK_COUNT - ffs_hdr::block_count
 * @param value [out] Pointer to output data
 * @return '0' on success, non-0 otherwise
 */
extern int ffs_info(ffs_t *, int, uint32_t *)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

/*!
 * @brief Close a @em FFS object, after writting any modified in-core
 *        data to the file (or device).
 * @memberof ffs
 * @param self [in] Pointer to ffs object
 * @return '0' on success, non-0 otherwise
 */
extern int ffs_close(ffs_t *)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

/*!
 * @brief Flush all modified in-core metadata of a @em FFS object,
 *        to the underlying file (or device).
 * @memberof ffs
 * @param self [in] Pointer to ffs object
 * @return '0' on success, non-0 otherwise
 */
extern int ffs_fsync(ffs_t *)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

/*!
 * @brief Pretty print the entries of a @em FFS partition table to
 *        stream 'out'
 * @memberof ffs
 * @param self [in] Pointer to ffs object
 * @param out [in] Pointer the output stream
 * @return '0' on success, non-0 otherwise
 */
extern int ffs_list_entries(ffs_t *, FILE *)
/*! @cond */ __nonnull ((1)) /*! @endcond */ ;

/*!
 * @brief Iterate over the entries of a @em FFS partition table and
 *        call a callback function 'func'
 * @note If the callback function returns non-0, the iteration function
 *       will return immediately.
 * @memberof ffs
 * @param self [in] Pointer to ffs object
 * @param func [in] Pointer the callback function
 * @return '0' on success, non-0 otherwise
 */
extern int ffs_iterate_entries(ffs_t *, int (*)(ffs_entry_t*))
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

/*!
 * @brief Find an entry in a @em FFS partition table and return
 *        a copy of the in 'entry'
 * @memberof ffs
 * @param self [in] Pointer to ffs object
 * @param path [in] Name of a partition entry
 * @param entry [out] Target entry object
 * @return '1' == found, '0' == not-found, error otherwise
 */
extern int ffs_entry_find(ffs_t *, const char *, ffs_entry_t *)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

/*!
 * @brief Find the parent entry in a @em FFS partition table and
 *        return a copy of the in 'entry'
 * @memberof ffs
 * @param self [in] Pointer to ffs object
 * @param path [in] Name of a partition entry
 * @param parent [out] Target entry object
 * @return '1' == found, '0' == not-found, error otherwise
 */
extern int ffs_entry_find_parent(ffs_t *, const char *, ffs_entry_t *)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

/*!
 * @brief Add a partition entry to a @em FFS partition table
 * @memberof ffs
 * @param self [in] Pointer to ffs object
 * @param path [in] Name of a partition entry
 * @param offset [in] Offset, in blocks, of the partition entry, from
 *        the beginning of the file (or device)
 * @param size [in] Size, in blocks, of the partition entry
 * @param type [in] Partition type.  FFS_TYPE_LOGICAL can be used to 
 *        to create a container for a set of partitions.  A logical partition
 *        can be thought of as a directory.  Use FFS_TYPE_DATA for data
 *        partitions.
 * @param flags [in] Partition flags.  FFS_FLAG_PROTECTED can be used to
 *        protect a partition from inadvertant updates.  The ffs related tools
 *        should *not* overwrite protected partitions unless a --protected
 *        flag is specified.
 * @return '0' on success, non-0 otherwise
 */
extern int ffs_entry_add(ffs_t *, const char *, off_t, size_t, ffs_type_t, uint32_t)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

/*!
 * @brief Delete a partition entry from the @em FFS partition table
 * @memberof ffs
 * @param self [in] Pointer to an ffs object
 * @param path [in] Name of a partition entry
 * @return '0' on success, non-0 otherwise
 */
extern int ffs_entry_delete(ffs_t *, const char *)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

/*!
 * @brief Get the value of a meta-data user word
 * @memberof ffs
 * @param self [in] Pointer to an ffs object
 * @param name [in] Name of a partition entry
 * @param user [in] User word number, in the range [0..FFS_USER_WORDS]
 * @param word [in] User word value, in the range [0..UINT32_MAX] (optional)
 * @note 'word' is optional, if omitted, the current value is returned
 * @return '0' on success, non-0 otherwise
 */
extern int ffs_entry_user_get(ffs_t *, const char *, uint32_t, uint32_t *)
/*! @cond */ __nonnull ((1,2,4)) /*! @endcond */ ;

/*!
 * @brief Set the value of a meta-data user word
 * @memberof ffs
 * @param self [in] Pointer to an ffs object
 * @param name [in] Name of a partition entry
 * @param user [in] User word number, in the range [0..FFS_USER_WORDS]
 * @param word [in] User word value, in the range [0..UINT32_MAX] (optional)
 * @note 'word' is optional, if omitted, the current value is returned
 * @return '0' on success, non-0 otherwise
 */
extern int ffs_entry_user_put(ffs_t *, const char *, uint32_t, uint32_t)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

/*!
 * @brief Hexdump the data contents of a partition entry to output stream
 *        'out'
 * @memberof ffs
 * @param self [in] Pointer to an ffs object
 * @param name [in] Name of a partition entry
 * @param out [in] Output stream
 * @return Negative on success, number of bytes written otherwise
 */
extern ssize_t ffs_entry_hexdump(ffs_t *, const char *, FILE *)
/*! @cond */ __nonnull ((1,2,3)) /*! @endcond */ ;

/*!
 * @brief Change the actual size of partition entry 'name'
 *        to 'offset' bytes from the beginning of the entry.
 * @memberof ffs
 * @param self [in] Pointer to an ffs object
 * @param name [in] Name of a partition entry
 * @param offset [in] Offset from the beginning of the partition (in bytes)
 * @return Negative on failure, zero otherwise
 */
extern ssize_t ffs_entry_truncate_no_pad(ffs_t *, const char *, off_t)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

/*!
 * @brief Change the actual size of partition entry 'name'
 *        to 'offset' bytes from the beginning of the entry.
 * @memberof ffs
 * @param self [in] Pointer to an ffs object
 * @param name [in] Name of a partition entry
 * @param offset [in] Offset from the beginning of the partition (in bytes)
 * @param pad [in] Pad character
 * @return Negative on failure, zero otherwise
 */
extern ssize_t ffs_entry_truncate(ffs_t *, const char *, off_t, uint8_t)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

/*!
 * @brief Read 'count' data bytes of partition entry 'name' to into 'buf' at
 *        offset 'offset' bytes from the beginning of the entry
 * @memberof ffs
 * @param self [in] Pointer to an ffs object
 * @param name [in] Name of a partition entry
 * @param buf [out] Output data buffer
 * @param offset [in] Offset from the beginning of the partition
 * @param count [in] Number of bytes to read
 * @return Negative on failure, number of bytes read otherwise
 */
extern ssize_t ffs_entry_read(ffs_t *, const char *, void *, off_t, size_t)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

/*!
 * @brief Write 'count' data bytes to partition entry 'name' at offset
 *        'offset' bytes from the beginning of the entry, from input
 *        buffer 'buf'
 * @memberof ffs
 * @param self [in] Pointer to an ffs object
 * @param name [in] Name of a partition entry
 * @param buf [in] Input data buffer
 * @param offset [in] Offset from the beginning of the partition
 * @param count [in] Number of bytes to write
 * @return Negative on failure, else number of bytes written otherwise
 */
extern ssize_t ffs_entry_write(ffs_t *, const char *, const void *, off_t, size_t)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

/*!
 * @brief Return an array of entry_t structures, one each partition that
 * 	exists in the partition table
 * @memberof ffs
 * @param self [in] Pointer to an ffs object
 * @param list [out] Array of entry_t structures allocated on the heap
 * @return Negative on failure, else number of entry_t's in the output
 *	array with a call to free()
 * @note Caller is responsible for freeing array with a call to free()
 */
extern ssize_t ffs_entry_list(ffs_t *, ffs_entry_t **)
/*! @cond */ __nonnull ((1,2)) /*! @endcond */ ;

#ifdef __cplusplus
}
#endif

#endif /* __LIBFFS2_H__ */
