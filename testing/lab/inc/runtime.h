/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_RUNTIME_H_
#define ULAB_RUNTIME_H_

#include <stdio.h>

#include "selector.h"

typedef struct {
    char provider[ULAB_MAX_REF];
    char script_dir[ULAB_MAX_PATH];
    char run_dir[ULAB_MAX_PATH];
    char repo[ULAB_MAX_PATH];
    FILE *logf;
    int service_enabled;
    int radio_enabled;
    int node_offline;
    int payment_failure_active;
    int software_failure_active;
    char node_version[ULAB_MAX_REF];
} runtime_t;

int runtime_init(runtime_t *rt, const char *provider,
                 const char *script_dir, const char *run_dir,
                 const char *repo);
void runtime_close(runtime_t *rt);
int runtime_ensure_network(runtime_t *rt, ulab_error_t *err);
int runtime_build_and_start_sites(const char *repo,
                                  runtime_t *rt,
                                  world_t *w,
                                  ulab_error_t *err);
int runtime_wait_nodes_ready(runtime_t *rt, const world_t *w,
                             const selector_result_t *nodes,
                             ulab_error_t *err);
int runtime_enable_pcrf_service(runtime_t *rt, const world_t *w,
                                ulab_error_t *err);
int runtime_ensure_media(runtime_t *rt, ulab_error_t *err);
int runtime_build_and_start_ues(const char *repo,
                                runtime_t *rt,
                                const world_t *w,
                                const selector_result_t *ues,
                                ulab_error_t *err);
int runtime_wait_ues_attached(runtime_t *rt, world_t *w,
                              const selector_result_t *ues,
                              ulab_error_t *err);
int runtime_generate_traffic(runtime_t *rt, const world_t *w,
                             const selector_result_t *ues,
                             uint64_t amount_mb, ulab_error_t *err);
int runtime_set_service(runtime_t *rt, int enabled, ulab_error_t *err);
int runtime_set_radio(runtime_t *rt, int enabled, ulab_error_t *err);
int runtime_mark_node_offline(runtime_t *rt, ulab_error_t *err);
int runtime_restore_nodes(runtime_t *rt, ulab_error_t *err);
int runtime_set_failure_control(runtime_t *rt,
                                const char *target,
                                int enabled,
                                ulab_error_t *err);
int runtime_restore_failure_controls(runtime_t *rt, ulab_error_t *err);
int runtime_set_node_version(runtime_t *rt, const char *version,
                             ulab_error_t *err);
int runtime_node_health_ok(const runtime_t *rt);
const char *runtime_node_version(const runtime_t *rt);
int runtime_restart_nodes(runtime_t *rt, const world_t *w,
                          const selector_result_t *nodes,
                          ulab_error_t *err);
int runtime_disconnect_nodes(runtime_t *rt, const world_t *w,
                             const selector_result_t *nodes,
                             ulab_error_t *err);
int runtime_reconnect_nodes(runtime_t *rt, const world_t *w,
                            const selector_result_t *nodes,
                            ulab_error_t *err);
int runtime_detach_ues(runtime_t *rt, const world_t *w, ulab_error_t *err);
int runtime_cleanup_ues(runtime_t *rt, const world_t *w, ulab_error_t *err);
int runtime_collect_cdr_diagnostics(runtime_t *rt, const world_t *w,
                                    ulab_error_t *err);
int runtime_collect_failure_logs(runtime_t *rt, const world_t *w,
                                 ulab_error_t *err);
int runtime_stop_ues(runtime_t *rt, const world_t *w, ulab_error_t *err);
int runtime_cleanup_infra(runtime_t *rt, const world_t *w, ulab_error_t *err);
int runtime_cleanup(runtime_t *rt, const world_t *w, ulab_error_t *err);


int runtime_virtual_build_and_start_sites(const char *repo,
                                          runtime_t *rt,
                                          world_t *w,
                                          ulab_error_t *err);
int runtime_virtual_wait_nodes_ready(runtime_t *rt, const world_t *w,
                                     const selector_result_t *nodes,
                                     ulab_error_t *err);
int runtime_virtual_restart_nodes(runtime_t *rt, const world_t *w,
                                  const selector_result_t *nodes,
                                  ulab_error_t *err);
int runtime_virtual_cleanup_infra(runtime_t *rt, const world_t *w,
                                  ulab_error_t *err);

#endif /* ULAB_RUNTIME_H_ */
