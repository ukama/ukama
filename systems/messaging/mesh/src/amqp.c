/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <stdarg.h>

#include "rabbitmq-c/amqp.h"

#include "mesh.h"
#include "u_amqp.h"
#include "nodeEvent.pb-c.h"
#include "meshRegisterEvent.pb-c.h"
#include "any.pb-c.h"

typedef Google__Protobuf__Any ANY;
typedef Ukama__Events__V1__NodeOnlineEvent  NodeOnlineEvent;
typedef Ukama__Events__V1__NodeOfflineEvent NodeOfflineEvent;

/* Mapping between Mesh.d internal state and AMQP routing key.
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
static char *create_routing_key(MeshEvent event, char *orgName);
static int object_type(MeshEvent event);
static void *serialize_register_event(char *ip, int port, size_t *outLen);
static void *serialize_node_online_event(char *nodeID, char *nodeIP, int nodePort,
                                         char *meshIP, int meshPort,
                                         size_t *outLen);
static void *serialize_node_offline_event(char *nodeID, size_t *outLen);
static void *serialize_any_packet(int eventType, size_t len, void *buff,
                                  size_t *outLen);

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

static void Free(void *ptr, ...) {
    va_list ap;
    void *p;

    if (!ptr) {
        return;
    }

    va_start(ap, ptr);

    /* free the first pointer */
    p = ptr;
    while (p != NULL) {
        free(p);
        p = va_arg(ap, void *);
    }

    va_end(ap);
}

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

static char *convert_object_to_str(MsgObject object) {

	char *str;

	switch(object) {

	case LINK:
		str = OBJECT_NODE_STR;
		break;

	case CERT:
		str = OBJECT_CERT_STR;
		break;

	default:
		return NULL;
	}

	return strdup(str);
}

static char *convert_state_to_str(ObjectState state) {

	char *str;

	switch(state) {

	case CONNECT:
    case ACTIVE:
        str = STATE_ONLINE_STR;
        break;

	case FAIL:
	case LOST:
	case END:
	case CLOSE:
		str = STATE_OFFLINE_STR;
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
		}
		break;

	default:
		log_error("%s: unknown error type", context);
		break;
	}
}

static WAMQPConn *init_amqp_connection(char *host, char *port, char *user, char *password) {

	int ret;
	WAMQPConn *conn=NULL;
	WAMQPSocket *socket=NULL;
	WAMQPReply reply;

	/* Sanity check */
	if (host == NULL || port == NULL || user == NULL || password == NULL) {
		log_error("Invalid AMQP connection parameters: host: %s port: %s user: %s password: ****",
				  host, port, user);
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
		log_error("Unable to connect with AMQP server at host: %s port: %s user: %s password: ****",
				  host, port, user);
		amqp_destroy_connection(conn);
		return NULL;
	}

	/* Login using user/password env */
	reply = amqp_login(conn, "/", 0, MAX_FRAME, 0, AMQP_SASL_METHOD_PLAIN,
					   user, password);
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

static char *create_routing_key(MeshEvent event, char *orgName) {

	int len;
	char *key=NULL;
	char *type=NULL, *source=NULL, *object=NULL, *state=NULL;

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
		Free(type, object, state, NULL);
		return NULL;
	}

	source = convert_source_to_str((MsgSource)CLOUD);

	len = strlen(type) + strlen(source) + strlen(LOCAL_AMQP) +
        strlen(orgName) + strlen(SYSTEM_NAME) + strlen(MSG_CONTAINER) +
		strlen(object) + strlen(state);

	key = (char *)malloc(len+7+1); /* 7 for '.' in the key and 1 for' \0' */
	if (key==NULL) {
		log_error("Error allocating memory of size: %d", len+1);
		Free(source, type, object, state, NULL);
		return NULL;
	}

    /* event.cloud.local.orgName.messaging.mesh.node.state */
	sprintf(key, "%s.%s.%s.%s.%s.%s.%s.%s",
            type,
            source,
            LOCAL_AMQP,
            orgName,
            SYSTEM_NAME,
            MSG_CONTAINER,
            object,
            state);

	Free(type, source, object, state, NULL);

	return key;
}

static void *serialize_any_packet(int eventType, size_t len, void *buff, size_t *outLen) {

    ANY anyEvent = GOOGLE__PROTOBUF__ANY__INIT;
    void *anyBuff = NULL;
    size_t anyLen;

    if (eventType == CONN_CLOSE || eventType == CONN_END ||
        eventType == CONN_LOST  || eventType == CONN_FAIL) {

        anyEvent.type_url = (char *)calloc(strlen(TYPE_URL_PREFIX) + 1 +
                     strlen(ukama__events__v1__node_offline_event__descriptor.name) + 1,
                                           sizeof(char));
        sprintf(anyEvent.type_url, "%s/%s",
                TYPE_URL_PREFIX,
                ukama__events__v1__node_offline_event__descriptor.name);
    } else if (eventType == CONN_CONNECT) {

        anyEvent.type_url = (char *)calloc(strlen(TYPE_URL_PREFIX) + 1 +
                      strlen(ukama__events__v1__node_online_event__descriptor.name) + 1,
                                           sizeof(char));
        sprintf(anyEvent.type_url, "%s/%s",
                TYPE_URL_PREFIX,
                ukama__events__v1__node_online_event__descriptor.name);
    } else if (eventType == MESH_REGISTER) {

        anyEvent.type_url = (char *)calloc(strlen(TYPE_URL_PREFIX) + 1 +
                                           strlen(mesh_register_event__descriptor.name) + 1,
                                           sizeof(char));
        sprintf(anyEvent.type_url, "%s/%s",
                TYPE_URL_PREFIX,
                mesh_register_event__descriptor.name);
    } else {
        return NULL;
    }

    anyEvent.value.len  = len;
    anyEvent.value.data = malloc(len);
    memcpy(anyEvent.value.data, buff, len);

    anyLen = google__protobuf__any__get_packed_size(&anyEvent);
    anyBuff = malloc(anyLen);
    if (anyBuff == NULL) {
        log_error("Error allocating buffer of size: %zu", anyLen);
        return NULL;
    }

    google__protobuf__any__pack(&anyEvent, anyBuff);

    if (outLen) {
        *outLen = anyLen;
    }

    free(anyEvent.type_url);
    free(anyEvent.value.data);

    return anyBuff;
}

static void *serialize_node_online_event(char *nodeID, char *nodeIP, int nodePort,
                                         char *meshIP, int meshPort,
                                         size_t *outLen) {

	NodeOnlineEvent nodeEvent = UKAMA__EVENTS__V1__NODE_ONLINE_EVENT__INIT;
	void *buff=NULL, *anyBuff = NULL;
	size_t len, anyLen = 0;

	if (nodeID == NULL || nodeIP == NULL || meshIP == NULL) return NULL;

	nodeEvent.nodeid   = strdup(nodeID);
	nodeEvent.nodeip   = strdup(nodeIP);
	nodeEvent.nodeport = nodePort;
	nodeEvent.meship   = strdup(meshIP);
	nodeEvent.meshport = meshPort;
	nodeEvent.meshhostname = strdup("localhost");

	len = ukama__events__v1__node_online_event__get_packed_size(&nodeEvent);

	buff = malloc(len);
	if (buff==NULL) {
		log_error("Error allocating buffer of size: %zu", len);
		return NULL;
	}

	ukama__events__v1__node_online_event__pack(&nodeEvent, buff);

    anyBuff = serialize_any_packet(CONN_CONNECT, len, buff, &anyLen);

    free(nodeEvent.nodeid);
	free(nodeEvent.nodeip);
	free(nodeEvent.meship);
    free(nodeEvent.meshhostname);
    free(buff);

    if (outLen) {
        *outLen = anyLen;
    }

	return anyBuff;
}

static void *serialize_node_offline_event(char *nodeID, size_t *outLen) {

	NodeOfflineEvent nodeEvent = UKAMA__EVENTS__V1__NODE_OFFLINE_EVENT__INIT;
	void *buff = NULL, *anyBuff = NULL;
	size_t len, anyLen = 0;

	if (nodeID == NULL) return NULL;

	nodeEvent.nodeid   = strdup(nodeID);
	len = ukama__events__v1__node_offline_event__get_packed_size(&nodeEvent);

	buff = malloc(len);
	if (buff==NULL) {
		log_error("Error allocating buffer of size: %zu", len);
		return NULL;
	}

	ukama__events__v1__node_offline_event__pack(&nodeEvent, buff);

    anyBuff = serialize_any_packet(CONN_CLOSE, len, buff, &anyLen);

    free(nodeEvent.nodeid);
    free(buff);

    if (outLen) {
        *outLen = anyLen;
    }

	return anyBuff;
}

static void *serialize_register_event(char *ip, int port, size_t *outLen) {

    MeshRegisterEvent registerEvent = MESH_REGISTER_EVENT__INIT;
    void *buff=NULL, *anyBuff=NULL;
    size_t len, anyLen = 0;

    registerEvent.ip   = strdup(ip);
    registerEvent.port = port;

    len = mesh_register_event__get_packed_size(&registerEvent);

    buff = malloc(len);
    if (buff == NULL) {
        log_error("Error allocating buffer of size: %zu", len);
        return NULL;
    }

    mesh_register_event__pack(&registerEvent, buff);
    anyBuff = serialize_any_packet(MESH_REGISTER, len, buff, &anyLen);

    free(registerEvent.ip);
    free(buff);

    if (outLen) {
        *outLen = anyLen;
    }

	return anyBuff;
}

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

static int publish_amqp_event(WAMQPConn *conn, char *exchange, MeshEvent event,
                              char *orgName, char *nodeID, char *nodeIP, int nodePort,
                              char *meshIP, int meshPort) {

	char *key=NULL;
	WAMQPProp prop;
	void *buff=NULL;
	size_t buffLen = 0;
	int ret;

	/* Step-1: build the routing key for the event. 
	 * <type>.<source>.<container>.<object>.<state>
	 */
	key = create_routing_key(event, orgName);
	if (key == NULL) {
		log_error("Error creating routing key. Ignoring the message");
		return FALSE;
	} else {
		log_debug("Routing key created: %s", key);
	}

	/* Step-2: setup AMQP message properties */
	prop._flags = AMQP_BASIC_CONTENT_TYPE_FLAG | AMQP_BASIC_DELIVERY_MODE_FLAG;
	prop.content_type = amqp_cstring_bytes("application/octet-stream");
	prop.delivery_mode = 2; /* persistent delivery mode */

	/* Step-3: protobuf msg. */
    if (event == CONN_CONNECT) {
		buff = serialize_node_online_event(nodeID, nodeIP, nodePort,
                                           meshIP, meshPort, &buffLen);
    } else if (event == CONN_CLOSE) {
        buff = serialize_node_offline_event(nodeID, &buffLen);
    }

    if (buff==NULL || buffLen == 0) {
        log_error("Error serializing Link packet for AMQP. Event: %d", event);
        free(key);
        return FALSE;
    }

	/* Step-4: send the message to AMQP broker */
	amqp_bytes_t body;
	body.len = buffLen;
	body.bytes = buff;

	ret = amqp_basic_publish(conn, 1, amqp_cstring_bytes(exchange),
							 amqp_cstring_bytes(key), 0, 0, &prop,
							 body);
	if (ret < 0) {
		ret = FALSE;
		log_error("Error sending AMQP message. Error: %s",
				  amqp_error_string2(ret));
	} else {
		ret = TRUE;
		log_debug("AMQP message successfully sent to default exchange");
	}

	free(buff);
	free(key);
	return ret;
}

int publish_event(MeshEvent event, char *orgName,
                  char *nodeID, char *nodeIP, int nodePort,
                  char *meshIP, int meshPort) {

    WAMQPConn *conn=NULL;
    char *amqpHost=NULL, *amqpPort=NULL, *amqpUser=NULL, *amqpPassword=NULL;

    amqpHost = getenv(ENV_AMQP_HOST);
    amqpPort = getenv(ENV_AMQP_PORT);
    amqpUser = getenv(ENV_AMQP_USER);
    amqpPassword = getenv(ENV_AMQP_PASSWORD);
    conn = init_amqp_connection(amqpHost, amqpPort, amqpUser, amqpPassword);
    if (conn == NULL) {
        log_error("Failed to connect with AMQP at %s:%s@%s:%s", amqpUser, amqpPassword, amqpHost, amqpPort);
        return FALSE;
    }

    if (object_type(event) == OBJECT_LINK) {
        publish_amqp_event(conn, DEFAULT_MESH_AMQP_EXCHANGE, event,
                           orgName, nodeID,
                           nodeIP, nodePort,
                           meshIP, meshPort);
    } else {
        log_error("Invalid event type. No publish");
    }

    amqp_channel_close(conn, 1, AMQP_REPLY_SUCCESS);
	amqp_connection_close(conn, AMQP_REPLY_SUCCESS);
	amqp_destroy_connection(conn);

    return TRUE;
}

int publish_register_event(char *exchange, int port) {

    WAMQPConn *conn=NULL;
    char *amqpHost=NULL, *amqpPort=NULL, *amqpUser=NULL, *amqpPassword=NULL;
    char *orgName=NULL, *ip=NULL;

    char key[MAX_BUFFER]={0};
    WAMQPProp prop;
    void *buff=NULL;
    size_t buffLen=0;
    int ret;

    amqpHost     = getenv(ENV_AMQP_HOST);
    amqpPort     = getenv(ENV_AMQP_PORT);
    amqpUser     = getenv(ENV_AMQP_USER);
    amqpPassword = getenv(ENV_AMQP_PASSWORD);
    orgName      = getenv(ENV_SYSTEM_ORG);
    ip           = getenv(ENV_BINDING_IP);

    conn = init_amqp_connection(amqpHost, amqpPort, amqpUser, amqpPassword);
    if (conn == NULL) {
        log_error("Failed to connect with AMQP at %s:%s@%s:%s",
                  amqpUser, amqpPassword, amqpHost, amqpPort);
        return FALSE;
    }

    /* set routing key */
    sprintf(key, "event.cloud.global.%s.messaging.mesh.ip.update", orgName);

    /* set AMQP delivery properties */
    prop._flags = AMQP_BASIC_CONTENT_TYPE_FLAG | AMQP_BASIC_DELIVERY_MODE_FLAG;
    prop.content_type = amqp_cstring_bytes("application/octet-stream");
    prop.delivery_mode = 2; /* persistent delivery mode */

    /* protobuf msg. */
    buff = serialize_register_event(ip, port, &buffLen);
    if (buff == NULL || buffLen == 0) {
        log_error("Error serializing boot packet for AMQP");
        return FALSE;
    }

    /* send the message to AMQP broker */
	amqp_bytes_t body;
	body.len   = buffLen;
	body.bytes = buff;

    ret = amqp_basic_publish(conn, 1, amqp_cstring_bytes(exchange),
                             amqp_cstring_bytes(key), 0, 0, &prop,
                             body);
    if (ret != AMQP_STATUS_OK) {
        free(buff);
        ret = FALSE;
        log_error("Error sending AMQP boot message. Error: %s",
                  amqp_error_string2(ret));
    } else {
        ret = TRUE;
        log_debug("AMQP boot message successfully sent to default exchange");
    }

    amqp_channel_close(conn, 1, AMQP_REPLY_SUCCESS);
    amqp_connection_close(conn, AMQP_REPLY_SUCCESS);
    amqp_destroy_connection(conn);

    free(buff);

    return ret;
}
