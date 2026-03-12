/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef WIMC_DB_H
#define WIMC_DB_H

#include <sqlite3.h>

int db_open_or_create(const char *dbPath, sqlite3 **db);
void update_local_db(sqlite3 *db, char *name, char *tag, char *path);
int db_insert_entry(sqlite3 *db, char *name, char *tag, char *status);
int db_update_status(sqlite3 *db, char *name, char *tag, char *status);
int db_update_path_status(sqlite3 *db, char *name, char *tag,
                          char *path, char *status);
int db_mark_old_downloads_failed(sqlite3 *db);
int db_read_status(sqlite3 *db, char *name, char *tag, char **status);

#endif
