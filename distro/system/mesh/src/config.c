/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Config.c
 *
 */

#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <stdio.h>

#include "mesh.h"
#include "config.h"
#include "toml.h"
#include "log.h"

static int parse_proxy_entries(Config *config, toml_table_t *proxyData);
static int parse_config_entries(int mode, int secure, Config *config,
				toml_table_t *configData);
static int read_line(char *buffer, int size, FILE *fp);

/*
 * print_config_data --
 *
 */

void print_config(Config *config) {

  if (config->mode == MODE_SERVER) {
    log_debug("Mode: SERVER");
  } else if (config->mode == MODE_CLIENT) {
    log_debug("Mode: CLIENT");
  }

  log_debug("Proxy: %s",  config->proxy  ? "enabled"         : "Disabled");
  log_debug("Secure: %s", config->secure ? "TLS/SSL enabled" : "Disabled");

  if (config->mode == MODE_SERVER) {
    log_debug("Remote accept port: %s", config->remoteAccept);
  } else if (config->mode == MODE_CLIENT) {
    log_debug("Remote connect port: %s", config->remoteConnect);
  }

  log_debug("Local accept port: %s", config->localAccept);

  if (config->secure) {
    log_debug("TLS/SSL key file: %s", config->keyFile);
    log_debug("TLS/SSL cert file: %s", config->certFile);
  }
}

/*
 * read_line -- read a line from file pointer.
 *
 */
static int read_line(char *buffer, int size, FILE *fp) {

  char *tmp;

  memset(buffer, 0, size);

  if (fgets(buffer, size, fp) == NULL) {
      *buffer = '\0';
      return FALSE;
  } else {
    /* remove newline */
    if ((tmp = strrchr(buffer, '\n')) != NULL) {
      *tmp = '\0';
    }
  }
  return TRUE;
}


/*
 * parse_proxy_entries -- handle reverse-proxy stuff.
 *
 */
static int parse_proxy_entries(Config *config, toml_table_t *proxyData) {

  toml_datum_t enable, httpPath, ip, port;

  enable = toml_string_in(proxyData, ENABLE);

  if (enable.ok) {
    if (strcasecmp(enable.u.s, "true")!=0) {
      config->reverseProxy = NULL;
      return TRUE;
    }
  } else {
    config->reverseProxy = NULL; /* disable by default. */
    return TRUE;
  }

  /* Will only come here if proxy is true. */
  httpPath = toml_string_in(proxyData, HTTP_PATH);
  ip       = toml_string_in(proxyData, CONNECT_IP);
  port     = toml_string_in(proxyData, CONNECT_PORT);

  if (!httpPath.ok && !ip.ok && !port.ok) {
    log_error("[%s] is missing required argument.", REVERSE_PROXY);
    return FALSE;
  }

  config->reverseProxy = (Proxy *)calloc(1, sizeof(Proxy));
  if (config->reverseProxy == NULL) {
    log_error("Error allocating memory of size: %s", sizeof(Proxy));
    return FALSE;
  }

  config->reverseProxy->enable    = TRUE;
  config->reverseProxy->httpPath = strdup(httpPath.u.s);
  config->reverseProxy->ip       = strdup(ip.u.s);
  config->reverseProxy->port     = strdup(port.u.s);

  free(httpPath.u.s);
  free(ip.u.s);
  free(port.u.s);
  if (enable.ok)
    free(enable.u.s);

  return TRUE;
}

/*
 * parse_config_entries -- Server/client stuff.
 *
 */

static int parse_config_entries(int mode, int secure, Config *config,
				toml_table_t *configData) {

  int ret=TRUE;
  FILE *fp;
  char *buffer=NULL;
  toml_datum_t remoteAccept, localAccept, remoteConnect, cert, key;
  toml_datum_t remoteIPFile;

  if (mode == MODE_SERVER) {
    remoteAccept = toml_string_in(configData, REMOTE_ACCEPT);
  } else if (mode == MODE_CLIENT) {
    remoteIPFile  = toml_string_in(configData, REMOTE_IP_FILE);
    remoteConnect = toml_string_in(configData, REMOTE_CONNECT);
  }

  localAccept = toml_string_in(configData, LOCAL_ACCEPT);

  cert  = toml_string_in(configData, CERT);
  key   = toml_string_in(configData, KEY);

  config->mode   = mode;
  config->secure = secure;

  if (config->mode == MODE_SERVER) {
    if (!remoteAccept.ok) {
      log_debug("[%s] is missing, setting to default: %s", REMOTE_ACCEPT,
		DEF_REMOTE_ACCEPT);
      config->remoteAccept = strdup(DEF_REMOTE_ACCEPT);
    } else {
      config->remoteAccept = strdup(remoteAccept.u.s);
    }
  }

  if (config->mode == MODE_CLIENT) {
    if (!remoteIPFile.ok) {
      log_debug("[%s] is missing. using default of 127.0.0.1", REMOTE_IP_FILE);
    } else {
      /* Read the content of the IP file. */
      fp = fopen(remoteIPFile.u.s, "r");
      if (fp == NULL) {
	log_error("[%s] Error opening file. Error: %s", remoteIPFile.u.s,
		  strerror(errno));
	ret = FALSE;
	goto done;
      }

      buffer = (char *)malloc(MAX_BUFFER);
      if (!buffer) {
	log_error("Error allocating memory of size: %s", MAX_BUFFER);
	ret = FALSE;
	goto done;
      }
      /* Read the file content. */
      if (read_line(buffer, MAX_BUFFER, fp)<=0) {
	log_error("[%s] Error reading file. Error: %s", remoteIPFile.u.s,
		  strerror(errno));
	ret = FALSE;
	free(buffer);
	buffer=NULL;
	goto done;
      }
    }

    if (!remoteConnect.ok || buffer == NULL) {
      log_debug("[%s] is missing, is mandatory", REMOTE_CONNECT);
      ret = FALSE;
      goto done;
    }

    config->remoteConnect = (char *)calloc(1, MAX_BUFFER);
    if (config->secure) {
      sprintf(config->remoteConnect, "wss://%s:%s/%s", buffer,
	      remoteConnect.u.s, PREFIX_WEBSOCKET);
    } else {
      sprintf(config->remoteConnect, "ws://%s:%s/%s", buffer,
	      remoteConnect.u.s, PREFIX_WEBSOCKET);
    }
    free(buffer); /* allocated by getline() */
  }

  if (!localAccept.ok) {
    log_debug("[%s] is missing, setting to default: %s", LOCAL_ACCEPT,
		DEF_LOCAL_ACCEPT);
    config->localAccept = strdup(DEF_LOCAL_ACCEPT);
  } else {
    config->localAccept = strdup(localAccept.u.s);
  }
  
  if (cert.ok) {
    config->certFile = strdup(cert.u.s);
  } else {
    config->certFile = strdup(DEF_SERVER_CERT);
  }

  if (key.ok) {
    config->keyFile = strdup(key.u.s);
  } else {
    config->keyFile = strdup(DEF_SERVER_KEY);
  }

 done:
  /* clear up toml allocations. */
  if (key.ok) free(key.u.s);
  if (cert.ok) free(cert.u.s);
  if (localAccept.ok) free(localAccept.u.s);
  if (remoteConnect.ok) free(remoteConnect.u.s);
  if (remoteIPFile.ok) free(remoteIPFile.u.s);
  if (remoteAccept.ok) free(remoteAccept.u.s);

  return ret;
}

/*
 * process_config_file -- read and parse the config file. 
 *                       
 *
 */
int process_config_file(int mode, int secure, int proxy, char *fileName,
			Config *config) {

  int ret=TRUE;
  FILE *fp;
  toml_table_t *fileData=NULL;
  toml_table_t *serverConfig=NULL, *clientConfig=NULL;
  toml_table_t *proxyConfig=NULL;
  
  char errBuf[MAX_BUFFER];

  /* Sanity check. */
  if (fileName == NULL || config == NULL)
    return FALSE;
  
  if ((fp = fopen(fileName, "r")) == NULL) {
    log_error("Error opening config file: %s: %s\n", fileName,
              strerror(errno));
    return FALSE;
  }

  /* Parse the TOML file entries. */
  fileData = toml_parse_file(fp, errBuf, sizeof(errBuf));
  
  fclose(fp);
 
  if (!fileData) {
    log_error("Error parsing the config file %s: %s\n", fileName, errBuf);
    return FALSE;
  }

  if (mode == MODE_SERVER) {

    serverConfig = toml_table_in(fileData, SERVER_CONFIG);

    if (serverConfig == NULL) {
      log_error("[%s] section parsing error in file: %s\n", SERVER_CONFIG,
		fileName);
      ret = FALSE;
      goto done;
    }
    ret = parse_config_entries(mode, secure, config, serverConfig);
    if (ret == FALSE) {
      goto done;
    }
  } else if (mode == MODE_CLIENT) {

    clientConfig = toml_table_in(fileData, CLIENT_CONFIG);

    if (clientConfig == NULL) {
      log_error("[%s] section parsing error in file: %s\n", CLIENT_CONFIG,
		fileName);
      ret = FALSE;
      goto done;
    }
    ret = parse_config_entries(mode, secure, config, clientConfig);
    if (ret == FALSE) {
      goto done;
    }
  }

  /* validate config entries for key and cert files. */
  if (secure) {
    if (config->certFile == NULL && config->keyFile == NULL) {
      ret = FALSE;
      goto done;
    }

    /* Make sure the cert and key are legit files. */
    if ((fp=fopen(config->certFile, "r")) == NULL) {
      log_error("Error with cert file: %s Error: %s", config->certFile,
		strerror(errno));
      ret = FALSE;
      goto done;
    }
    fclose(fp);

    if ((fp=fopen(config->keyFile, "r")) == NULL) {
      log_error("Error with key file: %s Error: %s", config->keyFile,
		strerror(errno));
      ret = FALSE;
      goto done;
    } 
    fclose(fp);
  }

  /* If proxies are enable */
  if (proxy) {

    proxyConfig = toml_table_in(fileData, REVERSE_PROXY);
    if (proxyConfig == NULL) {
      log_error("[%s] section parsing error in file: %s\n", SERVER_CONFIG,
		fileName);
      ret = FALSE;
      goto done;
    }
    ret = parse_proxy_entries(config, proxyConfig);
    if (ret == FALSE) {
      log_error("[%s] section parsing error in file: %s\n", REVERSE_PROXY,
		fileName);
      goto done;
    }
  }

 done:
  toml_free(fileData);
  return ret;
}

/*
 * clear_config --
 */

void clear_config(Config *config) {

  if (!config) return;

  if (config->mode == MODE_SERVER) {
    free(config->remoteAccept);
  }

  if (config->mode == MODE_CLIENT) {
    free(config->remoteConnect);
  }

  free(config->localAccept);
  free(config->certFile);
  free(config->keyFile);

  if (config->proxy) {
    free(config->reverseProxy->httpPath);
    free(config->reverseProxy->ip);
    free(config->reverseProxy->port);
    free(config->reverseProxy);
  }
}
