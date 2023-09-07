/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "collector.h"

#include "agents.h"
#include "config.h"
#include "log.h"
#include "server.h"

static int categoryCount = 0;
static MetricsConfig *metricsCfg = NULL;
bool collectionFlag = true;

agent_map_t agent_map[MAX_AGENTS] = {
    {.type = "sys_generic", .agentHandler = generic_stat_collector},
    {.type = "lte_agent", .agentHandler = lte_stack_collector},
    {.type = "rest_agent", .agentHandler = rest_collector},
    {.type = "sysfs_agent", .agentHandler = sysfs_collector},
    {.type = "cpu_agent", .agentHandler = cpu_collector},
    {.type = "memory_agent", .agentHandler = memory_collector},
    {.type = "network_agent", .agentHandler = network_collector},
    {.type = "ssd_agent", .agentHandler = ssd_collector},

};

CollectorFxn get_agent_handler_fxn(char *agent) {
  for (int idx = 0; idx < MAX_AGENTS; idx++) {
    if (strcmp(agent, agent_map[idx].type) == 0) {
      return agent_map[idx].agentHandler;
    }
  }
  return NULL;
}

int rest_collector(MetricsCatConfig *stat) {
  log_trace("Rest Agent started for source %s.", stat->source);

  return RETURN_OK;
}

int lte_stack_collector(MetricsCatConfig *stat) {
  log_trace("lte Agent started for source %s.", stat->source);
  int ret = fp_read_kpi_from_file(stat, metric_server_add_kpi_data);
  return ret;
}

int sysfs_collector(MetricsCatConfig *stat) {
  log_trace("sysfs Agent started for source %s.", stat->source);
  sysfs_collect_kpi(stat, metric_server_add_kpi_data);
  return RETURN_OK;
}

int cpu_collector(MetricsCatConfig *stat) {
  log_trace("CPU Agent started for source %s.", stat->source);
  int ret = sys_cpu_collect_stat(stat, metric_server_add_kpi_data);
  return ret;
}

int memory_collector(MetricsCatConfig *stat) {
  log_trace("Memory Agent started %s.", stat->source);
  int ret = sys_mem_collect_stat(stat, metric_server_add_kpi_data);
  return ret;
}

int network_collector(MetricsCatConfig *stat) {
  log_trace("Network Agent started for source %s.", stat->source);
  int ret = sys_net_collect_stat(stat, metric_server_add_kpi_data);
  return ret;
}

int ssd_collector(MetricsCatConfig *stat) {
  log_trace("SSD Agent started for source %s.", stat->source);
  int ret = sys_storage_collect_stat(stat, metric_server_add_kpi_data);
  return ret;
}

int generic_stat_collector(MetricsCatConfig *stat) {
  log_trace("Generic stat collection agent started for source %s.",
            stat->source);
  int ret = sys_gen_collect_stat(stat, metric_server_add_kpi_data);
  return ret;
}

int collector(char *cfg) {

  int ret = RETURN_OK;
  char *version;
  int scraping_time_period = 0;
  int server_port = 7001;

  /* Registry init */
  metric_server_registry_init();

  /* Parsing config */
  ret = toml_parse_config(cfg, &version, &scraping_time_period, &server_port,
                          &metricsCfg, &categoryCount);
  if (metricsCfg && (categoryCount <= 0)) {
    return RETURN_NOTOK;
  }
  MetricsCatConfig *stats = metricsCfg->metricsCategory;
  /* Active registry for HTTP Handler */
  metric_server_set_active_registry();

  /* Starting metric server */
  metric_server_start(server_port);

  /* Collect KPI metrics */
  while (collectionFlag) {

    for (int idx = 0; idx < categoryCount; idx++) {
      /* For each category */
      for (int cid = 0; cid < metricsCfg[idx].eachCategoryCount; cid++) {

        /* Start scraping metrics */
        CollectorFxn fxn =
            get_agent_handler_fxn(metricsCfg[idx].metricsCategory[cid].agent);
        if (fxn) {
          fxn(&metricsCfg[idx].metricsCategory[cid]);
        }
      }
    }

    sleep(scraping_time_period);
  }

  /* This means exit is called */
  log_info("Metrics: Clearing stats config.");
  if (stats) {
    free_stat_cfg(metricsCfg, categoryCount);
    metricsCfg = NULL;
    categoryCount = 0;
  }

  return ret;
}

/* Exit handler*/
void collector_exit(int sig) {
  log_info(
      " METRICS:: Signal %d caught Shutting down metrics collector client.");
  collectionFlag = false;
  metric_server_stop();
}
