/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdbool.h>
#include <string.h>
#include <unistd.h>

#include "collector.h"

#include "agents.h"
#include "config.h"
#include "server.h"

#include "usys_file.h"
#include "usys_log.h"
#include "usys_services.h"

typedef int (*CollectorFxn)(MetricsCatConfig *stat);

typedef struct {
    char *type;
    CollectorFxn agentHandler;
} AgentMap;

static int categoryCount          = 0;
static MetricsConfig *metricsCfg  = NULL;
static bool collectionFlag        = true;

static AgentMap agentMap[] = {
    {.type = "sys_generic",       .agentHandler = generic_stat_collector},
    {.type = "generic_agent",     .agentHandler = generic_stat_collector},
    {.type = "lte_agent",         .agentHandler = lte_stack_collector},
    {.type = "rest_agent",        .agentHandler = rest_collector},
    {.type = "sysfs_agent",       .agentHandler = sysfs_collector},
    {.type = "thermal_agent",     .agentHandler = thermal_collector},
    {.type = "cpu_agent",         .agentHandler = cpu_collector},
    {.type = "memory_agent",      .agentHandler = memory_collector},
    {.type = "network_agent",     .agentHandler = network_collector},
    {.type = "ssd_agent",         .agentHandler = ssd_collector},
    {.type = "backhaul_agent",    .agentHandler = backhaul_collector},
    {.type = "femd_agent",        .agentHandler = femd_collector},
    {.type = "fem_agent",         .agentHandler = femd_collector},
    {.type = "controllerd_agent", .agentHandler = controllerd_collector},
    {.type = "controller_agent",  .agentHandler = controllerd_collector},
    {.type = "switchd_agent",     .agentHandler = switchd_collector},
    {.type = "switch_agent",      .agentHandler = switchd_collector},
    {.type = "power_agent",       .agentHandler = power_collector},
    {.type = "powerd_agent",      .agentHandler = power_collector},
};

static CollectorFxn get_agent_handler(char *agent) {

    size_t idx          = 0;
    size_t handlerCount = 0;

    if (agent == NULL) {
        return NULL;
    }

    handlerCount = sizeof(agentMap) / sizeof(agentMap[0]);

    for (idx = 0; idx < handlerCount; idx++) {
        if (agentMap[idx].type == NULL) {
            continue;
        }

        if (strcmp(agent, agentMap[idx].type) == 0) {
            return agentMap[idx].agentHandler;
        }
    }

    return NULL;
}

int rest_collector(MetricsCatConfig *stat) {

    usys_log_trace("rest agent started for source %s", stat->source);
    return RETURN_OK;
}

int lte_stack_collector(MetricsCatConfig *stat) {

    usys_log_trace("lte agent started for source %s", stat->source);
    return fp_read_kpi_from_file(stat, metric_server_add_kpi_data);
}

int sysfs_collector(MetricsCatConfig *stat) {

    usys_log_trace("sysfs agent started for source %s", stat->source);
    sysfs_collect_kpi(stat, metric_server_add_kpi_data);
    return RETURN_OK;
}

int cpu_collector(MetricsCatConfig *stat) {

    usys_log_trace("cpu agent started for source %s", stat->source);
    return sys_cpu_collect_stat(stat, metric_server_add_kpi_data);
}

int memory_collector(MetricsCatConfig *stat) {

    usys_log_trace("memory agent started for source %s", stat->source);
    return sys_mem_collect_stat(stat, metric_server_add_kpi_data);
}

int network_collector(MetricsCatConfig *stat) {

    usys_log_trace("network agent started for source %s", stat->source);
    return sys_net_collect_stat(stat, metric_server_add_kpi_data);
}

int ssd_collector(MetricsCatConfig *stat) {

    usys_log_trace("ssd agent started for source %s", stat->source);
    return sys_storage_collect_stat(stat, metric_server_add_kpi_data);
}

int generic_stat_collector(MetricsCatConfig *stat) {

    usys_log_trace("generic agent started for source %s", stat->source);
    return sys_gen_collect_stat(stat, metric_server_add_kpi_data);
}

int thermal_collector(MetricsCatConfig *stat) {

    usys_log_trace("thermal agent started for source %s", stat->source);
    return thermal_collect_stat(stat, metric_server_add_kpi_data);
}

int backhaul_collector(MetricsCatConfig *stat) {

    usys_log_trace("backhaul agent started for source %s", stat->source);
    return backhaul_collect_stat(stat, metric_server_add_kpi_data);
}

int femd_collector(MetricsCatConfig *stat) {

    usys_log_trace("femd agent started for source %s", stat->source);
    return femd_collect_stat(stat, metric_server_add_kpi_data);
}

int controllerd_collector(MetricsCatConfig *stat) {

    usys_log_trace("controllerd agent started for source %s",
                   stat->source);
    return controllerd_collect_stat(stat, metric_server_add_kpi_data);
}

int switchd_collector(MetricsCatConfig *stat) {

    usys_log_trace("switchd agent started for source %s", stat->source);
    return switchd_collect_stat(stat, metric_server_add_kpi_data);
}

int power_collector(MetricsCatConfig *stat) {

    usys_log_trace("power agent started for source %s", stat->source);
    return power_collect_stat(stat, metric_server_add_kpi_data);
}

int collector(char *cfg) {

    int ret                 = RETURN_OK;
    int scrapingTimePeriod  = 0;
    int serverPort          = 0;
    int idx                 = 0;
    int categoryIndex       = 0;
    char *version           = NULL;
    MetricsCatConfig *stats = NULL;

    metric_server_registry_init();

    ret = toml_parse_config(cfg, &version, &scrapingTimePeriod,
                            &metricsCfg, &categoryCount);
    if ((ret != RETURN_OK) || (metricsCfg == NULL) || (categoryCount <= 0)) {
        usys_log_error("failed to parse metrics config");
        metric_server_registry_destroy();
        free(version);
        return RETURN_NOTOK;
    }

    stats = metricsCfg->metricsCategory;

    metric_server_set_active_registry();

    serverPort = usys_find_service_port(SERVICE_METRICS);
    if (serverPort <= 0) {
        usys_log_error("unable to determine metrics port");
        free_stat_cfg(metricsCfg, categoryCount);
        metricsCfg = NULL;
        categoryCount = 0;
        free(version);
        metric_server_registry_destroy();
        return RETURN_NOTOK;
    }

    ret = metric_server_start(serverPort);
    if (ret != RETURN_OK) {
        usys_log_error("failed to start metrics exporter on port %d",
                       serverPort);
        free_stat_cfg(metricsCfg, categoryCount);
        metricsCfg = NULL;
        categoryCount = 0;
        free(version);
        metric_server_registry_destroy();
        return RETURN_NOTOK;
    }

    while (collectionFlag == true) {
        for (idx = 0; idx < categoryCount; idx++) {
            for (categoryIndex = 0;
                 categoryIndex < metricsCfg[idx].eachCategoryCount;
                 categoryIndex++) {
                MetricsCatConfig *categoryConfig = NULL;
                CollectorFxn handler            = NULL;

                categoryConfig =
                    &metricsCfg[idx].metricsCategory[categoryIndex];

                usys_log_debug("collector source=%s agent=%s url=%s",
                               categoryConfig->source,
                               categoryConfig->agent,
                               categoryConfig->url);

                handler = get_agent_handler(categoryConfig->agent);
                if (handler == NULL) {
                    usys_log_error("no handler for agent %s source %s",
                                   categoryConfig->agent,
                                   categoryConfig->source);
                    continue;
                }

                usys_log_debug("calling handler for agent %s source %s",
                               categoryConfig->agent,
                               categoryConfig->source);

                handler(categoryConfig);
            }
        }

        sleep(scrapingTimePeriod);
    }

    usys_log_info("clearing metrics config");

    if (stats != NULL) {
        free_stat_cfg(metricsCfg, categoryCount);
        metricsCfg = NULL;
        categoryCount = 0;
    }

    free(version);
    metric_server_stop();
    metric_server_registry_destroy();

    return ret;
}

void collector_exit(int signum) {

    (void)signum;

    usys_log_info("shutting down metrics collector");
    collectionFlag = false;
}
