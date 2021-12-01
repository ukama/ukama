/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* wrapper function to setup linux networking using IP command
 */

#include <string.h>
#include <unistd.h>
#include <signal.h>
#include <sys/wait.h>
#include <sys/types.h>
#include <errno.h>
#include <stdlib.h>

#include "log.h"
#include "ipnet.h"

static int ipnet_run(char *exec, char *args, int testFlag);

/*
 * ipnet_run --
 */
static int ipnet_run(char *exec, char *args, int testFlag) {

  int ret=FALSE;
  sigset_t blockMask, origMask;
  pid_t childPid;
  char *cmd=NULL;
  int status, exitStatus, size;

  size = strlen(exec) + strlen(args) + 2;
  cmd = (char *)malloc(size);
  if (!cmd) {
    log_error("Error allocating memory of size: %d",size);
    return FALSE;
  }

  /* As per system(3), we gotta block SIGCHLD */
  sigemptyset(&blockMask);
  sigaddset(&blockMask, SIGCHLD);
  sigprocmask(SIG_BLOCK, &blockMask, &origMask);

  sprintf(cmd, "%s %s", exec, args);

  switch (childPid = fork()) {
  case -1: 
    status = -1;
    break;          /* Reset SIGCHLD */

  case 0: /* Child process */
    sigprocmask(SIG_SETMASK, &origMask, NULL);
    execl("/bin/sh", "sh", "-c", exec, args, (char *) NULL);
    _exit(127);

  default: /* Parent */
    while (waitpid(childPid, &status, 0) == -1) {
      if (errno != EINTR) {
	log_error("Script execution failure. Cmd: %s Args: %s Error: %s",
		  exec, args, strerror(errno));
	status = -1;
	break;
      }
    }
    break;
  }

  if (testFlag) {
    ret = status;
    goto done;
  }

  if (WIFEXITED(status)) { /* proper termination by calling exit */
    exitStatus = WEXITSTATUS(status);
    switch(exitStatus) {
    case 100:
      log_error("Invalid arguments for script: %s Args: %s", exec, args);
      break;
    case 101:
      log_error("Invalid PID for cspace. Does not exsist. Args: %s", args);
      break;
    case 102:
      log_error("Invalid Network Interface. Args: %s", args);
      break;
    case 0:
      log_debug("Script sucessfully executed. Cmd: %s Args: %s", cmd, args);
      ret = TRUE;
      break;
    default:
      log_error("Script return invalid exit code: %d Cmd: %s Args: %s",
		exitStatus, exec, args);
      break;
    }
  }

 done:
  /* reset SIGCHLD */
  sigprocmask(SIG_SETMASK, &origMask, NULL);
  free(cmd);

  return ret;
}

/*
 * ipnet_setup --
 */
int ipnet_setup(int type, char *brName, char *iface, char *spName, pid_t pid) {

  char *exec=NULL, *args=NULL, *dev=NULL;
  char *arg1=NULL, *arg2=NULL, *arg3=NULL;
  char pidStr[21]={0};
  int size, ret;

  size = strlen(NET_EXEC) + strlen(PATH) + 2;

  exec  = (char *)calloc(1, size);
  args = (char *)calloc(1, 1024);

  if (!exec || !args) {
    log_error("Error allocating memory of size: %d 1024", size);
    return FALSE;
  }

  sprintf(exec, "%s/%s", PATH, NET_EXEC);
  sprintf(pidStr, "%ld", (long)pid);

  if (type == IPNET_DEV_TYPE_BRIDGE) {
    /* setup_space_network br <name> */
    dev = IPNET_DEV_BRIDGE;
    arg1 = iface;
    arg2 = brName;
    arg3 = "";
  } else if (type == IPNET_DEV_TYPE_CSPACE) {
    /* setup_space_network ns <space_name> <pid> */
    dev  = IPNET_DEV_CSPACE;
    arg1 = pidStr;
    arg2 = spName;
    arg3 = brName;
  } else {
    log_error("Invalid type: %d Ignoring and not running %s", type, NET_EXEC);
    free(exec);
    free(args);
    return FALSE;
  }

  sprintf(args, "--add %s %s %s %s", dev, arg1, arg2, arg3);

  ret = ipnet_run(exec, args, FALSE);

  free(exec);
  free(args);

  return ret;
}

/*
 * ipnet_test -- For a given cspace, check if the networking is setup correctly
 *               It should be able to ping test IP of 172.217.6.78
 *
 */

int ipnet_test(char *spName) {

  char args[1024] = {0};
  int status, exitStatus, ret=FALSE;

  if (spName == NULL) return FALSE;

  sprintf(args, "netns exec %s %s %s", spName, PING_BIN, TEST_IP);

  status = ipnet_run(IP_BIN, args, TRUE);

  if (WIFEXITED(status)) { /* proper termination by calling exit */
    exitStatus = WEXITSTATUS(status);
    switch(exitStatus) {
    case 2:
      log_error("Unable to reach test IP: %s Args: %s", TEST_IP, args);
      break;
    case 1:
      log_error("Recevied no reply from test IP: %s Args: %s", TEST_IP, args);
      break;
    case 0:
      log_debug("Network setup is correct. Test IP %s is reachable",
		TEST_IP);
      ret = TRUE;
      break;
    default:
      log_error("Test return invalid exit code: %d Cmd: %s Args: %s",
		exitStatus, args);
      break;
    }
  }

  return ret;
}
