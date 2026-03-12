/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "jserdes.h"
#include "app.h"
#include "starterd.h"

#include <string.h>

static const char* app_state_str(AppState st) {

    switch (st) {
        case APP_STATE_STOPPED:  return "stopped";
        case APP_STATE_STARTING: return "starting";
        case APP_STATE_RUNNING:  return "running";
        case APP_STATE_STOPPING: return "stopping";
        case APP_STATE_FAILED:   return "failed";
        default:                 return "unknown";
    }
}

static const char* install_state_str(InstallState st) {

    switch (st) {
        case INSTALL_STATE_NONE:     return "none";
        case INSTALL_STATE_FETCHING: return "fetching";
        case INSTALL_STATE_STAGING:  return "staging";
        case INSTALL_STATE_SWITCHED: return "switched";
        case INSTALL_STATE_FAILED:   return "failed";
        default:                     return "unknown";
    }
}

json_t* jserdes_status_json(Space *spaceList) {

    json_t *root;
    json_t *spaces;
    json_t *js;
    json_t *apps;
    Space *s;
    App *a;

    root = json_object();
    spaces = json_array();
    json_object_set_new(root, "spaces", spaces);

    s = spaceList;
    while (s) {

        js = json_object();
        json_object_set_new(js, "name", json_string(s->name ? s->name : ""));

        apps = json_array();
        json_object_set_new(js, "apps", apps);

        a = s->appList;
        while (a) {
            json_t *ja;

            ja = json_object();
            json_object_set_new(ja, "name", json_string(a->name ? a->name : ""));
            json_object_set_new(ja, "tag", json_string(a->tag ? a->tag : ""));
            json_object_set_new(ja, "state", json_string(app_state_str(a->state)));
            json_object_set_new(ja, "installState", json_string(install_state_str(a->installState)));
            json_object_set_new(ja, "pid", json_integer((json_int_t)a->pid));
            json_object_set_new(ja, "pgid", json_integer((json_int_t)a->pgid));
            json_object_set_new(ja, "lastExitCode", json_integer((json_int_t)a->lastExitCode));
            json_object_set_new(ja, "lastExitSignal", json_integer((json_int_t)a->lastExitSignal));
            json_object_set_new(ja, "lastGoodTag", json_string(a->lastGoodTag ? a->lastGoodTag : ""));

            json_array_append_new(apps, ja);
            a = a->next;
        }

        json_array_append_new(spaces, js);
        s = s->next;
    }

    return root;
}
