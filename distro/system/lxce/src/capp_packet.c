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
#include <string.h>

#include "capp.h"
#include "capp_packet.h"
#include "log.h"

#define TRUE 1
#define FALSE 0

static void init_capp_packet(CAppPacket *packet);

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
static void init_capp_packet(CAppPacket *packet) {

  if (packet==NULL) return;

  packet->name = packet->tag = packet->path = NULL;

  packet->reqType  = CAPP_TYPE_NONE;
  packet->respType = CAPP_TYPE_NONE;
  packet->state    = CAPP_STATE_INVALID;

  uuid_clear(packet->uuid);
}

/*
 * create_capp_tx_packet --
 *
 */
int create_capp_tx_packet(CApp *capp, PacketList **list, int reqType) {

  PacketList *ptr = *list;
  CAppPacket *packet=NULL;

  if (capp==NULL || list==NULL) return FALSE;

  if (*list == NULL) { /* First entry */
    if (!init_capp_packet_list(list))
      return FALSE;
  }

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

  init_capp_packet(packet);
  packet->reqType = reqType;

  packet->name = strdup(capp->params->name);
  packet->tag  = strdup(capp->params->tag);
  packet->path = strdup(capp->params->path);

  uuid_clear(packet->uuid );

  /* Add to the list. */
  ptr->packet = packet;
  ptr->next   = NULL;

  return TRUE;
}
