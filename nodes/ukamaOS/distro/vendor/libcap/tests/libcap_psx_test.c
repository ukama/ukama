#ifndef _DEFAULT_SOURCE
#define _DEFAULT_SOURCE
#endif

#include <errno.h>
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/capability.h>
#include <sys/psx_syscall.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

static void *thread_fork_exit(void *data) {
    usleep(1234);
    pid_t pid = fork();
    cap_t start = cap_get_proc();
    if (pid == 0) {
	cap_set_proc(start);
	exit(0);
    }
    int res;
    if (waitpid(pid, &res, 0) != pid || res != 0) {
	printf("FAILED: pid=%d wait returned %d and/or error: %d\n",
	       pid, res, errno);
	exit(1);
    }
    cap_set_proc(start);
    return NULL;
}

int main(int argc, char **argv) {
    int i;
    printf("hello libcap and libpsx ");
    fflush(stdout);
    cap_t start = cap_get_proc();
    pthread_t ignored[10];
    for (i = 0; i < 10; i++) {
	pthread_create(&ignored[i], NULL, thread_fork_exit, NULL);
    }
    for (i = 0; i < 10; i++) {
	printf(".");     /* because of fork, this may print double */
	fflush(stdout);  /* try to limit the above effect */
	cap_set_proc(start);
	usleep(1000);
    }
    printf(" PASSED\n");
    exit(0);
}
