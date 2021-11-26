/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * capp creation related functions
 */
#define _GNU_SOURCE
#include <string.h>
#include <errno.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/capability.h>
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
#include <errno.h>
#include <jansson.h>

#include "space.h"
#include "cspace.h"
#include "log.h"
#include "capp.h"
#include "capp_config.h"
#include "capp_runtime.h"

static int exec_capp(CApp *capp);
static int capp_init_clone(void *arg);
static int create_capp(CApp *capp);

/*
 * create_and_run_capps -- Run all the apps from the PEND list
 *
 */

int create_and_run_capps(void *args) {

  CAppList *list=NULL;
  CApps *apps = args;
  CApp *capp=NULL;
  char *fileName;

  if (!apps) return FALSE;

  if (!apps->pend) return TRUE; /* Nothing to do. */

  for (list = apps->pend; list; list=list->next) {

    if (!list->capp) continue;

    /* Steps are:
     * 1. Check the existance of config.json
     * 2. Parse and process config.json content
     * 3. clone
     * 4. setup the right namespaces (as per config.json)
     * 5. exec into the process
     * 6. Parent process monitor the status of the process
     */

    capp = list->capp;

    if (!valid_path(capp->params->path)) {
      log_error("Invalid path for the capp. Name: %s Tag: %s Path: %s",
		capp->params->name, capp->params->tag, capp->params->path);
      /* Move this capp from PEND to ERROR. */

    }

    fileName = (char *)calloc(1, strlen(capp->params->path) +
			      strlen(DEF_CONFIG) + 2); /* +1 for null and '/'*/
    sprintf(fileName, "%s/%s", capp->params->path, DEF_CONFIG);

    if (!process_capp_config_file(capp->config, fileName)) {
      free(fileName);
      continue;
    }

    create_capp(capp);
  }

  return TRUE;
}

/*
 * exec_capp --
 *
 */
static int exec_capp(CApp *capp) {

  CAppProc *process;

  if (!capp || !capp->config) return FALSE;

  process = capp->config->process;

  if (!process->exec) return FALSE;

  execvpe(process->exec, /* binary */
	  process->argv, /* arguments to the binary */
	  process->env); /* environment variables */

  /* Only if there was an error */
  return errno;
}

/*
 * capp_init_clone --
 *
 */
static int capp_init_clone(void *arg) {

  int ret;
  CApp *capp = (CApp *)arg;
  char *hostName=NULL;

  if (capp->config->hostName) {
    hostName = capp->config->hostName;
  } else {
    hostName = CAPP_DEFAULT_HOSTNAME;
  }

  if (sethostname(hostName, strlen(hostName))) {
    log_error("CApp: %s Error setting host name: %s", capp->params->name,
              hostName);
    return FALSE;
  }

  /* execv into program */
  ret = exec_capp(capp);

  /* An error has occured. Inform the parent process over socket and exit. */
  if (write(capp->runtime->sockets[0], &ret, sizeof(ret)) != sizeof(ret)) {
    log_error("Capp: %s Error writing to parent socket. Value: %d Size: %d",
	      capp->params->name, ret, sizeof(ret));
  }

  exit(ret); /* Child exit */
}

/*
 * create_capp --
 *
 */
static int create_capp(CApp *capp) {

  if (capp == NULL) return FALSE;

  return create_space(AREA_TYPE_CAPP,
		      capp->runtime->sockets, capp->config->nameSpaces,
		      capp->params->name, &capp->runtime->pid,
		      capp_init_clone, (void *)capp);
}
