/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fcp/src/main.h $                                              */
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

#ifndef __FCP_H__
#define __FCP_H__

#include <sys/types.h>

#include <stdio.h>

#include <ffs/libffs.h>

#define FCP_MAJOR	0x01
#define FCP_MINOR	0x00
#define FCP_PATCH	0x00

#define TYPE_FILE	"file"
#define TYPE_RW		"rw"
#define TYPE_AA		"aa"
#define TYPE_SFC	"sfc"

#define verbose(fmt, args...) \
	({if (verbose) printf("%s: " fmt, __func__, ##args); })
#define debug(fmt, args...) \
	({if (debug) printf("%s: " fmt, __func__, ##args); })

typedef enum {
	c_ERROR = 0,
	c_PROBE = 'P',
	c_LIST = 'L',
	c_READ = 'R',
	c_WRITE = 'W',
	c_ERASE = 'E',
	c_COPY = 'C',
	c_TRUNC = 'T',
	c_COMPARE = 'M',
	c_USER = 'U',
} cmd_t;

typedef enum {
	o_ERROR = 0,
	o_OFFSET = 'o',
	o_BUFFER = 'b',
} option_t;

typedef enum {
	f_ERROR = 0,
	f_FORCE = 'f',
	f_PROTECTED = 'p',
	f_VERBOSE = 'v',
	f_DEBUG = 'd',
	f_HELP = 'h',
} flag_t;

typedef struct {
	const char *short_name;

	cmd_t cmd;
	char *src_type, *src_target, *src_name;
	char *dst_type, *dst_target, *dst_name;

	/* options */
	const char *offset;

	/* flags */
	flag_t force;
	flag_t protected;
	flag_t verbose;
	flag_t debug;

	const char **opt;
	int opt_sz, opt_nr;
} args_t;

extern args_t args;
extern size_t page_size;

extern int verbose;
extern int debug;

extern int fcp_read_entry(ffs_t *, const char *, FILE *);
extern int fcp_write_entry(ffs_t *, const char *, FILE *);
extern int fcp_erase_entry(ffs_t *, const char *, char);
extern int fcp_copy_entry(ffs_t *, const char *, ffs_t *, const char *);
extern int fcp_compare_entry(ffs_t *, const char *, ffs_t *, const char *);

extern int command_probe(args_t *);
extern int command_list(args_t *);
extern int command_read(args_t *);
extern int command_write(args_t *);
extern int command_erase(args_t *);
extern int command_copy_compare(args_t *);
extern int command_trunc(args_t *);
extern int command_compare(args_t *);
extern int command_user(args_t *);

#endif /* __FCP_H__ */
