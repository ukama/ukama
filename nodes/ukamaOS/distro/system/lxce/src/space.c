/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* Common functions related to cloning and space setup */

#define _GNU_SOURCE
#include <string.h>
#include <errno.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/prctl.h>
#include <sched.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <sys/mount.h>
#include <sys/syscall.h>
#include <fcntl.h>
#include <grp.h>
#include <signal.h>
#include <uuid/uuid.h>

#include "space.h"
#include "cspace.h"
#include "manifest.h"
#include "log.h"
#include "capp.h"
#include "capp_runtime.h"
#include "utils.h"

static int prepare_child_map_files(char *areaType, int sockets[2], pid_t pid,
				   char *name);
static int pivot_root(const char *new, const char *old);

/*
 * prepare_child_map_files -- setup map files (uid_map, gid_map and setgroups)
 *
 */
static int prepare_child_map_files(char *areaType,
				   int *sockets, pid_t pid,
				   char *name) {

  int ns=0, size, i, ret;
  int mapFd = 0;
  char *mapFiles[] = {"uid_map", "gid_map"};
  char mapPath[LXCE_MAX_PATH] = {0};
  char runMe[LXCE_MAX_PATH] = {0};

  if (strcmp(areaType, AREA_TYPE_CAPP) == 0) { /* XXX */
    goto proceed;
  }

  size = read(sockets[PARENT_SOCKET], &ns, sizeof(ns));
  if (size != sizeof(ns)) {
    log_error("%s: %s error reading client socket. Expected size: %d Got: %d",
	      areaType, name, sizeof(ns), size);
    return FALSE;
  }

  if (!ns) {
    log_error("%s: %s Child user namespace issue. unshare.", areaType, name);
    return FALSE;
  }

  sprintf(runMe, "/bin/mkdir -p %s/%s/proc/%d/", DEF_CSPACE_ROOTFS_PATH, name,
	  (int)pid);
  log_debug("Running command: %s", runMe);
  if ((ret = system(runMe)) < 0) {
    log_error("Unable to execute cmd %s for space: %s Code: %d", runMe, name,
	      ret);
    return FALSE;
  }

  for (i=0; i<2; i++) {
    sprintf(mapPath, "%s/%s/proc/%d/%s", DEF_CSPACE_ROOTFS_PATH, name, pid,
	    mapFiles[i]);

    if ((mapFd = open(mapPath, O_CREAT | O_WRONLY, 0644)) == -1) {
      log_error("%s: %s error opening map file: %s Error: %s", areaType,
		name, mapPath, strerror(errno));
      return FALSE;
    }

    if (dprintf(mapFd, "0 %d %d\n", USER_NS_OFFSET, USER_NS_COUNT) == -1) {
      log_error("%s: %s Error writing to map file: %s Error: %s", areaType,
		name, mapPath, strerror(errno));
      close(mapFd);
      return FALSE;
    }

    close(mapFd);
  }

 proceed:
  /* Inform child, it can proceed. */
  size = write(sockets[PARENT_SOCKET], &(int){TRUE}, sizeof(int));
  if (size != sizeof(int)) {
    log_error("%s: %s Error writing to child socket", areaType, name);
    log_error("Expected size: %d Wrote: %d", sizeof(int), size);
    return FALSE;
  }

  return TRUE;
}

/*
 * create_space --
 *
 */
int create_space(char *areaType,
		 int *sockets, int namespaces,
		 char *name, pid_t *pid,
		 int (*func)(void *), void *arg) {

  char *stack=NULL;

  if (!name || !areaType) return FALSE;

  /* Create socket pairs.
   * Re: SOCK_SEQPACKET:
   * http://urchin.earth.li/~twic/Sequenced_Packets_Over_Ordinary_TCP.html
   */
  if (socketpair(AF_LOCAL, SOCK_SEQPACKET, 0, sockets)) {
    log_error("%s: %s Error creating socket pair", areaType, name);
    return FALSE;
  }

  if (!(stack = malloc(SPACE_STACK_SIZE))) {
    log_error("%s: %s Error allocating stack of size: %d", areaType,
	      name, SPACE_STACK_SIZE);
    return FALSE;
  }

  /* clone with proper flags for namespaces */
  if (strcmp(areaType, AREA_TYPE_CAPP)==0) { /* XXX, Fix me */
    *pid = clone(func, stack + SPACE_STACK_SIZE, SIGCHLD, arg);
  } else {
    *pid = clone(func, stack + SPACE_STACK_SIZE, SIGCHLD | namespaces, arg);
  }

  if (*pid == -1) {
    log_error("%s: %s Unable to clone cInit. Error :%s", areaType, name,
	      strerror(errno));
    return FALSE;
  }

  if (*pid > 0 ) {
    /* Close child socket */
    close(sockets[CHILD_SOCKET]);

    /* prepare child process gid/uid map files. */
    if (prepare_child_map_files(areaType, sockets, *pid, name) == FALSE) {
      log_error("%s: error preparing map files for child. Terminating it",
		areaType);
      kill(*pid, SIGKILL); /* Kill child process */
      close(sockets[PARENT_SOCKET]);
      close(sockets[CHILD_SOCKET]);
      return FALSE;
    }
  }

  return TRUE;
}

/*
 * pivot_root -- wrapper for sys call
 *
 */
static int pivot_root(const char *new, const char *old) {
  return syscall(SYS_pivot_root, new, old);
}

/*
 * setup_mounts --
 *
 */
int setup_mounts(char *areaType, char *rootfs, char *name) {

  int ret=FALSE;
  char tempMount[] = "/tmp/tmp.ukama.XXXXXX"; /* last 6 char needs to be X */
  char oldRoot[]   = "/tmp/tmp.ukama.XXXXXX/oldroot.XXXXXX"; /* same */
  char *oldRootDir=NULL;

  if (mount(NULL, "/", NULL, MS_REC | MS_PRIVATE, NULL)) {
    log_error("%s: %s Failed to remount as MS_PRIVATE. Error: %s",
	      areaType, name, strerror(errno));
    return ret;
  }

  /* make temp and bind mount */
  if (!mkdtemp(tempMount)) {
    log_error("%s: %s Failed to make temp dir: %s. Error: %s",
	      name, areaType, tempMount, strerror(errno));
    return FALSE;
  }

  if (mount(rootfs, tempMount, NULL, MS_BIND | MS_PRIVATE, NULL)) {
    log_error("%s: %s Failed to do bind mount. %s %s Error: %s",
	      areaType, name, rootfs, tempMount, strerror(errno));
    return FALSE;
  }

  memcpy(oldRoot, tempMount, sizeof(tempMount) - 1);
  if (!mkdtemp(oldRoot)) {
    log_error("%s: Failed to create old Root directory. %s Error :%s",
	      areaType, oldRoot, strerror(errno));
    return FALSE;
  }

  /* pivot root */
  if (pivot_root(tempMount, oldRoot)) {
    log_error("%s: Failed to pivot_root from %s to %s. Error: %s",
	      areaType, oldRoot, rootfs, strerror(errno));
    return FALSE;
  }

  log_debug("%s: Pivot root sucessfully done to %s", areaType, rootfs);

  /* clean up */
  oldRootDir = basename(oldRoot);
  char rmv[sizeof(oldRoot) + 1] = { "/" };
  strcpy(&rmv[1], oldRootDir);

  if (chdir("/")) {
    log_error("%s: error changing director to / after pivot", areaType);
    return FALSE;
  }

  if (umount2(rmv, MNT_DETACH)) {
    log_error("%s: failed to umount/rm old root: %s. Error: %s",
	      areaType, oldRoot, strerror(errno));
    return FALSE;
  }

  if (rmdir(rmv)) {
    log_error("%s: failed to remove old root directory: %s", areaType, rmv);
    return FALSE;
  }

  return TRUE;
}
