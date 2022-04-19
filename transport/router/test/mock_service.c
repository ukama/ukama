/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* Mocking service */

#include <stdio.h>
#include <stdlib.h>
#include <ulfius.h>
#include <string.h>
#include <curl/curl.h>

#define TRUE 1
#define FALSE 0

#define MAX_LEN  1024
#define MAX_SIZE 1024

typedef struct _u_request  req_t;
typedef struct _u_response resp_t;
typedef struct _u_map      map_t;

#define REG_JSON \
  "{ \"pattern\" : %s , \"forward\": { \"ip\": \"%s\", \"port\" : \"%s\" } }"

static void print_map(map_t *map) {

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
static void print_request(const struct _u_request *request) {

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
  
  if (request->map_header) {
    fprintf(stdout, "Packet cookies: \n");
    print_map(request->map_cookie);
  }
  
  if (request->map_header) {
    fprintf(stdout, "Packet post body: \n");
    print_map(request->map_post_body);
  }
}

/* Callback function for the web application 
 *
 */
int callback_post(const req_t *request, resp_t *response, void *user_data) {

  char *str;

  print_request(request);
  str = (char *)user_data;

  ulfius_set_string_body_response(response, 200, str);

  return U_CALLBACK_CONTINUE;
}

/*
 * service_register --
 *
 */
int service_register(char *rIP, char *rPort, char *ip, char *port,
		     char *pattern) {
  int ret=FALSE;
  CURL *curl=NULL;
  char json[MAX_LEN] = {0};
  char url[MAX_LEN] = {0};
  char errBuffer[MAX_SIZE] = {0};
  CURLcode res = CURLE_FAILED_INIT;
  struct curl_slist *headers = NULL;

  sprintf(json, REG_JSON, pattern, ip, port);
  sprintf(url, "http://%s:%s/route", rIP, rPort);
  
  curl = curl_easy_init();
  if(!curl) {
    fprintf(stderr, "Error: curl_easy_init failed.\n");
    goto cleanup;
  }

  curl_easy_setopt(curl, CURLOPT_USERAGENT, "ukama/0.1");

  headers = curl_slist_append(headers, "Expect:");
  headers = curl_slist_append(headers, "Content-Type: application/json");
  curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
  curl_easy_setopt(curl, CURLOPT_POSTFIELDS, json);
  curl_easy_setopt(curl, CURLOPT_POSTFIELDSIZE, -1L);
  curl_easy_setopt(curl, CURLOPT_ERRORBUFFER, errBuffer);

  curl_easy_setopt(curl, CURLOPT_URL, url);

  fprintf(stdout, "Sending Registration JSON: %s\n", json);
  
  res = curl_easy_perform(curl);
  if (res != CURLE_OK) {
    fprintf(stderr, "error buffer: %s \n", errBuffer);
    fprintf(stderr, "curl error: %s \n", curl_easy_strerror(res));
    goto cleanup;
  } else {
    fprintf(stdout, "\nRegistration success. Status: %d \n", res);
  }

  ret = TRUE;
  
 cleanup:
  
  curl_slist_free_all(headers);
  curl_easy_cleanup(curl);
  return ret;
}

int callback_default(const req_t *request, resp_t *response, void *user_data) {

  print_request(request);

  ulfius_set_string_body_response(response, 404, "You are clearly high!\n");
  return U_CALLBACK_CONTINUE;
}

/*
 * Usage example:
 * ./mock_service 127.0.0.1 4444 4445 "hello" "{ \"key1\" : \"value1\", \"key2\" : \"value2\"}"
 *
 */
int main(int argc, char **argv) {

  char *kvPattern, *reply;
  char *rHost;
  char *port, *rPort;
  struct _u_instance inst;

  if (argc<6) {
    fprintf(stderr, "USAGE: %s router_host router_port port reply kv_pattern\n",
	    argv[0]);
    return 0;
  }

  /* Get args */
  rHost     = strdup(argv[1]);
  rPort     = strdup(argv[2]);
  port      = strdup(argv[3]);
  reply     = strdup(argv[4]);
  kvPattern = strdup(argv[5]);

  /* Initialize ulfius framework. */
  if (ulfius_init_instance(&inst, atoi(port), NULL, NULL) != U_OK) {
    fprintf(stderr, "Error ulfius_init_instance, abort\n");
    exit(1);
  }

  /* Endpoint list declaration for service. */
  ulfius_add_endpoint_by_val(&inst, "POST", "/service", NULL, 0,
                             &callback_post, (void *)reply);

  /* setup default. */
  ulfius_set_default_endpoint(&inst, &callback_default, NULL);

  /* Start the framework */
  if (ulfius_start_framework(&inst) == U_OK) {
    fprintf(stdout, "Web service started on port %d\n", inst.port);
  } else {
    fprintf(stderr, "Error starting web framework\n");
  }

  /* register the service to the router */
  service_register(rHost, rPort, "127.0.0.1", port, kvPattern);

  fprintf(stdout, "Press any key to exit ... \n");
  getchar();

  fprintf(stdout, "End service\n");

  ulfius_stop_framework(&inst);
  ulfius_clean_instance(&inst);

  return 0;
}
