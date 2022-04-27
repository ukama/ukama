/*
 * Copyright (c) 2008-11,16,19,2020 Andrew G. Morgan <morgan@kernel.org>
 *
 * This is a simple 'bash' (-DSHELL) wrapper program that can be used
 * to raise and lower both the bset and pI capabilities before
 * invoking /bin/bash.
 *
 * The --print option can be used as a quick test whether various
 * capability manipulations work as expected (or not).
 */

#ifndef _DEFAULT_SOURCE
#define _DEFAULT_SOURCE
#endif

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <sys/types.h>
#include <pwd.h>
#include <grp.h>
#include <errno.h>
#include <ctype.h>
#include <sys/capability.h>
#include <sys/prctl.h>
#include <sys/securebits.h>
#include <sys/wait.h>
#include <unistd.h>

#ifndef SHELL
#define SHELL "/bin/bash"
#endif /* ndef SHELL */

#define MAX_GROUPS       100   /* max number of supplementary groups for user */

static char *binary(unsigned long value)
{
    static char string[8*sizeof(unsigned long) + 1];
    unsigned i;

    i = sizeof(string);
    string[--i] = '\0';
    do {
	string[--i] = (value & 1) ? '1' : '0';
	value >>= 1;
    } while ((i > 0) && value);
    return string + i;
}

static void display_prctl_set(const char *name, int (*fn)(cap_value_t))
{
    unsigned cap;
    const char *sep;
    int set;

    printf("%s set =", name);
    for (sep = "", cap=0; (set = fn(cap)) >= 0; cap++) {
	char *ptr;
	if (!set) {
	    continue;
	}

	ptr = cap_to_name(cap);
	if (ptr == NULL) {
	    printf("%s%u", sep, cap);
	} else {
	    printf("%s%s", sep, ptr);
	    cap_free(ptr);
	}
	sep = ",";
    }
    if (!cap) {
	printf(" <unsupported>\n");
    } else {
	printf("\n");
    }
}

/* arg_print displays the current capability state of the process */
static void arg_print(void)
{
    long set;
    int status, j;
    cap_t all;
    char *text;
    const char *sep;
    struct group *g;
    gid_t groups[MAX_GROUPS], gid;
    uid_t uid, euid;
    struct passwd *u, *eu;
    cap_iab_t iab;

    all = cap_get_proc();
    text = cap_to_text(all, NULL);
    printf("Current: %s\n", text);
    cap_free(text);
    cap_free(all);

    display_prctl_set("Bounding", cap_get_bound);
    display_prctl_set("Ambient", cap_get_ambient);
    iab = cap_iab_get_proc();
    text = cap_iab_to_text(iab);
    printf("Current IAB: %s\n", text);
    cap_free(text);
    cap_free(iab);

    set = cap_get_secbits();
    if (set >= 0) {
	const char *b = binary(set);  /* verilog convention for binary string */
	printf("Securebits: 0%lo/0x%lx/%u'b%s\n", set, set,
	       (unsigned) strlen(b), b);
	printf(" secure-noroot: %s (%s)\n",
	       (set & SECBIT_NOROOT) ? "yes":"no",
	       (set & SECBIT_NOROOT_LOCKED) ? "locked":"unlocked");
	printf(" secure-no-suid-fixup: %s (%s)\n",
	       (set & SECBIT_NO_SETUID_FIXUP) ? "yes":"no",
	       (set & SECBIT_NO_SETUID_FIXUP_LOCKED) ? "locked":"unlocked");
	printf(" secure-keep-caps: %s (%s)\n",
	       (set & SECBIT_KEEP_CAPS) ? "yes":"no",
	       (set & SECBIT_KEEP_CAPS_LOCKED) ? "locked":"unlocked");
	if (CAP_AMBIENT_SUPPORTED()) {
	    printf(" secure-no-ambient-raise: %s (%s)\n",
		   (set & SECBIT_NO_CAP_AMBIENT_RAISE) ? "yes":"no",
		   (set & SECBIT_NO_CAP_AMBIENT_RAISE_LOCKED) ?
		   "locked":"unlocked");
	}
    } else {
	printf("[Securebits ABI not supported]\n");
	set = prctl(PR_GET_KEEPCAPS);
	if (set >= 0) {
	    printf(" prctl-keep-caps: %s (locking not supported)\n",
		   set ? "yes":"no");
	} else {
	    printf("[Keepcaps ABI not supported]\n");
	}
    }
    uid = getuid();
    u = getpwuid(uid);
    euid = geteuid();
    eu = getpwuid(euid);
    printf("uid=%u(%s) euid=%u(%s)\n", uid, u ? u->pw_name : "???", euid, eu ? eu->pw_name : "???");
    gid = getgid();
    g = getgrgid(gid);
    printf("gid=%u(%s)\n", gid, g ? g->gr_name : "???");
    printf("groups=");
    status = getgroups(MAX_GROUPS, groups);
    sep = "";
    for (j=0; j < status; j++) {
	g = getgrgid(groups[j]);
	printf("%s%u(%s)", sep, groups[j], g ? g->gr_name : "???");
	sep = ",";
    }
    printf("\n");
    cap_mode_t mode = cap_get_mode();
    printf("Guessed mode: %s (%d)\n", cap_mode_name(mode), mode);
}

static const cap_value_t raise_setpcap[1] = { CAP_SETPCAP };
static const cap_value_t raise_chroot[1] = { CAP_SYS_CHROOT };

static void push_pcap(cap_t *orig_p, cap_t *raised_for_setpcap_p)
{
    /*
     * We need to do this here because --inh=XXX may have reset
     * orig and it isn't until we are within the --drop code that
     * we know what the prevailing (orig) pI value is.
     */
    *orig_p = cap_get_proc();
    if (NULL == *orig_p) {
	perror("Capabilities not available");
	exit(1);
    }

    *raised_for_setpcap_p = cap_dup(*orig_p);
    if (NULL == *raised_for_setpcap_p) {
	fprintf(stderr, "modification requires CAP_SETPCAP\n");
	exit(1);
    }
    if (cap_set_flag(*raised_for_setpcap_p, CAP_EFFECTIVE, 1,
		     raise_setpcap, CAP_SET) != 0) {
	perror("unable to select CAP_SETPCAP");
	exit(1);
    }
}

static void pop_pcap(cap_t orig, cap_t raised_for_setpcap)
{
    cap_free(raised_for_setpcap);
    cap_free(orig);
}

static void arg_drop(const char *arg_names)
{
    char *ptr;
    cap_t orig, raised_for_setpcap;
    char *names;

    push_pcap(&orig, &raised_for_setpcap);
    if (strcmp("all", arg_names) == 0) {
	unsigned j = 0;
	while (CAP_IS_SUPPORTED(j)) {
	    int status;
	    if (cap_set_proc(raised_for_setpcap) != 0) {
		perror("unable to raise CAP_SETPCAP for BSET changes");
		exit(1);
	    }
	    status = cap_drop_bound(j);
	    if (cap_set_proc(orig) != 0) {
		perror("unable to lower CAP_SETPCAP post BSET change");
		exit(1);
	    }
	    if (status != 0) {
		char *name_ptr;

		name_ptr = cap_to_name(j);
		fprintf(stderr, "Unable to drop bounding capability [%s]\n",
			name_ptr);
		cap_free(name_ptr);
		exit(1);
	    }
	    j++;
	}
	pop_pcap(orig, raised_for_setpcap);
	return;
    }

    names = strdup(arg_names);
    if (NULL == names) {
	fprintf(stderr, "failed to allocate names\n");
	exit(1);
    }
    for (ptr = names; (ptr = strtok(ptr, ",")); ptr = NULL) {
	/* find name for token */
	cap_value_t cap;
	int status;

	if (cap_from_name(ptr, &cap) != 0) {
	    fprintf(stderr, "capability [%s] is unknown to libcap\n", ptr);
	    exit(1);
	}
	if (cap_set_proc(raised_for_setpcap) != 0) {
	    perror("unable to raise CAP_SETPCAP for BSET changes");
	    exit(1);
	}
	status = cap_drop_bound(cap);
	if (cap_set_proc(orig) != 0) {
	    perror("unable to lower CAP_SETPCAP post BSET change");
	    exit(1);
	}
	if (status != 0) {
	    fprintf(stderr, "failed to drop [%s=%u]\n", ptr, cap);
	    exit(1);
	}
    }
    pop_pcap(orig, raised_for_setpcap);
    free(names);
}

static void arg_change_amb(const char *arg_names, cap_flag_value_t set)
{
    char *ptr;
    cap_t orig, raised_for_setpcap;
    char *names;

    push_pcap(&orig, &raised_for_setpcap);
    if (strcmp("all", arg_names) == 0) {
	unsigned j = 0;
	while (CAP_IS_SUPPORTED(j)) {
	    int status;
	    if (cap_set_proc(raised_for_setpcap) != 0) {
		perror("unable to raise CAP_SETPCAP for AMBIENT changes");
		exit(1);
	    }
	    status = cap_set_ambient(j, set);
	    if (cap_set_proc(orig) != 0) {
		perror("unable to lower CAP_SETPCAP post AMBIENT change");
		exit(1);
	    }
	    if (status != 0) {
		char *name_ptr;

		name_ptr = cap_to_name(j);
		fprintf(stderr, "Unable to %s ambient capability [%s]\n",
			set == CAP_CLEAR ? "clear":"raise", name_ptr);
		cap_free(name_ptr);
		exit(1);
	    }
	    j++;
	}
	pop_pcap(orig, raised_for_setpcap);
	return;
    }

    names = strdup(arg_names);
    if (NULL == names) {
	fprintf(stderr, "failed to allocate names\n");
	exit(1);
    }
    for (ptr = names; (ptr = strtok(ptr, ",")); ptr = NULL) {
	/* find name for token */
	cap_value_t cap;
	int status;

	if (cap_from_name(ptr, &cap) != 0) {
	    fprintf(stderr, "capability [%s] is unknown to libcap\n", ptr);
	    exit(1);
	}
	if (cap_set_proc(raised_for_setpcap) != 0) {
	    perror("unable to raise CAP_SETPCAP for AMBIENT changes");
	    exit(1);
	}
	status = cap_set_ambient(cap, set);
	if (cap_set_proc(orig) != 0) {
	    perror("unable to lower CAP_SETPCAP post AMBIENT change");
	    exit(1);
	}
	if (status != 0) {
	    fprintf(stderr, "failed to %s ambient [%s=%u]\n",
		    set == CAP_CLEAR ? "clear":"raise", ptr, cap);
	    exit(1);
	}
    }
    pop_pcap(orig, raised_for_setpcap);
    free(names);
}

int main(int argc, char *argv[], char *envp[])
{
    pid_t child;
    unsigned i;
    const char *shell = SHELL;

    child = 0;

    char *temp_name = cap_to_name(cap_max_bits() - 1);
    if (temp_name[0] != 'c') {
	printf("WARNING: libcap needs an update (cap=%d should have a name).\n",
	       cap_max_bits() - 1);
    }
    cap_free(temp_name);

    for (i=1; i<argc; ++i) {
	if (!strncmp("--drop=", argv[i], 7)) {
	    arg_drop(argv[i]+7);
	} else if (!strncmp("--dropped=", argv[i], 10)) {
	    cap_value_t cap;
	    if (cap_from_name(argv[i]+10, &cap) < 0) {
		fprintf(stderr, "cap[%s] not recognized by library\n",
			argv[i] + 10);
		exit(1);
	    }
	    if (cap_get_bound(cap) > 0) {
		fprintf(stderr, "cap[%s] raised in bounding vector\n",
			argv[i]+10);
		exit(1);
	    }
	} else if (!strcmp("--has-ambient", argv[i])) {
	    if (!CAP_AMBIENT_SUPPORTED()) {
		fprintf(stderr, "ambient set not supported\n");
		exit(1);
	    }
	} else if (!strncmp("--addamb=", argv[i], 9)) {
	    arg_change_amb(argv[i]+9, CAP_SET);
	} else if (!strncmp("--delamb=", argv[i], 9)) {
	    arg_change_amb(argv[i]+9, CAP_CLEAR);
	} else if (!strncmp("--noamb", argv[i], 7)) {
	    if (cap_reset_ambient() != 0) {
		fprintf(stderr, "failed to reset ambient set\n");
		exit(1);
	    }
	} else if (!strncmp("--inh=", argv[i], 6)) {
	    cap_t all, raised_for_setpcap;
	    char *text;
	    char *ptr;

	    all = cap_get_proc();
	    if (all == NULL) {
		perror("Capabilities not available");
		exit(1);
	    }
	    if (cap_clear_flag(all, CAP_INHERITABLE) != 0) {
		perror("libcap:cap_clear_flag() internal error");
		exit(1);
	    }

	    raised_for_setpcap = cap_dup(all);
	    if ((raised_for_setpcap != NULL)
		&& (cap_set_flag(raised_for_setpcap, CAP_EFFECTIVE, 1,
				 raise_setpcap, CAP_SET) != 0)) {
		cap_free(raised_for_setpcap);
		raised_for_setpcap = NULL;
	    }

	    text = cap_to_text(all, NULL);
	    cap_free(all);
	    if (text == NULL) {
		perror("Fatal error concerning process capabilities");
		exit(1);
	    }
	    ptr = malloc(10 + strlen(argv[i]+6) + strlen(text));
	    if (ptr == NULL) {
		perror("Out of memory for inh set");
		exit(1);
	    }
	    if (argv[i][6] && strcmp("none", argv[i]+6)) {
		sprintf(ptr, "%s %s+i", text, argv[i]+6);
	    } else {
		strcpy(ptr, text);
	    }

	    all = cap_from_text(ptr);
	    if (all == NULL) {
		perror("Fatal error internalizing capabilities");
		exit(1);
	    }
	    cap_free(text);
	    free(ptr);

	    if (raised_for_setpcap != NULL) {
		/*
		 * This is only for the case that pP does not contain
		 * the requested change to pI.. Failing here is not
		 * indicative of the cap_set_proc(all) failing (always).
		 */
		(void) cap_set_proc(raised_for_setpcap);
		cap_free(raised_for_setpcap);
		raised_for_setpcap = NULL;
	    }

	    if (cap_set_proc(all) != 0) {
		perror("Unable to set inheritable capabilities");
		exit(1);
	    }
	    /*
	     * Since status is based on orig, we don't want to restore
	     * the previous value of 'all' again here!
	     */

	    cap_free(all);
	} else if (!strncmp("--caps=", argv[i], 7)) {
	    cap_t all, raised_for_setpcap;

	    raised_for_setpcap = cap_get_proc();
	    if (raised_for_setpcap == NULL) {
		perror("Capabilities not available");
		exit(1);
	    }

	    if ((raised_for_setpcap != NULL)
		&& (cap_set_flag(raised_for_setpcap, CAP_EFFECTIVE, 1,
				 raise_setpcap, CAP_SET) != 0)) {
		cap_free(raised_for_setpcap);
		raised_for_setpcap = NULL;
	    }

	    all = cap_from_text(argv[i]+7);
	    if (all == NULL) {
		fprintf(stderr, "unable to interpret [%s]\n", argv[i]);
		exit(1);
	    }

	    if (raised_for_setpcap != NULL) {
		/*
		 * This is only for the case that pP does not contain
		 * the requested change to pI.. Failing here is not
		 * indicative of the cap_set_proc(all) failing (always).
		 */
		(void) cap_set_proc(raised_for_setpcap);
		cap_free(raised_for_setpcap);
		raised_for_setpcap = NULL;
	    }

	    if (cap_set_proc(all) != 0) {
		fprintf(stderr, "Unable to set capabilities [%s]\n", argv[i]);
		exit(1);
	    }
	    /*
	     * Since status is based on orig, we don't want to restore
	     * the previous value of 'all' again here!
	     */

	    cap_free(all);
	} else if (!strcmp("--modes", argv[i])) {
	    cap_mode_t c;
	    printf("Supported modes:");
	    for (c = 1; ; c++) {
		const char *m = cap_mode_name(c);
		if (strcmp("UNKNOWN", m) == 0) {
		    break;
		}
		printf(" %s", m);
	    }
	    printf("\n");
	} else if (!strncmp("--mode=", argv[i], 7)) {
	    const char *target = argv[i]+7;
	    cap_mode_t c;
	    int found = 0;
	    for (c = 1; ; c++) {
		const char *m = cap_mode_name(c);
		if (!strcmp("UNKNOWN", m)) {
		    found = 0;
		    break;
		}
		if (!strcmp(m, target)) {
		    found = 1;
		    break;
		}
	    }
	    if (!found) {
		printf("unsupported mode: %s\n", target);
		exit(1);
	    }
	    int ret = cap_set_mode(c);
	    if (ret != 0) {
		printf("failed to set mode [%s]: %s\n",
		       target, strerror(errno));
		exit(1);
	    }
	} else if (!strncmp("--inmode=", argv[i], 9)) {
	    const char *target = argv[i]+9;
	    cap_mode_t c = cap_get_mode();
	    const char *m = cap_mode_name(c);
	    if (strcmp(m, target)) {
		printf("mismatched mode got=%s want=%s\n", m, target);
		exit(1);
	    }
	} else if (!strncmp("--keep=", argv[i], 7)) {
	    unsigned value;
	    int set;

	    value = strtoul(argv[i]+7, NULL, 0);
	    set = prctl(PR_SET_KEEPCAPS, value);
	    if (set < 0) {
		fprintf(stderr, "prctl(PR_SET_KEEPCAPS, %u) failed: %s\n",
			value, strerror(errno));
		exit(1);
	    }
	} else if (!strncmp("--chroot=", argv[i], 9)) {
	    int status;
	    cap_t orig, raised_for_chroot;

	    orig = cap_get_proc();
	    if (orig == NULL) {
		perror("Capabilities not available");
		exit(1);
	    }

	    raised_for_chroot = cap_dup(orig);
	    if (raised_for_chroot == NULL) {
		perror("Unable to duplicate capabilities");
		exit(1);
	    }

	    if (cap_set_flag(raised_for_chroot, CAP_EFFECTIVE, 1, raise_chroot,
			     CAP_SET) != 0) {
		perror("unable to select CAP_SET_SYS_CHROOT");
		exit(1);
	    }

	    if (cap_set_proc(raised_for_chroot) != 0) {
		perror("unable to raise CAP_SYS_CHROOT");
		exit(1);
	    }
	    cap_free(raised_for_chroot);

	    status = chroot(argv[i]+9);
	    if (cap_set_proc(orig) != 0) {
		perror("unable to lower CAP_SYS_CHROOT");
		exit(1);
	    }
	    /*
	     * Given we are now in a new directory tree, its good practice
	     * to start off in a sane location
	     */
	    status = chdir("/");

	    cap_free(orig);

	    if (status != 0) {
		fprintf(stderr, "Unable to chroot/chdir to [%s]", argv[i]+9);
		exit(1);
	    }
	} else if (!strncmp("--secbits=", argv[i], 10)) {
	    unsigned value;
	    int status;
	    value = strtoul(argv[i]+10, NULL, 0);
	    status = cap_set_secbits(value);
	    if (status < 0) {
		fprintf(stderr, "failed to set securebits to 0%o/0x%x\n",
			value, value);
		exit(1);
	    }
	} else if (!strncmp("--forkfor=", argv[i], 10)) {
	    unsigned value;
	    if (child != 0) {
		fprintf(stderr, "already forked\n");
		exit(1);
	    }
	    value = strtoul(argv[i]+10, NULL, 0);
	    if (value == 0) {
		goto usage;
	    }
	    child = fork();
	    if (child < 0) {
		perror("unable to fork()");
	    } else if (!child) {
		sleep(value);
		exit(0);
	    }
	} else if (!strncmp("--killit=", argv[i], 9)) {
	    int retval, status;
	    pid_t result;
	    unsigned value;

	    value = strtoul(argv[i]+9, NULL, 0);
	    if (!child) {
		fprintf(stderr, "no forked process to kill\n");
		exit(1);
	    }
	    retval = kill(child, value);
	    if (retval != 0) {
		perror("Unable to kill child process");
		exit(1);
	    }
	    result = waitpid(child, &status, 0);
	    if (result != child) {
		fprintf(stderr, "waitpid didn't match child: %u != %u\n",
			child, result);
		exit(1);
	    }
	    if (WTERMSIG(status) != value) {
		fprintf(stderr, "child terminated with odd signal (%d != %d)\n"
			, value, WTERMSIG(status));
		exit(1);
	    }
	    child = 0;
	} else if (!strncmp("--uid=", argv[i], 6)) {
	    unsigned value;
	    int status;

	    value = strtoul(argv[i]+6, NULL, 0);
	    status = setuid(value);
	    if (status < 0) {
		fprintf(stderr, "Failed to set uid=%u: %s\n",
			value, strerror(errno));
		exit(1);
	    }
	} else if (!strncmp("--cap-uid=", argv[i], 10)) {
	    unsigned value;
	    int status;

	    value = strtoul(argv[i]+10, NULL, 0);
	    status = cap_setuid(value);
	    if (status < 0) {
		fprintf(stderr, "Failed to cap_setuid(%u): %s\n",
			value, strerror(errno));
		exit(1);
	    }
	} else if (!strncmp("--gid=", argv[i], 6)) {
	    unsigned value;
	    int status;

	    value = strtoul(argv[i]+6, NULL, 0);
	    status = setgid(value);
	    if (status < 0) {
		fprintf(stderr, "Failed to set gid=%u: %s\n",
			value, strerror(errno));
		exit(1);
	    }
        } else if (!strncmp("--groups=", argv[i], 9)) {
	  char *ptr, *buf;
	  long length, max_groups;
	  gid_t *group_list;
	  int g_count;

	  length = sysconf(_SC_GETGR_R_SIZE_MAX);
	  buf = calloc(1, length);
	  if (NULL == buf) {
	    fprintf(stderr, "No memory for [%s] operation\n", argv[i]);
	    exit(1);
	  }

	  max_groups = sysconf(_SC_NGROUPS_MAX);
	  group_list = calloc(max_groups, sizeof(gid_t));
	  if (NULL == group_list) {
	    fprintf(stderr, "No memory for gid list\n");
	    exit(1);
	  }

	  g_count = 0;
	  for (ptr = argv[i] + 9; (ptr = strtok(ptr, ","));
	       ptr = NULL, g_count++) {
	    if (max_groups <= g_count) {
	      fprintf(stderr, "Too many groups specified (%d)\n", g_count);
	      exit(1);
	    }
	    if (!isdigit(*ptr)) {
	      struct group *g, grp;
	      getgrnam_r(ptr, &grp, buf, length, &g);
	      if (NULL == g) {
		fprintf(stderr, "Failed to identify gid for group [%s]\n", ptr);
		exit(1);
	      }
	      group_list[g_count] = g->gr_gid;
	    } else {
	      group_list[g_count] = strtoul(ptr, NULL, 0);
	    }
	  }
	  free(buf);
	  if (setgroups(g_count, group_list) != 0) {
	    fprintf(stderr, "Failed to setgroups.\n");
	    exit(1);
	  }
	  free(group_list);
	} else if (!strncmp("--user=", argv[i], 7)) {
	    struct passwd *pwd;
	    const char *user;
	    gid_t groups[MAX_GROUPS];
	    int status, ngroups;

	    user = argv[i] + 7;
	    pwd = getpwnam(user);
	    if (pwd == NULL) {
	      fprintf(stderr, "User [%s] not known\n", user);
	      exit(1);
	    }
	    ngroups = MAX_GROUPS;
	    status = getgrouplist(user, pwd->pw_gid, groups, &ngroups);
	    if (status < 1) {
	      perror("Unable to get group list for user");
	      exit(1);
	    }
	    status = cap_setgroups(pwd->pw_gid, ngroups, groups);
	    if (status != 0) {
		perror("Unable to set group list for user");
		exit(1);
	    }
	    status = cap_setuid(pwd->pw_uid);
	    if (status < 0) {
		fprintf(stderr, "Failed to set uid=%u(user=%s): %s\n",
			pwd->pw_uid, user, strerror(errno));
		exit(1);
	    }
	} else if (!strncmp("--decode=", argv[i], 9)) {
	    unsigned long long value;
	    unsigned cap;
	    const char *sep = "";

	    /* Note, if capabilities become longer than 64-bits we'll need
	       to fixup the following code.. */
	    value = strtoull(argv[i]+9, NULL, 16);
	    printf("0x%016llx=", value);

	    for (cap=0; (cap < 64) && (value >> cap); ++cap) {
		if (value & (1ULL << cap)) {
		    char *ptr;

		    ptr = cap_to_name(cap);
		    if (ptr != NULL) {
			printf("%s%s", sep, ptr);
			cap_free(ptr);
		    } else {
			printf("%s%u", sep, cap);
		    }
		    sep = ",";
		}
	    }
	    printf("\n");
        } else if (!strncmp("--supports=", argv[i], 11)) {
	    cap_value_t cap;

	    if (cap_from_name(argv[i] + 11, &cap) < 0) {
		fprintf(stderr, "cap[%s] not recognized by library\n",
			argv[i] + 11);
		exit(1);
	    }
	    if (!CAP_IS_SUPPORTED(cap)) {
		fprintf(stderr, "cap[%s=%d] not supported by kernel\n",
			argv[i] + 11, cap);
		exit(1);
	    }
	} else if (!strcmp("--print", argv[i])) {
	    arg_print();
	} else if ((!strcmp("--", argv[i])) || (!strcmp("==", argv[i]))) {
	    argv[i] = strdup(argv[i][0] == '-' ? shell : argv[0]);
	    argv[argc] = NULL;
	    execve(argv[i], argv+i, envp);
	    fprintf(stderr, "execve '%s' failed!\n", shell);
	    exit(1);
	} else if (!strncmp("--shell=", argv[i], 8)) {
	    shell = argv[i]+8;
	} else if (!strncmp("--has-p=", argv[i], 8)) {
	    cap_value_t cap;
	    cap_flag_value_t enabled;
	    cap_t orig;

	    if (cap_from_name(argv[i]+8, &cap) < 0) {
		fprintf(stderr, "cap[%s] not recognized by library\n",
			argv[i] + 8);
		exit(1);
	    }
	    orig = cap_get_proc();
	    if (cap_get_flag(orig, cap, CAP_PERMITTED, &enabled) || !enabled) {
		fprintf(stderr, "cap[%s] not permitted\n", argv[i]+8);
		exit(1);
	    }
	    cap_free(orig);
	} else if (!strncmp("--has-i=", argv[i], 8)) {
	    cap_value_t cap;
	    cap_flag_value_t enabled;
	    cap_t orig;

	    if (cap_from_name(argv[i]+8, &cap) < 0) {
		fprintf(stderr, "cap[%s] not recognized by library\n",
			argv[i] + 8);
		exit(1);
	    }
	    orig = cap_get_proc();
	    if (cap_get_flag(orig, cap, CAP_INHERITABLE, &enabled)
		|| !enabled) {
		fprintf(stderr, "cap[%s] not inheritable\n", argv[i]+8);
		exit(1);
	    }
	    cap_free(orig);
	} else if (!strncmp("--has-a=", argv[i], 8)) {
	    cap_value_t cap;
	    if (cap_from_name(argv[i]+8, &cap) < 0) {
		fprintf(stderr, "cap[%s] not recognized by library\n",
			argv[i] + 8);
		exit(1);
	    }
	    if (!cap_get_ambient(cap)) {
		fprintf(stderr, "cap[%s] not in ambient vector\n", argv[i]+8);
		exit(1);
	    }
	} else if (!strncmp("--is-uid=", argv[i], 9)) {
	    unsigned value;
	    uid_t uid;
	    value = strtoul(argv[i]+9, NULL, 0);
	    uid = getuid();
	    if (uid != value) {
		fprintf(stderr, "uid: got=%d, want=%d\n", uid, value);
		exit(1);
	    }
	} else if (!strncmp("--is-gid=", argv[i], 9)) {
	    unsigned value;
	    gid_t gid;
	    value = strtoul(argv[i]+9, NULL, 0);
	    gid = getgid();
	    if (gid != value) {
		fprintf(stderr, "gid: got=%d, want=%d\n", gid, value);
		exit(1);
	    }
	} else if (!strncmp("--iab=", argv[i], 6)) {
	    cap_iab_t iab = cap_iab_from_text(argv[i]+6);
	    if (iab == NULL) {
		fprintf(stderr, "iab: '%s' malformed\n", argv[i]+6);
		exit(1);
	    }
	    if (cap_iab_set_proc(iab)) {
		perror("unable to set IAP vectors");
		exit(1);
	    }
	    cap_free(iab);
	} else {
	usage:
	    printf("usage: %s [args ...]\n"
		   "  --help         this message (or try 'man capsh')\n"
		   "  --print        display capability relevant state\n"
		   "  --decode=xxx   decode a hex string to a list of caps\n"
		   "  --supports=xxx exit 1 if capability xxx unsupported\n"
		   "  --has-p=xxx    exit 1 if capability xxx not permitted\n"
		   "  --has-i=xxx    exit 1 if capability xxx not inheritable\n"
		   "  --drop=xxx     remove xxx,.. capabilities from bset\n"
		   "  --dropped=xxx  exit 1 unless bounding cap xxx dropped\n"
		   "  --has-ambient  exit 1 unless ambient vector supported\n"
		   "  --has-a=xxx    exit 1 if capability xxx not ambient\n"
		   "  --addamb=xxx   add xxx,... capabilities to ambient set\n"
		   "  --delamb=xxx   remove xxx,... capabilities from ambient\n"
		   "  --noamb        reset (drop) all ambient capabilities\n"
		   "  --caps=xxx     set caps as per cap_from_text()\n"
		   "  --inh=xxx      set xxx,.. inheritiable set\n"
		   "  --secbits=<n>  write a new value for securebits\n"
		   "  --iab=...      use cap_iab_from_text() to set iab\n"
		   "  --keep=<n>     set keep-capabability bit to <n>\n"
		   "  --uid=<n>      set uid to <n> (hint: id <username>)\n"
		   "  --cap-uid=<n>  libcap cap_setuid() to change uid\n"
		   "  --is-uid=<n>   exit 1 if uid != <n>\n"
		   "  --gid=<n>      set gid to <n> (hint: id <username>)\n"
		   "  --is-gid=<n>   exit 1 if gid != <n>\n"
		   "  --groups=g,... set the supplemental groups\n"
                   "  --user=<name>  set uid,gid and groups to that of user\n"
		   "  --chroot=path  chroot(2) to this path\n"
		   "  --modes        list libcap named capability modes\n"
		   "  --mode=<xxx>   set capability mode to <xxx>\n"
		   "  --inmode=<xxx> exit 1 if current mode is not <xxx>\n"
		   "  --killit=<n>   send signal(n) to child\n"
		   "  --forkfor=<n>  fork and make child sleep for <n> sec\n"
		   "  --shell=/xx/yy use /xx/yy instead of " SHELL " for --\n"
		   "  ==             re-exec(capsh) with args as for --\n"
		   "  --             remaining arguments are for " SHELL "\n"
		   "                 (without -- [%s] will simply exit(0))\n",
		   argv[0], argv[0]);

	    exit(strcmp("--help", argv[i]) != 0);
	}
    }

    exit(0);
}
