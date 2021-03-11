/*
 * ifhandler.h
 *
 *  Created on: Mar 8, 2021
 *      Author: vishal
 */

#ifndef SERVIFACE_IFHANDLER_H_
#define SERVIFACE_IFHANDLER_H_

#include <stdint.h>

#define ATTR_MAX_LEN 	(32)
#define MAX_LENGTH 		4096
#define LWM2M_GW_ADDRESS "127.0.0.0"
#define LWM2M_GW_PORT 3100
#define STATUS_OK 200
#define SADDR struct sockaddr


int connection_handler_create_sock();
int connection_handler_sock_connect(int sockfd, char* addr, uint32_t port);
void connection_handler_close_connection(int sockfd);
void connection_handler_stop(pthread_t thread);

pthread_t connection_handler_start(void *data);

#endif /* SERVIFACE_IFHANDLER_H_ */
