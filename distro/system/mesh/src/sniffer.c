/*
 * Sniff packets on the RAW socket binded to eth. 
 *
 */

#include "sniffer.h"
#include "log.h"

/* 
 * is_valid_iface -- Check if the given interface is up and running.
 *
 */

int is_valid_iface(char *ifname) {
  
  int sock, ret;
  struct ifreq if_req;
  
  sock = socket(AF_INET, SOCK_DGRAM, IPPROTO_IP);

  if (sock < 0) {
    log_error("Socket failed: %s", strerror(errno));
    return FALSE;
  }

  strncpy(if_req.ifr_name, ifname, sizeof(if_req.ifr_name));
  
  ret = ioctl(sock, SIOCGIFFLAGS, &if_req);
  close(sock);

  if ( ret < 0 ){
    log_error("Ioctl failed: %s\n", strerror(errno));
    return FALSE;
  }
  
  return (if_req.ifr_flags & IFF_UP) && (if_req.ifr_flags & IFF_RUNNING);
}

/* 
 * print_ip_header --
 *
 *
 */

void print_ip_header(unsigned char *buffer, int size) {
  
  struct iphdr *iph;
  struct sockaddr_in source, dest;
  
  iph = (struct iphdr *)buffer;

  memset(&source, 0, sizeof(source));
  source.sin_addr.s_addr = iph->saddr;
  
  memset(&dest, 0, sizeof(dest));
  dest.sin_addr.s_addr = iph->daddr;
  
  log_debug("IP Header");
  log_debug("   |-IP Version        : %d", (unsigned int)iph->version);
  log_debug("   |-IP Header Length  : %d DWORDS or %d bytes",
	    (unsigned int)iph->ihl,((unsigned int)(iph->ihl))*4);
  log_debug("   |-Type Of Service   : %d",(unsigned int)iph->tos);
  log_debug("   |-IP Total Length   : %d  Bytes(Size of Packet)\n",
	    ntohs(iph->tot_len));
  log_debug("   |-Identification    : %d",ntohs(iph->id));
  log_debug("   |-TTL               : %d",(unsigned int)iph->ttl);
  log_debug("   |-Protocol          : %d",(unsigned int)iph->protocol);
  log_debug("   |-Checksum          : %d",ntohs(iph->check));
  log_debug("   |-Source IP         : %s",inet_ntoa(source.sin_addr));
  log_debug("   |-Destination IP    : %s",inet_ntoa(dest.sin_addr));
}

/*
 * print_tcp_packet --
 *
 *
 */

void print_tcp_packet(unsigned char *buffer, int size) {
  
  unsigned short len;

  struct iphdr *iph;
  struct tcphdr *tcph;
  
  iph = (struct iphdr *)buffer;
  len = iph->ihl*4;
  
  tcph = (struct tcphdr*)(buffer + len);

  log_debug("\n ============= TCP Packet ============== \n");
  print_ip_header(buffer, size);

  log_debug("TCP Header");
  log_debug("   |-Source Port        : %u\n", ntohs(tcph->source));
  log_debug("   |-Destination Port   : %u\n", ntohs(tcph->dest));
  log_debug("   |-Sequence Number    : %u\n", ntohl(tcph->seq));
  log_debug("   |-Acknowledge Number : %u\n", ntohl(tcph->ack_seq));
  log_debug("   |-Header Length      : %d DWORDS or %d bytes",
	    (unsigned int)tcph->doff, (unsigned int)tcph->doff*4);
  log_debug("   |-Urgent Flag          : %d\n", (unsigned int)tcph->urg);
  log_debug("   |-Acknowledgement Flag : %d\n", (unsigned int)tcph->ack);
  log_debug("   |-Push Flag            : %d\n", (unsigned int)tcph->psh);
  log_debug("   |-Reset Flag           : %d\n", (unsigned int)tcph->rst);
  log_debug("   |-Synchronise Flag     : %d\n", (unsigned int)tcph->syn);
  log_debug("   |-Finish Flag          : %d\n", (unsigned int)tcph->fin);
  log_debug("   |-Window         : %d\n", ntohs(tcph->window));
  log_debug("   |-Checksum       : %d\n", ntohs(tcph->check));
  log_debug("   |-Urgent Pointer : %d\n", tcph->urg_ptr);
  
  log_debug("-----------------------------------------------------------");
}

/*
 * print_udp_packet --
 *
 *
 */
void print_udp_packet(unsigned char *buffer, int size) {
  
  unsigned short len;
  
  struct iphdr *iph;
  struct udphdr *udph;

  iph = (struct iphdr *)buffer;
  len = iph->ihl*4;

  udph = (struct udphdr*)(buffer + len);

  log_debug("\n ============= UDP Packet ============== \n");
  print_ip_header(buffer, size);

  log_debug("UDP Header");
  log_debug("   |-Source Port      : %d\n" , ntohs(udph->source));
  log_debug("   |-Destination Port : %d\n" , ntohs(udph->dest));
  log_debug("   |-UDP Length       : %d\n" , ntohs(udph->len));
  log_debug("   |-UDP Checksum     : %d\n" , ntohs(udph->check));
}

/*
 * is_valid_device_ip -- check if the passed IP is valid IP for the device.
 *
 */
int is_valid_device_ip(char *dIP) {

  int index=-1, size, i;
  ConnInfo *conn=NULL;
  
  /* IP is valid iff:
   * 1. there exists Node with 'dIP', and 
   * 2. there connection is live mesh.d and the Node. 
   */
  //  size = get_conn_info(conn); XXX

  for (i=0; i<size; i++) {
    if (strcmp(conn[i].node->ip, dIP) == 0) {
      if (conn[i].isLive) {
	index = i;
	log_debug("Matching Node found");
      } else {
	log_error("Matching Node found but is not Live: %s", dIP);
      }
      break;
    }
  }

  return index;
}

/* 
 * is_valid_packet -- Packet is defined valid if it is going to our 
 *                    device without any violation.
 *
 */

int is_valid_packet(unsigned char *buffer, int size, int pktType,
		    char *destIP) {

  struct iphdr *iph;
  struct sockaddr_in source, dest;
  char *sIP, *dIP;
  int index, ret;

  iph = (struct iphdr *)(buffer  + sizeof(struct ethhdr) );
	
  memset(&source, 0, sizeof(source));
  source.sin_addr.s_addr = iph->saddr;
	
  memset(&dest, 0, sizeof(dest));
  dest.sin_addr.s_addr = iph->daddr;
	
  sIP = strdup(inet_ntoa(source.sin_addr));
  dIP = strdup(inet_ntoa(dest.sin_addr));

  /* Packet is valid iff:
   * 1. dest IP matches with one of our devices IP in record.
   * 2. no policy violation from the source IP/port and dest port.
   */
  index = is_valid_device_ip(dIP);
  if (index < 0) {
    log_error("Trying to send packet to invalid device. sIP: %s dIP: %s .Drop!",
	      sIP, dIP);
    return FALSE;
  } /* ret is the index of the ConnInfo array */

  ret = run_policy_check(sIP, dIP, index, pktType);
  if (!ret) {
    log_error("Policy violated. sIP: %s dIP: %s .Ignoring", sIP, dIP);
    return FALSE;
  }

  /* copy the destination IP so the Cpool can move the packet to right queue.*/
  strncpy(dIP, destIP, strlen(dIP));
  
  return TRUE;
}

/*
 * filter_packet --
 *
 */
int filter_packet(unsigned char *buffer, int size, char *destIP) {

  struct iphdr *iph;
  int flag=0, ret=0;
  
  /* Get the IP header. */
  iph = (struct iphdr*)(buffer + sizeof(struct ethhdr));

  /* We only care for UDP and TCP packets, ignore rest. */
  switch (iph->protocol) {
  case 6:  /* TCP protocol. */
    flag = PKT_TCP;
    print_tcp_packet(buffer, size);
    break;
    
  case 17: /* UDP Protocol */
    flag = PKT_UDP;
    print_udp_packet(buffer, size);
    break;
    
  default: /* Everything else. */
    break;
  }

  if (flag == PKT_TCP || flag == PKT_UDP) {
    ret = is_valid_packet(buffer, size, flag, destIP);
  } else {
    return FALSE;
  }

  /* XXX */
  
  return ret;
}


/*
 * init_socket -- create RAW socket, bind to iface.
 *
 */ 

int init_socket(char *iface) {

  int rawSocket;
  int ret;

  rawSocket = socket(AF_PACKET, SOCK_RAW, htons(ETH_P_ALL)) ;

  if (rawSocket < 0) {
    log_error("Error creating RAW socket: %s", strerror(errno));
    return rawSocket;
  }

  /* Bind to the provided ethernet interface. This way, we will only listen
   * to the packets coming in from various microservices, as set by iptable.
   */

  ret = setsockopt(rawSocket, SOL_SOCKET, SO_BINDTODEVICE, iface,
		   strlen(iface)+1);

  if (ret < 0) {
    log_error("Error binding to interface %s: %s", iface, strerror(errno));
    close(rawSocket);
    return ret;
  }
  
  return rawSocket;
}

/*
 * inspect_packets -- read and filter packets. On the matching one(s), do CB.
 *
 */

int inspect_packets(int sock, unsigned char *buffer, int buffSize,
		    void (*cb)(int)) {

  int ret, size, len;
  struct sockaddr saddr;
  char destIP[32];

  size = sizeof(saddr);
  
  while (TRUE) {

    memset(destIP, 0, 32);
    
    /* Read a packet. */
    len = recvfrom(sock, buffer, buffSize, 0, &saddr, (socklen_t *)size);

    if (len < 0){
      log_error("recvfrom error: %s", strerror(errno));
      return FALSE;
    }

    ret = filter_packet(buffer, len, &destIP[0]);

    if (ret) { /* is valid packet. */
      /* XXX */
    }
  }

  /* Never would come here. */
  return TRUE;
}
