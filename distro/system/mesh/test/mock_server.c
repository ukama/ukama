/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* Mocking services end-points for testing. */

#include <stdio.h>
#include <stdlib.h>
#include <ulfius.h>

#define EP_REGISTRY "/registry/"

#define TRUE 1
#define FALSE 0

typedef struct _u_request  req_t;
typedef struct _u_response resp_t;
typedef struct _u_map      map_t;


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
int callback_registry(req_t *request, resp_t *response, void * user_data) {

  print_request(request);
  return U_CALLBACK_CONTINUE;
}

int callback_default(req_t *req, resp_t * resp, void * user_data) {

  return U_CALLBACK_CONTINUE;
}

int main(int argc, char **argv) {

  int port;
  struct _u_instance inst;

  if (argc<2) {
    fprintf(stderr, "USAGE: %s port\n", argv[0]);
    return 0;
  }

  port = atoi(argv[1]);

  /* Initialize ulfius framework. */
  if (ulfius_init_instance(&inst, port, NULL, NULL) != U_OK) {
    fprintf(stderr, "Error ulfius_init_instance, abort\n");
    exit(1);
  }

  /* Endpoint list declaration for registry service. */
  ulfius_add_endpoint_by_val(&inst, "GET", EP_REGISTRY, NULL, 0,
			     &callback_registry, NULL);
  ulfius_add_endpoint_by_val(&inst, "POST", EP_REGISTRY, NULL, 0,
			     &callback_registry, NULL);
  ulfius_add_endpoint_by_val(&inst, "PUT", EP_REGISTRY, NULL, 0,
			     &callback_registry, NULL);
   ulfius_add_endpoint_by_val(&inst, "DELETE", EP_REGISTRY, NULL, 0,
			     &callback_registry, NULL);

   /* setup default. */
   ulfius_set_default_endpoint(&inst, &callback_default, NULL);

   
  /* Start the framework */
  if (ulfius_start_framework(&inst) == U_OK) {
    fprintf(stdout, "Web service started on port %d\n", inst.port);
    fprintf(stdout, "Press any key to exit ... \n");
    getchar();
  }
  else {
    fprintf(stderr, "Error starting framework\n");
  }
  
  fprintf(stdout, "End framework\n");

  ulfius_stop_framework(&inst);
  ulfius_clean_instance(&inst);

  return 0;
}
