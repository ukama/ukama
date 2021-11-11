/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Function related to cApps.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <uuid/uuid.h>

#include "lxce_config.h"
#include "wimc.h"
#include "capp.h"
#include "log.h"
#include "manifest.h"
#include "csthreads.h"

/*
 * capp_init_params --
 *
 */
static int capp_init_params(CApp *capp, char *name, char *tag, char *path,
			    char *space) {

  capp->params = (CAppParams *)malloc(sizeof(CAppParams));
  if (capp->params == NULL) {
    log_error("Memory allocation error of size: %d", sizeof(CAppParams));
    return FALSE;
  }

  capp->params->name  = strdup(name);
  capp->params->tag   = strdup(tag);
  capp->params->path  = strdup(path);
  capp->params->space = strdup(space);

  uuid_clear(capp->params->uuid);

  return TRUE;
}

/*
 * capp_init_state --
 *
 */
static int capp_init_state(CApp *capp) {

  capp->state = (CAppState *)malloc(sizeof(CAppState));
  if (capp->state == NULL) {
    log_error("Memory allocation error of size: %d", sizeof(CAppState));
    return FALSE;
  }

  capp->state->state       = CAPP_STATE_INVALID;
  capp->state->exit_status = CAPP_STATE_INVALID;

  return TRUE;
}

/*
 * capp_init_policy --
 *
 */
static int capp_init_policy(CApp *capp, int flag) {

  capp->policy = (CAppPolicy *)malloc(sizeof(CAppPolicy));
  if (capp->policy == NULL) {
    log_error("Memory allocation error of size: %d", sizeof(CAppPolicy));
    return FALSE;
  }

  capp->policy->restart = flag;

  return TRUE;
}

/*
 * capp_free --
 *
 */
void clear_capp(CApp *capp) {

  if (!capp) return;

  if (capp->params) {
    free(capp->params->name);
    free(capp->params->tag);
    if (capp->params->path) free(capp->params->path);

    free(capp->params);
  }

  if (capp->state)
    free(capp->state);

  if (capp->policy)
    free(capp->policy);

  free(capp);
  capp=NULL;
}

/*
 * capp_init -- initialize cApp. If path is NULL will fetch from wimc
 *
 */
CApp *capp_init(Config *config, char *name, char *tag, char *path, char *space,
		int restart) {

  CApp *capp=NULL;

  if (!config || !name || !tag || !space)
    return FALSE;

#if 0
  if (path == NULL) {
    path = strdup(DEF_PATH);
  }
#endif

  if (path == NULL){
    if (get_capp_path(config, name, tag, path)==FALSE) {
      log_error("Error getting rootfs path from wimc. Name: %s Tag: %s",
		name, tag);
      return FALSE;
    } else {
      log_debug("Name: %s Tag: %s Path: %s", name, tag, path);
    }
  }

  /* We have valid name, tag and path. */
  capp = (CApp *)calloc(1, sizeof(CApp));
  if (!capp) {
    log_error("Memory allocation error of size: %d", sizeof(CApp));
    goto failure;
  }

  capp_init_params(capp, name, tag, path, space);
  capp_init_state(capp);
  capp_init_policy(capp, restart);

  capp->space = NULL;
  free(path);

  return capp;

 failure:
  if (path) free(path);
  clear_capp(capp);

  return NULL;
}

/*
 * capps_init -- for all the valid apps in manifest initialize their respective
 *               capps. Initially everything will be in 'pend'
 *
 */
int capps_init(CApps **capps, Config *config, Manifest *manifest) {

  CApp *app=NULL;
  CAppList *pend=NULL;
  ArrayElem *ptr=NULL;

  if (manifest == NULL || *capps) return FALSE;

  if (manifest->arrayElem == NULL) return FALSE;

  *capps = (CApps *)calloc(1, sizeof(CApps));
  if (*capps == NULL) {
    log_error("Memory allocation error of size: %d", sizeof(CApps));
    return FALSE;
  }

  for (ptr = manifest->arrayElem; ptr; ptr=ptr->next) {

    if (!ptr->name || !ptr->tag || !ptr->contained) {
      log_error("Invalid manifest entry. Ignoring");
      continue;
    }

    app = capp_init(config, ptr->name, ptr->tag, NULL, ptr->restart,
		    ptr->contained);
    if (app==NULL) {
      log_error("Error initializing the cApp. Name: %s Tag: %s Ignoring.",
		ptr->name, ptr->tag);
      continue;
    }

    /* Find space pointer */

    /* Add the capp to pend list */
    add_to_apps(*capps, app, PEND_LIST, NULL);
  }

  return TRUE;
}

/*
 * capp_start --
 *
 */
void capps_start(CApps *capps, CSpaceThread *threads) {

  CAppList *ptr;
  CApp *capp;
  CSpace *cspace;
  char *name, *tag, *path;
  ThreadShMem *shMem;

  /*
   * For each app in the pend list:
   *   1. find the thread handling the capp space.
   *   2. create capp_packet
   *   3. Send packet to space
   *   4. Get response (e.g., UUID or error)
   *   5. Move it to create
   *   Repeat
   */

  if (!capps) return FALSE;

  for(ptr=capps->pend; ptr; ptr=ptr->next) {

    capp = ptr->capp;

    if (capp==NULL) continue;

    name = capp->params->name;
    tag  = capp->params->tag;
    path = capp->params->path;

    shMem = find_matching_thread_shmem(capp->params->space);

    if (shMem == NULL) {
      log_error("No matching cspace for the capp. Name: %s tag: %s space: %s",
		capp->params->name, capp->params->tag, capp->params->space);
      continue;
    }

    /* create the tx packet and trigger conditional variable for the
     * target shared memory
     */
    create_capp_tx_packet(capp, &shMem->txList, CAPP_TYPE_REQ_CREATE);
  }

  /* broadcast item is available in the queue. */
  pthread_cond_broadcast(&(shMem->hasRX));

  /* unlock */
  pthread_mutex_unlock(&(shMem->rxMutex));

}

/*
 * add_to_apps --
 *
 */
void add_to_apps(CApps *capps, CApp *capp, int to, int from) {

  if (to == PEND_LIST && from == NULL) { /* New addition. */
    add_to_list(&(capps->pend), capp);
  } 
}

/*
 * add_to_list --
 *
 */
int add_to_list(CAppList **list, CApp *app) {

  CAppList *ptr;

  if (app == NULL) return FALSE;

  if (*list == NULL) { /* First entry */
    *list = (CAppList *)calloc(1, sizeof(CAppList));
    if (*list == NULL) return FALSE;
    ptr = *list;
  } else {
    (*list)->next = (CAppList *)calloc(1, sizeof(CAppList));
    if ((*list)->next == NULL) return FALSE;
    ptr = (*list)->next;
  }

  ptr->capp = app;

  return TRUE;
}

/*
 * clear_capps --
 *
 */
void clear_capps(CApps *capps, int flag) {

  CAppList *ptr=NULL, *tmp=NULL, *head=NULL;

  switch(flag) {
  case PEND_LIST:
    head = capps->pend;
    break;

  case CREATE_LIST:
    head = capps->create;
    break;

  case RUN_LIST:
    head = capps->run;
    break;

  case TERM_LIST:
    head = capps->term;
    break;

  case ERROR_LIST:
    head = capps->error;
    break;

  default:
    return;
  }

  ptr = head;
  while (ptr) {
    clear_capp(ptr->capp);
    tmp = ptr->next;
    free(ptr);
    ptr = tmp;
  }
}
