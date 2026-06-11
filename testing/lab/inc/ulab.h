/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_H_
#define ULAB_H_

#include <stdint.h>
#include <stddef.h>
#include <stdbool.h>

#include "version.h"

#define ULAB_VERSION       VERSION
#define ULAB_SCHEMA_VER    1
#define ULAB_MAX_NAME      256
#define ULAB_MAX_ID        512
#define ULAB_MAX_REF       128
#define ULAB_MAX_PATH      1024
#define ULAB_MAX_ERR       1024
#define ULAB_MAX_LINE      512

#define ULAB_OK            0
#define ULAB_ERR           1
#define ULAB_EUSAGE        2
#define ULAB_ESCENARIO     3
#define ULAB_EBFF          4
#define ULAB_ERUNTIME      5
#define ULAB_EINTERNAL     6

typedef enum {
    ULAB_FALSE = 0,
    ULAB_TRUE  = 1
} ulab_bool_t;

typedef struct {
    char msg[ULAB_MAX_ERR];
} ulab_error_t;

static inline void ulab_error_clear(ulab_error_t *err) {
    if (err != NULL) {
        err->msg[0] = '\0';
    }
}

#endif /* ULAB_H_ */
