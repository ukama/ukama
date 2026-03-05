/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "state_store.h"
#include "app.h"
#include "space.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#include <jansson.h>

#include "usys_log.h"

static bool ss_mkdir_p(const char *path) {

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

    if (mkdir(tmp, 0755) != 0) {
        /* ok if exists */
    }

    return true;
}

static char* ss_state_path(Config *config) {

    char *p;

    p = NULL;
    if (!config || !config->stateDir) return NULL;
    if (asprintf(&p, "%s/state.json", config->stateDir) < 0) p = NULL;
    return p;
}

bool state_store_load(Config *config, Space *spaceList) {

    char *path;
    json_error_t err;
    json_t *root;
    json_t *apps;
    const char *key;
    json_t *val;
    const char *space;
    const char *name;
    const char *desiredTag;
    const char *lastGoodTag;
    App *a;

    if (!config || !spaceList) return false;

    path = ss_state_path(config);
    if (!path) return false;

    root = json_load_file(path, 0, &err);
    free(path);

    if (!root) {
        return true;
    }

    apps = json_object_get(root, "apps");
    if (!json_is_object(apps)) {
        json_decref(root);
        return true;
    }

    json_object_foreach(apps, key, val) {

        if (!key || !val || !json_is_object(val)) continue;

        space = NULL;
        name = NULL;

        char *dot = strchr(key, '.');
        if (!dot) continue;

        char spaceBuf[128];
        size_t slen = (size_t)(dot - key);
        if (slen >= sizeof(spaceBuf)) continue;
        memcpy(spaceBuf, key, slen);
        spaceBuf[slen] = '\0';

        space = spaceBuf;
        name = dot + 1;

        a = app_find(spaceList, space, name);
        if (!a) continue;

        desiredTag = json_is_string(json_object_get(val, "desiredTag")) ?
                     json_string_value(json_object_get(val, "desiredTag")) : NULL;

        lastGoodTag = json_is_string(json_object_get(val, "lastGoodTag")) ?
                      json_string_value(json_object_get(val, "lastGoodTag")) : NULL;

        if (desiredTag && *desiredTag) {
            free(a->tag);
            a->tag = strdup(desiredTag);
        }

        if (lastGoodTag && *lastGoodTag) {
            free(a->lastGoodTag);
            a->lastGoodTag = strdup(lastGoodTag);
        }
    }

    json_decref(root);
    return true;
}

bool state_store_save(Config *config, Space *spaceList) {

    char *path;
    char *tmpPath;
    json_t *root;
    json_t *apps;
    Space *s;
    App *a;

    if (!config || !spaceList) return false;

    ss_mkdir_p(config->stateDir);

    path = ss_state_path(config);
    if (!path) return false;

    tmpPath = NULL;
    if (asprintf(&tmpPath, "%s.tmp.%d", path, getpid()) < 0) tmpPath = NULL;
    if (!tmpPath) {
        free(path);
        return false;
    }

    root = json_object();
    apps = json_object();
    json_object_set_new(root, "apps", apps);

    s = spaceList;
    while (s) {
        a = s->appList;
        while (a) {
            json_t *ja;
            char key[256];

            snprintf(key, sizeof(key), "%s.%s", s->name ? s->name : "", a->name ? a->name : "");

            ja = json_object();
            json_object_set_new(ja, "desiredTag", json_string(a->tag ? a->tag : ""));
            json_object_set_new(ja, "lastGoodTag", json_string(a->lastGoodTag ? a->lastGoodTag : ""));

            json_object_set_new(apps, key, ja);

            a = a->next;
        }
        s = s->next;
    }

    if (json_dump_file(root, tmpPath, JSON_INDENT(2) | JSON_SORT_KEYS) != 0) {
        usys_log_error("state: failed to write %s", tmpPath);
        json_decref(root);
        free(tmpPath);
        free(path);
        return false;
    }

    json_decref(root);

    if (rename(tmpPath, path) != 0) {
        usys_log_error("state: rename failed %s -> %s", tmpPath, path);
        unlink(tmpPath);
    }

    free(tmpPath);
    free(path);
    return true;
}
