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
#include <limits.h>
#include <net/if.h>
#include <sys/ioctl.h>
#include <unistd.h>
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

/* Accept a dotted IPv4 address followed by one decimal prefix, 0..32. */
bool ipv4_cidr_parse(const char *text, Ipv4Cidr *cidr) {

    char host[INET_ADDRSTRLEN];
    struct in_addr addr;
    const char *slash;
    const char *p;
    size_t len;
    unsigned int prefix = 0;

    if (text == NULL || cidr == NULL) return false;
    slash = strchr(text, '/');
    if (slash == NULL) return false;
    len = (size_t)(slash - text);
    if (len == 0 || len >= sizeof(host)) return false;
    memcpy(host, text, len);
    host[len] = '\0';
    if (inet_pton(AF_INET, host, &addr) != 1) return false;

    p = slash + 1;
    if (*p == '\0' || (*p == '0' && p[1] != '\0')) return false;
    for (; *p; p++) {
        if (*p < '0' || *p > '9') return false;
        prefix = prefix * 10 + (unsigned int)(*p - '0');
        if (prefix > 32) return false;
    }

    cidr->address = ntohl(addr.s_addr);
    cidr->mask = prefix == 0 ? 0 : UINT32_MAX << (32 - prefix);
    return true;
}

/* A UE pool needs a network, a broadcast and at least one host address. */
bool ue_subnet_parse(const char *text, Ipv4Cidr *pool) {

    return ipv4_cidr_parse(text, pool) &&
           (pool->address & ~pool->mask) == 0 && ~pool->mask > 1u;
}

int ip_in_cidr(const char *ip, const char *cidr) {

    uint32_t value;
    Ipv4Cidr subnet;

    return ip_to_uint32(ip, &value) && ipv4_cidr_parse(cidr, &subnet) &&
           (value & subnet.mask) == (subnet.address & subnet.mask);
}

/* SIOCGIFCONF includes secondary IPv4 addresses. Fail closed on an
 * incomplete snapshot; interface configuration can change between calls. */
static bool address_is_available(uint32_t value) {

    struct ifconf list = {0};
    struct ifreq *entries;
    const struct sockaddr_in *addr;
    size_t count;
    size_t i;
    int size;
    int fd;
    bool available = false;

    fd = socket(AF_INET, SOCK_DGRAM | SOCK_CLOEXEC, 0);
    if (fd < 0) return false;
    if (ioctl(fd, SIOCGIFCONF, &list) < 0 || list.ifc_len < 0 ||
        list.ifc_len > INT_MAX - (int)sizeof(*entries)) goto done;
    size = list.ifc_len + (int)sizeof(*entries);
    entries = malloc((size_t)size);
    if (entries == NULL) goto done;
    list.ifc_req = entries;
    list.ifc_len = size;
    if (ioctl(fd, SIOCGIFCONF, &list) == 0 && list.ifc_len >= 0 &&
        list.ifc_len <= size - (int)sizeof(*entries)) {
        available = true;
        count = (size_t)list.ifc_len / sizeof(*entries);
        for (i = 0; i < count; i++) {
            addr = (const struct sockaddr_in *)&entries[i].ifr_addr;
            if (addr->sin_family == AF_INET && ntohl(addr->sin_addr.s_addr) == value) {
                available = false;
                break;
            }
        }
    }
    free(entries);
done:
    close(fd);
    return available;
}

int ue_address_valid(const EpcemuConfig *config, const char *ip) {

    Ipv4Cidr pool;
    Ipv4Cidr iface;
    uint32_t value;

    if (config == NULL || !ue_subnet_parse(config->ueCidr, &pool) ||
        !ip_to_uint32(ip, &value) || (value & pool.mask) != pool.address ||
        value == pool.address || value == (pool.address | ~pool.mask) ||
        (value >> 24) == 0 || (value >> 24) == 127 ||
        (value >> 28) >= 14 || (value & 0xffff0000u) == 0xa9fe0000u) {
        return USYS_FALSE;
    }
    if ((ipv4_cidr_parse(config->tunAddr, &iface) && value == iface.address) ||
        (ipv4_cidr_parse(config->bridgeCidr, &iface) && value == iface.address) ||
        (ipv4_cidr_parse(config->gatewayCidr, &iface) && value == iface.address)) {
        return USYS_FALSE;
    }

    return address_is_available(value);
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
