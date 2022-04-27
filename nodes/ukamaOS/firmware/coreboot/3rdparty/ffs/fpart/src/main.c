/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fpart/src/main.c $                                            */
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
 *    File: main.c
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: FFS partition tool
 *    Date: 05/12/2012
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

args_t args;
size_t page_size;

static void usage(bool verbose)
{
	FILE *e = stderr;

	fprintf(e,
		"fpart - FFS Partition Tool v%d.%d.%d -- Authors: "
		"<shaun@us.ibm.com>\n", FPART_MAJOR, FPART_MINOR, FPART_PATCH);

	fprintf(e, "\nUsage:\n");
	fprintf(e, "  fpart <command> <options>...\n");

	fprintf(e, "\nExamples:\n");
	fprintf(e, "  fpart -C -t nor -s 64MiB -b 64k -p 0x3f0000,0x7f0000\n");
	fprintf(e, "  fpart -A -t nor -p 0x3f0000,0x7f0000 -s 1Mb -o 0M -g 0 "
		   "-n boot0 -l\n");
	fprintf(e, "  fpart --add --target nor --size 1mb --offset 1MiB "
		   "--flags 0x0 --name boot0/ipl\n");
	fprintf(e, "  fpart --delete --target nor nor --name boot0/ipl\n");
	fprintf(e, "  fpart --target nor --write ipl.bin --name boot1/ipl\n");
	fprintf(e, "  fpart --user 0 -t nor -n boot0/ipl --value 0xFF500FF5\n");
	fprintf(e, "  fpart --copy new_nor -t nor -n ipl\n");
	fprintf(e, "  fpart --compare new_nor -t nor -n bank0\n");

	fprintf(e, "\nCommands:\n");
	fprintf(e, "  -C, --create         [options]\n");
	if (verbose)
		fprintf(e, "\n  Create an empty partition table at each"
			" specified partition offest.\n\n");

	fprintf(e, "  -A, --add            [options]\n");
	if (verbose)
		fprintf(e, "\n  Add partition entry(s) to each specified"
			" partition offset.\n\n");

	fprintf(e, "  -D, --delete         [options]\n");
	if (verbose)
		fprintf(e, "\n  Delete partition entry(s) from each specified"
			" partition offset.\n\n");

	fprintf(e, "  -E, --erase          [options]\n");
	if (verbose)
		fprintf(e, "\n  Fill partition entry(s) from each specified"
			" partition offset.\n\n");

	fprintf(e, "  -L, --list           [options]\n");
	if (verbose)
		fprintf(e, "\n  Display the list of matching partition"
			" entry(s) of each specified\n"
			"  partition offset.\n\n");

	fprintf(e, "  -T, --trunc          [options]\n");
	if (verbose)
		fprintf(e, "\n  Truncate the actual data size of matching"
			" partition entry(s) of each\n  specified partition"
			" offset.\n\n");

	fprintf(e, "  -U, --user    <num>  [options]\n");
	if (verbose)
		fprintf(e, "\n  Read or write user words of matching partition"
			" entry(s) for each\n  specified partition offset.\n");

	/* =============================== */

	fprintf(e, "\nOptions:\n");
	fprintf(e, "  -p, --partition-offset <offset[,offset]>\n");
	if (verbose)
		fprintf(e, "\n  Specifies a comma (,) separated list of"
			" partition table offsets, in\n  bytes from the start"
			" of the target file (or device).\n\n");

	fprintf(e, "  -t, --target           <target>\n");
	if (verbose) {
		fprintf(e, "\n  Specifies the target of the operation.");
		fprintf(e, "  <target> can be one of:\n");
		fprintf(e, "\t    <port> : Aardvard USB port number [0..9]\n");
		fprintf(e, "\t<hostname> : RISCWatch probe hostname\n");
		fprintf(e, "\t    <path> : SFC device file, e.g. "
			"/dev/mtdblock/sfc.of\n\n");
	}

	fprintf(e, "  -n, --name             <name>\n");
	if (verbose)
		fprintf(e, "\n  Specifies the name of a partition entry(s)."
			"  For some commands\n  (--hexdump, --add)"
			" <name> must be a fully qualified name.\n  For other"
			" commands, name can specify a regular "
			"expression.\n\n");

	fprintf(e, "  -o, --offset           <offset>\n");
	if (verbose)
		fprintf(e, "\n  Specifies the offset of a partition entry,"
			" in bytes from the beginning\n  of the target file"
			" (or device).\n\n");

	fprintf(e, "  -s, --size             <size>\n");
	if (verbose)
		fprintf(e, "\n  Specifies the size of a partition entry (in"
			" bytes).\n  <size> is a decimal (or hex) number"
			" (optionally) followed by 'k', 'K',\n  'KB' or 'KiB'"
			" to specify.\n\n");

	fprintf(e, "  -b, --block-size       <size>\n");
	if (verbose)
		fprintf(e, "\n  Specifies the block size, in bytes.  <size>"
			" is a decimal (or hex)\n  number (optionally)"
			" followed by 'k', 'K', 'KB' or 'KiB' to"
			" specify.\n\n");

	fprintf(e, "  -u, --value            <value>\n");
	if (verbose)
		fprintf(e, "\n  Specifies the user word value.  <value> is a"
			" decimal (or hex) number.\n\n");

	fprintf(e, "  -g, --flags            <value>\n");
	if (verbose)
		fprintf(e, "\n  Specifies the partition flags value."
			"  <value> is a decimal (or hex)\n  number.\n\n");

	fprintf(e, "  -a, --pad              <value>\n");
	if (verbose)
		fprintf(e,
			"\n  Specifies the partition content initial value."
			"  <value> is a decimal\n  (or hex) number, default is"
			" 0xFF.\n\n");

	/* =============================== */

	fprintf(e, "\nFlags:\n");
	fprintf(e, "  -f, --force\n");
	if (verbose)
		fprintf(e, "\n  Override command safe guards.\n\n");

	fprintf(e, "  -l, --logical\n");
	if (verbose)
		fprintf(e, "\n  Specifies the partition entry is a logical"
			" partition instead of a\n  data partition.\n\n");

	fprintf(e, "  -v, --verbose\n");
	if (verbose)
		fprintf(e, "\n  Write progress messages to stdout\n\n");

	fprintf(e, "  -d, --debug\n");
	if (verbose)
		fprintf(e, "\n  Write debug messages to stdout\n\n");

	fprintf(e, "  -h, --help\n");
	if (verbose)
		fprintf(e, "\n  Write this help text to stdout and exit\n");

	fprintf(e, "\n");

	fprintf(e, "Report bugs to <https://bugzilla.linux.ibm.com/>"
		" (vendor='MCP for FSP*')\n");
}

static int process_argument(args_t * args, int opt, const char *optarg)
{
	assert(args != NULL);

	switch (opt) {
	case c_CREATE:		/* create */
	case c_ADD:		/* add */
	case c_DELETE:		/* delete */
	case c_LIST:		/* list */
	case c_TRUNC:		/* trunc */
	case c_ERASE:		/* erase */
	case c_USER:		/* user */
		if (args->cmd != c_ERROR) {
			UNEXPECTED("commands '%c' and '%c' are mutually "
				   "exclusive", args->cmd, opt);
			return -1;
		}

		args->cmd = (cmd_t) opt;
		if (args->cmd == c_USER)
			args->user = strdup(optarg);
		break;
	case o_POFFSET:		/* partition-offset */
		free(args->poffset);
		args->poffset = strdup(optarg);
		break;
	case o_TARGET:		/* target */
		free(args->target);
		args->target = strdup(optarg);
		break;
	case o_NAME:		/* name */
		free(args->name);
		args->name = strdup(optarg);
		break;
	case o_OFFSET:		/* offset */
		args->offset = strdup(optarg);
		break;
	case o_SIZE:		/* size */
		args->size = strdup(optarg);
		break;
	case o_BLOCK:		/* block */
		args->block = strdup(optarg);
		break;
	case o_VALUE:		/* value */
		args->value = strdup(optarg);
		break;
	case o_FLAGS:		/* flags */
		args->flags = strdup(optarg);
		break;
	case o_PAD:		/* pad */
		args->pad = strdup(optarg);
		break;
	case f_FORCE:		/* force */
		args->force = (flag_t) opt;
		break;
	case f_LOGICAL:		/* logical */
		args->logical = (flag_t) opt;
		break;
	case f_VERBOSE:		/* verbose */
		args->verbose = (flag_t) opt;
		break;
	case f_DEBUG:		/* debug */
		args->debug = (flag_t) opt;
		break;
	case f_HELP:		/* help */
		usage(true);
		exit(EXIT_SUCCESS);
		break;
	case '?':		/* error */
	default:
		usage(false);
		UNEXPECTED("unknown option '%c', please see --help for "
			   "details\n", opt);
		return -1;
	}

	return 0;
}

static int process_option(args_t * args, const char *opt)
{
	assert(args != NULL);
	assert(opt != NULL);

	if (args->opt_sz <= args->opt_nr) {
		size_t size = sizeof(*args->opt);

		args->opt_sz += 5;
		args->opt = (const char **)realloc(args->opt,
						   size * args->opt_sz);
		if (args->opt == NULL) {
			ERRNO(errno);
			return -1;
		}

		memset(args->opt + args->opt_nr, 0,
		       size * (args->opt_sz - args->opt_nr));
	}

	args->opt[args->opt_nr] = strdup(opt);
	args->opt_nr++;

	return 0;
}

static int validate_args(args_t * args)
{
	assert(args != NULL);

	if (args->cmd == c_ERROR) {
		usage(false);
		UNEXPECTED("no command specified, please "
			   "see --help for details\n");
		return -1;
	}

	if (args->poffset == NULL) {
		UNEXPECTED("--partition-offset is required for all '%s' "
			   "commands", args->short_name);
		return -1;
	}

	if (args->target == NULL) {
		UNEXPECTED("--target is required for all '%s' commands",
			   args->short_name);
		return -1;
	}

	#define REQUIRED(name,cmd)	({				\
	if (args->name == NULL) {					\
		UNEXPECTED("--%s is required for the --%s command",	\
			   #name, #cmd);				\
		return -1;						\
	}								\
					})

	#define UNSUPPORTED(name,cmd)	({				\
	if (args->name != NULL) {					\
		UNEXPECTED("--%s is unsupported for the --%s command",	\
			   #name, #cmd);				\
		return -1;						\
	}								\
					})

	if (args->cmd == c_CREATE) {
		REQUIRED(size, create);
		REQUIRED(block, create);
		REQUIRED(target, create);

		UNSUPPORTED(offset, create);
		UNSUPPORTED(flags, create);
		UNSUPPORTED(value, create);
		UNSUPPORTED(pad, create);
	} else if (args->cmd == c_ADD) {
		REQUIRED(name, add);
		REQUIRED(flags, add);

		if (args->logical != f_LOGICAL) {
			REQUIRED(size, add);
			REQUIRED(offset, add);
		}

		UNSUPPORTED(block, add);
		UNSUPPORTED(value, add);
	} else if (args->cmd == c_DELETE) {
		REQUIRED(name, delete);

		UNSUPPORTED(size, delete);
		UNSUPPORTED(offset, delete);
		UNSUPPORTED(block, delete);
		UNSUPPORTED(flags, delete);
		UNSUPPORTED(value, delete);
		UNSUPPORTED(pad, delete);
	} else if (args->cmd == c_LIST) {
		UNSUPPORTED(size, list);
		UNSUPPORTED(offset, list);
		UNSUPPORTED(block, list);
		UNSUPPORTED(flags, list);
		UNSUPPORTED(value, list);
		UNSUPPORTED(pad, list);
	} else if (args->cmd == c_TRUNC) {
		REQUIRED(name, trunc);

		UNSUPPORTED(offset, trunc);
		UNSUPPORTED(block, trunc);
		UNSUPPORTED(flags, trunc);
		UNSUPPORTED(value, trunc);
		UNSUPPORTED(pad, trunc);
	} else if (args->cmd == c_ERASE) {
		REQUIRED(name, erase);

		UNSUPPORTED(size, erase);
		UNSUPPORTED(offset, erase);
		UNSUPPORTED(flags, erase);
		UNSUPPORTED(value, erase);
	} else if (args->cmd == c_USER) {
		REQUIRED(name, user);

		UNSUPPORTED(size, user);
		UNSUPPORTED(offset, user);
		UNSUPPORTED(block, user);
		UNSUPPORTED(flags, user);
		UNSUPPORTED(pad, user);
	} else {
		UNEXPECTED("invalid command '%c'", args->cmd);
		return -1;
	}

	return 0;
}

static int process_args(args_t * args)
{
	assert(args != NULL);
	int rc = 0;

	switch (args->cmd) {
	case c_CREATE:
		rc = command_create(args);
		break;
	case c_ADD:
		rc = command_add(args);
		if (rc < 0)
			break;
		if (args->pad != NULL)
			rc = command_erase(args);
		break;
	case c_DELETE:
		rc = command_delete(args);
		break;
	case c_LIST:
		rc = command_list(args);
		break;
	case c_TRUNC:
		rc = command_trunc(args);
		break;
	case c_ERASE:
		rc = command_erase(args);
		break;
	case c_USER:
		rc = command_user(args);
		break;
	default:
		UNEXPECTED("NOT IMPLEMENTED YET => '%c'", args->cmd);
		rc = -1;
	}

	return rc;
}

static void free_args(args_t * args)
{
	free(args->short_name);
	free(args->name);
	free(args->target);
	free(args->offset);
	free(args->poffset);
	free(args->size);
	free(args->block);
	free(args->user);
	free(args->value);
	free(args->flags);
	free(args->pad);
}

static void args_dump(args_t * args)
{
	assert(args != NULL);

	if (args->short_name != NULL)
		printf("short_name[%s]\n", args->short_name);
	if (args->cmd != 0)
		printf("cmd[%c]\n", args->cmd);
	if (args->path != NULL)
		printf("path[%s]\n", args->path);
	if (args->target != NULL)
		printf("target[%s]\n", args->target);
	if (args->poffset != NULL)
		printf("poffset[%s]\n", args->poffset);
	if (args->name != NULL)
		printf("name[%s]\n", args->name);
	if (args->offset != NULL)
		printf("offset[%s]\n", args->offset);
	if (args->size != NULL)
		printf("size[%s]\n", args->size);
	if (args->block != NULL)
		printf("block[%s]\n", args->block);
	if (args->flags != NULL)
		printf("flags[%s]\n", args->flags);
	if (args->value != NULL)
		printf("value[%s]\n", args->value);
	if (args->pad != NULL)
		printf("pad[%s]\n", args->pad);
	for (int i = 0; i < args->opt_nr; i++) {
		if (args->opt[i] != NULL)
			printf("opt%d[%s]\n", i, args->opt[i]);
	}
	if (args->force != 0)
		printf("force[%c]\n", args->force);
	if (args->protected != 0)
		printf("protected[%c]\n", args->protected);
	if (args->logical != 0)
		printf("logical[%c]\n", args->logical);
	if (args->verbose != 0)
		printf("verbose[%c]\n", args->verbose);
	if (args->verbose != 0)
		printf("debug[%c]\n", args->debug);
}

int main(int argc, char *argv[])
{
	static const struct option long_opt[] = {
		/* commands */
		{"create", no_argument, NULL, c_CREATE},
		{"add", no_argument, NULL, c_ADD},
		{"delete", no_argument, NULL, c_DELETE},
		{"list", no_argument, NULL, c_LIST},
		{"trunc", no_argument, NULL, c_TRUNC},
		{"erase", no_argument, NULL, c_ERASE},
		{"user", required_argument, NULL, c_USER},
		/* options */
		{"partition-offset", required_argument, NULL, o_POFFSET},
		{"target", required_argument, NULL, o_TARGET},
		{"name", required_argument, NULL, o_NAME},
		{"offset", required_argument, NULL, o_OFFSET},
		{"size", required_argument, NULL, o_SIZE},
		{"block-size", required_argument, NULL, o_BLOCK},
		{"value", required_argument, NULL, o_VALUE},
		{"flags", required_argument, NULL, o_FLAGS},
		{"pad", required_argument, NULL, o_PAD},
		/* flags */
		{"force", no_argument, NULL, f_FORCE},
		{"protected", no_argument, NULL, f_PROTECTED},
		{"logical", no_argument, NULL, f_LOGICAL},
		{"verbose", no_argument, NULL, f_VERBOSE},
		{"debug", no_argument, NULL, f_DEBUG},
		{"help", no_argument, NULL, f_HELP},
		{0, 0, 0, 0}
	};

	static const char *short_opt;
	short_opt = "CADLTEU:p:t:n:o:s:b:u:g:a:frlvdh";

	int rc = EXIT_FAILURE;

	setlinebuf(stdout);

	if (argc == 1)
		usage(false), exit(rc);

	int opt = 0, idx = 0;
	while ((opt = getopt_long(argc, argv, short_opt, long_opt, &idx)) != -1)
		if (process_argument(&args, opt, optarg) < 0)
			goto error;

	/* getopt_long doesn't know what to do with orphans, */
	/* so we'll scoop them up here, and deal with them later */

	while (optind < argc)
		if (process_option(&args, argv[optind++]) < 0)
			goto error;

	if (args.debug == f_DEBUG)
		args_dump(&args);

	if (validate_args(&args) < 0)
		goto error;
	if (process_args(&args) < 0)
		goto error;

	rc = EXIT_SUCCESS;

	if (false) {
		err_t *err;
error:
		err = err_get();
		assert(err != NULL);

		fprintf(stderr, "%s: %s : %s(%d) : (code=%d) %.*s\n",
			program_invocation_short_name,
			err_type_name(err), err_file(err), err_line(err),
			err_code(err), err_size(err), (char *)err_data(err));
	}

	free_args(&args);
	return rc;
}

static void __ctor__(void) __constructor;
static void __ctor__(void)
{
	page_size = sysconf(_SC_PAGESIZE);

	args.short_name = strdup(program_invocation_short_name);
	args.poffset = strdup("0x3F0000,0x7F0000");
	args.target = strdup("/dev/mtdblock/sfc.of");
	args.name = strdup(".*");
}
