#include <stdio.h>
#include <string.h>
#include <sys/capability.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

/*
 * tests for cap_launch.
 */

#define MORE_THAN_ENOUGH 20
#define NO_MORE 1

struct test_case_s {
    int pass_on;
    const char *chroot;
    uid_t uid;
    gid_t gid;
    int ngroups;
    const gid_t groups[MORE_THAN_ENOUGH];
    const char *args[MORE_THAN_ENOUGH];
    const char **envp;
    const char *iab;
    cap_mode_t mode;
    int result;
};

#ifdef WITH_PTHREADS
#include <pthread.h>
#else /* WITH_PTHREADS */
#endif /* WITH_PTHREADS */

int main(int argc, char **argv) {
    static struct test_case_s vs[] = {
	{
	    .args = { "../progs/capsh", "--", "-c", "echo hello" },
	    .result = 0
	},
	{
	    .args = { "../progs/capsh", "--is-uid=123" },
	    .result = 256
	},
	{
	    .args = { "../progs/capsh", "--is-uid=123" },
	    .result = 0,
	    .uid = 123,
	},
	{
	    .args = { "../progs/capsh", "--is-gid=123" },
	    .result = 0,
	    .gid = 123,
	    .ngroups = 1,
	    .groups = { 456 },
	    .iab = "",
	},
	{
	    .args = { "../progs/capsh", "--dropped=cap_chown",
		      "--has-i=cap_chown" },
	    .result = 0,
	    .iab = "!%cap_chown"
	},
	{
	    .args = { "../progs/capsh", "--dropped=cap_chown",
		      "--has-i=cap_chown", "--is-uid=234",
		      "--has-a=cap_chown", "--has-p=cap_chown" },
	    .uid = 234,
	    .result = 0,
	    .iab = "!^cap_chown"
	},
	{
	    .args = { "../progs/capsh", "--inmode=NOPRIV" },
	    .result = 0,
	    .mode = CAP_MODE_NOPRIV
	},
	{
	    .args = { "/noop" },
	    .result = 0,
	    .chroot = ".",
	},
	{
	    .pass_on = NO_MORE
	},
    };

    cap_t orig = cap_get_proc();

    int success = 1, i;
    for (i=0; vs[i].pass_on != NO_MORE; i++) {
	const struct test_case_s *v = &vs[i];
	printf("[%d] test should %s\n", i,
	       v->result ? "generate error" : "work");
	cap_launch_t attr = cap_new_launcher(v->args[0], v->args, v->envp);
	if (v->chroot) {
	    cap_launcher_set_chroot(attr, v->chroot);
	}
	if (v->uid) {
	    cap_launcher_setuid(attr, v->uid);
	}
	if (v->gid) {
	    cap_launcher_setgroups(attr, v->gid, v->ngroups, v->groups);
	}
	if (v->iab) {
	    cap_iab_t iab = cap_iab_from_text(v->iab);
	    if (iab == NULL) {
		fprintf(stderr, "[%d] failed to decode iab [%s]", i, v->iab);
		perror(":");
		success = 0;
		continue;
	    }
	    cap_iab_t old = cap_launcher_set_iab(attr, iab);
	    if (cap_free(old)) {
		fprintf(stderr, "[%d] failed to decode iab [%s]", i, v->iab);
		perror(":");
		success = 0;
		continue;
	    }
	}
	if (v->mode) {
	    cap_launcher_set_mode(attr, v->mode);
	}

	pid_t child = cap_launch(attr, NULL);

	if (child <= 0) {
	    fprintf(stderr, "[%d] failed to launch", i);
	    perror(":");
	    success = 0;
	    continue;
	}
	if (cap_free(attr)) {
	    fprintf(stderr, "[%d] failed to free launcher", i);
	    perror(":");
	    success = 0;
	}
	int result;
	int ret = waitpid(child, &result, 0);
	if (ret != child) {
	    fprintf(stderr, "[%d] failed to wait", i);
	    perror(":");
	    success = 0;
	    continue;
	}
	if (result != v->result) {
	    fprintf(stderr, "[%d] bad result: got=%d want=%d", i, result,
		    v->result);
	    perror(":");
	    success = 0;
	    continue;
	}
    }

    cap_t final = cap_get_proc();
    if (cap_compare(orig, final)) {
	char *was = cap_to_text(orig, NULL);
	char *is = cap_to_text(final, NULL);
	printf("cap_launch_test: orig:'%s' != final:'%s'\n", was, is);
	cap_free(is);
	cap_free(was);
	success = 0;
    }
    cap_free(final);
    cap_free(orig);

    if (success) {
	printf("cap_launch_test: PASSED\n");
    } else {
	printf("cap_launch_test: FAILED\n");
    }
}
