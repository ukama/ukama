/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <string.h>
#include <stdlib.h>
#include <stdio.h>

#include "rabbitmq-c/amqp.h"

#include "mesh.h"
#include "u_amqp.h"
#include "link.pb-c.h"

/* 
 * AMQP Routing key:
 * <type>.<source>.<container>.<object>.<state>
 *
 * type:       event, request, response
 * source:     cloud, device
 * container:  mesh
 * object:     link, cert
 * state:      connect, fail, active, lost, end, close, valid, invalid, update
 *             expired
 *
 */

static char *convert_type_to_str(MsgType type);
static char *convert_source_to_str(MsgSource source);
static char *convert_object_to_str(MsgObject object);
static char *convert_state_to_str(ObjectState state);
static int is_valid_event(MeshEvent event);
static char *create_routing_key(MeshEvent event);
static void *serialize_link_msg(char *nodeID);
static int object_type(MeshEvent event);

/* Mapping between Mesh.d internal state and AMQP routing key. 
 *
 * internal    object state   type
 *
 * CONN_CONNECT link  connect  event
 * CONN_FAILED  link  fail     event
 * CONN_ACTIVE  link  active   event
 * CONN_LOST    link  lost     event
 * CONN_END     link  end      event
 * CONN_CLOSE   link  close    event
 *
 * CERT_OK       cert  valid    event
 * CERT_EXPIRED  cert  expired  event
 * CERT_INVALID  cert  invalid  event
 * CERT_REQUIRED cert  update   request
 *
 */
static AMQPRoutingKey routingKey[] = {

	[CONN_CONNECT]  = {.type=EVENT, .object=LINK, .state=CONNECT},
	[CONN_FAIL]     = {.type=EVENT, .object=LINK, .state=FAIL},
	[CONN_ACTIVE]   = {.type=EVENT, .object=LINK, .state=ACTIVE},
	[CONN_LOST]     = {.type=EVENT, .object=LINK, .state=LOST},
	[CONN_END]      = {.type=EVENT, .object=LINK, .state=END},
	[CONN_CLOSE]    = {.type=EVENT, .object=LINK, .state=CLOSE},

	[CERT_OK]       = {.type=EVENT, .object=CERT, .state=VALID},
	[CERT_EXPIRED]  = {.type=EVENT, .object=CERT, .state=EXPIRED},
	[CERT_INVALID]  = {.type=EVENT, .object=CERT, .state=INVALID},
	[CERT_REQUIRED] = {.type=REQUEST, .object=CERT, .state=UPDATE},
};

/* 
 * Free -- free variable list of arguments
 *
 */
static void Free(void *ptr, ... ) {

	void *p;
  
	if(!ptr)
		return;

	va_list list;
	va_start(list, ptr);

	p = va_arg(list , void *);
  
	while( p != ptr ) {
		if (p) free(p) ;
		p = va_arg(list, void *);
	}

	va_end(list);
}

/*
 * convert_type_to_str -- convert passed routing key elements into char*
 *
 */
static char *convert_type_to_str(MsgType type) {

	char *str;

	switch(type) {

	case EVENT:
		str=TYPE_EVENT_STR;
		break;

	case REQUEST:
		str=TYPE_REQUEST_STR;
		break;

	case RESPONSE:
		str=TYPE_RESPONSE_STR;
		break;

	default:
		return NULL;
	}

	return strdup(str);
}

/*
 * convert_source_to_str --
 *
 */
static char *convert_source_to_str(MsgSource source) {

	char *str;

	switch(source) {

	case DEVICE:
		str = SOURCE_DEVICE_STR;
		break;

	case CLOUD:
		str = SOURCE_CLOUD_STR;
		break;

	default:
		return NULL;
	}

	return strdup(str);
}

/*
 * convert_object_to_str --
 *
 */
static char *convert_object_to_str(MsgObject object) {

	char *str;

	switch(object) {

	case LINK:
		str = OBJECT_LINK_STR;
		break;

	case CERT:
		str = OBJECT_CERT_STR;
		break;

	default:
		return NULL;
	}

	return strdup(str);
}

/*
 * convert_state_to_str --
 *
 */
static char *convert_state_to_str(ObjectState state) {

	char *str;

	switch(state) {

	case CONNECT:
		str = STATE_CONNECT_STR;
		break;

	case FAIL:
		str = STATE_FAIL_STR;
		break;

	case ACTIVE:
		str = STATE_ACTIVE_STR;
		break;

	case LOST:
		str = STATE_LOST_STR;
		break;

	case END:
		str = STATE_END_STR;
		break;

	case CLOSE:
		str = STATE_CLOSE_STR;
		break;

	case VALID:
		str = STATE_VALID_STR;
		break;

	case INVALID:
		str = STATE_INVALID_STR;
		break;

	case EXPIRED:
		str = STATE_EXPIRED_STR;
		break;

	case UPDATE:
		str = STATE_UPDATE_STR;
		break;

	default:
		return NULL;
	}

	return strdup(str);
}

/*
 * is_valid_event --
 *
 */
static int is_valid_event(MeshEvent event) {

	int ret=FALSE;

	switch(event) {

	case CONN_CONNECT:
	case CONN_FAIL:
	case CONN_ACTIVE:
	case CONN_LOST:
	case CONN_END:
	case CONN_CLOSE:
	case CERT_OK:
	case CERT_EXPIRED:
	case CERT_INVALID:
	case CERT_REQUIRED:
		ret=TRUE;
		break;
	default:
		ret=FALSE;
	}
        
	return ret;
}

/*
 * log_amqp_response -- inspired from die_on_amqp_error()
 *
 */
static void log_amqp_response(WAMQPReply reply, const char *context) {

	switch (reply.reply_type) {
	case AMQP_RESPONSE_NORMAL:
		break;

	case AMQP_RESPONSE_NONE:
		log_error("%s: missing RPC reply type!", context);
		break;

	case AMQP_RESPONSE_LIBRARY_EXCEPTION:
		log_error("%s: %s\n", context, amqp_error_string2(reply.library_error));
		break;

	case AMQP_RESPONSE_SERVER_EXCEPTION:
		switch (reply.reply.id) {
		case AMQP_CONNECTION_CLOSE_METHOD:
			log_error("%s: server connection error", context);
			break;
      
		case AMQP_CHANNEL_CLOSE_METHOD:
			log_error("%s: server channel error", context);
			break;
      
		default:
			log_error("%s: unknown server error, method id 0x%08X", context,
					  reply.reply.id);
			break;
		} /*nested switch */
		break;

	default:
		log_error("%s: unknown error type", context);
		break;
	}
}

/*
 * init_amqp_connection -- 
 *
 */

WAMQPConn *init_amqp_connection(char *host, char *port) {

	int ret;
	WAMQPConn *conn=NULL;
	WAMQPSocket *socket=NULL;
	WAMQPReply reply;

	/* Sanity check */
	if (host == NULL && port == NULL) {
		return NULL;
	}

	/* Initialize AMQP state variable */
	conn = amqp_new_connection();
	if (!conn) {
		log_error("Unable to create AMQP, memory allocation issue");
		return NULL;
	}

	/* Create TCP socket */
	socket = amqp_tcp_socket_new(conn);
	if (!socket) {
		log_error("Unable to create AMQP TCP socket");
		amqp_destroy_connection(conn);
		return NULL;
	}

	/* Connect to the AMQP host */
	ret = amqp_socket_open(socket, host, atoi(port));
	if (ret) {
		log_error("Unable to connect with AMQP server at host: %s port: %s",
				  host, port);
		amqp_destroy_connection(conn);
		return NULL;
	}

	/* Login using guest/guest (for now - XXX) */
	reply = amqp_login(conn, "/", 0, MAX_FRAME, 0, AMQP_SASL_METHOD_PLAIN,
					   "guest", "guest");
	if (reply.reply_type != AMQP_RESPONSE_NORMAL) {
		log_amqp_response(reply, "AMQP login");
		amqp_channel_close(conn, 1, AMQP_CHANNEL_ERROR);
		amqp_destroy_connection(conn);
		return NULL;
	}
  
	/* Open the channel */
	amqp_channel_open(conn, 1);
	reply = amqp_get_rpc_reply(conn);
	if (reply.reply_type != AMQP_RESPONSE_NORMAL) {
		log_amqp_response(reply, "AMQP channel open");
		amqp_channel_close(conn, 1, AMQP_CHANNEL_ERROR);
		amqp_destroy_connection(conn);
		return NULL;
	}
  
	return conn;
}

/*
 * close_amqp_connection --
 *
 */
void close_amqp_connection(WAMQPConn *conn) {

	amqp_channel_close(conn, 1, AMQP_REPLY_SUCCESS);
	amqp_connection_close(conn, AMQP_REPLY_SUCCESS);
	amqp_destroy_connection(conn);
}

/*
 * create_routing_key --
 *
 */
static char *create_routing_key(MeshEvent event) {

	int len;
	char *key=NULL;
	char *type=NULL, *source=NULL, *container=NULL, *object=NULL, *state=NULL;

	/* Sanity check */
	if (!is_valid_event(event)) {
		return NULL;
	}

	/* Step-1: build the routing key for the event. 
	 * <type>.<source>.<container>.<object>.<state>
	 */
	type   = convert_type_to_str(routingKey[event].type);
	object = convert_object_to_str(routingKey[event].object);
	state  = convert_state_to_str(routingKey[event].state);

	if (type==NULL || object==NULL || state==NULL) {
		FREE(type, object, state);
		return NULL;
	}

	source = convert_source_to_str((MsgSource)CLOUD);

	len = strlen(type) + strlen(source) + strlen(MSG_CONTAINER) +
		strlen(object) + strlen(state);

	key = (char *)malloc(len+4+1); /* 4 for '.' in the key and 1 for' \0' */
	if (key==NULL) {
		log_error("Error allocating memory of size: %d", len+1);
		FREE(source, type, object, state);
		return NULL;
	}

	sprintf(key, "%s.%s.%s.%s.%s", type, source, MSG_CONTAINER, object, state);

	FREE(type, source, container, object, state);

	return key;
}

/*

 * serialize_link_msg -- Serialize the protobuf msg for the Link object
 *
 */
static void *serialize_link_msg(char *nodeID) {

	Link linkMsg = LINK__INIT;
	void *buff=NULL;
	size_t len, idLen=36+1;

	if (nodeID == NULL) {
		return NULL;
	}

    linkMsg.uuid = strdup(nodeID);

	len = link__get_packed_size(&linkMsg);

	buff = malloc(len);
	if (buff==NULL) {
		log_error("Error allocating buffer of size: %d", len);
		return NULL;
	}

	link__pack(&linkMsg, buff);

	free(linkMsg.uuid);

	return buff;
}

/*
 * object_type -- return object type of given event
 *
 */
static int object_type(MeshEvent event) {

	int type;

	switch(event) {

		/* Link object */
	case CONN_CONNECT:
	case CONN_FAIL:
	case CONN_ACTIVE:
	case CONN_LOST:
	case CONN_END:
	case CONN_CLOSE:
		type = OBJECT_LINK;
		break;

		/* Cert object */
	case CERT_OK:
	case CERT_EXPIRED:
	case CERT_INVALID:
	case CERT_REQUIRED:
		type = OBJECT_CERT;
		break;

	default:
		type = OBJECT_NONE;
		break;
	}

	return type;
}

/*
 * publish_amqp_event --
 *
 */
int publish_amqp_event(WAMQPConn *conn, char *exchange, MeshEvent event,
					   char *nodeID) {

	/* THREAD? XXX - Think about me*/
	char *key=NULL;
	WAMQPProp prop;
	void *buff=NULL;
	int ret;

	/* Sanity check */
	if (conn==NULL) {
		return FALSE;
	}

    if (nodeID == NULL) {
        return FALSE;
    }

	/* Step-1: build the routing key for the event. 
	 * <type>.<source>.<container>.<object>.<state>
	 */
	key = create_routing_key(event);
	if (key == NULL) {
		log_error("Error creating routing key. Ignoring the message");
		return FALSE;
	} else {
		log_debug("Routing key created: %s", key);
	}

	/* Step-2: setup AMQP message properties */
	prop._flags = AMQP_BASIC_CONTENT_TYPE_FLAG | AMQP_BASIC_DELIVERY_MODE_FLAG;
	prop.content_type = amqp_cstring_bytes("text/plain");
	prop.delivery_mode = 2; /* persistent delivery mode */

	/* Step-3: protobuf msg. */
	if (object_type(event) == OBJECT_LINK) {

		buff = serialize_link_msg(nodeID);
		if (buff==NULL) {
			log_error("Error serializing Link packet for AMQP. Event: %d",
					  event);
			free(key);
			return FALSE;
		}
	} else if (object_type(event) == OBJECT_NONE) {
		log_error("Invalid event type: %d", event);
		free(key);
		return FALSE;
	}

	/* Step-4: send the message to AMQP broker */
	ret = amqp_basic_publish(conn, 1, amqp_cstring_bytes(exchange),
							 amqp_cstring_bytes(key), 0, 0, &prop,
							 amqp_cstring_bytes(buff));
	if (ret < 0) {
		ret = FALSE;
		log_error("Error sending AMQP message. Error: %s",
				  amqp_error_string2(ret));
	} else {
		ret = TRUE;
		log_debug("AMQP message successfully sent to exchange");
	}

	free(buff);
	free(key);
	return ret;
}
