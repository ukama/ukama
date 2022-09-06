/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef MSG_CLIENT_AMQP_H
#define MSG_CLIENT_AMQP_H

#include <rabbitmq-c/amqp.h>

#define MSG_CONTAINER "MeshD" /* use in AMQP routing key */

#define MAX_FRAME  131072
#define MAX_EVENTS 10

#define TRUE  1
#define FALSE 0

/* Routing key format:
 * <type>.<source>.<container>.<object>.<state>
 */

/* type of object */
#define OBJECT_NONE    0
#define OBJECT_SERVICE 1
#define OBJECT_LINK    2
#define OBJECT_CERT    3

#define TYPE_EVENT_STR     "event"
#define TYPE_REQUEST_STR   "request"
#define TYPE_RESPONSE_STR  "response"

#define SOURCE_DEVICE_STR  "device"
#define SOURCE_SERVICE_STR "service"

#define OBJECT_GENERAL_STR "general"
#define OBJECT_LINK_STR    "link"
#define OBJECT_CERT_STR    "cert"

#define STATE_REGISTER_STR   "register"
#define STATE_UNREGISTER_STR "unregister"
#define STATE_CONNECT_STR    "connect"
#define STATE_FAIL_STR       "fail"
#define STATE_ACTIVE_STR     "active"
#define STATE_LOST_STR       "lost"
#define STATE_END_STR        "end"
#define STATE_CLOSE_STR      "close"
#define STATE_VALID_STR      "valid"
#define STATE_INVALID_STR    "invalid"
#define STATE_UPDATE_STR     "update"
#define STATE_EXPIRED_STR    "expired"



static int free_stop;
#define FREE( ... ) Free( &free_stop , __VA_ARGS__ , &free_stop )

typedef struct amqp_basic_properties_t_ WAMQPProp;
typedef struct amqp_connection_state_t_ WAMQPConn;
typedef struct amqp_socket_t_           WAMQPSocket;
typedef struct amqp_rpc_reply_t_        WAMQPReply;

/* AMQP Routing key format:
 * <type>.<source>.<container>.<object>.<state>
 */

/* Enum all events. */
typedef enum {

	SRVC_REGISTER=0,
	SRVC_UNREGISTER,

	CONN_CONNECT,
	CONN_FAIL,
	CONN_ACTIVE,
	CONN_LOST,
	CONN_END,
	CONN_CLOSE,

	CERT_OK,
	CERT_EXPIRED,
	CERT_INVALID,
	CERT_REQUIRED,
	
	MAX_EVENT=12, /* Always is the last and is total number of events. */
} MsgClientEvent;

typedef enum {

	DEVICE=1,
	SERVICE,
} MsgSource;

typedef enum {

	EVENT=1,
	REQUEST,
	RESPONSE,
} MsgType;

typedef enum {

	GENERAL=1,
	LINK,
	CERT,
} MsgObject;

typedef enum {

	REGISTER=1,
	UNREGISTER,
	CONNECT,
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
	MsgObject   object; /* Object msg is representing */
	ObjectState state;  /* State of the object. */
} AMQPRoutingKey;

WAMQPConn *init_amqp_connection(char *host, char *port, char *login,
								char *passwd);
void close_amqp_connection(WAMQPConn *conn);
int publish_amqp_event(WAMQPConn *conn, char *exchange, MsgClientEvent event,
					   char *name);

#endif /* MSG_CLIENT_AMQP_H */
