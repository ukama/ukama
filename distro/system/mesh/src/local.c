/*
 * Code to handle 'local' client on the host. 
 * currently, this does nothing. Sniffer takes care of forwarding the 
 * packets to the right TX queue.
 *
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <netdb.h>
#include <pthread.h>

#include <sys/types.h>
#include <sys/socket.h>
#include <sys/uio.h>
#include <netinet/in.h>
#include <arpa/inet.h>

#include "config.h"
#include "log.h"

#define FAILURE 0
#define SUCCESS 1
#define BACKLOG 100
#define MAX_RCV_BUFFER 2048

/* 
 * Local:
 *  Get 'local-port' and bind 'local-host' from the config file.
 *  Listen on the local-port and create thread for each incoming client.
 *  Get forwarded packet on the listening port: 'local-port'
 *  Ignore them.
 *
 * Sniffer:
 *  Determine the dest address.
 *  Lookup the right cpool TX thread for the dest address.
 *  Acquire lock on the TX queue.
 *  Add to the TX queue.
 *  Update stats.
 *  Close the local connection unless "keep-alive" is valid.
 *  Repeat forever.
 *
 */

/*
 * init_local_bind --
 *
 */

int init_local_bind(HostConfig *config) {

  struct addrinfo hints;
  struct addrinfo *result, *rp;
  int sock, ret;

  memset(&hints, 0, sizeof(struct addrinfo));
  hints.ai_family   = AF_UNSPEC;  /* Allow IPv4 or IPv6 */
  hints.ai_socktype = SOCK_STREAM; 
  hints.ai_flags    = AI_PASSIVE; /* For wildcard IP address */
  hints.ai_protocol = 0;          /* Any protocol */
  
  ret = getaddrinfo(config->localHostname, config->localPort, &hints, &result);
  if (ret != 0) {
    log_error("getaddrinfo failed: %s", gai_strerror(ret));
    return FAILURE;
  }
  
  /* getaddrinfo() returns a list of address structures.
   * Try each address until we successfully bind(2).
   * If socket(2) (or bind(2)) fails, we (close the socket
   * and) try the next address. 
   */
  
  for (rp = result; rp != NULL; rp = rp->ai_next) {
    
    sock = socket(rp->ai_family, rp->ai_socktype, rp->ai_protocol);
    
    if (sock == -1) {
      continue;
    }
    
    if (bind(sock, rp->ai_addr, rp->ai_addrlen) == 0) {
      break;
    }
    
    close(sock);
  }

  if (rp == NULL) { /* No address found. */
    log_error("Could not bind to %s:%d", config->localHostname,
	      config->localPort);
    return FAILURE;
  }

  freeaddrinfo(result);

#if 0 // XXX
  /* We want to know dest ip and port info for the incoming packet.
   * Setting IP_PKTINFO for setsocktop() will help us with this.
   */
  if (setsockopt(sock, IPPROTO_IP, DSTADD_SOCKOPT, &sockopt,
		 sizeof(sockopt)) == -1) {
    log_error("Error setting sock option: %s", strerror(errno));
    return FAILURE;
  }
#endif
  
  return sock;
}

/*
 * client_handler --
 *
 */

void *client_handler(void *client) {

  int sock, len;
  char buffer[MAX_RCV_BUFFER];
  
  sock = *(int *)client;

  /* Read the packet and do absolutely nothing with it! */
  while ( (len = recv(sock, buffer, MAX_RCV_BUFFER, 0)) > 0) {
    /* ignore it. */
    log_debug("Recevied packet of len: %d. Ignore. [Packet: %s]", len, buffer);
    memset(buffer, 0, MAX_RCV_BUFFER);
  }
  
  /* disconnected, thread exit. */
  return NULL;
}


/*
 * accept_local_conn --
 *
 */

void *accept_local_conn(void *arg) {

  int ret, sock;
  pthread_t id;  
  struct sockaddr_in clientaddr;
  socklen_t addrLen;

  sock = *(int *)arg;
  
  if (listen(sock, BACKLOG) != 0) {
    log_error("Listen error: %s", strerror(errno));
    return NULL;
  }

  /* forever! */
  while (1) {

    addrLen = sizeof(clientaddr);

    ret = accept(sock, (struct sockaddr *)&clientaddr, &addrLen);

    if (ret < 0) {
      log_error("Accept error: %s", strerror(errno));
      continue;
    }

    if (pthread_create(&id, NULL, client_handler, (void*) &sock) < 0) {
      log_error("Error creating thread for client");
      continue;
    }
  }
  
  close(sock);
  return NULL;
}

/*
 * enable_local_clients -- 
 *
 */

int enable_local_clients(HostConfig *config) {

  int sock;
  pthread_t id;
  
  /* Some basic sanity checks. */
  if (!config)
    return FAILURE;

  if (config->localHostname == NULL || config->localPort == NULL) {
    return FAILURE;
  }
  
  sock = init_local_bind(config);

  if (sock < 0) {
    log_error("Error enabling local client. Binding failed");
    return FALSE;
  }

  if (pthread_create(&id, NULL, accept_local_conn, (void *) &sock) < 0) {
    close(sock);
    log_error("Error creating thread to accept local clients.");
    return FAILURE;
  }

  return SUCCESS;
}
