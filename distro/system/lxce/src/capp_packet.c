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

#include "capp_packet.h"
#include "log.h"

#define TRUE 1
#define FALSE 0
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
