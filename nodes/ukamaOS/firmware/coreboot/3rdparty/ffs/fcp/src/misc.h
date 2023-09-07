/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fcp/src/misc.h $                                              */
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
 *    File: fcp_misc.h
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: Misc. helpers
 *    Date: 05/12/2012
 */

#ifndef __MISC__H__
#define __MISC__H__

#include <sys/types.h>

#include <stdio.h>
#include <regex.h>

#include <clib/list.h>

#include <ffs/libffs.h>

typedef struct entry_list entry_list_t;
struct entry_list {
	list_t list;
	ffs_t * ffs;
};

typedef struct entry_node entry_node_t;
struct entry_node {
	list_node_t node;
	ffs_entry_t entry;
};

extern entry_list_t * entry_list_create(ffs_t *);
extern entry_list_t * entry_list_create_by_regex(ffs_t *, const char *);
extern int entry_list_add(entry_list_t *, ffs_entry_t *);
extern int entry_list_add_child(entry_list_t *, ffs_entry_t *);
extern int entry_list_remove(entry_list_t *, entry_node_t *);
extern int entry_list_delete(entry_list_t *);
extern int entry_list_exists(entry_list_t *, ffs_entry_t *);
extern ffs_entry_t * entry_list_find(entry_list_t *, const char *);
extern int entry_list_dump(entry_list_t *, FILE *);

extern int parse_offset(const char *, off_t *);
extern int parse_size(const char *, uint32_t *);
extern int parse_number(const char *, uint32_t *);
extern int parse_path(const char *, char **, char **, char **);

extern int dump_errors(const char *, FILE *);
extern int check_file(const char *, FILE *, off_t);
extern int is_file(const char *, const char *, const char *);
extern int valid_type(const char *);

extern FILE *__fopen(const char *, const char *, const char *, int);

#endif /* __MISC__H__ */
