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
#include <sys/types.h>
#include <unistd.h>

#include "log.h"

#define PREFIX "/container"
#define DEF_LOG_LEVEL "TRACE"
#define DEF_STATUS_FILE "./status" 

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
 * update_status_file -- function to open/append/close the status
 *                       file. It is used to log PID and exit status.
 *                       If PID is 0, we write the exit status else pid.
 */

void update_status_file(char *fileName, pid_t pid, int status) {

  FILE *ptr = NULL;

  /* if file already exist and its a new start, overwrite it otherwise
   * append it.
   */

  if ( pid ){
    ptr = fopen(fileName, "w");
  } else {
    ptr = fopen(fileName, "a+");
  }

  if ( ptr ){
    if ( pid ) {
      fprintf(ptr, "%d", pid);
    } else {
      fprintf(ptr, ",%d", status);
    }
  } else {
    log_error("Unable to open file to update status.\n");
  }

  fclose(ptr);
}


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

/* 
 * usage -- Usage options for the microCE. 
 *
 *
 */

void usage() {
  
  printf("Usage: microCE [options] \n");
  printf("Options:\n");
  printf("--h, --help                         Help menu.\n");
  printf("--l, --level <TRACE | DEBUG | INFO>  Log level for the process.\n");
  printf("--p, --port                         Port to listen.\n");
  printf("--f, --file                         Status file\n");
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
  
  int ret=0, listen_port;
  char *debug = DEF_LOG_LEVEL;
  char *statusFile = DEF_STATUS_FILE;
  
  /* Parsing command line args. */
  while (true) {
    int opt = 0;
    int opdidx = 0;

    static struct option long_options[] = {
      { "port",    required_argument, 0, 'p'},
      { "level",   required_argument, 0, 'l'},
      { "file",    required_argument, 0, 'f'},
      { "help",    no_argument,       0, 'h'},
      { "version", no_argument,       0, 'V'},
      { 0,         0,                 0,  0}
    };

    
    opt = getopt_long(argc, argv, "f:v:p:hV:", long_options, &opdidx);
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
      
    case 'l':
      debug = optarg;
      set_log_level(debug);
      break;

    case 'f':
      statusFile = optarg;
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
    
    /* write PID to the status file. */
    update_status_file(statusFile, getpid(), 0);
    
    log_debug("Listening on port %d\n", instance.port);
    getchar(); /* For now. XXX */
  } else {
    log_error("Error starting framework\n");
  }
  
  log_debug("End World!\n");
  
  ulfius_stop_framework(&instance);
  ulfius_clean_instance(&instance);

  update_status_file(statusFile, 0, 0); /* 0=exit status */
  
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
