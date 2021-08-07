/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Mesh.d - L7-websocket based forward/reversed proxy
 *
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <getopt.h>
#include <ulfius.h>

#include "mesh.h"
#include "config.h"

#define VERSION "0.0.1"

/* Defined in network.c */
extern int start_web_services(Config *config, UInst *clientInst);
extern int start_websocket_server(Config *config, UInst *serverInst);
extern int start_websocket_client(Config *config,
				  struct _websocket_client_handler *handler);

/*
 * usage -- Usage options for the Mesh.d
 *
 *
 */

void usage() {

  printf("Usage: mesh.d [options] \n");
  printf("Options:\n");
  printf("--h, --help                         Help menu.\n");
  printf("--s, --secure                       Enable SSL/TLS \n");
  printf("--c, --config                       Config file.\n");
  printf("--m, --mode  <Server | Client>      Am I client or server?\n");
  printf("--l, --level <ERROR | DEBUG | INFO> Log level for the process.\n");
  printf("--V, --version                      Version.\n");
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

int main (int argc, char *argv[]) {

  int mode=MODE_SERVER, secure=FALSE;
  char *configFile=NULL;
  char *debug=DEF_LOG_LEVEL;
  Config *config=NULL;

  struct _u_instance serverInst;
  struct _u_instance clientInst;
  struct _websocket_client_handler websocket_client_handler = {NULL, NULL};

  /* Prase command line args. */
  while (TRUE) {

    int opt = 0;
    int opdidx = 0;

    static struct option long_options[] = {
      { "mode",      required_argument, 0, 'm'},
      { "secure",    no_argument,       0, 's'},
      { "config",    required_argument, 0, 'c'},
      { "level",     required_argument, 0, 'l'},
      { "help",      no_argument,       0, 'h'},
      { "version",   no_argument,       0, 'V'},
      { 0,           0,                 0,  0}
    };

    opt = getopt_long(argc, argv, "l:c:m::", long_options, &opdidx);
    if (opt == -1) {
      break;
    }

    switch (opt) {
    case 'h':
      usage();
      exit(0);
      break;

    case 'c':
      configFile = optarg;
      break;

    case 'l':
      debug = optarg;
      set_log_level(debug);
      break;

    case 's':
      secure=TRUE;
      break;

    case 'm':
      if (strcasecmp(optarg, "client")==0) {
	mode = MODE_CLIENT;
      } else if (strcasecmp(optarg, "server")==0) {
	mode = MODE_SERVER;
      } else {
	fprintf(stderr, "Invalid mode: %s\n", optarg);
	exit(1);
      }
      break;

    case 'V':
      fprintf(stdout, "Mesh.d - Version: %s\n", VERSION);
      exit(0);

    default:
      usage();
      exit(0);
    }
  } /* while */

  if (argc == 1 || configFile == NULL) {
    fprintf(stderr, "Missing required parameters\n");
    usage();
    exit(1);
  }

  config = (Config *)calloc(1, sizeof(Config));
  if (!config) {
    log_error("Memory allocation failure: %d", sizeof(Config));
    exit(1);
  }

  /* Step-1: read config file. */
  if (process_config_file(mode, secure, configFile, config) != TRUE) {
    fprintf(stderr, "Error parsing config file: %s. Exiting ... \n",
	    configFile);
    exit(1);
  }

  print_config(config);

  /* Step-2a: start webservice for local client. */
  if (start_web_services(config, &clientInst) != TRUE) {
    log_error("Webservice failed to setup for clients. Exiting.");
    exit(1);
  }

  /* Step-2b: if server, setup all endpoints, cb and run websocket. Wait. */
  if (config->mode == MODE_SERVER) {
    if (start_websocket_server(config, &serverInst) != TRUE) {
      log_error("Websocket failed to setup for server. Exiting...");
      exit(1);
    }
  }

  /* Step-2c: setup websocket client. */
  if (config->mode == MODE_CLIENT) {
    if (start_websocket_client(config, &websocket_client_handler) != TRUE) {
      log_error("Websocket failed to setup for client. Exiting...");
      exit(1);
    }
  }

  /* Wait here for ever. XXX */

  log_debug("Mesh.d running ...");

  getchar(); /* For now. */

  log_debug("UnMesh.d and Goodbye ... ");

  if (config->mode == MODE_CLIENT) {
    ulfius_websocket_client_connection_close(&websocket_client_handler);
  }

  if (config->mode == MODE_SERVER) {
    ulfius_stop_framework(&serverInst);
    ulfius_clean_instance(&serverInst);
  }

  clear_config(config);
  free(config);

  return 1;
}
