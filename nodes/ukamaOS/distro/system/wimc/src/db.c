/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * database helper functions.
 */

#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <sqlite3.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>

#include "log.h"
#include "wimc.h"

#define TRUE 1
#define FALSE 0

/*
 * insert_entry -- 
 *
 */
static int db_insert_entry(sqlite3 *db, char *name, char *tag, char *path) {

  int val=FALSE;
  char *buff, *err=NULL;

  buff = (char *)calloc(1, 2048);
  
  /* sanity checks. */
  if (db == NULL || name == NULL || tag == NULL || path == NULL ||
      buff == NULL) {
    goto failure;
  }

  sprintf(buff, "INSERT INTO Containers VALUES('%s', '%s', '%s', 'null', 'null');", name, tag, path);

  val = sqlite3_exec(db, buff, 0, 0, &err);

  if (val != SQLITE_OK) {
    log_error("SQL error: insert failure: %s", err);
    sqlite3_free(err);
    val = FALSE;
    goto failure;
  } else {
    log_debug("Succesfully insert in db. Name: %s Tag: %s Path: %s",
	      name, tag, path);
  }

  val = TRUE;

 failure:
  free(buff);
  return val;
}

/* 
 * parse_response -- query callback.
 *
 */

static int parse_response(void **arg, int argc, char **argv, char **colName) {
  
  int i;
  char **path = (char **)arg;
  
  for(i=0; i<argc; i++){
    if (strcmp(colName[i], "Path") == 0) {
      strcpy(*path, argv[1]);
      return 0;
    }
  }
  
  return 0;
}

/*
 * read_path -- for given container name and tag, return the path.
 *
 */
int db_read_path(sqlite3 *db, char *name, char *tag, char **path) {
  
  int val=FALSE;
  char buf[1024];
  char *err=NULL;
  
  /* sanity checks. */
  if (db == NULL || name == NULL || tag == NULL || path == NULL) {
    goto failure;
  }
  
  sprintf(buf, "SELECT * FROM Containers WHERE (Name='%s' AND Tag='%s');",
	  name, tag);

  val = sqlite3_exec(db, buf, parse_response, path, &err);

  if (val != SQLITE_OK) {
    log_error("SQL error: query failure: %s", err);
    sqlite3_free(err);
    val = FALSE;
    goto failure;
  } else {
    if (*path==NULL) {
      goto failure;
    }

    log_debug("db query. Name: %s Tag: %s Path: %s",
	      name, tag, *path);
  }
  val = TRUE;

 failure:
  return val;
}

/*
 * update_local_db -- update entries in local db.
 */

void update_local_db(sqlite3 *db, char *name, char *tag, char *path) {

  FILE *fp;
  struct stat sb;
  char fileName[WIMC_MAX_PATH_LEN]={0};

  /* sanity checks. */
  if (db == NULL || name == NULL || tag == NULL || path == NULL) {
    return;
  }

  /* Check if its a valid path and json exist. */
  if (stat(path, &sb) == -1) {
    log_error("Invalid path for db entry: %s", path);
    return;
  }

  /* Check to see if it was file. */
  if (!S_ISDIR(sb.st_mode)) {
    log_error("Not valid directory for db entry: %s", path);
    return;
  }

  sprintf(fileName, "%s/index.json", path);

  fp = fopen(fileName, "a");

  if (fp == NULL) {
    log_error("Failed to read index.json at: %s Error: %s", path,
	      strerror(errno));
    return;
  }

  /* All check passed. Add into the db for future generations. */
  db_insert_entry(db, name, tag, path);
}
