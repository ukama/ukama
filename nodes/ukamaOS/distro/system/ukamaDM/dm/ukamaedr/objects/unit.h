/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef OBJ_UNIT_INFO_H_
#define OBJ_UNIT_INFO_H_

#include "objects/objects.h"

#define OBJECT_ID_UNIT					34567

typedef struct __attribute__((__packed__)){
	uint16_t    instanceId;             // matches lwm2m_list_t::id
	char        uuid[MAX_LWM2M_OBJ_STR_LEN];
	char        name[MAX_LWM2M_OBJ_STR_LEN];
	int			class;
	char        skew[MAX_LWM2M_OBJ_STR_LEN];
	char        oemname[MAX_LWM2M_OBJ_STR_LEN];
	char        asmdate[MAX_LWM2M_OBJ_STR_LEN];
	char        mac[MAX_LWM2M_OBJ_STR_LEN];
	char        sw_version[MAX_LWM2M_OBJ_STR_LEN];
	char        psw_version[MAX_LWM2M_OBJ_STR_LEN];
	int 		module_count; //devices
}UnitObjInfo;

typedef struct _unit_info
{
    struct _unit_info * next;         // matches lwm2m_list_t::next
    UnitObjInfo data;
} unit_info_t;

// Resource Id's:
#define RES_M_UNIT_UUID                      0
#define RES_M_UNIT_NAME                      1
#define RES_M_UNIT_CLASS                     2
#define RES_M_SKEW                	 		 3
#define RES_M_UNIT_OEMNAME		             4
#define RES_M_UNIT_ASMDATE					 5
#define RES_M_UNIT_MAC                 		 6
#define RES_M_UNIT_SW_VERSION     		     7
#define RES_M_UNIT_PSW_VERSION				 8
#define RES_M_UNIT_MOD_COUNT            	 9

#endif

