/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>
#include <errno.h>

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

USysMutex mutex;

static void free_session(ConfigSession *session) {

    if (session == NULL) return;

    for (int index; index < session->receviedCount; index++) {

        usys_free(session->apps[index]->name);
        usys_free(session->apps[index]->fileName);
        usys_free(session->apps[index]->data);
        usys_free(session->apps[index]->version);
    }

    usys_free(session);
}

static ConfigSession *create_new_session(SessionData *sd) {

	ConfigSession *session = NULL;

    session = (ConfigSession*) usys_calloc(1, sizeof(ConfigSession));
    if (session == NULL) {
        usys_log_error("Unable to allocate memory of size: %d",
                       sizeof(ConfigSession));
        return NULL;
    }

    session->apps[0]->name     = strdup(sd->app);
    session->apps[0]->version  = strdup(sd->version);
    session->apps[0]->fileName = strdup(sd->fileName);
    session->apps[0]->data     = strdup(sd->data);
    session->apps[0]->reason   = sd->reason;

    session->timestamp     = sd->timestamp;
    session->expectedCount = sd->fileCount;
    session->receviedCount = 0;

	return session;
}

static void update_config_session(Config *c, SessionData *sd) {

    ConfigSession *s = NULL;
    int index = 0;

    s     = (ConfigSession *) c->updateSession;
    index = s->receviedCount;

    s->apps[index]->name     = strdup(sd->app);
    s->apps[index]->version  = strdup(sd->version);
    s->apps[index]->fileName = strdup(sd->fileName);
    s->apps[index]->data     = strdup(sd->data);
    s->apps[index]->reason   = sd->reason;

    s->receviedCount++;
}

static bool create_config_staging_area(const char *app, int timestamp) {

    char path[PATH_MAX]     = {0};
    char realPath[PATH_MAX] = {0};
    char destPath[PATH_MAX] = {0};

    ssize_t len = 0;

    snprintf(path, sizeof(path), "%s/%s/active", DEF_CONFIG_DIR, app);
    snprintf(destPath, sizeof(destPath), "%s/%d/%s",
             CONFIG_TMP_PATH, timestamp, app);

    // Resolve the symlink to find the actual path
    len = readlink(path, realPath, sizeof(realPath) - 1);
    if (len == -1) {
        usys_log_error("symlink to actual path for app: %s Error: %s",
                       app, strerror(errno));
        return USYS_FALSE;
    }
    realPath[len] = '\0';

    /* create (if needed) and copy 'active' config to staging area */
    if (clone_dir(realPath, destPath, false) == 0) {
        usys_log_debug("Staging area successful for app: %s", app);
    } else {
        usys_log_error("Unable to create staging area for app: %s", app);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

static bool update_symlinks(char *appName, int timestamp) {

    char basePath[MAX_PATH]           = {0};
    char activePath[MAX_FILE_PATH]    = {0};
    char previousPath[MAX_FILE_PATH]  = {0};
    char newActivePath[MAX_FILE_PATH] = {0};

    char currentActivePath[MAX_PATH], currentPreviousPath[MAX_PATH];

    snprintf(basePath,      sizeof(basePath),      "%s/%s", DEF_CONFIG_DIR, appName);
    snprintf(activePath,    sizeof(activePath),    "%s/active", basePath);
    snprintf(previousPath,  sizeof(previousPath),  "%s/previous", basePath);
    snprintf(newActivePath, sizeof(newActivePath), "%s/archive/%d", basePath, timestamp);

    if (realpath(activePath, currentActivePath) == NULL) {
        usys_log_error("Error reading active symlink for app: %s", appName);
        return USYS_FALSE;
    }

    if (realpath(previousPath, currentPreviousPath) == NULL) {
        usys_log_error("Error reading previous symlink for app: %s", appName);
        return USYS_FALSE;
    }

    /* Now create symlink, if fails revert back */
    if (symlink(currentActivePath, previousPath) != 0 ||
        symlink(newActivePath,     activePath) != 0 ){

        usys_log_error("Unable to change the active/previous for app: %s", appName);

        symlink(activePath, currentActivePath);
        symlink(previousPath, currentPreviousPath);

        return USYS_FALSE;
    }

    usys_log_debug("Symlink successfully updated for app: %s", appName);

    return USYS_TRUE;
}

static bool process_config_session(Config *config) {

    bool ret = USYS_TRUE;
    int index;
    ConfigSession *s = NULL;
    
    char srcPath[MAX_PATH]  = {0};
    char destPath[MAX_PATH] = {0};

    s = (ConfigSession*)config->updateSession;

    for (index=0; index < s->receviedCount; index++) {

        snprintf(srcPath, sizeof(srcPath), "%s/%d/%s",
                 CONFIG_TMP_PATH, s->timestamp, s->apps[index]->name);
        snprintf(destPath, sizeof(destPath), "%s/%s/archive/%d",
                 DEF_CONFIG_DIR, s->apps[index]->name, s->timestamp);

        /* copy the config from staging area to the app's config */
        clone_dir(srcPath, destPath, false);

        /* update the active and previous symlink */
        update_symlinks(s->apps[index]->name, s->timestamp);

        /* remove the staging area */
        remove_dir(srcPath);

        /* send message to starter.d to restart the app */
        if (wc_send_app_restart_request(config, s->apps[index]->name) == USYS_FALSE) {
            usys_log_error("Unable to restart the app: %s",
                           s->apps[index]->name);
            ret = USYS_FALSE;
            continue;
        }

        usys_log_debug("App restart accepted by starter.d: %s",
                       s->apps[index]->name);
    }

	free_session(s);
	config->updateSession = NULL;

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

    if (sd->reason != CONFIG_ADD    &&
        sd->reason != CONFIG_DELETE &&
        sd->reason != CONFIG_UPDATE) return USYS_FALSE;

    // TO-DO check against stater.d if the app is valid
    // /v1/status/:space/:name

    if (config->updateSession) {
        if (sd->timestamp < ((ConfigSession *)config->updateSession)->timestamp) {
            usys_log_error("Received config %s with timestamp %ld. "
                           "expecting config timestamp %d",
                           sd->timestamp,
                           ((ConfigSession *)config->updateSession)->timestamp);
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
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

bool process_received_config(JsonObj *json, Config *config) {

	SessionData   *sd      = NULL;
	ConfigSession *session = NULL;

    session = (ConfigSession *)config->updateSession;

	/* Deserialize incoming message from ukama */
	if (!json_deserialize_session_data(json, &sd)) {
		return USYS_FALSE;
	}

    /* Check if the recevied session data is valid */
    if (!is_valid_session_data(sd, config)) {
        return USYS_FALSE;
    }

    /* No on-going update session going */
    if (!session) {
        pthread_mutex_lock(&mutex);
        session = create_new_session(sd);
        if (!session) {
            usys_log_error("failed to create new session.");
            pthread_mutex_unlock(&mutex);
            return USYS_FALSE;
        }
        config->updateSession = session;
        pthread_mutex_unlock(&mutex);
    }

    if (!decode_data(sd)) {
        usys_log_error("Unable to decode recevied data");
        return USYS_FALSE;
    }

    /* create config staging area for valid session */
    create_config_staging_area(sd->app,
                               ((ConfigSession *)config->updateSession)->timestamp);

    switch(sd->reason) {
    case CONFIG_DELETE:
		if (!remove_config_file_from_staging_area(sd)) {
			usys_log_error("Failed to remove config for %s app version %s",
                           sd->app, sd->version);
		}
        break;
    case CONFIG_ADD:
    case CONFIG_UPDATE:
		pthread_mutex_lock(&mutex);
		if (!create_config_file_in_staging_area(sd)) {
			usys_log_error("Failed to create config for %s app version %s",
                           sd->app, sd->version);
		}
		pthread_mutex_unlock(&mutex);
        break;
    default:
        return USYS_FALSE;
	}

	/* Update session */
	update_config_session(config, sd);
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
}
