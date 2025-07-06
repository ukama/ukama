/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <getopt.h>
#include <sys/stat.h>

#include "config_app.h"
#include "log_app.h"

#define VERSION       "0.0.1"
#define DEF_LOG_LEVEL "TRACE"
#define MAX_BUFFER    256

enum {
  CAPP_CMD_NONE=0,
  CAPP_CMD_VERIFY,
  CAPP_CMD_CREATE,
};

extern int build_app(Config *config);
extern int create_app(Config *config);

void usage() {

  printf("Usage: [options] \n");
  printf("UKAMA_ROOT environment variable is required\n");
  printf("Options:\n");
  printf("--h, --help                         Help menu.\n");
  printf("--c, --create                       Create capp\n");
  printf("--v, --verify                       Verify config file\n");
  printf("--C, --config                       Config file name\n");
  printf("--l, --level <ERROR | DEBUG | INFO> Logging levels\n");
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

  int cmd=CAPP_CMD_NONE, ret=FALSE;
  int exitStatus=1;
  char *configFile=NULL;
  char cappFile[MAX_BUFFER] = {0};
  char *debug=DEF_LOG_LEVEL;
  struct stat st;
  Config *config=NULL;
  
  /* Prase command line args. */
  while (TRUE) {

    int opt = 0;
    int opdidx = 0;

    static struct option long_options[] = {
      { "create",    no_argument,       0, 'c'},
      { "verify",    no_argument,       0, 'v'},
      { "config",    required_argument, 0, 'C'},
      { "level",     required_argument, 0, 'l'},
      { "help",      no_argument,       0, 'h'},
      { "version",   no_argument,       0, 'V'},
      { 0,           0,                 0,  0}
    };

    opt = getopt_long(argc, argv, "C:cvhV:", long_options, &opdidx);
    if (opt == -1) {
      break;
    }

    switch (opt) {
    case 'h':
      usage();
      exit(0);
      break;

    case 'c':
      cmd = CAPP_CMD_CREATE;
      break;

    case 'v':
      cmd = CAPP_CMD_VERIFY;
      break;

    case 'C':
      configFile = optarg;
      break;
      
    case 'l':
      debug = optarg;
      set_log_level(debug);
      break;

    case 'V':
      fprintf(stdout, "capp - Version: %s\n", VERSION);
      exit(0);

    default:
      usage();
      exit(0);
    }
  } /* while */

  if (argc == 1 || cmd == CAPP_CMD_NONE) {
    fprintf(stderr, "Missing required parameters.\n");
    usage();
    exit(1);
  }

  if (getenv("UKAMA_ROOT") == NULL) {
      fprintf(stderr, "UKAMA_ROOT is not set.\n");
      usage();
      exit(1);
  }

  config = (Config *)calloc(1, sizeof(Config));
  if (!config) {
    log_error("Memory allocation error of size: %d Exiting", sizeof(Config));
    exit(1);
  }
  
  ret = read_config_file(config, configFile);
  if (!ret) {
    log_error("%s parsing error. Exiting.", configFile);
    clear_config(config, BUILD_ONLY & CAPP_ONLY);
    free(config);
    exit(1);
  } else {
    log_config(config);
  }

  if (!build_app(config)) {
    log_error("Error building the capp using: %s", configFile);
    goto done;
  }

  if (!create_app(config)) {
    log_error("Error creating the capp using: %s", configFile);
    goto done;
  }

  sprintf(cappFile, "%s/build/pkgs/%s_%s.tar.gz",
          getenv("UKAMA_ROOT"),
          config->capp->name,
          config->capp->version);
  stat(cappFile, &st);
  log_debug("All done. \n cApp: %s \n Size: %dK", cappFile, (int)st.st_size/1000);

  exitStatus = 0;

 done:
  clear_config(config, BUILD_ONLY | CAPP_ONLY);
  free(config);
  return exitStatus;
}
