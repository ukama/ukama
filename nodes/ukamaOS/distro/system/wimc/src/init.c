/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
  char *sql, *errMsg;

  sql = "DROP TABLE IF EXISTS Containers; CREATE TABLE Containers (Name TEXT, Tag TEXT, Path TEXT, Status TEXT, Flags INIT); INSERT INTO Containers VALUES('Null', 'latest', '/dev/null', 'Available', 0);"; 
  
  ret = sqlite3_exec(db, sql, NULL, NULL, &errMsg);
    
  if (ret != SQLITE_OK) {
        
    fprintf(stderr, "Error initializing the db: %s\n", errMsg);
    sqlite3_free(errMsg);

    return FALSE;
  } 
    
  return TRUE;
}

/*
 * open_db -- Open db. A new db is created and initialized with "containers"
 *            table if it doesn't exist.
 */

sqlite3 *open_db(char *dbFile, int flag) {

  sqlite3 *db;

  db = create_db(dbFile);
  if (db == NULL) {
    return NULL;
  }

  if (flag == WIMC_FLAG_CREATE_DB) {
    
    /* initialize db with new table. */
    if (!init_db(db)) {
      /* close the db and delete file. */
      sqlite3_close(db);
      delete_file(dbFile);
      return NULL;
    }
  } else {
    /* sanity check. Try to query the default table. */
    if (!table_exists(db)) {
      sqlite3_close(db);
      return NULL;
    }      
  }
  
  return db;
}
