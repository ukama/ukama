/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SWITCHD_SNMP_CLIENT_H
#define SWITCHD_SNMP_CLIENT_H

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

typedef enum {
    SNMP_VALUE_NONE = 0,
    SNMP_VALUE_INT,
    SNMP_VALUE_UINT,
    SNMP_VALUE_STRING,
    SNMP_VALUE_OID
} SnmpValueType;

typedef struct {
    SnmpValueType type;
    int64_t intValue;
    uint64_t uintValue;
    char stringValue[256];
    uint32_t oid[64];
    size_t oidLen;
} SnmpValue;

typedef struct {
    char host[64];
    int port;
    char community[64];
    int timeoutMs;
    int retries;
} SnmpSession;

typedef struct {
    uint32_t oid[64];
    size_t oidLen;
    SnmpValue value;
} SnmpVarBind;

int snmp_session_init(SnmpSession *s,
                      const char *host,
                      int port,
                      const char *community,
                      int timeoutMs,
                      int retries);
int snmp_get(SnmpSession *s,
             const uint32_t *oid,
             size_t oidLen,
             SnmpVarBind *out);
int snmp_set_integer(SnmpSession *s,
                     const uint32_t *oid,
                     size_t oidLen,
                     int32_t value);
int snmp_set_string(SnmpSession *s,
                    const uint32_t *oid,
                    size_t oidLen,
                    const char *value);
int snmp_walk(SnmpSession *s,
              const uint32_t *baseOid,
              size_t baseLen,
              int (*cb)(const SnmpVarBind *vb, void *arg),
              void *arg);
int snmp_get_next(SnmpSession *s,
                  const uint32_t *oid,
                  size_t oidLen,
                  SnmpVarBind *out);
int snmp_oid_from_string(const char *s, uint32_t *oid, size_t *oidLen);
int snmp_oid_has_prefix(const uint32_t *oid,
                        size_t oidLen,
                        const uint32_t *prefix,
                        size_t prefixLen);
char *snmp_oid_to_string(const uint32_t *oid,
                         size_t oidLen,
                         char *buf,
                         size_t bufLen);

#endif
