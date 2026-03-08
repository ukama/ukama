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

static int bind_text_or_null(sqlite3_stmt *stmt, int idx, char *val) {

    if (val == NULL) {
        return sqlite3_bind_null(stmt, idx);
    }

    return sqlite3_bind_text(stmt, idx, val, -1, SQLITE_TRANSIENT);
}

static int db_exec_simple(sqlite3 *db, const char *sql) {

    int ret;
    char *err = NULL;

    ret = sqlite3_exec(db, sql, NULL, NULL, &err);
    if (ret != SQLITE_OK) {
        log_error("SQL error: %s", err ? err : "unknown");
        if (err) {
            sqlite3_free(err);
        }
        return FALSE;
    }

    return TRUE;
}

static int db_upsert_entry(sqlite3 *db, char *name, char *tag,
                           char *path, char *status) {

    int ret;
    sqlite3_stmt *stmt = NULL;
    const char *sql =
        "INSERT INTO Containers(Name, Tag, Path, Status, Flags) "
        "VALUES(?, ?, COALESCE(?, 'null'), ?, 'null') "
        "ON CONFLICT(Name, Tag) DO UPDATE SET "
        "Path=COALESCE(excluded.Path, Containers.Path), "
        "Status=excluded.Status;";

    if (db == NULL || name == NULL || tag == NULL || status == NULL) {
        return FALSE;
    }

    ret = sqlite3_prepare_v2(db, sql, -1, &stmt, NULL);
    if (ret != SQLITE_OK) {
        log_error("SQL prepare error: %s", sqlite3_errmsg(db));
        return FALSE;
    }

    if (bind_text_or_null(stmt, 1, name) != SQLITE_OK ||
        bind_text_or_null(stmt, 2, tag) != SQLITE_OK ||
        bind_text_or_null(stmt, 3, path) != SQLITE_OK ||
        bind_text_or_null(stmt, 4, status) != SQLITE_OK) {
        log_error("SQL bind error: %s", sqlite3_errmsg(db));
        sqlite3_finalize(stmt);
        return FALSE;
    }

    ret = sqlite3_step(stmt);
    sqlite3_finalize(stmt);
    if (ret != SQLITE_DONE) {
        log_error("SQL step error: %s", sqlite3_errmsg(db));
        return FALSE;
    }

    log_debug("DB upsert. Name: %s Tag: %s Path: %s Status: %s",
              name, tag, path ? path : "(null)", status);

    return TRUE;
}

int db_insert_entry(sqlite3 *db, char *name, char *tag, char *status) {

    return db_upsert_entry(db, name, tag, NULL, status);
}

int db_update_status(sqlite3 *db, char *name, char *tag, char *status) {

    return db_upsert_entry(db, name, tag, NULL, status);
}

int db_update_path_status(sqlite3 *db, char *name, char *tag,
                          char *path, char *status) {

    return db_upsert_entry(db, name, tag, path, status);
}

int db_mark_old_downloads_failed(sqlite3 *db) {

    return db_exec_simple(db,
                          "UPDATE Containers SET Status='failed' "
                          "WHERE Status='download';");
}

int db_read_status(sqlite3 *db, char *name, char *tag, char **status) {

    int val = FALSE;
    sqlite3_stmt *stmt = NULL;
    const unsigned char *text = NULL;
    const char *sql = "SELECT Status FROM Containers WHERE Name=? AND Tag=?;";

    if (db == NULL || name == NULL || tag == NULL || status == NULL) {
        return FALSE;
    }

    *status = NULL;

    val = sqlite3_prepare_v2(db, sql, -1, &stmt, NULL);
    if (val != SQLITE_OK) {
        log_error("SQL prepare failure: %s", sqlite3_errmsg(db));
        return FALSE;
    }

    sqlite3_bind_text(stmt, 1, name, -1, SQLITE_TRANSIENT);
    sqlite3_bind_text(stmt, 2, tag, -1, SQLITE_TRANSIENT);

    val = sqlite3_step(stmt);
    if (val == SQLITE_ROW) {
        text = sqlite3_column_text(stmt, 0);
        if (text != NULL) {
            *status = strdup((const char *)text);
        }
        sqlite3_finalize(stmt);
        if (*status == NULL) {
            return FALSE;
        }
        log_debug("db query. Name: %s Tag: %s Status: %s", name, tag, *status);
        return TRUE;
    }

    sqlite3_finalize(stmt);
    return FALSE;
}

void update_local_db(sqlite3 *db, char *name, char *tag, char *path) {

    FILE *fp;
    struct stat sb;
    char fileName[WIMC_MAX_PATH_LEN] = {0};

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

    snprintf(fileName, sizeof(fileName), "%s/index.json", path);

    fp = fopen(fileName, "r");
    if (fp == NULL) {
        log_error("Failed to read index.json at: %s Error: %s", fileName,
                  strerror(errno));
        return;
    }
    fclose(fp);

    /* All checks passed. Add into the db for future generations. */
    db_update_path_status(db, name, tag, path, "available");
}
