/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef BUILDER_H_
#define BUILDER_H_

#define STATUS_OK              (0)
#define STATUS_NOK             (-1)

#define DEF_LOG_LEVEL    "TRACE"
#define BUILDER_VERSION  "0.0.1"

#define DEF_CONFIG_FILE  "ukama.json"

#define MAX_CONFIG_FILE_SIZE 4096
#define MAX_LINE_LENGTH      256
#define MAX_VARIABLES        64
#define MAX_BUFFER           1024

#define DELIMINATOR ","
#define UKAMA_AUTH  "ukama-auth"
#define DEF_NODE_ID "uk-sa0001-tnode-a1-1234"

#endif /* BUILDER_H_ */
