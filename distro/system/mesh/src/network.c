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

#include "callback.h"
#include "mesh.h"
#include "websocket.h"

/*
 * init_framework -- initializa ulfius framework.
 *
 */

int init_framework(struct _u_instance *inst, int port) {

  if (ulfius_init_instance(inst, port, NULL, NULL) != U_OK) {
    log_error("Error initializing instance for websocket remote port %d", port);
    return FALSE;
  }

  /* Set few params. */
  inst->max_post_body_size = 1024;
  
  return TRUE;
}
    
/*
 * setup_websocket_endpoints -- 
 *
 */

void setup_websocket_endpoints(Config *config, struct _u_instance *instance) {

  /* Endpoint list declaration. */
  ulfius_add_endpoint_by_val(instance, "GET", PREFIX_WEBSOCKET, NULL, 0,
			     &callback_websocket, config);
  ulfius_add_endpoint_by_val(instance, "POST", PREFIX_WEBSOCKET, NULL, 0,
			     &callback_websocket, config);
  ulfius_add_endpoint_by_val(instance, "PUT", PREFIX_WEBSOCKET, NULL, 0,
			     &callback_not_allowed, config);
  ulfius_add_endpoint_by_val(instance, "DELETE", PREFIX_WEBSOCKET, NULL, 0,
			     &callback_not_allowed, config);
  
  /* default endpoint. */
  ulfius_set_default_endpoint(instance, &callback_default, config);
}

/* 
 * start_framework --
 *
 */

int start_framework(Config *config, struct _u_instance *instance) {

  int ret;
  
  /* open HTTPS/HTTP connection. */
  if (config->secure) {
    ret = ulfius_start_secure_framework(instance, config->keyFile,
					config->certFile);
  } else {
    ret = ulfius_start_framework(instance);
  }

  if (ret != U_OK) {
    log_error("Error starting the webservice.");
    
    /* clean up. */
    ulfius_stop_framework(instance); /* don't think need this. XXX */
    ulfius_clean_instance(instance);
    
    return FALSE;
  }

  log_debug("Websocket succesfully started.");
  
  return TRUE;
}

/*
 * start_websocket_server -- start websocket server on the server port.
 *
 */

int start_websocket_server(Config *cfg, struct _u_instance *serverInst) {

  /* Initialize the admin and client webservices framework. */
  if (init_framework(serverInst, atoi(cfg->remoteAccept)) != TRUE) {
    log_error("Error initializing webservice framework");
    return FALSE;
  }

  /* setup endpoints and methods callback. */
  setup_websocket_endpoints(cfg, serverInst);
  
  /* open connection for both admin and client webservices */
  if (start_framework(cfg, serverInst)==FALSE) {
    log_error("Failed to start websocket at remote port %s", cfg->remoteAccept);
    return FALSE;
  }
  
  log_debug("Websocket on remote port %s: started.", cfg->remoteAccept);

  return TRUE;
}

/*
 * start_websocket_client -- Connect with remote server using websockets.
 *
 */

int start_websocket_client(Config *config,
			   struct _websocket_client_handler *handler) {

  struct _u_request request;
  struct _u_response response;

  if (ulfius_init_request(&request) != U_OK) {
    goto failure;
  }

  if (ulfius_init_response(&response) != U_OK) {
    goto failure;
  }

  /* Setup websocket request. */
  if (ulfius_set_websocket_request(&request, config->remoteConnect,
				   NULL, NULL) == U_OK) {
    /* Setup request parameters */
    u_map_put(request.map_header, "User-Agent",
	      MESH_CLIENT_AGENT "/" MESH_CLIENT_VERSION);
    ulfius_add_websocket_client_deflate_extension(handler);
    request.check_server_certificate = FALSE;

    /* Open websocket connection to remote host. */
    if (ulfius_open_websocket_client_connection(&request,
	                &websocket_manager_cb, (void *)NULL,
			&websocket_incoming_message_cb, (void *)NULL,
			&websocket_onclose_cb, (void*)NULL,		
		        handler, &response) == U_OK) {
      /* Success. The websocket will now run as seperate thread as cb */
      return TRUE;
    } else {
      log_error("Unable to open websocket connect to: %s",
		config->remoteConnect);
      goto failure;
    }
  } else {
    log_error("Error initializing the websocket client request");
    goto failure;
  }

 failure:
  ulfius_clean_request(&request);
  ulfius_clean_response(&response);

  return FALSE;
}
