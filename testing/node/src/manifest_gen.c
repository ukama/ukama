/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>

#include <jansson.h>

#include "config.h"
#include "log.h"
#include "manifest_gen.h"

#define STARTER_TARGET_DEFAULT  "virtual-node"
#define STARTER_VERSION_DEFAULT "0.1"

static int streq(const char *a, const char *b) {

    return (a && b && strcmp(a, b) == 0);
}

static int is_boot_app(const char *name) {

    if (!name) {
        return FALSE;
    }

    return streq(name, "noded")     ||
           streq(name, "bootstrap") ||
           streq(name, "meshd");
}

static int append_arg_tokens(json_t *jargv, const char *args) {

    char *dup  = NULL;
    char *tok  = NULL;
    char *save = NULL;

    if (!jargv) {
        return FALSE;
    }

    if (!args || *args == '\0') {
        return TRUE;
    }

    dup = strdup(args);
    if (!dup) {
        return FALSE;
    }

    tok = strtok_r(dup, " ", &save);
    while (tok) {
        if (json_array_append_new(jargv, json_string(tok)) != 0) {
            free(dup);
            return FALSE;
        }
        tok = strtok_r(NULL, " ", &save);
    }

    free(dup);
    return TRUE;
}

static int append_env_object(json_t *japp, EnvVar *env) {

    json_t *jenv = NULL;
    EnvVar *curr = NULL;

    if (!japp) {
        return FALSE;
    }

    if (!env) {
        return TRUE;
    }

    jenv = json_object();
    if (!jenv) {
        return FALSE;
    }

    for (curr = env; curr; curr = curr->next) {
        if (!curr->key || !curr->value) {
            continue;
        }

        if (json_object_set_new(jenv,
                                curr->key,
                                json_string(curr->value)) != 0) {
            json_decref(jenv);
            return FALSE;
        }
    }

    if (json_object_size(jenv) == 0) {
        json_decref(jenv);
        return TRUE;
    }

    if (json_object_set_new(japp, "env", jenv) != 0) {
        json_decref(jenv);
        return FALSE;
    }

    return TRUE;
}

static int manifest_add_app(json_t *apps, Config *config) {

    json_t *japp = NULL;
    json_t *jargv = NULL;
    char cmd[512] = {0};

    const char *name  = NULL;
    const char *tag   = NULL;
    const char *bin   = NULL;
    const char *binTo = NULL;

    if (!apps || !config || !config->capp || !config->build) {
        return FALSE;
    }

    if (!config->capp->name    ||
        !config->capp->version ||
        !config->capp->bin     ||
        !config->build->binTo) {
        return FALSE;
    }

    name  = config->capp->name;
    tag   = config->capp->version;
    bin   = config->capp->bin;
    binTo = config->build->binTo;

    snprintf(cmd,
             sizeof(cmd),
             "%s/%s",
             (binTo[0] == '/') ? binTo + 1 : binTo,
             bin);

    japp = json_object();
    if (!japp) {
        return FALSE;
    }

    json_object_set_new(japp, "name", json_string(name));
    json_object_set_new(japp, "tag",  json_string(tag));
    json_object_set_new(japp, "cmd",  json_string(cmd));

    jargv = json_array();
    if (!jargv) {
        json_decref(japp);
        return FALSE;
    }

    if (json_array_append_new(jargv, json_string(bin)) != 0) {
        json_decref(jargv);
        json_decref(japp);
        return FALSE;
    }

    if (!append_arg_tokens(jargv, config->capp->args)) {
        json_decref(jargv);
        json_decref(japp);
        return FALSE;
    }

    json_object_set_new(japp, "argv", jargv);

    if (!append_env_object(japp, config->capp->env)) {
        json_decref(japp);
        return FALSE;
    }

    if (json_array_append_new(apps, japp) != 0) {
        json_decref(japp);
        return FALSE;
    }

    return TRUE;
}

static Config* find_app_config(Configs *configs, const char *name) {

    Configs *ptr = NULL;

    if (!configs || !name) {
        return NULL;
    }

    for (ptr = configs; ptr; ptr = ptr->next) {
        if (!ptr->valid || !ptr->config || !ptr->config->capp) {
            continue;
        }

        if (ptr->config->capp->name &&
            strcmp(ptr->config->capp->name, name) == 0) {
            return ptr->config;
        }
    }

    return NULL;
}

static int append_boot_apps_in_order(json_t *bootApps, Configs *configs) {

    static const char *bootOrder[] = {
        "noded",
        "bootstrap",
        "meshd"
    };

    size_t i;
    Config *cfg = NULL;

    if (!bootApps || !configs) {
        return FALSE;
    }

    for (i = 0; i < sizeof(bootOrder) / sizeof(bootOrder[0]); i++) {
        cfg = find_app_config(configs, bootOrder[i]);
        if (!cfg) {
            continue;
        }

        if (!manifest_add_app(bootApps, cfg)) {
            return FALSE;
        }
    }

    return TRUE;
}

static int append_service_apps(json_t *svcApps, Configs *configs) {

    Configs *ptr = NULL;
    const char *name = NULL;

    if (!svcApps || !configs) {
        return FALSE;
    }

    for (ptr = configs; ptr; ptr = ptr->next) {
        if (!ptr->valid || !ptr->config || !ptr->config->capp) {
            continue;
        }

        name = ptr->config->capp->name;
        if (!name) {
            continue;
        }

        if (is_boot_app(name)) {
            continue;
        }

        if (!manifest_add_app(svcApps, ptr->config)) {
            return FALSE;
        }
    }

    return TRUE;
}

int create_manifest_config(Configs *configs) {

    FILE *fp = NULL;
    json_t *root = NULL;
    json_t *spaces = NULL;
    json_t *boot = NULL;
    json_t *services = NULL;
    json_t *reboot = NULL;
    json_t *bootApps = NULL;
    json_t *svcApps = NULL;
    json_t *rebootApps = NULL;
    char *out = NULL;

    if (!configs) {
        return FALSE;
    }

    root       = json_object();
    spaces     = json_array();
    boot       = json_object();
    services   = json_object();
    reboot     = json_object();
    bootApps   = json_array();
    svcApps    = json_array();
    rebootApps = json_array();

    if (!root       || !spaces   || !boot     || !services ||
        !reboot     || !bootApps || !svcApps  || !rebootApps) {
        goto fail;
    }

    json_object_set_new(root, "version", json_string(STARTER_VERSION_DEFAULT));
    json_object_set_new(root, "target",  json_string(STARTER_TARGET_DEFAULT));

    json_object_set_new(boot, "name", json_string("boot"));
    json_object_set_new(boot, "apps", bootApps);
    bootApps = NULL;

    json_object_set_new(services, "name", json_string("services"));
    json_object_set_new(services, "apps", svcApps);
    svcApps = NULL;

    json_object_set_new(reboot, "name", json_string("reboot"));
    json_object_set_new(reboot, "apps", rebootApps);
    rebootApps = NULL;

    if (!append_boot_apps_in_order(json_object_get(boot, "apps"), configs)) {
        goto fail;
    }

    if (!append_service_apps(json_object_get(services, "apps"), configs)) {
        goto fail;
    }

    if (json_array_append_new(spaces, boot) != 0) {
        goto fail;
    }
    boot = NULL;

    if (json_array_append_new(spaces, services) != 0) {
        goto fail;
    }
    services = NULL;

    if (json_array_append_new(spaces, reboot) != 0) {
        goto fail;
    }
    reboot = NULL;

    json_object_set_new(root, "spaces", spaces);
    spaces = NULL;

    out = json_dumps(root, JSON_INDENT(2));
    if (!out) {
        goto fail;
    }

    fp = fopen(MANIFEST_FILENAME, "w");
    if (!fp) {
        log_error("Unable to open manifest file: %s error: %s",
                  MANIFEST_FILENAME, strerror(errno));
        goto fail;
    }

    if (fwrite(out, strlen(out), 1, fp) <= 0) {
        log_error("Unable to write manifest file: %s error: %s",
                  MANIFEST_FILENAME, strerror(errno));
        fclose(fp);
        free(out);
        json_decref(root);
        return FALSE;
    }

    fclose(fp);
    free(out);
    json_decref(root);
    return TRUE;

fail:
    if (fp)         fclose(fp);
    if (out)        free(out);
    if (root)       json_decref(root);
    if (spaces)     json_decref(spaces);
    if (boot)       json_decref(boot);
    if (services)   json_decref(services);
    if (reboot)     json_decref(reboot);
    if (bootApps)   json_decref(bootApps);
    if (svcApps)    json_decref(svcApps);
    if (rebootApps) json_decref(rebootApps);
    return FALSE;
}

void purge_manifest_config(const char *fileName) {

    if (!fileName) return;

    if (remove(fileName) == 0) {
        log_debug("manifest config file removed: %s", fileName);
    }
}
