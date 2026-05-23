/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef DATA_PLANE_H_
#define DATA_PLANE_H_

#include <stdbool.h>
#include <stdint.h>

#include "epcemu.h"

typedef struct {
    bool running;
    bool ready;
    int udpFd;
    int tunFd;
    pthread_t thread;

    uint64_t uplinkPackets;
    uint64_t uplinkBytes;
    uint64_t downlinkPackets;
    uint64_t downlinkBytes;
    uint64_t droppedPackets;
    uint64_t droppedBytes;

    pthread_mutex_t lock;
} DataPlane;

int data_plane_start(DataPlane *dp,
                     EpcemuConfig *config,
                     EpcemuStatus *status);

void data_plane_stop(DataPlane *dp);

JsonObj *data_plane_json(DataPlane *dp, EpcemuConfig *config);

#endif /* DATA_PLANE_H_ */
