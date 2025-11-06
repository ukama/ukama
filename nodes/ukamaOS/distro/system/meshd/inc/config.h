#ifndef MESH_CONFIG_H
#define MESH_CONFIG_H

#include "mesh.h"

/* used in the config file and for parsing. */
#define LOCAL_CONFIG   "local-config"
#define LOCAL_ACCEPT   "local-accept"
#define LOCAL_HOSTNAME "local-hostname"
#define REMOTE_IP_FILE "remote-ip-file"
#define CERT "cert"
#define KEY  "key"

#define DEFAULT_LOCAL_ACCEPT   "5533"
#define DEFAULT_CERT           "certs/test.cert"
#define DEFAULT_KEY            "certs/server.key"
#define DEFAULT_REMOTE_PORT    "9275"
#define DEFAULT_LOCAL_HOSTNAME "localhost"
#define DEFAULT_ORG_FILENAME   "/ukama/org"

#define MAX_BUFFER 256
#define TRUE  1
#define FALSE 0

/* Struct to define the server and/or client host cfg. */
typedef struct {

    int  forwardPort;     /* Port on which to accept local clients */
    int  servicePort;
    char *localHostname;
    char *remoteConnect;  /* Client: hostname:port to connect with remotely */
    char *orgName;        /* org node belongs */
    char *certFile;       /* CA Cert file name. */
    char *keyFile;        /* Key file name.*/

    DeviceInfo *deviceInfo;   /* Device related info. */
} Config;

int process_config_file(Config *config, char *fileName);
void clear_config(Config *config);
void print_config(Config *config);
void split_strings(char *input, char **str1, char **str2, char *delimiter);

#endif /* MESH_CONFIG_H */
