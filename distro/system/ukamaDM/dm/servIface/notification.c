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
#include "notification.h"
#include "headers/errorcode.h"

#define DEBUG_GATEWAYIF
/* For storing outgoing notification.*/
static ListInfo evt_list;

// Delete the node from event list
static void free_evt_list_node(void* ip) {
	ListNode *node = (ListNode *)ip;
	if (node) {

		NotifyEvent* evt =(NotifyEvent*)node->data;
		if (evt->data) {
			free(evt->data);
			evt->data = NULL;
		}
		free(node->data);
		node->data = NULL;
	}
	free(node);
	node = NULL;
}

//Compare evt list nodes
static int compare_evt_list_node(void *ipt, void *sd) {
	NotifyEvent *ip = (NotifyEvent *)ipt;
	NotifyEvent *op = (NotifyEvent *)sd;
    int ret = 0;

    /* Compare request id's */
    if (ip->id == op->id) {
        ret = 1;
    }

    return ret;
}

//Copy event node
static void *copy_evt_node(void *pdata) {
	NotifyEvent *data = pdata;
	NotifyEvent *ndata = NULL;
	if (data) {
		ndata = malloc(sizeof(NotifyEvent));
		if (ndata) {
			memset(ndata, '\0', sizeof(NotifyEvent));
			memcpy(ndata->uuid, data->uuid, ATTR_MAX_LEN);
			memcpy(ndata->uri, data->uri, ATTR_MAX_LEN);
			ndata->id =  data->id;
			ndata->count = data->count;
			ndata->format = data->format;
			ndata->length = data->length;

			/* Try deep  copy for data now */
			ndata->data = malloc( sizeof(char) * ndata->length );
			if (ndata->data) {
				memcpy(ndata->data, data->data, ndata->length);
			}

		}

	}
	return ndata;
}

/* Add a new event to list*/
static void add_event(NotifyEvent* evt) {
#ifdef DEBUG_GATEWAYIF
	fprintf(stdout, "IFHANDLER:: Adding Event id %llu to list.\r\n", evt->id);
	fflush(stdout);
#endif
	list_append(&evt_list, evt);
}

/* Remove a event from the list */
static void remove_event(uint32_t id) {
	NotifyEvent evt;
	evt.id = id;
#ifdef DEBUG_GATEWAYIF
	fprintf(stdout, "IFHANDLER:: Removing event id %llu to list.\r\n", evt.id);
	fflush(stdout);
#endif
	/* Remove element */
	list_remove(&evt_list, &evt);
}

/* Search for event in list */
static NotifyEvent* search_evt_id(uint32_t id) {
	NotifyEvent evt;
	evt.id = id;
#ifdef DEBUG_GATEWAYIF
	fprintf(stdout, "IFHANDLER:: Searching event id %llu to list.\r\n", evt.id);
	fflush(stdout);
#endif
	return list_search(&evt_list, &evt);
}

/* Serialize event notification.*/
static char *evt_serialize(NotifyEvent *evt, size_t *sz) {
    char *data = NULL;
    if (evt) {
        *sz = sizeof(NotifyEvent) + evt->length;

        data = malloc(*sz);
        if (data) {
            memset(data, '\0', *sz);
            memcpy(data, evt, sizeof(NotifyEvent));
            memcpy(&data[sizeof(NotifyEvent)-sizeof(uint8_t*)], evt->data, evt->length);
        } else {
            *sz = 0;
        }
    }
    return data;
}

/* Deserialize event */
static NotifyEventResp *evt_deserialize(char *data) {
	NotifyEventResp *resp = NULL;
    if (data) {
        resp = malloc(sizeof(NotifyEventResp));
        if (resp) {
            memset(resp, '\0', sizeof(NotifyEventResp));
            memcpy(resp, data, sizeof(NotifyEventResp));
        }
    }
    return resp;
}

/* Receive response for events */
static char *evt_recv(int sockfd) {
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
				fprintf(stderr, "Err: IFHANDLER:: Timeout for the receive functions.\r\n");
			}
		}
	}
	return req;
}

/* Send notification*/
static int evt_send(int sockfd, void *data, size_t size) {
    /* and send that buffer to client */
    return send(sockfd, data, size, 0);
}

/* Function designed for Asynchronous request between client and server.*/
static NotifyEventResp* evt_handler_client_func(int sockfd, NotifyEvent *evt) {
	int ret = 0;
	size_t size = 0;
	NotifyEventResp* rmsg = NULL;
	/* Serialize */
	char *sdata = evt_serialize(evt, &size);
	if (sdata) {

		ret = evt_send(sockfd, sdata, size);
		if (ret > 0) {

			char *rdata = evt_recv(sockfd);
			if (rdata) {
				/* Deserialize message */
				rmsg = evt_deserialize(rdata);
			} else {
				ret = ERR_SOCK_RECV;
			}

			/* Free receive data */
			if(rdata){
				free(rdata);
			}
		} else {
			ret = ERR_SOCK_SEND;
		}

		/* Free serialized data */
		if(sdata){
			free(sdata);
		}
	} else {
		ret = ERR_IFMSG_SERIALIZATION;
	}
	return rmsg;
}

/* Send notification */
static int send_notification(NotifyEvent *evt ) {

	NotifyEventResp* rmsg = NULL;
	int ret = 0;
	/* Create socket */
	int clientsock = connection_handler_create_sock();
	if (clientsock > 0 ) {
		/* connect to server */
		if (connection_handler_sock_connect(clientsock, LWM2M_GW_ADDRESS, LWM2M_GW_PORT)) {
			ret = ERR_SOCK_CREATION;
			goto cleanup;
		}
	} else {
		ret = ERR_SOCK_CREATION;
		goto cleanup;
	}

	if (!ret) {
		rmsg = evt_handler_client_func(clientsock, evt);
	}

	/* If status is ok. remove the event from list.*/
	if (rmsg->status == STATUS_OK) {
#ifdef DEBUG_GATEWAYIF
		fprintf(stdout, "IFHANDLER:: Received response for Event ID %d with response %d.\r\n", rmsg->id, rmsg->status);
		fflush(stdout);
#endif
		remove_event(rmsg->id);
		ret = 0;
	}

	cleanup:
	/* After service close the socket */
	connection_handler_close_connection(clientsock);
	return ret;

}

/* Notify Handler */
int notify_handler( char* uuid, uint8_t* uri, uint32_t count, uint32_t format, uint8_t* data, int length) {
	NotifyEvent evt = {'\0'};
	evt.id = rand();
	evt.count = count;
	evt.format = format;
	evt.length = length;
	memcpy(evt.uuid, uuid, strlen(uuid));
	memcpy(evt.uri, uri, strlen(uri));
	evt.data = data;

	/* Add event to list */
	add_event(&evt);

	/* Send Notification */
	send_notification(&evt);

}

/* Initialize the event list */
int init_evt_list() {
	list_new(&evt_list, sizeof(NotifyEvent), free_evt_list_node, compare_evt_list_node, copy_evt_node);
}



