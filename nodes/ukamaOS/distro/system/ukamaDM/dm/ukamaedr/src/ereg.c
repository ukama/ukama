/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "inc/ereg.h"

#include "client/lwm2mclient.h"
#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "headers/edr/ifmsg.h"
#include "inc/msghandler.h"
#include "objects/objects.h"
#include "objects/alarm.h"

MsgFrame* ereg_send(MsgFrame* smsg, size_t * size, int* sflag) {
	MsgFrame* rmsg = NULL;
	int ret = 0;
	if (smsg) {
		rmsg = msghandler_client_send(smsg, size, sflag);
		if (!rmsg) {
			fprintf(stderr, "Err: EREG:: Failure in receiving message back from registry.");
		}
	}
	return rmsg;
}

int ereg_handledata(MsgFrame* rmsg, MsgFrame* smsg, void** data, int sflag) {
	int ret = 0;
	if (rmsg && smsg) {
		if (sflag) {
			/* Failure in processing request */
			fprintf(stderr, "Err(%d): EREG:: Failure while reading Inst: %d of Type 0x%x RId: %d",
					ret, rmsg->instance, rmsg->objecttype, rmsg->resourceId);
			UKAMA_FREE(*data);
		} else {
			if (rmsg) {
				if (!(rmsg->response)){
					/* No error */
					/* Copy response */
					if(rmsg->data && data) {
						if (*data) {
							memcpy(*data, rmsg->data, rmsg->datasize);
						}
					}
				} else {
					/* Failure in UkamaEdge registry */
					ret = rmsg->response;
					fprintf(stderr, "Err(%d): EREG:: Failure reading value in Ukama Edge Registry: %d of Type 0x%x RId: %d\n\r",
							ret, rmsg->instance, rmsg->objecttype, rmsg->resourceId);
				}
			}
		}
	}
	return ret;
}

int ereg_read_inst(uint16_t inst, uint16_t stype, uint16_t rid, void* data, size_t* size) {
	int ret = 0;
	MsgFrame* rmsg = NULL;
	MsgFrame* smsg = create_msgframe(MSG_TYPE_READ_REQ, inst, rid, stype, *size, true, data);
	if (smsg) {
		/* Send Msg */
		rmsg = ereg_send(smsg, size, &ret);
		/* Process Msg Response */
		if ( rmsg && (rmsg->msgtype == MSG_TYPE_READ_RESP) ) {
			ret = ereg_handledata(rmsg, smsg, &data, ret);
		} else {
			ret = ERR_UNEXPECTED_RESP_MSG ; /* ERR_UNEXPECTED_RESP_MSG */
		}
	} else {
		ret = ERR_UBSP_INVALID_POINTER;
	}
	free_msgframe(&smsg);
	free_msgframe(&rmsg);
	return ret;
}

int ereg_read_inst_count(uint16_t stype, void* data, size_t* size) {
	int ret = 0;
	MsgFrame* rmsg = NULL;
	MsgFrame* smsg = create_msgframe(MSG_TYPE_READ_INST_COUNT, 0, 0, stype, *size, true, data);
	if (smsg) {
		/* Send Msg */
		rmsg = ereg_send(smsg, size, &ret);

		/* Process Msg Response */
		if ( rmsg && (rmsg->msgtype == MSG_TYPE_READ_INST_COUNT_RESP) ) {
			ret = ereg_handledata(rmsg, smsg, &data, ret);
		} else {
			ret = ERR_UNEXPECTED_RESP_MSG ; /* ERR_UNEXPECTED_RESP_MSG */
		}

	} else {
		ret = ERR_UBSP_INVALID_POINTER;
	}
	free_msgframe(&smsg);
	free_msgframe(&rmsg);
	return ret;
}

int ereg_write_inst(uint16_t inst, uint16_t stype, uint16_t rid, void* data, size_t* size) {
	int ret = 0;
	MsgFrame* rmsg = NULL;
	MsgFrame* smsg = create_msgframe(MSG_TYPE_WRITE_REQ, inst, rid, stype, *size, true, data);
	if (smsg) {
		/* Send Msg */
		rmsg = ereg_send(smsg, size, &ret);
		/* Process Msg Response */
		if ( rmsg && (rmsg->msgtype == MSG_TYPE_WRITE_RESP) ) {
			ret = ereg_handledata(rmsg, smsg, &data, ret);
		} else {
			ret = ERR_UNEXPECTED_RESP_MSG ; /* ERR_UNEXPECTED_RESP_MSG */
		}

	} else {
		ret = ERR_UBSP_INVALID_POINTER;
	}

	free_msgframe(&smsg);
	free_msgframe(&rmsg);
	return ret;
}

int ereg_exec_sensor(uint16_t inst, uint16_t stype, uint16_t rid, void* data, size_t* size) {
	int ret = 0;
	MsgFrame* rmsg = NULL;
	MsgFrame* smsg = create_msgframe(MSG_TYPE_EXEC_REQ, inst, rid, stype, *size, true, data);
	if (smsg) {
		/* Send Msg */
		rmsg = ereg_send(smsg, size, &ret);
		/* Process Msg Response */
		if ( rmsg && (rmsg->msgtype == MSG_TYPE_EXEC_RESP) ) {
			ret = ereg_handledata(rmsg, smsg, &data, ret);
		} else {
			ret = ERR_UNEXPECTED_RESP_MSG ; /* ERR_UNEXPECTED_RESP_MSG */
		}
	} else {
		ret = ERR_UBSP_INVALID_POINTER;
	}
	free_msgframe(&smsg);
	free_msgframe(&rmsg);
	return ret;
}

int ereg_handle_alarm(void* ctx, void *objdata, size_t *size) {
	int ret = 0;
	AlarmObjInfo* obj = objdata;
	/* Create URI.*/
	uint16_t objid = OBJECT_ID_DEV_ALARM;
	uint16_t instid = obj->instanceId;
	uint16_t rsrcid = RES_M_AL_STATE;
	fprintf(stdout,
					"EREG:: Handling value %d change for Alarm /%d/%d/%d Sensor /%d/%d/%d.\r\n",
					obj->state, objid, instid, rsrcid, obj->sobjid, obj->sinstid, obj->srsrcid) ;
	handle_alarm_update(ctx, &obj->state, objid, instid, rsrcid, sizeof(uint16_t));
	return ret;
}
