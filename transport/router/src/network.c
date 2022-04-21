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
#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>

#include "router.h"
#include "callback.h"
#include "log.h"

/*
 * init_sockaddr --
 *
 */
static int init_sockaddr(struct sockaddr_in *name, char *hostName, char *port) {

  struct hostent *hostInfo=NULL;

  if (hostName==NULL || port==NULL) return FALSE;

  hostInfo = gethostbyname(hostName);
  if (hostInfo == NULL) {
    log_error("Unable to resolve hostName: %s", hostName);
    return FALSE;
  }
  
  name->sin_family = AF_INET;
  name->sin_port   = htons(atoi(port));
  name->sin_addr   = *(struct in_addr *)hostInfo->h_addr;

  return TRUE;
}

/* 
 * init_frameworks -- initializa ulfius framework and register various
 *                    callbacks.
 *
 */
static int init_frameworks(Config *config, struct _u_instance *webInst) {
  
  struct sockaddr_in sockAddr;
#if 0
  if (init_sockaddr(&sockAddr, config->hostName, config->port)) {
    log_error("Unable to init sockaddr for hostName %s port: %s",
	      config->hostName, config->port);
    return FALSE;
  }
#endif

  if (ulfius_init_instance(webInst, atoi(config->port), NULL, NULL)
      != U_OK) {
    log_error("Error initializing instance for hostName: %s port %s",
	      config->hostName, config->port);
    return FALSE;
  }

  /* Set few params. */
  webInst->max_post_body_size = 1024;

  return TRUE;
}

/*
 * setup_endpoints --
 *
 */
static void setup_endpoints(Router *router, struct _u_instance *instance) {

  /* GET, POST, DELETE /route */
  ulfius_add_endpoint_by_val(instance, "GET", EP_ROUTE, NULL, 0,
                             &callback_get_route, router);
  ulfius_add_endpoint_by_val(instance, "POST", EP_ROUTE, NULL, 0,
                             &callback_post_route, router);
  ulfius_add_endpoint_by_val(instance, "DELETE", EP_ROUTE, NULL, 0,
			     &callback_delete_route, router);
  /* GET /stats */
  ulfius_add_endpoint_by_val(instance, "GET", EP_STATS, NULL, 0,
                             &callback_get_stats, router);

  /* POST /service */
  ulfius_add_endpoint_by_val(instance, "POST", EP_SERVICE, NULL, 0,
                             &callback_post_service, router);

  /* default endpoint. */
  ulfius_set_default_endpoint(instance, &callback_default, router);
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

  log_debug("Webservice succesfully started");
  
  return TRUE;
}

/*
 * start_web_service --
 *
 */
int start_web_service(Router *router, struct _u_instance *webInst) {

  if (router == NULL)               return FALSE;
  if (router->config == NULL)       return FALSE;
  if (router->config->hostName == NULL ||
      router->config->port == NULL) return FALSE;

  /* Initialize the webservices framework. */
  if (init_frameworks(router->config, webInst) != TRUE) {
    log_error("Error initializing webservice framework");
    return FALSE;
  }

  /* setup endpoints and methods callback. */
  setup_endpoints(router, webInst);

  /* open connection for webservices */
  if (!start_framework(webInst)) {
    log_error("Failed to start webservices for hostName:%s port:%s",
	      router->config->hostName, router->config->port);
    return FALSE;
  }

  log_debug("Webservice on hostName:%s port:%s started.",
	    router->config->hostName, router->config->port);
  return TRUE;
}
