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
#include "switchd.h"
#include "tftp_server.h"
#include "utils.h"

#include "usys_log.h"

SwitchdContext gSwitchd;
static TftpServer gTftp;

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

    ctx = (SwitchdContext *)arg;
    lastInfo = 0;
    lastKpis = 0;

    while (!ctx->terminate) {
        uint64_t now;

        now = monotonic_msec();
        (void)switchd_refresh_ports(ctx);

        if (now - lastKpis >= (uint64_t)ctx->config.pollKpisSec * 1000ULL) {
            (void)switchd_refresh_kpis(ctx);
            lastKpis = now;
        }

        if (now - lastInfo >= (uint64_t)ctx->config.pollInfoSec * 1000ULL) {
            (void)switchd_refresh_info(ctx);
            lastInfo = now;
        }

        sleep((unsigned int)ctx->config.pollStatusSec);
    }

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
    if (driver_init(ctx) != SWITCHD_OK) {
        return SWITCHD_ERR_INTERNAL;
    }

    if (ctx->driver->ops.probe(ctx) != SWITCHD_OK) {
        ctx->state = SWITCHD_STATE_DEGRADED;
        ctx->info.reachable = false;
    } else {
        (void)switchd_refresh_info(ctx);
        (void)switchd_refresh_ports(ctx);
        (void)switchd_refresh_kpis(ctx);
        ctx->state = SWITCHD_STATE_READY;
    }

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

    pthread_mutex_lock(&ctx->driverMutex);
    ret = ctx->driver->ops.refresh_info(ctx, &info, &caps);
    pthread_mutex_unlock(&ctx->driverMutex);
    if (ret != SWITCHD_OK) {
        ctx->info.reachable = false;
        ctx->state = SWITCHD_STATE_DEGRADED;
        return ret;
    }

    ctx->info = info;
    ctx->caps = caps;
    ctx->portCount = info.portCount;
    return SWITCHD_OK;
}

int switchd_refresh_ports(SwitchdContext *ctx) {
    SwitchPortState ports[SWITCHD_MAX_PORTS];
    uint32_t count;
    int ret;

    count = 0;
    pthread_mutex_lock(&ctx->driverMutex);
    ret = ctx->driver->ops.refresh_ports(ctx, ports, &count);
    pthread_mutex_unlock(&ctx->driverMutex);
    if (ret != SWITCHD_OK) {
        ctx->info.reachable = false;
        ctx->state = SWITCHD_STATE_DEGRADED;
        return ret;
    }

    memcpy(ctx->ports, ports, sizeof(ports));
    ctx->portCount = count;
    ctx->info.reachable = true;
    if (ctx->state == SWITCHD_STATE_DEGRADED) {
        ctx->state = SWITCHD_STATE_READY;
    }
    return SWITCHD_OK;
}

int switchd_refresh_kpis(SwitchdContext *ctx) {
    SwitchKpis kpis;
    int ret;

    pthread_mutex_lock(&ctx->driverMutex);
    ret = ctx->driver->ops.refresh_kpis(ctx, &kpis);
    pthread_mutex_unlock(&ctx->driverMutex);
    if (ret != SWITCHD_OK) {
        return ret;
    }

    ctx->kpis = kpis;
    return SWITCHD_OK;
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
