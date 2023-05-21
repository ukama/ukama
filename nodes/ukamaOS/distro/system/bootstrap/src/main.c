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
#include <errno.h>
#include <curl/curl.h>

#include "nodeInfo.h"
#include "config.h"
#include "mesh_config.h"
#include "server.h"
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

/*
 * write_to_file --
 *
 */
static int write_to_file(char *fileName, char *buffer) {

	FILE *fp=NULL;
	size_t count;

	if (fileName == NULL || buffer == NULL) return FALSE;

    fp = fopen(fileName, "w");
    if(fp == NULL) {
		log_error("Error opening file for read: %s Error: %s", fileName,
				  strerror(errno));
		return FALSE;
    }

    count = fwrite(buffer, 1, strlen(buffer), fp);
    fclose(fp);

    return count;
}

/* bootstrap */

int main (int argc, char **argv) {

	Config *config=NULL;
	MeshConfig *meshConfig=NULL;
	ServerInfo *serverInfo=NULL;
	char *configFile=NULL;
	char *debug=DEF_LOG_LEVEL;
	char *nodeID=NULL;
	int opt, opdidx, ret=TRUE;

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
	}

	config = (Config *)calloc(1, sizeof(Config));
	if (config == NULL) {
		fprintf(stderr, "Error allocating memory of: %lu", sizeof(Config));
		exit(1);
	}

	serverInfo = (ServerInfo *)calloc(1, sizeof(ServerInfo));
	if (serverInfo == NULL) {
		fprintf(stderr, "Error allocating memory of: %lu", sizeof(ServerInfo));
		exit(1);
	}

	meshConfig = (MeshConfig *)calloc(1, sizeof(MeshConfig));
	if (meshConfig == NULL) {
		fprintf(stderr, "Error allocating memory of :%lu", sizeof(MeshConfig));
		exit(1);
	}

	if (configFile == NULL) {
		configFile = DEF_CONFIG_FILE;
	}

	/* Step-1 read the configuration file */
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

	/* Step-3: connect with the ukama bootstrap server */
    send_request_to_init_with_exponential_backoff(config->bootstrapServer,
                                                  nodeID, serverInfo);
	
	/* Step-4: read mesh config file, update server IP and certs */
	if (read_mesh_config_file(config->meshConfig, meshConfig) != TRUE) {
		log_error("Error reading mesh.d config file: %s", config->meshConfig);
		goto done;
	}

	/* Step-5: update mesh.d configuration with the recevied server info. */
	ret &= write_to_file(meshConfig->remoteIPFile, serverInfo->IP);
	ret &= write_to_file(meshConfig->certFile,     serverInfo->cert);
	if (ret == FALSE) {
		log_error("Error updating mesh.d configs. File: %s",
				  config->meshConfig);
		goto done;
	}

	/* Done. */
	log_debug("Mesh.d configuration update successfully.");

 done:
	log_debug("Bye World!\n");
	clear_config(config);
	clear_mesh_config(meshConfig);
	free_server_info(serverInfo);

	free(config);
	free(meshConfig);
	free(nodeID);
	free(serverInfo);

	return 1;
}
