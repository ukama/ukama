/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef RETAP_CODES_H_
#define RETAP_CODES_H_

#define RETAP_PROC_RESET_SOFTWARE          0x03
#define RETAP_PROC_GET_ALARM_STATUS        0x04
#define RETAP_PROC_GET_INFORMATION         0x05
#define RETAP_PROC_CLEAR_ACTIVE_ALARMS     0x06
#define RETAP_PROC_SELF_TEST               0x0A
#define RETAP_PROC_GET_DEVICE_DATA         0x0F
#define RETAP_PROC_ALARM_SUBSCRIBE         0x12
#define RETAP_PROC_CALIBRATE               0x31
#define RETAP_PROC_SEND_CONFIG_DATA        0x32
#define RETAP_PROC_SET_TILT                0x33
#define RETAP_PROC_GET_TILT                0x34

#define RETAP_RETURN_OK                    0x00
#define RETAP_RETURN_FAIL                  0x01

#define RETAP_MAX_PAYLOAD                  2048
#define RETAP_MAX_ALARMS                   32

#endif /* RETAP_CODES_H_ */
