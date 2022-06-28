/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* Mocking bootstrap server -- for testing only */

#include <stdio.h>
#include <stdlib.h>
#include <sqlite3.h>
#include <string.h>
#include <jansson.h>
#include <ulfius.h>

#define EP_SERVICE       "/service"
#define JSON_NODE        "node"
#define JSON_ORG         "org"
#define JSON_IP          "ip"
#define JSON_CERTIFICATE "certificate"

#define TRUE 1
#define FALSE 0

typedef struct _u_request  req_t;
typedef struct _u_response resp_t;
typedef struct _u_map      map_t;

/* Struct to define the server */
typedef struct {

	char *serverIP; /* Server's IPv4 for Mesh.d */
	char *cert;     /* Cert for connection with Server */
	char *orgName;  /* Organization */
} ServerInfo;

void print_map(map_t *map) {

	int i;

	if (map==NULL) return;

	for (i=0; i<map->nb_values; i++) {
		fprintf(stdout, "\t %s:%s \n", map->keys[i], map->values[i]);
	}
	
	fprintf(stdout, "----------------------\n");
}

/*
 * print_request -- print various parameters for the incoming request.
 *
 */
void print_request(const struct _u_request *request) {

	fprintf(stdout, "Recevied Packet: \n");
	fprintf(stdout, " Protocol: %s \n Method: %s \n Path: %s\n",
			request->http_protocol, request->http_verb, request->http_url);
	
	if (request->map_header) {
		fprintf(stdout, "Packet header: \n");
		print_map(request->map_header);
	}
	
	if (request->map_url) {
		fprintf(stdout, "Packet URL variables: \n");
		print_map(request->map_url);
	}
}

int is_valid_request(req_t *request, char **nodeID) {

	map_t *map;

	if (request == NULL) return FALSE;

	map = request->map_url;
	if (!map)                return FALSE;
	if (map->nb_values != 2) return FALSE;

	if (strcmp(map->keys[0], "node") == 0 &&
		strcmp(map->keys[1], "looking_for") == 0) {
		if (strcmp(map->values[1], "validation") == 0) {
			*nodeID = strdup(map->values[0]);
			return TRUE;
		}
	}

	return FALSE;
}

/* Callback function for the web application on /service
 * GET 'http://localhost:8091/service/?node=UK-SA2154-HNODE-A1-0001&looking_for=validation'
 *
 * check if the request has 1. nodeID and 2. looking_for=validation. Ignore 
 * anything else.
 */
int callback_get_service(req_t *request, resp_t *response, void * user_data) {

	char *str=NULL;
	char *nodeID=NULL;
	json_t *json=NULL;
	ServerInfo *serverInfo;

	print_request(request);

	if (!is_valid_request(request, &nodeID)) {
		fprintf(stderr, "Invalid request recevied. 400. Ignoring\n");
		ulfius_set_string_body_response(response, 400, "");
		goto done;
	}

	/* response with serverInfo in json format */
	json       = json_object();
	serverInfo = (ServerInfo *)user_data;

	json_object_set_new(json, JSON_NODE, json_string(nodeID));
	json_object_set_new(json, JSON_ORG,  json_string(serverInfo->orgName));
	json_object_set_new(json, JSON_IP,   json_string(serverInfo->serverIP));
	json_object_set_new(json, JSON_CERTIFICATE,
						json_string(serverInfo->cert));

	str = json_dumps(json, 0);
	fprintf(stdout, "ServerInfo response: %s\n", str);
	free(str);

	ulfius_set_json_body_response(response, 200, json);
	json_decref(json);

 done:
	if (nodeID) free(nodeID);
	return U_CALLBACK_CONTINUE;
}

int callback_default_ok(req_t *req, resp_t * resp, void * user_data) {

	ulfius_set_string_body_response(resp, 200, "OK\n");
	return U_CALLBACK_CONTINUE;
}

int read_cert_file(char *fileName, ServerInfo *serverInfo) {

	FILE *fp=NULL;
	long fSize=0;

	fp = fopen(fileName, "r");
	if (fp == NULL) {
		fprintf(stderr, "Unable to open cert file: %s", fileName);
		return FALSE;
	}

	fseek(fp, 0, SEEK_END);
	fSize = ftell(fp);
	fseek(fp, 0, SEEK_SET);

	serverInfo->cert = (char *)calloc(1, fSize+1);
	if (serverInfo->cert == NULL) {
		fprintf(stderr, "Unable to allocate memory of size: %lu", fSize+1);
		return FALSE;
	}

	fread(serverInfo->cert, fSize, 1, fp);
	fclose(fp);

	serverInfo->cert[fSize] = 0;

	return TRUE;
}

void free_server_info(ServerInfo *ptr) {

	if (ptr == NULL) return;

	if (ptr->serverIP) free(ptr->serverIP);
	if (ptr->cert)     free(ptr->cert);
	if (ptr->orgName)  free(ptr->orgName);

	free(ptr);
}

/*
 * {
 *   "node": "uk-sa2220-hnode-v0-dcf4",
 *   "org": "test",
 *   "ip": "192.168.0.1",
 *   "certificate": "aGVscG1lCg=="
 * }
 *
 */

/* ./mock_bootstrap --ip 192.168.0.1 --cert ./file.cert --org test */

int main(int argc, char **argv) {

	struct _u_instance inst;
	ServerInfo *serverInfo=NULL;
	int port;

	if (argc<2) {
		fprintf(stderr, "USAGE: %s port serverIP certFile orgName\n", argv[0]);
		return 0;
	}

	serverInfo = (ServerInfo *)calloc(1, sizeof(ServerInfo));
	if (serverInfo == NULL) {
		fprintf(stderr, "Error allocating memory of size: %lu\n",
				sizeof(ServerInfo));
		return 1;
	}

	port                 = atoi(argv[1]);
	serverInfo->serverIP = strdup(argv[2]);
	serverInfo->orgName  = strdup(argv[4]);

	if (!read_cert_file(argv[3], serverInfo)) {
		fprintf(stderr, "Unable to read cert file: %s", argv[3]);
		return 1;
	}
	
	if (ulfius_init_instance(&inst, port, NULL, NULL) != U_OK) {
		fprintf(stderr, "Error ulfius_init_instance, abort\n");
		return 1;
	}

	/* Endpoint list declaration */
	ulfius_add_endpoint_by_val(&inst, "GET", EP_SERVICE, NULL, 0,
							   &callback_get_service, serverInfo);

	/* Start the framework */
	if (ulfius_start_framework(&inst) == U_OK) {
		fprintf(stdout, "Famework start on port %d\n", inst.port);
		getchar();
	}
	else {
		fprintf(stderr, "Error starting framework\n");
	}

	fprintf(stdout, "End framework\n");

	ulfius_stop_framework(&inst);
	ulfius_clean_instance(&inst);

	free_server_info(serverInfo);
	
	return 0;
}
