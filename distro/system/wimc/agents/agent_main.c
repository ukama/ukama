/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Agent -- Basic agent framework 
 *
 */

#include <stdlib.h>
#include <string.h>
#include <sqlite3.h>
#include <getopt.h>
#include <ulfius.h>

#include "log.h"
#include "wimc.h"
#include "agent.h"

#include "agent/network.h"

#define TRUE 1
#define FALSE 0
#define VERSION "0.0.1"
#define DEF_LOG_LEVEL "TRACE"
/*
 * usage -- Usage options for the Agent.
 *
 *
 */

void usage() {

  printf("Agent: WIMC.d Agent to speak with service provider for contents\n");
  printf("Supported methods: Test, Chunk\n");
  printf("Usage: Agent [options] \n");
  printf("Options:\n");
  printf("--h, --help                         This help menu. \n");
  printf("--w, --wimc                         WIMC URL - http://host:port/\n");
  printf("--p, --port                         Client listening port. \n");
  printf("--m, --method                       Tx Method <Test | Chunk>.\n");
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

int set_method(char *method) {

  if (!strcmp(method, "TEST")) {
    return METHOD_TEST;
  } else if (!strcmp(method, "CHUNK")) {
    return METHOD_CHUNK;
  }

  log_debug("Setting method to default: TEST. Invalid input: %s", method);
  return METHOD_TEST;
}

/*
 * Agent main.
 */

int main(int argc, char **argv) {

  int id=0;
  long code;
  char *wimcURL = NULL, *port=NULL;
  char *debug = DEF_LOG_LEVEL;
  MethodType method = METHOD_TEST;
  struct _u_instance inst;
    
  if (argc == 1) {
    usage();
    exit(1);
  }

  /* Prase command line args. */
  while (TRUE) {
    
    int opt = 0;
    int opdidx = 0;

    static struct option long_options[] = {
      { "wimc",    required_argument, 0, 'w'},
      { "port",    required_argument, 0, 'p'},
      { "method",  required_argument, 0, 'm'},
      { "level",   required_argument, 0, 'l'},
      { "help",    no_argument,       0, 'h'},
      { "version", no_argument,       0, 'v'},
      { 0,         0,                 0,  0}
    };

    opt = getopt_long(argc, argv, "w:p:m:l:hV:", long_options, &opdidx);
    if (opt == -1) {
      break;
    }

    switch (opt) {
    case 'w':
      wimcURL = optarg;
      break;
      
    case 'p':
      port = optarg;
      break;

    case 'm':
      method = set_method(optarg);
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
      fprintf(stdout, "Agent - Version: %s\n", VERSION);
      exit(0);

    default:
      usage();
      exit(0);
    }
  } /* while */

  /* Step-1. setup the EP with respective CB and run webservice. */
  if (start_web_service(port, &method, &inst) != TRUE) {
    log_error("Failed to start webservice. Exiting.");
    exit(0);
  }
  
  /* Step2. register itself on the WIMC. */
  code = communicate_with_wimc(REQ_REG, wimcURL, port, method, &id);
  if (!code || code == 400) { /* Failure to register. */
    log_error("Failed to register to wimc.d. Exiting");
    goto cleanup;
  }
  
  /* Step3. trigger the modules using CB */
  

  /* Exit. unregister. */
  //  code = communicate_with_wimc(REQ_UNREG, wimcURL, port, NULL, &id);

  getchar(); /*. For now. xxx */
  
  log_debug("Ukama.\n");

 cleanup:
  ulfius_stop_framework(&inst);
  ulfius_clean_instance(&inst);
  
  return 1;
}
