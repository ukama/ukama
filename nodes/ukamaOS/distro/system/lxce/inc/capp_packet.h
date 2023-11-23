/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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

int create_capp_tx_packet(CApp *capp, PacketList **list, int reqType);


#endif /* CAPP_PACKET_H */
