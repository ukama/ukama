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

#define MAX_PATH      512
#define MAX_FILE_PATH 1024


int configd_process_incoming_config(const char *service,
                                    JsonObj *json,
                                    Config *config);
int configd_process_complete(Config *config);
int configd_trigger_update(Config* c);
int read_active_config(ConfigData **c);
void free_config_data(ConfigData *c);
#endif /* INC_NOTIFICATION_H_ */
