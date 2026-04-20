/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <ctype.h>
#include <errno.h>
#include <net/if.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "config.h"
#include "metrics.h"

#include "usys_log.h"

#define PROC_NET_DEV    "/proc/net/dev"
#define PROC_NET_ROUTE  "/proc/net/route"

typedef struct {
    char interface[IFNAMSIZ];
    unsigned long long rxBytes;
    unsigned long long rxPackets;
    unsigned long long rxErrors;
    unsigned long long rxDropped;
    unsigned long long rxFifoErrors;
    unsigned long long rxOverruns;
    unsigned long long rxCompressed;
    unsigned long long multicast;
    unsigned long long txBytes;
    unsigned long long txPackets;
    unsigned long long txErrors;
    unsigned long long txDropped;
    unsigned long long txFifoErrors;
    unsigned long long collisions;
    unsigned long long txCarrierErrors;
    unsigned long long txCompressed;
    unsigned long long linkspeed;
    unsigned long long latency;
    unsigned long long linkstatus;
} SysNetDevMetrics;

static int sys_net_dev_read_stat(SysNetDevMetrics *netDev, char *ifName);
static int sys_net_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                              SysNetDevMetrics *nstat,
                                              metricAddFunc addFunc);

static bool is_auto_source(char *source) {

    if (source == NULL) {
        return false;
    }

    if (!strcmp(source, "uplink") || !strcmp(source, "auto")) {
        return true;
    }

    return false;
}

static bool is_valid_runtime_iface(char *ifName) {

    if ((ifName == NULL) || (*ifName == '\0')) {
        return false;
    }

    if (!strcmp(ifName, "lo")) {
        return false;
    }

    return true;
}

static int resolve_default_route_interface(char *ifName, size_t len) {

    FILE *fp                 = NULL;
    char line[256];
    int ret                  = RETURN_NOTOK;

    if ((ifName == NULL) || (len == 0)) {
        return RETURN_NOTOK;
    }

    fp = fopen(PROC_NET_ROUTE, "r");
    if (fp == NULL) {
        usys_log_error("failed to open %s: %s",
                       PROC_NET_ROUTE, strerror(errno));
        return RETURN_NOTOK;
    }

    while (fgets(line, sizeof(line), fp) != NULL) {
        char iface[IFNAMSIZ];
        unsigned long dest = 0;
        unsigned long gw   = 0;
        int flags          = 0;

        if (sscanf(line, "%15s %lx %lx %X",
                   iface, &dest, &gw, &flags) != 4) {
            continue;
        }

        if ((dest == 0) && is_valid_runtime_iface(iface)) {
            snprintf(ifName, len, "%s", iface);
            ret = RETURN_OK;
            break;
        }
    }

    fclose(fp);

    return ret;
}

static int resolve_first_non_loopback_interface(char *ifName, size_t len) {

    FILE *fp              = NULL;
    char line[256];
    int ret               = RETURN_NOTOK;

    if ((ifName == NULL) || (len == 0)) {
        return RETURN_NOTOK;
    }

    fp = fopen(PROC_NET_DEV, "r");
    if (fp == NULL) {
        usys_log_error("failed to open %s: %s",
                       PROC_NET_DEV, strerror(errno));
        return RETURN_NOTOK;
    }

    while (fgets(line, sizeof(line), fp) != NULL) {
        char *p     = line;
        char *colon = NULL;

        while (isspace((unsigned char)*p)) {
            p++;
        }

        colon = strchr(p, ':');
        if (colon == NULL) {
            continue;
        }

        *colon = '\0';

        if (!is_valid_runtime_iface(p)) {
            continue;
        }

        snprintf(ifName, len, "%s", p);
        ret = RETURN_OK;
        break;
    }

    fclose(fp);

    return ret;
}

static int resolve_net_interface(char *source, char *ifName, size_t len) {

    if ((ifName == NULL) || (len == 0)) {
        return RETURN_NOTOK;
    }

    memset(ifName, 0, len);

    if (!is_auto_source(source)) {
        snprintf(ifName, len, "%s", source);
        return RETURN_OK;
    }

    if (resolve_default_route_interface(ifName, len) == RETURN_OK) {
        usys_log_info("network source %s resolved to default-route interface %s",
                      source, ifName);
        return RETURN_OK;
    }

    if (resolve_first_non_loopback_interface(ifName, len) == RETURN_OK) {
        usys_log_warn("network source %s fallback to interface %s",
                      source, ifName);
        return RETURN_OK;
    }

    usys_log_error("failed to resolve runtime interface for source %s", source);

    return RETURN_NOTOK;
}

static int read_uint64_from_path(char *path, unsigned long long *val) {

    FILE *fp = NULL;

    if ((path == NULL) || (val == NULL)) {
        return RETURN_NOTOK;
    }

    fp = fopen(path, "r");
    if (fp == NULL) {
        return RETURN_NOTOK;
    }

    if (fscanf(fp, "%llu", val) != 1) {
        fclose(fp);
        return RETURN_NOTOK;
    }

    fclose(fp);

    return RETURN_OK;
}

static int read_operstate(char *ifName, unsigned long long *val) {

    char path[256];
    char state[32];
    FILE *fp = NULL;

    if ((ifName == NULL) || (val == NULL)) {
        return RETURN_NOTOK;
    }

    snprintf(path, sizeof(path), "/sys/class/net/%s/operstate", ifName);
    fp = fopen(path, "r");
    if (fp == NULL) {
        return RETURN_NOTOK;
    }

    if (fgets(state, sizeof(state), fp) == NULL) {
        fclose(fp);
        return RETURN_NOTOK;
    }

    fclose(fp);

    if (!strncmp(state, "up", 2)) {
        *val = 1;
    } else {
        *val = 0;
    }

    return RETURN_OK;
}

static int read_speed(char *ifName, unsigned long long *val) {

    char path[256];

    if ((ifName == NULL) || (val == NULL)) {
        return RETURN_NOTOK;
    }

    snprintf(path, sizeof(path), "/sys/class/net/%s/speed", ifName);

    return read_uint64_from_path(path, val);
}

static int sys_net_dev_read_stat(SysNetDevMetrics *netDev, char *ifName) {

    int ret      = RETURN_NOTOK;
    FILE *fp     = NULL;
    char line[256];

    if ((netDev == NULL) || (ifName == NULL)) {
        return RETURN_NOTOK;
    }

    fp = fopen(PROC_NET_DEV, "r");
    if (fp == NULL) {
        usys_log_error("failed to open %s: %s",
                       PROC_NET_DEV, strerror(errno));
        return RETURN_NOTOK;
    }

    while (fgets(line, sizeof(line), fp) != NULL) {
        char *p     = line;
        char *colon = NULL;

        while (isspace((unsigned char)*p)) {
            p++;
        }

        colon = strchr(p, ':');
        if (colon == NULL) {
            continue;
        }

        *colon = '\0';
        if (strcmp(p, ifName) != 0) {
            continue;
        }

        memset(netDev, 0, sizeof(SysNetDevMetrics));
        snprintf(netDev->interface, sizeof(netDev->interface), "%s", ifName);

        sscanf(colon + 1,
               "%llu %llu %llu %llu %llu %llu %llu %llu "
               "%llu %llu %llu %llu %llu %llu %llu %llu",
               &netDev->rxBytes, &netDev->rxPackets,
               &netDev->rxErrors, &netDev->rxDropped,
               &netDev->rxFifoErrors, &netDev->rxOverruns,
               &netDev->rxCompressed, &netDev->multicast,
               &netDev->txBytes, &netDev->txPackets,
               &netDev->txErrors, &netDev->txDropped,
               &netDev->txFifoErrors, &netDev->collisions,
               &netDev->txCarrierErrors, &netDev->txCompressed);

        (void)read_speed(ifName, &netDev->linkspeed);
        (void)read_operstate(ifName, &netDev->linkstatus);

        /*
         * Placeholder until latency logic is implemented.
         */
        netDev->latency = 0;

        ret = RETURN_OK;
        break;
    }

    fclose(fp);

    return ret;
}

static int sys_net_push_stat_to_metric_server(MetricsCatConfig *cfgStat,
                                              SysNetDevMetrics *nstat,
                                              metricAddFunc addFunc) {

    int idx                = 0;
    unsigned long long val = 0;

    if ((cfgStat == NULL) || (nstat == NULL) || (addFunc == NULL)) {
        return RETURN_NOTOK;
    }

    for (idx = 0; idx < cfgStat->kpiCount; idx++) {
        KPIConfig *kpi = &(cfgStat->kpi[idx]);

        if ((kpi == NULL) || (kpi->fqname == NULL)) {
            continue;
        }

        if (strstr(kpi->fqname, "rx_bytes")) {
            val = nstat->rxBytes;
        } else if (strstr(kpi->fqname, "rx_error")) {
            val = nstat->rxErrors;
        } else if (strstr(kpi->fqname, "rx_dropped")) {
            val = nstat->rxDropped;
        } else if (strstr(kpi->fqname, "rx_overruns")) {
            val = nstat->rxOverruns;
        } else if (strstr(kpi->fqname, "rx_packets")) {
            val = nstat->rxPackets;
        } else if (strstr(kpi->fqname, "tx_bytes")) {
            val = nstat->txBytes;
        } else if (strstr(kpi->fqname, "tx_error")) {
            val = nstat->txErrors;
        } else if (strstr(kpi->fqname, "tx_dropped")) {
            val = nstat->txDropped;
        } else if (strstr(kpi->fqname, "tx_overruns")) {
            val = nstat->txFifoErrors;
        } else if (strstr(kpi->fqname, "tx_packets")) {
            val = nstat->txPackets;
        } else if (strstr(kpi->fqname, "linkspeed")) {
            val = nstat->linkspeed;
        } else if (strstr(kpi->fqname, "link")) {
            val = nstat->linkstatus;
        } else if (strstr(kpi->fqname, "latency")) {
            val = nstat->latency;
        } else {
            continue;
        }

        {
            double castVal = (double)val;
            addFunc(kpi, &castVal);
        }
    }

    return RETURN_OK;
}

int sys_net_collect_stat(MetricsCatConfig *cfgStat, metricAddFunc addFunc) {

    SysNetDevMetrics netDev;
    char ifName[IFNAMSIZ];

    if ((cfgStat == NULL) || (addFunc == NULL) || (cfgStat->source == NULL)) {
        return RETURN_NOTOK;
    }

    if (resolve_net_interface(cfgStat->source, ifName, sizeof(ifName))
        != RETURN_OK) {
        usys_log_error("failed to resolve network source %s",
                       cfgStat->source);
        return RETURN_NOTOK;
    }

    if (sys_net_dev_read_stat(&netDev, ifName) != RETURN_OK) {
        usys_log_error("failed to read network stats for %s", ifName);
        return RETURN_NOTOK;
    }

    if (sys_net_push_stat_to_metric_server(cfgStat, &netDev, addFunc)
        != RETURN_OK) {
        usys_log_error("failed to push network stats for %s", ifName);
        return RETURN_NOTOK;
    }

    return RETURN_OK;
}
