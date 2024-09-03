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

static bool is_valid_softlink(char *path) {

    struct stat pathStat;

    if (lstat(path, &pathStat) == -1) {
        usys_log_error("Unable to get info about config dir: %s. Error: %s",
                       path, strerror(errno));
        return USYS_FALSE;
    }

    if (S_ISLNK(pathStat.st_mode)) {

        char actualPath[MAX_PATH];

        if (realpath(path, actualPath) == NULL) {
            usys_log_error("Unable to resolve the actual path for: %s. Error: %s",
                           path, strerror(errno));
            return USYS_FALSE;
        }

        if (stat(actualPath, &pathStat) == -1) {
            usys_log_error("Unable to get stat about actual path: %s Error: %s",
                           actualPath, strerror(errno));
            return USYS_FALSE;
        }

        if (S_ISDIR(pathStat.st_mode)) {
            return USYS_TRUE;
        }
    }

    return USYS_FALSE;
}

static bool has_nonzero_metadata_file(const char *dirPath) {

    char filePath[PATH_MAX] = {0};
    struct stat fileStat;

    snprintf(filePath, sizeof(filePath), "%s/metadata.json", dirPath);

    if (stat(filePath, &fileStat) == -1) {
        if (errno == ENOENT) {
            usys_log_error("File metadata.json does not exist at: %s", dirPath);
        } else {
            usys_log_error("Unable to get stat about: %s Error: %s",
                           filePath, strerror(errno));
        }

        return USYS_FALSE;
    }

    if (fileStat.st_size > 0) {
        return USYS_TRUE;
    }

    usys_log_error("File metadata.json is empty at: %s", dirPath);

    return USYS_FALSE;
}

static void free_apps(AppState **apps, uint32_t count) {

	for (int i = 0; i < count; i++) {

		if (apps[i]) {
			usys_free(apps[i]->app);
			usys_free(apps[i]->fileName);
			usys_free(apps[i]);
		}
	}
}

void free_config_data(ConfigData *c) {

    if (c == NULL) return;

    usys_free(c->fileName);
    usys_free(c->app);
    usys_free(c->version);
    usys_free(c->data);
}

/* cleans a session for update */
void clean_session(ConfigSession *session) {

    char path[MAX_PATH] = {0};

    if (session == NULL) return;

    snprintf(path, sizeof(path), "%s/%s", CONFIG_TMP_PATH, session->version);
    remove_dir(path);

    if (session->version) usys_free(session->version);
    if (session->apps)    free_apps (session->apps, session->count);

    session->timestamp     = 0;
    session->count         = 0;
    session->expectedCount = 0;
    session->configdVer    = USYS_FALSE;

    usys_free(session);
}

/* creates a new session for update */
ConfigSession* create_new_update_session(ConfigData *cd) {

	ConfigSession *session = NULL;

    session = (ConfigSession*) usys_calloc(1, sizeof(ConfigSession));
    if (session == NULL) {
        usys_log_error("Unable to allocate memory of size: %d",
                       sizeof(ConfigSession));
        return NULL;
    }

    session->timestamp     = cd->timestamp;
    session->version       = usys_strdup(cd->version);
    session->expectedCount = cd->fileCount;
    session->count         = 0;
    session->configdVer    = USYS_FALSE;
    session->stored        = USYS_FALSE;
    /* Need to move this from here. Taking to long */
    //		if (prepare_for_new_config(cd) == 0) {
    //			usys_log_debug("New update config session created for commit %s and timestamp %ld", cd->version, cd->timestamp);
    //		} else {
    //			clean_session(session);
    //			session = NULL;
    //		}

	return session;
}

int prepare_copy_for_session(ConfigData *cd) {

	if (prepare_for_new_config(cd) == 0) {
		usys_log_debug("New update config session created for "
                       "commit %s and timestamp %ld",
                       cd->version,
                       cd->timestamp);
		return STATUS_OK;
	}

	return STATUS_NOK;
}

void update_session(Config* c, AppState* a) {

	ConfigSession *s = (ConfigSession *) c->updateSession;

	/* check if its a duplicate reception of file for session */
	for (int i = 0; i < s->count; i++) {
		if (!(usys_strcmp(a->app, s->apps[i]->app)) &&
            !(usys_strcmp(a->fileName, s->apps[i]->fileName))) {
			usys_log_debug("Received a duplicate config file %s for app %s.",
                           a->app, a->fileName);
			return;
		}
	}

	s->apps[s->count] = a;
	s->count++;

	/* version flag for configd/version.json */
	if (!(usys_strcmp(a->app, c->serviceName)) &&
        !(usys_strcmp(a->fileName, "version.json"))) {
		s->configdVer = USYS_TRUE;
	}
}

/* Validate commit and creates a new session if required */
int is_valid_commit(Config *c , ConfigData *cd, AppState **app) {

	int ret = USYS_FALSE;

	/* Discard config is older then current running config */
	ConfigData* rc = (ConfigData*) c->activeConfig;
	if (rc) {
		if (rc->timestamp > cd->timestamp) {
			usys_log_debug("Received config %s with timestamp %ld. "
                           "expecting config newer than %s with timestamp %d",
                           cd->version,
                           cd->timestamp,
                           rc->version,
                           rc->timestamp);
			goto response;
		}
	}

	pthread_mutex_lock(&mutex);
	ConfigSession* s = (ConfigSession*) c->updateSession;
	if ((cd->timestamp != s->timestamp) ||
        (usys_strcmp(cd->version, s->version))) {

		/* Newer config */
		if (cd->timestamp > s->timestamp) {
			usys_log_debug("Receiving new config %s with timestamp %ld. "
                           "Discarding old config %s with timestamp %d",
                           cd->version,
                           cd->timestamp,
                           s->version,
                           s->timestamp);

			clean_session(c->updateSession);
			c->updateSession = NULL;
			c->updateSession = create_new_update_session(cd);

			if (c->updateSession) {
				s = (ConfigSession*) c->updateSession;
				if (prepare_copy_for_session(cd) != STATUS_OK) {
					usys_log_error("Failed to prepare_copy for new session %s",
                                   cd->version);
					clean_session(c->updateSession);
					c->updateSession = NULL;
					goto response;
				}
			} else {
				goto response;
			}
		} else {
			/* Old rest request or wrog version */
			usys_log_error("Receiving config %s with timestamp %ld. "
                           "expecting config %s with timestamp %d",
                           cd->version,
                           cd->timestamp,
                           s->version,
                           s->timestamp);
			goto response;
		}
	}

	if (!(usys_strcmp(cd->version, s->version)) &&
        (cd->timestamp == s->timestamp) ) {

		AppState *as = (AppState*) usys_calloc(1, sizeof(AppState));
        if (as == NULL) {
            usys_log_error("Unable to allocate memory of size: %d",
                           sizeof(AppState));
            goto response;
        }

        as->app      = usys_strdup(cd->app);
        as->state    = STATE_UPDATE_AVAILABLE;
        as->fileName = usys_strdup(cd->fileName);
        *app = as;

		ret = USYS_TRUE;
	}

response:
	pthread_mutex_unlock(&mutex);
	return ret;
}

int process_config(JsonObj *json, Config *config) {

	int statusCode = STATUS_NOK;
	ConfigData *cd = NULL;
	ConfigSession *session = (ConfigSession*) config->updateSession;
	AppState *as = NULL;

	/* Deserialize incoming message from ukama */
	if (!json_deserialize_config_data(json, &cd)) {
		return STATUS_NOK;
	}

	/* get or create session */
	if (config) {
		if (!session) {
			pthread_mutex_lock(&mutex);
			session = create_new_update_session(cd);
			if (!session) {
				usys_log_error("failed to create update session.");
				pthread_mutex_unlock(&mutex);
				return STATUS_NOK;
			}
			config->updateSession = session;

			if (prepare_copy_for_session(cd) != STATUS_OK) {
				usys_log_error("Failed to prepare_copy for new session %s", cd->version);
				clean_session(config->updateSession);
				config->updateSession = NULL;
				pthread_mutex_unlock(&mutex);
				return STATUS_NOK;
			}
			pthread_mutex_unlock(&mutex);

		}
	} else {
		usys_log_error("invalid config for web service.");
		return STATUS_NOK;
	}

	if (cd->data) {
		int len = usys_strlen(cd->data);
		usys_log_debug("Config base64 [%d bytes] received is %s", len, cd->data);
		char *jc = usys_calloc(sizeof(char), len);
		if (jc) {
			base64_decode(jc, cd->data);
			usys_free(cd->data);
			cd->data=jc;
			usys_log_debug("Config text received is:\n  %s", cd->data);
		} else {
			usys_log_error("Memory exhausted for decoding request.");
			return STATUS_NOK;
		}

		/* Validate the json data */
		if (!is_valid_json(cd->data)) {
			return STATUS_NOK;
		}

	}

	/* Validate the commit*/
	if (!is_valid_commit(config, cd, &as)) {
		return STATUS_NOK;
	}

	if (cd->reason == CONFIG_DELETED){
		statusCode = remove_config(cd);
		if (statusCode != STATUS_OK ) {
			usys_log_error("Failed to remove config for %s app version %s",
                           cd->app, cd->version);
		}
	}
	else {
		pthread_mutex_lock(&mutex);
		statusCode =  create_config(cd);
		if (statusCode != STATUS_OK ) {
			usys_log_error("Failed to create config for %s app version %s",
                           cd->app, cd->version);
		}
		pthread_mutex_unlock(&mutex);
	}

	/* Update session */
	update_session(config, as);

	/* In case valid commit opened new update session */
	session = (ConfigSession*) config->updateSession;
	if (session->count == session->expectedCount) {
		if (session->configdVer) {
			usys_log_debug("Received all expected %d configs", session->expectedCount);
			statusCode = configd_process_complete(config);
		} else {
			usys_log_error("Received %d configs but version.json for "
                           "configd is missing", session->count);
			usys_log_error("Cleaning session.");
			clean_session(config->updateSession);
			config->updateSession = NULL;
			statusCode = STATUS_NOK;
		}
	} else {
		usys_log_debug("Received %d files and expected %d configs. Waiting for %d",
                       session->count,
                       session->expectedCount,
                       (session->expectedCount - session->count));
	}

	free_config_data(cd);

	return statusCode;
}

/* store incoming config file */
int configd_process_incoming_config(const char *service,
                                    JsonObj *json,
                                    Config *config){

	if ( process_config(json, config) != STATUS_OK ) {
		usys_log_error("Failed to process config message.");
        return STATUS_NOK;
	}

	return STATUS_OK;
}

int configd_process_complete(Config *config) {

	int statusCode = STATUS_NOK;
	ConfigSession* s = (ConfigSession*)config->updateSession;

	/* Store config */
	if (!(s->stored)) {
		s->stored = true;
		statusCode = store_config(s->version);
		if (statusCode != STATUS_OK ) {
			usys_log_error("Failed to store config %s", s->version);
			goto cleanup;
		}
	
		/* clean up empty dir in store */
//		char dir[512];
//		sprintf(dir,"%s/%s", CONFIG_TMP_PATH, s->version);
//		clean_empty_dir(dir);

		/* Trigger updates */
		statusCode = configd_trigger_update(config);

		/* Update active config */
		if (read_active_config((ConfigData**)&config->activeConfig)) {
			usys_log_error("Failed to update active config.");
		}
	}

cleanup:
	clean_session(config->updateSession);
	config->updateSession = NULL;

	return statusCode;
}

/* not monitoing anything app status for now */
int configd_trigger_update(Config* c) {

	int statusCode = STATUS_NOK;
	ConfigSession *s = (ConfigSession*)c->updateSession;

	for (int i = 0; i < s->count; i++) {
		usys_log_debug("Triggering update for %s app to version %s ",
                       s->apps[i]->app, s->version);

		if (usys_strcmp(s->apps[i]->app, c->serviceName)) {

			/* update runnig config */
			read_active_config(&(c->activeConfig));
		} else {
			/* send start message to restart app */
			statusCode = wc_send_restart_req(c, s->apps[i]->app);
			if (statusCode != STATUS_OK) {
				usys_log_error("Failed to exec app %s.", s->apps[i]->app);
				continue;
			}
		}
	}

	return statusCode;
}

int read_active_config(ConfigData **c) {

	ConfigData *cd = NULL;
    char path[MAX_PATH] = {0};
    char file[MAX_FILE_PATH] = {0};

    sprintf(path, "%s/%s/active", DEF_CONFIG_DIR, SERVICE_CONFIG);

    /* check if the path is valid - it is soft link to a valid
     * directory and have non-zero metadata.json file therein 
     */
    if (!is_valid_softlink(path) || !has_nonzero_metadata_file(path)) {
        return USYS_FALSE;
    }

    snprintf(file, sizeof(file), "%s/metadata.json", path);

	/* Deserialize 'active' config */
	if (!json_deserialize_active_config(file, &cd)) {
		usys_log_error("Failed to read active config %s", file);
		return STATUS_NOK;
	}

	/* clean and allocate */
	if (*c) {
		free_config_data(*c);
		*c = NULL;
	}

	if (cd) {
		*c = cd;
		usys_log_debug("Active config set to %s.", (*c)->version);
	}

	return STATUS_OK;
}
