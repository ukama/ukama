/*
 * Config.c
 *
 */

#include <stdlib.h>
#include <string.h>
#include <errno.h>

#include "config.h"
#include "toml.h"
#include "log.h"

/*
 * set_default_base_config --
 *
 */

static void set_default_base_config(BaseConfig *base) {
  
  base->admin = FALSE;
  base->maxRemoteClients = DEF_MAX_REMOTE_CLIENTS;
  base->maxLocalClients  = DEF_MAX_LOCAL_CLIENTS;
}

/*
 * prase_base_config --
 *
 */

static int parse_base_config(BaseConfig *base, toml_table_t *configData) {
  
  int ret=FALSE;
  toml_datum_t mode, admin, remoteClients, localClients;
  
  set_default_base_config(base);
  
  /* Read the config data from the config.toml and load into Config. */
  mode = toml_string_in(configData, MODE);
  admin = toml_string_in(configData, ADMIN);
  remoteClients = toml_string_in(configData, REMOTE_CLIENTS);
  localClients = toml_string_in(configData, LOCAL_CLIENTS);

  /* Mode is mandatory. For others default is used if not defined. */
  if (!mode.ok) {
    log_error("[%s] is mandatory but missing.", MODE);
    return FALSE;
  } else {
    if (!strcmp(MODE_SERVER_STR, mode.u.s)) {
      base->mode = MODE_SERVER;
    } else if (!strcmp(MODE_CLIENT_STR, mode.u.s)) {
      base->mode = MODE_CLIENT;
    } else if (!strcmp(MODE_DUAL_STR, mode.u.s)) {
      base->mode = MODE_DUAL;
    } else {
      log_error("[%s] invalid entry: %s", MODE, mode.u.s);
      goto exit;
    }
  }

  /* Admin. */
  if (admin.ok) {
    if (!strcmp(admin.u.s, "enable")) {
      base->admin = TRUE;
    } else if (!strcmp(admin.u.s, "disable")) {
      base->admin = FALSE;
    } else {
      log_error("[%s], invalid entry for %s. Setting to default", BASE_CONFIG,
		ADMIN);
    }
  }

  /* remote and local clients. */
  if (remoteClients.ok) {
    int val = atoi(remoteClients.u.s);

    if (val > MAX_REMOTE_CLIENTS || val < MIN_REMOTE_CLIENTS) {
      log_error("[%s], invalid range for: %s. Setting to default", BASE_CONFIG,
		REMOTE_CLIENTS);
      base->maxRemoteClients = DEF_MAX_REMOTE_CLIENTS;
    } else {
      base->maxRemoteClients = val;
    }
  }

  if (localClients.ok) {
    int val = atoi(localClients.u.s);

    if (val > MAX_LOCAL_CLIENTS || val < MIN_LOCAL_CLIENTS) {
      log_error("[%s], invalid range for: %s. Setting to default", BASE_CONFIG,
		LOCAL_CLIENTS);
      base->maxLocalClients = DEF_MAX_LOCAL_CLIENTS;
    } else {
      base->maxLocalClients = val;
    }
  }

  ret = TRUE;
  
 exit:
  return ret;
}

/*
 * prase_admin_config --
 *
 */

static int parse_admin_config(AdminConfig *admin, toml_table_t *configData) {
  
  toml_datum_t adminEP, statsEP, port;
  
  /* Read the config data from the config.toml and load into Config. */
  adminEP = toml_string_in(configData, ADMIN_ENDPOINT);
  statsEP = toml_string_in(configData, STATS_ENDPOINT);
  port = toml_string_in(configData, ADMIN_PORT);

  if (!adminEP.ok) {
    admin->adminEP = DEF_ADMIN_EP;
  } else {
    admin->adminEP = strdup(adminEP.u.s);
  }

  if (!statsEP.ok) {
    admin->statsEP = DEF_STATS_EP;
  } else {
    admin->statsEP = strdup(statsEP.u.s);
  }

  if (!port.ok) {
    admin->port = DEF_ADMIN_PORT;
  } else { 
    admin->port = atoi(port.u.s);
  }
  
  return TRUE;
}

/*
 * prase_host_config -- Server/client stuff.
 *
 */

static int parse_host_config(HostConfig *config, toml_table_t *configData,
			     int hFlag) {

  toml_datum_t localHost, localPort, remoteHost, remotePort, cert, key, proxy;

  localHost = toml_string_in(configData, LOCAL_HOST);
  localPort = toml_string_in(configData, LOCAL_PORT);
  
  remoteHost = toml_string_in(configData, REMOTE_HOST);
  remotePort = toml_string_in(configData, REMOTE_PORT);
  
  cert  = toml_string_in(configData, CERT);
  key   = toml_string_in(configData, KEY);
  proxy = toml_string_in(configData, PROXY);


  if (hFlag == 1) { /* Only server have local port, for now. XXX */

    config->type = 1;

    if (!localHost.ok) {
      log_debug("[%s] is missing, setting to default: localhost", LOCAL_HOST);
      config->localHostname = LOCALHOST;
    } else {
      config->localHostname = strdup(localHost.u.s);
    }

    if (!localPort.ok) {
      log_debug("[%s] is missing, setting to default: %d", LOCAL_PORT,
		DEF_LOCAL_PORT);
      config->localPort = DEF_LOCAL_PORT;
    } else {
      config->localPort = atoi(localPort.u.s);
    }
  } else {
    config->type = 2;
  }

  /* Common for both server and client. */
  
  /* Remote server/port. */
  if (!remoteHost.ok) {
    log_debug("[%s] is missing, setting to default: localhost", REMOTE_HOST);
    config->remoteHostname = LOCALHOST;
  } else {
    config->remoteHostname = strdup(remoteHost.u.s);
  }
  
  if (!remotePort.ok) {
    log_debug("[%s] is missing, setting to default: %s", REMOTE_PORT,
	      DEF_REMOTE_PORT);
    config->remotePort = DEF_REMOTE_PORT;
  } else {
    config->remotePort = strdup(remotePort.u.s);
  }
  
  /* Cert, key and proxy. */
  if (cert.ok) {
    config->certFile = strdup(cert.u.s);
  } else {
    config->certFile = DEF_SERVER_CERT;
  }

  if (key.ok) {
    config->keyFile = strdup(key.u.s);
  } else {
    config->keyFile = DEF_SERVER_KEY;
  }
  
  if (proxy.ok) {
    if (!strcmp(proxy.u.s, PROXY_NONE_STR)) {
      config->proxyType = PROXY_NONE;
    } else {
      log_error("%s: %s is not supported!", PROXY, proxy.u.s);
      return FALSE;
    }
  } else {
    config->proxyType = PROXY_NONE; /* default is no proxy. */
  }
  
  return TRUE;
}

/*
 * process_config_file -- read and parse the config file. 
 *                       
 *
 */
int process_config_file(char *fileName, Configs *config) {

  FILE *fp;

  toml_table_t *fileData;
  toml_table_t *baseConfig, *adminConfig, *serverConfig, *clientConfig;
  
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

  /* Parse the base-config. */
  baseConfig = toml_table_in(fileData, BASE_CONFIG);

  if (baseConfig == NULL) {
    log_error("[%s] section parsing error in file: %s\n", BASE_CONFIG, fileName);
    toml_free(fileData);
    return FALSE;
  }

  config->baseConfig = (BaseConfig *)calloc(sizeof(BaseConfig), 1);
  if (!config->baseConfig) {
    log_error("Error allocating memory", sizeof(BaseConfig));
    return FALSE;
  }
  
  parse_base_config(config->baseConfig, baseConfig);

  /* Parse the admin-config if admin=enable in [base-config] */
  if (config->baseConfig->admin) { 
  
    adminConfig = toml_table_in(fileData, ADMIN_CONFIG);

    if (adminConfig == NULL) {
         log_error("[%s] section parsing error in file: %s\n", ADMIN_CONFIG,
		   fileName);
	 toml_free(fileData);
	 return FALSE;
    }
    config->adminConfig = (AdminConfig *)calloc(sizeof(AdminConfig), 1);
    if (!config->adminConfig) {
      log_error("Error allocating memory", sizeof(AdminConfig));
      return FALSE;
    }
    parse_admin_config(config->adminConfig, adminConfig);
  } else {
    config->adminConfig = NULL;
  }

  /* Parse the server-config and client-config. */
  if (config->baseConfig->mode == MODE_SERVER ||
      config->baseConfig->mode == MODE_DUAL) { 
  
    serverConfig = toml_table_in(fileData, SERVER_CONFIG);

    if (serverConfig == NULL) {
      log_error("[%s] section parsing error in file: %s\n", SERVER_CONFIG, fileName);
      toml_free(fileData);
      return FALSE;
    }
    
    config->serverConfig = (HostConfig *)calloc(sizeof(HostConfig), 1);
    if (!config->serverConfig) {
      log_error("Error allocating memory", sizeof(HostConfig));
      return FALSE;
    }
    parse_host_config(config->serverConfig, serverConfig, 1);
  } else {
    config->serverConfig = NULL;
  }

  if (config->baseConfig->mode == MODE_CLIENT ||
      config->baseConfig->mode == MODE_DUAL) { 
    
    clientConfig = toml_table_in(fileData, CLIENT_CONFIG);

    if (clientConfig == NULL) {
      log_error("[%s] section parsing error in file: %s\n", CLIENT_CONFIG, fileName);
      toml_free(fileData);
      return FALSE;
    }
    
    config->clientConfig = (HostConfig *)calloc(sizeof(HostConfig), 1);
    if (!config->clientConfig) {
      log_error("Error allocating memory", sizeof(HostConfig));
      return FALSE;
    }
    parse_host_config(config->clientConfig, clientConfig, 2);
  } else {
    config->clientConfig = NULL;
  }
  
  toml_free(fileData);
  return TRUE;
}
