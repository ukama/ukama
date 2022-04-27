/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fcp/src/cmd_copy.c $                                          */
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
 *    File: cmd_copy.c
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: copy & compare implementation
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

static int validate_files(ffs_t * src, ffs_t * dst)
{
	assert(src != NULL);
	assert(dst != NULL);

	uint32_t s, d;

	if (__ffs_info(src, FFS_INFO_BLOCK_SIZE, &s) < 0)
		return -1;
	if (__ffs_info(dst, FFS_INFO_BLOCK_SIZE, &d) < 0)
		return -1;
	if (s != d) {
		UNEXPECTED("source '%s' and destination '%s' block "
			   "size differs, use --force to overwrite\n",
			   src->path, dst->path);
		return -1;
	}

	if (__ffs_info(src, FFS_INFO_BLOCK_COUNT, &s) < 0)
		return -1;
	if (__ffs_info(dst, FFS_INFO_BLOCK_COUNT, &d) < 0)
		return -1;
	if (s != d) {
		UNEXPECTED("source '%s' and destination '%s' block "
			   "count differs, use --force to overwrite\n",
			   src->path, dst->path);
		return -1;
	}

	if (__ffs_info(src, FFS_INFO_ENTRY_SIZE, &s) < 0)
		return -1;
	if (__ffs_info(dst, FFS_INFO_ENTRY_SIZE, &d) < 0)
		return -1;
	if (s != d) {
		UNEXPECTED("source '%s' and destination '%s' entry "
			   "size differs, use --force to overwrite\n",
			   src->path, dst->path);
		return -1;
	}

	if (__ffs_info(src, FFS_INFO_VERSION, &s) < 0)
		return -1;
	if (__ffs_info(dst, FFS_INFO_VERSION, &d) < 0)
		return -1;
	if (s != d) {
		UNEXPECTED("source '%s' and destination '%s' version "
			   "differs, use --force to overwrite\n",
			   src->path, dst->path);
		return -1;
	}

	return 0;
}

static int __copy_entry(args_t * args,
			ffs_t * src_ffs, ffs_entry_t * src_entry,
			ffs_t * dst_ffs, ffs_entry_t * dst_entry,
			entry_list_t * done_list)
{
	char full_src_name[page_size];
	if (__ffs_entry_name(src_ffs, src_entry, full_src_name,
		     sizeof full_src_name) < 0) {
		return -1;
	}

	char full_dst_name[page_size];
	if (__ffs_entry_name(src_ffs, dst_entry, full_dst_name,
			     sizeof full_dst_name) < 0) {
		return -1;
	}

	if (dst_entry->type == FFS_TYPE_LOGICAL) {
		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: %s: logical partition (skip)\n",
			        (long long)dst_ffs->offset, full_dst_name);
		return 0;
	}

	if (src_entry->type == FFS_TYPE_LOGICAL) {
		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: %s: logical partition (skip)\n",
			        (long long)src_ffs->offset, full_src_name);
		return 0;
	}

	if (dst_entry->type == FFS_TYPE_PARTITION) {
		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: %s: partition table (skip)\n",
			        (long long)dst_ffs->offset, full_dst_name);
		return 0;
	}

	if (src_entry->type == FFS_TYPE_PARTITION) {
		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: %s: partition table (skip)\n",
			        (long long)src_ffs->offset, full_src_name);
		return 0;
	}

	if (args->protected != f_PROTECTED) {
		if (dst_entry->flags & FFS_FLAGS_PROTECTED) {
			if (args->verbose == f_VERBOSE)
				fprintf(stderr, "%8llx: %s: protected "
					"partition (skip)\n", (long long)dst_ffs->offset,
					full_dst_name);
			return 0;
		}
		if (src_entry->flags & FFS_FLAGS_PROTECTED) {
			if (args->verbose == f_VERBOSE)
				fprintf(stderr, "%8llx: %s: protected "
					"partition (skip)\n", (long long)src_ffs->offset,
				        full_src_name);
			return 0;
		}
	}

	if (src_entry->size != dst_entry->size) {
		UNEXPECTED("%8llx: %s: source partition '%s' size mismatch "
			   "'%x', use --force to overwrite\n", (long long)dst_ffs->offset,
			   full_dst_name, full_src_name, dst_entry->size);
		return -1;
	}

	if (src_entry->type != dst_entry->type) {
		UNEXPECTED("%8llx: %s: source partition '%s' type mismatch "
			   "'%x', use --force to overwrite\n", (long long)dst_ffs->offset,
			   full_dst_name, full_src_name, dst_entry->size);
		return -1;
	}

	if (__ffs_entry_truncate(dst_ffs, full_dst_name,
				 src_entry->actual) < 0) {
		ERRNO(errno);
		return -1;
	}
	if (args->verbose == f_VERBOSE)
		fprintf(stderr, "%8llx: %s: trunc size '%x' (done)\n",
			(long long)dst_ffs->offset, full_dst_name, src_entry->actual);

	uint32_t src_val, dst_val;
	for (uint32_t i=0; i<FFS_USER_WORDS; i++) {
		__ffs_entry_user_get(src_ffs, full_src_name, i, &src_val);
		__ffs_entry_user_get(dst_ffs, full_dst_name, i, &dst_val);

                if (args->force != f_FORCE && i == USER_DATA_VOL)
                        continue;

		if (src_val != dst_val) {
			if (__ffs_entry_user_put(dst_ffs, full_dst_name,
					         i, src_val) < 0)
				return -1;
		}
	}
	if (args->verbose == f_VERBOSE)
		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: %s: copy user[] from '%s' "
				"(done)\n", (long long)dst_ffs->offset, full_dst_name,
				src_ffs->path);

	if (entry_list_exists(done_list, src_entry) == 1) {
		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: %s: copy from '%s' (skip)\n",
		       		(long long)dst_ffs->offset, full_dst_name, src_ffs->path);
		return 0;
	}

	if (entry_list_add(done_list, src_entry) < 0)
		return -1;

	if (fcp_copy_entry(src_ffs, full_src_name, dst_ffs, full_dst_name) < 0)
		return -1;

	if (args->verbose == f_VERBOSE)
		fprintf(stderr, "%8llx: %s: copy from '%s' (done)\n",
		        (long long)dst_ffs->offset, full_dst_name, src_ffs->path);

	return 0;
}

static int __compare_entry(args_t * args,
			   ffs_t * src_ffs, ffs_entry_t * src_entry,
			   ffs_t * dst_ffs, ffs_entry_t * dst_entry,
			   entry_list_t * done_list)
{
	char full_src_name[page_size];
	if (__ffs_entry_name(src_ffs, src_entry, full_src_name,
		     sizeof full_src_name) < 0) {
		return -1;
	}

	char full_dst_name[page_size];
	if (__ffs_entry_name(src_ffs, dst_entry, full_dst_name,
			     sizeof full_dst_name) < 0) {
		return -1;
	}

	if (dst_entry->type == FFS_TYPE_LOGICAL) {
		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: %s: logical partition (skip)\n",
			        (long long)dst_ffs->offset, full_dst_name);
		return 0;
	}

	if (src_entry->type == FFS_TYPE_LOGICAL) {
		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: %s: logical partition (skip)\n",
			        (long long)src_ffs->offset, full_src_name);
		return 0;
	}

	if (args->protected != f_PROTECTED) {
		if (dst_entry->flags & FFS_FLAGS_PROTECTED) {
			if (args->verbose == f_VERBOSE)
				fprintf(stderr, "%8llx: %s: protected "
					"partition (skip)\n", (long long)dst_ffs->offset,
					full_dst_name);
			return 0;
		}
		if (src_entry->flags & FFS_FLAGS_PROTECTED) {
			if (args->verbose == f_VERBOSE)
				fprintf(stderr, "%8llx: %s: protected "
					"partition (skip)\n", (long long)src_ffs->offset,
				       full_src_name);
			return 0;
		}
	}

	if (src_entry->size != dst_entry->size) {
		UNEXPECTED("%8llx: %s: source partition '%s' size mismatch "
			   "'%x', use --force to overwrite\n", (long long)dst_ffs->offset,
			   full_dst_name, full_src_name, dst_entry->size);
		return -1;
	}

	if (src_entry->type != dst_entry->type) {
		UNEXPECTED("%8llx: %s: source partition '%s' type mismatch "
			   "'%x', use --force to overwrite\n", (long long)dst_ffs->offset,
			   full_dst_name, full_src_name, dst_entry->size);
		return -1;
	}

	if (entry_list_exists(done_list, src_entry) == 1) {
		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: %s: compare from '%s' (skip)\n",
				(long long)dst_ffs->offset, full_dst_name, src_ffs->path);
		return 0;
	}

	if (entry_list_add(done_list, src_entry) < 0)
		return -1;

	if (fcp_compare_entry(src_ffs, full_src_name,
			      dst_ffs, full_dst_name) < 0)
		return -1;

	if (args->verbose == f_VERBOSE)
		fprintf(stderr, "%8llx: %s: compare from '%s' (done)\n",
			(long long)dst_ffs->offset, full_dst_name, src_ffs->path);

	return 0;
}

static int __force_part(ffs_t * src, FILE * dst)
{
	assert(src != NULL);
	assert(dst != NULL);

	uint32_t block_size;
	if (__ffs_info(src, FFS_INFO_BLOCK_SIZE, &block_size) < 0)
		return -1;

	char part[block_size];
	memset(part, 0, block_size);

	if (__ffs_entry_read(src, FFS_PARTITION_NAME, part,
			     0, block_size) < 0)
		return -1;

	uint32_t offset;
	if (__ffs_info(src, FFS_INFO_OFFSET, &offset) < 0)
		return -1;

	if (fseek(dst, offset, SEEK_SET) < 0) {
		ERRNO(errno);
		return -1;
	}

	if (fwrite(part, 1, block_size, dst) != block_size) {
		if (ferror(dst)) {
			ERRNO(errno);
			return -1;
		}
	}

	if (fflush(dst) == EOF) {
		ERRNO(errno);
		return -1;
	}

	return 0;
}

static int __copy_compare(args_t * args, off_t offset, entry_list_t * done_list)
{
	assert(args != NULL);
	assert(done_list != NULL);

	char * src_type = args->src_type;
	char * src_target = args->src_target;
	char * src_name = args->src_name;

	char * dst_type = args->dst_type;
	char * dst_target = args->dst_target;
	char * dst_name = src_name;

	if (args->dst_name != NULL && *args->dst_name != '\0')
		dst_name = args->dst_name;

	if (src_name == NULL)
		src_name = "*";
	if (dst_name == NULL)
		dst_name = "*";

	RAII(FILE*, src_file, __fopen(src_type, src_target, "r", debug),
	     fclose);
	if (src_file == NULL)
		return -1;
	if (check_file(src_target, src_file, offset) < 0)
		return -1;
	RAII(ffs_t*, src_ffs, __ffs_fopen(src_file, offset), __ffs_fclose);
	if (src_ffs == NULL)
		return -1;

	src_ffs->path = basename(src_target);
	done_list->ffs = src_ffs;

	RAII(FILE*, dst_file, __fopen(dst_type, dst_target, "r+", debug),
	     fclose);
	if (dst_file == NULL)
		return -1;

	if (args->force == f_FORCE && args->cmd == c_COPY) {
		if (__force_part(src_ffs, dst_file) < 0)
			return -1;

		if (args->verbose == f_VERBOSE)
			fprintf(stderr, "%8llx: partition table '%s' => '%s' "
				"(done)\n", (long long)offset, src_target, dst_target);
	}

	if (check_file(dst_target, dst_file, offset) < 0)
		return -1;

	RAII(ffs_t*, dst_ffs, __ffs_fopen(dst_file, offset), __ffs_fclose);
	if (dst_ffs == NULL)
		return -1;

	dst_ffs->path = basename(dst_target);

	if (validate_files(src_ffs, dst_ffs) < 0)
		return -1;

	if (src_ffs->count <= 0)		// fix me
		return 0;

	uint32_t block_size;
	if (__ffs_info(src_ffs, FFS_INFO_BLOCK_SIZE, &block_size) < 0)
		return -1;

	ffs_entry_t src_parent;
	ffs_entry_t dst_parent;

	if (strcmp(src_name, "*") == 0 && strcmp(dst_name, "*") == 0) {
		src_parent.type = FFS_TYPE_LOGICAL;
		src_parent.id = -1;
		dst_parent.type = FFS_TYPE_LOGICAL;
		dst_parent.id = -1;
	} else if ((src_name != NULL) && (dst_name != NULL)) {
		if (__ffs_entry_find(src_ffs, src_name, &src_parent) == false) {
			UNEXPECTED("%8llx: partition entry '%s' not found in "
				   "'%s'\n", (long long)src_ffs->offset, src_name,
				   src_target);
			return -1;
		}

		if (__ffs_entry_find(dst_ffs, dst_name, &dst_parent) == false) {
			UNEXPECTED("%8llx: partition entry '%s' not found in "
				   "'%s'\n", (long long)dst_ffs->offset, dst_name,
				   dst_target);
			return -1;
		}
	}

	if (src_parent.type == FFS_TYPE_DATA &&
	    dst_parent.type == FFS_TYPE_DATA) {
		if (args->cmd == c_COPY)
			return __copy_entry(args, src_ffs, &src_parent,
					    dst_ffs, &dst_parent, done_list);
		else
			return __compare_entry(args, src_ffs, &src_parent,
					       dst_ffs, &dst_parent, done_list);
	} else if (src_parent.type == FFS_TYPE_LOGICAL &&
		   dst_parent.type == FFS_TYPE_LOGICAL) {

		RAII(entry_list_t*, src_list, entry_list_create(src_ffs),
		     entry_list_delete);
		if (src_list == NULL)
			return -1;
		if (entry_list_add_child(src_list, &src_parent) < 0)
			return -1;

		RAII(entry_list_t*, dst_list, entry_list_create(dst_ffs),
		     entry_list_delete);
		if (dst_list == NULL)
			return -1;
		if (entry_list_add_child(dst_list, &dst_parent) < 0)
			return -1;

		while (list_empty(&src_list->list) == 0) {
			RAII(entry_node_t*, src_node,
			     container_of(list_head(&src_list->list),
					  entry_node_t, node), free);

			entry_list_remove(src_list, src_node);
			ffs_entry_t * src_entry = &src_node->entry;

			char full_src_name[page_size];
			if (__ffs_entry_name(src_ffs, src_entry, full_src_name,
					     sizeof full_src_name) < 0) {
				free(src_node);
				return -1;
			}

			if (list_empty(&dst_list->list)) {
				UNEXPECTED("source entry '%s' not a child of "
					   "destination entry '%s'\n",
					   full_src_name, dst_name);
				return -1;
			}

			RAII(entry_node_t*, dst_node,
			     container_of(list_head(&dst_list->list),
					  entry_node_t, node), free);

			entry_list_remove(dst_list, dst_node);
			ffs_entry_t * dst_entry = &dst_node->entry;

			char full_dst_name[page_size];
			if (__ffs_entry_name(dst_ffs, dst_entry, full_dst_name,
					     sizeof full_dst_name) < 0)
				return -1;

			if (src_entry->type == FFS_TYPE_LOGICAL) {
				if (entry_list_add_child(src_list,
							 src_entry) < 0)
					return -1;
			}

			if (dst_entry->type == FFS_TYPE_LOGICAL) {
				if (entry_list_add_child(dst_list,
							 dst_entry) < 0)
					return -1;
			}

			if (args->cmd == c_COPY) {
				if (__copy_entry(args, src_ffs, src_entry,
						 dst_ffs, dst_entry,
						 done_list) < 0)
					return -1;
			} else if (args->cmd == c_COMPARE) {
				if (__compare_entry(args, src_ffs, src_entry,
						 dst_ffs, dst_entry,
						 done_list) < 0) {
					return -1;
				}
			}
		}
	}

	return 0;
}

int command_copy_compare(args_t * args)
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

		rc = __copy_compare(args, offset, done_list);
		if (rc < 0)
			break;

		if (*end == '\0')
			break;
		end++;
	}

	return rc;
}
