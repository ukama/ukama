/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef INC_CONFIGD_H_
#define INC_CONFIGD_H_

#include "jserdes.h"
#include "session.h"

#include "usys_services.h"

#define SERVICE_NAME           SERVICE_CONFIG
#define STATUS_OK              (0)
#define STATUS_NOK             (-1)

#define MAX_APPS               32
#define MAX_PATH               512
#define MAX_URL                512
#define MAX_FILE_PATH          1024

#define DEF_LOG_LEVEL          "TRACE"
#define DEF_SPACE_NAME         "services"

#define CONFIG_VERSION         "0.0.0"

#define DEF_SERVICE_PORT       "8080"
#define DEF_NODED_HOST         "localhost"
#define DEF_NODED_PORT         "8095"
#define DEF_NODED_EP           "/v1/nodeinfo"
#define DEF_STARTER_HOST       "localhost"
#define DEF_STARTER_PORT       "8086"
#define DEF_STARTER_EP         "/v1"
#define DEF_NODE_ID            "ukama-aaa-bbbb-ccc-dddd"
#define ENV_CONFIG_DEBUG_MODE  "ENV_CONFIG_DEBUG_MODE"

#define DEF_CONFIG_DIR   "/ukama/configs"
#define CONFIG_TMP_PATH  "/tmp"


bool process_received_config(JsonObj *json, Config *config);
void free_session_data(SessionData *d);

#endif /* INC_CONFIGD_H_ */
