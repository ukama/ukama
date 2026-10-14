/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc. */
#ifndef ULAB_RACK_CLIENT_H_
#define ULAB_RACK_CLIENT_H_

#include <jansson.h>
#include "hardware.h"

typedef struct rack_client rack_client_t;

/* One session per scenario. The client keeps its lease alive during BFF waits. */
int rack_client_open(rack_client_t **out, const hardware_config_t *config, const char *run_dir,
                     ulab_error_t *err);
int rack_client_call(rack_client_t *client, const char *command, char request_id[65],
                     json_t **result, ulab_error_t *err);
int rack_client_healthy(rack_client_t *client, ulab_error_t *err);
int rack_client_close(rack_client_t **client, ulab_error_t *err);

#endif
