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

/* Standard IF-MIB/MIB-II */
static const uint32_t OID_ifNumber[]       = {1,3,6,1,2,1,2,1,0};
static const uint32_t OID_ifDescr[]        = {1,3,6,1,2,1,2,2,1,2};
static const uint32_t OID_ifSpeed[]        = {1,3,6,1,2,1,2,2,1,5};
static const uint32_t OID_ifAdminStatus[]  = {1,3,6,1,2,1,2,2,1,7};
static const uint32_t OID_ifOperStatus[]   = {1,3,6,1,2,1,2,2,1,8};
static const uint32_t OID_ifInOctets[]     = {1,3,6,1,2,1,2,2,1,10};
static const uint32_t OID_ifInUcast[]      = {1,3,6,1,2,1,2,2,1,11};
static const uint32_t OID_ifInDiscards[]   = {1,3,6,1,2,1,2,2,1,13};
static const uint32_t OID_ifInErrors[]     = {1,3,6,1,2,1,2,2,1,14};
static const uint32_t OID_ifOutOctets[]    = {1,3,6,1,2,1,2,2,1,16};
static const uint32_t OID_ifOutUcast[]     = {1,3,6,1,2,1,2,2,1,17};
static const uint32_t OID_ifOutDiscards[]  = {1,3,6,1,2,1,2,2,1,19};
static const uint32_t OID_ifOutErrors[]    = {1,3,6,1,2,1,2,2,1,20};
static const uint32_t OID_ifName[]         = {1,3,6,1,2,1,31,1,1,1,1};
static const uint32_t OID_ifHCInOctets[]   = {1,3,6,1,2,1,31,1,1,1,6};
static const uint32_t OID_ifHCOutOctets[]  = {1,3,6,1,2,1,31,1,1,1,10};

/* Tycon/IMI enterprise MIB: 1.3.6.1.4.1.12284.5 */
static const uint32_t OID_serialNumber[]             = {1,3,6,1,4,1,12284,5,1,7,0};
static const uint32_t OID_manufactureName[]          = {1,3,6,1,4,1,12284,5,1,8,0};
static const uint32_t OID_hardwareVersion[]          = {1,3,6,1,4,1,12284,5,1,9,0};
static const uint32_t OID_softwareVersion[]          = {1,3,6,1,4,1,12284,5,1,10,0};
static const uint32_t OID_commonLoadTftpAddress[]    = {1,3,6,1,4,1,12284,5,1,11,0};
static const uint32_t OID_commonLoadTftpFileName[]   = {1,3,6,1,4,1,12284,5,1,12,0};
static const uint32_t OID_commonLoadType[]           = {1,3,6,1,4,1,12284,5,1,13,0};
static const uint32_t OID_commonLoadExecute[]        = {1,3,6,1,4,1,12284,5,1,14,0};
static const uint32_t OID_commonLoadExecuteStatus[]  = {1,3,6,1,4,1,12284,5,1,15,0};
static const uint32_t OID_sysSaveToNvm[]             = {1,3,6,1,4,1,12284,5,1,1,0};
static const uint32_t OID_sysReset[]                 = {1,3,6,1,4,1,12284,5,1,2,0};

static const uint32_t OID_poeExist[]                 = {1,3,6,1,4,1,12284,5,2,1,1,2};
static const uint32_t OID_poeAdmin[]                 = {1,3,6,1,4,1,12284,5,2,1,1,3};
static const uint32_t OID_poeOperStatus[]            = {1,3,6,1,4,1,12284,5,2,1,1,4};
static const uint32_t OID_poePower[]                 = {1,3,6,1,4,1,12284,5,2,1,1,5};
static const uint32_t OID_poeCurrent[]               = {1,3,6,1,4,1,12284,5,2,1,1,6};
static const uint32_t OID_poeVoltage[]               = {1,3,6,1,4,1,12284,5,2,1,1,7};
static const uint32_t OID_poeClass[]                 = {1,3,6,1,4,1,12284,5,2,1,1,8};
static const uint32_t OID_poeTotalPowerConsumption[] = {1,3,6,1,4,1,12284,5,2,2,0};
static const uint32_t OID_poeTotalMaxPower[]         = {1,3,6,1,4,1,12284,5,2,3,0};

static const uint32_t OID_industrySystemTemperature[]     = {1,3,6,1,4,1,12284,5,6,3,0};
static const uint32_t OID_industryAmbientTemperature[]    = {1,3,6,1,4,1,12284,5,6,6,0};
static const uint32_t OID_industryPowerIn[]               = {1,3,6,1,4,1,12284,5,6,13,0};
static const uint32_t OID_industrySystemCurrent[]         = {1,3,6,1,4,1,12284,5,6,14,0};
static const uint32_t OID_industrySystemPower[]           = {1,3,6,1,4,1,12284,5,6,17,0};
static const uint32_t OID_industryOutAlarmPortLinkFail[]  = {1,3,6,1,4,1,12284,5,6,28,0};
static const uint32_t OID_industryOutAlarmPortPoeFail[]   = {1,3,6,1,4,1,12284,5,6,29,0};

typedef struct {
    SnmpSession session;
} TyconPriv;

typedef struct {
    SwitchPortState *ports;
    uint32_t max;
} WalkCtx;

static int get_string(SnmpSession *s, const uint32_t *oid, size_t oidLen, char *dst, size_t dstLen) {
    SnmpVarBind vb;
    int ret = snmp_get(s, oid, oidLen, &vb);
    if (ret != SWITCHD_OK) return ret;
    if (vb.value.type != SNMP_VALUE_STRING) return SWITCHD_ERR_PROTOCOL;
    if (dstLen == 0) return SWITCHD_ERR_INVAL;
    { size_t n = strlen(vb.value.stringValue); if (n >= dstLen) n = dstLen - 1; memcpy(dst, vb.value.stringValue, n); dst[n] = '\0'; }
    return SWITCHD_OK;
}

static int get_int(SnmpSession *s, const uint32_t *oid, size_t oidLen, int64_t *value) {
    SnmpVarBind vb;
    int ret = snmp_get(s, oid, oidLen, &vb);
    if (ret != SWITCHD_OK) return ret;
    if (vb.value.type == SNMP_VALUE_INT) *value = vb.value.intValue;
    else if (vb.value.type == SNMP_VALUE_UINT) *value = (int64_t)vb.value.uintValue;
    else return SWITCHD_ERR_PROTOCOL;
    return SWITCHD_OK;
}

static SwitchPortState *ensure_port(WalkCtx *w, uint32_t id) {
    uint32_t i;
    if (id == 0 || id > w->max) return NULL;
    for (i = 0; i < w->max; i++) {
        if (w->ports[i].id == id) return &w->ports[i];
        if (w->ports[i].id == 0) {
            w->ports[i].id = id;
            snprintf(w->ports[i].name, sizeof(w->ports[i].name), "port%u", id);
            snprintf(w->ports[i].media, sizeof(w->ports[i].media), "%s", (id <= 8) ? "copper" : "sfp");
            w->ports[i].present = true;
            return &w->ports[i];
        }
    }
    return NULL;
}

static int if_index_from_vb(const SnmpVarBind *vb, size_t baseLen, uint32_t *idx) {
    if (vb->oidLen != baseLen + 1) return -1;
    *idx = vb->oid[baseLen];
    return 0;
}

static int walk_if_name(const SnmpVarBind *vb, void *arg) {
    WalkCtx *w = (WalkCtx *)arg;
    uint32_t id;
    SwitchPortState *p;
    if (if_index_from_vb(vb, sizeof(OID_ifName)/sizeof(OID_ifName[0]), &id) != 0) return 0;
    p = ensure_port(w, id);
    if (p && vb->value.type == SNMP_VALUE_STRING) { size_t n = strlen(vb->value.stringValue); if (n >= sizeof(p->name)) n = sizeof(p->name) - 1; memcpy(p->name, vb->value.stringValue, n); p->name[n] = '\0'; }
    return 0;
}

static int walk_if_descr(const SnmpVarBind *vb, void *arg) {
    WalkCtx *w = (WalkCtx *)arg;
    uint32_t id;
    SwitchPortState *p;
    if (if_index_from_vb(vb, sizeof(OID_ifDescr)/sizeof(OID_ifDescr[0]), &id) != 0) return 0;
    p = ensure_port(w, id);
    if (p && p->name[0] == '\0' && vb->value.type == SNMP_VALUE_STRING) { size_t n = strlen(vb->value.stringValue); if (n >= sizeof(p->name)) n = sizeof(p->name) - 1; memcpy(p->name, vb->value.stringValue, n); p->name[n] = '\0'; }
    return 0;
}

#define WALK_NUM(NAME, FIELD, BASE) \
static int walk_##NAME(const SnmpVarBind *vb, void *arg) { \
    WalkCtx *w = (WalkCtx *)arg; uint32_t id; SwitchPortState *p; \
    if (if_index_from_vb(vb, sizeof(BASE)/sizeof((BASE)[0]), &id) != 0) return 0; \
    p = ensure_port(w, id); if (!p) return 0; \
    if (vb->value.type == SNMP_VALUE_INT) p->FIELD = (uint64_t)vb->value.intValue; \
    else if (vb->value.type == SNMP_VALUE_UINT) p->FIELD = vb->value.uintValue; \
    return 0; }

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
    WalkCtx *w = (WalkCtx *)arg; uint32_t id; SwitchPortState *p;
    if (if_index_from_vb(vb, sizeof(OID_ifAdminStatus)/sizeof(OID_ifAdminStatus[0]), &id) != 0) return 0;
    p = ensure_port(w, id); if (!p) return 0;
    p->adminUp = ((vb->value.type == SNMP_VALUE_INT ? vb->value.intValue : (int64_t)vb->value.uintValue) == 1);
    return 0;
}

static int walk_if_oper(const SnmpVarBind *vb, void *arg) {
    WalkCtx *w = (WalkCtx *)arg; uint32_t id; SwitchPortState *p;
    if (if_index_from_vb(vb, sizeof(OID_ifOperStatus)/sizeof(OID_ifOperStatus[0]), &id) != 0) return 0;
    p = ensure_port(w, id); if (!p) return 0;
    p->linkUp = ((vb->value.type == SNMP_VALUE_INT ? vb->value.intValue : (int64_t)vb->value.uintValue) == 1);
    return 0;
}

static int walk_poe_exist(const SnmpVarBind *vb, void *arg) {
    WalkCtx *w = (WalkCtx *)arg; uint32_t id; SwitchPortState *p;
    if (if_index_from_vb(vb, sizeof(OID_poeExist)/sizeof(OID_poeExist[0]), &id) != 0) return 0;
    p = ensure_port(w, id); if (!p) return 0;
    p->poeSupported = ((vb->value.type == SNMP_VALUE_INT ? vb->value.intValue : (int64_t)vb->value.uintValue) == 1);
    return 0;
}

static int walk_poe_admin(const SnmpVarBind *vb, void *arg) {
    WalkCtx *w = (WalkCtx *)arg; uint32_t id; SwitchPortState *p;
    if (if_index_from_vb(vb, sizeof(OID_poeAdmin)/sizeof(OID_poeAdmin[0]), &id) != 0) return 0;
    p = ensure_port(w, id); if (!p) return 0;
    p->poeEnabled = ((vb->value.type == SNMP_VALUE_INT ? vb->value.intValue : (int64_t)vb->value.uintValue) == 1);
    return 0;
}

static int walk_poe_oper(const SnmpVarBind *vb, void *arg) {
    WalkCtx *w = (WalkCtx *)arg; uint32_t id; SwitchPortState *p;
    if (if_index_from_vb(vb, sizeof(OID_poeOperStatus)/sizeof(OID_poeOperStatus[0]), &id) != 0) return 0;
    p = ensure_port(w, id); if (!p) return 0;
    p->poeOperational = ((vb->value.type == SNMP_VALUE_INT ? vb->value.intValue : (int64_t)vb->value.uintValue) == 1);
    if (p->poeSupported && p->poeEnabled && !p->poeOperational) snprintf(p->fault, sizeof(p->fault), "PoE enabled but not delivering power");
    return 0;
}

static int walk_poe_power(const SnmpVarBind *vb, void *arg) {
    WalkCtx *w = (WalkCtx *)arg; uint32_t id; SwitchPortState *p; double v;
    if (if_index_from_vb(vb, sizeof(OID_poePower)/sizeof(OID_poePower[0]), &id) != 0) return 0;
    p = ensure_port(w, id); if (!p) return 0;
    v = (vb->value.type == SNMP_VALUE_INT) ? (double)vb->value.intValue : (double)vb->value.uintValue;
    p->powerWatts = v;
    return 0;
}

static int walk_poe_current(const SnmpVarBind *vb, void *arg) {
    WalkCtx *w = (WalkCtx *)arg; uint32_t id; SwitchPortState *p; double v;
    if (if_index_from_vb(vb, sizeof(OID_poeCurrent)/sizeof(OID_poeCurrent[0]), &id) != 0) return 0;
    p = ensure_port(w, id); if (!p) return 0;
    v = (vb->value.type == SNMP_VALUE_INT) ? (double)vb->value.intValue : (double)vb->value.uintValue;
    p->currentAmps = v / 1000.0;
    return 0;
}

static int walk_poe_voltage(const SnmpVarBind *vb, void *arg) {
    WalkCtx *w = (WalkCtx *)arg; uint32_t id; SwitchPortState *p; double v;
    if (if_index_from_vb(vb, sizeof(OID_poeVoltage)/sizeof(OID_poeVoltage[0]), &id) != 0) return 0;
    p = ensure_port(w, id); if (!p) return 0;
    v = (vb->value.type == SNMP_VALUE_INT) ? (double)vb->value.intValue : (double)vb->value.uintValue;
    p->voltage = v / 10.0;
    return 0;
}

static int walk_poe_class(const SnmpVarBind *vb, void *arg) {
    WalkCtx *w = (WalkCtx *)arg; uint32_t id; SwitchPortState *p;
    if (if_index_from_vb(vb, sizeof(OID_poeClass)/sizeof(OID_poeClass[0]), &id) != 0) return 0;
    p = ensure_port(w, id); if (!p) return 0;
    p->poeClass = (int)((vb->value.type == SNMP_VALUE_INT) ? vb->value.intValue : (int64_t)vb->value.uintValue);
    return 0;
}

static int tycon_init(SwitchdContext *ctx) {
    TyconPriv *priv = calloc(1, sizeof(*priv));
    if (priv == NULL) return SWITCHD_ERR_NOMEM;
    snmp_session_init(&priv->session, ctx->config.snmpHost, ctx->config.snmpPort,
                      ctx->config.snmpCommunity, ctx->config.snmpTimeoutMs, ctx->config.snmpRetries);
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
    TyconPriv *priv = (TyconPriv *)ctx->driver->priv;
    SnmpVarBind vb;
    int ret = snmp_get(&priv->session, OID_softwareVersion, sizeof(OID_softwareVersion)/sizeof(OID_softwareVersion[0]), &vb);
    return ret;
}

static int tycon_refresh_info(SwitchdContext *ctx, SwitchInfo *info, SwitchCapabilities *caps) {
    TyconPriv *priv = (TyconPriv *)ctx->driver->priv;
    int64_t v;
    if (get_string(&priv->session, OID_manufactureName, sizeof(OID_manufactureName)/sizeof(OID_manufactureName[0]), info->vendor, sizeof(info->vendor)) != SWITCHD_OK) {
        snprintf(info->vendor, sizeof(info->vendor), "Tycon");
    }
    snprintf(info->model, sizeof(info->model), "TP-SW8GAT/BT/24-SFP");
    get_string(&priv->session, OID_serialNumber, sizeof(OID_serialNumber)/sizeof(OID_serialNumber[0]), info->serial, sizeof(info->serial));
    get_string(&priv->session, OID_hardwareVersion, sizeof(OID_hardwareVersion)/sizeof(OID_hardwareVersion[0]), info->hardwareVersion, sizeof(info->hardwareVersion));
    get_string(&priv->session, OID_softwareVersion, sizeof(OID_softwareVersion)/sizeof(OID_softwareVersion[0]), info->softwareVersion, sizeof(info->softwareVersion));
    if (get_int(&priv->session, OID_ifNumber, sizeof(OID_ifNumber)/sizeof(OID_ifNumber[0]), &v) != SWITCHD_OK) return SWITCHD_ERR_SNMP;
    info->portCount = (uint32_t)v;
    info->reachable = true;
    snprintf(info->managementAddress, sizeof(info->managementAddress), "%s", ctx->config.snmpHost);
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

static int tycon_refresh_ports(SwitchdContext *ctx, SwitchPortState *ports, uint32_t *count) {
    TyconPriv *priv = (TyconPriv *)ctx->driver->priv;
    WalkCtx w;
    uint32_t i;
    int ret;
    memset(&w, 0, sizeof(w));
    memset(ports, 0, sizeof(SwitchPortState) * SWITCHD_MAX_PORTS);
    w.ports = ports;
    w.max = SWITCHD_MAX_PORTS;
    ret = snmp_walk(&priv->session, OID_ifName, sizeof(OID_ifName)/sizeof(OID_ifName[0]), walk_if_name, &w);
    if (ret != SWITCHD_OK) ret = snmp_walk(&priv->session, OID_ifDescr, sizeof(OID_ifDescr)/sizeof(OID_ifDescr[0]), walk_if_descr, &w);
    snmp_walk(&priv->session, OID_ifAdminStatus, sizeof(OID_ifAdminStatus)/sizeof(OID_ifAdminStatus[0]), walk_if_admin, &w);
    snmp_walk(&priv->session, OID_ifOperStatus, sizeof(OID_ifOperStatus)/sizeof(OID_ifOperStatus[0]), walk_if_oper, &w);
    snmp_walk(&priv->session, OID_ifSpeed, sizeof(OID_ifSpeed)/sizeof(OID_ifSpeed[0]), walk_if_speed, &w);
    snmp_walk(&priv->session, OID_ifHCInOctets, sizeof(OID_ifHCInOctets)/sizeof(OID_ifHCInOctets[0]), walk_if_hc_in_octets, &w);
    snmp_walk(&priv->session, OID_ifHCOutOctets, sizeof(OID_ifHCOutOctets)/sizeof(OID_ifHCOutOctets[0]), walk_if_hc_out_octets, &w);
    snmp_walk(&priv->session, OID_ifInOctets, sizeof(OID_ifInOctets)/sizeof(OID_ifInOctets[0]), walk_if_in_octets, &w);
    snmp_walk(&priv->session, OID_ifOutOctets, sizeof(OID_ifOutOctets)/sizeof(OID_ifOutOctets[0]), walk_if_out_octets, &w);
    snmp_walk(&priv->session, OID_ifInUcast, sizeof(OID_ifInUcast)/sizeof(OID_ifInUcast[0]), walk_if_in_ucast, &w);
    snmp_walk(&priv->session, OID_ifOutUcast, sizeof(OID_ifOutUcast)/sizeof(OID_ifOutUcast[0]), walk_if_out_ucast, &w);
    snmp_walk(&priv->session, OID_ifInErrors, sizeof(OID_ifInErrors)/sizeof(OID_ifInErrors[0]), walk_if_in_err, &w);
    snmp_walk(&priv->session, OID_ifOutErrors, sizeof(OID_ifOutErrors)/sizeof(OID_ifOutErrors[0]), walk_if_out_err, &w);
    snmp_walk(&priv->session, OID_ifInDiscards, sizeof(OID_ifInDiscards)/sizeof(OID_ifInDiscards[0]), walk_if_in_drop, &w);
    snmp_walk(&priv->session, OID_ifOutDiscards, sizeof(OID_ifOutDiscards)/sizeof(OID_ifOutDiscards[0]), walk_if_out_drop, &w);

    snmp_walk(&priv->session, OID_poeExist, sizeof(OID_poeExist)/sizeof(OID_poeExist[0]), walk_poe_exist, &w);
    snmp_walk(&priv->session, OID_poeAdmin, sizeof(OID_poeAdmin)/sizeof(OID_poeAdmin[0]), walk_poe_admin, &w);
    snmp_walk(&priv->session, OID_poeOperStatus, sizeof(OID_poeOperStatus)/sizeof(OID_poeOperStatus[0]), walk_poe_oper, &w);
    snmp_walk(&priv->session, OID_poePower, sizeof(OID_poePower)/sizeof(OID_poePower[0]), walk_poe_power, &w);
    snmp_walk(&priv->session, OID_poeCurrent, sizeof(OID_poeCurrent)/sizeof(OID_poeCurrent[0]), walk_poe_current, &w);
    snmp_walk(&priv->session, OID_poeVoltage, sizeof(OID_poeVoltage)/sizeof(OID_poeVoltage[0]), walk_poe_voltage, &w);
    snmp_walk(&priv->session, OID_poeClass, sizeof(OID_poeClass)/sizeof(OID_poeClass[0]), walk_poe_class, &w);

    *count = 0;
    for (i = 0; i < SWITCHD_MAX_PORTS; i++) {
        if (ports[i].id == 0) continue;
        ports[i].updatedAt = time(NULL);
        (*count)++;
    }
    return SWITCHD_OK;
}

static int tycon_refresh_kpis(SwitchdContext *ctx, SwitchKpis *kpis) {
    TyconPriv *priv = (TyconPriv *)ctx->driver->priv;
    int64_t iv;
    char s[128];
    memset(kpis, 0, sizeof(*kpis));
    if (get_int(&priv->session, OID_poeTotalPowerConsumption, sizeof(OID_poeTotalPowerConsumption)/sizeof(OID_poeTotalPowerConsumption[0]), &iv) == SWITCHD_OK) {
        kpis->poeTotalPowerWatts = (double)iv;
    }
    if (get_int(&priv->session, OID_poeTotalMaxPower, sizeof(OID_poeTotalMaxPower)/sizeof(OID_poeTotalMaxPower[0]), &iv) == SWITCHD_OK) {
        kpis->poeMaxPowerWatts = (double)iv;
    }
    if (get_string(&priv->session, OID_industrySystemTemperature, sizeof(OID_industrySystemTemperature)/sizeof(OID_industrySystemTemperature[0]), s, sizeof(s)) == SWITCHD_OK) {
        kpis->systemTemperatureC = parse_prefixed_double(s);
    }
    if (get_string(&priv->session, OID_industryAmbientTemperature, sizeof(OID_industryAmbientTemperature)/sizeof(OID_industryAmbientTemperature[0]), s, sizeof(s)) == SWITCHD_OK) {
        kpis->ambientTemperatureC = parse_prefixed_double(s);
    }
    if (get_string(&priv->session, OID_industryPowerIn, sizeof(OID_industryPowerIn)/sizeof(OID_industryPowerIn[0]), s, sizeof(s)) == SWITCHD_OK) {
        kpis->inputVoltage = parse_prefixed_double(s);
    }
    if (get_string(&priv->session, OID_industrySystemCurrent, sizeof(OID_industrySystemCurrent)/sizeof(OID_industrySystemCurrent[0]), s, sizeof(s)) == SWITCHD_OK) {
        kpis->systemCurrentAmps = parse_prefixed_double(s);
    }
    if (get_string(&priv->session, OID_industrySystemPower, sizeof(OID_industrySystemPower)/sizeof(OID_industrySystemPower[0]), s, sizeof(s)) == SWITCHD_OK) {
        kpis->systemPowerWatts = parse_prefixed_double(s);
    }
    if (get_int(&priv->session, OID_industryOutAlarmPortLinkFail, sizeof(OID_industryOutAlarmPortLinkFail)/sizeof(OID_industryOutAlarmPortLinkFail[0]), &iv) == SWITCHD_OK) {
        kpis->inputLinkFailureAlarm = (iv != 0);
    }
    if (get_int(&priv->session, OID_industryOutAlarmPortPoeFail, sizeof(OID_industryOutAlarmPortPoeFail)/sizeof(OID_industryOutAlarmPortPoeFail[0]), &iv) == SWITCHD_OK) {
        kpis->inputPoeFailureAlarm = (iv != 0);
    }
    kpis->updatedAt = time(NULL);
    return SWITCHD_OK;
}

static int tycon_set_port_admin(SwitchdContext *ctx, uint32_t portId, bool up) {
    TyconPriv *priv = (TyconPriv *)ctx->driver->priv;
    uint32_t oid[32];
    size_t n = sizeof(OID_ifAdminStatus)/sizeof(OID_ifAdminStatus[0]);
    memcpy(oid, OID_ifAdminStatus, n * sizeof(uint32_t));
    oid[n++] = portId;
    return snmp_set_integer(&priv->session, oid, n, up ? 1 : 2);
}

static int tycon_set_port_poe(SwitchdContext *ctx, uint32_t portId, bool on) {
    TyconPriv *priv = (TyconPriv *)ctx->driver->priv;
    uint32_t oid[32];
    size_t n = sizeof(OID_poeAdmin)/sizeof(OID_poeAdmin[0]);
    memcpy(oid, OID_poeAdmin, n * sizeof(uint32_t));
    oid[n++] = portId;
    return snmp_set_integer(&priv->session, oid, n, on ? 1 : 0);
}

static int tycon_save_config(SwitchdContext *ctx) {
    TyconPriv *priv = (TyconPriv *)ctx->driver->priv;
    return snmp_set_integer(&priv->session, OID_sysSaveToNvm, sizeof(OID_sysSaveToNvm)/sizeof(OID_sysSaveToNvm[0]), 2);
}

static int tycon_reboot_switch(SwitchdContext *ctx) {
    TyconPriv *priv = (TyconPriv *)ctx->driver->priv;
    return snmp_set_integer(&priv->session, OID_sysReset, sizeof(OID_sysReset)/sizeof(OID_sysReset[0]), 2);
}

static int tycon_firmware_apply(SwitchdContext *ctx, const SwitchFirmware *fw) {
    TyconPriv *priv = (TyconPriv *)ctx->driver->priv;
    int ret;
    int64_t status;
    char tftpAddr[64];

    { size_t n = strlen(ctx->config.tftpBindIp); if (n >= sizeof(tftpAddr)) n = sizeof(tftpAddr) - 1; memcpy(tftpAddr, ctx->config.tftpBindIp, n); tftpAddr[n] = '\0'; }
    ret = snmp_set_string(&priv->session, OID_commonLoadTftpAddress, sizeof(OID_commonLoadTftpAddress)/sizeof(OID_commonLoadTftpAddress[0]), tftpAddr);
    if (ret != SWITCHD_OK) return ret;
    ret = snmp_set_string(&priv->session, OID_commonLoadTftpFileName, sizeof(OID_commonLoadTftpFileName)/sizeof(OID_commonLoadTftpFileName[0]), fw->tftpFilename);
    if (ret != SWITCHD_OK) return ret;
    ret = snmp_set_integer(&priv->session, OID_commonLoadType, sizeof(OID_commonLoadType)/sizeof(OID_commonLoadType[0]), 1);
    if (ret != SWITCHD_OK) return ret;
    ret = snmp_set_integer(&priv->session, OID_commonLoadExecute, sizeof(OID_commonLoadExecute)/sizeof(OID_commonLoadExecute[0]), 2);
    if (ret != SWITCHD_OK) return ret;
    /* Poll until not in progress. */
    for (int i = 0; i < 90; i++) {
        sleep(1);
        ret = get_int(&priv->session, OID_commonLoadExecuteStatus, sizeof(OID_commonLoadExecuteStatus)/sizeof(OID_commonLoadExecuteStatus[0]), &status);
        if (ret != SWITCHD_OK) continue;
        if (status == 2) continue;
        return (status == 3) ? SWITCHD_OK : SWITCHD_ERR_IO;
    }
    return SWITCHD_ERR_TIMEOUT;
}

int tycon_driver_attach(SwitchdContext *ctx) {
    SwitchDriver *d = calloc(1, sizeof(*d));
    if (d == NULL) return SWITCHD_ERR_NOMEM;
    d->name = "tycon_snmp";
    d->ops.init = tycon_init;
    d->ops.cleanup = tycon_cleanup;
    d->ops.probe = tycon_probe;
    d->ops.refresh_info = tycon_refresh_info;
    d->ops.refresh_ports = tycon_refresh_ports;
    d->ops.refresh_kpis = tycon_refresh_kpis;
    d->ops.set_port_admin = tycon_set_port_admin;
    d->ops.set_port_poe = tycon_set_port_poe;
    d->ops.save_config = tycon_save_config;
    d->ops.reboot_switch = tycon_reboot_switch;
    d->ops.firmware_apply = tycon_firmware_apply;
    ctx->driver = d;
    return d->ops.init(ctx);
}
