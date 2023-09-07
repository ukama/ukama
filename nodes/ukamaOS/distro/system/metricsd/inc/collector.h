/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_COLLECTOR_H_
#define INC_COLLECTOR_H_

#include "config.h"
#include "metrics.h"

#define MAX_AGENTS              8

typedef int (*CollectorFxn)(MetricsCatConfig *stat);

typedef struct {
    char *type;
    CollectorFxn agentHandler;
} agent_map_t;

int rest_collector(MetricsCatConfig *stat);
int lte_stack_collector(MetricsCatConfig *stat);
int sysfs_collector(MetricsCatConfig *stat);
int cpu_collector(MetricsCatConfig *stat);
int memory_collector(MetricsCatConfig *stat);
int network_collector(MetricsCatConfig *stat);
int ssd_collector(MetricsCatConfig *stat);
int generic_stat_collector(MetricsCatConfig *stat);

int collector(char *cfg);
void collector_exit();

#endif /* INC_COLLECTOR_H_ */
