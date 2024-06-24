/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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

int db_insert_entry(sqlite3 *db, char *name, char *tag, char *status) {

  int val=FALSE;
  char *buff, *err=NULL;

  buff = (char *)calloc(1, 2048);
  
  if (db == NULL || name == NULL || tag == NULL || status == NULL ||
      buff == NULL) {
    goto failure;
  }

  sprintf(buff, "INSERT INTO Containers VALUES('%s', '%s', 'null', '%s', 'null');",
          name, tag, status);

  val = sqlite3_exec(db, buff, 0, 0, &err);

  if (val != SQLITE_OK) {
    log_error("SQL error: insert failure: %s", err);
    sqlite3_free(err);
    val = FALSE;
    goto failure;
  } else {
    log_debug("Succesfully insert in db. Name: %s Tag: %s Status: %s",
	      name, tag, status);
  }

  val = TRUE;

 failure:
  free(buff);
  return val;
}

static int parse_status_response(void *arg, int argc, char **argv, char **colName) {

    int i;
    char **status = (char **)arg;

    for(i=0; i<argc; i++){
        if (strcmp(colName[i], "Status") == 0) {
            if (argv[i] != NULL) {
                *status = strdup(argv[i]);
            return 0;
            }
        }
    }

    return 0;
}

int db_read_status(sqlite3 *db, char *name, char *tag, char **status) {
  
    int val = FALSE;
    char buf[1024] = {0};
    char *err = NULL;

    if (db == NULL || name == NULL || tag == NULL || status == NULL) {
        return FALSE;
    }

    sprintf(buf, "SELECT Status FROM Containers WHERE Name='%s' AND Tag='%s';",
            name, tag);

    val = sqlite3_exec(db, buf, parse_status_response, status, &err);

    if (val != SQLITE_OK) {
        log_error("SQL error: query failure: %s", err);
        sqlite3_free(err);
        return FALSE;
    } else {
        if (*status == NULL) {
            return FALSE;
        }
    }

    log_debug("db query. Name: %s Tag: %s Status: %s", name, tag, *status);

    return TRUE;
}

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
