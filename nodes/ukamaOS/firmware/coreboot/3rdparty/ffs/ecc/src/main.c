/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: ecc/src/main.c $                                              */
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
 *   Descr: cmdline tool for ECC
 *    Date: 08/31/2012
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

#include <clib/attribute.h>
#include <clib/misc.h>
#include <clib/min.h>
#include <clib/ecc.h>
#include <clib/err.h>

#include "main.h"

#define ECC_SIZE	8

args_t args;

static void usage(const char *short_name, bool verbose)
{
	const char *n = short_name;
	FILE *e = stderr;

	fprintf(e, "NOR ECC Tool v%d.%d.%d -- Authors: <shaun@us.ibm.com>\n",
		ECC_MAJOR, ECC_MINOR, ECC_PATCH);

	fprintf(e, "\nUsage:\n");
	fprintf(e, "    %s <command> <path> <options>...\n", n);

	fprintf(e, "\nExamples:\n");
	fprintf(e, "    %s --inject sample.nor --output sample.nor.ecc\n", n);
	fprintf(e, "    %s --remove sample.nor.ecc --output sample.nor\n", n);
	fprintf(e, "    %s --hexdump sample.nor.ecc\n", n);

	fprintf(e, "\nCommands:\n");
	fprintf(e, "  -I, --inject <path> [options]\n");
	if (verbose)
		fprintf(e,
			"\n    Make a copy of the input file <path> and inject an ECC syndrom byte for\n"
			"    every 8 bytes of input data.\n\n");

	fprintf(e, "  -R, --remove <path> [options]\n");
	if (verbose)
		fprintf(e,
			"\n    Make a copy of the input file <path> and remove the ECC syndrom byte for\n"
			"    for every 9 byte of input data.\n\n");

	fprintf(e, "  -H, --hexdump <path> [options]\n");
	if (verbose)
		fprintf(e,
			"\n    Hex dump the contents of file <path> to stdout.\n\n");

	fprintf(e, "\nOptions:\n");
	fprintf(e, "  -o, --output <path>\n");
	if (verbose)
		fprintf(e, "\n    Specifies the output file path name.\n\n");

	fprintf(e, "  -h, --help\n");
	if (verbose)
		fprintf(e, "\n    Write this help text to stderr and exit\n");

	fprintf(e, "\nFlags:\n");
	fprintf(e, "  -p, --p8\n");
	if (verbose)
		fprintf(e,
			"\n    Invert the ECC bits for the P8 ECC engine.\n\n");

	fprintf(e, "\n");

	fprintf(e,
		"Report bugs to <https://bugzilla.linux.ibm.com/> (vendor='MCP for FSP*')\n");
}

static int process_argument(args_t * args, int opt, const char *optarg)
{
	assert(args != NULL);

	switch (opt) {
	case c_INJECT:		/* inject */
	case c_REMOVE:		/* remove */
	case c_HEXDUMP:		/* hexdump */
		if (args->cmd != c_ERROR) {
			UNEXPECTED("commands '%c' and '%c' are mutually "
				   "exclusive", args->cmd, opt);
			return -1;
		}

		args->cmd = (cmd_t) opt;
		args->path = strdup(optarg);

		break;
	case o_OUTPUT:		/* offset */
		args->file = strdup(optarg);
		break;
	case f_FORCE:		/* force */
		args->force = (flag_t) opt;
		break;
	case f_VERBOSE:		/* verbose */
		args->verbose = (flag_t) opt;
		break;
	case f_P8:		/* p8 */
		args->p8 = (flag_t) opt;
		break;
	case f_HELP:		/* help */
		usage(args->short_name, true);
		exit(EXIT_SUCCESS);
		break;
	case '?':		/* error */
	default:
		usage(args->short_name, false);
		UNEXPECTED("unknown option '%c', please see "
			   "--help for details\n", opt);
		return -1;
	}

	return 0;
}

static int process_option(args_t * args, const char *opt)
{
	assert(args != NULL);
	assert(opt != NULL);

	if (args->opt_sz <= args->opt_nr) {
		args->opt_sz += 5;
		args->opt = (const char **)realloc(args->opt,
						   sizeof(*args->opt) *
						   args->opt_sz);
		memset(args->opt + args->opt_nr, 0,
		       sizeof(*args->opt) * (args->opt_sz - args->opt_nr));
	}

	args->opt[args->opt_nr] = strdup(opt);
	args->opt_nr++;

	return 0;
}

bool check_extension(const char *path, const char *ext)
{
	int len = strlen(path), ext_len = strlen(ext);
	return (ext_len < len)
	    && (strncasecmp(path + len - ext_len, ext, ext_len) == 0);
}

static int validate_args(args_t * args)
{
	assert(args != NULL);

	if (args->cmd == c_ERROR) {
		usage(args->short_name, false);
		UNEXPECTED("no command specified, please see --help "
			   "for details\n");
		return -1;
	}

	if (args->path == NULL) {
		UNEXPECTED("no path specified, please see --help "
			   "for details");
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


	if (args->cmd == c_INJECT) {
		if (check_extension(args->path, ECC_EXT)) {
			UNEXPECTED("'%s' extension '%s' already exists "
				   "-- ignored", args->path, ECC_EXT);
			return -1;
		}

		if (args->file == NULL) {
			args->file = (char *)malloc(strlen(args->path) + 8);
			assert(args->file != NULL);
			sprintf((char *)args->file, "%s%s", args->path,
				ECC_EXT);

			if (args->force != f_FORCE)
				fprintf(stderr, "%s: --output <file> missing, "
				        "writing output to '%s'\n",
					args->short_name, args->file);
		}
	} else if (args->cmd == c_REMOVE) {
		if (!check_extension(args->path, ECC_EXT)) {
			UNEXPECTED("'%s' unknown extension, must be '%s' -- "
				   "ignored", args->path, ECC_EXT);
			return -1;
		}

		if (args->file == NULL) {
			args->file = strdup(args->path);
			assert(args->file != NULL);
			sprintf((char *)args->file, "%*s",
				(uint32_t)(strlen(args->path) - strlen(ECC_EXT)),
				args->path);

			fprintf(stderr, "%s: --output <file> missing, writing "
				"output to '%s'\n",
				args->short_name, args->file);
		}
	} else if (args->cmd == c_HEXDUMP) {
		if (!check_extension(args->path, ECC_EXT)) {
			UNEXPECTED("'%s' unknown extension, must be '%s' -- "
				   "ignored", args->path, ECC_EXT);
			return -1;
		}
	} else {
		UNEXPECTED("'%c' invalid command", args->cmd);
		return -1;
	}

	return 0;
}

static int command_inject(args_t * args)
{
	assert(args != NULL);

	struct stat st;
	if (stat(args->path, &st) != 0) {
		ERRNO(errno);
		return -1;
	}

	if (!S_ISREG(st.st_mode)) {
		ERRNO(errno);
		return -1;
	}

	FILE *i = fopen(args->path, "r");
	if (i == NULL) {
		ERRNO(errno);
		return-1;
	}

	FILE *o = fopen(args->file, "w");
	if (o == NULL) {
		ERRNO(errno);
		return -1;
	}

#define INPUT_SIZE		(4096 - (4096 / ECC_SIZE))
	char input[INPUT_SIZE];	// 4KB less 512 ECC bytes
#undef INPUT_SIZE

	size_t count = 0;
	while (count < st.st_size) {
		clearerr(i);

		size_t rc = fread(input, 1, sizeof input, i);
		if (rc == 0) {
			int err = ferror(i);
			if (err) {
				ERRNO(errno);
				return -1;
			} else
				break;
		}

		count += rc;

#define OUTPUT_SIZE		4096
		char output[OUTPUT_SIZE];
#undef OUTPUT_SIZE

		memset(output + sizeof input, 0, sizeof output - sizeof input);

		rc = (rc + 7) & ~7;	// 8-byte alignment

		ssize_t injected_size = 0;
		if (args->p8 == f_P8)
			injected_size =
			    p8_ecc_inject(output, sizeof output, input, rc);
		else
			injected_size =
			    sfc_ecc_inject(output, sizeof output, input, rc);
		if (injected_size < 0) {
			ERRNO(errno);
			return -1;
		}

		clearerr(o);
		rc = fwrite(output, 1, injected_size, o);
		if (rc == 0) {
			int err = ferror(o);
			if (err) {
				ERRNO(errno);
				return -1;
			}
		}
	}

	if (fclose(i) == EOF) {
		ERRNO(errno);
		return -1;
	}

	if (fclose(o) == EOF) {
		ERRNO(errno);
		return -1;
	}

	return 0;
}

static int command_remove(args_t * args)
{
	assert(args != NULL);

	struct stat st;
	if (stat(args->path, &st) != 0) {
		ERRNO(errno);
		return -1;
	}

	if (!S_ISREG(st.st_mode)) {
		ERRNO(errno);
		return -1;
	}

	FILE *i = fopen(args->path, "r");
	if (i == NULL) {
		ERRNO(errno);
		return -1;
	}

	FILE *o = fopen(args->file, "w");
	if (o == NULL) {
		ERRNO(errno);
		return -1;
	}

#define INPUT_SIZE		4086
	char input[INPUT_SIZE];	// multiple of 9-bytes
#undef INPUT_SIZE

	size_t count = 0;
	while (count < st.st_size) {
		clearerr(i);

		size_t rc = fread(input, 1, sizeof input, i);
		if (rc == 0) {
			int err = ferror(i);
			if (err) {
				ERRNO(errno);
				return -1;
			} else
				break;
		}

		count += rc;

    char output[(((sizeof input)/(ECC_SIZE+1))*ECC_SIZE) ];

		ssize_t removed_size;
		if (args->p8 == f_P8)
			removed_size =
			    p8_ecc_remove_size(output, sizeof output, input, rc);
		else
			removed_size =
			    sfc_ecc_remove(output, sizeof output, input, rc);
		if (removed_size < 0) {
			ERRNO(errno);
			return -1;
		}

		clearerr(o);
		rc = fwrite(output, 1, removed_size, o);
		if (rc == 0) {
			int err = ferror(o);
			if (err) {
				ERRNO(errno);
				return -1;
			}
		}
	}

	if (fclose(i) == EOF) {
		ERRNO(errno);
		return -1;
	}

	if (fclose(o) == EOF) {
		ERRNO(errno);
		return -1;
	}

	return 0;
}

static int command_hexdump(args_t * args)
{
	assert(args != NULL);

	struct stat st;
	if (stat(args->path, &st) != 0) {
		ERRNO(errno);
		return -1;
	}

	if (!S_ISREG(st.st_mode)) {
		ERRNO(errno);
		return -1;
	}

	FILE *i = fopen(args->path, "r");
	if (i == NULL) {
		ERRNO(errno);
		return -1;
	}

	FILE *o = stdout;
	if (args->file != NULL) {
		o = fopen(args->file, "w");
		if (o == NULL) {
			ERRNO(errno);
			return -1;
		}
	}
#define INPUT_SIZE 		4086
	char input[INPUT_SIZE];	// multiple of 9-bytes
#undef INPUT_SIZE

	if (setvbuf(i, NULL, _IOFBF, __round_pow2(sizeof input)) != 0) {
		ERRNO(errno);
		return -1;
	}

	size_t count = 0;
	while (count < st.st_size) {
		clearerr(i);

		size_t rc = fread(input, 1, sizeof input, i);
		if (rc == 0) {
			int err = ferror(i);
			if (err) {
				ERRNO(errno);
				return -1;
			} else
				break;
		}

		if (args->p8 == f_P8)
			p8_ecc_dump(o, count, input, rc);
		else
			sfc_ecc_dump(o, count, input, rc);

		count += rc;
	}

	if (fclose(i) == EOF) {
		ERRNO(errno);
		return -1;
	}

	if (o != stdout) {
		if (fclose(o) == EOF) {
			ERRNO(errno);
			return -1;
		}
	}

	return 0;
}

static int process_args(args_t * args)
{
	assert(args != NULL);

	switch (args->cmd) {
	case c_INJECT:
		command_inject(args);
		break;
	case c_REMOVE:
		command_remove(args);
		break;
	case c_HEXDUMP:
		command_hexdump(args);
		break;
	default:
		UNEXPECTED("NOT IMPLEMENTED YET => '%c'", args->cmd);
		return -1;
	}

	return 0;
}

static void args_dump(args_t * args)
{
	assert(args != NULL);

	printf("short_name[%s]\n", args->short_name);
	printf("path[%s]\n", args->path);
	printf("cmd[%d]\n", args->cmd);
	printf("output[%s]\n", args->file);
	printf("force[%d]\n", args->force);
	printf("p8[%d]\n", args->p8);
	printf("verbose[%d]\n", args->force);
}

int main(int argc, char *argv[])
{
	static const struct option long_opts[] = {
		/* commands */
		{"inject", required_argument, NULL, c_INJECT},
		{"remove", required_argument, NULL, c_REMOVE},
		{"hexdump", required_argument, NULL, c_HEXDUMP},
		/* options */
		{"output", required_argument, NULL, o_OUTPUT},
		/* flags */
		{"force", no_argument, NULL, f_FORCE},
		{"p8", no_argument, NULL, f_P8},
		{"verbose", no_argument, NULL, f_VERBOSE},
		{"help", no_argument, NULL, f_HELP},
		{0, 0, 0, 0}
	};

	static const char *short_opts = "I:R:H:o:fpvh";

	int rc = EXIT_FAILURE;

	if (argc == 1)
		usage(args.short_name, false), exit(rc);

	int opt = 0, idx = 0;
	while ((opt = getopt_long(argc, argv, short_opts, long_opts,
				  &idx)) != -1)
		if (process_argument(&args, opt, optarg) < 0)
			goto error;

	/* getopt_long doesn't know what to do with orphans, */
	/* so we'll scoop them up here, and deal with them later */

	while (optind < argc)
		if (process_option(&args, argv[optind++]) < 0)
			goto error;

	if (args.verbose == f_VERBOSE)
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

	return rc;
}

static void __ctor__(void) __constructor;
static void __ctor__(void)
{
	/* early initialization before main() is called the crt0 */
	args.short_name = program_invocation_short_name;
}
