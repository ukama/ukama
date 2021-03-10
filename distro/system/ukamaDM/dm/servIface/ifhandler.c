#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>
#include <pthread.h>

#include "list.h"
#include "ifhandler.h"

#define IF_LWM2M_SERVER_PORT 3000

#define CLIENT_REQ_MSG_LEN 512

#define DATA_OFFET  (4*sizeof(uint32_t))


static ListInfo req_list;

//Request handler function
void *request_handler(void *);

//Connection handler
void *connection_handler();

#ifdef DEBUG_GATEWAYIF
void ifhandler_print(char* data, int size) {
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
#endif

// Delete the list node
int free_req_list_node(void* ip) {
	ListNode *node = (ListNode *)ip;
	if (node) {
		RequestMsg* req =(RequestMsg*)node->data;
		if (req->msg) {
			free(req->msg);
			req->msg = NULL;
		}
		free(node->data);
		node->data = NULL;
	}
	free(node);
	node = NULL;
}

//Compare req list nodes
int compare_req_list_node(void *ipt, void *sd) {
	RequestMsg *ip = (RequestMsg *)ipt;
	RequestMsg *op = (RequestMsg *)sd;
    int ret = 0;
    /* Compare request id's */
    if (ip->reqid == op->reqid) {
        ret = 1;
    }
    return ret;
}

//Copy request node
void *copy_req_node(void *pdata) {
	RequestMsg *data = pdata;
	RequestMsg *ndata = NULL;
	if (data) {

		ndata = malloc(sizeof(RequestMsg));
		if (ndata) {
			ndata->reqid =  data->reqid;
			ndata->length =  data->length;
			ndata->sock = data->sock;
			/* Try deep  copy for properties of pdata now */
			ndata->msg = malloc( sizeof(char) * ndata->length );
			if (ndata->msg) {
				memcpy(ndata->msg, data->msg, ndata->length);
			}

		}

	}
	return ndata;
}

// Initialize the list
int init_req_list() {
	  list_new(&req_list, sizeof(RequestMsg), free_req_list_node, compare_req_list_node, copy_req_node);
}

// Add a new request to list
void add_request(RequestMsg* msg) {
#ifdef DEBUG_GATEWAYIF
	fprintf(stdout, "Adding Request id %llu to list.\r\n", msg->reqid);
	fflush(stdout);
#endif
	list_append(&req_list, msg);
}

//Remove a new request from the list
void remove_request(uint32_t id) {
	RequestMsg msg;
	msg.reqid = id;
#ifdef DEBUG_GATEWAYIF
	fprintf(stdout, "Removing Request id %llu to list.\r\n", msg.reqid);
	fflush(stdout);
#endif
	/* Remove element */
	list_remove(&req_list, &msg);
}

//Search for Request in list
RequestMsg* search_req_id(uint32_t id) {
	RequestMsg msg;
	msg.reqid = id;
#ifdef DEBUG_GATEWAYIF
	fprintf(stdout, "Searching Request id %llu to list.\r\n", msg.reqid);
	fflush(stdout);
#endif
	return list_search(&req_list, &msg);
}

#ifdef DEBUG_GATEWAYIF
// Print request node
void print_req_node(void* data) {
	if( data) {
		RequestMsg* msg = data;
		fprintf(stdout,"Request ID is: %llu\r\n", msg->reqid);
		fprintf(stdout,"Request message Length: %d\r\n", msg->length);
		fprintf(stdout,"Request message: %s\r\n", msg->msg);
		fflush(stdout);
	}
}

// Print request list
void print_req_list() {
    fprintf(stdout,
        "********************************************************************************\r\n");
    fprintf(stdout, "Length of the request list %llu.", req_list.logicalLength);
    fprintf(stdout,
        "****************************** Request list *************************************\r\n");
    fflush(stdout);
    if (req_list.logicalLength > 0) {

        list_for_each(&req_list, print_req_node);
        fprintf(stdout,
            "********************************************************************************\r\n");
        fflush(stdout);
    } else {
    	fprintf(stdout,"Request list is empty.\r\n");
    	fflush(stdout);
    }
}
#endif

// Read the response socket and respond
int get_response_socket(uint32_t id) {
	int sock = 0;
	RequestMsg* msg = search_req_id(id);
	if (msg) {
		sock =  msg->sock;
	} else {
		msg = -1;
	}
	return sock;
}

//Deserialize  request message
RequestMsg* deserailize_request(char* rmsg) {
	RequestMsg* req  = malloc(sizeof(RequestMsg));
	if (req) {
		req->reqid = *(uint32_t*)(rmsg);
		req->length = *(uint32_t*)(rmsg+sizeof(uint32_t));
		req->msg = malloc(sizeof(char)*req->length);
		if(req->msg) {
			memcpy(req->msg, rmsg+(2*sizeof(uint32_t)), req->length);
		}
	}
	return req;
}

//Serialize response for LwM2M gateway
char* serialize_response(ResponseMsg* resp, uint32_t *len){
	char *arr = NULL;
	if(resp) {
		*len = sizeof(ResponseMsg) + (sizeof(char)* resp->length);
		arr = malloc(*len);
		if (arr) {
			memset(arr, '\0', *len);
			memcpy(arr, resp, sizeof(ResponseMsg));
			if (resp->length) {
				memcpy(arr+(DATA_OFFET), resp->msg, (sizeof(char)* resp->length));
			}
		}
	}
	return arr;
}

// Transmitt response to LwM2M Gateway.
int send_response(ResponseMsg* resp) {
	uint32_t length = 0;
	int send_sock = -1;

	/* Get socket */
	send_sock = get_response_socket(resp->reqid);
	if (send_sock < 0) {
		fprintf(stdout, "Failed to respond for %llu request Id and status %d sent over sock %d sent.\r\n", resp->reqid,  resp->status, send_sock);
		fflush(stdout);
		goto cleanup;
	}

	/* Serialize response */
	char* respArray = serialize_response(resp, &length);
	if (respArray) {

#ifdef DEBUG_GATEWAYIF
		/* Debug message */
		ifhandler_print(respArray, length);
#endif
		/* Post response */
		int post = send(send_sock, respArray, length, 0);
		if (post>0) {
			fprintf(stdout, "Response message for %llu request Id with status %d type %d sent over sock %d sent.\r\n", resp->reqid,  resp->status, resp->format, send_sock);
			fflush(stdout);
		} else {
			fprintf(stdout, "Failed to respond for %llu request Id and status %d sent over sock %d sent.\r\n", resp->reqid,  resp->status, send_sock);
			fflush(stdout);
		}

	} else {
		fprintf(stdout, "No data available to respond for %llu request Id.\r\n", resp->reqid);
		fflush(stdout);
	}

	//close socket now.
	close(send_sock);

	cleanup:
	/* Remove Request id fom list */
	remove_request(resp->reqid);
#ifdef DEBUG_GATEWAYIF
	print_req_list();
#endif


}

// Serve request
void serve_request(char* msg, int sock) {
	uint32_t ret = 0;
	RequestMsg* rmsg = deserailize_request(msg);
	if (rmsg) {
		rmsg->sock = sock;
		fprintf(stdout, "Received request message id %llu, socket %d, length %d message %s.\r\n",
				rmsg->reqid, rmsg->sock, rmsg->length, rmsg->msg);
		fflush(stdout);

		//Add to queue
		add_request(rmsg);
	}

#ifdef DEBUG_GATEWAYIF
	//Print list
	print_req_list();
#endif
	//Handle the request
	ret = handle_server_req(rmsg);

  // If got error than it would have non null values
	if (ret) {
		fprintf(stdout, "Received request message is not parsed properly error %d.", ret);
		fflush(stdout);
		response_handler(rmsg->reqid, ret, 0, NULL, 0);
	}

	// Free memory
	if(rmsg) {
		if (rmsg->msg) {
			free(rmsg->msg);
			rmsg->msg = NULL;
		}
		free(rmsg);
		rmsg = NULL;
	}

}

// Prepare_response
int response_handler(uint32_t reqid, uint32_t status, uint32_t format, uint8_t* data, int length) {
	ResponseMsg msg = {0};
	msg.reqid = reqid;
	msg.status = status;
	msg.format = format;
	msg.length = length;
	msg.msg = data;
	send_response(&msg);
}


/*
  Handles the request.
  */
void *request_handler(void *socket_desc)
{
  //Get the socket descriptor
  int sock = *(int *)socket_desc;
  int n;

  char req_arr[CLIENT_REQ_MSG_LEN];
  fprintf(stdout, "Waiting for incoming message Request handler thread id %lld\r\n", pthread_self());
  fflush(stdout);
  while ((n = recv(sock, req_arr, CLIENT_REQ_MSG_LEN, 0)) > 0)
  {
#ifdef DEBUG_GATEWAYIF
	  ifhandler_print(req_arr, CLIENT_REQ_MSG_LEN);
#endif
	  // Server request
	  serve_request(req_arr, sock);
  }

  fprintf(stdout, "Exiting thread id %lld\r\n", pthread_self());
  fflush(stdout);
  pthread_exit(NULL);
}

/*
  Handles the incoming connection form the client
  */
void *connection_handler()
{
  int socket_desc, client_sock, c, *new_sock;
  struct sockaddr_in server, client;

  //Create socket
  socket_desc = socket(AF_INET, SOCK_STREAM, 0);
  if (socket_desc == -1)
  {
    fprintf(stdout, "Could not create socket\r\n");
  }
  fprintf(stdout, "Socket created\r\n");

  int flag = 1;
  if (-1 == setsockopt(socket_desc, SOL_SOCKET, (SO_REUSEADDR|SO_REUSEPORT), &flag, sizeof(flag))) {
  	fprintf(stderr, "Err: setsockopt fail for lwm2m server.");
  	return;
  }

  //Prepare the sockaddr_in structure
  server.sin_family = AF_INET;
  server.sin_addr.s_addr = INADDR_ANY;
  server.sin_port = htons(IF_LWM2M_SERVER_PORT);

  //Bind
  if (bind(socket_desc, (struct sockaddr *)&server, sizeof(server)) < 0)
  {
    //print the error message
    perror("bind failed. Error");
    return;
  }
  fprintf(stdout, "bind done\r\n");

  //Listen
  listen(socket_desc, 3);

  //Accept and incoming connection
  fprintf(stdout, "Waiting for incoming connections...\r\n");
  fflush(stdout);
  c = sizeof(struct sockaddr_in);
  while (true)
  {
    while (client_sock = accept(socket_desc, (struct sockaddr *)&client, (socklen_t *)&c))
    {
      fprintf(stdout, "Connection accepted\r\n");
      fflush(stdout);
      pthread_t subthread;
      new_sock = malloc(1);
      *new_sock = client_sock;

      if (pthread_create(&subthread, NULL, request_handler, (void *)new_sock) < 0)
      {
        perror("could not create thread");
        return 1;
      }

      fprintf(stdout, "Handler assigned\r\n");
      fflush(stdout);
    }
    if (client_sock < 0)
    {
      perror("accept failed");
    }
  }
}

/*
Start a connection handler thread.
*/
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
  }
  else
  {
    fprintf(stdout, "IFHANDLER:: Invalid lwm2mH context.\r\n");
    fflush(stdout);
  }
  return conn_id;
}

/*
Stop Connection handler
*/
void connection_handler_stop(pthread_t thread)
{
  fprintf(stdout, "IFHANDLER:: Canceling thread %lld.\r\n", thread);
  fflush(stdout);
  pthread_cancel(thread);
  pthread_join(thread, NULL);
}
