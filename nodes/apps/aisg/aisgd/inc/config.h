/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef CONFIG_H_
#define CONFIG_H_

#include <stdbool.h>

#include "usys_services.h"

#define AISGD_SERVICE_NAME     SERVICE_AISG
#define AISGD_APP_NAME         SERVICE_AISG
#define AISGD_MAX_STR          256
#define DEF_CONFIG_FILE        "/ukama/configs/aisgd/config.toml"
#define DEF_LOG_LEVEL          "TRACE"
#define DEF_CTRL_SOCKET        "/var/run/aisg-ctrl.sock"
#define DEF_CTRL_TIMEOUT_MS    3000
#define DEF_STATE_FILE         "/ukama/apps/data/aisgd/state.json"
#define EP_BS                  "/"
#define REST_API_VERSION       "v1"
#define URL_PREFIX             EP_BS REST_API_VERSION
#define API_RES_EP(RES)        EP_BS RES

typedef struct {
    char *serviceName;
    int  servicePort;
    char *controllerPath;
    int  controllerTimeoutMs;
    bool requireConfigBeforeCalibrate;
    bool requireCalibrateBeforeSetTilt;
    char *stateFile;
} Config;

void config_set_defaults(Config *config);
bool config_load_from_file(Config *config, const char *path);
void config_free(Config *config);

#endif /* CONFIG_H_ */
