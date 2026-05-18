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

static bool inst_file_exists_non_empty(const char *p) {

    struct stat st;

    if (!p || !*p) {
        return false;
    }

    if (stat(p, &st) != 0) {
        return false;
    }

    if (!S_ISREG(st.st_mode)) {
        return false;
    }

    return st.st_size > 0;
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

static bool inst_make_pkg_path(Config *config,
                               App *app,
                               const char *tag,
                               char *path,
                               size_t pathLen) {

    int ret;

    if (!config || !app || !tag || !path || pathLen == 0) {
        return false;
    }

    ret = snprintf(path,
                   pathLen,
                   "%s/%s_%s.tar.gz",
                   config->pkgsDir,
                   app->name,
                   tag);

    return ret > 0 && (size_t)ret < pathLen;
}

static bool inst_try_local_pkg(Config *config,
                               App *app,
                               const char *tag,
                               char **pathOut) {

    char path[512];

    if (!config || !app || !tag || !pathOut) {
        return false;
    }

    if (!inst_make_pkg_path(config, app, tag, path, sizeof(path))) {
        return false;
    }

    if (!inst_file_exists_non_empty(path)) {
        return false;
    }

    *pathOut = strdup(path);
    return *pathOut != NULL;
}

static bool inst_find_local_package(Config *config,
                                    App *app,
                                    const char *requestedTag,
                                    char **pathOut) {

    char altTag[256];

    if (pathOut != NULL) {
        *pathOut = NULL;
    }

    if (!config || !app || !requestedTag || !pathOut) {
        return false;
    }

    /*
     * Primary form:
     *   /ukama/apps/pkgs/example_latest.tar.gz
     *   /ukama/apps/pkgs/example_v1.0.0-abcdefgh.tar.gz
     */
    if (inst_try_local_pkg(config, app, requestedTag, pathOut)) {
        return true;
    }

    /*
     * Accept both v-prefixed and non-v-prefixed tags.
     */
    if (requestedTag[0] != 'v') {
        snprintf(altTag, sizeof(altTag), "v%s", requestedTag);
        if (inst_try_local_pkg(config, app, altTag, pathOut)) {
            return true;
        }
    } else if (requestedTag[1] != '\0') {
        snprintf(altTag, sizeof(altTag), "%s", requestedTag + 1);
        if (inst_try_local_pkg(config, app, altTag, pathOut)) {
            return true;
        }
    }

    return false;
}

static bool inst_actual_version_installed(Config *config,
                                          App *app,
                                          const char *actualVersion) {

    char *tagDir;
    bool ok;

    tagDir = NULL;
    ok = false;

    if (!config || !app || !actualVersion || !*actualVersion) {
        return false;
    }

    tagDir = inst_tag_dir(config, app, actualVersion);
    if (!tagDir) {
        return false;
    }

    ok = inst_path_exists(tagDir);
    free(tagDir);

    return ok;
}

static bool inst_install_from_package(Config *config,
                                      App *app,
                                      const char *requestedTag,
                                      const char *pkgPath) {

    char *stageDir;
    char *appDir;
    char *payloadDir;
    char *contentDir;
    char *pkgVersion;
    char *tagDir;
    bool ok;

    stageDir = NULL;
    appDir = NULL;
    payloadDir = NULL;
    contentDir = NULL;
    pkgVersion = NULL;
    tagDir = NULL;
    ok = false;

    if (!config || !app || !requestedTag || !pkgPath || !*pkgPath) {
        return false;
    }

    if (!inst_file_exists_non_empty(pkgPath)) {
        app->installState = INSTALL_STATE_FAILED;
        usys_log_error("install: package missing or empty %s", pkgPath);
        return false;
    }

    appDir = inst_app_dir(config, app);
    if (!appDir) {
        app->installState = INSTALL_STATE_FAILED;
        goto cleanup;
    }

    inst_mkdir_p(appDir);

    stageDir = inst_stage_dir(config, app);
    if (!stageDir) {
        app->installState = INSTALL_STATE_FAILED;
        goto cleanup;
    }

    inst_mkdir_p(stageDir);

    app->installState = INSTALL_STATE_STAGING;

    usys_log_info("install: unpacking local package %s for %s:%s",
                  pkgPath,
                  app->name,
                  requestedTag);

    if (!app_unpack_package(pkgPath, stageDir)) {
        app->installState = INSTALL_STATE_FAILED;
        usys_log_error("install: unpack failed %s", pkgPath);
        goto cleanup;
    }

    payloadDir = inst_find_single_child_dir(stageDir);
    contentDir = payloadDir ? payloadDir : stageDir;

    pkgVersion = inst_read_version_file(contentDir);
    if (pkgVersion == NULL || pkgVersion[0] == '\0' ||
        strcmp(pkgVersion, "-") == 0) {
        app->installState = INSTALL_STATE_FAILED;
        usys_log_error("install: invalid VERSION in pkg %s", pkgPath);
        goto cleanup;
    }

    if (!inst_set_app_tag(app, pkgVersion)) {
        app->installState = INSTALL_STATE_FAILED;
        goto cleanup;
    }

    tagDir = inst_tag_dir(config, app, app->tag);
    if (!tagDir) {
        app->installState = INSTALL_STATE_FAILED;
        goto cleanup;
    }

    /*
     * If another previous install already unpacked this VERSION, accept it.
     */
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
    if (payloadDir != NULL) {
        if (rename(payloadDir, tagDir) != 0) {
            app->installState = INSTALL_STATE_FAILED;
            usys_log_error("install: rename failed %s -> %s",
                           payloadDir,
                           tagDir);
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
                           stageDir,
                           tagDir);
            goto cleanup;
        }

        stageDir = NULL;
    }

    app->installState = INSTALL_STATE_SWITCHED;
    ok = true;

cleanup:
    free(pkgVersion);
    free(payloadDir);
    free(stageDir);
    free(appDir);
    free(tagDir);

    return ok;
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
    char *wimcVersion;
    bool ok;

    requestedTag = NULL;
    tagDir = NULL;
    pkgPath = NULL;
    wimcVersion = NULL;
    ok = false;

    if (config == NULL || app == NULL || app->name == NULL || app->tag == NULL) {
        return false;
    }

    requestedTag = strdup(app->tag);
    if (requestedTag == NULL) {
        return false;
    }

    /*
     * 1. Already installed exact requested tag.
     */
    tagDir = inst_tag_dir(config, app, requestedTag);
    if (tagDir == NULL) {
        goto cleanup;
    }

    if (inst_path_exists(tagDir)) {
        app->installState = INSTALL_STATE_SWITCHED;
        ok = true;
        goto cleanup;
    }

    free(tagDir);
    tagDir = NULL;

    /*
     * 2. Local package cache.
     *
     * WIMC owns writes/fetches into /ukama/apps/pkgs.
     * Starter is allowed to read and install from this cache.
     *
     * This solves first boot:
     *   starter can install example_latest.tar.gz even when wimc.d
     *   is not running yet.
     */
    if (inst_find_local_package(config, app, requestedTag, &pkgPath)) {
        usys_log_info("install: local package found %s:%s path=%s",
                      app->name,
                      requestedTag,
                      pkgPath);

        ok = inst_install_from_package(config, app, requestedTag, pkgPath);
        goto cleanup;
    }

    /*
     * 3. Ask WIMC only when local tarball is missing.
     *
     * If WIMC is unavailable, this becomes pending rather than fatal.
     * Boot supervisor decides whether pending should merely degrade boot
     * or fail an explicit update operation.
     */
    app->installState = INSTALL_STATE_FETCHING;

    if (!wc_fetch_package(config,
                          app->name,
                          requestedTag,
                          hub,
                          &pkgPath,
                          &wimcVersion)) {
        app->installState = INSTALL_STATE_PENDING;
        app->state = APP_STATE_STOPPED;
        app->reason = APP_REASON_PACKAGE_MISSING;

        usys_log_warn("install: package pending %s:%s "
                      "local package missing and wimc unavailable",
                      app->name,
                      requestedTag);
        goto cleanup;
    }

    /*
     * WIMC may resolve alias/latest to the actual VERSION.
     */
    if (wimcVersion != NULL && *wimcVersion != '\0') {
        if (inst_actual_version_installed(config, app, wimcVersion)) {
            if (!inst_set_app_tag(app, wimcVersion)) {
                app->installState = INSTALL_STATE_FAILED;
                goto cleanup;
            }

            app->installState = INSTALL_STATE_SWITCHED;
            ok = true;
            goto cleanup;
        }
    }

    if (pkgPath == NULL || *pkgPath == '\0') {
        app->installState = INSTALL_STATE_FAILED;
        usys_log_error("install: wimc returned empty package path %s:%s",
                       app->name,
                       requestedTag);
        goto cleanup;
    }

    ok = inst_install_from_package(config, app, requestedTag, pkgPath);

cleanup:
    free(wimcVersion);
    free(pkgPath);
    free(tagDir);
    free(requestedTag);

    return ok;
}
