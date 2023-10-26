/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "configd.h"

#include "util.h"
#include "web_client.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"

void free_apps(AppState **apps, uint32_t count) {
	for (uint32_t i = 0; i < count; i++) {
		if (apps[i]) free(apps[i]);
	}
}

void free_config_data(ConfigData *c) {
	if (c) {
		if (c->fileName) free(c->fileName);
		if (c->app) free(c->app);
		if (c->version) free(c->version);
		if (c->data) free(c->data);
	}
}

/* cleans a session for update */
void clean_session(ConfigSession *session) {

	if (session) {
		if (session->version) free(session->version);
		if (session->apps) free_apps (session->apps, session->count);
		session->timestamp = 0;
		session->count = 0;
		free(session);
	}
}

/* creates a new session for update */
ConfigSession* create_new_update_session(ConfigData *cd) {
	ConfigSession *session = (ConfigSession*) usys_calloc(1, sizeof(ConfigSession));
	if (session) {
		session->timestamp = cd->timestamp;
		session->version = usys_strdup(cd->version);
		session->count = 0;

		if (prepare_for_new_config(cd) == 0) {
			usys_log_debug("New update config session created for commit %s and timestamp %ld", cd->version, cd->timestamp);
		} else {
			clean_session(session);
			session = NULL;
		}
	}

	return session;
}

/* Validate commit and creates a new session if required */
int is_valid_commit(Config* c , ConfigData *cd) {

	/* Discard config is older then current running config */
	ConfigData* rc = (ConfigData*) c->runningConfig;
	if (rc) {
		if (rc->timestamp > cd->timestamp) {
			usys_log_debug("Received config %s with timestamp %ld. expecting config newer than %s with timestamp %d\n", cd->version, cd->timestamp, rc->version, rc->timestamp);
			return 0;
		}
	}

	ConfigSession* s = (ConfigSession*) c->updateSession;
	if ((cd->timestamp != s->timestamp) || (usys_strcmp(cd->version, s->version))) {
		/* Newer config */
		if (cd->timestamp > s->timestamp) {
			usys_log_debug("Receiving new config %s with timestamp %ld. Discarding old config %s with timestamp %d\n", cd->version, cd->timestamp, s->version, s->timestamp);
			clean_session(c->updateSession);
			c->updateSession = create_new_update_session(cd);
			if (c->updateSession) {
				s = (ConfigSession*) c->updateSession;
			} else {
				return 0;
			}
		} else {
			/* Old rest request or wrog version */
			usys_log_error("Receiving config %s with timestamp %ld. expecting config %s with timestamp %d\n", cd->version, cd->timestamp, s->version, s->timestamp);
			return 0;
		}
	}

	if (!(usys_strcmp(cd->version, s->version)) && (cd->timestamp == s->timestamp) ) {
		AppState* as = (AppState*) usys_calloc(1, sizeof(AppState));
		if (as) {
			as->app = usys_strdup(cd->app);
			as->state = STATE_UPDATE_AVAILABLE;
			s->apps[s->count] = as;
			s->count++;
		} else {
			perror("Memory failure");
			return 0;
		}
		return 1;
	}
	return 0;
}

int process_config(JsonObj *json, Config *config) {
	int statusCode = STATUS_NOK;
	ConfigData *cd = NULL;

	/* Deserialize incoming message from cloud */
	if (!json_deserialize_config_data(json, &cd)) {
		return STATUS_NOK;
	}

	/* Validate the json data */
	if (!is_valid_json(cd->data)) {
		return STATUS_NOK;
	}

	/* get or create session */
	if (config) {
		if (!config->updateSession) {
			ConfigSession *session = create_new_update_session(cd);
			if (!session) {
				usys_log_error("failed to create update session.");
				return STATUS_NOK;
			}
			config->updateSession = session;

		}
	} else {
		usys_log_error("invalid config for web service.");
		return STATUS_NOK;
	}

	/* Validate the commit*/
	if (!is_valid_commit(config, cd)) {
		return STATUS_NOK;
	}

	if (cd->reason == CONFIG_DELETED){
		statusCode =  remove_config(cd);
		if (statusCode != STATUS_OK ) {
			usys_log_error("Failed to remove config for %s app version %s", cd->app, cd->version);
		}
	}
	else {
		statusCode =  create_config(cd);
		if (statusCode != STATUS_OK ) {
			usys_log_error("Failed to create config for %s app version %s", cd->app, cd->version);
		}

	}
	free_config_data(cd);

	return statusCode;
}

/* store incoming config file */
int configd_process_incoming_config(const char *service,
		JsonObj *json, Config *config){
	int statusCode = STATUS_NOK;
	statusCode = process_config(json, config);
	if (statusCode != STATUS_OK ) {
		usys_log_error("Failed to process config message.");
	}

	return statusCode;
}

int configd_process_complete(const char *service,
		JsonObj *json, Config *config){
	int statusCode = STATUS_NOK;
	ConfigSession* s = (ConfigSession*)config->updateSession;
	if (!s) {
		usys_log_error("No configs pushed. Failed to process config complete message.");
		goto cleanup;
	}

	statusCode = process_config(json, config);
	if (statusCode != STATUS_OK ) {
		usys_log_error("Failed to process config complete message.");
		goto cleanup;
	}

	/* Store config*/
	statusCode = store_config(s->version);
	if (statusCode != STATUS_OK ) {
		usys_log_error("Failed to store config %s", s->version);
		goto cleanup;
	}

	/* Trigger updates */
	statusCode = configd_trigger_update(s);

	/* Update running config */
	if (configd_read_running_config(&config->runningConfig)) {
		usys_log_error("Failed to update running config.");
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
		usys_log_debug("Triggering update for %s app to version %s ", s->apps[i]->app, s->version);

		/* send start message to restart app */
		statusCode = wc_send_restart_req(c, s->apps[i]->app);
		if (statusCode != STATUS_OK) {
			usys_log_error("Failed to exec app %s.", s->apps[i]->app);
			continue;
		}

	}

	return statusCode;
}

int configd_read_running_config(ConfigData **c) {
	int statusCode = STATUS_NOK;
	ConfigData *cd = NULL;

	/* Read file */
	char* file[512]={'\0'};
	usys_strcpy(file, CONFIG_RUNNING);
	usys_strcat(file, CONFIGD);
	/* Deserialize running config */
	if (!json_deserialize_running_config(file, &cd)) {
		usys_log_error("Failed to read running config %s", file);
		return STATUS_NOK;
	}
	/* clean */
	if (*c) {
		free_config_data(*c);
		*c = NULL;
	}

	/* Allocate */
	if (cd) {
		*c = cd;
	}

	return STATUS_OK;
}
