/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <errno.h>
#include <limits.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

#include <sqlite3.h>

#include "db.h"
#include "log.h"
#include "package_cache.h"
#include "wimc.h"

#include "usys_log.h"

#define TRUE 1
#define FALSE 0

static int make_parent_dirs(const char *filePath) {

    char tmp[PATH_MAX];
    char *p;

    if (filePath == NULL || *filePath == '\0') {
        return -1;
    }

    memset(tmp, 0, sizeof(tmp));

    if (strlen(filePath) >= sizeof(tmp)) {
        usys_log_error("DB path too long: %s", filePath);
        return -1;
    }

    strncpy(tmp, filePath, sizeof(tmp) - 1);

    p = strrchr(tmp, '/');
    if (p == NULL) {
        return 0;
    }

    *p = '\0';
    if (tmp[0] == '\0') {
        return 0;
    }

    for (p = tmp + 1; *p != '\0'; p++) {
        if (*p == '/') {
            *p = '\0';
            if (mkdir(tmp, 0755) != 0 && errno != EEXIST) {
                usys_log_error("mkdir failed for %s: %s",
                               tmp, strerror(errno));
                return -1;
            }
            *p = '/';
        }
    }

    if (mkdir(tmp, 0755) != 0 && errno != EEXIST) {
        usys_log_error("mkdir failed for %s: %s", tmp, strerror(errno));
        return -1;
    }

    return 0;
}

static int db_exec(sqlite3 *db, const char *sql) {

    int rc;
    char *errMsg = NULL;

    rc = sqlite3_exec(db, sql, NULL, NULL, &errMsg);
    if (rc != SQLITE_OK) {
        usys_log_error("SQL error: %s", errMsg ? errMsg : "unknown");
        if (errMsg != NULL) {
            sqlite3_free(errMsg);
        }
        return -1;
    }

    return 0;
}

static int db_column_exists(sqlite3 *db, const char *table,
                            const char *column) {

    sqlite3_stmt *stmt = NULL;
    char sql[256];
    int rc;
    int found = 0;

    if (snprintf(sql, sizeof(sql), "PRAGMA table_info(%s);",
                 table) >= (int)sizeof(sql)) {
        return 0;
    }

    rc = sqlite3_prepare_v2(db, sql, -1, &stmt, NULL);
    if (rc != SQLITE_OK) {
        return 0;
    }

    while (sqlite3_step(stmt) == SQLITE_ROW) {
        const unsigned char *name;

        name = sqlite3_column_text(stmt, 1);
        if (name != NULL && strcmp((const char *)name, column) == 0) {
            found = 1;
            break;
        }
    }

    sqlite3_finalize(stmt);
    return found;
}

static int db_add_column_if_missing(sqlite3 *db, const char *column,
                                    const char *type) {

    char sql[256];

    if (db_column_exists(db, "Containers", column)) {
        return 0;
    }

    if (snprintf(sql, sizeof(sql), "ALTER TABLE Containers ADD COLUMN %s %s;",
                 column, type) >= (int)sizeof(sql)) {
        return -1;
    }

    return db_exec(db, sql);
}

static int db_init_schema(sqlite3 *db) {

    const char *sql =
        "CREATE TABLE IF NOT EXISTS Containers ("
        "  Name          TEXT NOT NULL,"
        "  Tag           TEXT NOT NULL,"
        "  Path          TEXT,"
        "  Status        TEXT NOT NULL,"
        "  Flags         TEXT,"
        "  ActualVersion TEXT,"
        "  Error         TEXT,"
        "  UpdatedAt     INTEGER DEFAULT (strftime('%s','now')),"
        "  UNIQUE(Name, Tag)"
        ");";

    if (db_exec(db, sql) != 0) {
        return -1;
    }

    if (db_add_column_if_missing(db, "ActualVersion", "TEXT") != 0 ||
        db_add_column_if_missing(db, "Error", "TEXT") != 0 ||
        db_add_column_if_missing(db, "UpdatedAt", "INTEGER") != 0) {
        return -1;
    }

    return 0;
}

static int bind_text_or_null(sqlite3_stmt *stmt, int idx, char *val) {

    if (val == NULL) {
        return sqlite3_bind_null(stmt, idx);
    }

    return sqlite3_bind_text(stmt, idx, val, -1, SQLITE_TRANSIENT);
}

int db_update_package_status(sqlite3 *db, char *name, char *tag,
                             char *path, char *status,
                             char *actualVersion, char *error) {

    int rc;
    sqlite3_stmt *stmt = NULL;
    const char *sql =
        "INSERT INTO Containers(Name, Tag, Path, Status, Flags, "
        "ActualVersion, Error, UpdatedAt) "
        "VALUES(?, ?, ?, ?, 'null', ?, ?, strftime('%s','now')) "
        "ON CONFLICT(Name, Tag) DO UPDATE SET "
        "Path=COALESCE(excluded.Path, Containers.Path), "
        "Status=excluded.Status, "
        "ActualVersion=COALESCE(excluded.ActualVersion, "
        "Containers.ActualVersion), "
        "Error=excluded.Error, "
        "UpdatedAt=strftime('%s','now');";

    if (db == NULL || name == NULL || tag == NULL || status == NULL) {
        return FALSE;
    }

    rc = sqlite3_prepare_v2(db, sql, -1, &stmt, NULL);
    if (rc != SQLITE_OK) {
        log_error("SQL prepare error: %s", sqlite3_errmsg(db));
        return FALSE;
    }

    if (bind_text_or_null(stmt, 1, name) != SQLITE_OK ||
        bind_text_or_null(stmt, 2, tag) != SQLITE_OK ||
        bind_text_or_null(stmt, 3, path) != SQLITE_OK ||
        bind_text_or_null(stmt, 4, status) != SQLITE_OK ||
        bind_text_or_null(stmt, 5, actualVersion) != SQLITE_OK ||
        bind_text_or_null(stmt, 6, error) != SQLITE_OK) {
        log_error("SQL bind error: %s", sqlite3_errmsg(db));
        sqlite3_finalize(stmt);
        return FALSE;
    }

    rc = sqlite3_step(stmt);
    sqlite3_finalize(stmt);
    if (rc != SQLITE_DONE) {
        log_error("SQL step error: %s", sqlite3_errmsg(db));
        return FALSE;
    }

    log_debug("DB package update. Name:%s Tag:%s Status:%s Version:%s",
              name, tag, status,
              actualVersion ? actualVersion : "(null)");

    return TRUE;
}

int db_insert_entry(sqlite3 *db, char *name, char *tag, char *status) {

    return db_update_package_status(db, name, tag, NULL, status, NULL, NULL);
}

int db_update_status(sqlite3 *db, char *name, char *tag, char *status) {

    return db_update_package_status(db, name, tag, NULL, status, NULL, NULL);
}

int db_update_path_status(sqlite3 *db, char *name, char *tag,
                          char *path, char *status) {

    return db_update_package_status(db, name, tag, path, status, NULL, NULL);
}

int db_mark_old_downloads_failed(sqlite3 *db) {

    return db_exec(db,
                   "UPDATE Containers "
                   "SET Status='failed', Error='stale download after restart', "
                   "UpdatedAt=strftime('%s','now') "
                   "WHERE Status='download' OR Status='downloading' OR "
                   "Status='queued';") == 0;
}

int db_read_status(sqlite3 *db, char *name, char *tag, char **status) {

    char *path = NULL;
    char *actualVersion = NULL;
    char *error = NULL;
    int ret;

    ret = db_read_package(db, name, tag, status, &path,
                          &actualVersion, &error);
    free(path);
    free(actualVersion);
    free(error);

    return ret;
}

int db_read_package(sqlite3 *db, char *name, char *tag, char **status,
                    char **path, char **actualVersion, char **error) {

    int rc;
    sqlite3_stmt *stmt = NULL;
    const char *sql =
        "SELECT Status, Path, ActualVersion, Error "
        "FROM Containers WHERE Name=? AND Tag=?;";

    if (db == NULL || name == NULL || tag == NULL || status == NULL ||
        path == NULL || actualVersion == NULL || error == NULL) {
        return FALSE;
    }

    *status = NULL;
    *path = NULL;
    *actualVersion = NULL;
    *error = NULL;

    rc = sqlite3_prepare_v2(db, sql, -1, &stmt, NULL);
    if (rc != SQLITE_OK) {
        log_error("SQL prepare failure: %s", sqlite3_errmsg(db));
        return FALSE;
    }

    sqlite3_bind_text(stmt, 1, name, -1, SQLITE_TRANSIENT);
    sqlite3_bind_text(stmt, 2, tag, -1, SQLITE_TRANSIENT);

    rc = sqlite3_step(stmt);
    if (rc == SQLITE_ROW) {
        const unsigned char *text;

        text = sqlite3_column_text(stmt, 0);
        if (text != NULL) *status = strdup((const char *)text);

        text = sqlite3_column_text(stmt, 1);
        if (text != NULL) *path = strdup((const char *)text);

        text = sqlite3_column_text(stmt, 2);
        if (text != NULL) *actualVersion = strdup((const char *)text);

        text = sqlite3_column_text(stmt, 3);
        if (text != NULL) *error = strdup((const char *)text);

        sqlite3_finalize(stmt);
        return *status != NULL;
    }

    sqlite3_finalize(stmt);
    return FALSE;
}

int db_count_status(sqlite3 *db, const char *status) {

    sqlite3_stmt *stmt = NULL;
    int rc;
    int count = 0;
    const char *sql = "SELECT COUNT(*) FROM Containers WHERE Status=?;";

    if (db == NULL || status == NULL) {
        return 0;
    }

    rc = sqlite3_prepare_v2(db, sql, -1, &stmt, NULL);
    if (rc != SQLITE_OK) {
        return 0;
    }

    sqlite3_bind_text(stmt, 1, status, -1, SQLITE_TRANSIENT);

    if (sqlite3_step(stmt) == SQLITE_ROW) {
        count = sqlite3_column_int(stmt, 0);
    }

    sqlite3_finalize(stmt);
    return count;
}

typedef struct _DbRow {
    char *name;
    char *tag;
    char *path;
    struct _DbRow *next;
} DbRow;

static void free_db_rows(DbRow *rows) {

    DbRow *tmp;

    while (rows != NULL) {
        tmp = rows->next;
        free(rows->name);
        free(rows->tag);
        free(rows->path);
        free(rows);
        rows = tmp;
    }
}

int db_revalidate_available(sqlite3 *db) {

    sqlite3_stmt *stmt = NULL;
    int rc;
    DbRow *rows = NULL;
    DbRow *row;
    const char *sql =
        "SELECT Name, Tag, Path FROM Containers WHERE Status='available';";

    rc = sqlite3_prepare_v2(db, sql, -1, &stmt, NULL);
    if (rc != SQLITE_OK) {
        return FALSE;
    }

    while (sqlite3_step(stmt) == SQLITE_ROW) {
        const unsigned char *name;
        const unsigned char *tag;
        const unsigned char *path;

        name = sqlite3_column_text(stmt, 0);
        tag = sqlite3_column_text(stmt, 1);
        path = sqlite3_column_text(stmt, 2);

        if (name == NULL || tag == NULL || path == NULL) {
            continue;
        }

        row = (DbRow *)calloc(1, sizeof(DbRow));
        if (row == NULL) {
            sqlite3_finalize(stmt);
            free_db_rows(rows);
            return FALSE;
        }

        row->name = strdup((const char *)name);
        row->tag = strdup((const char *)tag);
        row->path = strdup((const char *)path);
        row->next = rows;
        rows = row;
    }

    sqlite3_finalize(stmt);

    for (row = rows; row != NULL; row = row->next) {
        PackageInfo info;

        if (pkg_validate_tar(row->name, row->tag, row->path, &info)) {
            db_update_package_status(db, row->name, row->tag, row->path,
                                     WIMC_STATUS_AVAILABLE,
                                     info.actualVersion, NULL);
        } else {
            db_update_package_status(db, row->name, row->tag, row->path,
                                     WIMC_STATUS_CORRUPT,
                                     info.actualVersion[0] ?
                                     info.actualVersion : NULL,
                                     info.error[0] ? info.error :
                                     "invalid package");
        }
    }

    free_db_rows(rows);
    return TRUE;
}

void update_local_db(sqlite3 *db, char *name, char *tag, char *path) {

    PackageInfo info;

    if (db == NULL || name == NULL || tag == NULL || path == NULL) {
        return;
    }

    if (pkg_validate_tar(name, tag, path, &info)) {
        db_update_package_status(db, name, tag, path,
                                 WIMC_STATUS_AVAILABLE,
                                 info.actualVersion, NULL);
        db_update_package_status(db, name, info.actualVersion, path,
                                 WIMC_STATUS_AVAILABLE,
                                 info.actualVersion, NULL);
    } else {
        db_update_package_status(db, name, tag, path,
                                 WIMC_STATUS_CORRUPT,
                                 info.actualVersion[0] ?
                                 info.actualVersion : NULL,
                                 info.error[0] ? info.error :
                                 "invalid package");
    }
}

int db_open_or_create(const char *dbPath, sqlite3 **db) {

    int rc;

    if (dbPath == NULL || db == NULL) {
        return -1;
    }

    if (make_parent_dirs(dbPath) != 0) {
        usys_log_error("Failed to create parent dir for DB path: %s", dbPath);
        return -1;
    }

    rc = sqlite3_open(dbPath, db);
    if (rc != SQLITE_OK) {
        usys_log_error("sqlite3_open failed for %s: %s",
                       dbPath, *db ? sqlite3_errmsg(*db) : "unknown");
        if (*db != NULL) {
            sqlite3_close(*db);
            *db = NULL;
        }
        return -1;
    }

    (void)sqlite3_exec(*db, "PRAGMA journal_mode=WAL;", NULL, NULL, NULL);
    (void)sqlite3_exec(*db, "PRAGMA synchronous=NORMAL;", NULL, NULL, NULL);
    (void)sqlite3_busy_timeout(*db, 5000);

    if (db_init_schema(*db) != 0) {
        sqlite3_close(*db);
        *db = NULL;
        return -1;
    }

    usys_log_info("DB ready: %s", dbPath);

    return 0;
}
