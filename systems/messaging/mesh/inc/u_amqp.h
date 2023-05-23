/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef MESH_AMQP_H
#define MESH_AMQP_H

#include "mesh.h"
#include <rabbitmq-c/amqp.h>

#define MSG_CONTAINER "MeshD" /* use in AMQP routing key */

#define MAX_FRAME  131072
#define MAX_EVENTS 10

/* type of object */
#define OBJECT_NONE 0
#define OBJECT_LINK 1
#define OBJECT_CERT 2

#define TYPE_EVENT_STR    "event"
#define TYPE_REQUEST_STR  "request"
#define TYPE_RESPONSE_STR "response"

#define SOURCE_DEVICE_STR "device"
#define SOURCE_CLOUD_STR  "cloud"

#define OBJECT_LINK_STR "link"
#define OBJECT_CERT_STR "cert"

#define STATE_CONNECT_STR "connect"
#define STATE_FAIL_STR    "fail"
#define STATE_ACTIVE_STR  "active"
#define STATE_LOST_STR    "lost"
#define STATE_END_STR     "end"
#define STATE_CLOSE_STR   "close"
#define STATE_VALID_STR   "valid"
#define STATE_INVALID_STR "invalid"
#define STATE_UPDATE_STR  "update"
#define STATE_EXPIRED_STR "expired"

static int free_stop;
#define FREE( ... ) Free( &free_stop , __VA_ARGS__ , &free_stop )

typedef amqp_basic_properties_t WAMQPProp;
typedef struct amqp_connection_state_t_ WAMQPConn;
typedef struct amqp_socket_t_           WAMQPSocket;
typedef struct amqp_rpc_reply_t_        WAMQPReply;

/* AMQP Routing key format:
 * <type>.<source>.<container>.<object>.<state>
 */

/* Enum all Mesh.d internal events. */
typedef enum {

	CONN_CONNECT=0,
	CONN_FAIL,
	CONN_ACTIVE,
	CONN_LOST,
	CONN_END,
	CONN_CLOSE,

	CERT_OK,
	CERT_EXPIRED,
	CERT_INVALID,
	CERT_REQUIRED,

	MAX_EVENT=10, /* Always is the last and is total number of events. */
}MeshEvent;

typedef enum {

	DEVICE=1,
	CLOUD,
} MsgSource;

typedef enum {

	EVENT=1,
	REQUEST,
	RESPONSE,
} MsgType;

typedef enum {

	LINK=1,
	CERT,
} MsgObject;

typedef enum {

	CONNECT=1,
	FAIL,
	ACTIVE,
	LOST,
	END,
	CLOSE,
	VALID,
	INVALID,
	UPDATE,
	EXPIRED,
} ObjectState;

typedef struct _routing_key {

	MsgType     type;   /* Msg type is either event, request or response */
	MsgObject   object; /* Object msg is representing: link or cert */
	ObjectState state;  /* State of the object. */
} AMQPRoutingKey;

WAMQPConn *init_amqp_connection(char *host, char *port);
void close_amqp_connection(WAMQPConn *conn);
int publish_amqp_event(WAMQPConn *conn, char *exchange, MeshEvent event,
                       char *nodeID);

#endif /* MESH_AMQP_H */
