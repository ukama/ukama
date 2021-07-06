/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Network related stuff based on ulfius framework for the Agents.
 */

#include <ulfius.h>
#include <stdlib.h>

#include "log.h"
#include "callback.h"
#include "agent.h"
#include "wimc.h"
#include "agent/network.h"
#include "agent/jserdes.h"

#define AGENT_EP "/container"
#define STAT_EP  "/stats"

/*
 * setup_endpoints -- setup various EP for HTTP methods.
 *
 */
void setup_endpoints(MethodType *method, struct _u_instance *instance) {

  /* Endpoint decelrations. 
   * 1. /container 
   *     GET    - query an existing on-going fetch session with provider.
   *     POST   - initiate a new session.
   *     PUT    - update an existing session (inteval and CB URL).
   *     DELETE - cancel an on-going session with provider.
   */
  ulfius_add_endpoint_by_val(instance, "GET", AGENT_EP, NULL, 0,
                             &agent_callback_get, NULL);
  ulfius_add_endpoint_by_val(instance, "POST", AGENT_EP, NULL, 0,
                             &agent_callback_post, method);
  ulfius_add_endpoint_by_val(instance, "PUT", AGENT_EP, NULL, 0,
                             &agent_callback_put, NULL);
  ulfius_add_endpoint_by_val(instance, "DELETE", AGENT_EP, NULL, 0,
                             &agent_callback_delete, NULL);
  
  /* 2. /stats
   *     GET  - get Agent various internal stats.
   */
  ulfius_add_endpoint_by_val(instance, "GET", STAT_EP, NULL, 0,
                             &agent_callback_stats, NULL);

  /* default endpoint. */
  ulfius_set_default_endpoint(instance, &agent_callback_default, NULL);
}

/* 
 * start_framework --
 *
 */
int start_framework(struct _u_instance *instance) {

  int ret;
  
  /* open HTTP connection. */
  ret = ulfius_start_framework(instance);

  if (ret != U_OK) {
    log_error("Error starting the webservice.");
    
    /* clean up. */
    ulfius_stop_framework(instance); 
    ulfius_clean_instance(instance);
    
    return FALSE;
  }

  return TRUE;
}

/*
 * start_web_services -- start accepting REST clients on 127.0.0.1:port
 *
 */
int start_web_service(char *port, MethodType *method,
		      struct _u_instance *inst) {

  if (ulfius_init_instance(inst, atoi(port), NULL, NULL) != U_OK) {
    log_error("Error initializing instance for port %s", port);
    return FALSE;
  }

  /* Set few params. */
  u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");
  inst->max_post_body_size = 1024;

  /* Setup endpoints. */
  setup_endpoints(method, inst);
  
  /* open connection for WIMC call-to-action */
  if (!start_framework(inst)) {
    log_error("Failed to start webservices at :%s", port);
    return FALSE;
  }
  
  log_debug("Webservice on port %s started.", port);

  return TRUE;
}
