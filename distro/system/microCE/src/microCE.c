/**
 * Micro container engine.
 *
 */

#include <string.h>
#include <strings.h>
#include <ulfius.h>
#include <getopt.h>
#include <sys/types.h>
#include <unistd.h>
#include <errno.h>
#include <stdlib.h>
#include <stdio.h>

#include "log.h"
#include "toml.h"

#define TRUE 1
#define FALSE 0

#define PREFIX          "/container"
#define DEF_LOG_LEVEL   "TRACE"
#define DEF_STATUS_FILE "status"
#define DEF_CONFIG_FILE "config.toml"

#define CONFIG             "config"
#define BOOT_CONTAINER     "boot-container"
#define SERVICE_CONTAINER  "service-container"
#define SHUTDOWN_CONTAINER "shutdown-container"

#define PORT        "port"
#define STATUS_FILE "statusFile"
#define STDOUT      "stdout"
#define STDERR      "stderr"
#define LEVEL       "level"
#define ENABLE_SHUTDOWN "enableShutdown"

#define NAME         "name"
#define BUNDLE_PATH  "bundle-path"
#define VERSION      "version"
#define POLICY       "policy"
#define PID_FILE     "pid-file"
#define CONSOLE_SOCK "console-sock"
#define CONFIG_FILE  "config-file"
#define DELAY        "delay"

#define POLICY_ONE_TIME       "none"
#define POLICY_ALWAYS_RESTART "always"
#define POLICY_EXIT_0_RESTART "exit-0"

#define MAX_BUFFER 256

typedef struct {

  char *port; /* Port on which the CE is listening for REST calls. */
  char *statusFile;
  char *stdout;
  char *stderr;
  char *level;
  char *enableShutdown;
} Config;

typedef struct {

  char *name;
  char *bundlePath;
  char *version;
  char *policy;
  char *pidFile;
  char *configFile;
  char *console;
  int   delay;
} Container;

typedef struct _containers {

  Container *entry;

  struct _containers *next;
} Containers;

/*
 * callback functions declaration
 */

int callback_post_create_container (const struct _u_request *request,
				    struct _u_response *response,
				    void *user_data);
int callback_default (const struct _u_request *request,
		      struct _u_response *response, void *user_data);

/*
 * prase_containers --
 *
 *
 */

int parse_containers(Containers *containers, toml_array_t *array) {

  int ret = TRUE;
  toml_table_t* tab = NULL;
  Container *ent;
  Containers *tmp = containers;
  
  for (int i=0; (tab = toml_table_at(array, i)) != 0; i++){

    toml_datum_t name, path, version, policy, config, pidFile, delay, console;

    name = toml_string_in(tab, NAME);
    path = toml_string_in(tab, BUNDLE_PATH);
    version = toml_string_in(tab, VERSION);
    policy = toml_string_in(tab, POLICY);
    config = toml_string_in(tab, CONFIG_FILE);
    pidFile = toml_string_in(tab, PID_FILE);
    delay = toml_string_in(tab, DELAY);
    console = toml_string_in(tab, CONSOLE_SOCK);

    /* Check for mandatory options. */

    if (!name.ok || !path.ok || !version.ok || !policy.ok) {
      log_error("Missing mandatory key-value \n");
      return FALSE;
    } else {
      if (!containers->entry) {
	ent = (Container *)calloc(sizeof(Container), 1);
	if ( !ent ){
	  /* XXX */
	}
      } else {
	ent = containers->entry;
      }

      ent->name = strdup(name.u.s);
      ent->bundlePath = strdup(path.u.s);
      ent->version = strdup(version.u.s);
      ent->policy = strdup(policy.u.s);

      /* For optional, default them if not defined. */
      /* delay = 0 */

      if (!delay.ok) {
	ent->delay = 0;
      } else {
	ent->delay = atoi(delay.u.s);
      }

      if (console.ok)
	ent->console = strdup(console.u.s);

      if (pidFile.ok) {
	ent->pidFile = strdup(pidFile.u.s);
      } else {
	ent->pidFile = strdup("container.pid");
      }

      if (config.ok)
	ent->configFile = strdup(config.u.s);

      if (!containers->entry) {
	containers->entry = ent;
      }

      if (!containers->next) {
	containers->next = (Containers *)calloc(sizeof(Containers),1);
	/* XXX */
      }

      containers=containers->next;
    }
  }

  containers = tmp;
  return ret;
}


/*
 *  prase_config(configData, config);
 *
 */

int parse_config(Config *config, toml_table_t *configData) {

  int ret=FALSE;
  toml_datum_t port, statusFile, stdout, stderr, level, enableShutdown;
  
  /* Read the config data from the config.toml and load into Config. */

  port = toml_string_in(configData, PORT);
  statusFile = toml_string_in(configData, STATUS_FILE);
  stdout = toml_string_in(configData, STDOUT);
  stderr = toml_string_in(configData, STDERR);
  level = toml_string_in(configData, LEVEL);
  enableShutdown = toml_string_in(configData, ENABLE_SHUTDOWN);
  
  /* If any of the above entry is invalid, scream. */
  if (!port.ok || !statusFile.ok || !stdout.ok || !stderr.ok || !level.ok ||
      !enableShutdown.ok ){
    log_error("Cannot read valid config data \n");
  } else {

    config->port = strdup(port.u.s);
    config->statusFile = strdup(statusFile.u.s);
    config->stdout = strdup(stdout.u.s);
    config->stderr = strdup(stderr.u.s);
    config->level = strdup(level.u.s);
    config->enableShutdown = strdup(enableShutdown.u.s);

    ret = TRUE;
  }

  return ret;
}

/*
 * print_config_table -- print the config table. 
 *
 */

void print_config_table(Config *config) {


  if (config->port) {
    fprintf(stdout, "Port: %s\n", config->port);
  }

  if (config->statusFile) {
    fprintf(stdout, "StatusFile: %s\n", config->statusFile);
  }

  if (config->stdout) {
    fprintf(stdout, "stdout: %s\n", config->stdout);
  }    

  if (config->stderr) {
    fprintf(stdout, "stderr: %s\n", config->stderr);
  }    

  if (config->level) {
    fprintf(stdout, "level: %s\n", config->level);
  }    

  if (config->enableShutdown) {
    fprintf(stdout, "EnableShutdown: %s\n", config->enableShutdown);
  }

}
  
/*
 * process_config_file -- read and parse the config file. 
 *                       
 *
 */
int process_config_file(char *fileName, Config *config, Containers *boot,
			Containers *service, Containers *shutdown) {

  FILE *fp;

  toml_table_t *fileData, *configData;
  toml_array_t *containersData;
  
  char errBuf[MAX_BUFFER];
  
  if ((fp = fopen(fileName, "r")) == NULL) {
    log_error("Error opening config file: %s: %s\n", fileName,
	      strerror(errno));
    return FALSE;
  }

  /* Prase the TOML file entries. */
  fileData = toml_parse_file(fp, errBuf, sizeof(errBuf));
  
  fclose(fp);
 
  if (!fileData) {
    log_error("Error parsing the config file %s: %s\n", fileName, errBuf);
    return FALSE;
  }

  /* Parse the config. */
  configData = toml_table_in(fileData, CONFIG);

  if (configData == NULL) {
    log_error("[Config] section parsing error in file: %s\n", fileName);
    toml_free(fileData);
    return FALSE;
  }
    
  parse_config(config, configData);

  /* 1. boot containers. */
  containersData = toml_array_in(fileData, BOOT_CONTAINER);

  if (containersData == NULL ){
    log_error("[%s] section parsing error in file: %s\n", BOOT_CONTAINER,
	      fileName);
    toml_free(fileData);
    return FALSE;
  }
  parse_containers(boot, containersData);

  /* 2. service containers. */
  /* XXX code repition - fix this. */
  containersData = toml_array_in(fileData, SERVICE_CONTAINER);

  if (containersData == NULL ){
    log_error("[%s] section parsing error in file: %s\n", SERVICE_CONTAINER,
	      fileName);
    toml_free(fileData);
    return FALSE;
  }
  parse_containers(service, containersData);

  /* 3. shutdown containers. */
  containersData = toml_array_in(fileData, SHUTDOWN_CONTAINER);

  if (containersData == NULL ){
    log_error("[%s] section parsing error in file: %s\n", SHUTDOWN_CONTAINER,
	      fileName);
    toml_free(fileData);
    return FALSE;
  }
  parse_containers(shutdown, containersData);
  
  toml_free(fileData);
  return TRUE;
  
}

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
  printf("--c, --config                       Config file.\n");
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
  char *configFile = DEF_CONFIG_FILE;

  Config *config = NULL;
  Containers *boot=NULL, *service=NULL, *shutdown=NULL;
  
  /* Parsing command line args. */
  while (true) {
    int opt = 0;
    int opdidx = 0;

    static struct option long_options[] = {
      { "port",      required_argument, 0, 'p'},
      { "level",     required_argument, 0, 'l'},
      { "file",      required_argument, 0, 'f'},
      { "config",    required_argument, 0, 's'},
      { "help",      no_argument,       0, 'h'},
      { "version",   no_argument,       0, 'V'},
      { 0,           0,                 0,  0}
    };

    opt = getopt_long(argc, argv, "s:f:v:p:hV:", long_options, &opdidx);
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

    case 'c':
      configFile = optarg;
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

  /* Before we open the socket for REST, process the config file and
   * start them containers.
   */
  config = (Config *)calloc(sizeof(Config), 1);
  if (!config) {
    log_error("Memory allocation failure\n");
    exit(1);
  }

  boot     = (Containers *)calloc(sizeof(Containers), 1);
  service  = (Containers *)calloc(sizeof(Containers), 1);
  shutdown = (Containers *)calloc(sizeof(Containers), 1);
  /* XXX */
  
  if (process_config_file(configFile, config, boot, service, shutdown)
      != TRUE){
    log_error("Error processing the startup file\n");
    exit(1);
  }

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
