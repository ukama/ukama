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

#include <curl/curl.h>
#include <curl/easy.h>

#include "mesh.h"
#include "config.h"
#include "toml.h"
#include "log.h"

static int parse_proxy_entries(Config *config, toml_table_t *proxyData);
static int parse_config_entries(int mode, int secure, Config *config,
				toml_table_t *configData);
static int read_line(char *buffer, int size, FILE *fp);
static int parse_amqp_config(Config *config, toml_table_t *configData);
static int is_valid_url(char *name, char *port);

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
    log_debug("AMQP host: %s:%s", config->amqpHost, config->amqpPort);
    log_debug("AMQP exchange: %s", config->amqpExchange);
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
 * read_ip -- read IP (hostname) from the passed fileName
 *
 */

static char *read_ip(char *fileName) {

  FILE *fp=NULL;
  char *buffer=NULL;

  buffer = (char *)malloc(MAX_BUFFER);
  if (!buffer) {
    log_error("Error allocating memory of size: %s", MAX_BUFFER);
    return NULL;
  }

  fp = fopen(fileName, "r");
  if (fp == NULL) {
    log_error("[%s] Error opening file. Error: %s", fileName, strerror(errno));
    return NULL;
  }

  /* Read the file content. */
  if (read_line(buffer, MAX_BUFFER, fp)<=0) {
    log_error("[%s] Error reading file. Error: %s", fileName, strerror(errno));
    fclose(fp);
    free(buffer);
    return NULL;
  }

  fclose(fp);
  return buffer;
}

/*
 * valid_url -- validate URL via CURL request.
 */
static int is_valid_url(char *name, char *port) {

  char url[MAX_BUFFER]={0};
  CURL *curl;
  CURLcode response;

  if (name==NULL && port==NULL) {
    return FALSE;
  }

  sprintf(url, "%s:%s", name, port);

  curl = curl_easy_init();

  if (curl) {

    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_NOBODY, 1);

    response = curl_easy_perform(curl);

    curl_easy_cleanup(curl);
  }

  if (response != CURLE_WEIRD_SERVER_REPLY) { /* AMQP reply in binary */
    log_error("Error reaching AMQP at %s. Error: %s", url,
	      curl_easy_strerror(response));
    return FALSE;
  }

  return TRUE;
}

/*
 * parse_amqp_config --
 *
 */
static int parse_amqp_config(Config *config, toml_table_t *configData) {

  int ret=TRUE;
  toml_datum_t amqpHost, amqpPort, amqpExchange;

  if (config == NULL && configData == NULL) {
    return FALSE;
  }

  amqpHost     = toml_string_in(configData, AMQP_HOST);
  amqpPort     = toml_string_in(configData, AMQP_PORT);
  amqpExchange = toml_string_in(configData, AMQP_EXCHANGE);

  if (!amqpHost.ok) {
    log_debug("[%s] is missing but is mandatory.", AMQP_HOST);
    ret = FALSE;
    goto done;
  } else {
    config->amqpHost = strdup(amqpHost.u.s);
  }

  if (!amqpPort.ok) {
    log_debug("[%s] is missing but is mandatory.", AMQP_PORT);
    ret = FALSE;
    goto done;
  } else {
    config->amqpPort = strdup(amqpPort.u.s);
  }

  /* Check to see if the host:port is reachable. */
  if (is_valid_url(config->amqpHost, config->amqpPort)==FALSE) {
    log_error("Invalid host and/or port for AMQP. Host: %s Port %s",
	      config->amqpHost, config->amqpPort);
    ret = FALSE;
    goto done;
  }

  if (!amqpExchange.ok) {
    log_debug("[%s] is missing but is mandatory.", AMQP_EXCHANGE);
    ret = FALSE;
    goto done;
  } else {
    config->amqpExchange = strdup(amqpExchange.u.s);
  }

 done:
  if (amqpHost.ok) free(amqpHost.u.s);
  if (amqpPort.ok) free(amqpPort.u.s);
  if (amqpExchange.ok) free(amqpExchange.u.s);

  return ret;
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

  cert  = toml_string_in(configData, CFG_CERT);
  key   = toml_string_in(configData, CFG_KEY);

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

    /* Setup AMQP parameters. */
    if (parse_amqp_config(config, configData)==FALSE) {
      ret = FALSE;
      goto done;
    }
  }

  if (config->mode == MODE_CLIENT) {
    if (!remoteIPFile.ok) {
      log_debug("[%s] is missing. using default of 127.0.0.1", REMOTE_IP_FILE);
    } else {
      /* Read the content of the IP file. */
      buffer = read_ip(remoteIPFile.u.s);
      if (buffer == NULL) {
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
  if (mode == MODE_CLIENT) {
    if (remoteConnect.ok) free(remoteConnect.u.s);
    if (remoteIPFile.ok) free(remoteIPFile.u.s);
  }
  if (mode == MODE_SERVER) {
    if (remoteAccept.ok) free(remoteAccept.u.s);
  }
  if (buffer) free(buffer);

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
    free(config->amqpHost);
    free(config->amqpPort);
    free(config->amqpExchange);
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
