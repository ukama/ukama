/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "collector.h"
#include "file.h"
#include "log.h"

#include <errno.h>
#include <getopt.h>
#include <signal.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define METRIC_CONFIG "./config/metrics_config.toml"
#define DEF_LOG_LEVEL "TRACE"
#define VERSION "v0.0.1-"GITHASH

/* Terminate signal handler for Metrics collector */
void handle_sigint(int signum) {
  log_info("Metrics: Caught terminate signal.\n");

  /* Exiting Metrics */
  collector_exit(signum);
}

static struct option longOptions[] = {{"config", required_argument, 0, 'c'},
                                       {"logs", required_argument, 0, 'l'},
                                       {"help", no_argument, 0, 'h'},
                                       {"version", no_argument, 0, 'v'},
                                       {0, 0, 0, 0}};

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

/* Check if args supplied config file exist and have read permissions. */
void verify_file(char *file) {
  if (!file_exist(file)) {
    log_error("Metrics: File %s is missing.", file);
    exit(0);
  }
}

/* Usage options for the ukamaEDR */
void usage() {
  printf("Usage: metrics [options] \n");
  printf("Options:\n");
  printf("--h, --help                             Help menu.\n");
  printf("--l, --logs <TRACE> <DEBUG> <INFO>         Log level for the "
         "process.\n");
  printf("--c, --config <path>                       Config for the metrics "
         "collection.\n");
  printf("--v, --version                          Software Version.\n");
}

int main(int argc, char **argv) {
  int ret = 0;
  char *cfg = METRIC_CONFIG;
  char *debug = DEF_LOG_LEVEL;

  /* Parsing command line args. */
  while (true) {
    int opt = 0;
    int opdIdx = 0;

    opt = getopt_long(argc, argv, "c:l", longOptions, &opdIdx);
    if (opt == -1) {
      break;
    }

    switch (opt) {
    case 'h':
      usage();
      exit(0);
      break;

    case 'v':
      puts(VERSION);
      exit(0);

    case 'c':
      cfg = optarg;
      verify_file(cfg);
      break;

    case 'l':
      debug = optarg;
      set_log_level(debug);
      break;

    default:
      usage();
      exit(0);
    }
  }

  log_info("Metrics:: starting metrics collector.");

  /* Signal handler */
  signal(SIGINT, handle_sigint);

  /* Start metrics collector. */
  ret = collector(cfg);

  log_info("Metrics:: Stopping metrics collector.");
  return ret;
}
