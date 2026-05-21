/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <arpa/inet.h>
#include <stdio.h>
#include <string.h>
#include <time.h>

#include "ue.h"
#include "usys_mem.h"

static UeEntry gUes[EPCEMU_MAX_UES];
static pthread_mutex_t gLock;

static uint64_t now_sec(void) {
    return (uint64_t)time(NULL);
}

const char *ue_state_str(UeState state) {

    switch (state) {
    case UeStateFree:      return "free";
    case UeStateAttaching: return "attaching";
    case UeStateAttached:  return "attached";
    case UeStateDetaching: return "detaching";
    case UeStateFailed:    return "failed";
    default:               return "unknown";
    }
}

void ue_store_init(void) {

    memset(gUes, 0, sizeof(gUes));
    pthread_mutex_init(&gLock, NULL);
}

void ue_store_destroy(void) {

    pthread_mutex_destroy(&gLock);
}

static int find_idx_imsi_locked(const char *imsi) {

    int i;

    if (imsi == NULL) return -1;

    for (i = 0; i < EPCEMU_MAX_UES; i++) {
        if (gUes[i].state != UeStateFree &&
            strcmp(gUes[i].imsi, imsi) == 0) {
            return i;
        }
    }

    return -1;
}

static int find_idx_ip_locked(const char *ip) {

    int i;

    if (ip == NULL) return -1;

    for (i = 0; i < EPCEMU_MAX_UES; i++) {
        if (gUes[i].state != UeStateFree &&
            strcmp(gUes[i].ip, ip) == 0) {
            return i;
        }
    }

    return -1;
}

static int free_idx_locked(void) {

    int i;

    for (i = 0; i < EPCEMU_MAX_UES; i++) {
        if (gUes[i].state == UeStateFree) {
            return i;
        }
    }

    return -1;
}

static void copy_entry(UeEntry *dst, const UeEntry *src) {

    if (dst == NULL || src == NULL) return;

    memcpy(dst, src, sizeof(UeEntry));
}

int ue_attach_start(const char *imsi,
                    const char *iccid,
                    const char *ip,
                    const char *apn,
                    char *reason,
                    size_t reasonLen) {

    int idx;
    int existing;
    uint64_t now;

    if (imsi == NULL || ip == NULL) {
        snprintf(reason, reasonLen, "missing imsi or ip");
        return USYS_FALSE;
    }

    pthread_mutex_lock(&gLock);

    existing = find_idx_imsi_locked(imsi);
    if (existing >= 0) {
        if (strcmp(gUes[existing].ip, ip) == 0 &&
            gUes[existing].state == UeStateAttached) {
            pthread_mutex_unlock(&gLock);
            return USYS_TRUE;
        }

        snprintf(reason, reasonLen, "imsi already attached");
        pthread_mutex_unlock(&gLock);
        return USYS_FALSE;
    }

    existing = find_idx_ip_locked(ip);
    if (existing >= 0) {
        snprintf(reason, reasonLen, "ue ip already allocated");
        pthread_mutex_unlock(&gLock);
        return USYS_FALSE;
    }

    idx = free_idx_locked();
    if (idx < 0) {
        snprintf(reason, reasonLen, "max ue limit reached");
        pthread_mutex_unlock(&gLock);
        return USYS_FALSE;
    }

    now = now_sec();

    memset(&gUes[idx], 0, sizeof(UeEntry));
    snprintf(gUes[idx].imsi, sizeof(gUes[idx].imsi), "%s", imsi);
    snprintf(gUes[idx].iccid, sizeof(gUes[idx].iccid), "%s",
             iccid ? iccid : "");
    snprintf(gUes[idx].ip, sizeof(gUes[idx].ip), "%s", ip);
    snprintf(gUes[idx].apn, sizeof(gUes[idx].apn), "%s",
             apn ? apn : EPCEMU_DEF_APN);

    gUes[idx].state = UeStateAttaching;
    gUes[idx].attachedAt = now;
    gUes[idx].updatedAt = now;

    pthread_mutex_unlock(&gLock);
    return USYS_TRUE;
}

void ue_attach_complete(const char *imsi) {

    int idx;

    pthread_mutex_lock(&gLock);

    idx = find_idx_imsi_locked(imsi);
    if (idx >= 0) {
        gUes[idx].state = UeStateAttached;
        gUes[idx].updatedAt = now_sec();
    }

    pthread_mutex_unlock(&gLock);
}

void ue_attach_fail(const char *imsi, const char *reason) {

    int idx;

    (void)reason;

    pthread_mutex_lock(&gLock);

    idx = find_idx_imsi_locked(imsi);
    if (idx >= 0) {
        memset(&gUes[idx], 0, sizeof(UeEntry));
        gUes[idx].state = UeStateFree;
    }

    pthread_mutex_unlock(&gLock);
}

int ue_detach_start(const char *imsi, UeEntry *copy) {

    int idx;

    pthread_mutex_lock(&gLock);

    idx = find_idx_imsi_locked(imsi);
    if (idx < 0) {
        pthread_mutex_unlock(&gLock);
        return USYS_FALSE;
    }

    gUes[idx].state = UeStateDetaching;
    gUes[idx].updatedAt = now_sec();
    copy_entry(copy, &gUes[idx]);

    pthread_mutex_unlock(&gLock);
    return USYS_TRUE;
}

void ue_detach_complete(const char *imsi) {

    int idx;

    pthread_mutex_lock(&gLock);

    idx = find_idx_imsi_locked(imsi);
    if (idx >= 0) {
        memset(&gUes[idx], 0, sizeof(UeEntry));
        gUes[idx].state = UeStateFree;
    }

    pthread_mutex_unlock(&gLock);
}

int ue_get(const char *imsi, UeEntry *copy) {

    int idx;

    if (copy == NULL) return USYS_FALSE;

    pthread_mutex_lock(&gLock);

    idx = find_idx_imsi_locked(imsi);
    if (idx < 0) {
        pthread_mutex_unlock(&gLock);
        return USYS_FALSE;
    }

    copy_entry(copy, &gUes[idx]);

    pthread_mutex_unlock(&gLock);
    return USYS_TRUE;
}

int ue_bind_peer(const char *ip,
                 const struct sockaddr_storage *addr,
                 socklen_t addrLen) {

    int idx;

    if (ip == NULL || addr == NULL) return USYS_FALSE;

    pthread_mutex_lock(&gLock);

    idx = find_idx_ip_locked(ip);
    if (idx < 0 || gUes[idx].state != UeStateAttached) {
        pthread_mutex_unlock(&gLock);
        return USYS_FALSE;
    }

    memcpy(&gUes[idx].peerAddr, addr, addrLen);
    gUes[idx].peerLen = addrLen;
    gUes[idx].peerSet = true;
    gUes[idx].lastSeenAt = now_sec();

    pthread_mutex_unlock(&gLock);
    return USYS_TRUE;
}

int ue_find_by_src_ip(const char *ip, UeEntry *copy) {
    return ue_find_by_dst_ip(ip, copy);
}

int ue_find_by_dst_ip(const char *ip, UeEntry *copy) {

    int idx;

    if (ip == NULL || copy == NULL) return USYS_FALSE;

    pthread_mutex_lock(&gLock);

    idx = find_idx_ip_locked(ip);
    if (idx < 0 || gUes[idx].state != UeStateAttached) {
        pthread_mutex_unlock(&gLock);
        return USYS_FALSE;
    }

    copy_entry(copy, &gUes[idx]);

    pthread_mutex_unlock(&gLock);
    return USYS_TRUE;
}

void ue_record_uplink(const char *ip, uint64_t bytes) {

    int idx;

    pthread_mutex_lock(&gLock);

    idx = find_idx_ip_locked(ip);
    if (idx >= 0) {
        gUes[idx].uplinkPackets++;
        gUes[idx].uplinkBytes += bytes;
        gUes[idx].lastSeenAt = now_sec();
    }

    pthread_mutex_unlock(&gLock);
}

void ue_record_downlink(const char *ip, uint64_t bytes) {

    int idx;

    pthread_mutex_lock(&gLock);

    idx = find_idx_ip_locked(ip);
    if (idx >= 0) {
        gUes[idx].downlinkPackets++;
        gUes[idx].downlinkBytes += bytes;
        gUes[idx].lastSeenAt = now_sec();
    }

    pthread_mutex_unlock(&gLock);
}

static JsonObj *entry_json_locked(const UeEntry *ue) {

    JsonObj *obj;
    JsonObj *peer;
    char addr[128];
    uint16_t port;

    if (ue == NULL) return NULL;

    obj = json_object();
    if (obj == NULL) return NULL;

    json_object_set_new(obj, "imsi",       json_string(ue->imsi));
    json_object_set_new(obj, "iccid",      json_string(ue->iccid));
    json_object_set_new(obj, "ip",         json_string(ue->ip));
    json_object_set_new(obj, "apn",        json_string(ue->apn));
    json_object_set_new(obj, "state",      json_string(ue_state_str(ue->state)));
    json_object_set_new(obj, "attachedAt", json_integer(ue->attachedAt));
    json_object_set_new(obj, "updatedAt",  json_integer(ue->updatedAt));
    json_object_set_new(obj, "lastSeenAt", json_integer(ue->lastSeenAt));

    json_object_set_new(obj, "uplinkPackets",
                        json_integer(ue->uplinkPackets));
    json_object_set_new(obj, "uplinkBytes",
                        json_integer(ue->uplinkBytes));
    json_object_set_new(obj, "downlinkPackets",
                        json_integer(ue->downlinkPackets));
    json_object_set_new(obj, "downlinkBytes",
                        json_integer(ue->downlinkBytes));

    peer = json_object();
    json_object_set_new(peer, "bound", json_boolean(ue->peerSet));

    if (ue->peerSet && ue->peerAddr.ss_family == AF_INET) {
        struct sockaddr_in *sin;

        sin = (struct sockaddr_in *)&ue->peerAddr;
        inet_ntop(AF_INET, &sin->sin_addr, addr, sizeof(addr));
        port = ntohs(sin->sin_port);

        json_object_set_new(peer, "address", json_string(addr));
        json_object_set_new(peer, "port", json_integer(port));
    }

    json_object_set_new(obj, "peer", peer);

    return obj;
}

JsonObj *ue_get_json(const char *imsi) {

    int idx;
    JsonObj *obj;

    obj = NULL;

    pthread_mutex_lock(&gLock);

    idx = find_idx_imsi_locked(imsi);
    if (idx >= 0) {
        obj = entry_json_locked(&gUes[idx]);
    }

    pthread_mutex_unlock(&gLock);
    return obj;
}

JsonObj *ue_list_json(void) {

    JsonObj *arr;
    int i;

    arr = json_array();
    if (arr == NULL) return NULL;

    pthread_mutex_lock(&gLock);

    for (i = 0; i < EPCEMU_MAX_UES; i++) {
        if (gUes[i].state != UeStateFree) {
            json_array_append_new(arr, entry_json_locked(&gUes[i]));
        }
    }

    pthread_mutex_unlock(&gLock);
    return arr;
}

JsonObj *ue_summary_json(void) {

    JsonObj *obj;
    int i;
    int attached;
    int used;
    int failed;

    attached = 0;
    used = 0;
    failed = 0;

    obj = json_object();
    if (obj == NULL) return NULL;

    pthread_mutex_lock(&gLock);

    for (i = 0; i < EPCEMU_MAX_UES; i++) {
        if (gUes[i].state != UeStateFree) {
            used++;
        }
        if (gUes[i].state == UeStateAttached) {
            attached++;
        }
        if (gUes[i].state == UeStateFailed) {
            failed++;
        }
    }

    pthread_mutex_unlock(&gLock);

    json_object_set_new(obj, "max", json_integer(EPCEMU_MAX_UES));
    json_object_set_new(obj, "used", json_integer(used));
    json_object_set_new(obj, "attached", json_integer(attached));
    json_object_set_new(obj, "failed", json_integer(failed));

    return obj;
}

int ue_count_attached(void) {

    int i;
    int count;

    count = 0;

    pthread_mutex_lock(&gLock);

    for (i = 0; i < EPCEMU_MAX_UES; i++) {
        if (gUes[i].state == UeStateAttached) count++;
    }

    pthread_mutex_unlock(&gLock);
    return count;
}

int ue_count_total(void) {

    int i;
    int count;

    count = 0;

    pthread_mutex_lock(&gLock);

    for (i = 0; i < EPCEMU_MAX_UES; i++) {
        if (gUes[i].state != UeStateFree) count++;
    }

    pthread_mutex_unlock(&gLock);
    return count;
}

void ue_for_each_attached(int (*cb)(const UeEntry *ue, void *arg),
                          void *arg) {

    UeEntry copy[EPCEMU_MAX_UES];
    int i;
    int count;

    if (cb == NULL) return;

    count = 0;

    pthread_mutex_lock(&gLock);

    for (i = 0; i < EPCEMU_MAX_UES; i++) {
        if (gUes[i].state == UeStateAttached && count < EPCEMU_MAX_UES) {
            copy_entry(&copy[count], &gUes[i]);
            count++;
        }
    }

    pthread_mutex_unlock(&gLock);

    for (i = 0; i < count; i++) {
        cb(&copy[i], arg);
    }
}
