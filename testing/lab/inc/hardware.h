/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc. */

#ifndef ULAB_HARDWARE_H_
#define ULAB_HARDWARE_H_

#include "world.h"

#define ULAB_HARDWARE_MAX_SITES 16u
#define ULAB_HARDWARE_MAX_UES   2u

typedef struct {
    char ref[ULAB_MAX_REF];
    char tnode_id[ULAB_MAX_ID];
    char anode_id[ULAB_MAX_ID];
    char cnode_id[ULAB_MAX_ID];
} hardware_site_t;

typedef struct {
    char ref[ULAB_MAX_REF];
    char site_ref[ULAB_MAX_REF];
    char iccid[ULAB_MAX_ID];
    char imsi[ULAB_MAX_ID];
    uint32_t slot;
} hardware_ue_t;

typedef struct {
    char host[ULAB_MAX_NAME];
    char user[ULAB_MAX_REF];
    char ssh_path[ULAB_MAX_PATH];
    uint32_t ssh_port;
    uint32_t command_timeout_seconds;
    uint32_t operation_timeout_seconds;
    char identity_file[ULAB_MAX_PATH];
    char known_hosts_file[ULAB_MAX_PATH];
    char media_host[ULAB_MAX_NAME];
    uint32_t iperf_ports[ULAB_HARDWARE_MAX_UES];
    hardware_site_t sites[ULAB_HARDWARE_MAX_SITES];
    size_t site_count;
    hardware_ue_t ues[ULAB_HARDWARE_MAX_UES];
    size_t ue_count;
} hardware_config_t;

int hardware_config_load(const char *path, hardware_config_t *config,
                          ulab_error_t *err);
int hardware_check_capacity(const scenario_t *scenario,
                             const hardware_config_t *config,
                             ulab_error_t *err);
int hardware_execution_supported(const scenario_t *scenario, ulab_error_t *err);
/* Apply physical identities before any backend or runtime operation. */
int hardware_bind_world(const hardware_config_t *config, world_t *world,
                         ulab_error_t *err);

#endif /* ULAB_HARDWARE_H_ */
