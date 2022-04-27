/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef OBJ_MODULE_INFO_H_
#define OBJ_MODULE_INFO_H_

#include "objects/objects.h"

#define OBJECT_ID_MODULE		   34568

typedef struct __attribute__((__packed__)) {
	uint16_t    instanceId;             // matches lwm2m_list_t::id
	char        uuid[MAX_LWM2M_OBJ_STR_LEN];
	char        name[MAX_LWM2M_OBJ_STR_LEN];
	int			class;
	char        partnumber[MAX_LWM2M_OBJ_STR_LEN];
	char        mfgname[MAX_LWM2M_OBJ_STR_LEN];
	char        mfgdate[MAX_LWM2M_OBJ_STR_LEN];
	char        mac[MAX_LWM2M_OBJ_STR_LEN];
	char        sw_version[MAX_LWM2M_OBJ_STR_LEN];
	char        psw_version[MAX_LWM2M_OBJ_STR_LEN];
	char        hw_version[MAX_LWM2M_OBJ_STR_LEN];
	int 		device_count; //devices
} ModuleObjInfo;

typedef struct _module_info
{
    struct _module_info * next;         // matches lwm2m_list_t::next
    ModuleObjInfo data;
} module_info_t;

// Resource Id's:
#define RES_M_MOD_UUID                      0
#define RES_M_MOD_NAME                      1
#define RES_M_MOD_CLASS                     2
#define RES_M_PART_NUMBER               	3
#define RES_M_MOD_MFGNAME		           	4
#define RES_M_MOD_MFGDATE					5
#define RES_M_MOD_MAC                 		6
#define RES_M_MOD_SW_VERSION     		    7
#define RES_M_MOD_PSW_VERSION				8
#define RES_M_MOD_HW_VERSION               	9
#define RES_M_MOD_DEV_COUNT            		10
#endif
