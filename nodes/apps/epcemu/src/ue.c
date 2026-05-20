/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <time.h>

#include "epcemu.h"
#include "ue.h"

static UeEntry *gHead = NULL;
static pthread_mutex_t gLock;

static uint64_t now_sec(void) {

    return (uint64_t)time(NULL);
}

const char *ue_state_str(UeState state) {

    switch (state) {
    case UeStateAttaching: return "attaching";
    case UeStateAttached:  return "attached";
    case UeStateDetaching: return "detaching";
    case UeStateFailed:    return "failed";
    default:               return "unknown";
    }
}

static UeEntry *find_by_imsi(const char *imsi) {

    UeEntry *cur;

    for (cur = gHead; cur != NULL; cur = cur->next) {
        if (!strcmp(cur->imsi, imsi)) return cur;
    }

    return NULL;
}

static UeEntry *find_by_ip(const char *ip) {

    UeEntry *cur;

    for (cur = gHead; cur != NULL; cur = cur->next) {
        if (!strcmp(cur->ip, ip)) return cur;
    }

    return NULL;
}

static void copy_entry(UeEntry *dst, const UeEntry *src) {

    if (dst == NULL || src == NULL) return;

    memset(dst, 0, sizeof(UeEntry));
    snprintf(dst->imsi, sizeof(dst->imsi), "%s", src->imsi);
    snprintf(dst->ip, sizeof(dst->ip), "%s", src->ip);
    snprintf(dst->apn, sizeof(dst->apn), "%s", src->apn);
    dst->state = src->state;
    dst->attachedAt = src->attachedAt;
    dst->updatedAt = src->updatedAt;
}

static JsonObj *entry_json(const UeEntry *ue) {

    JsonObj *obj;

    if (ue == NULL) return NULL;

    obj = json_object();
    if (obj == NULL) return NULL;

    json_object_set_new(obj, "imsi", json_string(ue->imsi));
    json_object_set_new(obj, "ip", json_string(ue->ip));
    json_object_set_new(obj, "apn", json_string(ue->apn));
    json_object_set_new(obj, "state", json_string(ue_state_str(ue->state)));
    json_object_set_new(obj, "attachedAt", json_integer(ue->attachedAt));
    json_object_set_new(obj, "updatedAt", json_integer(ue->updatedAt));

    return obj;
}

void ue_store_init(void) {

    gHead = NULL;
    pthread_mutex_init(&gLock, NULL);
}

void ue_store_destroy(void) {

    UeEntry *cur;
    UeEntry *next;

    pthread_mutex_lock(&gLock);

    cur = gHead;
    while (cur != NULL) {
        next = cur->next;
        usys_free(cur);
        cur = next;
    }
    gHead = NULL;

    pthread_mutex_unlock(&gLock);
    pthread_mutex_destroy(&gLock);
}

int ue_attach_start(const char *imsi, const char *ip, const char *apn,
                    char *reason, size_t reasonLen) {

    UeEntry *ue;
    UeEntry *existing;

    if (imsi == NULL || ip == NULL) return USYS_FALSE;

    pthread_mutex_lock(&gLock);

    existing = find_by_imsi(imsi);
    if (existing != NULL) {
        snprintf(reason, reasonLen, "imsi already exists");
        pthread_mutex_unlock(&gLock);
        return USYS_FALSE;
    }

    existing = find_by_ip(ip);
    if (existing != NULL) {
        snprintf(reason, reasonLen, "ip already allocated");
        pthread_mutex_unlock(&gLock);
        return USYS_FALSE;
    }

    ue = usys_calloc(1, sizeof(UeEntry));
    if (ue == NULL) {
        snprintf(reason, reasonLen, "memory allocation failed");
        pthread_mutex_unlock(&gLock);
        return USYS_FALSE;
    }

    snprintf(ue->imsi, sizeof(ue->imsi), "%s", imsi);
    snprintf(ue->ip, sizeof(ue->ip), "%s", ip);
    snprintf(ue->apn, sizeof(ue->apn), "%s", apn ? apn : EPCEMU_DEF_APN);
    ue->state = UeStateAttaching;
    ue->attachedAt = now_sec();
    ue->updatedAt = ue->attachedAt;
    ue->next = gHead;
    gHead = ue;

    pthread_mutex_unlock(&gLock);
    return USYS_TRUE;
}

void ue_attach_complete(const char *imsi) {

    UeEntry *ue;

    pthread_mutex_lock(&gLock);

    ue = find_by_imsi(imsi);
    if (ue != NULL) {
        ue->state = UeStateAttached;
        ue->updatedAt = now_sec();
    }

    pthread_mutex_unlock(&gLock);
}

void ue_attach_fail(const char *imsi, const char *reason) {

    UeEntry *cur;
    UeEntry *prev;

    (void)reason;

    pthread_mutex_lock(&gLock);

    cur = gHead;
    prev = NULL;
    while (cur != NULL) {
        if (!strcmp(cur->imsi, imsi)) {
            if (prev == NULL) gHead = cur->next;
            else prev->next = cur->next;
            usys_free(cur);
            break;
        }
        prev = cur;
        cur = cur->next;
    }

    pthread_mutex_unlock(&gLock);
}

int ue_detach_start(const char *imsi, UeEntry *copy) {

    UeEntry *ue;

    if (imsi == NULL) return USYS_FALSE;

    pthread_mutex_lock(&gLock);

    ue = find_by_imsi(imsi);
    if (ue == NULL) {
        pthread_mutex_unlock(&gLock);
        return USYS_FALSE;
    }

    ue->state = UeStateDetaching;
    ue->updatedAt = now_sec();
    copy_entry(copy, ue);

    pthread_mutex_unlock(&gLock);
    return USYS_TRUE;
}

void ue_detach_complete(const char *imsi) {

    UeEntry *cur;
    UeEntry *prev;

    pthread_mutex_lock(&gLock);

    cur = gHead;
    prev = NULL;
    while (cur != NULL) {
        if (!strcmp(cur->imsi, imsi)) {
            if (prev == NULL) gHead = cur->next;
            else prev->next = cur->next;
            usys_free(cur);
            break;
        }
        prev = cur;
        cur = cur->next;
    }

    pthread_mutex_unlock(&gLock);
}

int ue_get(const char *imsi, UeEntry *copy) {

    UeEntry *ue;

    if (imsi == NULL || copy == NULL) return USYS_FALSE;

    pthread_mutex_lock(&gLock);

    ue = find_by_imsi(imsi);
    if (ue == NULL) {
        pthread_mutex_unlock(&gLock);
        return USYS_FALSE;
    }

    copy_entry(copy, ue);

    pthread_mutex_unlock(&gLock);
    return USYS_TRUE;
}

JsonObj *ue_get_json(const char *imsi) {

    UeEntry ue;

    if (!ue_get(imsi, &ue)) return NULL;

    return entry_json(&ue);
}

JsonObj *ue_list_json(void) {

    JsonObj *arr;
    UeEntry *cur;

    arr = json_array();
    if (arr == NULL) return NULL;

    pthread_mutex_lock(&gLock);

    for (cur = gHead; cur != NULL; cur = cur->next) {
        json_array_append_new(arr, entry_json(cur));
    }

    pthread_mutex_unlock(&gLock);

    return arr;
}

int ue_count_attached(void) {

    int count;
    UeEntry *cur;

    count = 0;
    pthread_mutex_lock(&gLock);

    for (cur = gHead; cur != NULL; cur = cur->next) {
        if (cur->state == UeStateAttached) count++;
    }

    pthread_mutex_unlock(&gLock);
    return count;
}

void ue_for_each_attached(int (*cb)(const UeEntry *ue, void *arg),
                          void *arg) {

    UeEntry *cur;
    UeEntry copy;

    if (cb == NULL) return;

    pthread_mutex_lock(&gLock);

    cur = gHead;
    while (cur != NULL) {
        if (cur->state == UeStateAttached) {
            copy_entry(&copy, cur);
            pthread_mutex_unlock(&gLock);
            cb(&copy, arg);
            pthread_mutex_lock(&gLock);
            cur = gHead;
            continue;
        }
        cur = cur->next;
    }

    pthread_mutex_unlock(&gLock);
}
