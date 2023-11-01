/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef INC_NODED_MACROS_H_
#define INC_NODED_MACROS_H_

#define MAX_NAME_LENGTH                    (24)
#define MAX_PATH_LENGTH                    (64)

/* Length */
#define UUID_LENGTH                         (32)
#define NAME_LENGTH                         MAX_NAME_LENGTH
#define PATH_LENGTH                         MAX_PATH_LENGTH
#define DATE_LENGTH                         (12)
#define MAC_LENGTH                          (18)
#define DESC_LENGTH                         (24)

#define MAX_JSON_DEVICE                     (32)
#define PROP_MAX_STR_LENGTH                 (64)

#define STATUS_OK                           (0)
#define STATUS_NOK                          (-1)

/* Limits */
#define MAX_NUMBER_MODULES_PER_UNIT         0x0007
#define MAX_NUMBER_DEVICES_PER_MODULE       0x0010

/* MAX_PAYLOAD_SIZE */
#define SCH_MAX_PAYLOAD_SIZE                0x1000

/*
 * Field Id These are used to identify configuration in the index tables
 */
#define FIELD_ID_NODE_INFO                  0x0001
#define FIELD_ID_NODE_CFG                   0x0002
#define FIELD_ID_MODULE_INFO                0x0003
#define FIELD_ID_MODULE_CFG                 0x0004
#define FIELD_ID_FACT_CFG                   0x0005
#define FIELD_ID_USER_CFG                   0x0006
#define FIELD_ID_FACT_CALIB                 0x0007
#define FIELD_ID_USER_CALIB                 0x0008
#define FIELD_ID_BS_CERTS                   0x0009
#define FIELD_ID_CLOUD_CERTS                0x000a

#endif /* INC_NODED_MACROS_H_ */
