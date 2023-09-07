/*
 * Copyright (c) 1997,2007 Andrew G. Morgan  <morgan@kernel.org>
 *
 * This displays the capabilities of a given file.
 */

#undef _XOPEN_SOURCE
#define _XOPEN_SOURCE 500

#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <dirent.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/capability.h>

#include <ftw.h>

static int verbose = 0;
static int recursive = 0;
static int namespace = 0;

static void usage(void)
{
    fprintf(stderr,
	    "usage: getcap [-v] [-r] [-h] [-n] <filename> [<filename> ...]\n"
	    "\n"
	    "\tdisplays the capabilities on the queried file(s).\n"
	);
    exit(1);
}

static int do_getcap(const char *fname, const struct stat *stbuf,
		     int tflag, struct FTW* ftwbuf)
{
    cap_t cap_d;
    char *result;
    uid_t rootid;

    if (tflag != FTW_F) {
	if (verbose) {
	    printf("%s (Not a regular file)\n", fname);
	}
	return 0;
    }

    cap_d = cap_get_file(fname);
    if (cap_d == NULL) {
	if (errno != ENODATA) {
	    fprintf(stderr, "Failed to get capabilities of file '%s' (%s)\n",
		    fname, strerror(errno));
	} else if (verbose) {
	    printf("%s\n", fname);
	}
	return 0;
    }

    result = cap_to_text(cap_d, NULL);
    if (!result) {
	fprintf(stderr,
		"Failed to get capabilities of human readable format at '%s' (%s)\n",
		fname, strerror(errno));
	cap_free(cap_d);
	return 0;
    }
    rootid = cap_get_nsowner(cap_d);
    if (namespace && (rootid+1 > 1)) {
	printf("%s %s [rootid=%d]\n", fname, result, rootid);
    } else {
	printf("%s %s\n", fname, result);
    }
    cap_free(cap_d);
    cap_free(result);

    return 0;
}

int main(int argc, char **argv)
{
    int i, c;

    while ((c = getopt(argc, argv, "rvhn")) > 0) {
	switch(c) {
	case 'r':
	    recursive = 1;
	    break;
	case 'v':
	    verbose = 1;
	    break;
	case 'n':
	    namespace = 1;
	    break;
	default:
	    usage();
	}
    }

    if (!argv[optind])
	usage();

    for (i=optind; argv[i] != NULL; i++) {
	struct stat stbuf;

	if (lstat(argv[i], &stbuf) != 0) {
	    fprintf(stderr, "%s (%s)\n", argv[i], strerror(errno));
	} else if (recursive) {
	    nftw(argv[i], do_getcap, 20, FTW_PHYS);
	} else {
	    int tflag = S_ISREG(stbuf.st_mode) ? FTW_F :
		(S_ISLNK(stbuf.st_mode) ? FTW_SL : FTW_NS);
	    do_getcap(argv[i], &stbuf, tflag, 0);
	}
    }

    return 0;
}
