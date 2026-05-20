/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <arpa/inet.h>
#include <ctype.h>
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
    if (len < 5 || len >= UE_IMSI_LEN) return USYS_FALSE;

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
