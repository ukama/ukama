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

#include "wimc.h"
#include "log.h"
#include "callback.h"

/* 
 * init_frameworks -- initializa ulfius framework and register various
 *                   callbacks.
 *
 */

int init_frameworks(struct _u_instance *adminInst,
		    struct _u_instance *clientInst,
		    int adminPort, int clientPort) {

  if (ulfius_init_instance(adminInst, adminPort, NULL, NULL) != U_OK) {
    log_error("Error initializing instance for admin port %d", adminPort);
    return FALSE;
  }

  if (ulfius_init_instance(clientInst, clientPort, NULL, NULL) != U_OK) {
    log_error("Error initializing instance for client port %d", clientPort);
    ulfius_clean_instance(adminInst);
    return FALSE;
  }

  /* Set few params. */
  u_map_put(adminInst->default_headers, "Access-Control-Allow-Origin", "*");
  adminInst->max_post_body_size = 1024;
  clientInst->max_post_body_size = 1024;
  
  return TRUE;
}
    
/*
 * setup_admin_endpoints -- 
 *
 */

void setup_admin_endpoints(WimcCfg *cfg, struct _u_instance *instance) {

  /* Endpoint decelrations. We have two endpoints:
   * 1. /admin (acting on 'containers' table)
   *     GET    - query the db.
   *     POST   - add new entry to db.
   *     PUT    - update an existing entry in db.
   *     DELETE - remove existing entry in db.
   */
  ulfius_add_endpoint_by_val(instance, "GET", WIMC_EP_ADMIN, NULL, 0,
                             &callback_get_container, cfg);
  ulfius_add_endpoint_by_val(instance, "POST", WIMC_EP_ADMIN, NULL, 0,
                             &callback_post_container, cfg);
  ulfius_add_endpoint_by_val(instance, "PUT", WIMC_EP_ADMIN, NULL, 0,
                             &callback_put_container, cfg);
  ulfius_add_endpoint_by_val(instance, "DELETE", WIMC_EP_ADMIN, NULL, 0,
                             &callback_delete_container, cfg);
  
  /* 2. /stats:
   *     GET  - get various WIMC.d internal stats.
   */
  ulfius_add_endpoint_by_val(instance, "GET", WIMC_EP_STATS, NULL, 0,
                             &callback_get_stats, cfg);

  /* default endpoint. */
  ulfius_set_default_endpoint(instance, &callback_default, cfg);
}

/*
 * setup_client_endpoints --
 *
 */

void setup_client_endpoints(WimcCfg *cfg, struct _u_instance *instance) {

  /* Endpoint decelrations. 
   * 1. /containers (acting on 'containers' table)
   *     GET    - query the db.
   */
  ulfius_add_endpoint_by_val(instance, "GET", WIMC_EP_CLIENT, NULL, 0,
                             &callback_get_container, cfg);
  ulfius_add_endpoint_by_val(instance, "POST", WIMC_EP_CLIENT, NULL, 0,
                             &callback_not_allowed, cfg);
  ulfius_add_endpoint_by_val(instance, "PUT", WIMC_EP_CLIENT, NULL, 0,
                             &callback_not_allowed, cfg);
  ulfius_add_endpoint_by_val(instance, "DELETE", WIMC_EP_CLIENT, NULL, 0,
                             &callback_not_allowed, cfg);
  
  /* default endpoint. */
  ulfius_set_default_endpoint(instance, &callback_default, cfg);
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
    ulfius_stop_framework(instance); /* don't think need this. XXX */
    ulfius_clean_instance(instance);
    
    return FALSE;
  }

  log_debug("Webservice succesfully started.");
  
  return TRUE;
}

/*
 * start_web_services -- start accepting REST clients on 127.0.0.1:port
 *
 */

int start_web_services(WimcCfg *cfg, struct _u_instance *adminInst,
		       struct _u_instance *clientInst) {
  
  /* Initialize the admin and client webservices framework. */
  if (init_frameworks(adminInst, clientInst, atoi(cfg->adminPort),
		      atoi(cfg->clientPort)) != TRUE) {
    log_error("Error initializing webservice framework");
    return FALSE;
  }
  
  /* setup endpoints and methods callback. */
  setup_admin_endpoints(cfg, adminInst);
  setup_client_endpoints(cfg, clientInst);
  
  /* open connection for both admin and client webservices */
  if (start_framework(adminInst)) {
    if (!start_framework(clientInst)) {

      /* if failsure, stop admin instance. */
      ulfius_stop_framework(adminInst); 
      ulfius_clean_instance(adminInst);
      
      return FALSE;
    }
  } else {
    log_error("Failed to start webservices for admin:%s and client: %s",
	      cfg->adminPort, cfg->clientPort);
    return FALSE;
  }
  
  log_debug("Webservice on admin port %s client port: %s started.",
	    cfg->adminPort, cfg->clientPort);

  return TRUE;
}
