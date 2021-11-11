/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* Functions related to capp requests from parent process to the threads. */

#include <stdlib.h>
#include <uuid/uuid.h>

#include "capp.h"
#include "capp_packet.h"
#include "log.h"

#define TRUE 1
#define FALSE 0

/*
 * init_capp_packet_list --
 *
 */
int init_capp_packet_list(PacketList **list) {

  (*list) = (PacketList *)malloc(sizeof(PacketList));
  if (*list == NULL) {
    log_error("Memory allocation error. Size: %s", sizeof(PacketList));
    return FALSE;
  }

  (*list)->packet = NULL;
  (*list)->next   = NULL;

  return TRUE;
}

/*
 * init_capp_packet --
 *
 */
int init_capp_packet(CAppPacket *packet) {

  if (packet==NULL) return FALSE;

  packet->name = packet->tag = packet->path = NULL;

  packet->req_type  = CAPP_TYPE_NONE;
  packet->resp_type = CAPP_TYPE_NONE;
  packet->state     = CAPP_STATE_INVALID;

  uuid_clear(packet->uuid);

  return TRUE;
}

/*
 * create_capp_tx_packet --
 *
 */
int create_capp_tx_packet(int reqType, CApp *capp, PacketList **list) {

  PacketList *ptr = list;
  CAppPacket *packet=NULL;

  if (capp==NULL || packet==NULL) return;

  while (ptr->next) {
    ptr=ptr->next;
  }

  ptr = (PacketList *)malloc(sizeof(PacketList));
  if (!ptr) {
    log_error("Memory allocation error. size: %d", sizeof(PacketList));
    return FALSE;
  }

  packet = (CAppPacket *)malloc(sizeof(CAppPacket));
  if (!packet) {
    log_error("Memory allocation error. size: %d", sizeof(CAppPacket));
    free(ptr);
    return FALSE;
  }

  packet->req_type = reqType;

  packet->name = strdup(capp->params->name);
  packet->tag  = strdup(capp->params->tag);
  packet->path = strdup(capp->params->path);

  uuid_clear(packet->uuid );

  /* Add to the list. */
  ptr->packet = packet;
  ptr->next   = NULL;

  return TRUE;
}
