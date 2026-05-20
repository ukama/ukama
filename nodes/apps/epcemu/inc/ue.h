/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef UE_H_
#define UE_H_

#include <stdint.h>
#include <stdbool.h>

#include "jansson.h"

#define UE_IMSI_LEN 32
#define UE_IP_LEN   64
#define UE_APN_LEN  64
#define UE_STATE_LEN 32

typedef json_t JsonObj;

typedef enum {
    UeStateAttaching = 0,
    UeStateAttached,
    UeStateDetaching,
    UeStateFailed
} UeState;

typedef struct UeEntry {
    char imsi[UE_IMSI_LEN];
    char ip[UE_IP_LEN];
    char apn[UE_APN_LEN];
    UeState state;
    uint64_t attachedAt;
    uint64_t updatedAt;
    struct UeEntry *next;
} UeEntry;

void ue_store_init(void);
void ue_store_destroy(void);
int ue_attach_start(const char *imsi, const char *ip, const char *apn,
                    char *reason, size_t reasonLen);
void ue_attach_complete(const char *imsi);
void ue_attach_fail(const char *imsi, const char *reason);
int ue_detach_start(const char *imsi, UeEntry *copy);
void ue_detach_complete(const char *imsi);
int ue_get(const char *imsi, UeEntry *copy);
JsonObj *ue_get_json(const char *imsi);
JsonObj *ue_list_json(void);
int ue_count_attached(void);
void ue_for_each_attached(int (*cb)(const UeEntry *ue, void *arg), void *arg);
const char *ue_state_str(UeState state);

#endif /* UE_H_ */
