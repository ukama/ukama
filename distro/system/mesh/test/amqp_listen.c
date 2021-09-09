/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* AMQP listening service. */

#include <stdio.h>
#include <string.h>
#include <amqp_socket.h>
#include <amqp.h>

#include "link.pb-c.h"

#define MAX_FRAME 131072
#define MAX_MSG_SIZE 1024

#define TRUE 1
#define FALSE 0

#define USER     "guest"
#define PASSWORD "guest"

typedef amqp_basic_properties_t WAMQPProp;
typedef struct amqp_connection_state_t_ WAMQPConn;
typedef struct amqp_socket_t_           WAMQPSocket;
typedef struct amqp_rpc_reply_t_        WAMQPReply;

static WAMQPConn *amqp_connection(char *host, int port) {

  int ret;
  WAMQPConn *conn=NULL;
  WAMQPSocket *socket=NULL;
  WAMQPReply reply;

  /* Sanity check */
  if (host == NULL && port <= 0) {
    return NULL;
  }

  /* Initialize AMQP state variable */
  conn = amqp_new_connection();
  if (!conn) {
    fprintf(stderr, "Unable to create AMQP, memory allocation issue\n");
    return NULL;
  }

  /* Create TCP socket */
  socket = amqp_tcp_socket_new(conn);
  if (!socket) {
    fprintf(stderr, "Unable to create AMQP TCP socket\n");
    amqp_destroy_connection(conn);
    return NULL;
  }

  /* Connect to the AMQP host */
  ret = amqp_socket_open(socket, host, port);
  if (ret) {
    fprintf(stderr, "Unable to connect with AMQP at host: %s port: %d\n",
	    host, port);
    amqp_destroy_connection(conn);
    return NULL;
  }

  /* Login using guest/guest (for now - XXX) */
  reply = amqp_login(conn, "/", 0, MAX_FRAME, 0, AMQP_SASL_METHOD_PLAIN,
		     USER, PASSWORD);
  if (reply.reply_type != AMQP_RESPONSE_NORMAL) {
    fprintf(stderr, "AMQP login failed!\n");
    amqp_channel_close(conn, 1, AMQP_CHANNEL_ERROR);
    amqp_destroy_connection(conn);
    return NULL;
  }

  /* Open the channel */
  amqp_channel_open(conn, 1);
  reply = amqp_get_rpc_reply(conn);
  if (reply.reply_type != AMQP_RESPONSE_NORMAL) {
    fprintf(stderr, "AMQP channel open error\n");
    amqp_channel_close(conn, 1, AMQP_CHANNEL_ERROR);
    amqp_destroy_connection(conn);
    return NULL;
  }

  return conn;
}

int main(int argc, char **argv) {

  int port;
  char *host, *exchange, *bindKey;
  WAMQPConn *conn=NULL;
  WAMQPReply reply;
  amqp_bytes_t queue;
  Link *msg;
  uint8_t buff[MAX_MSG_SIZE];

  if (argc<5) {
    fprintf(stderr, "USAGE: %s host port exchange key\n", argv[0]);
    return 1;
  }

  host     = argv[1];
  port     = atoi(argv[2]);
  exchange = argv[3];
  bindKey  = argv[4];

  /* Open connection with AMQP host */
  conn = amqp_connection(host, port);
  if (!conn) {
    fprintf(stderr, "Error connecting with AMQP at host: %s port %d\n", host,
	    port);
    return 1;
  }

  /* Setup Queue */
  amqp_queue_bind(conn, 1, queue, amqp_cstring_bytes(exchange),
                  amqp_cstring_bytes(bindKey), amqp_empty_table);
  reply = amqp_get_rpc_reply(conn);
  if (reply.reply_type != AMQP_RESPONSE_NORMAL) {
    fprintf(stderr, "Binding queue error.\n");
    return 1;
  }

  amqp_basic_consume(conn, 1, queue, amqp_empty_bytes, 0, 1, 0,
		     amqp_empty_table);
  reply = amqp_get_rpc_reply(conn);
  if (reply.reply_type != AMQP_RESPONSE_NORMAL) {
    fprintf(stderr, "AMQP consuming error.\n");
    return 1;
  }

  /* Loop forever to consume. */
  for (;;) {
    
    amqp_rpc_reply_t resp;
    amqp_envelope_t  envelope;

    amqp_maybe_release_buffers(conn);

    resp = amqp_consume_message(conn, &envelope, NULL, 0);

    if (resp.reply_type != AMQP_RESPONSE_NORMAL) {
      break;
    }

    fprintf(stdout, "Delivery %u, exchange %.*s routingkey %.*s\n",
	    (unsigned)envelope.delivery_tag, (int)envelope.exchange.len,
	    (char *)envelope.exchange.bytes, (int)envelope.routing_key.len,
	    (char *)envelope.routing_key.bytes);

    if (envelope.message.properties._flags & AMQP_BASIC_CONTENT_TYPE_FLAG) {
      fprintf(stdout, "Content-type: %.*s\n",
	      (int)envelope.message.properties.content_type.len,
	      (char *)envelope.message.properties.content_type.bytes);
    }

    fprintf(stdout, "----\n");

    fprintf(stdout, "Len: %ld Msg: %s\n", envelope.message.body.len,
	    (char *)envelope.message.body.bytes);

    if (envelope.message.body.len) {
      /* Try to unpack the msg using protobuf */
      memset(buff, 0, MAX_MSG_SIZE);
      msg = link__unpack(NULL, envelope.message.body.len, buff);

      if (msg == NULL) {
	fprintf(stderr, "Error unpacking incoming message\n");
	return 1;
      }

      fprintf(stdout, "Key: %s Msg recevied: %s \n",
	      (char *)envelope.routing_key.bytes, msg->uuid);
    }

    link__free_unpacked(msg, NULL);
    amqp_destroy_envelope(&envelope);
  }

  amqp_bytes_free(queue);
  amqp_channel_close(conn, 1, AMQP_REPLY_SUCCESS);
  amqp_connection_close(conn, AMQP_REPLY_SUCCESS);
  amqp_destroy_connection(conn);

  fprintf(stdout, "Ending listening session ... Connection terminated.\n");
  return 0;
}
