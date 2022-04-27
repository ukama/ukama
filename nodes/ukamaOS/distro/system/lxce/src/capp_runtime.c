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
#include "utils.h"
#include "capp_config.h"
#include "capp_runtime.h"

static int exec_capp(CApp *capp);
static int capp_init_clone(void *arg);
static int create_capp(CApp *capp);
static int token_count(char *str);
static int create_list(char *exec, char *arg, char ***argv);

/*
 * token_count --
 *
 */
static int token_count(char *str) {

  int count=0, i;

  if (!str) return count;
  if (!strlen(str)) return count;

  for (i=0; i<strlen(str); i++) {
    if (str[i] == ' ') {
      count++;
    }
  }

  return count++;
}

/*
 * create_list --
 *
 */
static int create_list(char *exec, char *arg, char ***argv) {

  int count=0, i=0;
  char *str, *token;

  if (arg == NULL) return 0;

  count = token_count(arg);
  if (exec) {
    count += 2; /* for binary and NULL, see exec(3) */
    log_debug("Creating args list for capp: %s", exec);
  } else {
    count += 1; /* for env variables. */
    log_debug("Creating env list for capp");
  }

  *argv = (char **)calloc(count, sizeof(char *));
  if (argv == NULL) {
    log_error("Error allocating memory of size: %d", count * sizeof(char *));
    return FALSE;
  }

  str = strdup(arg);

  if (exec) { /* only need for runtime arguments. */
    (*argv)[0] = strdup(exec);
    i++;
  }

  token = strtok(str, " ");
  while (token != NULL) {
    (*argv)[i] = strdup(token);
    token = strtok(NULL, " ");
    i++;
  }
  (*argv)[i] = (char *) NULL;

  return count;
}

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

  int i=0, argc=0, envc=0;
  CAppProc *process=NULL;
  char **argv=NULL, **env=NULL;

  if (!capp || !capp->config) return FALSE;

  process = capp->config->process;

  if (!process->exec) return FALSE;

  /* Runtime argument list */
  if ((argc = create_list(process->exec, process->argv, &argv)) == FALSE) {
    log_error("Error creating argument list for capp execution: %s",
	      process->exec);
    return FALSE;
  }

  /* Environment varaibles list */
  if (process->env) {
    if ((envc = create_list(NULL, process->env, &env)) == FALSE) {
      log_error("Error creating env list for capp execution: %s",
	      process->exec);
      return FALSE;
    }
  }

  log_debug("Executing capp: binary: %s argc: %d env: %d", process->exec,
	    argc, envc);
  for (i=0; i<argc; i++) {
    if (argv[i]) {
      log_debug("\t argc-%d: %s", i, argv[i]);
    }
  }

  for (i=0; i<envc; i++) {
    if (env[i]) {
      log_debug("\t envc-%d: %s", i, env[i]);
    }
  }

  execvpe(process->exec, argv, env);

  /* executed only if there was an error with exec */
  for (i=0; i<argc; i++) {
    if(argv[i]) free(argv[i]);
  }
  if (argv) free(argv);

  for (i=0; i<envc; i++) {
    if(env[i]) free(env[i]);
  }
  if (env) free(env);

  log_error("capp execution failed. Code: %d Error: %s", errno,
	    strerror(errno));
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

#if 0
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
#endif

  log_debug("Exec'ing into capp.");

  /* execv into program */
  ret = exec_capp(capp);

  /* An error has occured. Inform the parent process over socket and exit. */
  if (write(capp->runtime->sockets[CHILD_SOCKET], &ret,
	    sizeof(ret)) != sizeof(ret)) {
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

  log_debug("Creating the capp: %s", capp->params->name);

  return create_space(AREA_TYPE_CAPP,
		      capp->runtime->sockets, capp->config->nameSpaces,
		      capp->params->name, &capp->runtime->pid,
		      capp_init_clone, (void *)capp);
}
