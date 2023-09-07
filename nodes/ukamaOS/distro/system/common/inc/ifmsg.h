/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_IFMSG_H_
#define INC_IFMSG_H_

#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>

#define MISC_TYPE_UNIT					0x0001
#define MISC_TYPE_MODULE				0x0002

#define MSG_RESP_SUCCESS			0
#define MSG_RESP_FAILURE			(-2001)
#define MSG_RESP_INVALID_MSGTYPE	(-2002)
#define MSG_RESP_INVALID_DEVTYPE	(-2003)
#define MSG_RESP_INVALID_INSTAID	(-2004)
#define MSG_RESP_INAVLID_RESRCID	(-2005)
#define MSG_RESP_INAVLID_REQTOK		(-2006)
#define MSG_RESP_INAVLID_REQMSG		(-2007)
#define MSG_RESP_INSTANCE_MISSING	(-2008)


typedef enum {
	/* Resource */
	MSG_TYPE_READ_REQ = 1,
	MSG_TYPE_READ_RESP,
	MSG_TYPE_WRITE_REQ,
	MSG_TYPE_WRITE_RESP,
	MSG_TYPE_EXEC_REQ,
	MSG_TYPE_EXEC_RESP,
	MSG_TYPE_ALERT_REP,
	MSG_TYPE_ALERT_RESP,
	/* Instance Count */
	MSG_TYPE_READ_INST_COUNT,
	MSG_TYPE_READ_INST_COUNT_RESP
} MsgType;

typedef struct __attribute__((__packed__)) {
	MsgType msgtype;
	int reqtoken;
	uint16_t misc; 		 // For future use
	uint16_t objecttype; // Sensor and Module specific request.
	uint16_t instance;
	uint16_t resourceId;
	int16_t response;
	uint16_t datasize;
	void *data;
} MsgFrame;

int msgframe_validate(MsgFrame* rmsg, MsgFrame* smsg);
void free_msgframe(MsgFrame ** msg);
char* msgframe_serialize(MsgFrame* msg, size_t* size);

MsgFrame* msgframe_deserialize(char* msg);
MsgFrame *create_msgframe(MsgType msgtype, uint16_t inst, uint16_t rid,
                          uint16_t objtype, size_t size, bool token, void *data);

#endif /* INC_IFMSG_H_ */
