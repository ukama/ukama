/**
 * 
 * Ulfius Framework example program
 * 
 * This example program describes the main features 
 * that are available in a callback function
 * 
 * Copyright 2015-2017 Nicolas Mora <mail@babelouest.org>
 * 
 * License MIT
 *
 */

#include <string.h>
#include <ulfius.h>
#include <getopt.h>

#include "log.h"

#define PREFIX "/container"
#define DEF_LOG_LEVEL "TRACE"
#define VERSION "0.0.0"

/*
 * callback functions declaration
 */

int callback_post_create_container (const struct _u_request *request,
				    struct _u_response *response,
				    void *user_data);
int callback_default (const struct _u_request *request,
		      struct _u_response *response, void *user_data);

/*
 * decode a u_map into a string
 */
char * print_map(const struct _u_map * map) {
  char * line, * to_return = NULL;
  const char **keys, * value;
  int len, i;
  if (map != NULL) {
    keys = u_map_enum_keys(map);
    for (i=0; keys[i] != NULL; i++) {
      value = u_map_get(map, keys[i]);
      len = snprintf(NULL, 0, "key is %s, value is %s", keys[i], value);
      line = o_malloc((len+1)*sizeof(char));
      snprintf(line, (len+1), "key is %s, value is %s", keys[i], value);
      if (to_return != NULL) {
        len = o_strlen(to_return) + o_strlen(line) + 1;
        to_return = o_realloc(to_return, (len+1)*sizeof(char));
        if (o_strlen(to_return) > 0) {
          strcat(to_return, "\n");
        }
      } else {
        to_return = o_malloc((o_strlen(line) + 1)*sizeof(char));
        to_return[0] = 0;
      }
      strcat(to_return, line);
      o_free(line);
    }
    return to_return;
  } else {
    return NULL;
  }
}

static struct option long_options[] = {
    { "port", required_argument, 0, 'p' },
    { "verbose", required_argument, 0, 'v' },
    { "help", no_argument, 0, 'h' },
    { "version", no_argument, 0, 'V' },
    { 0, 0, 0, 0 }
};

/* 
 * usage -- Usage options for the microCE. 
 *
 *
 */

void usage() {
  
  printf("Usage: microCE [options] \n");
  printf("Options:\n");
  printf("--h, --help                         Help menu.\n");
  printf("--v, --verbose <TRACE | DEBUG | INFO>  Log level for the process.\n");
  printf("--p, --port                         Port to listen.\n");
  printf("--V, --version                      Version.\n");
}

/*
 * set_log_level -- Set verbose level.
 *
 *
 */

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

int main(int argc, char **argv) {
  
  int ret, listen_port;
  char *debug = DEF_LOG_LEVEL;
  
  /* Parsing command line args. */
  while (true) {
    int opt = 0;
    int opdidx = 0;
    
    opt = getopt_long(argc, argv, "h:v:p:V:", long_options, &opdidx);
    if (opt == -1) {
      break;
    }
    
    switch (opt) {
    case 'h':
      usage();
      exit(0);
      break;
      
    case 'p':
      listen_port = atoi(optarg);
      break;
      
    case 'v':
      debug = optarg;
      set_log_level(debug);
      break;

    case 'V':
      fprintf(stdout, "Version: %s\n", VERSION);
      exit(0);
      
    default:
      usage();
      exit(0);
    }
  }
  
  /* Set the framework port number */
  struct _u_instance instance;
  
  log_debug("Starting micro container engine ...\n");

  if (ulfius_init_instance(&instance, listen_port, NULL, NULL) != U_OK) {
    log_error("Error initializing ulfius instance. Exit!\n");
    exit(1);
  }
  
  u_map_put(instance.default_headers, "Access-Control-Allow-Origin", "*");
  instance.max_post_body_size = 1024;
  
  /* Endpoint declaration. */
  
  ulfius_add_endpoint_by_val(&instance, "POST", PREFIX, NULL, 0,
			     &callback_post_create_container, NULL);
  ulfius_set_default_endpoint(&instance, &callback_default, NULL);
  
  /* Open an http connection. World is never going to be same!*/
  ret = ulfius_start_framework(&instance);
  
  if (ret == U_OK) {
    log_debug("Listening on port %d\n", instance.port);
    getchar(); /* For now. XXX */
  } else {
    log_error("Error starting framework\n");
  }
  
  log_debug("End World!\n");
  
  ulfius_stop_framework(&instance);
  ulfius_clean_instance(&instance);
  
  return 0;
}

/*
 * callback_post_create_container -- callback function to create container. 
 *                                   For now, response by OK!
 * 
 */ 

int callback_post_create_container(const struct _u_request *request,
				   struct _u_response *response,
				   void *user_data) {
  
  char *post_params = print_map(request->map_post_body);
  char *response_body = msprintf("OK!\n%s", post_params);
  
  ulfius_set_string_body_response(response, 200, response_body);
  o_free(response_body);
  o_free(post_params);
  
  return U_CALLBACK_CONTINUE;
}

/*
 * callback_default -- default callback for no-match
 *
 *
 */
int callback_default(const struct _u_request *request,
		     struct _u_response *response, void *user_data) {
  
  ulfius_set_string_body_response(response, 404, "You are clearly high!");
  return U_CALLBACK_CONTINUE;
}
