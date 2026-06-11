/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_SCENARIO_H_
#define ULAB_SCENARIO_H_

#include "ulab.h"

#define ULAB_MAX_PACKAGES 16
#define ULAB_MAX_PHASES   32
#define ULAB_MAX_EVENTS   64
#define ULAB_MAX_CHECKS   64
#define ULAB_MAX_BUCKETS  16
#define ULAB_MAX_LIST     32

typedef enum {
    SEL_NONE = 0,
    SEL_ALL,
    SEL_SAMPLE_PER_SITE,
    SEL_CREATED_IN_PHASE,
    SEL_AFFECTED_BY_PHASE,
    SEL_NODE_TYPE,
    SEL_NODE_TYPE_COUNT_PER_NETWORK
} selector_kind_t;

typedef struct {
    selector_kind_t kind;
    char            value[ULAB_MAX_REF];
    uint32_t        count;
} selector_t;

typedef struct {
    char     ref[ULAB_MAX_REF];
    char     name[ULAB_MAX_NAME];
    uint64_t data_mb;
    uint32_t duration_hours;
    double   amount;
    uint32_t assign_percent;
} package_spec_t;

typedef struct {
    char     name[ULAB_MAX_REF];
    uint32_t percent;
    uint64_t amount_mb;
} profile_bucket_t;

typedef struct {
    char             name[ULAB_MAX_REF];
    profile_bucket_t buckets[ULAB_MAX_BUCKETS];
    size_t           bucket_count;
} profile_spec_t;

typedef enum {
    EVT_TRAFFIC = 0,
    EVT_TRAFFIC_BY_PROFILE,
    EVT_CREATE_UES,
    EVT_START_UES,
    EVT_WAIT_UES_ATTACHED,
    EVT_RESTART_NODES,
    EVT_WAIT_NODES_READY,
    EVT_CHECK
} event_type_t;

typedef enum {
    CHECK_COUNT = 0,
    CHECK_NODE_READY,
    CHECK_UE_ATTACHED,
    CHECK_USAGE_PER_SIM,
    CHECK_USAGE_SAMPLE,
    CHECK_PACKAGE_ACTIVE,
    CHECK_PACKAGE_REMAINING,
    CHECK_NODE_STATE,
    CHECK_DASHBOARD_LOADS,
    CHECK_BALANCE_NON_NEGATIVE
} check_type_t;

typedef struct {
    check_type_t type;
    char         target[ULAB_MAX_REF];
    selector_t   ues;
    selector_t   nodes;
    selector_t   sites;
    selector_t   networks;
    char         package_ref[ULAB_MAX_REF];
    char         expected[ULAB_MAX_REF];
    uint64_t     expected_used_mb;
    uint64_t     expected_remaining_mb;
    uint32_t     tolerance_percent;
    uint32_t     required;
} check_spec_t;

typedef struct {
    event_type_t type;
    char         name[ULAB_MAX_NAME];
    selector_t   ues;
    selector_t   nodes;
    selector_t   sites;
    uint64_t     amount_mb;
    char         profile[ULAB_MAX_REF];
    uint32_t     count_per_site;
    char         package_ref[ULAB_MAX_REF];
    check_spec_t checks[ULAB_MAX_CHECKS];
    size_t       check_count;
} event_spec_t;

typedef struct {
    char         name[ULAB_MAX_NAME];
    event_spec_t events[ULAB_MAX_EVENTS];
    size_t       event_count;
    check_spec_t checks[ULAB_MAX_CHECKS];
    size_t       check_count;
} phase_spec_t;

typedef struct {
    uint32_t networks;
    uint32_t sites_per_network;
    uint32_t tower_per_site;
    uint32_t amplifier_per_site;
    uint32_t controller_per_site;
    uint32_t ues_per_site;
} world_spec_t;

typedef struct {
    int create_networks;
    int create_sites;
    int create_nodes;
    int create_node_site_links;
    int create_packages;
    int create_subscribers;
    int create_sims;
} setup_spec_t;

typedef struct {
    int start_nodes;
    int start_ues;
    int wait_nodes_ready;
    int wait_ues_attached;
} runtime_spec_t;

typedef struct {
    uint32_t       version;
    char           name[ULAB_MAX_NAME];
    uint32_t       seed;
    world_spec_t   world;
    package_spec_t packages[ULAB_MAX_PACKAGES];
    size_t         package_count;
    setup_spec_t   setup;
    runtime_spec_t runtime;
    profile_spec_t profiles[ULAB_MAX_BUCKETS];
    size_t         profile_count;
    phase_spec_t   phases[ULAB_MAX_PHASES];
    size_t         phase_count;
    check_spec_t   final_checks[ULAB_MAX_CHECKS];
    size_t         final_check_count;
} scenario_t;

int scenario_load(const char *path, scenario_t *s, ulab_error_t *err);
int scenario_validate(const scenario_t *s, ulab_error_t *err);
void scenario_init(scenario_t *s);
const char *scenario_event_name(event_type_t type);
const char *scenario_check_name(check_type_t type);
int scenario_event_from_name(const char *name, event_type_t *out);
int scenario_check_from_name(const char *name, check_type_t *out);
void scenario_list_events(void);
void scenario_list_checks(void);

#endif /* ULAB_SCENARIO_H_ */
