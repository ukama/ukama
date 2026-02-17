/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef BOARD_CONFIG_H
#define BOARD_CONFIG_H

#define MAX_BOARD_APPS 128
#define MAX_APP_KEY    64

#define COMMON_CONFIG      "common.conf"
#define TRX_CONFIG         "trx.conf"
#define COM_CONFIG         "com.conf"
#define CONTROLLER_CONFIG  "controller.conf"

typedef enum {
    NODE_UNKNOWN = 0,
    NODE_TOWER,
    NODE_AMPLIFIER
} NodeType;

typedef struct {
    char key[MAX_APP_KEY];  /* e.g. METRICSD_APP */
    int  enabled;           /* 1=yes, 0=no */
} BoardAppEntry;

typedef struct {
    NodeType      type;
    BoardAppEntry entries[MAX_BOARD_APPS];
    int           count;
} BoardConfig;

/* Detect node type from nodeID */
NodeType detect_node_type(const char *nodeID);

/* Load common + board specific */
int board_config_load(BoardConfig *cfg,
                      const char *boardsDir,
                      NodeType type);

/* Query if app enabled (expects capp->name) */
int board_is_app_enabled(BoardConfig *cfg,
                         const char *appName);

#endif
