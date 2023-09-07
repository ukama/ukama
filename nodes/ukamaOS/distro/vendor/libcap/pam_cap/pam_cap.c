/*
 * Copyright (c) 1999,2007,19,20 Andrew G. Morgan <morgan@kernel.org>
 *
 * The purpose of this module is to enforce inheritable, bounding and
 * ambient capability sets for a specified user.
 */

/* #define DEBUG */

#ifndef _DEFAULT_SOURCE
#define _DEFAULT_SOURCE
#endif

#include <errno.h>
#include <grp.h>
#include <limits.h>
#include <pwd.h>
#include <stdarg.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <syslog.h>
#include <sys/capability.h>
#include <sys/types.h>
#include <linux/limits.h>

#include <security/pam_modules.h>
#include <security/_pam_macros.h>

#define USER_CAP_FILE           "/etc/security/capability.conf"
#define CAP_FILE_BUFFER_SIZE    4096
#define CAP_FILE_DELIMITERS     " \t\n"

struct pam_cap_s {
    int debug;
    const char *user;
    const char *conf_filename;
};

/*
 * load_groups obtains the list all of the groups associated with the
 * requested user: gid & supplemental groups.
 */
static int load_groups(const char *user, char ***groups, int *groups_n) {
    struct passwd *pwd;
    gid_t grps[NGROUPS_MAX];
    int ngrps = NGROUPS_MAX;

    *groups = NULL;
    *groups_n = 0;

    pwd = getpwnam(user);
    if (pwd == NULL) {
	return -1;
    }

    /* must include at least pwd->pw_gid, hence < 1 test. */
    if (getgrouplist(user, pwd->pw_gid, grps, &ngrps) < 1) {
	return -1;
    }

    *groups = calloc(ngrps, sizeof(char *));
    int g_n = 0, i;
    for (i = 0; i < ngrps; i++) {
	const struct group *g = getgrgid(grps[i]);
	if (g == NULL) {
	    continue;
	}
	D(("noting [%s] is a member of [%s]", user, g->gr_name));
	(*groups)[g_n++] = strdup(g->gr_name);
    }

    *groups_n = g_n;
    return 0;
}

/* obtain the inheritable capabilities for the current user */

static char *read_capabilities_for_user(const char *user, const char *source)
{
    char *cap_string = NULL;
    char buffer[CAP_FILE_BUFFER_SIZE], *line;
    char **groups;
    int groups_n;
    FILE *cap_file;

    if (load_groups(user, &groups, &groups_n)) {
	D(("unknown user [%s]", user));
	return NULL;
    }

    cap_file = fopen(source, "r");
    if (cap_file == NULL) {
	D(("failed to open capability file"));
	goto defer;
    }

    int found_one = 0;
    while (!found_one &&
	   (line = fgets(buffer, CAP_FILE_BUFFER_SIZE, cap_file))) {
	const char *cap_text;

	char *next = NULL;
	cap_text = strtok_r(line, CAP_FILE_DELIMITERS, &next);

	if (cap_text == NULL) {
	    D(("empty line"));
	    continue;
	}
	if (*cap_text == '#') {
	    D(("comment line"));
	    continue;
	}

	/*
	 * Explore whether any of the ids are a match for the current
	 * user.
	 */
	while ((line = strtok_r(next, CAP_FILE_DELIMITERS, &next))) {
	    if (strcmp("*", line) == 0) {
		D(("wildcard matched"));
		found_one = 1;
		break;
	    }

	    if (strcmp(user, line) == 0) {
		D(("exact match for user"));
		found_one = 1;
		break;
	    }

	    if (line[0] != '@') {
		D(("user [%s] is not [%s] - skipping", user, line));
	    }

	    int i;
	    for (i=0; i < groups_n; i++) {
		if (!strcmp(groups[i], line+1)) {
		    D(("user group matched [%s]", line));
		    found_one = 1;
		    break;
		}
	    }
	    if (found_one) {
		break;
	    }
	}

	if (found_one) {
	    cap_string = strdup(cap_text);
	    D(("user [%s] matched - caps are [%s]", user, cap_string));
	}

	cap_text = NULL;
	line = NULL;
    }

    fclose(cap_file);

defer:
    memset(buffer, 0, CAP_FILE_BUFFER_SIZE);

    int i;
    for (i = 0; i < groups_n; i++) {
	char *g = groups[i];
	_pam_overwrite(g);
	_pam_drop(g);
    }
    if (groups != NULL) {
	memset(groups, 0, groups_n * sizeof(char *));
	_pam_drop(groups);
    }

    return cap_string;
}

/*
 * Set capabilities for current process to match the current
 * permitted+executable sets combined with the configured inheritable
 * set.
 */
static int set_capabilities(struct pam_cap_s *cs)
{
    cap_t cap_s;
    char *conf_caps;
    int ok = 0;
    cap_iab_t iab;

    cap_s = cap_get_proc();
    if (cap_s == NULL) {
	D(("your kernel is capability challenged - upgrade: %s",
	   strerror(errno)));
	return 0;
    }

    conf_caps =	read_capabilities_for_user(cs->user,
					   cs->conf_filename
					   ? cs->conf_filename:USER_CAP_FILE );
    if (conf_caps == NULL) {
	D(("no capabilities found for user [%s]", cs->user));
	goto cleanup_cap_s;
    }

    ssize_t conf_caps_length = strlen(conf_caps);
    if (!strcmp(conf_caps, "all")) {
	/*
	 * all here is interpreted as no change/pass through, which is
	 * likely to be the same as none for sensible system defaults.
	 */
	ok = 1;
	goto cleanup_conf;
    }

    if (!strcmp(conf_caps, "none")) {
	/* clearing CAP_INHERITABLE will also clear the ambient caps,
	 * but for legacy reasons we do not alter the bounding set. */
	cap_clear_flag(cap_s, CAP_INHERITABLE);
	if (!cap_set_proc(cap_s)) {
	    ok = 1;
	}
	goto cleanup_cap_s;
    }

    iab = cap_iab_from_text(conf_caps);
    if (iab == NULL) {
	D(("unable to parse the IAB [%s] value", conf_caps));
	goto cleanup_conf;
    }

    if (!cap_iab_set_proc(iab)) {
	D(("able to set the IAB [%s] value", conf_caps));
	ok = 1;
    }
    cap_free(iab);

cleanup_conf:
    memset(conf_caps, 0, conf_caps_length);
    _pam_drop(conf_caps);

cleanup_cap_s:
    if (cap_s) {
	cap_free(cap_s);
	cap_s = NULL;
    }
    return ok;
}

/* log errors */

static void _pam_log(int err, const char *format, ...)
{
    va_list args;

    va_start(args, format);
    openlog("pam_cap", LOG_CONS|LOG_PID, LOG_AUTH);
    vsyslog(err, format, args);
    va_end(args);
    closelog();
}

static void parse_args(int argc, const char **argv, struct pam_cap_s *pcs)
{
    /* step through arguments */
    for (; argc-- > 0; ++argv) {
	if (!strcmp(*argv, "debug")) {
	    pcs->debug = 1;
	} else if (!strncmp(*argv, "config=", 7)) {
	    pcs->conf_filename = 7 + *argv;
	} else {
	    _pam_log(LOG_ERR, "unknown option; %s", *argv);
	}
    }
}

/*
 * pam_sm_authenticate parses the config file with respect to the user
 * being authenticated and determines if they are covered by any
 * capability inheritance rules.
 */
int pam_sm_authenticate(pam_handle_t *pamh, int flags,
			int argc, const char **argv)
{
    int retval;
    struct pam_cap_s pcs;
    char *conf_caps;

    memset(&pcs, 0, sizeof(pcs));
    parse_args(argc, argv, &pcs);

    retval = pam_get_user(pamh, &pcs.user, NULL);
    if (retval == PAM_CONV_AGAIN) {
	D(("user conversation is not available yet"));
	memset(&pcs, 0, sizeof(pcs));
	return PAM_INCOMPLETE;
    }

    if (retval != PAM_SUCCESS) {
	D(("pam_get_user failed: %s", pam_strerror(pamh, retval)));
	memset(&pcs, 0, sizeof(pcs));
	return PAM_AUTH_ERR;
    }

    conf_caps =	read_capabilities_for_user(pcs.user,
					   pcs.conf_filename
					   ? pcs.conf_filename:USER_CAP_FILE );
    memset(&pcs, 0, sizeof(pcs));

    if (conf_caps) {
	D(("it appears that there are capabilities for this user [%s]",
	   conf_caps));

	/* We could also store this as a pam_[gs]et_data item for use
	   by the setcred call to follow. As it is, there is a small
	   race associated with a redundant read. Oh well, if you
	   care, send me a patch.. */

	_pam_overwrite(conf_caps);
	_pam_drop(conf_caps);

	return PAM_SUCCESS;

    } else {

	D(("there are no capabilities restrctions on this user"));
	return PAM_IGNORE;

    }
}

/*
 * pam_sm_setcred applies inheritable capabilities loaded by the
 * pam_sm_authenticate pass for the user.
 */
int pam_sm_setcred(pam_handle_t *pamh, int flags,
		   int argc, const char **argv)
{
    int retval;
    struct pam_cap_s pcs;

    if (!(flags & (PAM_ESTABLISH_CRED | PAM_REINITIALIZE_CRED))) {
	D(("we don't handle much in the way of credentials"));
	return PAM_IGNORE;
    }

    memset(&pcs, 0, sizeof(pcs));
    parse_args(argc, argv, &pcs);

    retval = pam_get_item(pamh, PAM_USER, (const void **)&pcs.user);
    if ((retval != PAM_SUCCESS) || (pcs.user == NULL) || !(pcs.user[0])) {
	D(("user's name is not set"));
	return PAM_AUTH_ERR;
    }

    retval = set_capabilities(&pcs);
    memset(&pcs, 0, sizeof(pcs));

    return (retval ? PAM_SUCCESS:PAM_IGNORE );
}
