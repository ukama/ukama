/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "installer.h"
#include "web_client.h"

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#include "usys_log.h"

extern bool app_unpack_package(const char *tarPath, const char *dstDir);

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
    if (asprintf(&p, "%s/%s-%s.tar.gz",
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

    char *tagDir;
    char *pkgPath;
    char *stageDir;
    char *appDir;
    bool ok;

    tagDir = NULL;
    pkgPath = NULL;
    stageDir = NULL;
    appDir = NULL;
    ok = false;

    if (!config || !app || !hub || !*hub) return false;

    tagDir = inst_tag_dir(config, app, app->tag);
    if (!tagDir) return false;

    if (inst_path_exists(tagDir)) {
        free(tagDir);
        return true;
    }

    inst_mkdir_p(config->pkgsDir);

    pkgPath = inst_pkg_path(config, app, app->tag);
    if (!pkgPath) {
        free(tagDir);
        return false;
    }

    if (!inst_path_exists(pkgPath)) {
        app->installState = INSTALL_STATE_FETCHING;
        if (!wc_fetch_package(config, app->name, app->tag, hub, pkgPath)) {
            app->installState = INSTALL_STATE_FAILED;
            usys_log_error("install: fetch failed %s:%s hub=%s",
                           app->name, app->tag, hub);
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
        goto cleanup;
    }

    if (rename(stageDir, tagDir) != 0) {
        app->installState = INSTALL_STATE_FAILED;
        usys_log_error("install: rename failed %s -> %s", stageDir, tagDir);
        goto cleanup;
    }

    stageDir = NULL;
    app->installState = INSTALL_STATE_SWITCHED;

    ok = true;

cleanup:
    if (stageDir) {
        app_remove_dir_recursive(stageDir);
        free(stageDir);
    }
    free(appDir);
    free(pkgPath);
    free(tagDir);
    return ok;
}
