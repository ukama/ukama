/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Network related stuff based on ulfius framework.
 */

#include <ulfius.h>
#include <stdlib.h>
#include <string.h>

#include "lxce_callback.h"
#include "lxce_config.h"
#include "log.h"

#define PREFIX_WEBSERVICE "*"

/*
 * init_framework -- initializa ulfius framework
 *
 */
static int init_framework(UInst *inst, int port) {

  if (ulfius_init_instance(inst, port, NULL, NULL) != U_OK) {
    log_error("Error initializing instance for websocket remote port %d", port);
    return FALSE;
  }

  /* Set few params. */
  u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");

  return TRUE;
}

/*
 * setup_webservice_endpoints --
 *
 */
static void setup_webservice_endpoints(Config *config, UInst *instance) {

  /* Endpoint list declaration. */
  ulfius_add_endpoint_by_val(instance, "GET", PREFIX_WEBSERVICE, NULL, 0,
			     &callback_webservice, config);
  ulfius_add_endpoint_by_val(instance, "POST", PREFIX_WEBSERVICE, NULL, 0,
			     &callback_webservice, config);
  ulfius_add_endpoint_by_val(instance, "PUT", PREFIX_WEBSERVICE, NULL, 0,
			     &callback_webservice, config);
  ulfius_add_endpoint_by_val(instance, "DELETE", PREFIX_WEBSERVICE, NULL, 0,
			     &callback_not_allowed, config);

  /* default endpoint. */
  ulfius_set_default_endpoint(instance, &callback_default, config);
}

/*
 * start_framework --
 *
 */
static int start_framework(Config *config, UInst *instance) {

  int ret;

  ret = ulfius_start_framework(instance);
  if (ret != U_OK) {
    log_error("Error starting the webservice/websocket.");

    /* clean up. */
    ulfius_stop_framework(instance); /* don't think need this. XXX */
    ulfius_clean_instance(instance);

    return FALSE;
  }

  log_debug("Webservice sucessfully started.");

  return TRUE;
}

/*
 * start_web_services -- start accepting REST clients on 127.0.0.1:port
 *
 */
int start_web_services(Config *config, UInst *clientInst) {

  /* Initialize the admin and client webservices framework. */
  if (init_framework(clientInst, atoi(config->localAccept)) != TRUE){
    log_error("Error initializing webservice framework");
    return FALSE;
  }

  /* setup endpoints and methods callback. */
  setup_webservice_endpoints(config, clientInst);

  /* open connection for client webservices */
  if (!start_framework(config, clientInst)) {
    log_error("Failed to start webservices for client: %s",
	      config->localAccept);
    return FALSE;
  }

  log_debug("Webservice on client port: %s started.", config->localAccept);

  return TRUE;
}
