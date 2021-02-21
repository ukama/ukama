/*
 *
 */

#ifndef MESH_H
#define MESH_H

#include <getopt.h>

#include "log.h"
#include "config.h"
#include "ssl.h"

#define DEF_FILENAME "cert.crt"
#define DEF_CA_FILE  ""
#define DEF_CRL_FILE ""
#define DEF_CA_PATH  ""
#define DEF_SERVER_NAME "localhost"
#define DEF_TLS_SERVER_PORT "4444"
#define DEF_LOG_LEVEL "TRACE"
#define DEF_CLOUD_SERVER_NAME "localhost"
#define DEF_CLOUD_SERVER_PORT "4444"
#define DEF_CLOUD_SERVER_CERT "certs/test.crt"

#define TRUE 1
#define FALSE 0

#define PROXY_NONE    0x01
#define PROXY_FORWARD 0x02
#define PROXY_REVERSE 0x04

extern int connect_to_secure_server(Connection *conn, const char *serverName,
				    const char *portNumber,
				    const char *certFile);
extern int process_config_file(char *fileName, Configs *config);

#endif /* MESH_H */
