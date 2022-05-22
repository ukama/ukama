/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Ukama Node bootstrap --
 *
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <getopt.h>
#include <ulfius.h>
#include <curl/curl.h>

#include "nodeInfo.h"
#include "config.h"
#include "log.h"

#define DEF_LOG_LEVEL "TRACE"

#define VERSION "0.0.1"

/*
 * usage -- Usage options
 *
 *
 */
static void usage() {

	printf("bootstrap: ukama's node bootstrap client \n");
	printf("Usage: bootstrap [options] \n");
	printf("Options:\n");
	printf("--h, --help                         This help menu. \n");
	printf("--c, --config                       Configuration file \n");
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

/* bootstrap */

int main (int argc, char **argv) {

	Config *config=NULL;
	char *configFile=NULL;
	char *debug=DEF_LOG_LEVEL;
	char *nodeID=NULL;
	int opt, opdidx;

	/* Prase command line args. */
	while (TRUE) {

		opt    = 0;
		opdidx = 0;

		static struct option long_options[] = {
			{ "config",  required_argument, 0, 'c'},
			{ "level",   required_argument, 0, 'l'},
			{ "help",    no_argument,       0, 'h'},
			{ "version", no_argument,       0, 'v'},
			{ 0,         0,                 0,  0}
		};

		opt = getopt_long(argc, argv, "c:l:hv:", long_options, &opdidx);
		if (opt == -1) {
			break;
		}

		switch (opt) {
		case 'c':
			configFile = optarg;
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
			fprintf(stdout, "bootstrap - Version: %s\n", VERSION);
			exit(0);

		default:
			usage();
			exit(0);
		}
	} /* while */

	config = (Config *)calloc(1, sizeof(Config));
	if (config == NULL) {
		fprintf(stderr, "Error allocating memory of: %lu", sizeof(Config));
		exit(1);
	}
	
	if (configFile == NULL) {
		configFile = DEF_CONFIG_FILE;
	}

	/* Step-1 read the configuration file. */
	if (process_config_file(configFile, config) != TRUE) {
		log_error("Error processing the config file: %s", configFile);
		exit(1);
	}
	print_config(config);

	/* Step-2: request node.d for NodeID */
	if (get_nodeID_from_noded(&nodeID, config->nodedHost, config->nodedPort)
		!= TRUE) {
	    log_error("Error retreiving NodeID from noded.d at %s:%s",
				  config->nodedHost, config->nodedPort);
		goto done;
	}

	/* Step-3: connect with ukama bootstrap server */
	
	
	/* Step-4: update config file for mesh.d */
	
	getchar(); /* For now. XXX */

 done:
	log_debug("Bye World!\n");
	clear_config(config);
	free(config);
	free(nodeID);
  
	return 1;
}
