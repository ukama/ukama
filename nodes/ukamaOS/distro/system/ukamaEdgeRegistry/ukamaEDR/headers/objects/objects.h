/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#ifndef OBJ_OBJECTS_H_
#define OBJ_OBJECTS_H_

#include <stdbool.h>
#include <stdint.h>

/* OBJECT TYPE */
#define OBJ_TYPE_NULL                   0x0000
#define OBJ_TYPE_UNIT                   0x0001
#define OBJ_TYPE_MOD                   	0x0002
#define OBJ_TYPE_TMP					0x0003
#define OBJ_TYPE_VOLT					0x0004
#define OBJ_TYPE_CURR					0x0005
#define OBJ_TYPE_PWR					0x0006
#define OBJ_TYPE_DIP					0x0007
#define OBJ_TYPE_DOP					0x0008
#define OBJ_TYPE_LED					0x0009
#define OBJ_TYPE_ADC					0x000A
#define OBJ_TYPE_ATT					0x000B
#define OBJ_TYPE_ALARM					0x000C
#define OBJ_TYPE_MAX					0x000D

typedef uint16_t ObjectType;

#define ALL_RESOURCE_ID					0xFFFF

#define MAX_LWM2M_OBJ_STR_LEN			256
#endif /* OBJ_OBJECTS_H_ */
