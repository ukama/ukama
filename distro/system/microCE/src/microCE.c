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

#include "log.h"

#define TRUE 1
#define FALSE 0

#define PREFIX "/container"
#define DEF_LOG_LEVEL "TRACE"
#define DEF_STATUS_FILE "status"
#define DEF_STARTUP_FILE "startup.cfg"

#define VERSION "0.0.0"

#define MAXL 1024
#define MAXC 128

#define TYPE_CENT_INIT    1
#define TYPE_CENT_ONBOOT  2
#define TYPE_CENT_SERVICE 3

#define HEADER_TOKEN "#startup"
#define INIT_TOKEN "[init]"
#define ONBOOT_TOKEN "[onboot]"
#define SERVICE_TOKEN "[service]"

#define STARTUP_VERSION 1

#define DELIM ";"

typedef struct {

  char *name;
  char *image;
  char *version;
  char *path;
  char *json;
  char *policy;

  int type; /* init, onboot, service. */
} sFileEntries;

typedef struct _sType {

  sFileEntries *entry;
  struct _sType *next;

} sType;


/*
 * callback functions declaration
 */

int callback_post_create_container (const struct _u_request *request,
				    struct _u_response *response,
				    void *user_data);
int callback_default (const struct _u_request *request,
		      struct _u_response *response, void *user_data);

/*
 * type_to_str -- Convert the type into string for printing.
 *
 *
 */

char* type_to_str(int type) {

  switch(type) {
  case TYPE_CENT_INIT:
    return INIT_TOKEN;
    break;

  case TYPE_CENT_ONBOOT:
    return ONBOOT_TOKEN;
    break;

  case TYPE_CENT_SERVICE:
    return SERVICE_TOKEN;
    break;

  default:
     return "";
  }

  return "";
}

/*
 * valid_config_file -- check if the config file had valid header, if so
 *                      return the version number else -1
 *
 *                      #startup:<version>
 */

int valid_config_file(char *line) {

  int ret = -1;

  /* We expect [0]=#, [1-6]="config", [7]=":" and [8-]=digits */

  char *token = strtok(line, ":");

  if (strcmp(token, HEADER_TOKEN)==0){
    char *ver = strtok(NULL, ":");
    ret = atoi(ver); /*XXX - no error handling, better to use strtol() */
  } else {
    log_error("Invalid header for the startup file\n");
  }

  return ret;
}

/*
 * is_comment -- check to see if the line is a comment.
 *
 *
 */

int is_comment(char *line) {

  const char whitespace[] = " \f\n\r\t\v";
  int index;

  index = strspn(line, whitespace);

  if (line[index] == '#') {
    return TRUE;
  }

  return FALSE;
}

/*
 * is_token -- check to see if the line is a token.
 *
 *
 */

int is_token(char *line) {

  const char whitespace[] = " \f\n\r\t\v";

  char *token = line + strspn(line, whitespace);

  if (strncmp(token, INIT_TOKEN, strlen(INIT_TOKEN))==0) {
    return TYPE_CENT_INIT;
  } else if (strncmp(token, ONBOOT_TOKEN, strlen(ONBOOT_TOKEN))==0) {
    return TYPE_CENT_ONBOOT;
  } else if (strncmp(token, SERVICE_TOKEN, strlen(SERVICE_TOKEN))==0) {
    return TYPE_CENT_SERVICE;
  }

  return FALSE;
}

/*
 * allocate_entry -- allocate memory to the entry structure.
 *
 *
 */

int allocate_entry(sFileEntries *ptr) {

  int ret=FALSE;

  if (ptr) {

    ptr->name = (char *)malloc(MAXC);
    ptr->image = (char *)malloc(MAXC);
    ptr->version = (char *)malloc(MAXC);
    ptr->path = (char *)malloc(MAXC);
    ptr->policy = (char *)malloc(MAXC);
    ptr->json = (char *)malloc(MAXC);
    /* XXX - Fix me. */

    ptr->type = -1;

    ret = TRUE;
  }

  return ret;
}

/*
 * process_line -- process the configuration line and update the respective
 *                 list.
 *
 *
 */

int process_line(char *line, int tokenType,  sType *initCon,  sType *onbootCon,
			  sType *serviceCon) {

  sType *ptr=NULL;
  sFileEntries *ePtr = NULL;
  char *token;

  if (tokenType == TYPE_CENT_INIT) {
    ptr = initCon;
  } else if (tokenType == TYPE_CENT_ONBOOT) {
    ptr = onbootCon;
  } else if (tokenType == TYPE_CENT_SERVICE) {
    ptr = serviceCon;
  }

  /* Forward to the last entry in the list */
  while (ptr->next) {
    ptr = ptr->next;
  }

  ptr->entry = (sFileEntries *)malloc(sizeof(sFileEntries));
  if (ptr->entry == NULL){
    /* XXX Fix me. */
    return 0;
  }

  ePtr = ptr->entry;

  if (allocate_entry(ptr->entry) == FALSE) {
    /* Fix me. */
    free(ptr->entry);
    return FALSE;
  }

  /* Process each token in given order. */

  /* 1. Name */
  token = strtok(line, DELIM);
  strcpy(ePtr->name, token);

  /* 2. Image */
  token = strtok(NULL, DELIM);
  strcpy(ePtr->image, token);

  /* 3. version */
  token = strtok(NULL, DELIM);
  strcpy(ePtr->version, token);

  /* 4. path */
  token = strtok(NULL, DELIM);
  strcpy(ePtr->path, token);

  /* 5. policy */
  token = strtok(NULL, DELIM);
  strcpy(ePtr->policy, token);

  /* 6. JSON file */
  token = strtok(NULL, DELIM);
  strcpy(ePtr->json, token);

  ePtr->type = tokenType;

  return TRUE;
}

/*
 * print_startup_file -- Print the info on the stdout.
 *
 */
void print_startupFile_info (sType *ptr) {

  int count=1;

  while (ptr) {

    sFileEntries *sPtr = ptr->entry;

    if (sPtr) {
      fprintf(stdout, "Entry:%d Type: %s\n", count, type_to_str(sPtr->type));
      fprintf(stdout, "Name: %s\n", sPtr->name);
      fprintf(stdout, "Image: %s\n", sPtr->image);
      fprintf(stdout, "Version: %s\n", sPtr->version);
      fprintf(stdout, "Path: %s\n", sPtr->path);
      fprintf(stdout, "Policy: %s\n", sPtr->policy);
      fprintf(stdout, "Json: %s\n", sPtr->json);
    } else {
      return;
    }

    count++;
    ptr = ptr->next;
  }

}

/*
 * process_startup_file -- read and parse the startup file.
 *
 *
 */
int process_startup_file(char *fileName, sType *initCon,  sType *onbootCon,
			 sType *serviceCon) {

  int ret=FALSE, ver=-1, lineCount=0;
  int curr_token, nread;
  int processedToken[4] = {0, 0, 0, 0};

  FILE *fp;
  size_t len=0;
  char *line=NULL;


  if ((fp = fopen(fileName, "r")) == NULL) {
    log_error("Error opening startup file: %s: %s\n", fileName,
	      strerror(errno));
    return FALSE;
  }

  /* Frist line needs to be "#startup:<version> */
  if ((nread=getline(&line, &len, fp)) != -1){

    ver = valid_config_file(line);
    lineCount++;

    if (ver==-1) {
      log_error("%s:%d Invalid header\n");
      goto done;
    }

    if (ver > STARTUP_VERSION){
      log_error("%s:%d Startup file version mismatch. Expected:%d, Got:%d\n",
		fileName, lineCount, STARTUP_VERSION, ver);
      goto done;
    }
   } else {
    log_error("%s:%d Error reading startup file\n", fileName, lineCount);
    goto done;
  }

  while ((nread=getline(&line, &len, fp)) != -1){ /* read each line in file  */

    int token=0;

    lineCount++;

    /* Skip empty lines. */
    if (line[0] == '\n' || line[0] == 0)
      continue;

    /* Skip comment lines. */
    if (is_comment(line))
      continue;

    /* Process token */
    token = is_token(line);

    if (token) {
      curr_token = token;
      processedToken[curr_token]++;

      /* If this was a duplicate token, throw error and exit. */

      if (processedToken[curr_token] > 1){
	log_error("%s:%d Duplicate token of type %s found!\n");
	return FALSE; /* XXX memory LEAK. */
      }
      continue;
    }

    /* Line is valid of type "curr_token" */
    if (process_line(line, curr_token, initCon, onbootCon, serviceCon)
	== FALSE){
      log_error("%s:%d. Error processing the line\n", fileName, lineCount);
      goto done;
    }
  }

  ret = TRUE;

 done:

  free(line);
  fclose(fp);

  return ret;
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
  printf("--s, --startFile                      Startup File.\n");
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
  char *startupFile = DEF_STARTUP_FILE;

  /* Three type of containers in the startup file. */
  sType *initCon=NULL;
  sType *onbootCon=NULL;
  sType *serviceCon=NULL;
  
  /* Parsing command line args. */
  while (true) {
    int opt = 0;
    int opdidx = 0;

    static struct option long_options[] = {
      { "port",      required_argument, 0, 'p'},
      { "level",     required_argument, 0, 'l'},
      { "file",      required_argument, 0, 'f'},
      { "startFile", required_argument, 0, 's'},
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

    case 's':
      startupFile = optarg;
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

  /* Before we open the socket for REST, process the startup file and
   * start them containers.
   */
  initCon = (sType *)calloc(sizeof(sType), 1);
  onbootCon = (sType *)calloc(sizeof(sType), 1);
  serviceCon = (sType *)calloc(sizeof(sType), 1);

  if (initCon == NULL || onbootCon == NULL || serviceCon == NULL) {
    log_error("Memory allocation failure\n");
    exit(1);
  }

  if (process_startup_file(startupFile, initCon, onbootCon, serviceCon)
      != TRUE){
    log_error("Error processing the startup file\n");
    exit(1);
  }

  print_startupFile_info(initCon);
  print_startupFile_info(onbootCon);
  print_startupFile_info(serviceCon);

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
