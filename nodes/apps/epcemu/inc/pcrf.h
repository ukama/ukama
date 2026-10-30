/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef PCRF_H_
#define PCRF_H_

#include "epcemu.h"
#include "ue.h"

typedef struct {
    char instanceId[EPCEMU_ID_LEN];
    char state[32];
} PcrfEpc;

typedef struct {
    char state[32];
    char dataState[16];
    char reason[EPCEMU_MAX_REASON];
} PcrfPdn;

/* Return the exact HTTP status, or zero for an uncertain/invalid response. */
int pcrf_get_epc(EpcemuConfig *config, PcrfEpc *epc);
int pcrf_register_epc(EpcemuConfig *config, const char *id,
                      const char *expected);
int pcrf_create_session(EpcemuConfig *config, const UeEntry *ue,
                        PcrfPdn *pdn);
int pcrf_get_pdn(EpcemuConfig *config, const UeEntry *ue, PcrfPdn *pdn);
int pcrf_delete_session(EpcemuConfig *config, const UeEntry *ue);

#endif /* PCRF_H_ */
