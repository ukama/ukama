/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Callback functions for various endpoints and REST methods.
 */

#include <ulfius.h>
#include <string.h>
#include <jansson.h>

#include "callback.h"
#include "mesh.h"
#include "log.h"
#include "work.h"
#include "jserdes.h"

/* define in websocket.c */
extern void websocket_manager(const URequest *request, WSManager *manager,
			      void *data);
extern void websocket_incoming_message(const URequest *request,
				       WSManager *manager, WSMessage *message,
				       void *data);
extern void  websocket_onclose(const URequest *request, WSManager *manager,
			       void *data);

/*
 *
 */

static int is_valid_request(URequest *request) {

  return TRUE;
}

/*
 * Ulfius main callback function that simply calls the websocket manager
 * and closes
 */
int callback_websocket (const URequest *request, UResponse *response,
			void *data) {

  int ret;

  if ((ret = ulfius_set_websocket_response(response, NULL, NULL,
					   &websocket_manager,
					   data,
					   &websocket_incoming_message,
					   data,
					   &websocket_onclose,
					   data)) == U_OK) {
    ulfius_add_websocket_deflate_extension(response);
    return U_CALLBACK_CONTINUE;
  } else {
    return U_CALLBACK_ERROR;
  }
}

/*
 * callback_not_allowed -- 
 *
 */
int callback_not_allowed(const URequest *request, UResponse *response,
			 void *user_data) {
  
  ulfius_set_string_body_response(response, 403, "Operation not allowed\n");
  return U_CALLBACK_CONTINUE;
}

/*
 * callback_default_websocket -- default callback for no-match
 *
 */
int callback_default_websocket(const URequest *request, UResponse *response,
			       void *user_data) {

  ulfius_set_string_body_response(response, 404, "You are clearly high!\n");
  return U_CALLBACK_CONTINUE;
}

/*
 * callback_default -- default callback for no-match
 *
 */
int callback_default_webservice(const URequest *request, UResponse *response,
				void *data) {

  ulfius_set_string_body_response(response, 404, "You are clearly high!\n");
  return U_CALLBACK_CONTINUE;
}

/*
 * callback_webservice --
 *
 */
int callback_webservice(const URequest *request, UResponse *response,
			void *data) {

  json_t *jReq=NULL;
  Config *config;
  int ret;

  config = (Config *)data;
  
  /* For every incoming request, do following:
   *
   * 1. Sanity check.
   * 2. Convert request into JSON.
   * 3. Send request to Ukama proxy via websocket.
   * 4. Process websocket response.
   * 5. Wait for the response from server.
   * 6. Process response.
   * 7. Send response back to the client.
   * 8. Done
   */

  if (is_valid_request(request)==FALSE) {
    goto fail;
  }

  ret = serialize_forward_request(request, &jReq, config);
  if (ret == FALSE && jReq == NULL) {
    log_error("Failed to convert request to JSON");
    goto fail;
  }

  /* Add work for the websocket for transmission. */
  if (jReq != NULL) {

    /* No pre/post transmission func. This will block. */
    add_work_to_queue(get_transmit(), (Packet)jReq, NULL, 0, NULL, 0);
  }

 fail:
  /* Send response back to the callee */
  return U_CALLBACK_CONTINUE;
}
