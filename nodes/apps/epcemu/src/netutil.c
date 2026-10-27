/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <arpa/inet.h>
#include <ctype.h>
#include <errno.h>
#include <sys/random.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include "ue.h"
#include "epcemu.h"
#include "netutil.h"

int imsi_valid(const char *imsi) {

    int i;
    int len;

    if (imsi == NULL) return USYS_FALSE;

    len = strlen(imsi);
    if (len < 5 || len > 15) return USYS_FALSE;

    for (i = 0; i < len; i++) {
        if (!isdigit((unsigned char)imsi[i])) return USYS_FALSE;
    }

    return USYS_TRUE;
}

int ip_to_uint32(const char *ip, uint32_t *out) {

    unsigned char buf[4];

    if (ip == NULL || out == NULL) return USYS_FALSE;

    if (inet_pton(AF_INET, ip, buf) != 1) return USYS_FALSE;

    *out = ((uint32_t)buf[0] << 24) |
           ((uint32_t)buf[1] << 16) |
           ((uint32_t)buf[2] << 8)  |
           ((uint32_t)buf[3]);

    return USYS_TRUE;
}

static int cidr_to_netmask(int prefix, uint32_t *mask) {

    if (mask == NULL || prefix < 0 || prefix > 32) return USYS_FALSE;

    if (prefix == 0) {
        *mask = 0;
    } else {
        *mask = 0xffffffffu << (32 - prefix);
    }

    return USYS_TRUE;
}

int ip_in_cidr(const char *ip, const char *cidr) {

    char cidrCopy[EPCEMU_MAX_STR];
    char *slash;
    int prefix;
    uint32_t ipValue;
    uint32_t netValue;
    uint32_t mask;

    if (ip == NULL || cidr == NULL) return USYS_FALSE;

    snprintf(cidrCopy, sizeof(cidrCopy), "%s", cidr);

    slash = strchr(cidrCopy, '/');
    if (slash == NULL) return USYS_FALSE;

    *slash = '\0';
    prefix = atoi(slash + 1);

    if (!ip_to_uint32(ip, &ipValue)) return USYS_FALSE;
    if (!ip_to_uint32(cidrCopy, &netValue)) return USYS_FALSE;
    if (!cidr_to_netmask(prefix, &mask)) return USYS_FALSE;

    return ((ipValue & mask) == (netValue & mask)) ? USYS_TRUE : USYS_FALSE;
}

JsonObj *imsi_to_json_array(const char *imsi) {

    JsonObj *arr;
    int i;
    int len;

    if (!imsi_valid(imsi)) return NULL;

    arr = json_array();
    if (arr == NULL) return NULL;

    len = strlen(imsi);
    for (i = 0; i < len; i++) {
        if (json_array_append_new(arr, json_integer(imsi[i] - '0')) != 0) {
            json_decref(arr);
            return NULL;
        }
    }

    return arr;
}

int apn_normalize(const char *input, size_t len, char *output) {

    size_t i;
    size_t label;
    unsigned char ch;
    bool alnum;

    if (input == NULL || output == NULL ||
        len == 0 || len >= EPCEMU_APN_LEN) return USYS_FALSE;

    label = 0;
    for (i = 0; i < len; i++) {
        ch = (unsigned char)input[i];
        if (ch >= 'A' && ch <= 'Z') ch += 'a' - 'A';
        alnum = (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9');
        if (ch == '.') {
            if (label == 0 || input[i - 1] == '-') return USYS_FALSE;
            label = 0;
        } else {
            if (!alnum && ch != '-') return USYS_FALSE;
            if ((label == 0 && !alnum) || ++label > 63) return USYS_FALSE;
        }
        output[i] = (char)ch;
    }
    if (label == 0 || input[len - 1] == '-') return USYS_FALSE;
    output[len] = '\0';
    return USYS_TRUE;
}

int identity_create(char *id) {

    unsigned char bytes[16];
    size_t offset;
    ssize_t count;

    if (id == NULL) return USYS_FALSE;

    offset = 0;
    while (offset < sizeof(bytes)) {
        count = getrandom(bytes + offset, sizeof(bytes) - offset, 0);
        if (count < 0 && errno == EINTR) continue;
        if (count <= 0) return USYS_FALSE;
        offset += (size_t)count;
    }
    bytes[6] = (bytes[6] & 0x0f) | 0x40;
    bytes[8] = (bytes[8] & 0x3f) | 0x80;
    snprintf(id, EPCEMU_ID_LEN,
             "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-"
             "%02x%02x%02x%02x%02x%02x",
             bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5],
             bytes[6], bytes[7], bytes[8], bytes[9], bytes[10], bytes[11],
             bytes[12], bytes[13], bytes[14], bytes[15]);
    return USYS_TRUE;
}

int identity_valid(const char *id) {

    size_t i;
    bool nonzero;
    char ch;

    if (id == NULL || strlen(id) != EPCEMU_ID_LEN - 1) return USYS_FALSE;
    nonzero = false;
    for (i = 0; i < EPCEMU_ID_LEN - 1; i++) {
        ch = id[i];
        if (i == 8 || i == 13 || i == 18 || i == 23) {
            if (ch != '-') return USYS_FALSE;
        } else {
            if (!((ch >= '0' && ch <= '9') ||
                  (ch >= 'a' && ch <= 'f'))) return USYS_FALSE;
            if (ch != '0') nonzero = true;
        }
    }
    return nonzero ? USYS_TRUE : USYS_FALSE;
}
