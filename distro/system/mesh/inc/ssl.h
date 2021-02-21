/*
 *
 *
 */

#ifndef SSL_H
#define SSL_H


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

/* mbed TLS includes */
#include "mbedtls/entropy.h"
#include "mbedtls/ctr_drbg.h"
#include "mbedtls/net_sockets.h"
#include "mbedtls/ssl.h"
#include "mbedtls/x509.h"
#include "mbedtls/debug.h"

/* certs. */
#include "mbedtls/certs.h"

/* Others. */
#include "log.h"

#define TRUE 1
#define FALSE 0

typedef struct {

  int cloud; /* TRUE if current context is for cloud connection. */

  mbedtls_net_context fd;
  mbedtls_entropy_context entropy;
  mbedtls_ctr_drbg_context ctr_drbg;
  mbedtls_ssl_context ssl;
  mbedtls_ssl_config conf;
  mbedtls_x509_crt cert;
}Connection;

#endif /* SSL_H */
 
