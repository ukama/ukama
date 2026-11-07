/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef NETUTIL_H_
#define NETUTIL_H_

#include <stdint.h>
#include <stdbool.h>

#include "epcemu.h"

typedef json_t JsonObj;

typedef struct {
    uint32_t address;
    uint32_t mask;
} Ipv4Cidr;

bool ipv4_cidr_parse(const char *text, Ipv4Cidr *cidr);
bool ue_subnet_parse(const char *text, Ipv4Cidr *pool);
int apn_normalize(const char *input, size_t len, char *output);
int identity_create(char *id);
int identity_valid(const char *id);
int imsi_valid(const char *imsi);
int ip_to_uint32(const char *ip, uint32_t *out);
int ue_address_valid(const EpcemuConfig *config, const char *ip);
int ip_in_cidr(const char *ip, const char *cidr);
JsonObj *imsi_to_json_array(const char *imsi);

#endif /* NETUTIL_H_ */

