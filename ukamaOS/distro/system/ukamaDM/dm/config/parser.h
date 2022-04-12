/*
 * parser.h
 *
 *  Created on: Jun 7, 2021
 *      Author: vishal
 */

#ifndef CONFIG_PARSER_H_
#define CONFIG_PARSER_H_

#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <stdlib.h>

#include "../config/toml.h"

#define SERVER_STR								"BootstrapServer"
#define SERVER_ADD_STR							"address"
#define SERVER_PORT_STR							"port"

#define FILE_STORE_STR							"FileStore"
#define CERTS_STORE_STR							"certs"
#define IP_STORE_STR							"ip"

typedef struct {
	char* addr;
	int port;
} server_cfg_t;


typedef struct {
	char* certs;
	char* addr;
} file_store_t;

typedef struct
{
	server_cfg_t* server;
	file_store_t* file_store;
} client_config_t;

int parse_config(char* cfgName);

#endif /* CONFIG_PARSER_H_ */
