/*
 * SSL/TLS based connection.
 *
 *
 */

#include "ssl.h"

static void my_debug(void *ctx, int level, const char *file, int line,
                      const char *str) {
  ((void) level);
  
  mbedtls_fprintf( (FILE *) ctx, "%s:%04d: %s", file, line, str );
  fflush(  (FILE *) ctx  );
}

#if 0
/*
 * start_non_TLS_server -- server which will accept forwarded client packets.
 *
 *
 */

int start_non_TLS_server(char *port) {

  int ret=FALSE;
  struct addrinfo hints, *res, *p;
  int listenfd;
  
  // getaddrinfo for host
  memset (&hints, 0, sizeof(hints));
  hints.ai_family = AF_INET;
  hints.ai_socktype = SOCK_STREAM;
  hints.ai_flags = AI_PASSIVE;
  
  if (getaddrinfo(NULL, port, &hints, &res) != 0) {
    log_error("Error getting addr info");
    return ret;
  }
  
  // socket and bind
  for (p = res; p != NULL; p=p->ai_next) {
    
    listenfd = socket (p->ai_family, p->ai_socktype, 0);
    if (listenfd == -1) 
      continue;
    
    if (bind(listenfd, p->ai_addr, p->ai_addrlen) == 0)
      break;
  }
  
  if ( p == NULL) {
    log_error("Socket/bind error ...");
    return FALSE;
  }
    
  freeaddrinfo(res);
  
  /* Listen for nonTLS clients. */


  
  return ret;
}
#endif

/*
 * init_connection -- Initialize the Connection struct. 
 *
 */

void init_connection(Connection *conn) {

  mbedtls_net_init(&conn->fd);
  mbedtls_ssl_init(&conn->ssl);
  mbedtls_ssl_config_init(&conn->conf);
  mbedtls_x509_crt_init(&conn->cert);
  mbedtls_ctr_drbg_init(&conn->ctr_drbg);
  mbedtls_entropy_init(&conn->entropy);

  conn->cloud = TRUE; /* XXX - for now. */
  
}

/*
 * free_connection --
 *
 *
 */

void free_connection(Connection *conn) {

  mbedtls_net_free(&conn->fd);
  mbedtls_x509_crt_free(&conn->cert);
  mbedtls_ssl_free(&conn->ssl);
  mbedtls_ssl_config_free(&conn->conf);
  mbedtls_ctr_drbg_free(&conn->ctr_drbg);
  mbedtls_entropy_free(&conn->entropy);
}

/*
 * connect_to_secure_server -- Connect to SSL/TLS server and return SSL context.
 *
 */

int connect_to_secure_server(Connection *conn, const char *serverName,
			     const char *portNumber, const char *certFile) {
  
  int ret;

  init_connection(conn);

  ret = mbedtls_ctr_drbg_seed(&conn->ctr_drbg, mbedtls_entropy_func,
			      &conn->entropy, NULL, 0);
  
  if (ret != 0) {
    log_error("RNG seeding failed: %d", ret);
    goto done;
  }

#if defined(TEST_EMBED_CERT)
  log_debug("Loading embed cert and key");

  ret = mbedtls_x509_crt_parse(&conn->cert,
			       (const unsigned char *) mbedtls_test_cas_pem,
			       mbedtls_test_cas_pem_len );
  if(ret != 0) {
    log_error("Error loading cert!");
    goto done;
  }
#else
  ret = mbedtls_x509_crt_parse_file(&conn->cert, certFile);
  if (ret != 0){
    log_error("CRT parsing failed for file: %s with error: %d", certFile,
	      ret);
    goto done;
  }
#endif /* TEST_EMBED_CERT */

  /* Start connecting to the server. */
  log_debug("Connecting to SSL/TLS server at %s:%s", serverName, portNumber);
  
  ret = mbedtls_net_connect(&conn->fd, serverName, portNumber,
			    MBEDTLS_NET_PROTO_TCP);

  if (ret != 0) {
    log_error("Failed connecting to server: %s at port: %s", serverName,
	      portNumber);
    goto done;
  }
  
  ret = mbedtls_ssl_config_defaults(&conn->conf, MBEDTLS_SSL_IS_CLIENT,
				    MBEDTLS_SSL_TRANSPORT_STREAM,
				    MBEDTLS_SSL_PRESET_DEFAULT);
  
  if (ret != 0) {
    log_error("Failed to setup SSL/TLS structure: %s", ret);
    goto done;
  }

  mbedtls_ssl_conf_authmode(&conn->conf, MBEDTLS_SSL_VERIFY_OPTIONAL); /* XXX */
  mbedtls_ssl_conf_ca_chain(&conn->conf, &conn->cert, NULL);
  mbedtls_ssl_conf_rng(&conn->conf, mbedtls_ctr_drbg_random, &conn->ctr_drbg);
  mbedtls_ssl_conf_dbg(&conn->conf, my_debug, stdout);

  ret = mbedtls_ssl_setup(&conn->ssl, &conn->conf);
  if (ret != 0) {
    log_error("Failed to setup SSL");
    goto done;
  }

  ret = mbedtls_ssl_set_hostname(&conn->ssl, serverName);
  if (ret != 0) {
    log_error("Failed to setup hostname");
    goto done;
  }

  mbedtls_ssl_set_bio(&conn->ssl, &conn->fd, mbedtls_net_send,
		      mbedtls_net_recv, NULL );
  
  /* Perform SSL/TLS handshake with the server. */
  while ((ret = mbedtls_ssl_handshake(&conn->ssl)) != 0) {
    if (ret != MBEDTLS_ERR_SSL_WANT_READ &&
	ret != MBEDTLS_ERR_SSL_WANT_WRITE) {
      log_error("Failed to handshake with server: %d", ret);
      goto done;
    }
  }

  log_debug("Handshake with %s:%d succesful", serverName, portNumber);
  
  /* Verify. */
  ret = mbedtls_ssl_get_verify_result(&conn->ssl);
  if (ret != 0) {
    log_error("Cert verification failed!");
    goto done;
  }
     
  log_debug("Cert verified, all systems are go");

  /* ssl will be use to/from read/write to the server. */

  return TRUE;
  
 done:
    
  free_connection(conn);
  return FALSE;
}
