/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>
#include <dirent.h>

#include "installer.h"
#include "web_client.h"

#include "usys_log.h"

extern bool app_unpack_package(const char *tarPath, const char *dstDir);

static bool inst_is_dot_name(const char *name) {

    if (!name) return false;
    return (strcmp(name, ".") == 0 || strcmp(name, "..") == 0);
}

static bool inst_set_app_tag(App *app, const char *tag) {

    char *dup;

    if (!app || !tag || !*tag) return false;

    dup = strdup(tag);
    if (!dup) return false;

    free(app->tag);
    app->tag = dup;
    return true;
}

static void inst_trim(char *s) {

    char *start;
    char *end;
    size_t len;

    if (!s || !*s) return;

    start = s;
    while (*start == ' ' || *start == '\t' ||
           *start == '\r' || *start == '\n') {
        start++;
    }

    if (start != s) {
        memmove(s, start, strlen(start) + 1);
    }

    len = strlen(s);
    if (len == 0) return;

    end = s + len - 1;
    while (end >= s &&
           (*end == ' ' || *end == '\t' ||
            *end == '\r' || *end == '\n')) {
        *end = '\0';
        end--;
    }
}

static char* inst_read_version_file(const char *dirPath) {

    char path[512];
    char buf[256];
    FILE *fp;
    char *out;

    fp  = NULL;
    out = NULL;

    if (!dirPath) return NULL;

    snprintf(path, sizeof(path), "%s/VERSION", dirPath);

    fp = fopen(path, "r");
    if (!fp) return NULL;

    if (!fgets(buf, sizeof(buf), fp)) {
        fclose(fp);
        return NULL;
    }

    fclose(fp);

    inst_trim(buf);
    if (!buf[0]) return NULL;

    out = strdup(buf);
    return out;
}

static char* inst_find_latest_pkg_path(Config *config, App *app) {

    DIR *dir;
    struct dirent *de;
    struct stat st;
    char prefix[256];
    char *bestPath;
    time_t bestMtime;
    char fullPath[512];
    size_t prefixLen;
    size_t nameLen;
    const char *name;

    dir       = NULL;
    bestPath  = NULL;
    bestMtime = 0;

    if (!config ||
        !app ||
        !config->pkgsDir ||
        !app->name) {
        return NULL;
    }

    snprintf(prefix, sizeof(prefix), "%s_", app->name);
    prefixLen = strlen(prefix);

    dir = opendir(config->pkgsDir);
    if (!dir) return NULL;

    while ((de = readdir(dir)) != NULL) {
        name = de->d_name;
        if (inst_is_dot_name(name)) {
            continue;
        }

        nameLen = strlen(name);
        if (nameLen <= prefixLen + 7) {
            continue;
        }

        if (strncmp(name, prefix, prefixLen) != 0) {
            continue;
        }

        if (strcmp(name + nameLen - 7, ".tar.gz") != 0) {
            continue;
        }

        snprintf(fullPath, sizeof(fullPath), "%s/%s", config->pkgsDir, name);

        if (stat(fullPath, &st) != 0 || !S_ISREG(st.st_mode)) {
            continue;
        }

        if (!bestPath || st.st_mtime > bestMtime) {
            free(bestPath);
            bestPath  = strdup(fullPath);
            bestMtime = st.st_mtime;
        }
    }

    closedir(dir);
    return bestPath;
}

static char* inst_find_single_child_dir(const char *path) {

    DIR *dir;
    struct dirent *de;
    char *childName;
    char *childPath;
    struct stat st;
    int entries;

    dir       = NULL;
    childName = NULL;
    childPath = NULL;
    entries   = 0;

    if (!path) return NULL;

    dir = opendir(path);
    if (!dir) return NULL;

    while ((de = readdir(dir)) != NULL) {
        if (inst_is_dot_name(de->d_name)) {
            continue;
        }

        entries++;

        if (entries > 1) {
            free(childName);
            closedir(dir);
            return NULL;
        }

        childName = strdup(de->d_name);
        if (!childName) {
            closedir(dir);
            return NULL;
        }
    }

    closedir(dir);

    if (entries != 1 || !childName) {
        free(childName);
        return NULL;
    }

    if (asprintf(&childPath, "%s/%s", path, childName) < 0) {
        childPath = NULL;
    }

    free(childName);

    if (!childPath) return NULL;

    if (stat(childPath, &st) != 0 || !S_ISDIR(st.st_mode)) {
        free(childPath);
        return NULL;
    }

    return childPath;
}

static bool inst_mkdir_p(const char *path) {

    char tmp[512];
    size_t len;
    size_t i;

    if (!path || !*path) return false;

    snprintf(tmp, sizeof(tmp), "%s", path);
    len = strlen(tmp);

    for (i = 1; i < len; i++) {
        if (tmp[i] == '/') {
            tmp[i] = '\0';
            mkdir(tmp, 0755);
            tmp[i] = '/';
        }
    }

    mkdir(tmp, 0755);
    return true;
}

static bool inst_path_exists(const char *p) {

    struct stat st;

    if (!p) {
        return false;
    }

    return stat(p, &st) == 0;
}

static char* inst_app_dir(Config *config, App *app) {

    char *p;

    p = NULL;
    if (asprintf(&p, "%s/%s/%s",
                 config->appsRoot,
                 app->space,
                 app->name) < 0) p = NULL;
    return p;
}

static char* inst_tag_dir(Config *config, App *app, const char *tag) {

    char *p;

    p = NULL;
    if (asprintf(&p, "%s/%s/%s/%s",
                 config->appsRoot,
                 app->space,
                 app->name,
                 tag) < 0) p = NULL;
    return p;
}

static char* inst_stage_dir(Config *config, App *app) {

    char *p;

    p = NULL;
    if (asprintf(&p, "%s/%s/%s/.stage.%d",
                 config->appsRoot,
                 app->space,
                 app->name,
                 getpid()) < 0) p = NULL;
    return p;
}

static char* inst_pkg_path(Config *config, App *app, const char *tag) {

    char *p;

    p = NULL;
    if (asprintf(&p, "%s/%s_%s.tar.gz",
                 config->pkgsDir,
                 app->name,
                 tag) < 0) p = NULL;
    return p;
}

static bool inst_write_symlink_atomic(const char *linkPath, const char *target) {

    char tmp[512];

    if (!linkPath || !target) return false;

    snprintf(tmp, sizeof(tmp), "%s.tmp.%d", linkPath, getpid());
    unlink(tmp);

    if (symlink(target, tmp) != 0) {
        return false;
    }

    if (rename(tmp, linkPath) != 0) {
        unlink(tmp);
        return false;
    }

    return true;
}

bool installer_switch_current(Config *config, App *app) {

    char *appDir;
    char *curLink;
    bool ok;

    appDir = NULL;
    curLink = NULL;
    ok = false;

    if (!config || !app) return false;

    appDir = inst_app_dir(config, app);
    if (!appDir) {
        return false;
    }

    curLink = NULL;
    if (asprintf(&curLink, "%s/current", appDir) < 0) curLink = NULL;
    if (!curLink) {
        free(appDir);
        return false;
    }

    ok = inst_write_symlink_atomic(curLink, app->tag);

    free(curLink);
    free(appDir);
    return ok;
}

bool installer_revert_to_last_good(Config *config, App *app) {

    char *savedTag;

    if (!config || !app) return false;
    if (!app->lastGoodTag || !*app->lastGoodTag) return false;

    savedTag = strdup(app->tag ? app->tag : "");
    free(app->tag);
    app->tag = strdup(app->lastGoodTag);

    if (!installer_switch_current(config, app)) {
        free(app->tag);
        app->tag = savedTag;
        return false;
    }

    free(savedTag);
    return true;
}

bool installer_ensure_installed(Config *config, App *app, const char *hub) {

    char *requestedTag;
    char *tagDir;
    char *pkgPath;
    char *stageDir;
    char *appDir;
    char *payloadDir;
    char *contentDir;
    char *pkgVersion;
    char *localLatestPkg;
    bool ok;

    requestedTag   = NULL;
    tagDir         = NULL;
    pkgPath        = NULL;
    stageDir       = NULL;
    appDir         = NULL;
    payloadDir     = NULL;
    contentDir     = NULL;
    pkgVersion     = NULL;
    localLatestPkg = NULL;
    ok             = false;

    if (!config || !app) return false;

    requestedTag = strdup(app->tag ? app->tag : "");
    if (!requestedTag) return false;

    tagDir = inst_tag_dir(config, app, app->tag);
    if (!tagDir) goto cleanup;

    if (inst_path_exists(tagDir)) {
        app->installState = INSTALL_STATE_SWITCHED;
        ok = true;
        goto cleanup;
    }

    inst_mkdir_p(config->pkgsDir);

    pkgPath = inst_pkg_path(config, app, requestedTag);
    if (!pkgPath) {
        goto cleanup;
    }

    /*
     * Backward-safe local lookup:
     * if requested tag is "latest" but local packages are now versioned,
     * pick the newest matching local package.
     */
    if (!inst_path_exists(pkgPath) && strcmp(requestedTag, "latest") == 0) {
        localLatestPkg = inst_find_latest_pkg_path(config, app);
        if (localLatestPkg) {
            free(pkgPath);
            pkgPath = localLatestPkg;
            localLatestPkg = NULL;
        }
    }

    if (!inst_path_exists(pkgPath)) {
        app->installState = INSTALL_STATE_FETCHING;
        if (!wc_fetch_package(config, app->name, requestedTag, hub, pkgPath)) {
            app->installState = INSTALL_STATE_FAILED;
            usys_log_error("install: fetch failed %s:%s",
                           app->name, requestedTag);
            goto cleanup;
        }
    }

    appDir = inst_app_dir(config, app);
    if (!appDir) goto cleanup;

    inst_mkdir_p(appDir);

    stageDir = inst_stage_dir(config, app);
    if (!stageDir) goto cleanup;

    unlink(stageDir);
    inst_mkdir_p(stageDir);

    app->installState = INSTALL_STATE_STAGING;
    if (!app_unpack_package(pkgPath, stageDir)) {
        app->installState = INSTALL_STATE_FAILED;
        usys_log_error("install: unpack failed %s", pkgPath);
        goto cleanup;
    }

    payloadDir = inst_find_single_child_dir(stageDir);
    contentDir = payloadDir ? payloadDir : stageDir;

    pkgVersion = inst_read_version_file(contentDir);
    if (!pkgVersion || !pkgVersion[0] || strcmp(pkgVersion, "-") == 0) {
        app->installState = INSTALL_STATE_FAILED;
        usys_log_error("install: invalid VERSION in pkg %s", pkgPath);
        goto cleanup;
    }

    if (!inst_set_app_tag(app, pkgVersion)) {
        app->installState = INSTALL_STATE_FAILED;
        goto cleanup;
    }

    free(tagDir);
    tagDir = inst_tag_dir(config, app, app->tag);
    if (!tagDir) {
        app->installState = INSTALL_STATE_FAILED;
        goto cleanup;
    }

    if (inst_path_exists(tagDir)) {
        app->installState = INSTALL_STATE_SWITCHED;
        ok = true;
        goto cleanup;
    }

    /*
     * Flatten packages that unpack as:
     *   <stageDir>/<name>_<version>/...
     * into:
     *   <tagDir>/...
     */
    if (payloadDir) {
        if (rename(payloadDir, tagDir) != 0) {
            app->installState = INSTALL_STATE_FAILED;
            usys_log_error("install: rename failed %s -> %s",
                           payloadDir, tagDir);
            goto cleanup;
        }

        if (rmdir(stageDir) != 0) {
            usys_log_warn("install: unable to remove stage dir %s",
                          stageDir);
        }
    } else {
        if (rename(stageDir, tagDir) != 0) {
            app->installState = INSTALL_STATE_FAILED;
            usys_log_error("install: rename failed %s -> %s",
                           stageDir, tagDir);
            goto cleanup;
        }

        stageDir = NULL;
    }

    app->installState = INSTALL_STATE_SWITCHED;
    ok = true;

cleanup:
    free(localLatestPkg);
    free(pkgVersion);
    if (contentDir != payloadDir) {
        free(contentDir);
    }
    free(payloadDir);
    free(stageDir);
    free(appDir);
    free(pkgPath);
    free(tagDir);
    free(requestedTag);
    return ok;
}
