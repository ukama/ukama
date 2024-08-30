/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef CONFIG_MACROS_H_
#define CONFIG_MACROS_H_

#define SERVICE_NAME           "configd"
#define STATUS_OK              (0)
#define STATUS_NOK             (-1)

#define MAX_SERVICE_COUNT      (32)

#define DEF_LOG_LEVEL          "TRACE"

#define CONFIG_VERSION         "0.0.0"

#define DEF_SERVICE_PORT       "8080"
#define DEF_NODED_HOST         "localhost"
#define DEF_NODED_PORT         "8095"
#define DEF_NODED_EP           "/v1"
#define DEF_STARTER_HOST       "localhost"
#define DEF_STARTER_PORT       "8086"
#define DEF_STARTER_EP         "/v1"
#define DEF_NODE_ID            "ukama-aaa-bbbb-ccc-dddd"
#define ENV_CONFIG_DEBUG_MODE  "ENV_CONFIG_DEBUG_MODE"

#define CONFIG_TMP_PATH "/tmp"

#define CONFIG_STORE_PATH "/etc/config"
#define CONFIG_RUNNING "/etc/config/running"
#define CONFIG_BACKUP "/etc/config/backup"
#define CONFIG_OLD "/etc/config/old"

#define CONFIGD  "/configd/version.json"

#endif /* INC_CONFIG_MACROS_H_ */
