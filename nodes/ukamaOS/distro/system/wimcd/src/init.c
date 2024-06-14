/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

/*
 * Open, create and initialize WIMC db with 'Containers' table
 */

#include <stdio.h>
#include <stdlib.h>
#include <sqlite3.h>

#include "log.h"

#define WIMC_FLAG_CREATE_DB 1
#define FALSE 0
#define TRUE  1

/* 
 * delete_file -- Delete db file. 
 *
 */

void delete_file(char *fileName) {

  if (remove(fileName)) {
    log_debug("File removed: %s", fileName);
  } else {
    log_error("Error removing file: %s", fileName);
  }
}

/* 
 * table_exists -- Check if default table exists in the db. 
 *
 */

int table_exists(sqlite3 *db) {

  int ret;
  char *sql=NULL, *errMsg=NULL;
  
  sql = "SELECT 1 FROM sqlite_master where type='table' and name='Containers'";

  ret = sqlite3_exec(db, sql, NULL, NULL, &errMsg);

  if (ret) {
    return TRUE;
  } else {
    log_debug("Table exists query returned error: %s", errMsg);
    sqlite3_free(errMsg);
    return FALSE;
  }
}
  
/* 
 * create_db -- 
 *
 */

sqlite3 *create_db(char *dbFile) {
  
  int ret;
  sqlite3 *db;

  ret = sqlite3_open(dbFile, &db);

  if (ret != SQLITE_OK) {
    log_error("Can not open db at: %s Error: %s", dbFile, sqlite3_errmsg(db));
    sqlite3_close(db);
    return NULL;
  } else {
    log_debug("db opened: %s", dbFile);
  }
  
  return db;
}

/*
 * init_db -- create table in the db. 
 *            Table is 'Containers' and have following entries:
 *            Name   - name of the container, e.g, neutron
 *            Tag    - version e.g. latest or 8.7.9
 *            Path   - unbundle location of this container. This is where the 
 *                     rootfs and config.json reside.
 *            Status - Download, Unbundle, Available, Obsolete
 *            Flags  - any flags associated with container (internal)
 */

int init_db(sqlite3 *db) {

  int ret;
  char *errMsg = NULL;
  const char *sql = "CREATE TABLE IF NOT EXISTS Containers ("
      "Name TEXT NOT NULL,"
      "Tag TEXT NOT NULL,"
      "Path TEXT NOT NULL,"
      "Status TEXT NOT NULL,"
      "Flags TEXT);";

  ret = sqlite3_exec(db, sql, NULL, NULL, &errMsg);
    
  if (ret != SQLITE_OK) {
      log_error("Error initializing the db: %s", errMsg);
      sqlite3_free(errMsg);

      return FALSE;
  }
  return TRUE;
}

void open_db(sqlite3 **db, char *dbFile, int flag) {

    *db = create_db(dbFile);
    if (*db == NULL) {
        log_error("Failed to created db: %s", dbFile);
        return;
    }

    if (flag == WIMC_FLAG_CREATE_DB) {
        if (!init_db(*db)) {
            log_error("Failed to initialize db: %s", dbFile);
            sqlite3_close(*db);
            delete_file(dbFile);
            *db = NULL;
            return;
        }
    } else {
        if (!table_exists(*db)) {
            log_error("Default table does not exist in db: %s", dbFile);
            sqlite3_close(*db);
            *db = NULL;
            return;
        }
    }
}
