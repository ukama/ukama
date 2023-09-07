/*
 * Copyright (c) 1997,2008 Andrew G. Morgan  <morgan@kernel.org>
 *
 * This displays the capabilities of given target process(es).
 */

#include <sys/types.h>
#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <sys/capability.h>

static void usage(int exiter)
{
    fprintf(stderr,
"usage: getcaps <pid> [<pid> ...]\n\n"
"  This program displays the capabilities on the queried process(es).\n"
"  The capabilities are displayed in the cap_from_text(3) format.\n\n"
"  Optional arguments:\n"
"     --help or --usage     display this message.\n"
"     --verbose             use a more verbose output format.\n"
"     --ugly or --legacy    use the archaic legacy output format.\n\n"
"[Copyright (c) 1997-8,2007,2019 Andrew G. Morgan  <morgan@kernel.org>]\n"
	);
    exit(exiter);
}

int main(int argc, char **argv)
{
    int retval = 0;
    int verbose = 0;

    if (argc < 2) {
	usage(1);
    }

    for ( ++argv; --argc > 0; ++argv ) {
	ssize_t length;
	int pid;
	cap_t cap_d;

	if (!strcmp(argv[0], "--help") || !strcmp(argv[0], "--usage")) {
	    usage(0);
	} else if (!strcmp(argv[0], "--verbose")) {
	    verbose = 1;
	    continue;
	} else if (!strcmp(argv[0], "--ugly") || !strcmp(argv[0], "--legacy")) {
	    verbose = 2;
	    continue;
	}

	pid = atoi(argv[0]);

	cap_d = cap_get_pid(pid);
	if (cap_d == NULL) {
		fprintf(stderr, "Failed to get cap's for process %d:"
			" (%s)\n", pid, strerror(errno));
		retval = 1;
		continue;
	} else {
	    char *result = cap_to_text(cap_d, &length);
	    if (verbose == 1) {
		printf("Capabilities for '%s': %s\n", *argv, result);
	    } else if (verbose == 2) {
		fprintf(stderr, "Capabilities for `%s': %s\n", *argv, result);
	    } else {
		printf("%s: %s\n", *argv, result);
	    }
	    cap_free(result);
	    result = NULL;
	    cap_free(cap_d);
	}
    }

    return retval;
}
