#define _GNU_SOURCE
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <sys/reboot.h>
#include <errno.h>
#include <stdlib.h>
#include <string.h>

/*
 * LD_PRELOAD library used ONLY during tests.
 * Goal: keep device.d behavior intact but prevent dangerous / environment-dependent actions:
 *  - reboot()
 *  - setuid(0)
 *  - fork/execv/waitpid used by actions_tower.c to run absolute-path scripts
 *
 * We simulate successful command execution without spawning processes.
 */

static pid_t g_fake_pid = 4242;

pid_t fork(void) {
    /* Pretend we're the parent and that a child was spawned. */
    g_fake_pid++;
    return g_fake_pid;
}

int execv(const char *path, char *const argv[]) {
    (void)path;
    (void)argv;
    errno = ENOENT;
    /* If someone does call execv unexpectedly, fail safely. */
    return -1;
}

pid_t waitpid(pid_t pid, int *status, int options) {
    (void)options;
    if (pid <= 0) {
        errno = EINVAL;
        return -1;
    }

    if (status) {
        /* Encode exit status 0 */
        *status = 0;
    }

    return pid;
}

int reboot(int cmd) {
    (void)cmd;
    /* Simulate success without rebooting */
    return 0;
}

int setuid(uid_t uid) {
    (void)uid;
    return 0;
}

void sync(void) {
    /* no-op */
}
