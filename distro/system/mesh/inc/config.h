/*
 * config.h
 */

#ifndef MESH_CONFIG_H
#define MESH_CONFIG_H

#include <uuid/uuid.h>

/* used in the config file and for parsing. */
#define SERVER_CONFIG "server-config"
#define CLIENT_CONFIG "client-config"
#define REVERSE_PROXY "reverse-proxy"

#define REMOTE_ACCEPT  "remote-accept"
#define LOCAL_ACCEPT   "local-accept"
#define REMOTE_CONNECT "remote-connect"
#define CONNECT_IP     "connect-ip"
#define CONNECT_PORT   "connect-port"
#define HTTP_PATH      "http-path"

#define CERT           "cert"
#define KEY            "key"
#define ENABLE         "enable"

/* Some default */
#define DEF_REMOTE_ACCEPT  "5534"
#define DEF_LOCAL_ACCEPT   "5533"
/* default for client. */
#define DEF_REMOTE_CONNECT "ws://localhost:5534/"
#define DEF_REMOTE_SECURE_CONNECT "wss://localhost:5534/"

#define MODE_SERVER_STR "server"
#define MODE_CLIENT_STR "client"
#define MODE_DUAL_STR   "dual"

#define MODE_SERVER 1
#define MODE_CLIENT 2
#define MODE_DUAL   3

#define MAX_BUFFER 256

#define DEF_SERVER_CERT "certs/test.cert"
#define DEF_SERVER_KEY  "cert/server.key"

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

  int mode;             /* client or server. */
  int secure;           /* enable SSL/TLS for remote accept */
  int proxy;            /* reverse-proxy enabled (true | false) */

  char *remoteAccept;   /* Server: Port on which to accept remote clients */
  char *localAccept;    /* Both: Port on which to accept local clients */
  char *remoteConnect;  /* Client: hostname:port to connect with remotely */

  char *certFile;       /* CA Cert file name. */
  char *keyFile;        /* Key file name.*/
  uuid_t uuid;          /* Device UUID. */

  Proxy *reverseProxy;         /* define any reverse proxy */
} Config;

int process_config_file(int mode, int secure, int proxy, char *fileName,
			Config *config);
void clear_config(Config *config);
void print_config(Config *config);

#endif /* MESH_CONFIG_H */
