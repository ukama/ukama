/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef WIMC_PACKAGE_CACHE_H
#define WIMC_PACKAGE_CACHE_H

#include <stdbool.h>
#include <sqlite3.h>
#include <stddef.h>

#include "wimc.h"

#define WIMC_STATUS_MISSING     "missing"
#define WIMC_STATUS_QUEUED      "queued"
#define WIMC_STATUS_DOWNLOAD    "download"
#define WIMC_STATUS_DOWNLOADING "downloading"
#define WIMC_STATUS_AVAILABLE   "available"
#define WIMC_STATUS_FAILED      "failed"
#define WIMC_STATUS_CORRUPT     "corrupt"

#define WIMC_ALIAS_LATEST       "latest"

typedef struct {
    char name[WIMC_MAX_NAME_LEN];
    char tag[WIMC_MAX_NAME_LEN];
    char path[WIMC_MAX_PATH_LEN];
    char actualVersion[WIMC_MAX_NAME_LEN];
    char error[WIMC_MAX_ERR_STR];
    bool exists;
    bool valid;
    bool alias;
} PackageInfo;

bool pkg_is_valid_identifier(const char *value);
bool pkg_is_alias_tag(const char *tag);

int pkg_path_for_tag(const char *name, const char *tag,
                     char *path, size_t pathLen);
int pkg_tmp_tar_path(const char *uuidStr, const char *name,
                     const char *tag, char *path, size_t pathLen);
int pkg_extract_path(const char *uuidStr, const char *name,
                     const char *tag, char *path, size_t pathLen);

int pkg_ensure_cache_dirs(void);
int pkg_cleanup_tmp(void);
int pkg_read_version_from_dir(const char *dir, char *version,
                              size_t versionLen);
int pkg_read_version_from_tar(const char *path, char *version,
                              size_t versionLen);
int pkg_validate_tar(const char *name, const char *tag, const char *path,
                     PackageInfo *info);
int pkg_publish_from_dir(const char *name, const char *tag,
                         const char *uuidStr, const char *dir,
                         char *publishedPath, size_t publishedPathLen,
                         char *actualVersion,
                         size_t actualVersionLen);
int pkg_publish_tar(const char *name, const char *tag,
                    const char *tmpTar,
                    char *publishedPath, size_t publishedPathLen,
                    char *actualVersion, size_t actualVersionLen);
int pkg_reconcile_startup(sqlite3 *db, const char *pkgDir);

#endif /* WIMC_PACKAGE_CACHE_H */
