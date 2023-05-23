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
static int parse_config_entries(int secure, Config *config,
								toml_table_t *configData);
static int read_line(char *buffer, int size, FILE *fp);

/*
 * print_config --
 *
 */
void print_config(Config *config) {

	log_debug("Proxy: %s",  config->proxy  ? "enabled"         : "Disabled");
	log_debug("Secure: %s", config->secure ? "TLS/SSL enabled" : "Disabled");
	log_debug("Remote connect port: %s", config->remoteConnect);

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
 * split_strings --
 *
 */
static void split_strings(char *input, char **str1, char **str2) {

    char *delimiter=";", *token=NULL;

    token = strtok(input, delimiter);

    if (token != NULL) {
        *str1 = strdup(token);

        token = strtok(NULL, delimiter);
        if (token != NULL) {
            *str2 = strdup(token);
        }
    }
}

/*
 * read_hostname_and_nodeid -- read hostname (ip:port) and nodeID from the
 *                             passed file
 *
 */
static int read_hostname_and_nodeid(char *fileName, char **hostname,
                                    char **nodeID) {

    int ret=TRUE;
	FILE *fp=NULL;
	char *buffer=NULL;

	buffer = (char *)malloc(MAX_BUFFER);
	if (!buffer) {
		log_error("Error allocating memory of size: %s", MAX_BUFFER);
		return FALSE;
	}

	fp = fopen(fileName, "r");
	if (fp == NULL) {
		log_error("[%s] Error opening file. Error: %s", fileName,
				  strerror(errno));
		return FALSE;
	}

	/* Read the file content. */
	if (read_line(buffer, MAX_BUFFER, fp)<=0) {
		log_error("[%s] Error reading file. Error: %s", fileName,
				  strerror(errno));
        ret = FALSE;
	} else {
        split_strings(buffer, hostname, nodeID);
    }

	fclose(fp);
    free(buffer);

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

static int parse_config_entries(int secure, Config *config,
								toml_table_t *configData) {

	int ret=TRUE;
	char *hostname=NULL, *nodeID=NULL;
	toml_datum_t localAccept, cert, key;
	toml_datum_t remoteIPFile;

	remoteIPFile  = toml_string_in(configData, REMOTE_IP_FILE);
	localAccept   = toml_string_in(configData, LOCAL_ACCEPT);
	cert          = toml_string_in(configData, CFG_CERT);
	key           = toml_string_in(configData, CFG_KEY);

	config->secure = secure;

	if (!remoteIPFile.ok) {
		log_error("[%s] is missing but is mandatory", REMOTE_IP_FILE);
        ret=FALSE;
        goto done;
	} else {
		/* Read the content of the IP file. */
		if (read_hostname_and_nodeid(remoteIPFile.u.s, &hostname, &nodeID)
            == FALSE) {
			goto done;
		}
	}

	config->remoteConnect = (char *)calloc(1, MAX_BUFFER);
	if (config->secure) {
		sprintf(config->remoteConnect, "wss://%s/%s", hostname,
                PREFIX_WEBSOCKET);
	} else {
		sprintf(config->remoteConnect, "ws://%s/%s", hostname,
				PREFIX_WEBSOCKET);
	}

	config->deviceInfo = (DeviceInfo *)malloc(sizeof(DeviceInfo));
	if (config->deviceInfo == NULL) {
		log_error("Error allocating memory of size: %d", sizeof(DeviceInfo));
		goto done;
	}
    config->deviceInfo->nodeID = strdup(nodeID);

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
	if (key.ok)           free(key.u.s);
	if (cert.ok)          free(cert.u.s);
	if (localAccept.ok)   free(localAccept.u.s);
	if (remoteIPFile.ok)  free(remoteIPFile.u.s);
    if (hostname)         free(hostname);
    if (nodeID)           free(nodeID);

	return ret;
}

/*
 * process_config_file -- read and parse the config file. 
 *                       
 *
 */
int process_config_file(int secure, int proxy, char *fileName, Config *config) {

	int ret=TRUE;
	FILE *fp;
	toml_table_t *fileData=NULL;
	toml_table_t *clientConfig=NULL;
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

	clientConfig = toml_table_in(fileData, CLIENT_CONFIG);

	if (clientConfig == NULL) {
		log_error("[%s] section parsing error in file: %s\n", CLIENT_CONFIG,
				  fileName);
		ret = FALSE;
		goto done;
	}
	ret = parse_config_entries(secure, config, clientConfig);
	if (ret == FALSE) {
		goto done;
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
			log_error("[%s] section parsing error in file: %s\n", REVERSE_PROXY,
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

	free(config->remoteConnect);
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
