/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include "manifest.h"
#include "space.h"
#include "app.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <jansson.h>

#include "usys_log.h"

static bool m_is_valid_name(const char *s) {

    const char *p;

    if (!s || !*s) return false;

    p = s;
    while (*p) {
        if ((*p >= 'a' && *p <= 'z') ||
            (*p >= 'A' && *p <= 'Z') ||
            (*p >= '0' && *p <= '9') ||
            (*p == '-') || (*p == '_') || (*p == '.')) {
            p++;
            continue;
        }
        return false;
    }

    return true;
}

static char** m_parse_str_array(json_t *arr, int *countOut) {

    size_t i;
    size_t n;
    char **out;
    json_t *v;

    if (countOut) *countOut = 0;
    if (!arr || !json_is_array(arr)) return NULL;

    n = json_array_size(arr);
    if (n == 0) return NULL;

    out = calloc(n + 1, sizeof(char*));
    if (!out) return NULL;

    for (i = 0; i < n; i++) {
        v = json_array_get(arr, i);
        if (!json_is_string(v)) continue;
        out[i] = strdup(json_string_value(v));
    }

    out[n] = NULL;
    if (countOut) *countOut = (int)n;
    return out;
}

static char** m_parse_env_object(json_t *obj, int *countOut) {

    const char *key;
    json_t *val;
    size_t n;
    size_t i;
    char **out;
    char *kv;
    const char *vs;

    if (countOut) *countOut = 0;
    if (!obj || !json_is_object(obj)) return NULL;

    n = json_object_size(obj);
    if (n == 0) return NULL;

    out = calloc(n + 1, sizeof(char*));
    if (!out) return NULL;

    i = 0;
    json_object_foreach(obj, key, val) {
        if (!key) continue;
        if (json_is_string(val)) {
            vs = json_string_value(val);
        } else if (json_is_integer(val)) {
            static char tmp[64];
            snprintf(tmp, sizeof(tmp), "%lld", (long long)json_integer_value(val));
            vs = tmp;
        } else if (json_is_boolean(val)) {
            vs = json_is_true(val) ? "true" : "false";
        } else {
            continue;
        }

        kv = NULL;
        if (asprintf(&kv, "%s=%s", key, vs) < 0) kv = NULL;
        out[i++] = kv;
        if (i >= n) break;
    }

    out[i] = NULL;
    if (countOut) *countOut = (int)i;
    return out;
}

static void m_free_argv(char **argv) {

    int i;

    if (!argv) return;
    for (i = 0; argv[i]; i++) free(argv[i]);
    free(argv);
}

static void m_free_envp(char **envp) {

    int i;

    if (!envp) return;
    for (i = 0; envp[i]; i++) free(envp[i]);
    free(envp);
}

static App* m_parse_app(const char *spaceName, json_t *j) {

    App *a;
    json_t *v;
    const char *name;
    const char *tag;
    const char *cmd;
    const char *workdir;

    if (!json_is_object(j)) return NULL;

    v = json_object_get(j, "name");
    name = json_is_string(v) ? json_string_value(v) : NULL;

    v = json_object_get(j, "tag");
    tag = json_is_string(v) ? json_string_value(v) : NULL;

    v = json_object_get(j, "cmd");
    cmd = json_is_string(v) ? json_string_value(v) : NULL;

    v = json_object_get(j, "workdir");
    workdir = json_is_string(v) ? json_string_value(v) : NULL;

    if (!m_is_valid_name(name) || !m_is_valid_name(spaceName)) {
        usys_log_error("manifest: invalid app/space name");
        return NULL;
    }

    if (!tag || !*tag) tag = "latest";
    if (!cmd || !*cmd) {
        usys_log_error("manifest: missing cmd for %s/%s", spaceName, name);
        return NULL;
    }

    a = calloc(1, sizeof(*a));
    if (!a) return NULL;

    a->space = strdup(spaceName);
    a->name  = strdup(name);
    a->tag   = strdup(tag);
    a->cmd   = strdup(cmd);
    a->workdir = workdir ? strdup(workdir) : NULL;

    v = json_object_get(j, "argv");
    a->argv = m_parse_str_array(v, &a->argc);
    if (!a->argv) {
        a->argv = calloc(2, sizeof(char*));
        if (a->argv) {
            a->argv[0] = strdup(cmd);
            a->argv[1] = NULL;
            a->argc = 1;
        }
    }

    v = json_object_get(j, "env");
    a->envp = m_parse_env_object(v, &a->envc);

    v = json_object_get(j, "port");
    a->port = json_is_integer(v) ? (int)json_integer_value(v) : 0;

    a->state        = APP_STATE_STOPPED;
    a->installState = INSTALL_STATE_NONE;
    a->pid  = -1;
    a->pgid = -1;
    a->lastGoodTag = strdup(tag);

    return a;
}

static void m_free_app(App *a) {

    if (!a) return;

    free(a->space);
    free(a->name);
    free(a->tag);
    free(a->cmd);
    m_free_argv(a->argv);
    m_free_envp(a->envp);
    free(a->workdir);
    free(a->lastGoodTag);
    free(a);
}

static void m_free_space(Space *s) {

    App *a;
    App *n;

    if (!s) return;

    a = s->appList;
    while (a) {
        n = a->next;
        m_free_app(a);
        a = n;
    }

    free(s->name);
    free(s);
}

bool manifest_load(Config *config, Space **spaceListOut) {

    json_error_t err;
    json_t *root;
    json_t *spaces;
    size_t i;
    size_t n;
    json_t *js;
    json_t *apps;
    json_t *capps;
    json_t *arr;
    const char *spaceName;
    Space *head;
    Space *tail;
    Space *s;
    App *a;
    App *alast;

    if (!config || !spaceListOut) return false;

    *spaceListOut = NULL;

    root = json_load_file(config->manifestPath, 0, &err);
    if (!root) {
        usys_log_error("manifest: load failed %s:%d %s", err.source, err.line, err.text);
        return false;
    }

    spaces = json_object_get(root, "spaces");
    if (!json_is_array(spaces)) {
        usys_log_error("manifest: missing spaces array");
        json_decref(root);
        return false;
    }

    head = NULL;
    tail = NULL;

    n = json_array_size(spaces);
    for (i = 0; i < n; i++) {

        js = json_array_get(spaces, i);
        if (!json_is_object(js)) continue;

        spaceName = NULL;
        if (json_is_string(json_object_get(js, "name"))) {
            spaceName = json_string_value(json_object_get(js, "name"));
        }

        if (!m_is_valid_name(spaceName)) {
            usys_log_error("manifest: invalid space name");
            continue;
        }

        s = calloc(1, sizeof(*s));
        if (!s) continue;

        s->name = strdup(spaceName);
        s->appList = NULL;
        s->next = NULL;

        apps = json_object_get(js, "apps");
        capps = json_object_get(js, "capps");
        arr = json_is_array(apps) ? apps : (json_is_array(capps) ? capps : NULL);
        if (!arr) {
            usys_log_error("manifest: space %s missing apps/capps array", spaceName);
            m_free_space(s);
            continue;
        }

        alast = NULL;
        for (size_t j = 0; j < json_array_size(arr); j++) {
            a = m_parse_app(spaceName, json_array_get(arr, j));
            if (!a) continue;

            if (!s->appList) s->appList = a;
            if (alast) alast->next = a;
            alast = a;
        }

        if (!s->appList) {
            usys_log_error("manifest: space %s has no valid apps", spaceName);
            m_free_space(s);
            continue;
        }

        if (!head) head = s;
        if (tail) tail->next = s;
        tail = s;
    }

    json_decref(root);

    if (!head) {
        usys_log_error("manifest: no valid spaces");
        return false;
    }

    *spaceListOut = head;
    return true;
}

void manifest_free(Space *spaceList) {

    Space *s;
    Space *n;

    s = spaceList;
    while (s) {
        n = s->next;
        m_free_space(s);
        s = n;
    }
}
