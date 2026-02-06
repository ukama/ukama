/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
int backhaul_collector(MetricsCatConfig *stat);

int collector(char *cfg);
void collector_exit();

#endif /* INC_COLLECTOR_H_ */
