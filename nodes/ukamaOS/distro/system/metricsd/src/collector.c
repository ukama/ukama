/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include "collector.h"

#include "agents.h"
#include "config.h"
#include "server.h"

#include "usys_log.h"
#include "usys_file.h"
#include "usys_services.h"

static int categoryCount = 0;
static MetricsConfig *metricsCfg = NULL;
bool collectionFlag = true;

agent_map_t agent_map[MAX_AGENTS] = {
    {.type = "sys_generic",    .agentHandler = generic_stat_collector},
    {.type = "lte_agent",      .agentHandler = lte_stack_collector},
    {.type = "rest_agent",     .agentHandler = rest_collector},
    {.type = "sysfs_agent",    .agentHandler = sysfs_collector},
    {.type = "cpu_agent",      .agentHandler = cpu_collector},
    {.type = "memory_agent",   .agentHandler = memory_collector},
    {.type = "network_agent",  .agentHandler = network_collector},
    {.type = "ssd_agent",      .agentHandler = ssd_collector},
    {.type = "backhaul_agent", .agentHandler = backhaul_collector},
    {.type = "femd_agent",     .agentHandler = femd_collector},
};

CollectorFxn get_agent_handler_fxn(char *agent) {
  for (int idx = 0; idx < MAX_AGENTS; idx++) {
    if (!agent_map[idx].type) {
      continue;
    }
    if (strcmp(agent, agent_map[idx].type) == 0) {
      return agent_map[idx].agentHandler;
    }
  }
  return NULL;
}

int rest_collector(MetricsCatConfig *stat) {
    usys_log_trace("Rest Agent started for source %s.", stat->source);
    return RETURN_OK;
}

int lte_stack_collector(MetricsCatConfig *stat) {
  usys_log_trace("lte Agent started for source %s.", stat->source);
  return fp_read_kpi_from_file(stat, metric_server_add_kpi_data);
}

int sysfs_collector(MetricsCatConfig *stat) {
  usys_log_trace("sysfs Agent started for source %s.", stat->source);
  sysfs_collect_kpi(stat, metric_server_add_kpi_data);
  return RETURN_OK;
}

int cpu_collector(MetricsCatConfig *stat) {
  usys_log_trace("CPU Agent started for source %s.", stat->source);
  return sys_cpu_collect_stat(stat, metric_server_add_kpi_data);
}

int memory_collector(MetricsCatConfig *stat) {
  usys_log_trace("Memory Agent started %s.", stat->source);
  return sys_mem_collect_stat(stat, metric_server_add_kpi_data);
}

int network_collector(MetricsCatConfig *stat) {
  usys_log_trace("Network Agent started for source %s.", stat->source);
  return sys_net_collect_stat(stat, metric_server_add_kpi_data);
}

int ssd_collector(MetricsCatConfig *stat) {
  usys_log_trace("SSD Agent started for source %s.", stat->source);
  return sys_storage_collect_stat(stat, metric_server_add_kpi_data);
}

int generic_stat_collector(MetricsCatConfig *stat) {
  usys_log_trace("Generic stat collection agent started for source %s.",
                 stat->source);
  return sys_gen_collect_stat(stat, metric_server_add_kpi_data);
}

int backhaul_collector(MetricsCatConfig *stat) {
  usys_log_trace("Backhaul Agent started for source %s.", stat->source);
  return backhaul_collect_stat(stat, metric_server_add_kpi_data);
}

int femd_collector(MetricsCatConfig *stat) {
    usys_log_trace("Femd Agent started for source %s", stat->source);
    return femd_collect_stat(stat, metric_server_add_kpi_data);
}

int collector(char *cfg) {

  int ret = RETURN_OK;
  char *version;
  int scraping_time_period = 0;
  int server_port;

  /* Registry init */
  metric_server_registry_init();

  server_port = usys_find_service_port(SERVICE_METRICS);

  /* Parsing config */
  ret = toml_parse_config(cfg, &version, &scraping_time_period,
                          &metricsCfg, &categoryCount);
  if (metricsCfg && (categoryCount <= 0)) {
    return RETURN_NOTOK;
  }
  MetricsCatConfig *stats = metricsCfg->metricsCategory;
  /* Active registry for HTTP Handler */
  metric_server_set_active_registry();

  /* Starting metric server */
  server_port = usys_find_service_port(SERVICE_METRICS);
  metric_server_start(server_port);

  /* Collect KPI metrics */
  while (collectionFlag) {
      for (int idx = 0; idx < categoryCount; idx++) {
          /* For each category */
          for (int cid = 0; cid < metricsCfg[idx].eachCategoryCount; cid++) {
              usys_log_debug("collector: src=%s agent='%s' url=%s",
                             metricsCfg[idx].metricsCategory[cid].source,
                             metricsCfg[idx].metricsCategory[cid].agent,
                             metricsCfg[idx].metricsCategory[cid].url);
              CollectorFxn fxn =
                  get_agent_handler_fxn(metricsCfg[idx].metricsCategory[cid].agent);

              if (!fxn) {
                  usys_log_error("collector: NO HANDLER for agent='%s' src=%s",
                                 metricsCfg[idx].metricsCategory[cid].agent,
                                 metricsCfg[idx].metricsCategory[cid].source);
              } else {
                  usys_log_debug("collector: calling handler for agent='%s' src=%s",
                                 metricsCfg[idx].metricsCategory[cid].agent,
                                 metricsCfg[idx].metricsCategory[cid].source);
                  fxn(&metricsCfg[idx].metricsCategory[cid]);
              }
          }
      }

      sleep(scraping_time_period);
  }

  /* This means exit is called */
  usys_log_info("Clearing stats config.");
  if (stats) {
    free_stat_cfg(metricsCfg, categoryCount);
    metricsCfg = NULL;
    categoryCount = 0;
  }

  if (version) {
      free(version);
  }

  return ret;
}

/* Exit handler*/
void collector_exit(int sig) {
  usys_log_info("Signal %d caught Shutting down metrics collector client.");
  collectionFlag = false;
  metric_server_stop();
}
