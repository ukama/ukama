/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * capp_packet.h
 */

#ifndef CAPP_PACKET_H
#define CAPP_PACKET_H

#include <pthread.h>
#include <uuid/uuid.h>

/* requests */
#define CAPP_TYPE_REQ_CREATE 0x01
#define CAPP_TYPE_REQ_RUN    0x02
#define CAPP_TYPE_REQ_STATUS 0x03
#define CAPP_TYPE_REQ_TERM   0x04

/* responses */
#define CAPP_TYPE_RESP_UUID   0x10
#define CAPP_TYPE_RESP_STATUS 0x11

#define CAPP_TYPE_NONE        0xff

/* For capp state */
#define CAPP_STATE_PENDING 0x01
#define CAPP_STATE_CREATE  0x02
#define CAPP_STATE_RUN     0x03
#define CAPP_STATE_TERM    0x04
#define CAPP_STATE_INVALID 0xff

typedef struct capp_packet_t {

  int    reqType;   /* create, run, status, stop, term */
  int    respType;  /* uuid, state */

  /* Trio for capp creation. */
  char   *name;      /* Name of the capp */
  char   *tag;       /* tag for the capp, e.g., latest. */
  char   *path;      /* complete path of the contained app, rootfs */

  uuid_t uuid;        /* UUID associated with the capp */

  int    state;       /* Current state of the capp */
  int    exitStatus; /* If capp terminate, its exit status */
}CAppPacket;

typedef struct packet_list_ {

  CAppPacket *packet;

  struct packet_list *next;
}PacketList;

#endif /* CAPP_PACKET_H */
