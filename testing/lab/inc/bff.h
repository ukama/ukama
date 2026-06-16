/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_BFF_H_
#define ULAB_BFF_H_

#include <stdint.h>
#include <stdio.h>

#include "world.h"
#include "ulab.h"

typedef struct {
    char url[ULAB_MAX_URL];
    char pauth_url[ULAB_MAX_URL];
    char bff_base_url[ULAB_MAX_URL];
    char session_token[4096];
    char token[4096];
    int  authenticated;

    char access_id[ULAB_MAX_ID];
    char backhaul_id[ULAB_MAX_ID];
    char power_id[ULAB_MAX_ID];
    char spectrum_id[ULAB_MAX_ID];
    char switch_id[ULAB_MAX_ID];

    FILE *logf;
} bff_client_t;

typedef struct {
    char state[ULAB_MAX_REF];
    char connectivity[ULAB_MAX_REF];
} bff_node_state_t;

int bff_init(bff_client_t *c,
             const char *url,
             const char *run_dir);

int bff_login(bff_client_t *c,
              const char *identifier,
              const char *password,
              ulab_error_t *err);

void bff_close(bff_client_t *c);

int bff_add_network(bff_client_t *c,
                    network_t *n,
                    ulab_error_t *err);

int bff_wait_site_anchor_online(bff_client_t *c,
                                site_t *site,
                                ulab_error_t *err);

int bff_add_site(bff_client_t *c,
                 site_t *s,
                 const network_t *n,
                 ulab_error_t *err);

int bff_add_package(bff_client_t *c,
                    package_t *p,
                    ulab_error_t *err);

int bff_add_subscriber(bff_client_t *c,
                       subscriber_t *sub,
                       const network_t *net,
                       ulab_error_t *err);

int bff_upload_sims_from_csv(bff_client_t *c,
                             const char *csv_path,
                             const char *sim_type,
                             ulab_error_t *err);

int bff_get_sims_from_pool(bff_client_t *c,
                           const char *sim_type,
                           char iccids[][ULAB_MAX_ID],
                           size_t max_iccids,
                           size_t *iccid_count,
                           ulab_error_t *err);

int bff_allocate_sim_from_pool(bff_client_t *c,
                               ue_t *ue,
                               const subscriber_t *sub,
                               const network_t *net,
                               const package_t *pkg,
                               const char *sim_type,
                               ulab_error_t *err);

int bff_allocate_sim(bff_client_t *c,
                     ue_t *ue,
                     const subscriber_t *sub,
                     const network_t *net,
                     const package_t *pkg,
                     ulab_error_t *err);

int bff_get_sim_usage(bff_client_t *c,
                      const ue_t *ue,
                      uint64_t *used_mb,
                      ulab_error_t *err);

int bff_get_packages_for_sim(bff_client_t *c,
                             const ue_t *ue,
                             const char *package_id,
                             int *active,
                             ulab_error_t *err);

int bff_get_node_state(bff_client_t *c,
                       const node_t *node,
                       bff_node_state_t *state,
                       ulab_error_t *err);

int bff_network_overview_loads(bff_client_t *c,
                               const network_t *net,
                               ulab_error_t *err);

int bff_site_view_loads(bff_client_t *c,
                        const site_t *site,
                        ulab_error_t *err);

int bff_query_count(bff_client_t *c,
                    const char *target,
                    const world_t *w,
                    size_t *count,
                    ulab_error_t *err);

#endif /* ULAB_BFF_H_ */
