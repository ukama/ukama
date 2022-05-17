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
#include <jansson.h>

#define TRUE 1
#define FALSE 0

#define MAX_LEN  1024
#define MAX_SIZE 1024

#define JSON_UUID "uuid"

typedef struct _u_request  req_t;
typedef struct _u_response resp_t;
typedef struct _u_map      map_t;

struct Response {
  char *buffer;
  size_t size;
};

/*
 * {
 *      "name": "service_name",
 *	"patterns": [{
 *			"key1": "value1",
 *			"key1": "value2",
 *			"path": "/abc"
 *		},
 *		{
 *			"key1": "value1",
 *			"path": "/abv/xcv"
 *		}
 *	],
 *	"forward": {
 *		"ip": "10.0.0.1",
 *		"port": 8080,
 *		"default_path": "/abc"
 *	}
 * }
 *
 */

#define DEL_JSON "{ \"uuid\" : \"%s\" }"

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
int callback_service(const req_t *request, resp_t *response, void *userData) {

  char *str;
  char buffer[MAX_LEN] = {0};

  print_request(request);
  str = (char *)userData;

  sprintf(buffer, "%s: %s\n", request->http_verb, str);

  ulfius_set_string_body_response(response, 200, buffer);

  return U_CALLBACK_CONTINUE;
}

/* Callback function for /ping
 *
 */
int callback_ping(const req_t *request, resp_t *response, void *userData) {

  char *str;
  char buffer[MAX_LEN] = {0};

  print_request(request);
  str = (char *)userData;

  sprintf(buffer, "%s: %s\n", request->http_verb, str);

  ulfius_set_string_body_response(response, 200, buffer);

  return U_CALLBACK_CONTINUE;
}

/*
 * response_callback --
 */
static size_t response_callback(void *contents, size_t size, size_t nmemb,
                                void *userp) {

  size_t realsize = size * nmemb;
  struct Response *response = (struct Response *)userp;

  response->buffer = realloc(response->buffer, response->size + realsize + 1);

  if(response->buffer == NULL) {
    fprintf(stderr, "Not enough memory to realloc of size: %ld",
              response->size + realsize + 1);
    return 0;
  }

  memcpy(&(response->buffer[response->size]), contents, realsize);
  response->size += realsize;
  response->buffer[response->size] = 0; /* Null terminate. */

  return realsize;
}

/*
 * service_unregister --
 *
 */
static int service_unregister(char *rIP, char *rPort, char *uuidStr) {

  int ret=FALSE;
  CURL *curl=NULL;
  char json[MAX_LEN] = {0};
  char url[MAX_LEN] = {0};
  char errBuffer[MAX_SIZE] = {0};
  CURLcode res = CURLE_FAILED_INIT;
  struct curl_slist *headers = NULL;

  if (uuidStr == NULL) return FALSE;

  sprintf(json, DEL_JSON, uuidStr);
  sprintf(url, "http://%s:%s/routes", rIP, rPort);

  curl = curl_easy_init();
  if(!curl) {
    fprintf(stderr, "Error: curl_easy_init failed.\n");
    goto cleanup;
  }

  curl_easy_setopt(curl, CURLOPT_USERAGENT, "ukama/0.1");

  headers = curl_slist_append(headers, "Expect:");
  headers = curl_slist_append(headers, "Content-Type: application/json");

  curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "DELETE");
  curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
  curl_easy_setopt(curl, CURLOPT_POSTFIELDS, json);
  curl_easy_setopt(curl, CURLOPT_POSTFIELDSIZE, -1L);
  curl_easy_setopt(curl, CURLOPT_ERRORBUFFER, errBuffer);

  curl_easy_setopt(curl, CURLOPT_URL, url);

  fprintf(stdout, "Sending un-register JSON: %s\n", json);

  res = curl_easy_perform(curl);
  if (res != CURLE_OK) {
    fprintf(stderr, "error buffer: %s \n", errBuffer);
    fprintf(stderr, "curl error: %s \n", curl_easy_strerror(res));
    goto cleanup;
  } else {
    fprintf(stdout, "\n un-register success. Status: %d \n", res);
  }

  ret = TRUE;

 cleanup:

  curl_slist_free_all(headers);
  curl_easy_cleanup(curl);

  return ret;
}

/*
 * read_pattern_from_file --
 *
 */
static int read_pattern_from_file(char *fileName, char **buffer) {

  FILE *fp;
  long fSize=0;

  fp = fopen(fileName, "r");
  if (fp==NULL) {
    fprintf(stderr, "Error opening file: %s", fileName);
    return FALSE;
  }

  fseek(fp, 0, SEEK_END);
  fSize = ftell(fp);
  fseek(fp, 0, SEEK_SET);

  *buffer = (char *)malloc(fSize + 1);
  fread(*buffer, fSize, 1, fp);
  fclose(fp);

  return TRUE;
}

/*
 * service_register --
 *
 */
static int service_register(char *name, char *rIP, char *rPort, char *ip,
          int port, char *path, char *fileName,
          char **uuidStr) {

  int ret=FALSE;
  CURL *curl=NULL;
  char *json=NULL;
  char url[MAX_LEN] = {0};
  char errBuffer[MAX_SIZE] = {0};
  CURLcode res = CURLE_FAILED_INIT;
  struct curl_slist *headers = NULL;
  struct Response response;
  json_t *jRoot=NULL, *jID=NULL;

  if (!read_pattern_from_file(fileName, &json)) {
    exit(1);
  }

  sprintf(url, "http://%s:%s/routes", rIP, rPort);

  curl = curl_easy_init();
  if(!curl) {
    fprintf(stderr, "Error: curl_easy_init failed.\n");
    goto cleanup;
  }

  response.buffer = (char *)malloc(1);
  response.size   = 0;

  curl_easy_setopt(curl, CURLOPT_USERAGENT, "ukama/0.1");

  headers = curl_slist_append(headers, "Expect:");
  headers = curl_slist_append(headers, "Content-Type: application/json");
  curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
  curl_easy_setopt(curl, CURLOPT_POSTFIELDS, json);
  curl_easy_setopt(curl, CURLOPT_POSTFIELDSIZE, -1L);
  curl_easy_setopt(curl, CURLOPT_ERRORBUFFER, errBuffer);
  curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
  curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);

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

  /* get UUID */
  jRoot = json_loads(response.buffer, JSON_DECODE_ANY, NULL);
  if (!jRoot) {
    fprintf(stderr, "Can not load str into JSON object. Str: %s",
      response.buffer);
    goto cleanup;
  }

  jID = json_object_get(jRoot, JSON_UUID);
  if (jID == NULL) {
    fprintf(stderr, "Unable to find %s in response", JSON_UUID);
    goto cleanup;
  }

  *uuidStr = strdup(json_string_value(jID));

  ret = TRUE;

 cleanup:

  curl_slist_free_all(headers);
  curl_easy_cleanup(curl);
  json_decref(jRoot);
  if (json) free(json);

  return ret;
}

int callback_default(const req_t *request, resp_t *response, void *userData) {

  char *name;
  char buffer[MAX_LEN] = {0};

  name = (char *)userData;

  print_request(request);

  sprintf(buffer, "%s: not implemented. End-point: %s \n", name,
    request->url_path);

  ulfius_set_string_body_response(response, 404, buffer);
  return U_CALLBACK_CONTINUE;
}

/*
 * Usage example:
 * ./mock_service 127.0.0.1 4444 4445 "/service" "hello" ./pattern.sample
 *
 */
int main(int argc, char **argv) {

  char *name;
  char *reply;
  char *rHost;
  int  port;
  char *rPort;
  char *uuidStr=NULL;
  char *path;
  char *file;
  struct _u_instance inst;

  if (argc<6) {
    fprintf(stderr, "USAGE: %s router_host router_port port reply kv_pattern\n",
      argv[0]);
    return 0;
  }

  /* Get args */
  name      = strdup(argv[0]);
  rHost     = strdup(argv[1]);
  rPort     = strdup(argv[2]);
  port      = atoi(argv[3]);
  path      = strdup(argv[4]);
  reply     = strdup(argv[5]);
  file      = strdup(argv[6]);

  /* Initialize ulfius framework. */
  if (ulfius_init_instance(&inst, port, NULL, NULL) != U_OK) {
    fprintf(stderr, "Error ulfius_init_instance, abort\n");
    exit(1);
  }

  /* Endpoint list declaration for service. */
  ulfius_add_endpoint_by_val(&inst, "GET", path, NULL, 0,
                             &callback_service, (void *)reply);
  ulfius_add_endpoint_by_val(&inst, "POST", path, NULL, 0,
                             &callback_service, (void *)reply);
  ulfius_add_endpoint_by_val(&inst, "PUT", path, NULL, 0,
                             &callback_service, (void *)reply);
  ulfius_add_endpoint_by_val(&inst, "DELETE", path, NULL, 0,
                             &callback_service, (void *)reply);

  /* /ping */
  ulfius_add_endpoint_by_val(&inst, "GET", "/ping", NULL, 0,
                             &callback_ping, "pong");

  /* setup default. */
  ulfius_set_default_endpoint(&inst, &callback_default, name);

  /* Start the framework */
  if (ulfius_start_framework(&inst) == U_OK) {
    fprintf(stdout, "Web service started on port %d\n", inst.port);
  } else {
    fprintf(stderr, "Error starting web framework\n");
  }

  /* register the service to the router */
  service_register(name, rHost, rPort, "127.0.0.1", port, path, file,
       &uuidStr);
  fprintf(stdout, "UUID: %s\n", uuidStr);

  fprintf(stdout, "Press any key to exit ... \n");
  getchar();

  /*unregister the service from the router */
  service_unregister(rHost, rPort, uuidStr);
  fprintf(stdout, "Service un-registered\n");

  fprintf(stdout, "End service\n");

  ulfius_stop_framework(&inst);
  ulfius_clean_instance(&inst);

  return 0;
}
