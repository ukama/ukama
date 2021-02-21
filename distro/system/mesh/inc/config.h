/*
 * config.h
 */

#ifndef CONFIG_H
#define CONFIG_H

/* used in the config file and for parsing. */
#define BASE_CONFIG  "base-config"
#define ADMIN_CONFIG  "admin-config"
#define SERVER_CONFIG "server-config"
#define CLIENT_CONFIG "client-config"
#define REVERSE       "reverse-proxy"

/* defs for base-config. */
#define MODE              "mode"
#define ADMIN             "admin"
#define REMOTE_CLIENTS    "remote-clients"
#define LOCAL_CLIENTS     "local-clients"

/* defs for admin-config. */
#define ADMIN_ENDPOINT "admin-endpoint"
#define STATS_ENDPOINT "stats-endpoint"
#define ADMIN_PORT     "admin-port"

#define DEF_ADMIN_EP "/admin"
#define DEF_STATS_EP "/stats"
#define DEF_ADMIN_PORT 5533

/* defs for server-config and client-config */
#define LOCAL_HOST  "local-host"
#define LOCAL_PORT  "local-port"
#define REMOTE_HOST "remote-host"
#define REMOTE_PORT "remote-port"
#define CERT        "cert"
#define KEY         "key"
#define PROXY       "proxy"

#define LOCALHOST   "localhost"
#define DEF_LOCAL_PORT  5534
#define DEF_REMOTE_PORT "5535" /* String due to mbedTLS lib. */

/* defs related to proxy. */
#define PROXY_NONE    0x01
#define PROXY_FORWARD 0x02
#define PROXY_REVERSE 0x04
#define PROXY_NONE_STR    "none"
#define PROXY_FORWARD_STR "forward"
#define PROXY_REVERSE_STR "reverse"
#define PROXY_ALL_STR     "all"

#define MODE_SERVER_STR "server"
#define MODE_CLIENT_STR "client"
#define MODE_DUAL_STR   "dual"

#define MODE_SERVER 1
#define MODE_CLIENT 2
#define MODE_DUAL   3

#define DEF_SERVER_CERT "certs/test.cert"
#define DEF_SERVER_KEY  "cert/server.key"

#define MAX_BUFFER 256

#define MAX_LOCAL_CLIENTS  1000
#define MAX_REMOTE_CLIENTS 1000

#define MIN_LOCAL_CLIENTS  1
#define MIN_REMOTE_CLIENTS 1

#define DEF_MAX_LOCAL_CLIENTS  50
#define DEF_MAX_REMOTE_CLIENTS 50


#define TRUE 1
#define FALSE 0

typedef struct {
  
  int mode;             /* Mode of mesh.d */
  int admin;            /* Flag to enable or disable admin port. */  
  int maxRemoteClients; /* Max number of remote client, (server or dual mode) */
  int maxLocalClients;  /* Max number of local client (device or cluster) */
}BaseConfig;

typedef struct {

  char *adminEP; /* REST end-point for the admin. */
  char *statsEP; /* REST end-point for prometheous stats collector. */
  int port;      /* admin port. */
}AdminConfig;


/* Struct to define the server and/or client host cfg. */
typedef struct {

  int type;             /* client or server. */
  
  char *localHostname;  /* Binding local host. */
  int localPort;        /* port for local clients. */
  char *remoteHostname; /* Binding remote host. */
  char *remotePort;     /* port for remote clients (devices) */
  char *certFile;       /* CA Cert file name. */
  char *keyFile;        /* Key file name.*/ 
  int proxyType;        /* none, forward, reverse or both. */
}HostConfig;

typedef struct {
  
  BaseConfig   *baseConfig;
  AdminConfig  *adminConfig;
  HostConfig   *serverConfig;
  HostConfig   *clientConfig;

  /* XXX - reverse proxy. */
}Configs;


#endif /* CONFIG_H */
