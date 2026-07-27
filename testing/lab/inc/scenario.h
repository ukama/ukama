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
    SEL_REF,
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
    uint32_t duration_days;
    uint32_t duration_minutes;
    double   amount;
    char     currency[ULAB_MAX_REF];
    char     country[ULAB_MAX_REF];
    char     scope[ULAB_MAX_REF];
    char     network_ref[ULAB_MAX_REF];
    int      active;
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
    EVT_WAIT,
    EVT_RESTART_NODES,
    EVT_WAIT_NODE_CONNECTIVITY,
    EVT_WAIT_NODES_READY,
    EVT_ADD_PACKAGE_TO_SIM,
    EVT_PURCHASE_PACKAGE,
    EVT_PURCHASE_PACKAGES_PARALLEL,
    EVT_ALLOCATE_SIM,
    EVT_CREATE_INVALID_PACKAGE,
    EVT_WAIT_PACKAGE_BOUNDARY,
    EVT_SET_PACKAGE_ACTIVE,
    EVT_REMOVE_PACKAGE_FROM_SIM,
    EVT_SET_SIM_STATUS,
    EVT_TOGGLE_SERVICE,
    EVT_TOGGLE_RADIO,
    EVT_RESTART_SITE,
    EVT_PROMOTE_RELEASE,
    EVT_SOFTWARE_UPDATE,
    EVT_DISCONNECT_NODES,
    EVT_RECONNECT_NODES,
    EVT_MARK_NODE_OFFLINE,
    EVT_RESTORE_NODE,
    EVT_CHECK
} event_type_t;

typedef enum {
    CHECK_BACKEND_COUNT = 0,
    CHECK_LIST_CONTAINS,
    CHECK_LIST_EXCLUDES,
    CHECK_STATUS_EQUALS,
    CHECK_TRAFFIC_ALLOWED,
    CHECK_TRAFFIC_BLOCKED,
    CHECK_NODE_READY,
    CHECK_UE_ATTACHED,
    CHECK_USAGE_PER_SIM,
    CHECK_USAGE_SAMPLE,
    CHECK_PACKAGE_ACTIVE,
    CHECK_PACKAGE_REMAINING,
    CHECK_PACKAGE_STATE,
    CHECK_PACKAGE_ASSIGNMENT_COUNT,
    CHECK_PACKAGE_ASSIGNMENT_CHAIN,
    CHECK_PACKAGE_CATALOG_EQUALS,
    CHECK_PACKAGE_VISIBLE,
    CHECK_PACKAGE_HIDDEN,
    CHECK_PACKAGE_NAME_AVAILABLE,
    CHECK_PACKAGE_BUSINESS_METRICS,
    CHECK_SIM_UNALLOCATED,
    CHECK_PAYMENT_EQUALS,
    CHECK_PAYMENT_COUNT,
    CHECK_KPI_VALUE,
    CHECK_KPI_TREND,
    CHECK_KPI_CONTRACT,
    CHECK_KPI_ROLLUP_CONSISTENCY,
    CHECK_PERFORMANCE_REPORT_CELL,
    CHECK_PERFORMANCE_REPORT_ROW,
    CHECK_REVENUE_SUMMARY,
    CHECK_SUBSCRIBER_BILLING_SUMMARY,
    CHECK_PAYMENT_ENTITLEMENT_RECONCILES,
    CHECK_PACKAGE_DASHBOARD_METRIC,
    CHECK_NETWORK_OVERVIEW_METRIC,
    CHECK_CONSOLE_INVENTORY_RECONCILES,
    CHECK_USAGE_AGGREGATE,
    CHECK_NODE_STATE,
    CHECK_DASHBOARD_LOADS,
    CHECK_DASHBOARD_SECTION_OK,
    CHECK_NODE_VERSION_EQUALS,
    CHECK_NODE_HEALTH_OK,
    CHECK_RELEASE_UNAVAILABLE,
    CHECK_HISTORY_PRESERVED,
    CHECK_AUDIT_EVENT_EXISTS,
    CHECK_RELATIONSHIP_EXISTS,
    CHECK_RELATIONSHIP_ENDED,
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
    char         other_package_ref[ULAB_MAX_REF];
    char         variant[ULAB_MAX_REF];
    char         view[ULAB_MAX_REF];
    char         ref[ULAB_MAX_REF];
    char         entity[ULAB_MAX_REF];
    char         status[ULAB_MAX_REF];
    char         expected[ULAB_MAX_REF];
    char         app[ULAB_MAX_NAME];
    char         key[ULAB_MAX_REF];
    char         span[ULAB_MAX_REF];
    char         op[ULAB_MAX_REF];
    char         comparator[ULAB_MAX_REF];
    char         report[ULAB_MAX_REF];
    char         column[ULAB_MAX_REF];
    char         scope_key[ULAB_MAX_REF];
    char         scope_value[ULAB_MAX_ID];
    char         trend_direction[ULAB_MAX_REF];
    char         currency[ULAB_MAX_REF];
    char         payment_method[ULAB_MAX_REF];
    double       expected_value;
    double       tolerance_value;
    uint32_t     expected_count;
    int          has_expected_count;
    uint32_t     timeout_seconds;
    uint32_t     poll_seconds;
    int          has_expected_value;
    int          expected_partial;
    int          has_expected_partial;
    int          require_computed_at;
    int          require_scope;
    int          require_trend_consistency;
    uint64_t     expected_used_mb;
    uint64_t     expected_remaining_mb;
    uint32_t     tolerance_percent;
    uint32_t     required;
    int          immediate;
} check_spec_t;

typedef struct {
    event_type_t type;
    char         name[ULAB_MAX_NAME];
    selector_t   ues;
    selector_t   nodes;
    selector_t   sites;
    uint64_t     amount_mb;
    int64_t      offset_seconds;
    char         profile[ULAB_MAX_REF];
    char         expect_result[ULAB_MAX_REF];
    char         error_contains[ULAB_MAX_ERR];
    char         status[ULAB_MAX_REF];
    char         app[ULAB_MAX_NAME];
    char         tag[ULAB_MAX_REF];
    char         currency[ULAB_MAX_REF];
    char         payer_email[ULAB_MAX_NAME];
    char         payer_phone[ULAB_MAX_REF];
    char         idempotency_key[ULAB_MAX_ID];
    char         other_package_ref[ULAB_MAX_REF];
    char         variant[ULAB_MAX_REF];
    double       amount;
    int          has_amount;
    int          active;
    int          has_active;
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
    uint32_t sims_per_subscriber;
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
    char type[ULAB_MAX_REF];
} provider_spec_t;

typedef struct {
    int start_nodes;
    int start_ues;
    int wait_nodes_ready;
    int wait_ues_attached;
} runtime_spec_t;

typedef struct {
    uint32_t       version;
    char           name[ULAB_MAX_NAME];
    char           description[ULAB_MAX_LINE];
    uint32_t       seed;
    char           suite[ULAB_MAX_REF];
    char           priority[ULAB_MAX_REF];
    char           tags[ULAB_MAX_LINE];
    char           status[ULAB_MAX_REF];
    int            generated;
    char           entity[ULAB_MAX_REF];
    char           action[ULAB_MAX_REF];
    world_spec_t   world;
    package_spec_t packages[ULAB_MAX_PACKAGES];
    size_t         package_count;
    setup_spec_t   setup;
    provider_spec_t provider;
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
