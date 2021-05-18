#include <errno.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <time.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <pthread.h>

#include "list.h"
#include "ifhandler.h"
#include "request.h"
#include "notification.h"
#include "headers/errorcode.h"

#define DEBUG_GATEWAYIF
#define IF_LWM2M_SERVER_PORT 3000

#define CLIENT_REQ_MSG_LEN 512
#define MAX_LENGTH 		4096

#define SADDR struct sockaddr

//Request handler function
void *request_handler(void *);

//Connection handler
void *connection_handler();

/* Close socket */
void connection_handler_close_connection(int sockfd) {
    close(sockfd);
    fprintf(stdout, "IFHANDLER::Closed client socket created %d\r\n .", sockfd );
    sockfd = -1;
}

/* Create socket */
int connection_handler_create_sock() {
    int ret = 0;
    int sockfd, connfd;
    /* socket create and verification*/
    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd == -1) {
        fprintf(stderr, "Err(%d): IFHandler:: Socket creation failed...\r\n", ret);
    } else {
        fprintf(stdout, "IFHandler:: Socket %d successfully created \r\n", sockfd);
    }
    return sockfd;
}

/* Connect socket */
int connection_handler_sock_connect(int sockfd, char* addr, uint32_t port) {
    int ret = 0;
    struct sockaddr_in servaddr, cli;
    bzero(&servaddr, sizeof(servaddr));
    /* assign IP, PORT*/
    servaddr.sin_family = AF_INET;
    servaddr.sin_addr.s_addr = inet_addr(addr);
    servaddr.sin_port = htons(port);

    /* connect the client socket to server socket */
    if (connect(sockfd, (SADDR *)&servaddr, sizeof(servaddr)) < 0) {
        ret = -1;
        fprintf(stderr,
            "Err(%d): IFHandler:: Connection with the server failed...\r\n",
            ret);
    } else {
        fprintf(stdout, "IFHandler:: Connected to the server..\r\n");
    }
    return ret;
}

/* Handles the request. */
void *request_handler(void *socket_desc)
{
  /* Get the socket descriptor */
  int sock = *(int *)socket_desc;
  int n;

  char req_arr[CLIENT_REQ_MSG_LEN];
  fprintf(stdout, "IFHANDLER:: Waiting for incoming message Request handler thread id %lld\r\n", pthread_self());
  fflush(stdout);
  while ((n = recv(sock, req_arr, CLIENT_REQ_MSG_LEN, 0)) > 0)
  {
#ifdef DEBUG_GATEWAYIF
	  ifhandler_print(req_arr, CLIENT_REQ_MSG_LEN);
#endif
	  /* Service request */
	  serve_request(req_arr, sock);
  }

  fprintf(stdout, "IFHANDLER:: Exiting thread id %lld\r\n", pthread_self());
  fflush(stdout);
  pthread_exit(NULL);
}



/* Handles the incoming connection form the client */
void *connection_handler()
{
	int ret = 0;
	int client_sock, c, *new_sock;
	struct sockaddr_in server, client;

	/* Create socket */
	int socket_desc = connection_handler_create_sock();
	if (socket_desc < 0 ) {
		ret = -1;
		goto cleanup;
	}

	int flag = 1;
	if (-1 == setsockopt(socket_desc, SOL_SOCKET, (SO_REUSEADDR|SO_REUSEPORT), &flag, sizeof(flag))) {
		fprintf(stderr, "IFHANDLER:: Err: setsockopt fail for lwm2m gateway server.");
		ret = -1;
		goto cleanup;
	}

	/* assign IP, PORT */
	server.sin_family = AF_INET;
	server.sin_addr.s_addr = htonl(INADDR_ANY);
	server.sin_port = htons(IF_LWM2M_SERVER_PORT);

	 /* Binding newly created socket to given IP and verification */
	if (bind(socket_desc, (struct sockaddr *)&server, sizeof(server)) != 0) {
        ret = -1;
        perror("Socket Failure.");
        fprintf(stderr, "Err(%d): IFHANDLER:: Socket bind failed...\r\n", ret);
        goto cleanup;
    } else {
        fprintf(stdout, "IFHANDLER:: Socket successfully binded..\r\n");
    }

    /* Now server is ready to listen and verification */
    if ((listen(socket_desc, 5)) != 0) {
        ret = -1;
        fprintf(stderr, "Err(%d): IFHANDLER:: Listen failed...\r\n", ret);
        goto cleanup;
    } else {
        fprintf(stdout, "IFHANDLER::  Server is listening..\n");
    }

	c = sizeof(struct sockaddr_in);
	while (true)
	{
		while (client_sock = accept(socket_desc, (struct sockaddr *)&client, (socklen_t *)&c))
		{
			fprintf(stdout, "IFHANDLER:: Connection accepted\r\n");
			fflush(stdout);
			pthread_t subthread;
			new_sock = malloc(sizeof(int));
			*new_sock = client_sock;

			if (pthread_create(&subthread, NULL, request_handler, (void *)new_sock) < 0)
			{
				perror("IFHANDLER::  Could not create thread for client request.");
				ret = -1;
				break;
			}

			fprintf(stdout, "IFHANDLER:: handler assigned\r\n");
			fflush(stdout);
		}

		if (client_sock < 0) {
			perror("IFHANDLER:: Connection handler failed to accept message");
		}

		if (ret) {
			perror("IFHANDLER:: Connection failed to create thread for request.");
		}
		connection_handler_close_connection(client_sock);

	}

	cleanup:
	if (ret) {
		/* After service close the socket */
		connection_handler_close_connection(socket_desc);
	}

}

/* Start a connection handler thread. */
pthread_t connection_handler_start(void *data)
{
  pthread_t conn_id = 0;
  if (data)
  {
    if (pthread_create(&conn_id, NULL, &connection_handler, data))
    {
      /*Thread creation failed*/
      conn_id = 0;
    }

    // Create request list
    init_req_list();

    //Create notify event list
    init_evt_list();

    // initialize random number generator
    srand(time(NULL));
  }
  else
  {
    fprintf(stdout, "IFHANDLER:: Invalid lwm2mH context.\r\n");
    fflush(stdout);
  }
  return conn_id;
}

/* Stop Connection handler */
void connection_handler_stop(pthread_t thread)
{
  fprintf(stdout, "IFHANDLER:: Canceling thread %lld.\r\n", thread);
  fflush(stdout);
  pthread_cancel(thread);
  pthread_join(thread, NULL);
}
