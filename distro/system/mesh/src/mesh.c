/*
 * Mesh.d
 *
 */

#include "mesh.h"
#include "mbedtls/certs.h"

#define VERSION "0.0.1"

#if !defined(MBEDTLS_BIGNUM_C) || !defined(MBEDTLS_ENTROPY_C) ||  \
    !defined(MBEDTLS_SSL_TLS_C) || !defined(MBEDTLS_SSL_CLI_C) || \
    !defined(MBEDTLS_NET_C) || !defined(MBEDTLS_RSA_C) ||         \
    !defined(MBEDTLS_CERTS_C) || !defined(MBEDTLS_PEM_PARSE_C) || \
    !defined(MBEDTLS_CTR_DRBG_C) || !defined(MBEDTLS_X509_CRT_PARSE_C)
int main( void )
{
    mbedtls_printf("MBEDTLS_BIGNUM_C and/or MBEDTLS_ENTROPY_C and/or "
           "MBEDTLS_SSL_TLS_C and/or MBEDTLS_SSL_CLI_C and/or "
           "MBEDTLS_NET_C and/or MBEDTLS_RSA_C and/or "
           "MBEDTLS_CTR_DRBG_C and/or MBEDTLS_X509_CRT_PARSE_C "
           "not defined.\n");
    mbedtls_exit( 0 );
}
#else


static void my_debug(void *ctx, int level, const char *file, int line,
                      const char *str) {
  ((void) level);
  
  mbedtls_fprintf( (FILE *) ctx, "%s:%04d: %s", file, line, str );
  fflush(  (FILE *) ctx  );
}

/*
 * usage -- Usage options for the Mesh.
 *
 *
 */

void usage() {

  printf("Usage: mesh.d [options] \n");
  printf("Options:\n");
  printf("--h, --help                         Help menu.\n");
  printf("--P, --port                         Port to listen SSL/TLS\n");
  printf("--C, --ca                           CRT file\n");
  printf("--k, --key                          Key file\n");
  printf("--l, --level <TRACE | DEBUG | INFO> Log level for the process.\n");
  printf("--V, --version                      Version.\n");
}

/* Set the verbosity level for logs. */
void set_log_level(char *slevel) {

  int ilevel = LOG_TRACE;

  if (!strcmp(slevel, "DEBUG")) {
    ilevel = LOG_DEBUG;
  } else if (!strcmp(slevel, "INFO")) {
    ilevel = LOG_INFO;
  } else if (!strcmp(slevel, "ERROR")) {
    ilevel = LOG_ERROR;
  }

  log_set_level(ilevel);
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

/* 
 * Cloud-config-file. 
 *
 * [TLS-server-config]
 *   CA path.
 *   CA file.
 *   CRT file.
 *   TLS server listening port.
 *   debug-level
 * 
 * [nonTLS-client-config]
 *   listen port
 *   debug-level
 *
 * [cloud-config]
 *   CA cert file path
 *
 */

int main (int argc, char **argv) {

  int ret = 0;
  char *tlsListenPort;
  char *keyFile=NULL, *caFile=NULL, *crlFile=NULL;
  char *caPath = NULL;
  char *debug = DEF_LOG_LEVEL;

  int server, client; 

  mbedtls_net_context tlsListenFd, tlsClientFd;
  mbedtls_entropy_context entropy;
  mbedtls_ctr_drbg_context ctr_drbg;
  mbedtls_ssl_context ssl;
  mbedtls_ssl_config conf;
  mbedtls_x509_crt srvcert;
  mbedtls_pk_context key;

  Connection cloud;
  
  /* Initalize some values. */
  mbedtls_net_init(&tlsListenFd);
  mbedtls_net_init(&tlsClientFd);
  mbedtls_ctr_drbg_init(&ctr_drbg);
  mbedtls_ssl_init(&ssl);
  mbedtls_ssl_config_init(&conf);
  mbedtls_x509_crt_init(&srvcert);
  mbedtls_pk_init(&key);
  mbedtls_entropy_init(&entropy);

  if (argc == 1) {
    fprintf(stderr, "Missing required parameters\n");
    usage();
    exit(1);
  }

  /* Prase command line args. */
  while (TRUE) {
    
    int opt = 0;
    int opdidx = 0;

    static struct option long_options[] = {
      { "port",      required_argument, 0, 'P'},
      { "level",     required_argument, 0, 'l'},
      { "key",       required_argument, 0, 'k'},
      { "ca",        required_argument, 0, 'A'},
      { "crl",       required_argument, 0, 'c'},
      { "path",      required_argument, 0, 'p'},
      { "server",    no_argument,       0, 'S'},
      { "client",    no_argument,       0, 'C'},
      { "help",      no_argument,       0, 'h'},
      { "version",   no_argument,       0, 'V'},
      { 0,           0,                 0,  0}
    };

    opt = getopt_long(argc, argv, "P:l:k:A:c:p:hVSC:", long_options, &opdidx);
    if (opt == -1) {
      break;
    }

    switch (opt) {
    case 'h':
      usage();
      exit(0);
      break;

    case 'S':
      server = TRUE;
      break;

    case 'C':
      client = TRUE;
      break;
      
    case 'P':
      tlsListenPort = optarg;
      break;

    case 'l':
      debug = optarg;
      set_log_level(debug);
      break;
      
    case 'k':
      keyFile = optarg;
      break;

    case 'A':
      caFile = optarg;
      break;

    case 'c':
      crlFile = optarg;
      break;

    case 'p':
      caPath = optarg;
      break;
      
    case 'V':
      fprintf(stdout, "Mesh.d - Version: %s\n", VERSION);
      exit(0);

    default:
      usage();
      exit(0);
    }
  } /* while */

  if (server) {
  
    log_debug("Starting mesh data plane ... [Server]");

#if defined(TEST_EMBED_CERT)
    log_debug("Loading embed cert and key.");
    ret = mbedtls_x509_crt_parse(&srvcert,
				 (const unsigned char *) mbedtls_test_srv_crt,
				 mbedtls_test_srv_crt_len);
    if(ret != 0) {
      log_error("Loading server cert and key failed");
      goto exit;
    }

    ret = mbedtls_x509_crt_parse(&srvcert,
				 (const unsigned char *) mbedtls_test_cas_pem,
				 mbedtls_test_cas_pem_len);
    if(ret != 0) {
      log_error("Loading server cert and key failed.");
      goto exit;
    }
    
    ret =  mbedtls_pk_parse_key(&key, (const unsigned char *) mbedtls_test_srv_key,
				mbedtls_test_srv_key_len, NULL, 0 );
    if(ret != 0) {
      log_error("Loading key file failed");
      goto exit;
    }
#else
    /* Load the cert and private key. */
    if (caFile) {
      ret = mbedtls_x509_crt_parse_file(&srvcert, caFile);
      if (ret != 0){
	log_error("CRT parsing failed: %d", ret);
	goto exit;
      }
    }

    if (keyFile) {
      ret = mbedtls_pk_parse_keyfile(&key, keyFile, NULL);
      if (ret != 0){
	log_error("Key file parsing failed: %d", ret );
	goto exit;
      }
  }
#endif  /* TEST_EMBED_CERT */
    
    /*
     * 2. Setup the listening TCP socket
     */
    ret = mbedtls_net_bind(&tlsListenFd, NULL, tlsListenPort,
			   MBEDTLS_NET_PROTO_TCP);
    
    if (ret != 0) {
      log_error("Failed to bind on port: %d. Code: %d", tlsListenPort, ret);
      goto exit;
    }
    
    /*
     * 3. Seed the RNG
     */
    ret = mbedtls_ctr_drbg_seed(&ctr_drbg, mbedtls_entropy_func, &entropy,
				NULL, /* XXX this could be device UUID. */
				0);
    
    if (ret != 0) {
      log_error("Failed to seed the random number generator.");
      goto exit;
    }
    
    ret = mbedtls_ssl_config_defaults(&conf, MBEDTLS_SSL_IS_SERVER,
				      MBEDTLS_SSL_TRANSPORT_STREAM,
				      MBEDTLS_SSL_PRESET_DEFAULT);
    
    if (ret != 0) {
      log_error("Failed to setup SSL data.");
      goto exit;
    }
    
    mbedtls_ssl_conf_rng(&conf, mbedtls_ctr_drbg_random, &ctr_drbg );
    mbedtls_ssl_conf_dbg(&conf, my_debug, stdout );
    
    mbedtls_ssl_conf_ca_chain(&conf, srvcert.next, NULL );
    
    ret = mbedtls_ssl_conf_own_cert(&conf, &srvcert, &key);
    
    if (ret !=0 ) {
      log_error("Failed to setup SSL data.");
      goto exit;
    }
    
    ret = mbedtls_ssl_setup(&ssl, &conf);
    
    if (ret != 0) {
      log_error("Failed to setup SSL data.");
      goto exit;
    }
    
  reset: /* Will come back here if handshake failed. */
    
    mbedtls_net_free(&tlsClientFd);
    mbedtls_ssl_session_reset(&ssl);
    
    /* Wait for client connection ... */
    
    log_debug("Waiting for client on port: %s ...", tlsListenPort);
    
    ret = mbedtls_net_accept(&tlsListenFd, &tlsClientFd, NULL, 0, NULL);
    
    if (ret != 0) {
      log_error("Accept failed: %d", ret);
      goto exit;
    }
    
    mbedtls_ssl_set_bio(&ssl, &tlsClientFd, mbedtls_net_send,
			mbedtls_net_recv, NULL);
    
    while ((ret = mbedtls_ssl_handshake(&ssl)) != 0) {
      if (ret != MBEDTLS_ERR_SSL_WANT_READ &&
	  ret != MBEDTLS_ERR_SSL_WANT_WRITE) {
	log_error("Handshake failed: %d", ret);
	goto reset;
      }
    }
  } /* if (server) */

  if (client) {
    
    /* Connect to the secure cloud server at given port. Mesh.d can be server
     * and client at same time.
     */
    connect_to_secure_server(&cloud, DEF_CLOUD_SERVER_NAME,
			     DEF_CLOUD_SERVER_PORT,
			     DEF_CLOUD_SERVER_CERT);
  }
  
  /* Connection established. */

  log_debug("All done. Exiting ...");
  return 0;

 exit:

  log_error("FAIL!. Exiting ...");
  
  mbedtls_net_free(&tlsListenFd);
  mbedtls_ctr_drbg_free(&ctr_drbg);
  mbedtls_entropy_free(&entropy);
  
  return 1;
}

#endif /* MBEDTLS_BIGNUM_C && MBEDTLS_ENTROPY_C && MBEDTLS_SSL_TLS_C &&
          MBEDTLS_SSL_CLI_C && MBEDTLS_NET_C && MBEDTLS_RSA_C &&
          MBEDTLS_CERTS_C && MBEDTLS_PEM_PARSE_C && MBEDTLS_CTR_DRBG_C &&
          MBEDTLS_X509_CRT_PARSE_C */
