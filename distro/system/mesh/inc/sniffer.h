/*
 * sniffer.h
 *
 */

#ifndef SNIFFER_H
#define SNIFFER_H


#include <errno.h>
#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <linux/if_ether.h>

#include <sys/types.h>
#include <sys/socket.h>
#include <sys/ioctl.h>
#include <net/if.h>
#include <arpa/inet.h>

#include <netinet/in.h>
#include <netinet/ip.h>
#include <netinet/tcp.h>
#include <netinet/udp.h>


#define TRUE 1
#define FALSE 0

#define PKT_UDP 1
#define PKT_TCP 2

/*
 * NodeInfo utilized by mesh.d for housekeeping, via MQQT.
 *
 */
typedef struct {

  char *uuid; 
  char *ip;
  int  state;
}NodeInfo;

/* Some basic stats. */
typedef struct {
  
  long txBytes;
  long rxBytes;
  long dropPkts;
  long updPkts;
  long tcpPkts;
}ConnStats;

/*
 * For each Node, from the dReg, mesh.d maintains connection status.
 */
typedef struct {

  NodeInfo  *node;
  ConnStats *stats;
  int       isLive;
}ConnInfo;

extern int get_conn_info(ConnInfo *conn);
extern int run_policy_check(char *sIP, char *dIP, int index, int pkyType);

#endif /* SNIFFER_H */
