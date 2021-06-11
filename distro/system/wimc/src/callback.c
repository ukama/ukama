/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Callback functions for various endpoints and REST methods.
 */

#include "callback.h"
#include "wimc.h"

/*
 * decode a u_map into a string
 */

static char *print_map(const struct _u_map * map) {

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
 * validate_path --
 *
 */

static int validate_path(char *path) {

  int val=FALSE;
  struct stat sb;
  char buf[256];

  sprintf(buf, "%s/config.json", path); 
  
  if (stat(path, &sb) == 0 && S_ISDIR(sb.st_mode)) {
    /* Now check for config.json file. */
    if (access(&buf[0], F_OK) == 0) {
      log_debug("Config.json found at: %s", path);
      val = TRUE;
    } else {
      log_debug("Valid path but config.json NOT found!");
    }      
  } else {
    log_debug("Invalid path: %s", path);
  }
  
  return val;
}

/*
 * container_name_and_tag -- 
 *
 */
static int container_name_and_tag(char *str, char *name, char *tag) { 

  char *c, *t;

  c = strtok(str, ":");
  t = strtok(NULL, ":");

  if (c == NULL || t == NULL) {
    goto failure;
  }

  /* sanity check. */
  if (strlen(c) > WIMC_MAX_NAME_LEN ||
      strlen(t) > WIMC_MAX_TAG_LEN) {
    goto failure;
  }

  strncpy(name, c, strlen(c));
  strncpy(tag, t, strlen(t));

  return TRUE;
 failure:
  return FALSE;
}

/* 
 * get_key_value -- extract key value pair from the passed string. It is in
 *                  format key=value. If keyName is valid, also compared it 
 *                  against it.
 *
 */

static int get_key_value(char *str, char *key, char *value, char *keyName,
			 int maxKeyLen) {

  int val = FALSE;
  char *k, *v;
  
  k = strtok(str, "=");
  v = strtok(NULL, "=");
  
  if (k == NULL || v == NULL) {
    goto failure;
  }
  
  if ( keyName != NULL && maxKeyLen > 0) { /* Validate. */
    if ((strcasecmp(keyName, k)==0) || (strlen(v) < maxKeyLen)) {
      val = TRUE;
    }
  }
  
  /* We still copy even if the keyname doesn't match so we can log this. */
  strncpy(value, v, strlen(v));
  strncpy(key, k, strlen(k));
  
 failure:
  return val;
}

/*
 * get_container_name -- return the requested container name from the url.
 *
 */

static int get_container_name(char *http, char *name, char *tag) {

  int val=FALSE;
  char *token1, *token2;

  token1 = strtok(http, "?"); /* will return /container */
  token2 = strtok(NULL, "?");
  
  if (token1 == NULL || token2 == NULL) {
    goto failure;
  }		  
  
  if (strcasecmp(WIMC_EP_CONTAINER, token1) == 0){
    val = container_name_and_tag(token2, name, tag); /* name:tag */

    if (val == FALSE) {
      goto failure;
    }
  }

 failure:
  return val;
}

/*
 * get_post_params -- get params for the POST method.
 *                    param is of format name=container:tag&path=/some/path
 *
 */

static int get_post_params(char *params, char *name, char *tag, char *path) {

  int val=FALSE;
  char *token1, *token2;
  char key[256]={0}, value[256]={0};

  /* Extract the first token: container name and tag. */
  token1 = strtok(params, "&");
  token2 = strtok(NULL, "&"); 
  
  if (token1 == NULL || token2 == NULL) {
    goto failure;
  }
  
  val = get_key_value(token1, &key[0], &value[0], WIMC_PARAM_CONTAINER_NAME,
		      WIMC_MAX_NAME_LEN + WIMC_MAX_PATH_LEN + 1);
  
  if (val == TRUE) {
    /* Extract the container name and its tag. It is in format
     * container_name:tag 
     */
    val = container_name_and_tag(value, name, tag);
    if (!val) {
      log_debug("Invalid container name and/or tag. Ignoring");
      goto failure;
    }
  } else {
    log_debug("Invalid key %s. Expected %s", key, WIMC_PARAM_CONTAINER_NAME);
    goto failure;
  }
  
  /* Process the path. */
  memset(key, 0, sizeof(char)*256);
  memset(value, 0, sizeof(char)*256);
  
  val = get_key_value(token2, &key[0], &value[0], WIMC_PARAM_CONTAINER_PATH,
		      WIMC_MAX_PATH_LEN);
  if (val == TRUE){
    strncpy(path, value, strlen(value));
  } else {
    goto failure;
  }

  log_debug("Valid POST params recevied. Name: %s Tag: %s Path: %s",
	    name, tag, path);

 failure:
  return val;
}

/* 
 * log_request -- log various parameters for the incoming request. 
 *
 */

static void log_request(const struct _u_request *request) {

  log_debug("Recevied: %s %s %s", request->http_protocol, request->http_verb,
	    request->http_url);
}


/*
 * callback_get_container --
 *
 */

/* 
   {http_protocol = 0x7fffec000ee0 "HTTP/1.1", 
   http_verb = 0x7fffec000f00 "GET", 
   http_url = 0x7fffec000ea0 "/container?container_name:tag", 
   url_path = 0x7fffec000ec0 "/container", 
*/

int callback_get_container(const struct _u_request *request,
			   struct _u_response *response,
			   void *user_data) {

  int val=FALSE, found=FALSE;
  
  char name[WIMC_MAX_NAME_LEN] = {0};
  char tag[WIMC_MAX_TAG_LEN]  = {0};
  char path[WIMC_MAX_PATH_LEN] = {0};
  char *params, *post_params, *response_body;
  WimcCfg *cfg;
  
  cfg = (WimcCfg *)user_data;
  
  log_request(request);
  params = (char *)request->http_url;
  
  val = get_container_name(params, &name[0], &tag[0]);

  if (val) {
    log_debug("Valid GET request for %s: name[%s] tag[%s]",
	      WIMC_EP_CONTAINER, name, tag);
    if (db_read_path(cfg->db, name, tag, &path[0])) {
      /* we have the contents stored locally. Return the location. */

      /* XXX - TODO. */
    } else {

      /* Ask cloud based service provider for the content. 
       * dB will be kept updated with the status.
       * Special CB URL will be return to the callee. 
       */
      
      fetch_content_from_service_provider(cfg, name, tag, WIMC_TYPE_CONTAINER);
    }
  } else {
    log_debug("Invalid GET request received for %s. Ignoring!",
	      WIMC_EP_CONTAINER);
    /* XXX */
  }
  
  post_params = print_map(request->map_post_body);
  response_body = msprintf("OK!\n%s", post_params);

  ulfius_set_string_body_response(response, 200, response_body);
  o_free(response_body);
  o_free(post_params);

  return U_CALLBACK_CONTINUE;
}

/*
 * callback_post_container --
 *
 */
int callback_post_container(const struct _u_request *request,
			    struct _u_response *response,
			    void *user_data) {
  int ret;
  char *params, *post_params, *response_body;
  WimcCfg *cfg;

  cfg = (WimcCfg *)user_data;
    
  char name[WIMC_MAX_NAME_LEN] = {0};
  char tag[WIMC_MAX_TAG_LEN]   = {0};
  char path[WIMC_MAX_PATH_LEN] = {0};

  log_request(request);
  
  params = (char *)request->binary_body;

  ret = get_post_params(params, &name[0], &tag[0], &path[0]);

  /* Validate the path is correct: 
   * 1. path exists and 2. config.json exists in the root.
   */
  if (validate_path(path)) {
    /* Add entry to the db and return OK. */
    if (db_insert_entry(cfg->db, name, tag, path)) {
      response_body = msprintf("OK!\n%s", post_params);
    } else {
      response_body = msprintf("Error inserting into db!\n%s", post_params);
    }
  } else {
    response_body = msprintf("Invalid path. Ignoring\n%s", post_params);
  }
  
  post_params = print_map(request->map_post_body);

  if (!ret) {
    log_debug("Invalid POST request recevied for %s. Ignoring!",
	      WIMC_EP_ADMIN);
  }
  
  ulfius_set_string_body_response(response, 200, response_body);
  o_free(response_body);
  o_free(post_params);

  return U_CALLBACK_CONTINUE;
}

/*
 * callback_put_container --
 *
 */
int callback_put_container(const struct _u_request *request,
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
 * callback_delete_container --
 *
 */
int callback_delete_container(const struct _u_request *request,
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
 * callback_get_stats --
 *
 */
int callback_get_stats(const struct _u_request *request,
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
 * callback_not_allowed -- 
 *
 */
int callback_not_allowed(const struct _u_request *request,
			 struct _u_response *response, void *user_data) {

  ulfius_set_string_body_response(response, 403, "Operation not allowed\n");
  return U_CALLBACK_CONTINUE;
}

/*
 * callback_default -- default callback for no-match
 *
 */
int callback_default(const struct _u_request *request,
                     struct _u_response *response, void *user_data) {

  ulfius_set_string_body_response(response, 404, "You are clearly high!\n");
  return U_CALLBACK_CONTINUE;
}
