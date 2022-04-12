/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "inc/msghandler.h"

#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "inc/msghandlerproc.h"

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/time.h>
#include <arpa/inet.h>
#include <unistd.h>

#define MAX_LENGTH 		4096

#define SADDR struct sockaddr

void msghandler_print(char* data, int size) {
	int iter  = 0;
	while(iter<size) {
		if (iter%16 == 0) {
			fprintf(stdout, "\r\n 0x%02X", (data[iter]&0xFF));
		}
		else {
			fprintf(stdout, "\t 0x%02X", (data[iter]&0xFF));
		}
		iter++;
	}
	printf("\r\n");
	fflush(stdout);
}

char *msghandler_recv(int sockfd) {
	int bytes = 0;
	char *req = malloc(sizeof(char) * MAX_LENGTH);
	if (req) {
		bzero(req, MAX_LENGTH);
		/* read the message from client and copy it in buffer*/
		bytes = recv(sockfd, req, MAX_LENGTH, 0);
		if (bytes < 0 ) {
			free(req) ;
			req =  NULL;
			if (errno ==  EWOULDBLOCK) {
				fprintf(stderr, "Err: MSGHANDLER:: Timeout for the receive functions.\r\n");
			}
		}
	}
	return req;
}

int msghandler_send(int sockfd, void *data, size_t size) {
    /* and send that buffer to client */
    return send(sockfd, data, size, 0);
}

void msghandler_close_connection(int sockfd) {
    close(sockfd);
    fprintf(stdout, "MSGHANDLER::Closed client socket created %d\r\n .", sockfd );
    sockfd = -1;
}

/* Function designed for service req. between client and server.*/
void msghandler_server_func(void* ctx, int sockfd) {
    char *resp = NULL;
    size_t respsize = 0;
    /* Timeout for recv */
    struct timeval tv;
    tv.tv_sec = RECV_MSG_TIMEOUT;
    tv.tv_usec = 0;
    setsockopt(sockfd, SOL_SOCKET, SO_RCVTIMEO, &tv, sizeof(tv));

    char *req = msghandler_recv(sockfd);
    if (req) {
#ifdef IFMSG_DEBUG
        /* print buffer which contains the client contents */
        fprintf(stdout, "MSGHANDLER:: From client to server.\r\n");
        msghandler_print(req, 64);
#endif
        resp = msghandler_proc(ctx, req, &respsize);
        if (resp) {
#ifdef IFMSG_DEBUG
            /* print buffer which contains the client contents */
            fprintf(stdout, "MSGHANDLER:: From server to client.\r\n");
            msghandler_print(req, respsize);
#endif
            msghandler_send(sockfd, resp, respsize);
        } else {
            fprintf(stdout, "MSGHANDLER:: No response message...\r\n!");
        }
        UKAMA_FREE(resp);
    }
    UKAMA_FREE(req);
}

// Function designed for Asynchronous request between client and server.
MsgFrame* msghandler_client_func(int sockfd, MsgFrame *smsg, size_t* size, int *sflag) {
	int ret = 0;
	MsgFrame* rmsg = NULL;
	/* Serialize */
	char *sdata = msgframe_serialize(smsg, size);
	if (sdata) {
#ifdef IFMSG_DEBUG
		/* print buffer which contains the client contents */
		fprintf(stdout, "MSGHANDLER:: From server to client.\r\n");
		msghandler_print(sdata, *size);
#endif
		ret = msghandler_send(sockfd, sdata, *size);
		if (ret > 0) {
			char *rdata = msghandler_recv(sockfd);
			if (rdata) {
#ifdef IFMSG_DEBUG
				/* print buffer which contains the client contents */
				fprintf(stdout, "MSGHANDLER:: From client to server.\r\n");
				msghandler_print(rdata, 1026);
#endif
				/* Deserialize message */
				rmsg = msgframe_deserialize(rdata);

				/* Validate the message */
				ret = msgframe_validate(rmsg, smsg);

			} else {
				ret = ERR_SOCK_RECV;
			}
			UKAMA_FREE(rdata);
		} else {
			ret = ERR_SOCK_SEND;
		}
		UKAMA_FREE(sdata);
	} else {
		ret = ERR_IFMSG_SERIALIZATION;
	}
	*sflag = ret;
	return rmsg;
}

void msghandler_init() {

}

int msg_handler_create_sock() {
    int ret = 0;
    int sockfd, connfd;
    /* socket create and verification*/
    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd == -1) {
        fprintf(stderr, "Err(%d): MSGHANDLER:: Socket creation failed...\r\n", ret);
    } else {
        fprintf(stdout, "MSGHANDLER:: Socket %d successfully created \r\n", sockfd);
    }
    return sockfd;
}

int msghandler_sock_connect(int sockfd) {
    int ret = 0;
    struct sockaddr_in servaddr, cli;
    bzero(&servaddr, sizeof(servaddr));
    /* assign IP, PORT*/
    servaddr.sin_family = AF_INET;
    servaddr.sin_addr.s_addr = inet_addr("127.0.0.1");
    servaddr.sin_port = htons(UKAMAREGPORT);

    /* connect the client socket to server socket */
    if (connect(sockfd, (SADDR *)&servaddr, sizeof(servaddr)) != 0) {
        ret = -1;
        fprintf(stderr,
            "Err(%d): MSGHANDLER:: Connection with the server failed...\r\n",
            ret);
    } else {
        fprintf(stdout, "MSGHANDLER:: Connected to the server..\r\n");
    }
    return ret;
}

/* lwm2m send message to the ukamaEDR*/
MsgFrame* msghandler_client_send(MsgFrame *msg, size_t* size, int* sflag) {
	MsgFrame* rmsg = NULL;
	int ret = 0;
	/* Create socket */
	int clientsock = msg_handler_create_sock();
	if (clientsock > 0 ) {
		/* connect to server */
		if (msghandler_sock_connect(clientsock)) {
			ret = ERR_SOCK_CONNECT;
			goto cleanup;
		}
	} else {
		ret = ERR_SOCK_CREATION;
		goto cleanup;
	}

	*sflag = ret;
	if (!ret) {
		rmsg = msghandler_client_func(clientsock, msg, size, sflag);
	}

	cleanup:

	/* After service close the socket */
	msghandler_close_connection(clientsock);

	return rmsg;
}

/* Server for the request/response messages from the client.*/
int msghandler_server(void* ctx) {
    int ret = 0;
    int sockfd, connfd;
    unsigned int len;
    struct sockaddr_in servaddr, cli;

    /* socket create and verification */
    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd == -1) {
        ret = -1;
        fprintf(stderr, "Err(%d): MSGHANDLER:: Socket creation failed...\r\n", ret);
        goto cleanup;
    } else
        fprintf(stdout, "MSGHANDLER:: Socket successfully created..\r\n");
    bzero(&servaddr, sizeof(servaddr));

    /* assign IP, PORT */
    servaddr.sin_family = AF_INET;
    servaddr.sin_addr.s_addr = htonl(INADDR_ANY);
    servaddr.sin_port = htons(UKAMALWM2MPORT);

    int flag = 1;
    if (-1 == setsockopt(sockfd, SOL_SOCKET, (SO_REUSEADDR|SO_REUSEPORT), &flag, sizeof(flag))) {
    	fprintf(stderr, "Err: setsockopt fail for UkamaEDR.\r\n");
    	goto cleanup;
    }

    /* Binding newly created socket to given IP and verification */
    if ((bind(sockfd, (SADDR *)&servaddr, sizeof(servaddr))) != 0) {
        ret = -1;
        perror("Socket Failure.");
        fprintf(stderr, "Err(%d): MSGHANDLER:: Socket bind failed...\r\n", ret);
        goto cleanup;
    } else
        fprintf(stdout, "MSGHANDLER:: Socket successfully binded..\r\n");

    /* Now server is ready to listen and verification */
    if ((listen(sockfd, 5)) != 0) {
        ret = -1;
        fprintf(stderr, "Err(%d): MSGHANDLER:: Listen failed...\r\n", ret);
        goto cleanup;
    } else
        fprintf(stdout, "Server listening..\r\n");
    len = sizeof(cli);

    while (TRUE) {
        /* Accept the data packet from client and verification */
        connfd = accept(sockfd, (SADDR *)&cli, &len);
        if (connfd < 0) {
            ret = -1;
            fprintf(stderr, "Err(%d): MSGHANDLER:: Server accept failed...\r\n", ret);
            goto cleanup;
        } else {
            fprintf(stdout, "MSGHANDLER:: Server accept the client...\r\n");
        }

        /* Function for service request*/
        msghandler_server_func(ctx, connfd);

    }

cleanup:
    if (ret) {
        /* After service close the socket */
        close(sockfd);
    }
    return ret;
}

void *msghandler_service(void *args) {
    msghandler_server(args);
    return NULL;
}

pthread_t msghandler_start(void* data) {
	pthread_t serviceid = 0;
	if (data) {

		if (pthread_create(&serviceid, NULL, &msghandler_service, data)) {
			/*Thread creation failed*/
			serviceid = 0;
		}
	} else {
		fprintf(stdout, "MSGHANDLER:: Invalid lwm2mH context.\r\n");
	}
	return serviceid;
}

void msghandler_stop(pthread_t thread) {
    pthread_cancel(thread);
    pthread_join(thread, NULL);
}
