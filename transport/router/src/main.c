/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Service router --
 *
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <getopt.h>
#include <ulfius.h>
#include <curl/curl.h>

#include "router.h"
#include "network.h"
#include "log.h"
#include "pattern.h"

#define DEF_LOG_LEVEL "TRACE"

#define VERSION "0.0.1"

/*
 * usage -- Usage options
 *
 *
 */
static void usage() {
  
  printf("srvc_router: Microservice message router\n");
  printf("Usage: srvc_router [options] \n");
  printf("Options:\n");
  printf("--h, --help                         This help menu. \n");
  printf("--H, --host                         Bind host address \n");
  printf("--p, --port                         Listneing port. \n");
  printf("--l, --level <ERROR | DEBUG | INFO> Log level for the process. \n");
  printf("--v, --version                      Version. \n");
}

/* Set the verbosity level for logs. */
void set_log_level(char *slevel) {
  
  int ilevel = LOG_TRACE;

  if (!strcmp(slevel, "DEBUG")) {
    ilevel = LOG_DEBUG;
  } else if (!strcmp(slevel, "INFO")) {
    ilevel = LOG_INFO;
  } else if (!strcmp(slevel, "ERROR")) {
    ilevel = LOG_ERROR;
  }

  log_set_level(ilevel);
}

void free_router(Router *router) {

  Service *sPtr=NULL, *next=NULL;

  if (router->services) {
    sPtr = router->services;

    while (sPtr) {
      next = sPtr->next;
      free_service(sPtr);
      sPtr = next;
    }
  }

  if (router->config) {
    free(router->config);
  }

  free(router);
  router=NULL;
}

/* srvc_router */

int main (int argc, char **argv) {

  Router *router=NULL;
  Config *config=NULL;
  char *debug=DEF_LOG_LEVEL;
  int opt, opdidx;
  struct _u_instance webInst;

  router = (Router *)calloc(1, sizeof(Router));
  if (router == NULL) {
    fprintf(stderr, "Error allocating memory of: %lu", sizeof(Router));
    exit(1);
  }

  router->config = (Config *)calloc(1, sizeof(Config));
  config = router->config;
  if (router->config == NULL) {
    fprintf(stderr, "Error allocating memory of: %lu", sizeof(Config));
    exit(1);
  }
  
  /* Prase command line args. */
  while (TRUE) {
    
    opt    = 0;
    opdidx = 0;

    static struct option long_options[] = {
      { "host",    required_argument, 0, 'H'},
      { "port",    required_argument, 0, 'p'},
      { "level",   required_argument, 0, 'l'},
      { "help",    no_argument,       0, 'h'},
      { "version", no_argument,       0, 'v'},
      { 0,         0,                 0,  0}
    };

    opt = getopt_long(argc, argv, "H:p:l:hv:", long_options, &opdidx);
    if (opt == -1) {
      break;
    }

    switch (opt) {
    case 'H':
      config->hostName = optarg;
      break;
      
    case 'p':
      config->port = optarg;
      break;

    case 'h':
      usage();
      exit(0);
      break;

    case 'l':
      debug = optarg;
      set_log_level(debug);
      break;

    case 'v':
      fprintf(stdout, "srvc_router - Version: %s\n", VERSION);
      exit(0);

    default:
      usage();
      exit(0);
    }
  } /* while */

  if (argc == 1 || config->hostName == NULL || config->port == NULL) {
    fprintf(stderr, "Missing required parameters\n");
    usage();
    exit(1);
  }

  if (start_web_service(router, &webInst) != TRUE) {
    log_error("Webservice failed to setup. Exiting");
    exit(0);
  }

  getchar(); /* For now. XXX */

  log_debug("Bye World!\n");
  
  ulfius_stop_framework(&webInst);
  ulfius_clean_instance(&webInst);
  curl_global_cleanup();

  free_router(router);
  
  return 1;
}
