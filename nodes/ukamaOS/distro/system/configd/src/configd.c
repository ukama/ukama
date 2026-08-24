/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <sys/types.h>
#include <sys/stat.h>
#include <stdio.h>
#include <unistd.h>
#include <errno.h>
#include <libgen.h>

#include "configd.h"
#include "base64.h"
#include "util.h"
#include "web_client.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"
#include "usys_services.h"

USysMutex mutex = PTHREAD_MUTEX_INITIALIZER;

static void config_status_set(Config *config,
                              ConfigApplyState state,
                              const char *requestId,
                              bool replaceRequestId) {

    if (!config) return;

    pthread_mutex_lock(&mutex);
    config->applyState = state;
    if (replaceRequestId) {
        snprintf(config->requestId,
                 sizeof(config->requestId),
                 "%s",
                 requestId ? requestId : "");
    }
    pthread_mutex_unlock(&mutex);
}

void config_status_snapshot(Config *config,
                            ConfigApplyState *state,
                            char *requestId,
                            size_t requestIdSize) {

    if (!config || !state || !requestId || requestIdSize == 0) return;

    pthread_mutex_lock(&mutex);
    *state = config->applyState;
    snprintf(requestId, requestIdSize, "%s", config->requestId);
    pthread_mutex_unlock(&mutex);
}

static void free_session(ConfigSession *session) {

    if (session == NULL) return;

    for (int index = 0; index < session->appCount; index++) {

        usys_free(session->apps[index].name);
        usys_free(session->apps[index].fileName);
        usys_free(session->apps[index].data);
        usys_free(session->apps[index].version);
    }

    usys_free(session->requestId);
    usys_free(session);
}

void config_session_clear(Config *config) {

    ConfigSession *session;

    if (!config) return;

    pthread_mutex_lock(&mutex);
    session = (ConfigSession *)config->updateSession;
    config->updateSession = NULL;
    pthread_mutex_unlock(&mutex);

    free_session(session);
}

static ConfigSession *create_new_session(SessionData *sd) {

    ConfigSession *session;

    session = (ConfigSession *)usys_calloc(1, sizeof(ConfigSession));
    if (session == NULL) {
        usys_log_error("Unable to allocate memory of size: %zu",
                       sizeof(ConfigSession));
        return NULL;
    }

    session->requestId = sd->requestId ? strdup(sd->requestId) : NULL;
    if (sd->requestId && !session->requestId) {
        usys_free(session);
        return NULL;
    }
    session->timestamp = sd->timestamp;
    session->expectedCount = sd->fileCount;

    return session;
}

static int session_app_index(ConfigSession *session, const char *app) {

    int index;

    if (!session || !app) return -1;

    for (index = 0; index < session->appCount; index++) {
        if (session->apps[index].name &&
            strcmp(session->apps[index].name, app) == 0) {
            return index;
        }
    }

    return -1;
}

static bool session_add_app(ConfigSession *session, const char *app) {

    char *name;

    if (!session || !app) return USYS_FALSE;
    if (session_app_index(session, app) >= 0) return USYS_TRUE;
    if (session->appCount >= MAX_APPS) return USYS_FALSE;

    name = strdup(app);
    if (!name) return USYS_FALSE;

    session->apps[session->appCount].name = name;
    session->appCount++;
    return USYS_TRUE;
}

static bool create_config_staging_area(const char *app, int timestamp) {

    char path[PATH_MAX]     = {0};
    char linkPath[128]      = {0};
    char realPath[PATH_MAX] = {0};
    char destPath[PATH_MAX] = {0};
    char resolvedLinkPath[PATH_MAX] = {0};

    ssize_t len = 0;

    snprintf(path, sizeof(path), "%s/%s/active", DEF_CONFIG_DIR, app);
    snprintf(destPath, sizeof(destPath), "%s/%d/%s", CONFIG_TMP_PATH, timestamp, app);

    /* Resolve the symlink to find the actual path */
    len = readlink(path, linkPath, sizeof(linkPath) - 1);
    if (len == -1) {
        usys_log_error("Error reading symlink for app: %s: %s", app, strerror(errno));
        return USYS_FALSE;
    }
    linkPath[len] = '\0';

    /* Check if the symlink is relative or absolute */
    if (linkPath[0] == '/') {
        snprintf(resolvedLinkPath, sizeof(resolvedLinkPath), "%s", linkPath);
    } else {
        char pathCopy[PATH_MAX] = {0};
        snprintf(pathCopy, sizeof(pathCopy), "%s", path);
        char *dir = dirname(pathCopy);

        snprintf(resolvedLinkPath, sizeof(resolvedLinkPath), "%s/%s", dir, linkPath);

        if (realpath(resolvedLinkPath, realPath) == NULL) {
            usys_log_error("Error resolving real path for: %s. Error: %s",
                           resolvedLinkPath, strerror(errno));
            return USYS_FALSE;
        }

        snprintf(resolvedLinkPath, sizeof(resolvedLinkPath), "%s", realPath);
    }

    /* create (if needed) and copy 'active' config to staging area */
    if (clone_dir(resolvedLinkPath, destPath, false) == 0) {
        usys_log_debug("Staging area successful for app: %s", app);
    } else {
        usys_log_error("Unable to create staging area for app: %s", app);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

static bool update_symlinks(const char *appName, int timestamp) {

    char basePath[MAX_PATH]           = {0};
    char activePath[MAX_FILE_PATH]    = {0};
    char previousPath[MAX_FILE_PATH]  = {0};
    char newActivePath[MAX_FILE_PATH] = {0};

    char currentActivePath[MAX_PATH] = {0}, currentPreviousPath[MAX_PATH] = {0};

    snprintf(basePath,      sizeof(basePath),      "%s/%s", DEF_CONFIG_DIR, appName);
    snprintf(activePath,    sizeof(activePath),    "%s/active", basePath);
    snprintf(previousPath,  sizeof(previousPath),  "%s/backup", basePath);
    snprintf(newActivePath, sizeof(newActivePath), "%s/archive/%d", basePath, timestamp);

    if (realpath(activePath, currentActivePath) == NULL) {
        usys_log_error("Error reading active symlink for app: %s. Error: %s",
                       appName, strerror(errno));
        return USYS_FALSE;
    }

    if (realpath(previousPath, currentPreviousPath) == NULL) {
        usys_log_error("Error reading previous symlink for app: %s. Error: %s",
                       appName, strerror(errno));
        return USYS_FALSE;
    }

    if (unlink(previousPath) != 0 && errno != ENOENT) {
        usys_log_error("Failed to remove old previous symlink for app: %s. Error: %s",
                       appName, strerror(errno));
        return USYS_FALSE;
    }

    if (symlink(currentActivePath, previousPath) != 0) {
        usys_log_error("Unable to create new previous symlink for app: %s. Error: %s",
                       appName, strerror(errno));
        return USYS_FALSE;
    }

    if (unlink(activePath) != 0 && errno != ENOENT) {
        usys_log_error("Failed to remove old active symlink for app: %s. Error: %s",
                       appName, strerror(errno));

        /* Rollback: restore the old previous symlink */
        unlink(previousPath);
        symlink(currentPreviousPath, previousPath);

        return USYS_FALSE;
    }

    if (symlink(newActivePath, activePath) != 0) {
        usys_log_error("Unable to create new active symlink for app: %s. Error: %s",
                       appName, strerror(errno));

        /* Rollback: restore the old symlinks */
        unlink(previousPath);
        symlink(currentPreviousPath, previousPath);
        symlink(currentActivePath, activePath);

        return USYS_FALSE;
    }

    usys_log_debug("Symlink successfully updated for app: %s", appName);

    return USYS_TRUE;
}

static bool process_config_session(Config *config) {

    bool ret = USYS_TRUE;
    int index;
    ConfigSession *s;
    
    char srcPath[MAX_PATH]  = {0};
    char destPath[MAX_PATH] = {0};

    s = (ConfigSession *)config->updateSession;

    for (index = 0; index < s->appCount; index++) {

        snprintf(srcPath, sizeof(srcPath), "%s/%d/%s",
                 CONFIG_TMP_PATH, s->timestamp, s->apps[index].name);
        snprintf(destPath, sizeof(destPath), "%s/%s/archive/%d",
                 DEF_CONFIG_DIR, s->apps[index].name, s->timestamp);

        /* copy the config from staging area to the app's config */
        if (clone_dir(srcPath, destPath, false) != 0) {
            usys_log_error("Unable to archive config for app: %s",
                           s->apps[index].name);
            ret = USYS_FALSE;
            continue;
        }

        /* update the active and previous symlink */
        if (!update_symlinks(s->apps[index].name, s->timestamp)) {
            ret = USYS_FALSE;
            continue;
        }

        /* remove the staging area */
        remove_dir(srcPath);

        /* send message to starter.d to restart the app */
        if (wc_send_app_restart_request(config,
                                        s->apps[index].name) == USYS_FALSE) {
            usys_log_error("Unable to restart the app: %s",
                           s->apps[index].name);
            ret = USYS_FALSE;
            continue;
        }

        usys_log_debug("App restart completed by starter.d: %s",
                       s->apps[index].name);
    }

    config_status_set(config,
                      ret ? CONFIG_APPLY_APPLIED : CONFIG_APPLY_FAILED,
                      s->requestId,
                      true);

    config_session_clear(config);

    return ret;
}

static bool is_valid_session_data(SessionData *sd, Config *config) {

    if (sd == NULL)           return USYS_FALSE;
    if (sd->timestamp <= 0)   return USYS_FALSE;
    if (sd->fileCount <= 0)   return USYS_FALSE;
    if (sd->app == NULL)      return USYS_FALSE;
    if (sd->fileName == NULL) return USYS_FALSE;
    if (sd->version == NULL)  return USYS_FALSE;
    if (sd->data == NULL)     return USYS_FALSE;
    if (sd->requestId &&
        strlen(sd->requestId) >= CONFIG_REQUEST_ID_LEN) return USYS_FALSE;

    if (sd->reason != CONFIG_ADD    &&
        sd->reason != CONFIG_DELETE &&
        sd->reason != CONFIG_UPDATE) return USYS_FALSE;

    if (wc_is_app_valid(config, sd->app) == USYS_FALSE) {
        return USYS_FALSE;
    }

    if (config->updateSession) {
        ConfigSession *session;

        session = (ConfigSession *)config->updateSession;
        if (sd->timestamp != session->timestamp) {
            usys_log_error("Received config with timestamp %d; "
                           "expecting timestamp %d",
                           sd->timestamp,
                           session->timestamp);
            return USYS_FALSE;
        }

        if (sd->fileCount != session->expectedCount) {
            usys_log_error("Config file_count does not match active session");
            return USYS_FALSE;
        }

        if ((session->requestId && !sd->requestId) ||
            (!session->requestId && sd->requestId) ||
            (session->requestId && sd->requestId &&
             strcmp(session->requestId, sd->requestId) != 0)) {
            usys_log_error("Config requestId does not match active session");
            return USYS_FALSE;
        }

        if (session->receviedCount >= session->expectedCount) {
            usys_log_error("Config session already received all files");
            return USYS_FALSE;
        }
    }

    return USYS_TRUE;
}

static bool decode_data(SessionData *sd) {

    int len;
    char *jc = NULL;

    if (!sd->data) return USYS_TRUE;

    len = usys_strlen(sd->data);
    usys_log_debug("Config base64 [%d bytes] received is %s", len, sd->data);

    jc = usys_calloc(sizeof(char), len);
    if (jc == NULL) {
        usys_log_error("Memory exhausted for decoding request. Size: %d", len);
        return USYS_FALSE;
    }

    base64_decode(jc, sd->data);
    usys_free(sd->data);
    sd->data = jc;
    usys_log_debug("Config text received\n:  %s", sd->data);

    if (!is_valid_json(sd->data)) {
        usys_free(sd->data);
        sd->data = NULL;
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

bool process_received_config(JsonObj *json, Config *config) {

    SessionData *sd;
    ConfigSession *session;
    bool firstForApp;

    if (!json || !config) return USYS_FALSE;

    sd = NULL;

    session = (ConfigSession *)config->updateSession;

    /* Deserialize incoming message from ukama */
    if (!json_deserialize_session_data(json, &sd)) {
        if (!config->updateSession) {
            config_status_set(config, CONFIG_APPLY_FAILED, NULL, true);
        }
        return USYS_FALSE;
    }

    /* Check if the recevied session data is valid */
    if (!is_valid_session_data(sd, config)) {
        if (!config->updateSession) {
            config_status_set(config,
                              CONFIG_APPLY_FAILED,
                              sd->requestId,
                              true);
        }
        free_session_data(sd);
        return USYS_FALSE;
    }

    /* No on-going update session going */
    if (!session) {
        pthread_mutex_lock(&mutex);
        session = create_new_session(sd);
        if (!session) {
            usys_log_error("failed to create new session.");
            pthread_mutex_unlock(&mutex);
            config_status_set(config,
                              CONFIG_APPLY_FAILED,
                              sd->requestId,
                              true);
            free_session_data(sd);
            return USYS_FALSE;
        }
        config->updateSession = session;
        pthread_mutex_unlock(&mutex);

        config_status_set(config,
                          CONFIG_APPLY_IN_PROGRESS,
                          session->requestId,
                          true);
    }

    if (!decode_data(sd)) {
        usys_log_error("Unable to decode recevied data");
        config_status_set(config,
                          CONFIG_APPLY_FAILED,
                          NULL,
                          false);
        config_session_clear(config);
        free_session_data(sd);
        return USYS_FALSE;
    }

    firstForApp = session_app_index(session, sd->app) < 0;

    /* create one staging tree per app in this session */
    if (firstForApp &&
        !create_config_staging_area(sd->app, session->timestamp)) {
        usys_log_error("Unable to create staging area for app");
        config_status_set(config,
                          CONFIG_APPLY_FAILED,
                          NULL,
                          false);
        config_session_clear(config);
        free_session_data(sd);
        return USYS_FALSE;
    }

    switch (sd->reason) {
    case CONFIG_DELETE:
        if (!remove_config_file_from_staging_area(sd)) {
            usys_log_error("Failed to remove config for %s app version %s",
                           sd->app, sd->version);
            config_status_set(config,
                              CONFIG_APPLY_FAILED,
                              NULL,
                              false);
            config_session_clear(config);
            free_session_data(sd);
            return USYS_FALSE;
        }
        break;
    case CONFIG_ADD:
    case CONFIG_UPDATE:
        pthread_mutex_lock(&mutex);
        if (!create_config_file_in_staging_area(sd)) {
            usys_log_error("Failed to create config for %s app version %s",
                           sd->app, sd->version);
            pthread_mutex_unlock(&mutex);
            config_status_set(config,
                              CONFIG_APPLY_FAILED,
                              NULL,
                              false);
            config_session_clear(config);
            free_session_data(sd);
            return USYS_FALSE;
        }
        pthread_mutex_unlock(&mutex);
        break;
    default:
        config_status_set(config,
                          CONFIG_APPLY_FAILED,
                          NULL,
                          false);
        config_session_clear(config);
        free_session_data(sd);
        return USYS_FALSE;
    }

    /* Update session */
    if (!session_add_app(session, sd->app)) {
        config_status_set(config,
                          CONFIG_APPLY_FAILED,
                          NULL,
                          false);
        config_session_clear(config);
        free_session_data(sd);
        return USYS_FALSE;
    }
    session->receviedCount++;
    free_session_data(sd);

    /* if this was the last data, process the session */
    if (session->expectedCount == session->receviedCount) {
        return process_config_session(config);
    }

    usys_log_debug("Received %d files and expected %d configs. Waiting for %d",
                   session->receviedCount,
                   session->expectedCount,
                   (session->expectedCount - session->receviedCount));

    return USYS_TRUE;
}

void free_session_data(SessionData *s) {

    if (s == NULL) return;

    usys_free(s->fileName);
    usys_free(s->app);
    usys_free(s->version);
    usys_free(s->data);
    usys_free(s->requestId);

    usys_free(s);
}
