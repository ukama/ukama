/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <stdint.h>
#include <unistd.h>

#include "driver.h"
#include "snmp_client.h"
#include "utils.h"
#include "log.h"

#include "usys_log.h"

/* Standard IF-MIB/MIB-II */
static const uint32_t OID_ifNumber[] = {
    1, 3, 6, 1, 2, 1, 2, 1, 0
};

static const uint32_t OID_ifDescr[] = {
    1, 3, 6, 1, 2, 1, 2, 2, 1, 2
};

static const uint32_t OID_ifSpeed[] = {
    1, 3, 6, 1, 2, 1, 2, 2, 1, 5
};

static const uint32_t OID_ifAdminStatus[] = {
    1, 3, 6, 1, 2, 1, 2, 2, 1, 7
};

static const uint32_t OID_ifOperStatus[] = {
    1, 3, 6, 1, 2, 1, 2, 2, 1, 8
};

static const uint32_t OID_ifInOctets[] = {
    1, 3, 6, 1, 2, 1, 2, 2, 1, 10
};

static const uint32_t OID_ifInUcast[] = {
    1, 3, 6, 1, 2, 1, 2, 2, 1, 11
};

static const uint32_t OID_ifInDiscards[] = {
    1, 3, 6, 1, 2, 1, 2, 2, 1, 13
};

static const uint32_t OID_ifInErrors[] = {
    1, 3, 6, 1, 2, 1, 2, 2, 1, 14
};

static const uint32_t OID_ifOutOctets[] = {
    1, 3, 6, 1, 2, 1, 2, 2, 1, 16
};

static const uint32_t OID_ifOutUcast[] = {
    1, 3, 6, 1, 2, 1, 2, 2, 1, 17
};

static const uint32_t OID_ifOutDiscards[] = {
    1, 3, 6, 1, 2, 1, 2, 2, 1, 19
};

static const uint32_t OID_ifOutErrors[] = {
    1, 3, 6, 1, 2, 1, 2, 2, 1, 20
};

static const uint32_t OID_ifName[] = {
    1, 3, 6, 1, 2, 1, 31, 1, 1, 1, 1
};

static const uint32_t OID_ifHCInOctets[] = {
    1, 3, 6, 1, 2, 1, 31, 1, 1, 1, 6
};

static const uint32_t OID_ifHCOutOctets[] = {
    1, 3, 6, 1, 2, 1, 31, 1, 1, 1, 10
};

/* Tycon/IMI enterprise MIB: 1.3.6.1.4.1.12284.5 */
static const uint32_t OID_serialNumber[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 1, 7, 0
};

static const uint32_t OID_manufactureName[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 1, 8, 0
};

static const uint32_t OID_hardwareVersion[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 1, 9, 0
};

static const uint32_t OID_softwareVersion[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 1, 10, 0
};

static const uint32_t OID_commonLoadTftpAddress[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 1, 11, 0
};

static const uint32_t OID_commonLoadTftpFileName[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 1, 12, 0
};

static const uint32_t OID_commonLoadType[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 1, 13, 0
};

static const uint32_t OID_commonLoadExecute[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 1, 14, 0
};

static const uint32_t OID_commonLoadExecuteStatus[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 1, 15, 0
};

static const uint32_t OID_sysSaveToNvm[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 1, 1, 0
};

static const uint32_t OID_sysReset[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 1, 2, 0
};

static const uint32_t OID_poeExist[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 2, 1, 1, 2
};

static const uint32_t OID_poeAdmin[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 2, 1, 1, 3
};

static const uint32_t OID_poeOperStatus[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 2, 1, 1, 4
};

static const uint32_t OID_poePower[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 2, 1, 1, 5
};

static const uint32_t OID_poeCurrent[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 2, 1, 1, 6
};

static const uint32_t OID_poeVoltage[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 2, 1, 1, 7
};

static const uint32_t OID_poeClass[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 2, 1, 1, 8
};

static const uint32_t OID_poeTotalPowerConsumption[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 2, 2, 0
};

static const uint32_t OID_poeTotalMaxPower[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 2, 3, 0
};

static const uint32_t OID_industrySystemTemperature[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 6, 3, 0
};

static const uint32_t OID_industryAmbientTemperature[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 6, 6, 0
};

static const uint32_t OID_industryPowerIn[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 6, 13, 0
};

static const uint32_t OID_industrySystemCurrent[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 6, 14, 0
};

static const uint32_t OID_industrySystemPower[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 6, 17, 0
};

static const uint32_t OID_industryOutAlarmPortLinkFail[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 6, 28, 0
};

static const uint32_t OID_industryOutAlarmPortPoeFail[] = {
    1, 3, 6, 1, 4, 1, 12284, 5, 6, 29, 0
};

typedef struct {
    SnmpSession session;
} TyconPriv;

typedef struct {
    SwitchPortState *ports;
    uint32_t max;
} WalkCtx;

static int get_string(SnmpSession *s,
                      const uint32_t *oid,
                      size_t oidLen,
                      char *dst,
                      size_t dstLen) {

    SnmpVarBind vb;
    int ret;
    size_t n;

    ret = snmp_get(s, oid, oidLen, &vb);
    if (ret != SWITCHD_OK) {
        return ret;
    }

    if (vb.value.type != SNMP_VALUE_STRING) {
        return SWITCHD_ERR_PROTOCOL;
    }

    if (dstLen == 0) {
        return SWITCHD_ERR_INVAL;
    }

    n = strlen(vb.value.stringValue);
    if (n >= dstLen) {
        n = dstLen - 1;
    }

    memcpy(dst, vb.value.stringValue, n);
    dst[n] = '\0';

    return SWITCHD_OK;
}

static int get_int(SnmpSession *s,
                   const uint32_t *oid,
                   size_t oidLen,
                   int64_t *value) {

    SnmpVarBind vb;
    int ret;

    ret = snmp_get(s, oid, oidLen, &vb);
    if (ret != SWITCHD_OK) {
        return ret;
    }

    if (vb.value.type == SNMP_VALUE_INT) {
        *value = vb.value.intValue;
    } else if (vb.value.type == SNMP_VALUE_UINT) {
        *value = (int64_t)vb.value.uintValue;
    } else {
        return SWITCHD_ERR_PROTOCOL;
    }

    return SWITCHD_OK;
}

static SwitchPortState *ensure_port(WalkCtx *w, uint32_t id) {

    uint32_t i;
    SwitchPortState *port;

    if (id == 0 || id > w->max) {
        return NULL;
    }

    for (i = 0; i < w->max; i++) {
        port = &w->ports[i];

        if (port->id == id) {
            return port;
        }

        if (port->id == 0) {
            port->id = id;

            snprintf(port->name, sizeof(port->name), "port%u", id);
            snprintf(port->media,
                     sizeof(port->media),
                     "%s",
                     (id <= 8) ? "copper" : "sfp");

            port->present = true;
            return port;
        }
    }

    return NULL;
}

static int if_index_from_vb(const SnmpVarBind *vb,
                            size_t baseLen,
                            uint32_t *idx) {

    if (vb->oidLen != baseLen + 1) {
        return -1;
    }

    *idx = vb->oid[baseLen];
    return 0;
}

static void copy_snmp_string(char *dst, size_t dstLen, const char *src) {

    size_t n;

    if (!dst || dstLen == 0 || !src) {
        return;
    }

    n = strlen(src);
    if (n >= dstLen) {
        n = dstLen - 1;
    }

    memcpy(dst, src, n);
    dst[n] = '\0';
}

static int walk_if_name(const SnmpVarBind *vb, void *arg) {

    WalkCtx *w;
    SwitchPortState *p;
    uint32_t id;

    w = (WalkCtx *)arg;

    if (if_index_from_vb(vb,
                         sizeof(OID_ifName) / sizeof(OID_ifName[0]),
                         &id) != 0) {
        return 0;
    }

    p = ensure_port(w, id);
    if (p && vb->value.type == SNMP_VALUE_STRING) {
        copy_snmp_string(p->name,
                         sizeof(p->name),
                         vb->value.stringValue);
    }

    return 0;
}

static int walk_if_descr(const SnmpVarBind *vb, void *arg) {

    WalkCtx *w;
    SwitchPortState *p;
    uint32_t id;

    w = (WalkCtx *)arg;

    if (if_index_from_vb(vb,
                         sizeof(OID_ifDescr) / sizeof(OID_ifDescr[0]),
                         &id) != 0) {
        return 0;
    }

    p = ensure_port(w, id);
    if (p && p->name[0] == '\0' &&
        vb->value.type == SNMP_VALUE_STRING) {
        copy_snmp_string(p->name,
                         sizeof(p->name),
                         vb->value.stringValue);
    }

    return 0;
}

#define WALK_NUM(NAME, FIELD, BASE)                                      \
static int walk_##NAME(const SnmpVarBind *vb, void *arg) {               \
    WalkCtx *w;                                                          \
    SwitchPortState *p;                                                  \
    uint32_t id;                                                         \
                                                                         \
    w = (WalkCtx *)arg;                                                  \
                                                                         \
    if (if_index_from_vb(vb,                                             \
                         sizeof(BASE) / sizeof((BASE)[0]),               \
                         &id) != 0) {                                    \
        return 0;                                                        \
    }                                                                    \
                                                                         \
    p = ensure_port(w, id);                                              \
    if (!p) {                                                            \
        return 0;                                                        \
    }                                                                    \
                                                                         \
    if (vb->value.type == SNMP_VALUE_INT) {                              \
        p->FIELD = (uint64_t)vb->value.intValue;                         \
    } else if (vb->value.type == SNMP_VALUE_UINT) {                      \
        p->FIELD = vb->value.uintValue;                                  \
    }                                                                    \
                                                                         \
    return 0;                                                            \
}

WALK_NUM(if_speed, speedBps, OID_ifSpeed)
WALK_NUM(if_in_octets, rxBytes, OID_ifInOctets)
WALK_NUM(if_out_octets, txBytes, OID_ifOutOctets)
WALK_NUM(if_hc_in_octets, rxBytes, OID_ifHCInOctets)
WALK_NUM(if_hc_out_octets, txBytes, OID_ifHCOutOctets)
WALK_NUM(if_in_ucast, rxPackets, OID_ifInUcast)
WALK_NUM(if_out_ucast, txPackets, OID_ifOutUcast)
WALK_NUM(if_in_err, rxErrors, OID_ifInErrors)
WALK_NUM(if_out_err, txErrors, OID_ifOutErrors)
WALK_NUM(if_in_drop, rxDrops, OID_ifInDiscards)
WALK_NUM(if_out_drop, txDrops, OID_ifOutDiscards)

static int walk_if_admin(const SnmpVarBind *vb, void *arg) {

    WalkCtx *w;
    SwitchPortState *p;
    uint32_t id;
    int64_t value;

    w = (WalkCtx *)arg;
    value = 0;

    if (if_index_from_vb(vb,
                         sizeof(OID_ifAdminStatus) /
                         sizeof(OID_ifAdminStatus[0]),
                         &id) != 0) {
        return 0;
    }

    p = ensure_port(w, id);
    if (!p) {
        return 0;
    }

    value = (vb->value.type == SNMP_VALUE_INT) ?
            vb->value.intValue : (int64_t)vb->value.uintValue;

    p->adminUp = (value == 1);

    return 0;
}

static int walk_if_oper(const SnmpVarBind *vb, void *arg) {

    WalkCtx *w;
    SwitchPortState *p;
    uint32_t id;
    int64_t value;

    w = (WalkCtx *)arg;
    value = 0;

    if (if_index_from_vb(vb,
                         sizeof(OID_ifOperStatus) /
                         sizeof(OID_ifOperStatus[0]),
                         &id) != 0) {
        return 0;
    }

    p = ensure_port(w, id);
    if (!p) {
        return 0;
    }

    value = (vb->value.type == SNMP_VALUE_INT) ?
            vb->value.intValue : (int64_t)vb->value.uintValue;

    p->linkUp = (value == 1);

    return 0;
}

static int walk_poe_exist(const SnmpVarBind *vb, void *arg) {

    WalkCtx *w;
    SwitchPortState *p;
    uint32_t id;
    int64_t value;

    w = (WalkCtx *)arg;
    value = 0;

    if (if_index_from_vb(vb,
                         sizeof(OID_poeExist) / sizeof(OID_poeExist[0]),
                         &id) != 0) {
        return 0;
    }

    p = ensure_port(w, id);
    if (!p) {
        return 0;
    }

    value = (vb->value.type == SNMP_VALUE_INT) ?
            vb->value.intValue : (int64_t)vb->value.uintValue;

    p->poeSupported = (value == 1);

    return 0;
}

static int walk_poe_admin(const SnmpVarBind *vb, void *arg) {

    WalkCtx *w;
    SwitchPortState *p;
    uint32_t id;
    int64_t value;

    w = (WalkCtx *)arg;
    value = 0;

    if (if_index_from_vb(vb,
                         sizeof(OID_poeAdmin) / sizeof(OID_poeAdmin[0]),
                         &id) != 0) {
        return 0;
    }

    p = ensure_port(w, id);
    if (!p) {
        return 0;
    }

    value = (vb->value.type == SNMP_VALUE_INT) ?
            vb->value.intValue : (int64_t)vb->value.uintValue;

    p->poeEnabled = (value == 1);

    return 0;
}

static int walk_poe_oper(const SnmpVarBind *vb, void *arg) {

    WalkCtx *w;
    SwitchPortState *p;
    uint32_t id;
    int64_t value;

    w = (WalkCtx *)arg;
    value = 0;

    if (if_index_from_vb(vb,
                         sizeof(OID_poeOperStatus) /
                         sizeof(OID_poeOperStatus[0]),
                         &id) != 0) {
        return 0;
    }

    p = ensure_port(w, id);
    if (!p) {
        return 0;
    }

    value = (vb->value.type == SNMP_VALUE_INT) ?
            vb->value.intValue : (int64_t)vb->value.uintValue;

    p->poeOperational = (value == 1);

    if (p->poeSupported && p->poeEnabled && !p->poeOperational) {
        snprintf(p->fault,
                 sizeof(p->fault),
                 "PoE enabled but not delivering power");
    }

    return 0;
}

static int walk_poe_power(const SnmpVarBind *vb, void *arg) {

    WalkCtx *w;
    SwitchPortState *p;
    uint32_t id;
    double value;

    w = (WalkCtx *)arg;
    value = 0.0;

    if (if_index_from_vb(vb,
                         sizeof(OID_poePower) / sizeof(OID_poePower[0]),
                         &id) != 0) {
        return 0;
    }

    p = ensure_port(w, id);
    if (!p) {
        return 0;
    }

    value = (vb->value.type == SNMP_VALUE_INT) ?
            (double)vb->value.intValue : (double)vb->value.uintValue;

    p->powerWatts = value;

    return 0;
}

static int walk_poe_current(const SnmpVarBind *vb, void *arg) {

    WalkCtx *w;
    SwitchPortState *p;
    uint32_t id;
    double value;

    w = (WalkCtx *)arg;
    value = 0.0;

    if (if_index_from_vb(vb,
                         sizeof(OID_poeCurrent) /
                         sizeof(OID_poeCurrent[0]),
                         &id) != 0) {
        return 0;
    }

    p = ensure_port(w, id);
    if (!p) {
        return 0;
    }

    value = (vb->value.type == SNMP_VALUE_INT) ?
            (double)vb->value.intValue : (double)vb->value.uintValue;

    p->currentAmps = value / 1000.0;

    return 0;
}

static int walk_poe_voltage(const SnmpVarBind *vb, void *arg) {

    WalkCtx *w;
    SwitchPortState *p;
    uint32_t id;
    double value;

    w = (WalkCtx *)arg;
    value = 0.0;

    if (if_index_from_vb(vb,
                         sizeof(OID_poeVoltage) /
                         sizeof(OID_poeVoltage[0]),
                         &id) != 0) {
        return 0;
    }

    p = ensure_port(w, id);
    if (!p) {
        return 0;
    }

    value = (vb->value.type == SNMP_VALUE_INT) ?
            (double)vb->value.intValue : (double)vb->value.uintValue;

    p->voltage = value / 10.0;

    return 0;
}

static int walk_poe_class(const SnmpVarBind *vb, void *arg) {

    WalkCtx *w;
    SwitchPortState *p;
    uint32_t id;

    w = (WalkCtx *)arg;

    if (if_index_from_vb(vb,
                         sizeof(OID_poeClass) / sizeof(OID_poeClass[0]),
                         &id) != 0) {
        return 0;
    }

    p = ensure_port(w, id);
    if (!p) {
        return 0;
    }

    p->poeClass = (int)((vb->value.type == SNMP_VALUE_INT) ?
                  vb->value.intValue : (int64_t)vb->value.uintValue);

    return 0;
}

static int tycon_init(SwitchdContext *ctx) {

    TyconPriv *priv;

    priv = calloc(1, sizeof(*priv));
    if (priv == NULL) {
        return SWITCHD_ERR_NOMEM;
    }

    snmp_session_init(&priv->session,
                      ctx->config.snmpHost,
                      ctx->config.snmpPort,
                      ctx->config.snmpCommunity,
                      ctx->config.snmpTimeoutMs,
                      ctx->config.snmpRetries);

    ctx->driver->priv = priv;

    return SWITCHD_OK;
}

static void tycon_cleanup(SwitchdContext *ctx) {

    if (ctx && ctx->driver && ctx->driver->priv) {
        free(ctx->driver->priv);
        ctx->driver->priv = NULL;
    }
}

static int tycon_probe(SwitchdContext *ctx) {

    TyconPriv *priv;
    SnmpVarBind vb;
    int ret;

    if (!ctx || !ctx->driver || !ctx->driver->priv) {
        return SWITCHD_ERR_INVAL;
    }

    priv = (TyconPriv *)ctx->driver->priv;
    memset(&vb, 0, sizeof(vb));

    usys_log_debug("tycon: probe snmp=%s:%d community=%s timeoutMs=%d "
                   "retries=%d",
                   priv->session.host,
                   priv->session.port,
                   priv->session.community,
                   priv->session.timeoutMs,
                   priv->session.retries);

    ret = snmp_get(&priv->session,
                   OID_softwareVersion,
                   sizeof(OID_softwareVersion) /
                   sizeof(OID_softwareVersion[0]),
                   &vb);

    if (ret == SWITCHD_OK) {
        usys_log_info("tycon: probe ok swType=%d sw='%s'",
                      vb.value.type,
                      vb.value.stringValue);
    } else {
        usys_log_error("tycon: probe failed ret=%d", ret);
    }

    return ret;
}

static int tycon_refresh_info(SwitchdContext *ctx,
                              SwitchInfo *info,
                              SwitchCapabilities *caps) {

    TyconPriv *priv;
    int64_t value;
    int ret;

    priv = (TyconPriv *)ctx->driver->priv;
    value = 0;

    ret = get_string(&priv->session,
                     OID_manufactureName,
                     sizeof(OID_manufactureName) /
                     sizeof(OID_manufactureName[0]),
                     info->vendor,
                     sizeof(info->vendor));
    if (ret != SWITCHD_OK) {
        snprintf(info->vendor, sizeof(info->vendor), "Tycon");
    }

    snprintf(info->model, sizeof(info->model), "TP-SW8GAT/BT/24-SFP");

    get_string(&priv->session,
               OID_serialNumber,
               sizeof(OID_serialNumber) / sizeof(OID_serialNumber[0]),
               info->serial,
               sizeof(info->serial));

    get_string(&priv->session,
               OID_hardwareVersion,
               sizeof(OID_hardwareVersion) / sizeof(OID_hardwareVersion[0]),
               info->hardwareVersion,
               sizeof(info->hardwareVersion));

    get_string(&priv->session,
               OID_softwareVersion,
               sizeof(OID_softwareVersion) / sizeof(OID_softwareVersion[0]),
               info->softwareVersion,
               sizeof(info->softwareVersion));

    ret = get_int(&priv->session,
                  OID_ifNumber,
                  sizeof(OID_ifNumber) / sizeof(OID_ifNumber[0]),
                  &value);
    if (ret != SWITCHD_OK) {
        return SWITCHD_ERR_SNMP;
    }

    info->portCount = (uint32_t)value;
    info->reachable = true;
    snprintf(info->managementAddress,
             sizeof(info->managementAddress),
             "%s",
             ctx->config.snmpHost);

    info->updatedAt = time(NULL);

    memset(caps, 0, sizeof(*caps));
    caps->supportsPortAdmin = true;
    caps->supportsPoeControl = true;
    caps->supportsPoeCycle = true;
    caps->supportsPortCounters = true;
    caps->supportsPowerMetrics = true;
    caps->supportsSystemMetrics = true;
    caps->supportsFirmwareUpdate = true;
    caps->supportsSaveConfig = true;
    caps->maxPorts = info->portCount;

    return SWITCHD_OK;
}

static int tycon_refresh_ports(SwitchdContext *ctx,
                               SwitchPortState *ports,
                               uint32_t *count) {

    TyconPriv *priv;
    WalkCtx w;
    uint32_t i;
    int ret;

    priv = (TyconPriv *)ctx->driver->priv;

    memset(&w, 0, sizeof(w));
    memset(ports, 0, sizeof(SwitchPortState) * SWITCHD_MAX_PORTS);

    w.ports = ports;
    w.max = SWITCHD_MAX_PORTS;

    ret = snmp_walk(&priv->session,
                    OID_ifName,
                    sizeof(OID_ifName) / sizeof(OID_ifName[0]),
                    walk_if_name,
                    &w);
    if (ret != SWITCHD_OK) {
        ret = snmp_walk(&priv->session,
                        OID_ifDescr,
                        sizeof(OID_ifDescr) / sizeof(OID_ifDescr[0]),
                        walk_if_descr,
                        &w);
    }

    snmp_walk(&priv->session,
              OID_ifAdminStatus,
              sizeof(OID_ifAdminStatus) / sizeof(OID_ifAdminStatus[0]),
              walk_if_admin,
              &w);

    snmp_walk(&priv->session,
              OID_ifOperStatus,
              sizeof(OID_ifOperStatus) / sizeof(OID_ifOperStatus[0]),
              walk_if_oper,
              &w);

    snmp_walk(&priv->session,
              OID_ifSpeed,
              sizeof(OID_ifSpeed) / sizeof(OID_ifSpeed[0]),
              walk_if_speed,
              &w);

    snmp_walk(&priv->session,
              OID_ifHCInOctets,
              sizeof(OID_ifHCInOctets) / sizeof(OID_ifHCInOctets[0]),
              walk_if_hc_in_octets,
              &w);

    snmp_walk(&priv->session,
              OID_ifHCOutOctets,
              sizeof(OID_ifHCOutOctets) / sizeof(OID_ifHCOutOctets[0]),
              walk_if_hc_out_octets,
              &w);

    snmp_walk(&priv->session,
              OID_ifInOctets,
              sizeof(OID_ifInOctets) / sizeof(OID_ifInOctets[0]),
              walk_if_in_octets,
              &w);

    snmp_walk(&priv->session,
              OID_ifOutOctets,
              sizeof(OID_ifOutOctets) / sizeof(OID_ifOutOctets[0]),
              walk_if_out_octets,
              &w);

    snmp_walk(&priv->session,
              OID_ifInUcast,
              sizeof(OID_ifInUcast) / sizeof(OID_ifInUcast[0]),
              walk_if_in_ucast,
              &w);

    snmp_walk(&priv->session,
              OID_ifOutUcast,
              sizeof(OID_ifOutUcast) / sizeof(OID_ifOutUcast[0]),
              walk_if_out_ucast,
              &w);

    snmp_walk(&priv->session,
              OID_ifInErrors,
              sizeof(OID_ifInErrors) / sizeof(OID_ifInErrors[0]),
              walk_if_in_err,
              &w);

    snmp_walk(&priv->session,
              OID_ifOutErrors,
              sizeof(OID_ifOutErrors) / sizeof(OID_ifOutErrors[0]),
              walk_if_out_err,
              &w);

    snmp_walk(&priv->session,
              OID_ifInDiscards,
              sizeof(OID_ifInDiscards) / sizeof(OID_ifInDiscards[0]),
              walk_if_in_drop,
              &w);

    snmp_walk(&priv->session,
              OID_ifOutDiscards,
              sizeof(OID_ifOutDiscards) / sizeof(OID_ifOutDiscards[0]),
              walk_if_out_drop,
              &w);

    snmp_walk(&priv->session,
              OID_poeExist,
              sizeof(OID_poeExist) / sizeof(OID_poeExist[0]),
              walk_poe_exist,
              &w);

    snmp_walk(&priv->session,
              OID_poeAdmin,
              sizeof(OID_poeAdmin) / sizeof(OID_poeAdmin[0]),
              walk_poe_admin,
              &w);

    snmp_walk(&priv->session,
              OID_poeOperStatus,
              sizeof(OID_poeOperStatus) /
              sizeof(OID_poeOperStatus[0]),
              walk_poe_oper,
              &w);

    snmp_walk(&priv->session,
              OID_poePower,
              sizeof(OID_poePower) / sizeof(OID_poePower[0]),
              walk_poe_power,
              &w);

    snmp_walk(&priv->session,
              OID_poeCurrent,
              sizeof(OID_poeCurrent) / sizeof(OID_poeCurrent[0]),
              walk_poe_current,
              &w);

    snmp_walk(&priv->session,
              OID_poeVoltage,
              sizeof(OID_poeVoltage) / sizeof(OID_poeVoltage[0]),
              walk_poe_voltage,
              &w);

    snmp_walk(&priv->session,
              OID_poeClass,
              sizeof(OID_poeClass) / sizeof(OID_poeClass[0]),
              walk_poe_class,
              &w);

    *count = 0;
    for (i = 0; i < SWITCHD_MAX_PORTS; i++) {
        if (ports[i].id == 0) {
            continue;
        }

        ports[i].updatedAt = time(NULL);
        (*count)++;
    }

    return SWITCHD_OK;
}

static double normalize_watts(int64_t value) {

    if (value > 1000) {
        return ((double)value) / 1000.0;
    }

    return (double)value;
}

static double normalize_voltage(int64_t value) {

    if (value > 1000) {
        return ((double)value) / 1000.0;
    }

    return (double)value;
}

static double normalize_current(int64_t value) {

    if (value > 50) {
        return ((double)value) / 1000.0;
    }

    return (double)value;
}

static int get_double_any(SnmpSession *s,
                          const uint32_t *oid,
                          size_t oidLen,
                          double *out,
                          double (*normalize)(int64_t)) {

    SnmpVarBind vb;
    int ret;

    if (!s || !oid || !out) {
        return SWITCHD_ERR_INVAL;
    }

    ret = snmp_get(s, oid, oidLen, &vb);
    if (ret != SWITCHD_OK) {
        return ret;
    }

    if (vb.value.type == SNMP_VALUE_STRING) {
        *out = parse_prefixed_double(vb.value.stringValue);
        return SWITCHD_OK;
    }

    if (vb.value.type == SNMP_VALUE_INT) {
        *out = normalize ? normalize(vb.value.intValue) :
                           (double)vb.value.intValue;
        return SWITCHD_OK;
    }

    if (vb.value.type == SNMP_VALUE_UINT) {
        *out = normalize ? normalize((int64_t)vb.value.uintValue) :
                           (double)vb.value.uintValue;
        return SWITCHD_OK;
    }

    return SWITCHD_ERR_PROTOCOL;
}

static int get_bool_any(SnmpSession *s,
                        const uint32_t *oid,
                        size_t oidLen,
                        bool *out) {

    SnmpVarBind vb;
    int ret;

    if (!s || !oid || !out) {
        return SWITCHD_ERR_INVAL;
    }

    ret = snmp_get(s, oid, oidLen, &vb);
    if (ret != SWITCHD_OK) {
        return ret;
    }

    if (vb.value.type == SNMP_VALUE_INT) {
        *out = (vb.value.intValue != 0);
        return SWITCHD_OK;
    }

    if (vb.value.type == SNMP_VALUE_UINT) {
        *out = (vb.value.uintValue != 0);
        return SWITCHD_OK;
    }

    if (vb.value.type == SNMP_VALUE_STRING) {
        *out = (parse_prefixed_double(vb.value.stringValue) != 0.0);
        return SWITCHD_OK;
    }

    return SWITCHD_ERR_PROTOCOL;
}

static int read_double_metric(SnmpSession *session,
                              const char *name,
                              const uint32_t *oid,
                              size_t oidLen,
                              double *value,
                              double (*normalize)(int64_t),
                              int *okCount) {

    int ret;

    ret = get_double_any(session, oid, oidLen, value, normalize);
    if (ret == SWITCHD_OK) {
        (*okCount)++;
        usys_log_debug("tycon: kpi %-28s ok value=%.3f",
                       name,
                       *value);
    } else {
        usys_log_error("tycon: kpi %-28s failed ret=%d", name, ret);
    }

    return ret;
}

static int read_bool_metric(SnmpSession *session,
                            const char *name,
                            const uint32_t *oid,
                            size_t oidLen,
                            bool *value,
                            int *okCount) {

    int ret;

    ret = get_bool_any(session, oid, oidLen, value);
    if (ret == SWITCHD_OK) {
        (*okCount)++;
        usys_log_debug("tycon: kpi %-28s ok value=%d", name, *value);
    } else {
        usys_log_error("tycon: kpi %-28s failed ret=%d", name, ret);
    }

    return ret;
}

static int tycon_refresh_kpis(SwitchdContext *ctx, SwitchKpis *kpis) {

    TyconPriv *priv;
    double dv;
    bool bv;
    int okCount;

    if (!ctx || !ctx->driver || !ctx->driver->priv || !kpis) {
        return SWITCHD_ERR_INVAL;
    }

    priv    = (TyconPriv *)ctx->driver->priv;
    okCount = 0;

    memset(kpis, 0, sizeof(*kpis));

    usys_log_debug("tycon: refresh_kpis begin snmp=%s:%d",
                   priv->session.host,
                   priv->session.port);

    dv = 0.0;
    if (read_double_metric(&priv->session,
                           "poe_total_power_watts",
                           OID_poeTotalPowerConsumption,
                           sizeof(OID_poeTotalPowerConsumption) /
                           sizeof(OID_poeTotalPowerConsumption[0]),
                           &dv,
                           normalize_watts,
                           &okCount) == SWITCHD_OK) {
        kpis->poeTotalPowerWatts = dv;
    }

    dv = 0.0;
    if (read_double_metric(&priv->session,
                           "poe_max_power_watts",
                           OID_poeTotalMaxPower,
                           sizeof(OID_poeTotalMaxPower) /
                           sizeof(OID_poeTotalMaxPower[0]),
                           &dv,
                           normalize_watts,
                           &okCount) == SWITCHD_OK) {
        kpis->poeMaxPowerWatts = dv;
    }

    dv = 0.0;
    if (read_double_metric(&priv->session,
                           "system_temperature_c",
                           OID_industrySystemTemperature,
                           sizeof(OID_industrySystemTemperature) /
                           sizeof(OID_industrySystemTemperature[0]),
                           &dv,
                           NULL,
                           &okCount) == SWITCHD_OK) {
        kpis->systemTemperatureC = dv;
    }

    dv = 0.0;
    if (read_double_metric(&priv->session,
                           "ambient_temperature_c",
                           OID_industryAmbientTemperature,
                           sizeof(OID_industryAmbientTemperature) /
                           sizeof(OID_industryAmbientTemperature[0]),
                           &dv,
                           NULL,
                           &okCount) == SWITCHD_OK) {
        kpis->ambientTemperatureC = dv;
    }

    dv = 0.0;
    if (read_double_metric(&priv->session,
                           "input_voltage",
                           OID_industryPowerIn,
                           sizeof(OID_industryPowerIn) /
                           sizeof(OID_industryPowerIn[0]),
                           &dv,
                           normalize_voltage,
                           &okCount) == SWITCHD_OK) {
        kpis->inputVoltage = dv;
    }

    dv = 0.0;
    if (read_double_metric(&priv->session,
                           "system_current_amps",
                           OID_industrySystemCurrent,
                           sizeof(OID_industrySystemCurrent) /
                           sizeof(OID_industrySystemCurrent[0]),
                           &dv,
                           normalize_current,
                           &okCount) == SWITCHD_OK) {
        kpis->systemCurrentAmps = dv;
    }

    dv = 0.0;
    if (read_double_metric(&priv->session,
                           "system_power_watts",
                           OID_industrySystemPower,
                           sizeof(OID_industrySystemPower) /
                           sizeof(OID_industrySystemPower[0]),
                           &dv,
                           normalize_watts,
                           &okCount) == SWITCHD_OK) {
        kpis->systemPowerWatts = dv;
    }

    bv = false;
    if (read_bool_metric(&priv->session,
                         "input_link_failure_alarm",
                         OID_industryOutAlarmPortLinkFail,
                         sizeof(OID_industryOutAlarmPortLinkFail) /
                         sizeof(OID_industryOutAlarmPortLinkFail[0]),
                         &bv,
                         &okCount) == SWITCHD_OK) {
        kpis->inputLinkFailureAlarm = bv;
    }

    bv = false;
    if (read_bool_metric(&priv->session,
                         "input_poe_failure_alarm",
                         OID_industryOutAlarmPortPoeFail,
                         sizeof(OID_industryOutAlarmPortPoeFail) /
                         sizeof(OID_industryOutAlarmPortPoeFail[0]),
                         &bv,
                         &okCount) == SWITCHD_OK) {
        kpis->inputPoeFailureAlarm = bv;
    }

    if (okCount == 0) {
        usys_log_error("tycon: refresh_kpis failed: no KPI OID read");
        return SWITCHD_ERR_SNMP;
    }

    kpis->updatedAt = time(NULL);

    usys_log_info("tycon: refresh_kpis ok okCount=%d poe=%.2fW "
                  "budget=%.2fW temp=%.2fC ambient=%.2fC "
                  "sysPower=%.2fW vin=%.2fV current=%.3fA "
                  "linkAlarm=%d poeAlarm=%d",
                  okCount,
                  kpis->poeTotalPowerWatts,
                  kpis->poeMaxPowerWatts,
                  kpis->systemTemperatureC,
                  kpis->ambientTemperatureC,
                  kpis->systemPowerWatts,
                  kpis->inputVoltage,
                  kpis->systemCurrentAmps,
                  kpis->inputLinkFailureAlarm,
                  kpis->inputPoeFailureAlarm);

    return SWITCHD_OK;
}

static int tycon_set_port_admin(SwitchdContext *ctx,
                                uint32_t portId,
                                bool up) {

    TyconPriv *priv;
    uint32_t oid[32];
    size_t n;

    priv = (TyconPriv *)ctx->driver->priv;
    n    = sizeof(OID_ifAdminStatus) / sizeof(OID_ifAdminStatus[0]);

    memcpy(oid, OID_ifAdminStatus, n * sizeof(uint32_t));
    oid[n++] = portId;

    return snmp_set_integer(&priv->session, oid, n, up ? 1 : 2);
}

static int tycon_set_port_poe(SwitchdContext *ctx,
                              uint32_t portId,
                              bool on) {

    TyconPriv *priv;
    uint32_t oid[32];
    size_t n;

    priv = (TyconPriv *)ctx->driver->priv;
    n    = sizeof(OID_poeAdmin) / sizeof(OID_poeAdmin[0]);

    memcpy(oid, OID_poeAdmin, n * sizeof(uint32_t));
    oid[n++] = portId;

    return snmp_set_integer(&priv->session, oid, n, on ? 1 : 0);
}

static int tycon_save_config(SwitchdContext *ctx) {

    TyconPriv *priv;

    priv = (TyconPriv *)ctx->driver->priv;

    return snmp_set_integer(&priv->session,
                            OID_sysSaveToNvm,
                            sizeof(OID_sysSaveToNvm) /
                            sizeof(OID_sysSaveToNvm[0]),
                            2);
}

static int tycon_reboot_switch(SwitchdContext *ctx) {

    TyconPriv *priv;

    priv = (TyconPriv *)ctx->driver->priv;

    return snmp_set_integer(&priv->session,
                            OID_sysReset,
                            sizeof(OID_sysReset) /
                            sizeof(OID_sysReset[0]),
                            2);
}

static int tycon_firmware_apply(SwitchdContext *ctx,
                                const SwitchFirmware *fw) {

    TyconPriv *priv;
    int ret;
    int64_t status;
    char tftpAddr[64];

    priv = (TyconPriv *)ctx->driver->priv;
    ret = SWITCHD_OK;
    status = 0;

    copy_snmp_string(tftpAddr, sizeof(tftpAddr), ctx->config.tftpBindIp);

    ret = snmp_set_string(&priv->session,
                          OID_commonLoadTftpAddress,
                          sizeof(OID_commonLoadTftpAddress) /
                          sizeof(OID_commonLoadTftpAddress[0]),
                          tftpAddr);
    if (ret != SWITCHD_OK) {
        return ret;
    }

    ret = snmp_set_string(&priv->session,
                          OID_commonLoadTftpFileName,
                          sizeof(OID_commonLoadTftpFileName) /
                          sizeof(OID_commonLoadTftpFileName[0]),
                          fw->tftpFilename);
    if (ret != SWITCHD_OK) {
        return ret;
    }

    ret = snmp_set_integer(&priv->session,
                           OID_commonLoadType,
                           sizeof(OID_commonLoadType) /
                           sizeof(OID_commonLoadType[0]),
                           1);
    if (ret != SWITCHD_OK) {
        return ret;
    }

    ret = snmp_set_integer(&priv->session,
                           OID_commonLoadExecute,
                           sizeof(OID_commonLoadExecute) /
                           sizeof(OID_commonLoadExecute[0]),
                           2);
    if (ret != SWITCHD_OK) {
        return ret;
    }

    for (int i = 0; i < 90; i++) {
        sleep(1);

        ret = get_int(&priv->session,
                      OID_commonLoadExecuteStatus,
                      sizeof(OID_commonLoadExecuteStatus) /
                      sizeof(OID_commonLoadExecuteStatus[0]),
                      &status);
        if (ret != SWITCHD_OK) {
            continue;
        }

        if (status == 2) {
            continue;
        }

        return (status == 3) ? SWITCHD_OK : SWITCHD_ERR_IO;
    }

    return SWITCHD_ERR_TIMEOUT;
}

int tycon_driver_attach(SwitchdContext *ctx) {

    SwitchDriver *driver;

    driver = calloc(1, sizeof(*driver));
    if (driver == NULL) {
        return SWITCHD_ERR_NOMEM;
    }

    driver->name               = "tycon_snmp";
    driver->ops.init           = tycon_init;
    driver->ops.cleanup        = tycon_cleanup;
    driver->ops.probe          = tycon_probe;
    driver->ops.refresh_info   = tycon_refresh_info;
    driver->ops.refresh_ports  = tycon_refresh_ports;
    driver->ops.refresh_kpis   = tycon_refresh_kpis;
    driver->ops.set_port_admin = tycon_set_port_admin;
    driver->ops.set_port_poe   = tycon_set_port_poe;
    driver->ops.save_config    = tycon_save_config;
    driver->ops.reboot_switch  = tycon_reboot_switch;
    driver->ops.firmware_apply = tycon_firmware_apply;

    ctx->driver = driver;

    return driver->ops.init(ctx);
}
