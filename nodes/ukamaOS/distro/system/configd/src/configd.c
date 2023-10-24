/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "configd.h"
#include "web_client.h"
#include "usys_error.h"
#include "usys_log.h"

void free_config_data(ConfigData *c) {
	if (c) {
		if (c->fileName); free(c->fileName);
		if (c->app); free(c->app);
		if (c->version); free(c->version);
		if (c->data);free(c->data);
	}
}

int configd_process_incoming_config(const char *service,
                                         JsonObj *json, Config *config){

	int statusCode=-1;
	ConfigData *config=NULL;

	/* Deserialize incoming message from noded */
	if (!json_deserialize_config(json, &config)) {
		return STATUS_NOK;
	}

	/* Validate the jsosn data */
	 if (is_valid_json(config->data)) {
		 statusCode = create_config(config);
	 }

	 free_config_data(config);

	return statusCode;
}


int is_valid_commit(ConfigData *c) {

	if ((c->timestamp > timestamp) && (usys_strcmp(c->version, version))) {
		usys_log_debug("Receiving new config %s with timestamp %ld. Discarding old config %s with timestamp %d\n", c->version, c->timestamp, version, timestamp);
		free_config_state_handler();
		create_ne_config_state_handler();
	}

	if (!(usys_strcmp(c->version, version)) && (c->timestamp == timestamp) ) {
		return 1;
	}
	return 0;
}

int configd_process_complete(const char *service,
						JsonObj *json, Config *config){
	int statusCode=-1;
	ConfigData *cd = NULL;

	/* Deserialize incoming message from noded */
	if (!json_deserialize_config(json, &cd)) {
		return STATUS_NOK;
	}

	/* Validate the jsosn data */
	 if ( is_valid_json(cd->data) && is_valid_commit(cd) ) {
		 statusCode = store_config(cd);
	 }


	 /* Trigger updates */


	 free_config_data(cd);

	return statusCode;
}

