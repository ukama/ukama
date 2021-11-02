/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * lxce.d - Ukama's light-weight container engine.
 *
 */

#include <string.h>
#include <strings.h>
#include <ulfius.h>
#include <getopt.h>
#include <sys/types.h>
#include <unistd.h>
#include <errno.h>
#include <stdlib.h>
#include <stdio.h>

#include "log.h"
#include "toml.h"
#include "lxce_config.h"
#include "manifest.h"
#include "cspace.h"
#include "csthreads.h"
#include "lxce_callback.h"

#define VERSION "0.0.1"

#define TRUE 1
#define FALSE 0

#define DEF_LOG_LEVEL     "TRACE"
#define DEF_CONFIG_FILE   "config.toml"
#define DEF_MANIFEST_FILE "manifest.json"

/* Default files for various contained setup */
#define DEF_BOOT_CONTAINED_FILE     "boot_contained.json"
#define DEF_SERVICE_CONTAINED_FILE  "service_contained.json"
#define DEF_SHUTDOWN_CONTAINED_FILE "shutdown_contained.json"

/* from lxce_network */
extern int start_web_services(Config *config, UInst *clientInst);

/*
 * callback functions declaration
 */

int callback_post_create_container (const struct _u_request *request,
				    struct _u_response *response,
				    void *user_data);
int callback_default (const struct _u_request *request,
		      struct _u_response *response, void *user_data);

/*
 * decode a u_map into a string
 */
char * print_map(const struct _u_map * map) {
  char * line, * to_return = NULL;
  const char **keys, * value;
  int len, i;
  if (map != NULL) {
    keys = u_map_enum_keys(map);
    for (i=0; keys[i] != NULL; i++) {
      value = u_map_get(map, keys[i]);
      len = snprintf(NULL, 0, "key is %s, value is %s", keys[i], value);
      line = o_malloc((len+1)*sizeof(char));
      snprintf(line, (len+1), "key is %s, value is %s", keys[i], value);
      if (to_return != NULL) {
        len = o_strlen(to_return) + o_strlen(line) + 1;
        to_return = o_realloc(to_return, (len+1)*sizeof(char));
        if (o_strlen(to_return) > 0) {
          strcat(to_return, "\n");
        }
      } else {
        to_return = o_malloc((o_strlen(line) + 1)*sizeof(char));
        to_return[0] = 0;
      }
      strcat(to_return, line);
      o_free(line);
    }
    return to_return;
  } else {
    return NULL;
  }
}

/* 
 * usage -- Usage options for the lxce.d
 *
 *
 */

void usage() {
  
  printf("Usage: lxce.d [options] \n");
  printf("Options:\n");
  printf("--h, --help                         Help menu.\n");
  printf("--c, --config                       Config file.\n");
  printf("--m, --manifest                     Manifest file.\n");
  printf("--l, --level <TRACE | DEBUG | INFO> Log level for the process.\n");
  printf("--V, --version                      Version.\n");
}

/*
 * set_log_level -- Set verbose level.
 *
 *
 */

/* Set the verbosity level for logs. */
void set_log_level(char *slevel) {
  
  int ilevel = LOG_TRACE;

  if (!strcmp(slevel, "TRACE")) {
    ilevel = LOG_TRACE;
  } else if (!strcmp(slevel, "DEBUG")) {
    ilevel = LOG_DEBUG;
  } else if (!strcmp(slevel, "INFO")) {
    ilevel = LOG_INFO;
  }
  
  log_set_level(ilevel);
}

int main(int argc, char **argv) {
  
  int i;
  char *debug = DEF_LOG_LEVEL;
  char *configFile = DEF_CONFIG_FILE;
  char *manifestFile = DEF_MANIFEST_FILE;
  Config *config = NULL;
  Manifest *manifest = NULL;
  CSpace *cSpaces, *cPtr;
  CSpaceThread *csThread=NULL;
  struct _u_instance clientInst;

  /* Parsing command line args. */
  while (true) {
    int opt = 0;
    int opdidx = 0;

    static struct option long_options[] = {
      { "level",     required_argument, 0, 'l'},
      { "config",    required_argument, 0, 'c'},
      { "manifest",  required_argument, 0, 'm'},
      { "help",      no_argument,       0, 'h'},
      { "version",   no_argument,       0, 'V'},
      { 0,           0,                 0,  0}
    };

    opt = getopt_long(argc, argv, "l:c:m:hV:", long_options, &opdidx);
    if (opt == -1) {
      break;
    }
    
    switch (opt) {
    case 'h':
      usage();
      exit(0);
      break;

    case 'm':
      manifestFile = optarg;
      break;
      
    case 'l':
      debug = optarg;
      set_log_level(debug);
      break;

    case 'c':
      configFile = optarg;
      break;

    case 'V':
      fprintf(stdout, "Version: %s\n", VERSION);
      exit(0);
      
    default:
      usage();
      exit(0);
    }
  }
  
  log_debug("Starting lxce.d ... \n");

  /* Initialize cspace threads list */
  init_cspace_thread_list();

  /* Before we open the socket for REST, process the config file and
   * start them containers.
   */
  config = (Config *)calloc(1, sizeof(Config));
  if (!config) {
    log_error("Memory allocation failure. Size: %d", sizeof(Config));
    exit(1);
  }

  /* Step-1: read configuration file. */
  if (process_config_file(configFile, config) != TRUE){
    log_error("Error processing the startup file");
    exit(1);
  }
  print_config(config);

  /* Step-2: read cspace config file(s) */
  if (config->cSpaceConfigs) {
    cSpaces = (CSpace *)calloc(1, sizeof(CSpace));
    if (!cSpaces) {
      log_error("Memory allocation failure. Size: %d", sizeof(CSpace));
      exit(1);
    }
    cPtr = cSpaces;
  }

  for (i=0; i<config->cSpaceCount; i++) {

    if (config->cSpaceConfigs[i]) {
      if(!process_cspace_config(config->cSpaceConfigs[i], cPtr)) {
	log_error("Error processing cSpace config file: %s",
		  config->cSpaceConfigs[i]);
	exit(1);
      }
    }

    if (i+1 != config->cSpaceCount) {
      cPtr->next =  (CSpace *)calloc(1, sizeof(CSpace));
    } else {
      cPtr->next = NULL;
    }
    cPtr = cPtr->next;
  }

  /* Step-3: setup cSpaces */
  /* For each space, we create a thread which would clone and parent
   * would wait for the space to exit. Space is currently active until the
   * device restarts.
   */
  cPtr = cSpaces;

  /* Go over the cSpaces, start thread and create actual contained spaces. */
  for (cPtr=cSpaces; cPtr; cPtr=cPtr->next) {

    csThread = init_cspace_thread(cPtr->name, cPtr);

    if (add_to_cspace_thread_list(csThread)) {
      if (pthread_create(&(csThread->tid), NULL, cspace_thread_start,
			 csThread)) {
	log_error("Error creating pthread for cSpaces. Name: %s", cPtr->name);
	exit(1);
      }
      log_debug("Thread created for cspace: %s", cPtr->name);
    } else {
      log_error("Failed to create cspace thread for: %s", cPtr->name);
    }
  }

  /* Step-4: process manifest.json file. */
  manifest = (Manifest *)calloc(1, sizeof(Manifest));
  if (!manifest) {
    log_error("Memory allocation failure. Size: %d", sizeof(Manifest));
    exit(1);
  }
  if (process_manifest(manifestFile, manifest, cSpaces) != TRUE) {
    log_error("Error process the manifest file: %s", manifestFile);
    exit(1);
  }

  /* Step-5: get manifest.json containers path from wimc */
  // get_containers_local_path(manifest, config);

  /* Step-6: open REST interface. */
  if (!start_web_services(config, &clientInst)) {
    log_error("Webservice failed to setup for clients. Exiting.");
    goto done;
  }
  
  log_debug("lxce.d running ....");
  getchar(); /* For now. XXX */

 done:
  log_debug("End World!\n");
  
  ulfius_stop_framework(&clientInst);
  ulfius_clean_instance(&clientInst);
  
  clear_config(config);
  clear_manifest(manifest);

  free(config);
  free(manifest);

  return 1;
}

