/*
 *
 */

#ifndef MESH_H
#define MESH_H


#include "mbedtls/config.h"
#include "mbedtls/platform.h"

#include <arpa/inet.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <signal.h>
#include <fcntl.h>
#include <getopt.h>

/* mbed TLS includes */
#include "mbedtls/entropy.h"
#include "mbedtls/ctr_drbg.h"
#include "mbedtls/net_sockets.h"
#include "mbedtls/ssl.h"
#include "mbedtls/x509.h"
#include "mbedtls/debug.h"

/* Others. */
#include "log.h"
#include "config.h"

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

typedef struct {

  int cloud; /* TRUE if current context is for cloud connection. */

  mbedtls_net_context fd;
  mbedtls_entropy_context entropy;
  mbedtls_ctr_drbg_context ctr_drbg;
  mbedtls_ssl_context ssl;
  mbedtls_ssl_config conf;
  mbedtls_x509_crt cert;
}Connection;

extern int process_config_file(char *fileName, Configs *config);

#endif /* MESH_H */
