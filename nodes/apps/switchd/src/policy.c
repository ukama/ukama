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
#include <sys/stat.h>
#include <unistd.h>

#include <jansson.h>

#include "json_types.h"
#include "policy.h"
#include "utils.h"

#include "usys_log.h"

#define POLICY_TMP_SUFFIX ".tmp"
#define POLICY_HASH_PREFIX "fnv64:"

static void set_err(char *err, size_t errLen, const char *msg) {

    if (err && errLen > 0) {
        snprintf(err, errLen, "%s", msg ? msg : "");
    }
}

const char *policy_state_to_str(SwitchPolicyState state) {

    switch (state) {
    case SWITCH_POLICY_STATE_MISSING:
        return "missing";
    case SWITCH_POLICY_STATE_LOADED:
        return "loaded";
    case SWITCH_POLICY_STATE_INVALID:
        return "invalid";
    default:
        return "unknown";
    }
}

const char *policy_type_to_str(SwitchPortPolicyType type) {

    switch (type) {
    case SWITCH_PORT_POLICY_PROTECTED:
        return "protected";
    case SWITCH_PORT_POLICY_FREE:
        return "free";
    case SWITCH_PORT_POLICY_NEVER_OFF_REMOTE:
        return "never_off_remote";
    case SWITCH_PORT_POLICY_DISABLED:
        return "disabled";
    case SWITCH_PORT_POLICY_UNKNOWN:
    default:
        return "unknown";
    }
}

const char *policy_action_to_str(SwitchPolicyAction action) {

    switch (action) {
    case SWITCH_POLICY_ACTION_ADMIN_UP:
        return "admin_up";
    case SWITCH_POLICY_ACTION_ADMIN_DOWN:
        return "admin_down";
    case SWITCH_POLICY_ACTION_POE_ON:
        return "poe_on";
    case SWITCH_POLICY_ACTION_POE_OFF:
        return "poe_off";
    case SWITCH_POLICY_ACTION_POE_CYCLE:
        return "poe_cycle";
    default:
        return "unknown";
    }
}

static SwitchPortPolicyType policy_type_from_str(const char *value) {

    if (value == NULL) {
        return SWITCH_PORT_POLICY_UNKNOWN;
    }

    if (strcmp(value, "protected") == 0) {
        return SWITCH_PORT_POLICY_PROTECTED;
    }

    if (strcmp(value, "free") == 0) {
        return SWITCH_PORT_POLICY_FREE;
    }

    if (strcmp(value, "never_off_remote") == 0) {
        return SWITCH_PORT_POLICY_NEVER_OFF_REMOTE;
    }

    if (strcmp(value, "disabled") == 0) {
        return SWITCH_PORT_POLICY_DISABLED;
    }

    return SWITCH_PORT_POLICY_UNKNOWN;
}

static bool action_turns_off(SwitchPolicyAction action) {

    return (action == SWITCH_POLICY_ACTION_ADMIN_DOWN ||
            action == SWITCH_POLICY_ACTION_POE_OFF ||
            action == SWITCH_POLICY_ACTION_POE_CYCLE);
}

static void policy_clear(SwitchPolicy *policy) {

    if (!policy) {
        return;
    }

    memset(policy, 0, sizeof(*policy));
    policy->state = SWITCH_POLICY_STATE_MISSING;
}

static void policy_hash_body(const char *body,
                             size_t len,
                             char *dst,
                             size_t dstLen) {

    size_t i;
    unsigned long long hash;

    if (dst == NULL || dstLen == 0) {
        return;
    }

    hash = 1469598103934665603ULL;

    if (body != NULL) {
        for (i = 0; i < len; i++) {
            hash ^= (unsigned char)body[i];
            hash *= 1099511628211ULL;
        }
    }

    snprintf(dst, dstLen, POLICY_HASH_PREFIX "%016llx", hash);
}

static int read_policy_file(const char *path,
                            char **body,
                            size_t *bodyLen,
                            time_t *mtime) {

    FILE *fp;
    struct stat st;
    char *buffer;
    size_t nread;

    if (body == NULL || bodyLen == NULL || path == NULL) {
        return SWITCHD_ERR_INVAL;
    }

    *body = NULL;
    *bodyLen = 0;

    if (stat(path, &st) != 0) {
        return SWITCHD_ERR_NOTFOUND;
    }

    fp = fopen(path, "rb");
    if (fp == NULL) {
        return SWITCHD_ERR_IO;
    }

    buffer = calloc(1, (size_t)st.st_size + 1);
    if (buffer == NULL) {
        fclose(fp);
        return SWITCHD_ERR_NOMEM;
    }

    nread = fread(buffer, 1, (size_t)st.st_size, fp);
    fclose(fp);

    if (nread != (size_t)st.st_size) {
        free(buffer);
        return SWITCHD_ERR_IO;
    }

    buffer[nread] = '\0';

    *body = buffer;
    *bodyLen = nread;

    if (mtime) {
        *mtime = st.st_mtime;
    }

    return SWITCHD_OK;
}

static int ensure_parent_dir(const char *path) {

    char dir[SWITCHD_STAGE_PATH_LEN];
    char *slash;

    if (path == NULL || *path == '\0') {
        return SWITCHD_ERR_INVAL;
    }

    snprintf(dir, sizeof(dir), "%s", path);

    slash = strrchr(dir, '/');
    if (slash == NULL) {
        return SWITCHD_OK;
    }

    *slash = '\0';
    if (dir[0] == '\0') {
        return SWITCHD_OK;
    }

    return (mkdir_p(dir, 0755) == 0) ? SWITCHD_OK : SWITCHD_ERR_IO;
}

static int write_policy_file(const char *path, const char *body, size_t len) {

    char tmp[SWITCHD_STAGE_PATH_LEN + 8];
    FILE *fp;

    if (!path || !body) {
        return SWITCHD_ERR_INVAL;
    }

    if (ensure_parent_dir(path) != SWITCHD_OK) {
        return SWITCHD_ERR_IO;
    }

    snprintf(tmp, sizeof(tmp), "%s%s", path, POLICY_TMP_SUFFIX);

    fp = fopen(tmp, "wb");
    if (fp == NULL) {
        return SWITCHD_ERR_IO;
    }

    if (fwrite(body, 1, len, fp) != len) {
        fclose(fp);
        unlink(tmp);
        return SWITCHD_ERR_IO;
    }

    if (fclose(fp) != 0) {
        unlink(tmp);
        return SWITCHD_ERR_IO;
    }

    if (rename(tmp, path) != 0) {
        unlink(tmp);
        return SWITCHD_ERR_IO;
    }

    return SWITCHD_OK;
}

static bool copy_json_str(JsonObj *root,
                          const char *key,
                          char *dst,
                          size_t dstLen) {

    JsonObj *entry;
    const char *value;

    if (!root || !key || !dst || dstLen == 0) {
        return false;
    }

    entry = json_object_get(root, key);
    if (entry == NULL) {
        dst[0] = '\0';
        return true;
    }

    if (!json_is_string(entry)) {
        return false;
    }

    value = json_string_value(entry);
    snprintf(dst, dstLen, "%s", value ? value : "");
    return true;
}

static int parse_policy(JsonObj *root,
                        const char *path,
                        const char *hash,
                        time_t fileMtime,
                        SwitchPolicy *out,
                        char *err,
                        size_t errLen) {

    JsonObj *ports;
    JsonObj *item;
    size_t index;

    if (!root || !out) {
        set_err(err, errLen, "invalid policy");
        return SWITCHD_ERR_INVAL;
    }

    policy_clear(out);

    out->state     = SWITCH_POLICY_STATE_LOADED;
    out->loadedAt  = time(NULL);
    out->fileMtime = fileMtime;
    out->checkedAt = time(NULL);

    snprintf(out->path, sizeof(out->path), "%s", path ? path : "");
    snprintf(out->hash, sizeof(out->hash), "%s", hash ? hash : "");

    if (!copy_json_str(root, "site_id", out->siteId, sizeof(out->siteId)) ||
        !copy_json_str(root, JTAG_SOURCE, out->source, sizeof(out->source)) ||
        !copy_json_str(root, "updated_at", out->updatedAt,
                       sizeof(out->updatedAt))) {
        set_err(err, errLen, "bad policy metadata");
        return SWITCHD_ERR_INVAL;
    }

    ports = json_object_get(root, JTAG_PORTS);
    if (ports == NULL || !json_is_array(ports)) {
        set_err(err, errLen, "policy missing ports array");
        return SWITCHD_ERR_INVAL;
    }

    json_array_foreach(ports, index, item) {
        SwitchPortPolicy port;
        const char *policy;
        JsonObj *field;
        json_int_t portId;

        if (!json_is_object(item)) {
            set_err(err, errLen, "bad port policy entry");
            return SWITCHD_ERR_INVAL;
        }

        memset(&port, 0, sizeof(port));

        field = json_object_get(item, "port");
        if (field == NULL || !json_is_integer(field)) {
            field = json_object_get(item, JTAG_PORT_ID);
        }

        if (field == NULL || !json_is_integer(field)) {
            set_err(err, errLen, "port policy missing port");
            return SWITCHD_ERR_INVAL;
        }

        portId = json_integer_value(field);
        if (portId <= 0 || portId > SWITCHD_MAX_PORTS) {
            set_err(err, errLen, "port policy has invalid port");
            return SWITCHD_ERR_INVAL;
        }

        if (out->ports[portId - 1].present) {
            set_err(err, errLen, "duplicate port policy");
            return SWITCHD_ERR_INVAL;
        }

        port.port = (uint32_t)portId;
        port.present = true;

        if (!copy_json_str(item, "role", port.role, sizeof(port.role)) ||
            !copy_json_str(item, JSON_KEY_NODE_ID, port.nodeId,
                           sizeof(port.nodeId)) ||
            !copy_json_str(item, "class", port.klass, sizeof(port.klass))) {
            set_err(err, errLen, "bad port policy strings");
            return SWITCHD_ERR_INVAL;
        }

        field = json_object_get(item, "policy");
        if (field == NULL || !json_is_string(field)) {
            set_err(err, errLen, "port policy missing policy");
            return SWITCHD_ERR_INVAL;
        }

        policy = json_string_value(field);
        port.policy = policy_type_from_str(policy);
        if (port.policy == SWITCH_PORT_POLICY_UNKNOWN) {
            set_err(err, errLen, "unknown port policy");
            return SWITCHD_ERR_INVAL;
        }

        out->ports[port.port - 1] = port;
    }

    return SWITCHD_OK;
}

static int parse_policy_body(const char *path,
                             const char *body,
                             size_t bodyLen,
                             time_t fileMtime,
                             SwitchPolicy *out,
                             char *err,
                             size_t errLen) {

    JsonErrObj jerr;
    JsonObj *root;
    char hash[SWITCHD_SHA256_LEN];
    int ret;

    if (body == NULL || bodyLen == 0) {
        set_err(err, errLen, "empty policy body");
        return SWITCHD_ERR_INVAL;
    }

    memset(&jerr, 0, sizeof(jerr));

    root = json_loadb(body, bodyLen, 0, &jerr);
    if (root == NULL) {
        set_err(err, errLen, "invalid policy json");
        return SWITCHD_ERR_INVAL;
    }

    policy_hash_body(body, bodyLen, hash, sizeof(hash));

    ret = parse_policy(root,
                       path,
                       hash,
                       fileMtime,
                       out,
                       err,
                       errLen);

    json_decref(root);
    return ret;
}

int policy_load(SwitchdContext *ctx) {

    SwitchPolicy loaded;
    char *body;
    size_t bodyLen;
    time_t mtime;
    char err[SWITCHD_OP_DETAIL_LEN];
    int ret;

    if (ctx == NULL) {
        return SWITCHD_ERR_INVAL;
    }

    policy_clear(&ctx->policy);
    snprintf(ctx->policy.path, sizeof(ctx->policy.path), "%s",
             ctx->config.policyPath);

    body = NULL;
    bodyLen = 0;
    mtime = 0;

    ret = read_policy_file(ctx->config.policyPath, &body, &bodyLen, &mtime);
    if (ret != SWITCHD_OK) {
        ctx->policy.state = SWITCH_POLICY_STATE_MISSING;
        ctx->policy.checkedAt = time(NULL);
        snprintf(ctx->policy.error, sizeof(ctx->policy.error), "missing");
        usys_log_warn("switchd: no policy file at %s", ctx->config.policyPath);
        return ret;
    }

    memset(err, 0, sizeof(err));

    ret = parse_policy_body(ctx->config.policyPath,
                            body,
                            bodyLen,
                            mtime,
                            &loaded,
                            err,
                            sizeof(err));

    free(body);

    if (ret != SWITCHD_OK) {
        ctx->policy.state = SWITCH_POLICY_STATE_INVALID;
        ctx->policy.checkedAt = time(NULL);
        snprintf(ctx->policy.error, sizeof(ctx->policy.error), "%s", err);
        usys_log_error("switchd: policy validation failed: %s", err);
        return ret;
    }

    ctx->policy = loaded;

    usys_log_info("switchd: loaded policy path=%s state=%s hash=%s",
                  ctx->config.policyPath,
                  policy_state_to_str(ctx->policy.state),
                  ctx->policy.hash);

    return SWITCHD_OK;
}

int policy_reload_if_changed(SwitchdContext *ctx) {

    struct stat st;
    SwitchPolicy loaded;
    char *body;
    size_t bodyLen;
    char err[SWITCHD_OP_DETAIL_LEN];
    int ret;

    if (ctx == NULL) {
        return SWITCHD_ERR_INVAL;
    }

    ctx->policy.checkedAt = time(NULL);

    if (stat(ctx->config.policyPath, &st) != 0) {
        if (ctx->policy.state == SWITCH_POLICY_STATE_LOADED) {
            snprintf(ctx->policy.error,
                     sizeof(ctx->policy.error),
                     "policy_file_missing");
            return SWITCHD_OK;
        }

        ctx->policy.state = SWITCH_POLICY_STATE_MISSING;
        snprintf(ctx->policy.path,
                 sizeof(ctx->policy.path),
                 "%s",
                 ctx->config.policyPath);
        snprintf(ctx->policy.error, sizeof(ctx->policy.error), "missing");
        return SWITCHD_ERR_NOTFOUND;
    }

    if (ctx->policy.state == SWITCH_POLICY_STATE_LOADED &&
        ctx->policy.fileMtime == st.st_mtime) {
        return SWITCHD_OK;
    }

    body = NULL;
    bodyLen = 0;

    ret = read_policy_file(ctx->config.policyPath, &body, &bodyLen, NULL);
    if (ret != SWITCHD_OK) {
        snprintf(ctx->policy.error,
                 sizeof(ctx->policy.error),
                 "policy_read_failed");
        return ret;
    }

    memset(err, 0, sizeof(err));

    ret = parse_policy_body(ctx->config.policyPath,
                            body,
                            bodyLen,
                            st.st_mtime,
                            &loaded,
                            err,
                            sizeof(err));

    free(body);

    if (ret != SWITCHD_OK) {
        if (ctx->policy.state == SWITCH_POLICY_STATE_LOADED) {
            snprintf(ctx->policy.error,
                     sizeof(ctx->policy.error),
                     "new_policy_invalid");
            usys_log_error("switchd: new policy invalid, keeping old: %s",
                           err);
            return SWITCHD_OK;
        }

        ctx->policy.state = SWITCH_POLICY_STATE_INVALID;
        snprintf(ctx->policy.error, sizeof(ctx->policy.error), "%s", err);
        return ret;
    }

    ctx->policy = loaded;

    usys_log_info("switchd: reloaded policy path=%s hash=%s",
                  ctx->config.policyPath,
                  ctx->policy.hash);

    return SWITCHD_OK;
}

int policy_apply_body(SwitchdContext *ctx,
                      const char *body,
                      size_t bodyLen,
                      char *err,
                      size_t errLen) {

    SwitchPolicy loaded;
    int ret;

    if (ctx == NULL || body == NULL || bodyLen == 0) {
        set_err(err, errLen, "empty policy body");
        return SWITCHD_ERR_INVAL;
    }

    ret = parse_policy_body(ctx->config.policyPath,
                            body,
                            bodyLen,
                            time(NULL),
                            &loaded,
                            err,
                            errLen);
    if (ret != SWITCHD_OK) {
        return ret;
    }

    ret = write_policy_file(ctx->config.policyPath, body, bodyLen);
    if (ret != SWITCHD_OK) {
        set_err(err, errLen, "failed to write policy file");
        return ret;
    }

    ret = policy_reload_if_changed(ctx);
    if (ret != SWITCHD_OK) {
        ctx->policy = loaded;
    }

    usys_log_info("switchd: applied policy path=%s hash=%s",
                  ctx->config.policyPath,
                  ctx->policy.hash);

    return SWITCHD_OK;
}

const SwitchPortPolicy *policy_get_port(const SwitchdContext *ctx,
                                        uint32_t portId) {

    if (ctx == NULL || portId == 0 || portId > SWITCHD_MAX_PORTS) {
        return NULL;
    }

    if (!ctx->policy.ports[portId - 1].present) {
        return NULL;
    }

    return &ctx->policy.ports[portId - 1];
}

int policy_check_action(SwitchdContext *ctx,
                        uint32_t portId,
                        SwitchPolicyAction action,
                        const char *source,
                        char *err,
                        size_t errLen) {

    const SwitchPortPolicy *port;
    bool fromSiteController;

    if (ctx == NULL) {
        set_err(err, errLen, "no context");
        return SWITCHD_ERR_INVAL;
    }

    (void)policy_reload_if_changed(ctx);

    if (ctx->policy.state != SWITCH_POLICY_STATE_LOADED) {
        set_err(err, errLen, "port policy is not loaded");
        return SWITCHD_ERR_AUTH;
    }

    port = policy_get_port(ctx, portId);
    if (port == NULL) {
        set_err(err, errLen, "port is not in policy");
        return SWITCHD_ERR_AUTH;
    }

    fromSiteController = (source != NULL &&
                          strcmp(source,
                                 SWITCHD_POLICY_SOURCE_SITE_CONTROLLER) == 0);

    switch (port->policy) {
    case SWITCH_PORT_POLICY_FREE:
        return SWITCHD_OK;

    case SWITCH_PORT_POLICY_PROTECTED:
        if (fromSiteController) {
            return SWITCHD_OK;
        }
        set_err(err, errLen, "port is protected by site-controller");
        return SWITCHD_ERR_AUTH;

    case SWITCH_PORT_POLICY_NEVER_OFF_REMOTE:
        if (action_turns_off(action)) {
            set_err(err, errLen, "port cannot be disabled remotely");
            return SWITCHD_ERR_AUTH;
        }

        if (fromSiteController) {
            return SWITCHD_OK;
        }

        set_err(err, errLen, "port is protected by site-controller");
        return SWITCHD_ERR_AUTH;

    case SWITCH_PORT_POLICY_DISABLED:
        set_err(err, errLen, "port is disabled by policy");
        return SWITCHD_ERR_AUTH;

    case SWITCH_PORT_POLICY_UNKNOWN:
    default:
        set_err(err, errLen, "unknown port policy");
        return SWITCHD_ERR_AUTH;
    }
}

JsonObj *policy_serialize_overlay(const SwitchdContext *ctx,
                                  uint32_t portId) {

    const SwitchPortPolicy *port;
    JsonObj *json;

    json = json_object();
    if (ctx == NULL) {
        return json;
    }

    json_object_set_new(json,
                        JTAG_STATE,
                        json_string(policy_state_to_str(ctx->policy.state)));
    json_object_set_new(json, "hash", json_string(ctx->policy.hash));

    port = policy_get_port(ctx, portId);
    if (port == NULL) {
        json_object_set_new(json, "policy", json_string("unknown"));
        return json;
    }

    json_object_set_new(json, "role", json_string(port->role));
    json_object_set_new(json, JSON_KEY_NODE_ID, json_string(port->nodeId));
    json_object_set_new(json, "class", json_string(port->klass));
    json_object_set_new(json,
                        "policy",
                        json_string(policy_type_to_str(port->policy)));

    return json;
}

JsonObj *policy_serialize(const SwitchdContext *ctx) {

    JsonObj *root;
    JsonObj *ports;
    uint32_t i;

    root = json_object();
    ports = json_array();

    if (ctx == NULL) {
        json_object_set_new(root, JTAG_PORTS, ports);
        return root;
    }

    json_object_set_new(root,
                        JTAG_STATE,
                        json_string(policy_state_to_str(ctx->policy.state)));
    json_object_set_new(root, "site_id", json_string(ctx->policy.siteId));
    json_object_set_new(root, JTAG_SOURCE, json_string(ctx->policy.source));
    json_object_set_new(root, "updated_at", json_string(ctx->policy.updatedAt));
    json_object_set_new(root, JTAG_PATH, json_string(ctx->policy.path));
    json_object_set_new(root, JTAG_ERROR, json_string(ctx->policy.error));
    json_object_set_new(root, "hash", json_string(ctx->policy.hash));
    json_object_set_new(root,
                        "loaded_at",
                        json_integer((json_int_t)ctx->policy.loadedAt));
    json_object_set_new(root,
                        "checked_at",
                        json_integer((json_int_t)ctx->policy.checkedAt));
    json_object_set_new(root,
                        "file_mtime",
                        json_integer((json_int_t)ctx->policy.fileMtime));

    for (i = 0; i < SWITCHD_MAX_PORTS; i++) {
        SwitchPortPolicy *port;
        JsonObj *item;

        port = (SwitchPortPolicy *)&ctx->policy.ports[i];
        if (!port->present) {
            continue;
        }

        item = json_object();
        json_object_set_new(item, "port", json_integer(port->port));
        json_object_set_new(item, "role", json_string(port->role));
        json_object_set_new(item, JSON_KEY_NODE_ID, json_string(port->nodeId));
        json_object_set_new(item, "class", json_string(port->klass));
        json_object_set_new(item,
                            "policy",
                            json_string(policy_type_to_str(port->policy)));
        json_array_append_new(ports, item);
    }

    json_object_set_new(root, JTAG_PORTS, ports);
    return root;
}
