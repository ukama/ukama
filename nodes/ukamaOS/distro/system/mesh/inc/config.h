/*
 * config.h
 */

#ifndef MESH_CONFIG_H
#define MESH_CONFIG_H

#include <uuid/uuid.h>

#include "mesh.h"

/* used in the config file and for parsing. */
#define CLIENT_CONFIG "client-config"
#define REVERSE_PROXY "reverse-proxy"

#define LOCAL_ACCEPT   "local-accept"
#define REMOTE_IP_FILE "remote-ip-file"
#define CONNECT_IP     "connect-ip"
#define CONNECT_PORT   "connect-port"
#define HTTP_PATH      "http-path"

#define CFG_CERT "cert"
#define CFG_KEY  "key"
#define ENABLE   "enable"

/* Some default */
#define DEF_LOCAL_ACCEPT   "5533"
#define DEF_REMOTE_CONNECT "ws://localhost:5534/"
#define DEF_REMOTE_SECURE_CONNECT "wss://localhost:5534/"

#define MODE_CLIENT_STR "client"

#define MAX_BUFFER 256

#define DEF_SERVER_CERT "certs/test.cert"
#define DEF_SERVER_KEY  "cert/server.key"
#define DEF_PORT        "9275"

#define TRUE 1
#define FALSE 0

/* Struct to define the reverse proxies  */
typedef struct {

	int  enable;

	char *httpPath;
	char *ip;
	char *port;
} Proxy;

/* Struct to define the server and/or client host cfg. */
typedef struct {

	int secure;           /* enable SSL/TLS for remote accept */
	int proxy;            /* reverse-proxy enabled (true | false) */

	char *localAccept;    /* Both: Port on which to accept local clients */
	char *remoteConnect;  /* Client: hostname:port to connect with remotely */

	char *certFile;       /* CA Cert file name. */
	char *keyFile;        /* Key file name.*/

	DeviceInfo *deviceInfo;   /* Device related info. */
	Proxy      *reverseProxy; /* define any reverse proxy */
} Config;

int process_config_file(int secure, int proxy, char *fileName, Config *config);
void clear_config(Config *config);
void print_config(Config *config);

#endif /* MESH_CONFIG_H */
