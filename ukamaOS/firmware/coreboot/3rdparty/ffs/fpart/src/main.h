/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fpart/src/main.h $                                            */
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
 *    File: main.h
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: cmdline tool for libffs.so
 *    Date: 05/12/2012
 */

#ifndef __MAIN_H__
#define __MAIN_H__

#include <stdio.h>

#include <ffs/libffs.h>

#define FPART_MAJOR	0x01
#define FPART_MINOR	0x00
#define FPART_PATCH	0x00

typedef enum {
	c_ERROR = 0,
	c_CREATE = 'C',
	c_ADD = 'A',
	c_DELETE = 'D',
	c_ERASE = 'E',
	c_LIST = 'L',
	c_TRUNC = 'T',
	c_USER = 'U',
} cmd_t;

typedef enum {
	o_ERROR = 0,
	o_POFFSET = 'p',
	o_TARGET = 't',
	o_NAME = 'n',
	o_OFFSET = 'o',
	o_SIZE = 's',
	o_BLOCK = 'b',
	o_VALUE = 'u',
	o_FLAGS = 'g',
	o_PAD = 'a',
} option_t;

typedef enum {
	f_ERROR = 0,
	f_FORCE = 'f',
	f_PROTECTED = 'r',
	f_LOGICAL = 'l',
	f_VERBOSE = 'v',
	f_DEBUG = 'd',
	f_HELP = 'h',
} flag_t;

typedef struct {
	char *short_name;

	/* target */
	const char *path;

	/* command */
	cmd_t cmd;

	/* options */
	char *name, *target;
	char *offset, *poffset;
	char *size, *block;
	char *user, *value;
	char *flags, *pad;

	/* flags */
	flag_t force, logical;
	flag_t verbose, debug;
	flag_t protected;

	const char **opt;
	int opt_sz, opt_nr;
} args_t;

typedef struct ffs_entry_node ffs_entry_node_t;
struct ffs_entry_node {
	list_node_t node;
	ffs_entry_t entry;
};

extern args_t args;
extern size_t page_size;

extern int parse_offset(const char *, off_t *);
extern int parse_size(const char *, uint32_t *);
extern int parse_number(const char *, uint32_t *);

extern bool check_extension(const char *, const char *);
extern int create_regular_file(const char *, size_t, char);
extern FILE *fopen_generic(const char *, const char *, int);

extern int command(args_t *, int (*)(args_t *, off_t));
extern int verify_operation(const char *, ffs_t *, ffs_entry_t *,
				          ffs_t *, ffs_entry_t *);

extern int command_create(args_t *);
extern int command_add(args_t *);
extern int command_delete(args_t *);
extern int command_list(args_t *);
extern int command_trunc(args_t *);
extern int command_erase(args_t *);
extern int command_user(args_t *);

#endif /* __MAIN_H__ */
