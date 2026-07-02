/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef RETAP_CODES_H_
#define RETAP_CODES_H_

/*
 * 3GPP TS 25.463 RETAP elementary procedure codes.
 * Ukama supports common procedures plus single-antenna RET procedures only.
 */
#define RETAP_PROC_RESET_SOFTWARE          0x03
#define RETAP_PROC_GET_ERROR_STATUS        0x04
#define RETAP_PROC_GET_ALARM_STATUS        RETAP_PROC_GET_ERROR_STATUS /* legacy API name */
#define RETAP_PROC_GET_INFORMATION         0x05
#define RETAP_PROC_CLEAR_ACTIVE_ALARMS     0x06
#define RETAP_PROC_ALARM                   0x07  /* secondary -> primary, class 2 */
#define RETAP_PROC_SELF_TEST               0x0A
#define RETAP_PROC_SET_DEVICE_DATA         0x0E
#define RETAP_PROC_GET_DEVICE_DATA         0x0F
#define RETAP_PROC_READ_USER_DATA          0x10
#define RETAP_PROC_WRITE_USER_DATA         0x11
#define RETAP_PROC_ALARM_SUBSCRIBE         0x12

#define RETAP_PROC_CALIBRATE               0x31
#define RETAP_PROC_SEND_CONFIG_DATA        0x32
#define RETAP_PROC_SET_TILT                0x33
#define RETAP_PROC_GET_TILT                0x34
#define RETAP_PROC_SINGLE_RET_ALARM        0x35

/*
 * RETAP response data starts with OK or FAIL for single-antenna devices.
 * Return codes are from TS 25.463 Annex A.
 */
#define RETAP_RETURN_OK                    0x00
#define RETAP_RETURN_FAIL                  0x0B

#define RETAP_RC_OK                        0x00
#define RETAP_RC_ACTUATOR_DETECTION_FAIL   0x01
#define RETAP_RC_ACTUATOR_JAM_PERMANENT    0x02
#define RETAP_RC_ACTUATOR_JAM_TEMPORARY    0x03
#define RETAP_RC_BLOCK_NUMBER_SEQUENCE     0x04
#define RETAP_RC_BUSY                      0x05
#define RETAP_RC_CHECKSUM_ERROR            0x06
#define RETAP_RC_PROCEDURE_SEQUENCE_ERROR  0x07
#define RETAP_RC_DATA_ERROR                0x08
#define RETAP_RC_DEVICE_DISABLED           0x09
#define RETAP_RC_EEPROM_ERROR              0x0A
#define RETAP_RC_FAIL                      0x0B
#define RETAP_RC_FLASH_ERASE_ERROR         0x0C
#define RETAP_RC_FLASH_ERROR               0x0D
#define RETAP_RC_NOT_CALIBRATED            0x0E
#define RETAP_RC_NOT_SCALED                0x0F
#define RETAP_RC_OTHER_HARDWARE_ERROR      0x11
#define RETAP_RC_OTHER_SOFTWARE_ERROR      0x12
#define RETAP_RC_OUT_OF_RANGE              0x13
#define RETAP_RC_POSITION_LOST             0x14
#define RETAP_RC_RAM_ERROR                 0x15
#define RETAP_RC_SEGMENT_NUMBER_SEQUENCE   0x16
#define RETAP_RC_UART_ERROR                0x17
#define RETAP_RC_UNKNOWN_PROCEDURE         0x19
#define RETAP_RC_READ_ONLY                 0x1D
#define RETAP_RC_UNKNOWN_PARAMETER         0x1E
#define RETAP_RC_UNKNOWN_ANTENNA_NUMBER    0x1F
#define RETAP_RC_TOO_MUCH_DATA             0x20
#define RETAP_RC_WORKING_SOFTWARE_MISSING  0x21
#define RETAP_RC_INVALID_SOFTWARE          0x22
#define RETAP_RC_DOWNLOAD_IN_PROGRESS      0x23

#define RETAP_HEADER_LEN                   3
#define RETAP_MAX_PAYLOAD                  2045
#define RETAP_MAX_ENCODED                  (RETAP_HEADER_LEN + RETAP_MAX_PAYLOAD)
#define RETAP_MAX_ALARMS                   32

#define RETAP_DEFAULT_TIMEOUT_MS           3000
#define RETAP_SET_TILT_TIMEOUT_MS          120000
#define RETAP_CALIBRATE_TIMEOUT_MS         240000
#define RETAP_CONFIG_TIMEOUT_MS            3000
#define RETAP_CONFIG_SEGMENT_MAX           70

#endif /* RETAP_CODES_H_ */
