/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * WIMC.d --
 *
 */

#include <stdlib.h>
#include <string.h>
#include <sqlite3.h>
#include <getopt.h>

#include "log.h"
#include "ulfius.h"
#include "wimc.h"
#include "agent.h"

#define ENV_CLIENT_PORT "WIMC_CLIENT_PORT"
#define ENV_ADMIN_PORT  "WMIN_ADMIN_PORT"
#define ENV_DB_FILE     "WIMC_DB"
#define ENV_UKAMA_CLOUD "WIMC_UKAMA_CLOUD"
#define DEF_LOG_LEVEL   "TRACE"

#define TRUE 1
#define FALSE 0
#define VERSION "0.0.1"

#define WIMC_FLAG_CREATE_DB 1

/*
 * usage -- Usage options for the WIMC.d.
 *
 *
 */

void usage() {

  printf("WIMC.d: Service to answer \"Where Is My Content?\"\n");
  printf("Usage: wimc.d [options] \n");
  printf("Options:\n");
  printf("--h, --help                         This help menu. \n");
  printf("--u, --url                          Cloud URL - http://host:port/\n");
  printf("--d, --dbFile                       Full path for db file. \n");
  printf("--p, --port                         Client listening port. \n");
  printf("--P, --aPort                        Admin listneing port. \n");
  printf("--a, --agent                        Max. number of agents. \n");
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

/* Ensure all the required parameters in CFG are valid. Otherwise exit. */
void is_valid_config(WimcCfg *cfg) {

  /* Make sure we have ports and dbFile defined. Environment variable 
   * takes precedence over the passed arguments.
   */
  
  if (cfg->clientPort == NULL) {
    log_error("Error: client port must be defined, via env or arg");
    usage();
    exit(0);
  }
  
  if (cfg->dbFile == NULL) {
    log_error("Error: dbFile must be defined, via env or arg");
    usage();
    exit(0);
  }
  
  if (cfg->adminPort == NULL) {
    log_error("Error: admin port must be defined, via env or arg");
    usage();
    exit(0);
  }
  
  /* Make sure the URL is in right format. */
  if (is_valid_url(cfg->cloud) != TRUE) {
    log_error("Invalid URL: %s\n", cfg->cloud);
    usage();
    exit(0);
  }

  log_debug("Using following configuration: \n clientPort: %s \n dbFile: %s \n adminPort: %s \n Cloud-URL: %s", cfg->clientPort, cfg->dbFile, cfg->adminPort,
	    cfg->cloud);
}

/*
 * wimc.d needs three main environment variables. Environment variables always
 * override the passed arguments. 
 *
 * WIMC_DB -- db to store the data
 * WIMC_CLINET_PORT -- Port for handling clients request.
 * WIMC_ADMIN_PORT  -- Admin port.
 */

int main (int argc, char **argv) {

  sqlite3 *db=NULL;
  struct _u_instance adminInst, clientInst;
  WimcCfg *cfg=NULL;
  Agent *agents=NULL;

  char *debug = DEF_LOG_LEVEL;

  cfg = (WimcCfg *)calloc(sizeof(WimcCfg), 1);
  if (cfg == NULL) {
    exit(1);
  }
  
  /* Check if the environment variables are set. */
  cfg->clientPort = getenv(ENV_CLIENT_PORT);
  cfg->adminPort  = getenv(ENV_ADMIN_PORT);
  cfg->dbFile     = getenv(ENV_DB_FILE);
  cfg->cloud      = getenv(ENV_UKAMA_CLOUD);
  cfg->maxAgents  = 1;
  
  if (argc == 1 ||
      (cfg->clientPort == NULL && cfg->adminPort == NULL
       && cfg->dbFile == NULL)) {
    fprintf(stderr, "Missing required parameters\n");
    usage();
    exit(1);
  }

  /* Prase command line args. */
  while (TRUE) {
    
    int opt = 0;
    int opdidx = 0;

    static struct option long_options[] = {
      { "url",       required_argument, 0, 'u'},
      { "dbFile",    required_argument, 0, 'd'},
      { "port",      required_argument, 0, 'p'},
      { "aPort",     required_argument, 0, 'P'},
      { "level",     required_argument, 0, 'l'},
      { "agents",    required_argument, 0, 'a'},
      { "help",      no_argument,       0, 'h'},
      { "version",   no_argument,       0, 'v'},
      { 0,           0,                 0,  0}
    };

    opt = getopt_long(argc, argv, "u:d:p:a:d:l:hV:", long_options, &opdidx);
    if (opt == -1) {
      break;
    }

    switch (opt) {
    case 'c':
      if (cfg->cloud == NULL) {
	cfg->cloud = optarg;
      }
      break;
      
    case 'd':
      if (cfg->dbFile == NULL) { /* Ignore otherwise. */
	cfg->dbFile = optarg;
      }
      break;
      
    case 'p':
      if (cfg->clientPort == NULL) { /* ignore this option otherwise. */
	cfg->clientPort = optarg;
      }
      break;

    case 'P':
      if (cfg->adminPort == NULL){ /* ignore this otherwise. */
	cfg->adminPort = optarg;
      }
      break;

    case 'a':
      cfg->maxAgents = atoi(optarg);
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
      fprintf(stdout, "WIMC.d - Version: %s\n", VERSION);
      exit(0);

    default:
      usage();
      exit(0);
    }
  } /* while */

  agents = (Agent *)calloc(cfg->maxAgents, sizeof(Agent));
  if (agents == NULL) {
    log_error("Error allocating memory: %d", (int)sizeof(Agent)*cfg->maxAgents);
    exit(1);
  }

  cfg->agents = agents;
  is_valid_config(cfg);
  
  /* Steps are as follows:
   * 1. Check if db exits, otherwise create one.
   * 2. Open port to accept REST calls. Register cb for:
   *    GET  - query the db.
   *    POST - add new entry to db.
   *    PUT  - update an existing entry in db.
   *    DELETE - remove existing entry in db.
   */

  /* Example usage:
   * curl GET localhost:port:/container/nginx:4.1.1
   * Response: nginx:4.1.1:/path/to/bundle-1/:Available:1
   * curl GET localhost:port:/container/nginx
   * Response: nginx:latest:/other/path/to/bundle-2/:Available:1
   */
  
  /* Step-1 */
  db = open_db(cfg->dbFile, WIMC_FLAG_CREATE_DB);

  if (db == NULL) {
    log_error("Error creating db at: %s. Exiting", cfg->dbFile);
    exit(0);
  }

  /* Step-2, setup all endpoints, cb and run webservice at ports */
  if (start_web_services(cfg, agents, &adminInst, &clientInst) != TRUE) {
    log_error("Webservice failed to setup for admin/clients. Exiting");
    exit(0);
  }

  getchar(); /* For now. XXX */

  log_debug("End World!\n");
  
  ulfius_stop_framework(&adminInst);
  ulfius_clean_instance(&adminInst);

  ulfius_stop_framework(&clientInst);
  ulfius_clean_instance(&clientInst);
  
  return 1;
}
