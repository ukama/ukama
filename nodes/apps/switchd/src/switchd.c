/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#include "alarms.h"
#include "config.h"
#include "driver.h"
#include "policy.h"
#include "switchd.h"
#include "tftp_server.h"
#include "utils.h"

#include "usys_log.h"

SwitchdContext gSwitchd;
static TftpServer gTftp;

typedef struct {
    uint64_t pollerLoops;
    uint64_t infoAttempts;
    uint64_t infoSuccess;
    uint64_t infoFailure;
    uint64_t kpiAttempts;
    uint64_t kpiSuccess;
    uint64_t kpiFailure;
    uint64_t portAttempts;
    uint64_t portSuccess;
    uint64_t portFailure;

    int lastProbeRet;
    int lastInfoRet;
    int lastKpisRet;
    int lastPortsRet;

    time_t pollerStartedAt;
    time_t lastPollerAt;
    time_t lastInfoAttemptAt;
    time_t lastInfoSuccessAt;
    time_t lastKpisAttemptAt;
    time_t lastKpisSuccessAt;
    time_t lastPortsAttemptAt;
    time_t lastPortsSuccessAt;

    char lastStage[32];
} SwitchdDebugStats;

static SwitchdDebugStats gDebug;

static void debug_set_stage(const char *stage) {
    snprintf(gDebug.lastStage, sizeof(gDebug.lastStage), "%s",
             stage ? stage : "");
}

static void debug_reset(void) {
    memset(&gDebug, 0, sizeof(gDebug));
    gDebug.lastProbeRet = SWITCHD_ERR_INTERNAL;
    gDebug.lastInfoRet = SWITCHD_ERR_INTERNAL;
    gDebug.lastKpisRet = SWITCHD_ERR_INTERNAL;
    gDebug.lastPortsRet = SWITCHD_ERR_INTERNAL;
    debug_set_stage("init");
}


static void op_reset(SwitchOperation *op) {
    memset(op, 0, sizeof(*op));
    op->state = SWITCHD_OP_STATE_IDLE;
}

static int op_begin(SwitchdContext *ctx,
                    SwitchdOperationType type,
                    uint32_t portId,
                    const char *detail) {
    if (pthread_mutex_trylock(&ctx->opMutex) != 0) {
        return SWITCHD_ERR_BUSY;
    }

    ctx->op.id = ++ctx->nextOpId;
    ctx->op.type = type;
    ctx->op.state = SWITCHD_OP_STATE_RUNNING;
    ctx->op.portId = portId;
    ctx->op.progress = 0;
    ctx->op.error = SWITCHD_OK;
    ctx->op.startedAt = time(NULL);
    ctx->op.endedAt = 0;
    snprintf(ctx->op.detail,
             sizeof(ctx->op.detail),
             "%s",
             detail ? detail : "");
    ctx->state = (type == SWITCHD_OP_FW_APPLY) ?
                 SWITCHD_STATE_UPDATING : SWITCHD_STATE_BUSY;
    return SWITCHD_OK;
}

static void op_end(SwitchdContext *ctx, int error, const char *detail) {
    ctx->op.state = (error == SWITCHD_OK) ?
                    SWITCHD_OP_STATE_DONE : SWITCHD_OP_STATE_FAILED;
    ctx->op.error = error;
    ctx->op.progress = (error == SWITCHD_OK) ? 100 : ctx->op.progress;
    ctx->op.endedAt = time(NULL);
    if (detail != NULL) {
        snprintf(ctx->op.detail, sizeof(ctx->op.detail), "%s", detail);
    }
    ctx->state = ctx->info.reachable ?
                 SWITCHD_STATE_READY : SWITCHD_STATE_DEGRADED;
    pthread_mutex_unlock(&ctx->opMutex);
}

static void *poller_main(void *arg) {
    SwitchdContext *ctx;
    uint64_t lastInfo;
    uint64_t lastKpis;
    uint64_t lastPorts;

    ctx = (SwitchdContext *)arg;
    lastInfo = 0;
    lastKpis = 0;
    lastPorts = 0;

    gDebug.pollerStartedAt = time(NULL);
    debug_set_stage("poller");

    while (!ctx->terminate) {
        uint64_t now;
        int ret;

        now = monotonic_msec();
        gDebug.pollerLoops++;
        gDebug.lastPollerAt = time(NULL);

        if (now - lastKpis >= (uint64_t)ctx->config.pollKpisSec * 1000ULL) {
            usys_log_debug("switchd: poll kpis begin");
            ret = switchd_refresh_kpis(ctx);
            gDebug.lastKpisRet = ret;

            if (ret != SWITCHD_OK) {
                usys_log_error("switchd: poll kpis failed: %d", ret);
            } else {
                usys_log_debug("switchd: poll kpis ok updatedAt=%ld",
                               (long)ctx->kpis.updatedAt);
            }

            lastKpis = now;
        }

        if (now - lastInfo >= (uint64_t)ctx->config.pollInfoSec * 1000ULL) {
            usys_log_debug("switchd: poll info begin");
            ret = switchd_refresh_info(ctx);
            gDebug.lastInfoRet = ret;

            if (ret != SWITCHD_OK) {
                usys_log_error("switchd: poll info failed: %d", ret);
            } else {
                usys_log_debug("switchd: poll info ok portCount=%u sw=%s",
                               ctx->info.portCount,
                               ctx->info.softwareVersion);
            }

            lastInfo = now;
        }

        if (now - lastPorts >=
            (uint64_t)ctx->config.pollStatusSec * 1000ULL) {
            usys_log_debug("switchd: poll ports begin");
            ret = switchd_refresh_ports(ctx);
            gDebug.lastPortsRet = ret;

            if (ret != SWITCHD_OK) {
                usys_log_error("switchd: poll ports failed: %d", ret);
            } else {
                usys_log_debug("switchd: poll ports ok count=%u",
                               ctx->portCount);
            }

            lastPorts = now;
        }

        sleep((unsigned int)ctx->config.pollStatusSec);
    }

    debug_set_stage("poller-exit");
    return NULL;
}

static void *alarm_main(void *arg) {
    SwitchdContext *ctx;

    ctx = (SwitchdContext *)arg;
    while (!ctx->terminate) {
        (void)alarms_scan(ctx);
        sleep((unsigned int)ctx->config.alarmScanSec);
    }

    return NULL;
}

int switchd_init(SwitchdContext *ctx) {
    debug_reset();
    memset(ctx, 0, sizeof(*ctx));
    pthread_mutex_init(&ctx->stateMutex, NULL);
    pthread_mutex_init(&ctx->opMutex, NULL);
    pthread_mutex_init(&ctx->alarmMutex, NULL);
    pthread_mutex_init(&ctx->driverMutex, NULL);

    ctx->state = SWITCHD_STATE_INIT;
    op_reset(&ctx->op);
    ctx->fw.state = SWITCHD_FW_IDLE;

    if (config_load(&ctx->config) != SWITCHD_OK) {
        return SWITCHD_ERR_INTERNAL;
    }
    if (mkdir_p(ctx->config.tftpRoot, 0755) != 0) {
        return SWITCHD_ERR_IO;
    }
    (void)policy_load(ctx);

    if (driver_init(ctx) != SWITCHD_OK) {
        return SWITCHD_ERR_INTERNAL;
    }

    debug_set_stage("probe");
    gDebug.lastProbeRet = ctx->driver->ops.probe(ctx);
    if (gDebug.lastProbeRet != SWITCHD_OK) {
        usys_log_error("switchd: probe failed: %d", gDebug.lastProbeRet);
        ctx->state = SWITCHD_STATE_DEGRADED;
        ctx->info.reachable = false;
    } else {
        usys_log_info("switchd: probe ok");

        debug_set_stage("initial-kpis");
        gDebug.lastKpisRet = switchd_refresh_kpis(ctx);
        usys_log_debug("switchd: initial kpis ret=%d updatedAt=%ld",
                       gDebug.lastKpisRet,
                       (long)ctx->kpis.updatedAt);

        debug_set_stage("initial-info");
        gDebug.lastInfoRet = switchd_refresh_info(ctx);
        usys_log_debug("switchd: initial info ret=%d portCount=%u",
                       gDebug.lastInfoRet,
                       ctx->portCount);

        debug_set_stage("initial-ports");
        gDebug.lastPortsRet = switchd_refresh_ports(ctx);
        usys_log_debug("switchd: initial ports ret=%d count=%u",
                       gDebug.lastPortsRet,
                       ctx->portCount);

        ctx->state = ctx->info.reachable ?
                     SWITCHD_STATE_READY : SWITCHD_STATE_DEGRADED;
    }

    debug_set_stage("ready");
    return SWITCHD_OK;
}

int switchd_start(SwitchdContext *ctx) {
    if (pthread_create(&ctx->pollerThread, NULL, poller_main, ctx) != 0) {
        return SWITCHD_ERR_INTERNAL;
    }
    ctx->pollerRunning = true;

    if (pthread_create(&ctx->alarmThread, NULL, alarm_main, ctx) != 0) {
        ctx->terminate = true;
        pthread_join(ctx->pollerThread, NULL);
        return SWITCHD_ERR_INTERNAL;
    }
    ctx->alarmRunning = true;

    return SWITCHD_OK;
}

void switchd_request_terminate(SwitchdContext *ctx) {
    if (ctx != NULL) {
        ctx->terminate = true;
        ctx->state = SWITCHD_STATE_TERMINATING;
    }
}

void switchd_stop(SwitchdContext *ctx) {
    switchd_request_terminate(ctx);

    if (ctx->pollerRunning) {
        pthread_join(ctx->pollerThread, NULL);
    }
    if (ctx->alarmRunning) {
        pthread_join(ctx->alarmThread, NULL);
    }

    tftp_server_stop(&gTftp);
}

void switchd_cleanup(SwitchdContext *ctx) {
    driver_cleanup(ctx);
    pthread_mutex_destroy(&ctx->stateMutex);
    pthread_mutex_destroy(&ctx->opMutex);
    pthread_mutex_destroy(&ctx->alarmMutex);
    pthread_mutex_destroy(&ctx->driverMutex);
}

int switchd_refresh_info(SwitchdContext *ctx) {
    SwitchInfo info;
    SwitchCapabilities caps;
    int ret;

    memset(&info, 0, sizeof(info));
    memset(&caps, 0, sizeof(caps));

    gDebug.infoAttempts++;
    gDebug.lastInfoAttemptAt = time(NULL);
    debug_set_stage("refresh-info");

    pthread_mutex_lock(&ctx->driverMutex);
    ret = ctx->driver->ops.refresh_info(ctx, &info, &caps);
    pthread_mutex_unlock(&ctx->driverMutex);

    gDebug.lastInfoRet = ret;
    if (ret != SWITCHD_OK) {
        gDebug.infoFailure++;
        ctx->info.reachable = false;
        ctx->state = SWITCHD_STATE_DEGRADED;
        usys_log_error("switchd: refresh_info failed ret=%d", ret);
        return ret;
    }

    gDebug.infoSuccess++;
    gDebug.lastInfoSuccessAt = time(NULL);

    ctx->info = info;
    ctx->caps = caps;
    ctx->portCount = info.portCount;
    ctx->info.reachable = true;

    usys_log_debug("switchd: refresh_info ok vendor=%s serial=%s sw=%s ports=%u",
                   ctx->info.vendor,
                   ctx->info.serial,
                   ctx->info.softwareVersion,
                   ctx->info.portCount);

    return SWITCHD_OK;
}

int switchd_refresh_ports(SwitchdContext *ctx) {
    SwitchPortState ports[SWITCHD_MAX_PORTS];
    uint32_t count;
    int ret;

    count = 0;
    memset(ports, 0, sizeof(ports));

    gDebug.portAttempts++;
    gDebug.lastPortsAttemptAt = time(NULL);
    debug_set_stage("refresh-ports");

    pthread_mutex_lock(&ctx->driverMutex);
    ret = ctx->driver->ops.refresh_ports(ctx, ports, &count);
    pthread_mutex_unlock(&ctx->driverMutex);

    gDebug.lastPortsRet = ret;
    if (ret != SWITCHD_OK) {
        gDebug.portFailure++;
        ctx->info.reachable = false;
        ctx->state = SWITCHD_STATE_DEGRADED;
        usys_log_error("switchd: refresh_ports failed ret=%d", ret);
        return ret;
    }

    gDebug.portSuccess++;
    gDebug.lastPortsSuccessAt = time(NULL);

    memcpy(ctx->ports, ports, sizeof(ports));
    ctx->portCount = count;
    ctx->info.reachable = true;
    if (ctx->state == SWITCHD_STATE_DEGRADED) {
        ctx->state = SWITCHD_STATE_READY;
    }

    usys_log_debug("switchd: refresh_ports ok count=%u", count);
    return SWITCHD_OK;
}

int switchd_refresh_kpis(SwitchdContext *ctx) {
    SwitchKpis kpis;
    int ret;

    memset(&kpis, 0, sizeof(kpis));

    gDebug.kpiAttempts++;
    gDebug.lastKpisAttemptAt = time(NULL);
    debug_set_stage("refresh-kpis");

    pthread_mutex_lock(&ctx->driverMutex);
    ret = ctx->driver->ops.refresh_kpis(ctx, &kpis);
    pthread_mutex_unlock(&ctx->driverMutex);

    gDebug.lastKpisRet = ret;
    if (ret != SWITCHD_OK) {
        gDebug.kpiFailure++;
        usys_log_error("switchd: refresh_kpis failed ret=%d", ret);
        return ret;
    }

    gDebug.kpiSuccess++;
    gDebug.lastKpisSuccessAt = time(NULL);

    ctx->kpis = kpis;
    ctx->info.reachable = true;

    if (ctx->state == SWITCHD_STATE_DEGRADED) {
        ctx->state = SWITCHD_STATE_READY;
    }

    usys_log_debug("switchd: refresh_kpis ok updatedAt=%ld poe=%.2f "
                   "budget=%.2f temp=%.2f ambient=%.2f power=%.2f "
                   "vin=%.2f current=%.3f linkAlarm=%d poeAlarm=%d",
                   (long)ctx->kpis.updatedAt,
                   ctx->kpis.poeTotalPowerWatts,
                   ctx->kpis.poeMaxPowerWatts,
                   ctx->kpis.systemTemperatureC,
                   ctx->kpis.ambientTemperatureC,
                   ctx->kpis.systemPowerWatts,
                   ctx->kpis.inputVoltage,
                   ctx->kpis.systemCurrentAmps,
                   ctx->kpis.inputLinkFailureAlarm,
                   ctx->kpis.inputPoeFailureAlarm);

    return SWITCHD_OK;
}

JsonObj *switchd_debug_status_json(SwitchdContext *ctx) {
    JsonObj *root;
    JsonObj *polls;
    JsonObj *last;
    JsonObj *cfg;
    JsonObj *cache;

    root = json_object();
    polls = json_object();
    last = json_object();
    cfg = json_object();
    cache = json_object();

    if (ctx == NULL) {
        return root;
    }

    json_object_set_new(root, "state", json_string(state_to_str(ctx->state)));
    json_object_set_new(root, "stage", json_string(gDebug.lastStage));
    json_object_set_new(root, "reachable", json_boolean(ctx->info.reachable));
    json_object_set_new(root, "driver",
                        json_string(ctx->driver ? ctx->driver->name : ""));
    json_object_set_new(root, "pollerRunning",
                        json_boolean(ctx->pollerRunning));
    json_object_set_new(root, "pollerLoops",
                        json_integer((json_int_t)gDebug.pollerLoops));

    json_object_set_new(cfg, "snmpHost", json_string(ctx->config.snmpHost));
    json_object_set_new(cfg, "snmpPort", json_integer(ctx->config.snmpPort));
    json_object_set_new(cfg, "snmpTimeoutMs",
                        json_integer(ctx->config.snmpTimeoutMs));
    json_object_set_new(cfg, "snmpRetries",
                        json_integer(ctx->config.snmpRetries));
    json_object_set_new(cfg, "pollKpisSec",
                        json_integer(ctx->config.pollKpisSec));
    json_object_set_new(cfg, "pollInfoSec",
                        json_integer(ctx->config.pollInfoSec));
    json_object_set_new(cfg, "pollStatusSec",
                        json_integer(ctx->config.pollStatusSec));
    json_object_set_new(root, "config", cfg);

    json_object_set_new(polls, "kpiAttempts",
                        json_integer((json_int_t)gDebug.kpiAttempts));
    json_object_set_new(polls, "kpiSuccess",
                        json_integer((json_int_t)gDebug.kpiSuccess));
    json_object_set_new(polls, "kpiFailure",
                        json_integer((json_int_t)gDebug.kpiFailure));
    json_object_set_new(polls, "infoAttempts",
                        json_integer((json_int_t)gDebug.infoAttempts));
    json_object_set_new(polls, "infoSuccess",
                        json_integer((json_int_t)gDebug.infoSuccess));
    json_object_set_new(polls, "infoFailure",
                        json_integer((json_int_t)gDebug.infoFailure));
    json_object_set_new(polls, "portAttempts",
                        json_integer((json_int_t)gDebug.portAttempts));
    json_object_set_new(polls, "portSuccess",
                        json_integer((json_int_t)gDebug.portSuccess));
    json_object_set_new(polls, "portFailure",
                        json_integer((json_int_t)gDebug.portFailure));
    json_object_set_new(root, "polls", polls);

    json_object_set_new(last, "probeRet", json_integer(gDebug.lastProbeRet));
    json_object_set_new(last, "kpisRet", json_integer(gDebug.lastKpisRet));
    json_object_set_new(last, "infoRet", json_integer(gDebug.lastInfoRet));
    json_object_set_new(last, "portsRet", json_integer(gDebug.lastPortsRet));
    json_object_set_new(last, "pollerStartedAt",
                        json_integer((json_int_t)gDebug.pollerStartedAt));
    json_object_set_new(last, "lastPollerAt",
                        json_integer((json_int_t)gDebug.lastPollerAt));
    json_object_set_new(last, "lastKpisAttemptAt",
                        json_integer((json_int_t)gDebug.lastKpisAttemptAt));
    json_object_set_new(last, "lastKpisSuccessAt",
                        json_integer((json_int_t)gDebug.lastKpisSuccessAt));
    json_object_set_new(last, "lastInfoAttemptAt",
                        json_integer((json_int_t)gDebug.lastInfoAttemptAt));
    json_object_set_new(last, "lastInfoSuccessAt",
                        json_integer((json_int_t)gDebug.lastInfoSuccessAt));
    json_object_set_new(last, "lastPortsAttemptAt",
                        json_integer((json_int_t)gDebug.lastPortsAttemptAt));
    json_object_set_new(last, "lastPortsSuccessAt",
                        json_integer((json_int_t)gDebug.lastPortsSuccessAt));
    json_object_set_new(root, "last", last);

    json_object_set_new(cache, "kpisUpdatedAt",
                        json_integer((json_int_t)ctx->kpis.updatedAt));
    json_object_set_new(cache, "infoUpdatedAt",
                        json_integer((json_int_t)ctx->info.updatedAt));
    json_object_set_new(cache, "portCount", json_integer(ctx->portCount));
    json_object_set_new(cache, "poeTotalPowerWatts",
                        json_real(ctx->kpis.poeTotalPowerWatts));
    json_object_set_new(cache, "poeMaxPowerWatts",
                        json_real(ctx->kpis.poeMaxPowerWatts));
    json_object_set_new(cache, "systemTemperatureC",
                        json_real(ctx->kpis.systemTemperatureC));
    json_object_set_new(cache, "ambientTemperatureC",
                        json_real(ctx->kpis.ambientTemperatureC));
    json_object_set_new(cache, "systemPowerWatts",
                        json_real(ctx->kpis.systemPowerWatts));
    json_object_set_new(cache, "inputVoltage",
                        json_real(ctx->kpis.inputVoltage));
    json_object_set_new(cache, "systemCurrentAmps",
                        json_real(ctx->kpis.systemCurrentAmps));
    json_object_set_new(root, "cache", cache);
    json_object_set_new(root, "policy", policy_serialize(ctx));

    return root;
}

SwitchPortState *switchd_get_port(SwitchdContext *ctx, uint32_t portId) {
    uint32_t i;

    for (i = 0; i < ctx->portCount && i < SWITCHD_MAX_PORTS; i++) {
        if (ctx->ports[i].id == portId) {
            return &ctx->ports[i];
        }
    }

    return NULL;
}

int switchd_set_port_admin(SwitchdContext *ctx, uint32_t portId, bool up) {
    int ret;

    if (switchd_get_port(ctx, portId) == NULL) {
        return SWITCHD_ERR_NOTFOUND;
    }

    ret = op_begin(ctx,
                   SWITCHD_OP_PORT_ADMIN_SET,
                   portId,
                   up ? "admin up" : "admin down");
    if (ret != SWITCHD_OK) {
        return ret;
    }

    pthread_mutex_lock(&ctx->driverMutex);
    ret = ctx->driver->ops.set_port_admin(ctx, portId, up);
    if (ret == SWITCHD_OK && ctx->config.saveAfterWrite &&
        ctx->driver->ops.save_config) {
        (void)ctx->driver->ops.save_config(ctx);
    }
    pthread_mutex_unlock(&ctx->driverMutex);

    if (ret == SWITCHD_OK) {
        (void)switchd_refresh_ports(ctx);
    }
    op_end(ctx, ret, (ret == SWITCHD_OK) ? "ok" : "set admin failed");
    return ret;
}

int switchd_set_port_poe(SwitchdContext *ctx, uint32_t portId, bool on) {
    int ret;
    SwitchPortState *port;

    port = switchd_get_port(ctx, portId);
    if (port == NULL) {
        return SWITCHD_ERR_NOTFOUND;
    }
    if (!port->poeSupported) {
        return SWITCHD_ERR_UNSUPPORTED;
    }

    ret = op_begin(ctx,
                   SWITCHD_OP_PORT_POE_SET,
                   portId,
                   on ? "poe on" : "poe off");
    if (ret != SWITCHD_OK) {
        return ret;
    }

    pthread_mutex_lock(&ctx->driverMutex);
    ret = ctx->driver->ops.set_port_poe(ctx, portId, on);
    if (ret == SWITCHD_OK && ctx->config.saveAfterWrite &&
        ctx->driver->ops.save_config) {
        (void)ctx->driver->ops.save_config(ctx);
    }
    pthread_mutex_unlock(&ctx->driverMutex);

    if (ret == SWITCHD_OK) {
        (void)switchd_refresh_ports(ctx);
    }
    op_end(ctx, ret, (ret == SWITCHD_OK) ? "ok" : "set poe failed");
    return ret;
}

int switchd_cycle_port_poe(SwitchdContext *ctx, uint32_t portId, int offMs) {
    SwitchPortState *port;
    struct timespec ts;
    int ret;

    port = switchd_get_port(ctx, portId);
    if (port == NULL) {
        return SWITCHD_ERR_NOTFOUND;
    }
    if (!port->poeSupported) {
        return SWITCHD_ERR_UNSUPPORTED;
    }

    ret = op_begin(ctx, SWITCHD_OP_PORT_POE_CYCLE, portId, "poe cycle");
    if (ret != SWITCHD_OK) {
        return ret;
    }

    pthread_mutex_lock(&ctx->driverMutex);
    ret = ctx->driver->ops.set_port_poe(ctx, portId, false);
    pthread_mutex_unlock(&ctx->driverMutex);

    if (ret == SWITCHD_OK) {
        ts.tv_sec = offMs / 1000;
        ts.tv_nsec = (long)(offMs % 1000) * 1000000L;
        nanosleep(&ts, NULL);

        pthread_mutex_lock(&ctx->driverMutex);
        ret = ctx->driver->ops.set_port_poe(ctx, portId, true);
        if (ret == SWITCHD_OK && ctx->config.saveAfterWrite &&
            ctx->driver->ops.save_config) {
            (void)ctx->driver->ops.save_config(ctx);
        }
        pthread_mutex_unlock(&ctx->driverMutex);
    }

    if (ret == SWITCHD_OK) {
        (void)switchd_refresh_ports(ctx);
    }
    op_end(ctx, ret, (ret == SWITCHD_OK) ? "ok" : "poe cycle failed");
    return ret;
}

int switchd_stage_firmware(SwitchdContext *ctx,
                           const char *path,
                           const char *version,
                           const char *sha256) {
    char dst[512];
    const char *base;
    struct stat st;
    int ret;
    size_t len;

    ret = op_begin(ctx, SWITCHD_OP_FW_STAGE, 0, "firmware stage");
    if (ret != SWITCHD_OK) {
        return ret;
    }

    if (path == NULL || access(path, R_OK) != 0) {
        op_end(ctx, SWITCHD_ERR_NOTFOUND, "firmware file not readable");
        return SWITCHD_ERR_NOTFOUND;
    }

    base = strrchr(path, '/');
    base = (base == NULL) ? path : base + 1;

    snprintf(dst, sizeof(dst), "%s/%s", ctx->config.tftpRoot, base);
    if (copy_file(path, dst) != 0 || stat(dst, &st) != 0) {
        op_end(ctx, SWITCHD_ERR_IO, "failed to stage firmware");
        return SWITCHD_ERR_IO;
    }

    memset(&ctx->fw, 0, sizeof(ctx->fw));

    len = strlen(dst);
    if (len >= sizeof(ctx->fw.path)) {
        len = sizeof(ctx->fw.path) - 1;
    }
    memcpy(ctx->fw.path, dst, len);
    ctx->fw.path[len] = '\0';

    len = strlen(base);
    if (len >= sizeof(ctx->fw.tftpFilename)) {
        len = sizeof(ctx->fw.tftpFilename) - 1;
    }
    memcpy(ctx->fw.tftpFilename, base, len);
    ctx->fw.tftpFilename[len] = '\0';

    if (version != NULL) {
        snprintf(ctx->fw.version, sizeof(ctx->fw.version), "%s", version);
    }
    if (sha256 != NULL) {
        snprintf(ctx->fw.sha256, sizeof(ctx->fw.sha256), "%s", sha256);
    }

    ctx->fw.size = st.st_size;
    ctx->fw.state = SWITCHD_FW_STAGED;
    ctx->fw.executeStatus = 1;
    ctx->fw.stagedAt = time(NULL);
    ctx->fw.updatedAt = time(NULL);

    op_end(ctx, SWITCHD_OK, "staged");
    return SWITCHD_OK;
}

int switchd_apply_firmware(SwitchdContext *ctx) {
    time_t deadline;
    int ret;

    if (ctx->fw.state != SWITCHD_FW_STAGED) {
        return SWITCHD_ERR_STATE;
    }

    ret = op_begin(ctx, SWITCHD_OP_FW_APPLY, 0, "firmware apply");
    if (ret != SWITCHD_OK) {
        return ret;
    }

    ctx->fw.state = SWITCHD_FW_APPLYING;
    ctx->fw.updatedAt = time(NULL);

    ret = tftp_server_start(&gTftp,
                            ctx->config.tftpBindIp,
                            ctx->config.tftpPort,
                            ctx->config.tftpRoot,
                            ctx->fw.tftpFilename);
    if (ret != SWITCHD_OK) {
        ctx->fw.state = SWITCHD_FW_FAILED;
        snprintf(ctx->fw.detail,
                 sizeof(ctx->fw.detail),
                 "failed to start TFTP server");
        op_end(ctx, ret, ctx->fw.detail);
        return ret;
    }

    pthread_mutex_lock(&ctx->driverMutex);
    ret = ctx->driver->ops.firmware_apply(ctx, &ctx->fw);
    pthread_mutex_unlock(&ctx->driverMutex);
    tftp_server_stop(&gTftp);

    if (ret != SWITCHD_OK) {
        ctx->fw.state = SWITCHD_FW_FAILED;
        snprintf(ctx->fw.detail,
                 sizeof(ctx->fw.detail),
                 "switch reported firmware transfer failure");
        op_end(ctx, ret, ctx->fw.detail);
        return ret;
    }

    ctx->fw.state = SWITCHD_FW_RECONNECTING;
    ctx->state = SWITCHD_STATE_RECOVERING;

    deadline = time(NULL) + ctx->config.firmwareReconnectSec;
    while (time(NULL) < deadline) {
        sleep(2);
        if (switchd_refresh_info(ctx) == SWITCHD_OK) {
            break;
        }
    }

    ctx->fw.state = SWITCHD_FW_VERIFYING;
    if (switchd_refresh_info(ctx) != SWITCHD_OK ||
        switchd_refresh_ports(ctx) != SWITCHD_OK ||
        switchd_refresh_kpis(ctx) != SWITCHD_OK) {
        ctx->fw.state = SWITCHD_FW_FAILED;
        snprintf(ctx->fw.detail,
                 sizeof(ctx->fw.detail),
                 "switch did not recover after update");
        op_end(ctx, SWITCHD_ERR_TIMEOUT, ctx->fw.detail);
        return SWITCHD_ERR_TIMEOUT;
    }

    ctx->fw.state = SWITCHD_FW_DONE;
    ctx->fw.executeStatus = 3;
    snprintf(ctx->fw.detail,
             sizeof(ctx->fw.detail),
             "firmware apply complete");
    op_end(ctx, SWITCHD_OK, ctx->fw.detail);
    return SWITCHD_OK;
}
