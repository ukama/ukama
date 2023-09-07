/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "msghandler.h"
#include "msghandlerproc.h"

#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "headers/ubsp/devices.h"
#include "headers/utils/list.h"
#include "headers/utils/log.h"

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>
#ifdef HAVE_SYS_TIME_H
#include <sys/time.h>
#endif

static pthread_t msghandlerid;

#define MAX_LENGTH 4096

#define SADDR struct sockaddr

#define IFMSG_DEBUG

int g_server_sockfd = 0;

void msghandler_print(char *data, int size) {
    int iter = 0;
    while (iter < size) {
        if (iter % 16 == 0) {
            printf("\n 0x%02X", (data[iter] & 0xFF));
        } else {
            printf("\t 0x%02X", (data[iter] & 0xFF));
        }
        iter++;
    }
    printf("\n");
    fflush(stdout);
}

char *msghandler_recv(int sockfd) {
    char *req = NULL;
    int bytes = 0;
    req = malloc(sizeof(char) * MAX_LENGTH);
    if (req) {
        bzero(req, MAX_LENGTH);

        /* read the message from client and copy it in buffer*/
        bytes = recv(sockfd, req, MAX_LENGTH, 0);
        if (bytes < 0) {
            free(req);
            req = NULL;
            if (errno == EWOULDBLOCK) {
                log_error(
                    "Err(%d): MSGHANDLER:: Timeout for the receive functions.");
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
    log_trace("MSGHANDLER::Closed client socket created %d .", sockfd);
    sockfd = -1;
}

/* Function designed for service req. between client and server.*/
void msghandler_server_func(int sockfd) {
    char *resp = NULL;
    char *req = NULL;
    size_t respsize = 0;
    /* Timeout for recv */
    struct timeval tv;
    tv.tv_sec = RECV_MSG_TIMEOUT;
    tv.tv_usec = 0;
    setsockopt(sockfd, SOL_SOCKET, SO_RCVTIMEO, &tv, sizeof(tv));

    req = msghandler_recv(sockfd);
    if (req) {
        /* print buffer which contains the client contents */
        log_trace("MSGHANDLER:: From client to UkamaEDR.", req);
        msghandler_print(req, 64);

        resp = msghandler_proc(req, &respsize);
        if (resp) {
            /* print buffer which contains the client contents */
            log_trace("MSGHANDLER:: From UkamaEDR to client.");
            msghandler_print(resp, respsize);

            msghandler_send(sockfd, resp, respsize);
        } else {
            log_trace("MSGHANDLER:: No response message...!");
        }
        UKAMA_FREE(resp);
    } else {
        log_error("Err: MSGHANDLER:: Error while receiving messages.");
    }
    UKAMA_FREE(req);
}

// Function designed for Asynchronous request between client and server.
MsgFrame *msghandler_client_func(int sockfd, MsgFrame *smsg, size_t *size,
                                 int *sflag) {
    int ret = 0;
    MsgFrame *rmsg = NULL;
    /* Serialize */
    char *sdata = msgframe_serialize(smsg, size);
    if (sdata) {
        /* Timeout for recv */
        struct timeval tv;
        tv.tv_sec = RECV_MSG_TIMEOUT;
        tv.tv_usec = 0;
        setsockopt(sockfd, SOL_SOCKET, SO_RCVTIMEO, &tv, sizeof(tv));

#ifdef IFMSG_DEBUG
        /* print buffer which contains the client contents */
        log_debug("MSGHANDLER:: From UkamaEDR to client.");
        msghandler_print(sdata, *size);
#endif
        ret = msghandler_send(sockfd, sdata, *size);
        if (ret > 0) {
            char *rdata = msghandler_recv(sockfd);
            if (rdata) {
#ifdef IFMSG_DEBUG
                /* print buffer which contains the client contents */
                log_debug("MSGHANDLER:: From client to UkamaEDR.");
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

void msghandler_exit() {
    /* Closing socket*/
    log_debug("MSGHANDLER:: Closing UkamaEDR server socket.");
    close(g_server_sockfd);
    /* Exiting thread*/
    msghandler_stop();
}

int msg_handler_create_sock() {
    int ret = 0;
    int sockfd, connfd;
    /* socket create and verification*/
    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd == -1) {
        log_error("Err(%d): MSGHANDLER:: Socket creation failed...\n", ret);
    } else {
        log_trace("MSGHANDLER:: Socket successfully created..\n");
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
    servaddr.sin_port = htons(UKAMALWM2MPORT);

    /* connect the client socket to server socket */
    if (connect(sockfd, (SADDR *)&servaddr, sizeof(servaddr)) != 0) {
        ret = -1;
        log_error(
            "Err(%d): MSGHANDLER:: connection with the server failed...\n",
            ret);
    } else {
        log_trace("MSGHANDLER:: Connected to the server..\n");
    }
    return ret;
}

/* Client for the Asynchronous messages like event and  Alerts.*/
MsgFrame *msghandler_client_send(MsgFrame *msg, size_t *size, int *sflag) {
    MsgFrame *rmsg = NULL;
    int ret = 0;
    /* Create socket */
    int clientsock = msg_handler_create_sock();
    if (clientsock > 0) {
        /* connect to server */
        if (msghandler_sock_connect(clientsock)) {
            ret = ERR_SOCK_CONNECT;
            goto cleanup;
        }
    } else {
        ret = ERR_SOCK_CREATION;
        goto cleanup;
    }

    /* Send the message */
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
int msghandler_server() {
    int ret = 0;
    int sockfd, connfd;
    unsigned int len;
    struct sockaddr_in servaddr, cli;

    /* socket create and verification */
    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd == -1) {
        ret = -1;
        log_error("Err(%d): MSGHANDLER:: Socket creation failed...\n", ret);
        goto cleanup;
    } else {
        g_server_sockfd = sockfd;
        log_trace("MSGHANDLER:: Socket successfully created..\n");
    }
    bzero(&servaddr, sizeof(servaddr));

    /* assign IP, PORT */
    servaddr.sin_family = AF_INET;
    servaddr.sin_addr.s_addr = htonl(INADDR_ANY);
    servaddr.sin_port = htons(UKAMAREGPORT);

    int flag = 1;
    if (-1 == setsockopt(sockfd, SOL_SOCKET, (SO_REUSEADDR | SO_REUSEPORT),
                         &flag, sizeof(flag))) {
        log_error("Err(%d): setsockopt fail for UkamaEDR.");
        goto cleanup;
    }

    /* Binding newly created socket to given IP and verification */
    if ((bind(sockfd, (SADDR *)&servaddr, sizeof(servaddr))) != 0) {
        ret = -1;
        perror("Socket Failure.");
        log_error("Err(%d): MSGHANDLER:: Socket bind failed.\n", ret);
        goto cleanup;
    }

    log_trace("MSGHANDLER:: Socket successfully binded..\n");

    /* Now server is ready to listen and verification */
    if ((listen(sockfd, 5)) != 0) {
        ret = -1;
        log_error("Err(%d): MSGHANDLER:: Listen failed...\n", ret);
        goto cleanup;
    } else
        log_debug("MSGHANDLER:: UkamaEDR Server listening..\n");
    len = sizeof(cli);

    while (TRUE) {
        /* Accept the data packet from client and verification */
        connfd = accept(sockfd, (SADDR *)&cli, &len);
        if (connfd < 0) {
            ret = -1;
            log_error("Err(%d): MSGHANDLER:: Server accept failed...\n", ret);
            goto cleanup;
        } else {
            log_trace("MSGHANDLER:: Server accept the client...\n");
        }

        /* Function for service request*/
        msghandler_server_func(connfd);
    }

cleanup:
    if (ret) {
        /* After service close the socket */
        close(sockfd);
        g_server_sockfd = 0;
    }
    return ret;
}

void *msghandler_service(void *data) {
    msghandler_server();
    return NULL;
}

void msghandler_start() {
    pthread_t serviceid = 0;
    if (pthread_create(&serviceid, NULL, &msghandler_service, NULL)) {
        /*Thread creation failed*/
        msghandlerid = 0;
    } else {
        msghandlerid = serviceid;
    }
}

void msghandler_stop() {
    log_debug("MSGHANDLER:: Exiting UkamaEDR server thread %ld.", msghandlerid);
    pthread_cancel(msghandlerid);
    pthread_join(msghandlerid, NULL);
    msghandlerid = 0;
}
