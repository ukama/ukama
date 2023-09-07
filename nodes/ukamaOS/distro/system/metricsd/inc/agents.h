/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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

#endif /* INC_AGENTS_H_ */
