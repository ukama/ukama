/*
 * config.h
 */

#ifndef CONFIG_H
#define CONFIG_H

/* used in the config file and for parsing. */
#define SERVER_CONFIG "server-config"
#define CLIENT_CONFIG "client-config"
#define REVERSE       "reverse-proxy"
#define REMOTE_ACCEPT "remote-accept"
#define LOCAL_ACCEPT  "local-accept"
#define REMOTE_CONNECT "remote-connect"
#define CERT           "cert"
#define KEY            "key"

/* Some default */
#define DEF_REMOTE_ACCEPT  "5534"
#define DEF_REMOTE_CONNECT "5534"
#define DEF_LOCAL_ACCEPT   "5533"

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

/* Struct to define the server and/or client host cfg. */
typedef struct {

  int mode;             /* client or server. */
  int secure;           /* enable SSL/TLS for remote accept */
  
  char *remoteAccept;   /* Server: Port on which to accept remote clients */
  char *localAccept;    /* Both: Port on which to accept local clients */
  char *remoteConnect;  /* Client: hostname:port to connect with remotely */

  char *certFile;       /* CA Cert file name. */
  char *keyFile;        /* Key file name.*/ 
} Config;

int process_config_file(int mode, int secure, char *fileName, Config *config);
void clear_config(Config *config);
void print_config(Config *config);

#endif /* CONFIG_H */
