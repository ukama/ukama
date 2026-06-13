/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_WORLD_H_
#define ULAB_WORLD_H_

#include "scenario.h"

#define ULAB_NODE_TOWER      "tower"
#define ULAB_NODE_AMPLIFIER  "amplifier"
#define ULAB_NODE_CONTROLLER "controller"

#define ULAB_NODE_KIND_TOWER      "tnode"
#define ULAB_NODE_KIND_AMPLIFIER  "anode"
#define ULAB_NODE_KIND_CONTROLLER "cnode"

typedef struct {
    char ref[ULAB_MAX_REF];
    char id[ULAB_MAX_ID];
    char name[ULAB_MAX_NAME];
    char bff_id[ULAB_MAX_ID];
} network_t;

typedef struct {
    char ref[ULAB_MAX_REF];
    char id[ULAB_MAX_ID];
    char name[ULAB_MAX_NAME];
    char network_ref[ULAB_MAX_REF];
    char latitude[ULAB_MAX_REF];
    char longitude[ULAB_MAX_REF];
    char bff_id[ULAB_MAX_ID];
} site_t;

typedef struct {
    char ref[ULAB_MAX_REF];
    char id[ULAB_MAX_ID];
    char name[ULAB_MAX_NAME];
    char type[ULAB_MAX_REF];
    char site_ref[ULAB_MAX_REF];
    char network_ref[ULAB_MAX_REF];
    /*
     * Real NodeID selected by factory and used by runtime/BFF node registry,
     * e.g. uk-sa2602-tnode-v0-344c. id/ref remain logical lab identifiers.
     */
    char runtime_id[ULAB_MAX_ID];
    char latitude[ULAB_MAX_REF];
    char longitude[ULAB_MAX_REF];
    char bff_id[ULAB_MAX_ID];
} node_t;

typedef struct {
    char ref[ULAB_MAX_REF];
    char id[ULAB_MAX_ID];
    char name[ULAB_MAX_NAME];
    char email[ULAB_MAX_NAME];
    char phone[ULAB_MAX_REF];
    char network_ref[ULAB_MAX_REF];
    char site_ref[ULAB_MAX_REF];
    char bff_id[ULAB_MAX_ID];
} subscriber_t;

typedef struct {
    char ref[ULAB_MAX_REF];
    char id[ULAB_MAX_ID];
    char iccid[ULAB_MAX_ID];
    char imsi[ULAB_MAX_ID];
    char subscriber_ref[ULAB_MAX_REF];
    char network_ref[ULAB_MAX_REF];
    char site_ref[ULAB_MAX_REF];
    char package_ref[ULAB_MAX_REF];
    char bff_id[ULAB_MAX_ID];
    int  started;
    int  attached;
} ue_t;

typedef struct {
    char     ref[ULAB_MAX_REF];
    char     name[ULAB_MAX_NAME];
    uint64_t data_mb;
    uint32_t duration_days;
    double   amount;
    char     bff_id[ULAB_MAX_ID];
} package_t;

typedef struct {
    char         run_id[ULAB_MAX_ID];
    uint32_t     seed;
    network_t    *networks;
    size_t       network_count;
    site_t       *sites;
    size_t       site_count;
    node_t       *nodes;
    size_t       node_count;
    subscriber_t *subscribers;
    size_t       subscriber_count;
    ue_t         *ues;
    size_t       ue_count;
    package_t    *packages;
    size_t       package_count;
} world_t;

int world_generate(const scenario_t *s,
                   const char *run_id,
                   world_t *w,
                   ulab_error_t *err);
void world_free(world_t *w);
network_t *world_network_by_ref(world_t *w, const char *ref);
site_t *world_site_by_ref(world_t *w, const char *ref);
node_t *world_node_by_ref(world_t *w, const char *ref);
ue_t *world_ue_by_ref(world_t *w, const char *ref);
package_t *world_package_by_ref(world_t *w, const char *ref);
int world_add_ues(world_t *w,
                  const scenario_t *s,
                  const char *phase,
                  const selector_t *sites,
                  uint32_t count_per_site,
                  const char *package_ref,
                  ulab_error_t *err);
int world_write_json(const world_t *w, const char *path);
void world_print(const world_t *w);

#endif /* ULAB_WORLD_H_ */
