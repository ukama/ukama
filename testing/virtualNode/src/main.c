/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * virtualNode -- tool to create, destory, info stuff related to ukama's 
 *                virtual node
 *
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <getopt.h>

#include "node.h"
#include "config.h"
#include "log.h"
#include "supervisor.h"
#include "jserdes.h"

#define VERSION       "0.0.1"
#define DEF_LOG_LEVEL "TRACE"

#define CMD_CREATE  "create"
#define CMD_INSPECT "inspect"
#define CMD_DELETE  "delete"
#define CMD_VERIFY  "verify"

#define ENV_VNODE_METADATA "VNODE_METADATA"

enum {
	VNODE_CMD_NONE=0,
	VNODE_CMD_VERIFY,
	VNODE_CMD_CREATE,
	VNODE_CMD_INSPECT,
	VNODE_CMD_DELETE
};

extern int build_capp(Config *config);

/*
 * usage --
 *
 */
void usage() {

	printf("Usage: [options] \n");
	printf("Options:\n");
	printf("--h, --help                      help menu.\n");
	printf("--x, --exec                      command to execute\n");
	printf("                                 [create delete inspect verify]\n");
	printf("--c, --capps                     capps config folder\n");
	printf("--r, --registry                  registry URL\n");
	printf("--l, --level <ERROR | DEBUG | INFO> logging levels\n");
	printf("--V, --version                   version.\n");
}

/* 
 * set_log_level -- set the verbosity level for logs
 *
 */
static void set_log_level(char *slevel) {

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
 * get_cmd_type --
 *
 */
static int get_cmd_type(char *arg) {

	if (strcmp(arg, CMD_CREATE) == 0) {
		return VNODE_CMD_CREATE;
	} else if (strcmp(arg, CMD_DELETE) == 0) {
		return VNODE_CMD_DELETE;
	} else if (strcmp(arg, CMD_VERIFY) == 0) {
		return VNODE_CMD_VERIFY;
	} else if (strcmp(arg, CMD_INSPECT) == 0) {
		return VNODE_CMD_INSPECT;
	}

	return VNODE_CMD_NONE;
}

/*
 * Usage: virtualNode --exec [commands] [command-options] --registry [URL]
 *
 * create  --apps ./path/to/configs --registry localhost:port/name:tag
 * inspect --registry localhost:port/name:tag
 * delete  --registry localhost:port/name:tag
 *
 */
int main (int argc, char *argv[]) {

	int cmd=VNODE_CMD_NONE;
	char *configDir=NULL, *registryURL=NULL;
	char *debug=DEF_LOG_LEVEL;
	char *envVNodeMetaData=NULL;
	Configs *configs=NULL, *ptr=NULL;
	Node *node=NULL;
	json_t *jNode=NULL;

	while (TRUE) {

		int opt = 0;
		int opdidx = 0;

		static struct option long_options[] = {
			{ "exec",      required_argument, 0, 'x'},
			{ "capp",      required_argument, 0, 'c'},
			{ "registry",  required_argument, 0, 'r'},
			{ "level",     required_argument, 0, 'l'},
			{ "help",      no_argument,       0, 'h'},
			{ "version",   no_argument,       0, 'V'},
			{ 0,           0,                 0,  0}
		};

		opt = getopt_long(argc, argv, "x:c:r:l:hV:", long_options, &opdidx);
		if (opt == -1) {
			break;
		}

		switch (opt) {
		case 'h':
			usage();
			exit(0);
			break;

		case 'x':
			cmd = get_cmd_type(optarg);
			break;

		case 'c':
			configDir = optarg;
			break;

		case 'r':
			registryURL = optarg;
			break;

		case 'l':
			debug = optarg;
			set_log_level(debug);
			break;

		case 'V':
			fprintf(stdout, "%s - version: %s\n", argv[0], VERSION);
			exit(0);

		default:
			usage();
			exit(0);
		}
	} /* while */

	if (argc == 1 || cmd == VNODE_CMD_NONE) {
		fprintf(stderr, "Must specify command\n");
		usage();
		exit(1);
	}

	envVNodeMetaData = getenv(ENV_VNODE_METADATA);
	if (envVNodeMetaData == NULL) {
	  log_error("Env variable: %s not set \n Exiting.", ENV_VNODE_METADATA);
	  exit(1);
	}

	jNode = json_loads(envVNodeMetaData, JSON_DECODE_ANY, NULL);
	if (!jNode) {
	  log_error("Invalid JSON for in env variable: %s\n Exiting",
				ENV_VNODE_METADATA);
	  exit(1);
	}

	if (!deserialize_node(&node, jNode)) {
	  log_error("Unable to deserialize env variable: %s\n Exiting.",
				ENV_VNODE_METADATA);
	  exit(1);
	}

	if (!read_config_files(&configs, configDir)) {
	    log_error("Parsing error reading configs from %s \n. Exiting.",
				  configDir);
		free_configs(configs);
		exit(1);
	}

	/* Build all them capps */
	ptr = configs;
	while (ptr) {
	  if (ptr->valid && ptr->config) {
		if (!build_capp(ptr->config)) {
		  log_error("Error building capp %s:%s using config file: %s",
					ptr->config->capp->name, ptr->config->capp->version,
					ptr->fileName);
		  free_configs(configs);
		  /* XXX clean up build dir */
		  exit(1);
		}
	  }
	  ptr = ptr->next;
	}

	/* Create config file supervisor.d */
	if (!create_supervisor_config(configs)) {
	  log_error("Unable to create configuration file for supervisor.d");
	  purge_supervisor_config(SVISOR_FILENAME);
	  free_configs(configs);
	  exit(1);
	}

 done:
	free_configs(configs);
	json_decref(jNode);
	free_node(node);
	exit(0);
}
