/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * capp.h
 */

#ifndef LXCE_CAPP_H
#define LXCE_CAPP_H

#include <uuid/uuid.h>

#include "lxce_config.h"
#include "capp_config.h"
#include "capp_runtime.h"
#include "manifest.h"

/* For capp state */
#define CAPP_STATE_PENDING 0x01
#define CAPP_STATE_CREATE  0x02
#define CAPP_STATE_RUN     0x03
#define CAPP_STATE_TERM    0x04
#define CAPP_STATE_INVALID 0xff

/* commands */
#define CAPP_CMD_CREATE "create"
#define CAPP_CMD_RUN    "run"
#define CAPP_CMD_STATUS "status"

/* action state */
#define CAPP_CMD_STATE_WAIT  1
#define CAPP_CMD_STATE_DONE  2
#define CAPP_CMD_STATE_ERROR 3

#define CAPP_MAX_BUFFER   1024
#define CAPP_READ_ERROR   1
#define CAPP_READ_TIMEOUT 2

/* List type */
enum { START_LIST=0,
       PEND_LIST,
       CREATE_LIST, 
       RUN_LIST,
       TERM_LIST,
       ERROR_LIST,
       END_LIST
};

/* Default various capps path */
#define DEF_CAPP_PATH        "/capps/pkgs"
#define DEF_CAPP_UNPACK_PATH "/capps/pkgs/unpack"
#define DEF_CONFIG "config.json"

typedef struct capp_params_t {

  char   *name;  /* capp name */
  char   *tag;   /* capp tag */
  char   *path;  /* path to rootfs */
  char   *space; /*cspace name */
  uuid_t uuid;   /* UUID per its cspace */
} CAppParams;

typedef struct capp_state_ {

  int state;       /* State of capp. CAPP_STATE_XXX */
  int exit_status; /* Exit status of the capp if terminated. */
} CAppState;

typedef struct capp_policy_ {

  int restart;     /* restart of capp terminates? */
} CAppPolicy;

typedef struct capp_t_ {

  CAppParams *params;
  CAppState  *state;  /* capp state */
  CAppPolicy *policy; /* capp assocated policy */
  void       *space;  /* space the capp belongs to */
  CAppConfig *config; /* config.json */

  CAppRuntime *runtime; /* runtime stuff */
} CApp;

typedef struct capp_list_ {

  CApp *capp;
  
  struct capp_list_ *next;
} CAppList;

typedef struct capp_t {

  CAppList *pend;    /* yet to be created on its space */
  CAppList *create;  /* cspace is currently creating this capp */
  CAppList *run;     /* capp is running within its cspace */
  CAppList *term;    /* capp is terminated (stop, term, killed) */
  CAppList *error;   /* capp has an error */
} CApps;

int capps_init(CApps **capps, Config *config, Manifest *manifest, void *space);
void capps_start(CApps *capps);
void clear_capp(CApp *capp);
void clear_capps(CApps *capps, int flag);
void add_to_apps(CApps *capps, CApp *capp, int to, int from);

#endif /* LXCE_CAPP_H */

