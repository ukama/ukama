/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_BFF_H_
#define ULAB_BFF_H_

#include <stdio.h>

#include "world.h"

typedef struct {
    char url[ULAB_MAX_PATH];
    FILE *logf;
} bff_client_t;

typedef struct {
    char state[ULAB_MAX_REF];
    char connectivity[ULAB_MAX_REF];
} bff_node_state_t;

int bff_init(bff_client_t *c,
             const char *url,
             const char *run_dir);

void bff_close(bff_client_t *c);

int bff_add_network(bff_client_t *c,
                    network_t *n,
                    ulab_error_t *err);

int bff_add_site(bff_client_t *c,
                 site_t *s,
                 const network_t *n,
                 ulab_error_t *err);

int bff_add_node(bff_client_t *c,
                 node_t *n,
                 ulab_error_t *err);

int bff_add_node_to_site(bff_client_t *c,
                         const node_t *n,
                         const site_t *s,
                         const network_t *net,
                         ulab_error_t *err);

int bff_add_package(bff_client_t *c,
                    package_t *p,
                    ulab_error_t *err);

int bff_add_subscriber(bff_client_t *c,
                       subscriber_t *sub,
                       const network_t *net,
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
