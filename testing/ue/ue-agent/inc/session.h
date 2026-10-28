/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef UE_SESSION_H
#define UE_SESSION_H

#include <stdbool.h>

#ifndef UE_STATUS_DIR
#define UE_STATUS_DIR "/run/ukama-ue"
#endif
#define UE_STATUS_FILE UE_STATUS_DIR "/status.json"

typedef struct {
    const char *url;
    const char *imsi;
    const char *iccid;
    const char *ip;
    const char *requestedApn;
    const char *dataHost;
    int dataPort;
    char instanceId[37];
    char pdnId[37];
    char apn[100];
    char state[16];
    char dataState[16];
    char reason[128];
    long httpStatus;
    bool probe;
    bool uncertain;
    bool terminal;
} UeSession;

int session_init(UeSession *session);
int session_attach(UeSession *session);
int session_poll(UeSession *session);
int session_detach(UeSession *session);

#endif
