/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: fcp/src/main.c $                                              */
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
 *    File: fcp_main.c
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: cp for FFS files
 *    Date: 04/25/2013
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

args_t args;
size_t page_size;

int verbose;
int debug;

static void usage(bool verbose)
{
	FILE *e = stderr;

	fprintf(e,
		"fcp - FFS copy v%d.%d.%d -- Authors: "
		"<shaun@us.ibm.com>\n", FCP_MAJOR, FCP_MINOR, FCP_PATCH);

	fprintf(e, "\n");
	fprintf(e, "Usage:\n");
	fprintf(e," fcp [<src_type>:]<src_target>[:<src_name>] -PLETU"
		"\n     [-b <size>] [-o <offset,...>] [-fpvdh]\n");
	fprintf(e," fcp [<src_type>:]<src_target>[:<src_name>] "
		  "[<dst_type>:]<dst_target>[:<dst_name>]  -RWCM"
		  "\n     [-b <size>] [-o <offset,...>] [-fpvdh]\n");
	fprintf(e, "\n");
	fprintf(e, "    <type>\n");
	fprintf(e, "       'aa' : Aardvark USB probe\n");
	fprintf(e, "       'rw' : RISCWatch Ethernet probe\n");
	fprintf(e, "      'sfc' : FSP SFC character device\n");
	fprintf(e, "     'file' : UNIX regular file\n");
	fprintf(e, "  <target>\n");
	fprintf(e, "       'aa' : <number> USB device number [0..9]\n");
	fprintf(e, "       'rw' : <hostname>@<port> RISCwatch probe\n");
	fprintf(e, "      'sfc' : <path> to SFC char device\n");
	fprintf(e, "     'file' : <path> to UNIX regular file\n");
	fprintf(e, "    <name>\n");
	fprintf(e, "            : FFS name, e.g. bank0/bootenv/card\n");
	fprintf(e, "\n");

	if (verbose) {
		fprintf(e, "Examples:\n");
		fprintf(e, " fcp -P aa:0\n");
		fprintf(e, " fcp -P rw:riscwatch.ibm.com\n");
		fprintf(e, "\n");
		fprintf(e, " fcp -L aa:0\n");
		fprintf(e, " fcp -L rw:riscwatch.ibm.com\n");
		fprintf(e, " fcp -L nor.mif\n");
		fprintf(e, "\n");
		fprintf(e, " fcp -E 0 nor.mif:bank0\n");
		fprintf(e, " fcp -E 0xff nor.mif:bank0/spl\n");
		fprintf(e, "\n");
		fprintf(e, " fcp -W ipl.bin file:nor:bank0/ipl\n");
		fprintf(e, " fcp -R file:nor:bank0/ipl ipl.bin\n");
		fprintf(e, "\n");
		fprintf(e, " cat ipl.bin | fcp -R file:nor:bank0/ipl -\n");
		fprintf(e, " fcp -W file:nor:bank0/ipl - | cat > ipl.bin\n");
		fprintf(e, "\n");
		fprintf(e, " fcp -C file:nor:bank0 rw:host.ibm.com:bank0\n");
		fprintf(e, " fcp -C file:nor rw:host.ibm.com@6470\n");
		fprintf(e, "\n");
		fprintf(e, " fcp -M file:nor:bank0 rw:host.ibm.com:bank0\n");
		fprintf(e, " fcp -M file:nor:* rw:host.ibm.com@6470:*\n");
		fprintf(e, "\n");
		fprintf(e, " fcp -C nor.mif:bank0 nor.mif:bank1\n");
		fprintf(e, " fcp -C rw:host:bank0/ipl aa:0:bank1/ipl\n");
		fprintf(e, "\n");
		fprintf(e, " fcp -M nor.mif:bank0 nor.mif:bank1\n");
		fprintf(e, " fcp -M rw:host:bank0/ipl aa:0:bank1/ipl\n");
		fprintf(e, "\n");
		fprintf(e, " fcp -T nor.mif:bank0/spl\n");
		fprintf(e, " fcp -T 0 nor.mif:bank0/spl\n");
		fprintf(e, "\n");
		fprintf(e, " fcp -U nor.mif:bank0/spl\n");
		fprintf(e, " fcp -U 0 1 2 nor.mif:bank0/spl\n");
		fprintf(e, " fcp -U 0=0xffffffff 1=0 nor.mif:bank0/spl\n");
		fprintf(e, "\n");
	}

	/* =============================== */

	fprintf(e, "Commands:\n");

	fprintf(e, "  -P, --probe\n");
	if (verbose)
		fprintf(e,
			"\n  Read the SPI NOR device Id\n\n");

	fprintf(e, "  -L, --list\n");
	if (verbose)
		fprintf(e,
			"\n  List the contents of a partition table\n\n");

	fprintf(e, "  -R, --read\n");
	if (verbose)
		fprintf(e,
			"\n  Read the contents of a partition and write "
			"the data to an output file (use\n  '-' for stdout)."
			"\n\n");

	fprintf(e, "  -W, --write\n");
	if (verbose)
		fprintf(e,
			"\n  Read the contents of an input file (use '-' for "
			"stdin) and write the data\n  to a partition.\n\n");

	fprintf(e, "  -E, --erase  <value>\n");
	if (verbose)
		fprintf(e,
			"\n  Fill the contents of a partition with <value>.  "
			"<value> is a decimal (or hex)\n  number, default is "
			"0xFF.\n\n");

	fprintf(e, "  -C, --copy\n");
	if (verbose)
		fprintf(e,
			"\n  Copy source partition(s) to destination "
			"partition(s).  Both source and\n  destination name(s) "			"can specify either 'data' or 'logical' partitions."
			"\n\n");

	fprintf(e, "  -T, --trunc  <size>\n");
	if (verbose)
		fprintf(e,
			"\n  Truncate the actual size of partition(s) to "
			"<size> size bytes.  <size> is\n  a decimal (or hex) "
			"number, default is the partition entry size\n\n");

	fprintf(e, "  -M, --compare\n");
	if (verbose)
		fprintf(e,
			"\n  Compare source partition(s) to destination "
			"partition(s).  Both source and\n  destination name(s) "
			"can specify either 'data' or 'logical' partitions."
			"\n\n");

	fprintf(e, "  -U, --user   [<word>[=<value>] ...]\n");
	if (verbose)
		fprintf(e,
			"\n  Get or set a user word.  <word> and <value> are "
			"decimal (or hex) numbers.\n");

	fprintf(e, "\n");

	fprintf(e, "Options:\n");
	fprintf(e, "  -o, --offset <offset[,offset]>\n");
	if (verbose)
		fprintf(e, "\n  Specifies a comma (,) separated list of"
			" partition table offsets, in bytes\n  from the start"
			" of the target file (or device).\n\n");

	fprintf(e, "  -b, --buffer <value>\n");
	if (verbose)
		fprintf(e,
			"\n  Ignored.\n\n");
	fprintf(e, "\n");

	/* =============================== */

	fprintf(e, "Flags:\n");

	fprintf(e, "  -f, --force\n");
	if (verbose)
		fprintf(e, "\n  Override command safe guards\n\n");

	fprintf(e, "  -p, --protected\n");
	if (verbose)
		fprintf(e, "\n  Do not ignore protected partition "
			"entries\n\n");

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
	case c_PROBE:		/* probe */
	case c_LIST:		/* list */
	case c_READ:		/* read */
	case c_WRITE:		/* write */
	case c_ERASE:		/* erase */
	case c_COPY:		/* copy */
	case c_TRUNC:		/* trunc */
	case c_COMPARE:		/* compare */
	case c_USER:		/* user */
		if (args->cmd != c_ERROR) {
			UNEXPECTED("commands '%c' and '%c' are mutually "
				   "exclusive", args->cmd, opt);
			return -1;
		}
		args->cmd = (cmd_t) opt;
		break;
	case o_OFFSET:		/* offset */
		args->offset = strdup(optarg);
		break;
	case o_BUFFER:		/* buffer */
		/* We ignore it, it's useless but kept for backwards compat */
		break;
	case f_FORCE:		/* force */
		args->force = (flag_t) opt;
		break;
	case f_PROTECTED:	/* protected */
		args->protected = (flag_t) opt;
		break;
	case f_VERBOSE:		/* verbose */
		verbose = 1;
		args->verbose = (flag_t) opt;
		break;
	case f_DEBUG:		/* debug */
		debug = 1;
		args->debug = (flag_t) opt;
		break;
	case f_HELP:		/* help */
		usage(true);
		exit(EXIT_SUCCESS);
		break;
	case '?':		/* error */
	default:
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

/*
 * probe:
 *	fcp [<type>:]<target>[:<path>] -P
 * list:
 * 	fcp <path> -L
 *	fcp [<type>:]<target>[:<path>] -L
 * erase:
 * 	fcp [<type>:]<target>[:<path>] -E <pad>
 * trunc:
 * 	fcp [<type>:]<target>[:<path>] -T <size>
 * user:
 * 	fcp [<type>:]<target>[:<path>] -U <word>[=<value>] ...
 * write:
 * 	fcp <path>			[<type>:]<target>:<path> -W
 * read:
 * 	fcp [<type>:]<target>:<path>	<path> 			 -R
 * copy:
 * 	fcp <type>:]<target>:<source>	[<type>:]<target>:<dest>      -C
 * 	fcp [<type>:]<source>:<path> 	[<type>:]<destination>:<path> -C
 * compare:
 * 	fcp <type>:]<target>:<source>	[<type>:]<target>:<dest>      -M
 * 	fcp [<type>:]<source>:<path> 	[<type>:]<destination>:<path> -M
 */
static int validate_args(args_t * args)
{
	assert(args != NULL);

	if (args->offset == NULL) {
		UNEXPECTED("--offset is required for all '%s' "
			   "commands", args->short_name);
		return -1;
	}

	switch (args->cmd) {
	case c_PROBE:
	case c_LIST:
	case c_ERASE:
	case c_TRUNC:
	case c_USER:
		if (args->opt_nr < 1) {
			UNEXPECTED("invalid options, please see --help for "
				   "details");
			return -1;
		}

		if (parse_path(args->opt[0], &args->dst_type,
			       &args->dst_target, &args->dst_name) < 0)
			return -1;

		break;
	case c_READ:
	case c_WRITE:
	case c_COPY:
	case c_COMPARE:
		if (args->opt_nr < 2) {
			UNEXPECTED("invalid options, please see --help for "
				   "details");
			return -1;
		}

		if (parse_path(args->opt[0], &args->src_type,
			       &args->src_target, &args->src_name) < 0)
			return -1;
		if (parse_path(args->opt[1], &args->dst_type,
			       &args->dst_target, &args->dst_name) < 0)
			return -1;

		break;
	case 0:
		UNEXPECTED("specify a command, or please see --help "
			   "for details");
		return -1;
	default:
		UNEXPECTED("invalid command '%c', please see --help for "
			   "details", args->cmd);
		return -1;

	}

	debug("cmd[%c]\n", args->cmd);

	debug("  src_type: '%s'\n", args->src_type);
	debug("src_target: '%s'\n", args->src_target);
	debug("  src_name: '%s'\n", args->src_name);

	debug("  dst_type: '%s'\n", args->dst_type);
	debug("dst_target: '%s'\n", args->dst_target);
	debug("  dst_name: '%s'\n", args->dst_name);

	#define REQ_OPT(name,cmd)	({				\
	if (args->name == 0) {						\
		UNEXPECTED("--%s is required for the --%s command",	\
			   #name, #cmd);				\
		syntax();						\
		return -1;						\
	}								\
					})

	#define REQ_FIELD(name,cmd)	({				\
	if (args->name == 0) {						\
		UNEXPECTED("<%s> is required for the --%s command",	\
			   #name, #cmd);				\
		syntax();						\
		return -1;						\
	}								\
					})

	#define UNSUP_OPT(name,cmd)	({				\
	if (args->name != 0) {						\
		UNEXPECTED("--%s is unsupported for the --%s command",	\
			   #name, #cmd);				\
		syntax();						\
		return -1;						\
	}								\
					})

	if (args->cmd == c_PROBE) {
		void syntax(void) {
			fprintf(stderr, "Syntax: %s [<dst_type>:]<dst_target>"
				"[:<dst_name>] --probe [--verbose]\n",
				args->short_name);
		}
		if (args->opt_nr != 1) {
			syntax();
			UNEXPECTED("syntax error");
			return -1;
		}

	} else if (args->cmd == c_LIST) {
		void syntax(void) {
			fprintf(stderr, "Syntax: %s [<dst_type>:]<dst_target>"
				"[:<dst_name>] --list [--verbose]\n",
				args->short_name);
		}
		if (args->opt_nr != 1) {
			syntax();
			UNEXPECTED("syntax error");
			return -1;
		}

	} else if (args->cmd == c_READ) {
		void syntax(void) {
			fprintf(stderr, "Syntax: %s [<src_type>:]<src_source>"
				":<src_name> <path> --read [--verbose] "
				"[--force] [--protected] [--buffer <value>]\n",
				args->short_name);
		}
		if (args->opt_nr != 2) {
			syntax();
			UNEXPECTED("syntax error");
			return -1;
		}

		REQ_FIELD(src_name, read);
	} else if (args->cmd == c_WRITE) {
		void syntax(void) {
			fprintf(stderr, "Syntax: %s <path> [<dst_type>:]"
				"<dst_target>:<dst_name> --write [--verbose] "
				"[--force] [--protected] [--buffer <value>]\n",
				args->short_name);
		}
		if (args->opt_nr != 2) {
			syntax();
			UNEXPECTED("syntax error");
			return -1;
		}


		REQ_FIELD(dst_name, write);
	} else if (args->cmd == c_ERASE) {
		void syntax(void) {
			fprintf(stderr, "Syntax: %s [<dst_type>:]<dst_target>"
				":<dst_name> --erase <value> [--verbose] "
				"[--force] [--protected] [--buffer <value>]\n",
				args->short_name);
		}
		if (args->opt_nr != 2) {
			syntax();
			UNEXPECTED("syntax error");
			return -1;
		}

		REQ_FIELD(dst_name, erase);
	} else if (args->cmd == c_TRUNC) {
		void syntax(void) {
			fprintf(stderr, "Syntax: %s [<dst_type>:]<dst_target>"
				":<dst_name> --trunc <value> [--verbose] "
				"[--force] [--protected]\n", args->short_name);
		}
		if (2 < args->opt_nr) {
			syntax();
			UNEXPECTED("syntax error");
			return -1;
		}

		REQ_FIELD(dst_name, trunc);

	} else if (args->cmd == c_USER) {
		void syntax(void) {
			fprintf(stderr, "Syntax: %s [<dst_type>:]<dst_target>"
				":<dst_name> --user <value>[=<value] "
				"[--verbose] [--force] [--protected]\n",
				args->short_name);
		}

		REQ_FIELD(dst_name, user);

	} else if (args->cmd == c_COPY) {
		void syntax(void) {
			fprintf(stderr, "Syntax: %s [<src_type>:]<src_target>"
				"[:<src_name>] [<dst_type>:]<dst_target>"
				"[:<dst_name>] --copy [--verbose] [--force] "
				"[--protected] [--buffer <value>]\n",
				args->short_name);
		}
		if (args->opt_nr != 2) {
			syntax();
			UNEXPECTED("syntax error");
			return -1;
		}
	} else if (args->cmd == c_COMPARE) {
		void syntax(void) {
			fprintf(stderr, "Syntax: %s [<src_type>:]<src_target>"
				"[:<src_name>] [<dst_type>:]<dst_target>"
				"[:<dst_name>] --compare [--verbose] [--force] "
				"[--protected] [--buffer <value>]\n",
				args->short_name);
		}
		if (args->opt_nr != 2) {
			syntax();
			UNEXPECTED("syntax error");
			return -1;
		}
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
	case c_PROBE:
		//rc = command_probe(args);
		break;
	case c_LIST:
		rc = command_list(args);
		break;
	case c_READ:
		rc = command_read(args);
		break;
	case c_WRITE:
		rc = command_write(args);
		break;
	case c_ERASE:
		rc = command_erase(args);
		break;
	case c_TRUNC:
		rc = command_trunc(args);
		break;
	case c_USER:
		rc = command_user(args);
		break;
	case c_COPY:
	case c_COMPARE:
		rc = command_copy_compare(args);
		break;
	default:
		UNEXPECTED("NOT IMPLEMENTED YET => '%c'", args->cmd);
		rc = -1;
	}

	return rc;
}

static void args_dump(args_t * args)
{
	assert(args != NULL);

	if (args->short_name != NULL)
		printf("short_name[%s]\n", args->short_name);
	printf("cmd[%c]\n", args->cmd);
	if (args->offset != NULL)
		printf("offset[%s]\n", args->offset);
	if (args->force != 0)
		printf("force[%c]\n", args->force);
	if (args->protected != 0)
		printf("protected[%c]\n", args->protected);
	if (args->verbose != 0)
		printf("verbose[%c]\n", args->verbose);
	if (args->debug != 0)
		printf("debug[%c]\n", args->debug);

	for (int i = 0; i < args->opt_nr; i++) {
		if (args->opt[i] != NULL)
			printf("opt%d[%s]\n", i, args->opt[i]);
	}
}

int main(int argc, char *argv[])
{
	static const struct option long_opt[] = {
		/* commands */
		{"probe", no_argument, NULL, c_PROBE},
		{"list", no_argument, NULL, c_LIST},
		{"read", no_argument, NULL, c_READ},
		{"write", no_argument, NULL, c_WRITE},
		{"erase", no_argument, NULL, c_ERASE},
		{"copy", no_argument, NULL, c_COPY},
		{"trunc", no_argument, NULL, c_TRUNC},
		{"compare", no_argument, NULL, c_COMPARE},
		{"user", no_argument, NULL, c_USER},
		/* options */
		{"offset", required_argument, NULL, o_OFFSET},
		{"buffer", required_argument, NULL, o_BUFFER},
		/* flags */
		{"force", no_argument, NULL, f_FORCE},
		{"protected", no_argument, NULL, f_PROTECTED},
		{"verbose", no_argument, NULL, f_VERBOSE},
		{"debug", no_argument, NULL, f_DEBUG},
		{"help", no_argument, NULL, f_HELP},
		{0, 0, 0, 0}
	};

	static const char *short_opt;
	short_opt = "PLRWECTMUo:b:fpvdh";

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
error:
		dump_errors(args.short_name, stderr);
	}

	return rc;
}

static void __ctor__(void) __constructor;
static void __ctor__(void)
{
	page_size = sysconf(_SC_PAGESIZE);

	args.short_name = program_invocation_short_name;
	args.offset = "0x3F0000,0x7F0000";
}
