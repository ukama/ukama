/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef INC_AGENTS_H_
#define INC_AGENTS_H_

#include "collector.h"
#include "config.h"
#include "metrics.h"

int fp_read_kpi_from_file(MetricsCatConfig *stat, metricAddFunc addFunc);
int sys_cpu_collect_stat(MetricsCatConfig *stat, metricAddFunc addFunc);
int sys_gen_collect_stat(MetricsCatConfig *stat, metricAddFunc addFunc);
int sys_mem_collect_stat(MetricsCatConfig *stat, metricAddFunc addFunc);
int sys_net_collect_stat(MetricsCatConfig *stat, metricAddFunc addFunc);
int sys_storage_collect_stat(MetricsCatConfig *stat, metricAddFunc addFunc);
int sysfs_collect_kpi(MetricsCatConfig *stat, metricAddFunc addFunc);
int backhaul_collect_stat(MetricsCatConfig *stat, metricAddFunc addFunc);

#endif /* INC_AGENTS_H_ */
