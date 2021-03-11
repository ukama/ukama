/*
 * request.c
 *
 *  Created on: Mar 10, 2021
 *      Author: vishal
 */
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

#include "headers/errorcode.h"

#define DATA_OFFSET  (4*sizeof(uint32_t))

/* For storing incoming request */
static ListInfo req_list;


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

/* Delete the list node */
static void free_req_list_node(void* ip) {
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

/* Compare req list nodes */
static int compare_req_list_node(void *ipt, void *sd) {
	RequestMsg *ip = (RequestMsg *)ipt;
	RequestMsg *op = (RequestMsg *)sd;
	int ret = 0;

	/* Compare request id's */
	if (ip->reqid == op->reqid) {
		ret = 1;
	}

	return ret;
}

/* Copy request node */
static void *copy_req_node(void *pdata) {
	RequestMsg *data = pdata;
	RequestMsg *ndata = NULL;
	if (data) {

		ndata = malloc(sizeof(RequestMsg));
		if (ndata) {
			ndata->reqid =  data->reqid;
			ndata->length =  data->length;
			ndata->sock = data->sock;

			/* Try deep  copy for data now */
			ndata->msg = malloc( sizeof(char) * ndata->length );
			if (ndata->msg) {
				memcpy(ndata->msg, data->msg, ndata->length);
			}

		}

	}
	return ndata;
}



/* Add a new request to list */
static void add_request(RequestMsg* msg) {
#ifdef DEBUG_GATEWAYIF
	fprintf(stdout, "Adding Request id %llu to list.\r\n", msg->reqid);
	fflush(stdout);
#endif
	list_append(&req_list, msg);
}

/* Remove a request from the list */
static void remove_request(uint32_t id) {
	RequestMsg msg;
	msg.reqid = id;
#ifdef DEBUG_GATEWAYIF
	fprintf(stdout, "Removing Request id %llu to list.\r\n", msg.reqid);
	fflush(stdout);
#endif
	/* Remove element */
	list_remove(&req_list, &msg);
}

/* Search for Request in list */
static RequestMsg* search_req_id(uint32_t id) {
	RequestMsg msg;
	msg.reqid = id;
#ifdef DEBUG_GATEWAYIF
	fprintf(stdout, "Searching Request id %llu to list.\r\n", msg.reqid);
	fflush(stdout);
#endif
	return list_search(&req_list, &msg);
}

#ifdef DEBUG_GATEWAYIF
/* Print request node */
static void print_req_node(void* data) {
	if( data) {
		RequestMsg* msg = data;
		fprintf(stdout,"Request ID is: %llu\r\n", msg->reqid);
		fprintf(stdout,"Request message Length: %d\r\n", msg->length);
		fprintf(stdout,"Request message: %s\r\n", msg->msg);
		fflush(stdout);
	}
}

/* Print request list */
static void print_req_list() {
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



/* Read the response socket and respond */
static int get_response_socket(uint32_t id) {
	int sock = 0;

	RequestMsg* msg = search_req_id(id);
	if (msg) {
		sock =  msg->sock;
	} else {
		sock = -1;
	}

	return sock;
}

/* Deserialize  request message */
static RequestMsg* deserailize_request(char* rmsg) {
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

/* Serialize response for LwM2M gateway */
static char* serialize_response(ResponseMsg* resp, uint32_t *len){
	char *arr = NULL;
	if(resp) {
		*len = sizeof(ResponseMsg) + (sizeof(char)* resp->length);

		arr = malloc(*len);
		if (arr) {
			memset(arr, '\0', *len);
			memcpy(arr, resp, sizeof(ResponseMsg));
			if (resp->length) {
				memcpy(arr+(DATA_OFFSET), resp->msg, (sizeof(char)* resp->length));
			}
		}
	}
	return arr;
}

/* Transmitt response to LwM2M Gateway.*/
static int send_response(ResponseMsg* resp) {
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

/* Serve request */
void serve_request(char* msg, int sock) {
	uint32_t ret = 0;

	RequestMsg* rmsg = deserailize_request(msg);
	if (rmsg) {
		rmsg->sock = sock;
		fprintf(stdout, "Received request message id %llu, socket %d, length %d message %s.\r\n",
				rmsg->reqid, rmsg->sock, rmsg->length, rmsg->msg);
		fflush(stdout);

		/* Add to queue */
		add_request(rmsg);
	}

#ifdef DEBUG_GATEWAYIF
	/* Print list */
	print_req_list();
#endif
	/* Handle the request */
	ret = handle_server_req(rmsg);

	/* If got error than it would have non null values */
	if (ret) {
		fprintf(stdout, "Received request message is not parsed properly error %d.", ret);
		fflush(stdout);
		response_handler(rmsg->reqid, ret, 0, NULL, 0);
	}

	/* Free memory */
	if(rmsg) {
		if (rmsg->msg) {
			free(rmsg->msg);
			rmsg->msg = NULL;
		}
		free(rmsg);
		rmsg = NULL;
	}

}

/* Prepare_response */
int response_handler(uint32_t reqid, uint32_t status, uint32_t format, uint8_t* data, int length) {
	ResponseMsg msg = {0};
	msg.reqid = reqid;
	msg.status = status;
	msg.format = format;
	msg.length = length;
	msg.msg = data;
	send_response(&msg);
}

/* Initialize the request list */
int init_req_list() {
	list_new(&req_list, sizeof(RequestMsg), free_req_list_node, compare_req_list_node, copy_req_node);
}

