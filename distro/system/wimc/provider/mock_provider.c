/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* Mocking cloud-based service provider for testing only. */

#include <stdio.h>
#include <stdlib.h>
#include <sqlite3.h>
#include <string.h>
#include "ulfius.h"

#define CONTAINERS "/containers"

#define TRUE 1
#define FALSE 0
#define MAX_ENT 256

typedef struct _u_request req_t;
typedef struct _u_response resp_t;

/* DB: name, tag, method, CB-URL
 */

typedef struct { 
  char *name;   /* */
  char *tag;    /* */
  int  type;    /* */
  char *method; /* Mechanisim supported by service at the url. */
  char *url;    /* callback URL for the agent. */
} dbEntry;

typedef struct {
  sqlite3 *dbPtr;
  int     numEnt;
  dbEntry ent[MAX_ENT];
} DB;

/* 
 * table_exists -- Check if default table exists in the db. 
 *
 */

int table_exists(sqlite3 *db) {

  int ret;
  char *sql=NULL, *errMsg=NULL;
  
  sql = "SELECT 1 FROM sqlite_master where type='table' and name='Containers'";

  ret = sqlite3_exec(db, sql, NULL, NULL, &errMsg);

  if (ret != SQLITE_OK) {
    fprintf(stdout, "Table does not exit. Query returned error: %s\n", errMsg);
    sqlite3_free(errMsg);
    return FALSE;
  } else {
    return TRUE;
  }
}

/*
 * open_db -- Open db. A new db is created and initialized with "containers"
 *            table if it doesn't exist.
 */

static sqlite3 *open_db(char *dbFile) {

  sqlite3 *db=NULL;
  int ret;
  
  ret = sqlite3_open(dbFile, &db);
  
  if (ret) {
    fprintf(stderr, "Error opening the dbFile: %s\n", dbFile);
    exit(1);
  } else {
    fprintf(stdout, "db opened: %s\n", dbFile);
  }

  /* sanity check. Try to query the default table. */
  if (!table_exists(db)) {
    sqlite3_close(db);
    return NULL;
  } else {
    fprintf(stdout, "Found right table in db: %s\n", dbFile);
  }
  
  return db;
}

/* Callback function for the web application on /validation url call */
int callback_get_containers(req_t *req, resp_t * resp, void * user_data) {

  DB *db = (DB *)user_data;
  
  return U_CALLBACK_CONTINUE;
}

int callback_default_ok(req_t *req, resp_t * resp, void * user_data) {

  DB *db = (DB *)user_data;

  ulfius_set_string_body_response(resp, 200, "OK\n");
  return U_CALLBACK_CONTINUE;
}

static int read_entries(void *arg, int argc, char **argv, char **colName) {
  
  int i, ent;
  DB *db = (DB *)arg;
  
  if (db == NULL) {
    fprintf(stderr, "Memory failure\n");
    exit(1);
  }

  ent = db->numEnt;
  
  for(i=0; i<argc; i++){
    
    if (strcmp(colName[i], "Name") == 0) {
      db->ent[ent].name = strdup(argv[i]);
    } else if (strcmp(colName[i], "Tag") == 0 ) {
      db->ent[ent].tag = strdup(argv[i]);
    } else if (strcmp(colName[i], "Type") == 0) {
      db->ent[ent].type = atoi(argv[i]);
    } else if (strcmp(colName[i], "Method") == 0) {
      db->ent[ent].method = strdup(argv[i]);
    } else if (strcmp(colName[i], "URL") == 0) {
      db->ent[ent].url = strdup(argv[i]);
    }
  }
  
  if (db->ent[db->numEnt].name)
    db->numEnt++;
  
  return 0;
}

static int read_all_db_entries(DB *db) {
  
  int val=FALSE;
  char buf[128];
  char *err=NULL;
 
  /* sanity checks. */
  if (db == NULL || db->dbPtr == NULL) {
    return FALSE;
  }
 
  sprintf(buf, "SELECT * FROM Containers;");

  val = sqlite3_exec(db->dbPtr, buf, read_entries, db, &err);

  if (val != SQLITE_OK) {
    fprintf(stderr, "SQL read error, query failure: %s\n", err);
    sqlite3_free(err);
  } else {
    fprintf(stdout, " Query: %s\n Response ok", buf);
    val = TRUE;
  }
  
  return val;
}

static void create_db_entries_endpoint(DB *db, struct _u_instance *inst) {

  int i;

  /* Register all URL to default ok */
  for (i=0; i<db->numEnt; i++) {
    ulfius_add_endpoint_by_val(inst, "GET", db->ent[i].url, NULL, 0,
			       &callback_default_ok, db);
  }
}

static void free_db_entries(DB *db) {

  int i;

  for (i=0; i<db->numEnt; i++) {
    free(db->ent[i].name);
    free(db->ent[i].tag);
    free(db->ent[i].method);
    free(db->ent[i].url);
  }

  free(db);
}


int main(int argc, char **argv) {

  int port;
  struct _u_instance inst;
  DB *db;
  
  if (argc<2) {
    fprintf(stderr, "USAGE: %s port dbFile\n", argv[0]);
    return 0;
  }

  db = (DB *)calloc(sizeof(DB), 1);
  if (!db) {
    fprintf(stderr, "Memory allocation issue of size: %ld\n", sizeof(DB));
    exit(1);
  }
  
  port = atoi(argv[1]);
  db->dbPtr = open_db(argv[2]);

  if (!db->dbPtr) {
    fprintf(stderr, "Error opening db file: %s\n", argv[2]);
    return 1;
  }

  if (ulfius_init_instance(&inst, port, NULL, NULL) != U_OK) {
    fprintf(stderr, "Error ulfius_init_instance, abort\n");
    return 1;
  }
  
  /* Endpoint list declaration */
  ulfius_add_endpoint_by_val(&inst, "GET", CONTAINERS, NULL, 0,
			     &callback_get_containers, db);

  /* Read all db entries and process them. */
  if (read_all_db_entries(db)) {
    create_db_entries_endpoint(db, &inst);
  }

  /* Start the framework */
  if (ulfius_start_framework(&inst) == U_OK) {
    fprintf(stdout, "Famework start on port %d\n", inst.port);
    getchar();
  }
  else {
    fprintf(stderr, "Error starting framework\n");
  }
  
  fprintf(stdout, "End framework\n");

  ulfius_stop_framework(&inst);
  ulfius_clean_instance(&inst);

  free_db_entries(db);
  
  return 0;
}
