/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef UE_H_
#define UE_H_

#include <netinet/in.h>
#include <stdbool.h>
#include <stdint.h>
#include <sys/socket.h>
#include <time.h>

#include "jansson.h"
#include "epcemu.h"

typedef json_t JsonObj;

#define UE_IMSI_LEN   32
#define UE_ICCID_LEN  32
#define UE_IP_LEN     64
#define UE_APN_LEN    64
#define UE_STATE_LEN  32

typedef enum {
    UeStateFree = 0,
    UeStateAttaching,
    UeStateAttached,
    UeStateDetaching,
    UeStateFailed
} UeState;

typedef struct {
    char imsi[UE_IMSI_LEN];
    char iccid[UE_ICCID_LEN];
    char ip[UE_IP_LEN];
    char apn[UE_APN_LEN];

    UeState state;

    struct sockaddr_storage peerAddr;
    socklen_t peerLen;
    bool peerSet;

    uint64_t attachedAt;
    uint64_t updatedAt;
    uint64_t lastSeenAt;

    uint64_t uplinkPackets;
    uint64_t uplinkBytes;
    uint64_t downlinkPackets;
    uint64_t downlinkBytes;
} UeEntry;

void ue_store_init(void);
void ue_store_destroy(void);

int ue_attach_start(const char *imsi,
                    const char *iccid,
                    const char *ip,
                    const char *apn,
                    char *reason,
                    size_t reasonLen);

void ue_attach_complete(const char *imsi);
void ue_attach_fail(const char *imsi, const char *reason);

int ue_detach_start(const char *imsi, UeEntry *copy);
void ue_detach_complete(const char *imsi);

int ue_get(const char *imsi, UeEntry *copy);
int ue_bind_peer(const char *ip,
                 const struct sockaddr_storage *addr,
                 socklen_t addrLen);

int ue_find_by_src_ip(const char *ip, UeEntry *copy);
int ue_find_by_dst_ip(const char *ip, UeEntry *copy);

void ue_record_uplink(const char *ip, uint64_t bytes);
void ue_record_downlink(const char *ip, uint64_t bytes);

JsonObj *ue_get_json(const char *imsi);
JsonObj *ue_list_json(void);
JsonObj *ue_summary_json(void);

int ue_count_attached(void);
int ue_count_total(void);

void ue_for_each_attached(int (*cb)(const UeEntry *ue, void *arg),
                          void *arg);

const char *ue_state_str(UeState state);

#endif /* UE_H_ */
