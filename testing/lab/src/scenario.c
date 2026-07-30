/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "scenario.h"
#include "log.h"
#include "util.h"


static int expand_environment(const char *input,
                              char *output,
                              size_t output_len,
                              char *missing,
                              size_t missing_len) {
    size_t in_pos;
    size_t out_pos;

    if (input == NULL || output == NULL || output_len == 0) {
        return ULAB_ERR;
    }

    in_pos = 0;
    out_pos = 0;
    while (input[in_pos] != '\0') {
        if (input[in_pos] == '$' && input[in_pos + 1] == '{') {
            char name[ULAB_MAX_REF];
            const char *value;
            size_t name_len;
            size_t value_len;
            size_t end;

            end = in_pos + 2;
            while (input[end] != '\0' && input[end] != '}') {
                end++;
            }
            if (input[end] != '}') {
                return ULAB_ERR;
            }

            name_len = end - (in_pos + 2);
            if (name_len == 0 || name_len >= sizeof(name)) {
                return ULAB_ERR;
            }
            memcpy(name, input + in_pos + 2, name_len);
            name[name_len] = '\0';

            value = getenv(name);
            if (value == NULL) {
                if (missing != NULL && missing_len > 0) {
                    ulab_copy(missing, missing_len, name);
                }
                return ULAB_ERR;
            }

            value_len = strlen(value);
            if (out_pos + value_len + 1 > output_len) {
                return ULAB_ERR;
            }
            memcpy(output + out_pos, value, value_len);
            out_pos += value_len;
            in_pos = end + 1;
            continue;
        }

        if (out_pos + 2 > output_len) {
            return ULAB_ERR;
        }
        output[out_pos++] = input[in_pos++];
    }

    output[out_pos] = '\0';
    return ULAB_OK;
}

typedef enum {
    SEC_NONE = 0,
    SEC_WORLD,
    SEC_NODES_PER_SITE,
    SEC_PACKAGES,
    SEC_SETUP,
    SEC_SETUP_LIST,
    SEC_PROVIDER,
    SEC_RUNTIME,
    SEC_PROFILES,
    SEC_PROFILE_ONE,
    SEC_PROFILE_BUCKET,
    SEC_PHASES,
    SEC_PHASE_EVENTS,
    SEC_EVENT_EXPECT,
    SEC_PHASE_CHECKS,
    SEC_FINAL_CHECKS
} parse_sec_t;

void scenario_init(scenario_t *s) {
    memset(s, 0, sizeof(*s));
    snprintf(s->suite, sizeof(s->suite), "default");
    snprintf(s->priority, sizeof(s->priority), "p2");
    snprintf(s->status, sizeof(s->status), "active");
    snprintf(s->provider.type, sizeof(s->provider.type), "virtual");
}

const char *scenario_event_name(event_type_t type) {
    switch (type) {
    case EVT_TRAFFIC: return "traffic";
    case EVT_TRAFFIC_BY_PROFILE: return "traffic_by_profile";
    case EVT_CREATE_UES: return "create_ues";
    case EVT_START_UES: return "start_ues";
    case EVT_WAIT_UES_ATTACHED: return "wait_ues_attached";
    case EVT_WAIT: return "wait";
    case EVT_RESTART_NODES: return "restart_nodes";
    case EVT_WAIT_NODE_CONNECTIVITY: return "wait_node_connectivity";
    case EVT_WAIT_NODES_READY: return "wait_nodes_ready";
    case EVT_ADD_PACKAGE_TO_SIM: return "add_package_to_sim";
    case EVT_PURCHASE_PACKAGE: return "purchase_package";
    case EVT_PURCHASE_PACKAGES_PARALLEL:
        return "purchase_packages_parallel";
    case EVT_ALLOCATE_SIM: return "allocate_sim";
    case EVT_CREATE_INVALID_PACKAGE: return "create_invalid_package";
    case EVT_WAIT_PACKAGE_BOUNDARY: return "wait_package_boundary";
    case EVT_SET_PACKAGE_ACTIVE: return "set_package_active";
    case EVT_REMOVE_PACKAGE_FROM_SIM: return "remove_package_from_sim";
    case EVT_SET_SIM_STATUS: return "set_sim_status";
    case EVT_TOGGLE_SERVICE: return "toggle_service";
    case EVT_TOGGLE_RADIO: return "toggle_radio";
    case EVT_RESTART_SITE: return "restart_site";
    case EVT_PROMOTE_RELEASE: return "promote_release";
    case EVT_SOFTWARE_UPDATE: return "software_update";
    case EVT_DISCONNECT_NODES: return "disconnect_nodes";
    case EVT_RECONNECT_NODES: return "reconnect_nodes";
    case EVT_MARK_NODE_OFFLINE: return "mark_node_offline";
    case EVT_RESTORE_NODE: return "restore_node";
    case EVT_CHECK: return "check";
    default: return "unknown";
    }
}

const char *scenario_check_name(check_type_t type) {
    switch (type) {
    case CHECK_BACKEND_COUNT: return "backend_count";
    case CHECK_LIST_CONTAINS: return "list_contains";
    case CHECK_LIST_EXCLUDES: return "list_excludes";
    case CHECK_STATUS_EQUALS: return "status_equals";
    case CHECK_TRAFFIC_ALLOWED: return "traffic_allowed";
    case CHECK_TRAFFIC_BLOCKED: return "traffic_blocked";
    case CHECK_NODE_READY: return "node_ready";
    case CHECK_UE_ATTACHED: return "ue_attached";
    case CHECK_USAGE_PER_SIM: return "usage_per_sim";
    case CHECK_USAGE_SAMPLE: return "usage_sample";
    case CHECK_PACKAGE_ACTIVE: return "package_active";
    case CHECK_PACKAGE_REMAINING: return "package_remaining";
    case CHECK_PACKAGE_STATE: return "package_state";
    case CHECK_PACKAGE_ASSIGNMENT_COUNT: return "package_assignment_count";
    case CHECK_PACKAGE_ASSIGNMENT_CHAIN: return "package_assignment_chain";
    case CHECK_PACKAGE_CATALOG_EQUALS: return "package_catalog_equals";
    case CHECK_PACKAGE_VISIBLE: return "package_visible";
    case CHECK_PACKAGE_HIDDEN: return "package_hidden";
    case CHECK_PACKAGE_NAME_AVAILABLE: return "package_name_available";
    case CHECK_PACKAGE_BUSINESS_METRICS:
        return "package_business_metrics";
    case CHECK_SIM_UNALLOCATED: return "sim_unallocated";
    case CHECK_PAYMENT_EQUALS: return "payment_equals";
    case CHECK_PAYMENT_COUNT: return "payment_count";
    case CHECK_KPI_VALUE: return "kpi_value";
    case CHECK_KPI_TREND: return "kpi_trend";
    case CHECK_KPI_CONTRACT: return "kpi_contract";
    case CHECK_KPI_ROLLUP_CONSISTENCY:
        return "kpi_rollup_consistency";
    case CHECK_PERFORMANCE_REPORT_CELL: return "performance_report_cell";
    case CHECK_PERFORMANCE_REPORT_ROW: return "performance_report_row";
    case CHECK_REVENUE_SUMMARY: return "revenue_summary";
    case CHECK_SUBSCRIBER_BILLING_SUMMARY:
        return "subscriber_billing_summary";
    case CHECK_PAYMENT_ENTITLEMENT_RECONCILES:
        return "payment_entitlement_reconciles";
    case CHECK_PACKAGE_DASHBOARD_METRIC:
        return "package_dashboard_metric";
    case CHECK_NETWORK_OVERVIEW_METRIC:
        return "network_overview_metric";
    case CHECK_CONSOLE_INVENTORY_RECONCILES:
        return "console_inventory_reconciles";
    case CHECK_USAGE_AGGREGATE: return "usage_aggregate";
    case CHECK_NODE_STATE: return "node_state";
    case CHECK_DASHBOARD_LOADS: return "dashboard_loads";
    case CHECK_DASHBOARD_SECTION_OK: return "dashboard_section_ok";
    case CHECK_NODE_VERSION_EQUALS: return "node_version_equals";
    case CHECK_NODE_HEALTH_OK: return "node_health_ok";
    case CHECK_RELEASE_UNAVAILABLE: return "release_unavailable";
    case CHECK_HISTORY_PRESERVED: return "history_preserved";
    case CHECK_AUDIT_EVENT_EXISTS: return "audit_event_exists";
    case CHECK_RELATIONSHIP_EXISTS: return "relationship_exists";
    case CHECK_RELATIONSHIP_ENDED: return "relationship_ended";
    case CHECK_BALANCE_NON_NEGATIVE: return "balance_non_negative";
    default: return "unknown";
    }
}

int scenario_event_from_name(const char *name, event_type_t *out) {
    if (ulab_streq(name, "traffic")) *out = EVT_TRAFFIC;
    else if (ulab_streq(name, "traffic_by_profile")) {
        *out = EVT_TRAFFIC_BY_PROFILE;
    } else if (ulab_streq(name, "create_ues")) *out = EVT_CREATE_UES;
    else if (ulab_streq(name, "start_ues")) *out = EVT_START_UES;
    else if (ulab_streq(name, "wait_ues_attached")) {
        *out = EVT_WAIT_UES_ATTACHED;
    } else if (ulab_streq(name, "wait")) {
        *out = EVT_WAIT;
    } else if (ulab_streq(name, "restart_nodes")) {
        *out = EVT_RESTART_NODES;
    } else if (ulab_streq(name, "wait_node_connectivity")) {
        *out = EVT_WAIT_NODE_CONNECTIVITY;
    } else if (ulab_streq(name, "wait_nodes_ready")) {
        *out = EVT_WAIT_NODES_READY;
    } else if (ulab_streq(name, "add_package_to_sim")) {
        *out = EVT_ADD_PACKAGE_TO_SIM;
    } else if (ulab_streq(name, "purchase_package")) {
        *out = EVT_PURCHASE_PACKAGE;
    } else if (ulab_streq(name, "purchase_packages_parallel")) {
        *out = EVT_PURCHASE_PACKAGES_PARALLEL;
    } else if (ulab_streq(name, "allocate_sim")) {
        *out = EVT_ALLOCATE_SIM;
    } else if (ulab_streq(name, "create_invalid_package")) {
        *out = EVT_CREATE_INVALID_PACKAGE;
    } else if (ulab_streq(name, "wait_package_boundary")) {
        *out = EVT_WAIT_PACKAGE_BOUNDARY;
    } else if (ulab_streq(name, "set_package_active")) {
        *out = EVT_SET_PACKAGE_ACTIVE;
    } else if (ulab_streq(name, "remove_package_from_sim")) {
        *out = EVT_REMOVE_PACKAGE_FROM_SIM;
    } else if (ulab_streq(name, "set_sim_status")) {
        *out = EVT_SET_SIM_STATUS;
    } else if (ulab_streq(name, "toggle_service")) {
        *out = EVT_TOGGLE_SERVICE;
    } else if (ulab_streq(name, "toggle_radio")) {
        *out = EVT_TOGGLE_RADIO;
    } else if (ulab_streq(name, "restart_site")) {
        *out = EVT_RESTART_SITE;
    } else if (ulab_streq(name, "promote_release")) {
        *out = EVT_PROMOTE_RELEASE;
    } else if (ulab_streq(name, "software_update")) {
        *out = EVT_SOFTWARE_UPDATE;
    } else if (ulab_streq(name, "disconnect_nodes")) {
        *out = EVT_DISCONNECT_NODES;
    } else if (ulab_streq(name, "reconnect_nodes")) {
        *out = EVT_RECONNECT_NODES;
    } else if (ulab_streq(name, "mark_node_offline")) {
        *out = EVT_MARK_NODE_OFFLINE;
    } else if (ulab_streq(name, "restore_node")) {
        *out = EVT_RESTORE_NODE;
    } else if (ulab_streq(name, "check")) *out = EVT_CHECK;
    else return ULAB_ERR;
    return ULAB_OK;
}

int scenario_check_from_name(const char *name, check_type_t *out) {
    if (ulab_streq(name, "count") || ulab_streq(name, "backend_count")) {
        *out = CHECK_BACKEND_COUNT;
    } else if (ulab_streq(name, "list_contains")) {
        *out = CHECK_LIST_CONTAINS;
    } else if (ulab_streq(name, "list_excludes")) {
        *out = CHECK_LIST_EXCLUDES;
    } else if (ulab_streq(name, "status_equals")) {
        *out = CHECK_STATUS_EQUALS;
    } else if (ulab_streq(name, "traffic_allowed")) {
        *out = CHECK_TRAFFIC_ALLOWED;
    } else if (ulab_streq(name, "traffic_blocked")) {
        *out = CHECK_TRAFFIC_BLOCKED;
    } else if (ulab_streq(name, "node_ready")) *out = CHECK_NODE_READY;
    else if (ulab_streq(name, "ue_attached")) *out = CHECK_UE_ATTACHED;
    else if (ulab_streq(name, "usage_per_sim")) {
        *out = CHECK_USAGE_PER_SIM;
    } else if (ulab_streq(name, "usage_sample")) {
        *out = CHECK_USAGE_SAMPLE;
    } else if (ulab_streq(name, "package_active")) {
        *out = CHECK_PACKAGE_ACTIVE;
    } else if (ulab_streq(name, "package_remaining")) {
        *out = CHECK_PACKAGE_REMAINING;
    } else if (ulab_streq(name, "package_state")) {
        *out = CHECK_PACKAGE_STATE;
    } else if (ulab_streq(name, "package_assignment_count")) {
        *out = CHECK_PACKAGE_ASSIGNMENT_COUNT;
    } else if (ulab_streq(name, "package_assignment_chain")) {
        *out = CHECK_PACKAGE_ASSIGNMENT_CHAIN;
    } else if (ulab_streq(name, "package_catalog_equals")) {
        *out = CHECK_PACKAGE_CATALOG_EQUALS;
    } else if (ulab_streq(name, "package_visible")) {
        *out = CHECK_PACKAGE_VISIBLE;
    } else if (ulab_streq(name, "package_hidden")) {
        *out = CHECK_PACKAGE_HIDDEN;
    } else if (ulab_streq(name, "package_name_available")) {
        *out = CHECK_PACKAGE_NAME_AVAILABLE;
    } else if (ulab_streq(name, "package_business_metrics")) {
        *out = CHECK_PACKAGE_BUSINESS_METRICS;
    } else if (ulab_streq(name, "sim_unallocated")) {
        *out = CHECK_SIM_UNALLOCATED;
    } else if (ulab_streq(name, "payment_equals")) {
        *out = CHECK_PAYMENT_EQUALS;
    } else if (ulab_streq(name, "payment_count")) {
        *out = CHECK_PAYMENT_COUNT;
    } else if (ulab_streq(name, "kpi_value")) {
        *out = CHECK_KPI_VALUE;
    } else if (ulab_streq(name, "kpi_trend")) {
        *out = CHECK_KPI_TREND;
    } else if (ulab_streq(name, "kpi_contract")) {
        *out = CHECK_KPI_CONTRACT;
    } else if (ulab_streq(name, "kpi_rollup_consistency")) {
        *out = CHECK_KPI_ROLLUP_CONSISTENCY;
    } else if (ulab_streq(name, "performance_report_cell")) {
        *out = CHECK_PERFORMANCE_REPORT_CELL;
    } else if (ulab_streq(name, "performance_report_row")) {
        *out = CHECK_PERFORMANCE_REPORT_ROW;
    } else if (ulab_streq(name, "revenue_summary")) {
        *out = CHECK_REVENUE_SUMMARY;
    } else if (ulab_streq(name, "subscriber_billing_summary")) {
        *out = CHECK_SUBSCRIBER_BILLING_SUMMARY;
    } else if (ulab_streq(name, "payment_entitlement_reconciles")) {
        *out = CHECK_PAYMENT_ENTITLEMENT_RECONCILES;
    } else if (ulab_streq(name, "package_dashboard_metric")) {
        *out = CHECK_PACKAGE_DASHBOARD_METRIC;
    } else if (ulab_streq(name, "network_overview_metric")) {
        *out = CHECK_NETWORK_OVERVIEW_METRIC;
    } else if (ulab_streq(name, "console_inventory_reconciles")) {
        *out = CHECK_CONSOLE_INVENTORY_RECONCILES;
    } else if (ulab_streq(name, "usage_aggregate")) {
        *out = CHECK_USAGE_AGGREGATE;
    } else if (ulab_streq(name, "node_state")) *out = CHECK_NODE_STATE;
    else if (ulab_streq(name, "dashboard_loads")) {
        *out = CHECK_DASHBOARD_LOADS;
    } else if (ulab_streq(name, "dashboard_section_ok")) {
        *out = CHECK_DASHBOARD_SECTION_OK;
    } else if (ulab_streq(name, "node_version_equals")) {
        *out = CHECK_NODE_VERSION_EQUALS;
    } else if (ulab_streq(name, "node_health_ok")) {
        *out = CHECK_NODE_HEALTH_OK;
    } else if (ulab_streq(name, "release_unavailable")) {
        *out = CHECK_RELEASE_UNAVAILABLE;
    } else if (ulab_streq(name, "history_preserved")) {
        *out = CHECK_HISTORY_PRESERVED;
    } else if (ulab_streq(name, "audit_event_exists")) {
        *out = CHECK_AUDIT_EVENT_EXISTS;
    } else if (ulab_streq(name, "relationship_exists")) {
        *out = CHECK_RELATIONSHIP_EXISTS;
    } else if (ulab_streq(name, "relationship_ended")) {
        *out = CHECK_RELATIONSHIP_ENDED;
    } else if (ulab_streq(name, "balance_non_negative")) {
        *out = CHECK_BALANCE_NON_NEGATIVE;
    } else return ULAB_ERR;
    return ULAB_OK;
}

static int indent_of(const char *s) {
    int n = 0;

    while (*s == ' ') {
        n++;
        s++;
    }
    return n;
}

static void strip_comment(char *s) {
    int quote = 0;

    while (*s) {
        if (*s == '"' || *s == '\'') {
            quote = !quote;
        }
        if (!quote && *s == '#') {
            *s = '\0';
            return;
        }
        s++;
    }
}

static int split_kv(char *s, char **key, char **val) {
    char *p;

    p = strchr(s, ':');
    if (p == NULL) {
        return ULAB_ERR;
    }
    *p = '\0';
    *key = ulab_trim(s);
    *val = ulab_trim(p + 1);
    return ULAB_OK;
}

static int parse_selector_value(selector_t *sel, const char *key,
                                const char *val) {
    if (!ulab_streq(key, "ues") && !ulab_streq(key, "nodes") &&
        !ulab_streq(key, "sites") && !ulab_streq(key, "networks")) {
        return ULAB_ERR;
    }

    memset(sel, 0, sizeof(*sel));
    if (ulab_streq(val, "all")) {
        sel->kind = SEL_ALL;
        return ULAB_OK;
    }

    if (val != NULL && val[0] != '\0') {
        sel->kind = SEL_REF;
        return ulab_copy(sel->value, sizeof(sel->value), val);
    }

    return ULAB_ERR;
}

static int parse_inline_list(const char *v, const char *item) {
    char buf[ULAB_MAX_LINE];
    char *p;
    char *tok;

    if (!ulab_starts(v, "[") || !ulab_ends(v, "]")) {
        return 0;
    }
    if (ulab_copy(buf, sizeof(buf), v + 1) != ULAB_OK) {
        return 0;
    }
    buf[strlen(buf) - 1] = '\0';
    tok = strtok_r(buf, ",", &p);
    while (tok != NULL) {
        tok = ulab_trim(tok);
        if (ulab_streq(tok, item)) {
            return 1;
        }
        tok = strtok_r(NULL, ",", &p);
    }
    return 0;
}

static event_spec_t *new_event(phase_spec_t *p, const char *type,
                               ulab_error_t *err) {
    event_spec_t *e;

    if (p->event_count >= ULAB_MAX_EVENTS) {
        snprintf(err->msg, sizeof(err->msg), "too many events");
        return NULL;
    }
    e = &p->events[p->event_count++];
    memset(e, 0, sizeof(*e));
    if (scenario_event_from_name(type, &e->type) != ULAB_OK) {
        snprintf(err->msg, sizeof(err->msg), "unknown event type: %s", type);
        return NULL;
    }
    return e;
}

static check_spec_t *new_check(check_spec_t *arr, size_t *cnt,
                               const char *type, ulab_error_t *err) {
    check_spec_t *c;

    if (*cnt >= ULAB_MAX_CHECKS) {
        snprintf(err->msg, sizeof(err->msg), "too many checks");
        return NULL;
    }
    c = &arr[(*cnt)++];
    memset(c, 0, sizeof(*c));
    c->tolerance_percent = 2;
    c->tolerance_value = 0.000001;
    c->timeout_seconds = 300;
    c->poll_seconds = 5;
    snprintf(c->comparator, sizeof(c->comparator), "equals");
    if (scenario_check_from_name(type, &c->type) != ULAB_OK) {
        snprintf(err->msg, sizeof(err->msg), "unknown check type: %s", type);
        return NULL;
    }

    /*
     * Package lifecycle checks should fail quickly when the backend state
     * does not transition. Dashboard metrics can lag slightly longer, but
     * should not hold a P0 run for five minutes by default. Scenario YAML
     * may still override both values explicitly.
     */
    if (c->type == CHECK_PACKAGE_STATE) {
        c->timeout_seconds = 30;
    } else if (c->type == CHECK_PACKAGE_BUSINESS_METRICS) {
        c->timeout_seconds = 60;
        c->poll_seconds = 10;
    }
    return c;
}

static int apply_check_field(check_spec_t *c, const char *key,
                             const char *val) {
    if (ulab_streq(key, "target")) return ulab_copy(c->target,
        sizeof(c->target), val);
    if (ulab_streq(key, "expected")) return ulab_copy(c->expected,
        sizeof(c->expected), val);
    if (ulab_streq(key, "package")) return ulab_copy(c->package_ref,
        sizeof(c->package_ref), val);
    if (ulab_streq(key, "after_package") ||
        ulab_streq(key, "other_package")) {
        return ulab_copy(c->other_package_ref,
                         sizeof(c->other_package_ref), val);
    }
    if (ulab_streq(key, "variant")) return ulab_copy(c->variant,
        sizeof(c->variant), val);
    if (ulab_streq(key, "view")) return ulab_copy(c->view,
        sizeof(c->view), val);
    if (ulab_streq(key, "ref")) return ulab_copy(c->ref,
        sizeof(c->ref), val);
    if (ulab_streq(key, "entity")) return ulab_copy(c->entity,
        sizeof(c->entity), val);
    if (ulab_streq(key, "status")) return ulab_copy(c->status,
        sizeof(c->status), val);
    if (ulab_streq(key, "app")) return ulab_copy(c->app,
        sizeof(c->app), val);
    if (ulab_streq(key, "key")) return ulab_copy(c->key,
        sizeof(c->key), val);
    if (ulab_streq(key, "span")) return ulab_copy(c->span,
        sizeof(c->span), val);
    if (ulab_streq(key, "op")) return ulab_copy(c->op,
        sizeof(c->op), val);
    if (ulab_streq(key, "comparator")) return ulab_copy(c->comparator,
        sizeof(c->comparator), val);
    if (ulab_streq(key, "report")) return ulab_copy(c->report,
        sizeof(c->report), val);
    if (ulab_streq(key, "column")) return ulab_copy(c->column,
        sizeof(c->column), val);
    if (ulab_streq(key, "scope_key")) return ulab_copy(c->scope_key,
        sizeof(c->scope_key), val);
    if (ulab_streq(key, "scope_value")) return ulab_copy(c->scope_value,
        sizeof(c->scope_value), val);
    if (ulab_streq(key, "trend_direction")) {
        return ulab_copy(c->trend_direction,
                         sizeof(c->trend_direction), val);
    }
    if (ulab_streq(key, "currency")) return ulab_copy(c->currency,
        sizeof(c->currency), val);
    if (ulab_streq(key, "payment_method")) {
        return ulab_copy(c->payment_method,
                         sizeof(c->payment_method), val);
    }
    if (ulab_streq(key, "expected_value")) {
        if (ulab_parse_double(val, &c->expected_value)) return ULAB_ERR;
        c->has_expected_value = 1;
        return ULAB_OK;
    }
    if (ulab_streq(key, "tolerance")) {
        return ulab_parse_double(val, &c->tolerance_value);
    }
    if (ulab_streq(key, "expected_count")) {
        if (ulab_parse_u32(val, &c->expected_count)) return ULAB_ERR;
        c->has_expected_count = 1;
        return ULAB_OK;
    }
    if (ulab_streq(key, "timeout_seconds")) {
        return ulab_parse_u32(val, &c->timeout_seconds);
    }
    if (ulab_streq(key, "poll_seconds")) {
        return ulab_parse_u32(val, &c->poll_seconds);
    }
    if (ulab_streq(key, "expected_partial")) {
        c->expected_partial = ulab_streq(val, "true") ||
            ulab_streq(val, "1");
        c->has_expected_partial = 1;
        return ULAB_OK;
    }
    if (ulab_streq(key, "require_computed_at")) {
        c->require_computed_at = ulab_streq(val, "true") ||
            ulab_streq(val, "1");
        return ULAB_OK;
    }
    if (ulab_streq(key, "require_scope")) {
        c->require_scope = ulab_streq(val, "true") ||
            ulab_streq(val, "1");
        return ULAB_OK;
    }
    if (ulab_streq(key, "require_trend_consistency")) {
        c->require_trend_consistency = ulab_streq(val, "true") ||
            ulab_streq(val, "1");
        return ULAB_OK;
    }
    if (ulab_streq(key, "section") || ulab_streq(key, "version") ||
        ulab_streq(key, "tag")) {
        return ulab_copy(c->expected, sizeof(c->expected), val);
    }
    if (ulab_streq(key, "amount_mb")) {
        return ulab_parse_u64(val, &c->expected_used_mb);
    }
    if (ulab_streq(key, "expected_used_mb")) {
        return ulab_parse_u64(val, &c->expected_used_mb);
    }
    if (ulab_streq(key, "expected_remaining_mb")) {
        return ulab_parse_u64(val, &c->expected_remaining_mb);
    }
    if (ulab_streq(key, "tolerance_percent")) {
        return ulab_parse_u32(val, &c->tolerance_percent);
    }
    if (ulab_streq(key, "required")) {
        c->required = ulab_streq(val, "true") || ulab_streq(val, "1");
        return ULAB_OK;
    }
    if (ulab_streq(key, "immediate")) {
        c->immediate = ulab_streq(val, "true") || ulab_streq(val, "1");
        return ULAB_OK;
    }
    if (parse_selector_value(&c->ues, key, val) == ULAB_OK &&
        ulab_streq(key, "ues")) return ULAB_OK;
    if (parse_selector_value(&c->nodes, key, val) == ULAB_OK &&
        ulab_streq(key, "nodes")) return ULAB_OK;
    if (parse_selector_value(&c->sites, key, val) == ULAB_OK &&
        ulab_streq(key, "sites")) return ULAB_OK;
    if (parse_selector_value(&c->networks, key, val) == ULAB_OK &&
        ulab_streq(key, "networks")) return ULAB_OK;
    if (ulab_streq(key, "sample_per_site")) {
        c->ues.kind = SEL_SAMPLE_PER_SITE;
        return ulab_parse_u32(val, &c->ues.count);
    }
    if (ulab_streq(key, "created_in_phase")) {
        c->ues.kind = SEL_CREATED_IN_PHASE;
        return ulab_copy(c->ues.value, sizeof(c->ues.value), val);
    }
    if (ulab_streq(key, "type_selector")) {
        if (ulab_copy(c->nodes.value, sizeof(c->nodes.value), val)) {
            return ULAB_ERR;
        }
        c->nodes.kind = c->nodes.count > 0 ?
            SEL_NODE_TYPE_COUNT_PER_NETWORK : SEL_NODE_TYPE;
        return ULAB_OK;
    }
    if (ulab_streq(key, "count_per_network")) {
        if (ulab_parse_u32(val, &c->nodes.count) || c->nodes.count == 0) {
            return ULAB_ERR;
        }
        c->nodes.kind = SEL_NODE_TYPE_COUNT_PER_NETWORK;
        return ULAB_OK;
    }
    return ULAB_ERR;
}

static int apply_event_field(event_spec_t *e, const char *key,
                             const char *val) {
    if (ulab_streq(key, "name")) return ulab_copy(e->name,
        sizeof(e->name), val);
    if (ulab_streq(key, "amount_mb") ||
        ulab_streq(key, "seconds")) {
        return ulab_parse_u64(val, &e->amount_mb);
    }
    if (ulab_streq(key, "profile")) return ulab_copy(e->profile,
        sizeof(e->profile), val);
    if (ulab_streq(key, "count_per_site")) {
        return ulab_parse_u32(val, &e->count_per_site);
    }
    if (ulab_streq(key, "package")) return ulab_copy(e->package_ref,
        sizeof(e->package_ref), val);
    if (ulab_streq(key, "other_package")) {
        return ulab_copy(e->other_package_ref,
                         sizeof(e->other_package_ref), val);
    }
    if (ulab_streq(key, "variant")) return ulab_copy(e->variant,
        sizeof(e->variant), val);
    if (ulab_streq(key, "app")) return ulab_copy(e->app,
        sizeof(e->app), val);
    if (ulab_streq(key, "tag")) return ulab_copy(e->tag,
        sizeof(e->tag), val);
    if (ulab_streq(key, "amount")) {
        if (ulab_parse_double(val, &e->amount)) return ULAB_ERR;
        e->has_amount = 1;
        return ULAB_OK;
    }
    if (ulab_streq(key, "offset_seconds")) {
        char *end;
        long long parsed;

        end = NULL;
        parsed = strtoll(val, &end, 10);
        if (end == val || *end != '\0') return ULAB_ERR;
        e->offset_seconds = (int64_t)parsed;
        return ULAB_OK;
    }
    if (ulab_streq(key, "currency")) return ulab_copy(e->currency,
        sizeof(e->currency), val);
    if (ulab_streq(key, "payer_email")) {
        return ulab_copy(e->payer_email, sizeof(e->payer_email), val);
    }
    if (ulab_streq(key, "payer_phone")) {
        return ulab_copy(e->payer_phone, sizeof(e->payer_phone), val);
    }
    if (ulab_streq(key, "idempotency_key")) {
        return ulab_copy(e->idempotency_key,
                         sizeof(e->idempotency_key), val);
    }
    if (ulab_streq(key, "active")) {
        e->active = ulab_streq(val, "true") || ulab_streq(val, "1");
        e->has_active = 1;
        return ULAB_OK;
    }
    if (ulab_streq(key, "status") || ulab_streq(key, "state") ||
        ulab_streq(key, "connectivity") ||
        ulab_streq(key, "version")) return ulab_copy(e->status,
        sizeof(e->status), val);
    if (ulab_streq(key, "expect_result")) return ulab_copy(e->expect_result,
        sizeof(e->expect_result), val);
    if (ulab_streq(key, "error_contains")) return ulab_copy(e->error_contains,
        sizeof(e->error_contains), val);
    if (parse_selector_value(&e->ues, key, val) == ULAB_OK &&
        ulab_streq(key, "ues")) return ULAB_OK;
    if (parse_selector_value(&e->sites, key, val) == ULAB_OK &&
        ulab_streq(key, "sites")) return ULAB_OK;
    if (parse_selector_value(&e->nodes, key, val) == ULAB_OK &&
        ulab_streq(key, "nodes")) return ULAB_OK;
    if (ulab_streq(key, "created_in_phase")) {
        e->ues.kind = SEL_CREATED_IN_PHASE;
        return ulab_copy(e->ues.value, sizeof(e->ues.value), val);
    }
    if (ulab_streq(key, "affected_by_phase")) {
        e->nodes.kind = SEL_AFFECTED_BY_PHASE;
        return ulab_copy(e->nodes.value, sizeof(e->nodes.value), val);
    }
    if (ulab_streq(key, "type_selector")) {
        if (ulab_copy(e->nodes.value, sizeof(e->nodes.value), val)) {
            return ULAB_ERR;
        }
        e->nodes.kind = e->nodes.count > 0 ?
            SEL_NODE_TYPE_COUNT_PER_NETWORK : SEL_NODE_TYPE;
        return ULAB_OK;
    }
    if (ulab_streq(key, "count_per_network")) {
        if (ulab_parse_u32(val, &e->nodes.count) || e->nodes.count == 0) {
            return ULAB_ERR;
        }
        e->nodes.kind = SEL_NODE_TYPE_COUNT_PER_NETWORK;
        return ULAB_OK;
    }
    return ULAB_ERR;
}


static int apply_event_expect_field(event_spec_t *e, const char *key,
                                    const char *val) {
    if (ulab_streq(key, "result")) {
        return ulab_copy(e->expect_result, sizeof(e->expect_result), val);
    }
    if (ulab_streq(key, "error_contains")) {
        return ulab_copy(e->error_contains, sizeof(e->error_contains), val);
    }
    return ULAB_ERR;
}

static int parse_item_value(char *line, char **key, char **val) {
    char *p = ulab_trim(line);

    if (!ulab_starts(p, "- ")) {
        return ULAB_ERR;
    }
    p += 2;
    return split_kv(p, key, val);
}

int scenario_load(const char *path, scenario_t *s, ulab_error_t *err) {
    FILE *fp;
    char line[ULAB_MAX_LINE];
    parse_sec_t sec = SEC_NONE;
    phase_spec_t *phase = NULL;
    event_spec_t *event = NULL;
    check_spec_t *check = NULL;
    package_spec_t *pkg = NULL;
    profile_spec_t *prof = NULL;
    profile_bucket_t *bucket = NULL;
    int lineno = 0;

    scenario_init(s);
    fp = fopen(path, "r");
    if (fp == NULL) {
        snprintf(err->msg, sizeof(err->msg), "unable to open %s", path);
        return ULAB_ERR;
    }

    while (fgets(line, sizeof(line), fp) != NULL) {
        char expanded[ULAB_MAX_LINE];
        char missing[ULAB_MAX_REF];
        char *key;
        char *val;
        char *p;
        int ind;

        lineno++;
        strip_comment(line);
        memset(missing, 0, sizeof(missing));
        if (expand_environment(line, expanded, sizeof(expanded),
                               missing, sizeof(missing)) != ULAB_OK) {
            if (missing[0] != '\0') {
                snprintf(err->msg, sizeof(err->msg),
                         "line %d: missing environment variable %s",
                         lineno, missing);
            } else {
                snprintf(err->msg, sizeof(err->msg),
                         "line %d: invalid environment expansion", lineno);
            }
            fclose(fp);
            return ULAB_ERR;
        }
        if (ulab_copy(line, sizeof(line), expanded) != ULAB_OK) {
            snprintf(err->msg, sizeof(err->msg),
                     "line %d: expanded line too long", lineno);
            fclose(fp);
            return ULAB_ERR;
        }
        p = ulab_trim(line);
        if (*p == '\0') {
            continue;
        }
        ind = indent_of(line);
        if (!ulab_starts(p, "- ") &&
            split_kv(p, &key, &val) != ULAB_OK) {
            snprintf(err->msg, sizeof(err->msg), "line %d: bad syntax", lineno);
            fclose(fp);
            return ULAB_ERR;
        }

        if (ind == 0 && !ulab_starts(p, "- ")) {
            if (ulab_streq(key, "version")) {
                if (ulab_parse_u32(val, &s->version) != ULAB_OK) goto bad;
            } else if (ulab_streq(key, "name")) {
                if (ulab_copy(s->name, sizeof(s->name), val) != ULAB_OK) {
                    goto bad;
                }
            } else if (ulab_streq(key, "description")) {
                if (ulab_copy(s->description, sizeof(s->description), val)) {
                    goto bad;
                }
            } else if (ulab_streq(key, "seed")) {
                if (ulab_parse_u32(val, &s->seed) != ULAB_OK) goto bad;
            } else if (ulab_streq(key, "suite")) {
                if (ulab_copy(s->suite, sizeof(s->suite), val)) goto bad;
            } else if (ulab_streq(key, "priority")) {
                if (ulab_copy(s->priority, sizeof(s->priority), val)) goto bad;
            } else if (ulab_streq(key, "tags")) {
                if (ulab_copy(s->tags, sizeof(s->tags), val)) goto bad;
            } else if (ulab_streq(key, "status")) {
                if (ulab_copy(s->status, sizeof(s->status), val)) goto bad;
            } else if (ulab_streq(key, "generated")) {
                s->generated = ulab_streq(val, "true") || ulab_streq(val, "1");
            } else if (ulab_streq(key, "entity")) {
                if (ulab_copy(s->entity, sizeof(s->entity), val)) goto bad;
            } else if (ulab_streq(key, "action")) {
                if (ulab_copy(s->action, sizeof(s->action), val)) goto bad;
            } else if (ulab_streq(key, "world")) sec = SEC_WORLD;
            else if (ulab_streq(key, "packages")) sec = SEC_PACKAGES;
            else if (ulab_streq(key, "setup")) sec = SEC_SETUP;
            else if (ulab_streq(key, "provider")) sec = SEC_PROVIDER;
            else if (ulab_streq(key, "runtime")) sec = SEC_RUNTIME;
            else if (ulab_streq(key, "profiles")) sec = SEC_PROFILES;
            else if (ulab_streq(key, "phases")) sec = SEC_PHASES;
            else if (ulab_streq(key, "final_checks")) sec = SEC_FINAL_CHECKS;
            else goto unknown;
            continue;
        }

        if (sec == SEC_WORLD) {
            if (ind == 2 && ulab_streq(key, "networks")) {
                if (ulab_parse_u32(val, &s->world.networks) != ULAB_OK) {
                    goto bad;
                }
            } else if (ind == 2 && ulab_streq(key, "sites_per_network")) {
                if (ulab_parse_u32(val, &s->world.sites_per_network)) goto bad;
            } else if (ind == 2 && ulab_streq(key, "ues_per_site")) {
                if (ulab_parse_u32(val, &s->world.ues_per_site)) goto bad;
            } else if (ind == 2 && ulab_streq(key, "sims_per_subscriber")) {
                if (ulab_parse_u32(val, &s->world.sims_per_subscriber)) {
                    goto bad;
                }
            } else if (ind == 2 && ulab_streq(key, "nodes_per_site")) {
                sec = SEC_NODES_PER_SITE;
            } else goto unknown;
            continue;
        }
        if (sec == SEC_NODES_PER_SITE) {
            if (ind == 4 && ulab_streq(key, "tower")) {
                if (ulab_parse_u32(val, &s->world.tower_per_site)) goto bad;
            } else if (ind == 4 && ulab_streq(key, "amplifier")) {
                if (ulab_parse_u32(val, &s->world.amplifier_per_site)) goto bad;
            } else if (ind == 4 && ulab_streq(key, "controller")) {
                if (ulab_parse_u32(val, &s->world.controller_per_site)) goto bad;
            } else if (ind == 2 && ulab_streq(key, "ues_per_site")) {
                if (ulab_parse_u32(val, &s->world.ues_per_site)) goto bad;
                sec = SEC_WORLD;
            } else if (ind == 2 && ulab_streq(key, "sims_per_subscriber")) {
                if (ulab_parse_u32(val, &s->world.sims_per_subscriber)) {
                    goto bad;
                }
                sec = SEC_WORLD;
            } else if (ind == 2 && ulab_streq(key, "networks")) {
                if (ulab_parse_u32(val, &s->world.networks)) goto bad;
                sec = SEC_WORLD;
            } else if (ind == 2 && ulab_streq(key, "sites_per_network")) {
                if (ulab_parse_u32(val, &s->world.sites_per_network)) goto bad;
                sec = SEC_WORLD;
            } else goto unknown;
            continue;
        }
        if (sec == SEC_PACKAGES) {
            if (ind == 2 && ulab_starts(p, "- ")) {
                if (s->package_count >= ULAB_MAX_PACKAGES) goto many;
                pkg = &s->packages[s->package_count++];
                memset(pkg, 0, sizeof(*pkg));
                snprintf(pkg->currency, sizeof(pkg->currency), "USD");
                snprintf(pkg->country, sizeof(pkg->country), "USA");
                snprintf(pkg->scope, sizeof(pkg->scope), "network");
                pkg->active = 1;
                if (parse_item_value(p, &key, &val) ||
                    !ulab_streq(key, "ref")) goto bad;
                if (ulab_copy(pkg->ref, sizeof(pkg->ref), val)) goto bad;
            } else if (ind == 4 && pkg != NULL) {
                if (ulab_streq(key, "name")) {
                    if (ulab_copy(pkg->name, sizeof(pkg->name), val)) goto bad;
                } else if (ulab_streq(key, "data_mb")) {
                    if (ulab_parse_u64(val, &pkg->data_mb)) goto bad;
                } else if (ulab_streq(key, "duration_days")) {
                    if (ulab_parse_u32(val, &pkg->duration_days)) goto bad;
                } else if (ulab_streq(key, "duration_minutes")) {
                    if (ulab_parse_u32(val, &pkg->duration_minutes)) goto bad;
                } else if (ulab_streq(key, "duration_hours")) {
                    uint32_t hours = 0;
                    if (ulab_parse_u32(val, &hours)) goto bad;
                    if (hours > UINT32_MAX / 60u) goto bad;
                    pkg->duration_minutes = hours * 60u;
                } else if (ulab_streq(key, "amount")) {
                    if (ulab_parse_double(val, &pkg->amount)) goto bad;
                } else if (ulab_streq(key, "currency")) {
                    if (ulab_copy(pkg->currency, sizeof(pkg->currency), val)) {
                        goto bad;
                    }
                } else if (ulab_streq(key, "country")) {
                    if (ulab_copy(pkg->country, sizeof(pkg->country), val)) {
                        goto bad;
                    }
                } else if (ulab_streq(key, "active")) {
                    pkg->active = ulab_streq(val, "true") ||
                        ulab_streq(val, "1");
                } else if (ulab_streq(key, "assign_percent")) {
                    if (ulab_parse_u32(val, &pkg->assign_percent)) goto bad;
                } else if (ulab_streq(key, "scope")) {
                    if (ulab_copy(pkg->scope, sizeof(pkg->scope), val)) {
                        goto bad;
                    }
                } else if (ulab_streq(key, "network")) {
                    if (ulab_copy(pkg->network_ref,
                                  sizeof(pkg->network_ref), val)) {
                        goto bad;
                    }
                } else goto unknown;
            } else goto unknown;
            continue;
        }
        if (sec == SEC_SETUP) {
            if (ind == 2 && ulab_streq(key, "create_via_bff")) {
                if (val[0] != '\0') {
                    s->setup.create_networks =
                        parse_inline_list(val, "networks");
                    s->setup.create_sites =
                        parse_inline_list(val, "sites");
                    s->setup.create_nodes =
                        parse_inline_list(val, "nodes");
                    s->setup.create_node_site_links =
                        parse_inline_list(val, "node_site_links");
                    s->setup.create_packages =
                        parse_inline_list(val, "packages");
                    s->setup.create_subscribers =
                        parse_inline_list(val, "subscribers");
                    s->setup.create_sims =
                        parse_inline_list(val, "sims");
                } else {
                    sec = SEC_SETUP_LIST;
                }
            } else goto unknown;
            continue;
        }
        if (sec == SEC_SETUP_LIST) {
            if (ind == 4 && ulab_starts(p, "- ")) {
                char *item = ulab_trim(p + 2);
                if (ulab_streq(item, "networks")) s->setup.create_networks = 1;
                else if (ulab_streq(item, "sites")) s->setup.create_sites = 1;
                else if (ulab_streq(item, "nodes")) s->setup.create_nodes = 1;
                else if (ulab_streq(item, "node_site_links")) {
                    s->setup.create_node_site_links = 1;
                } else if (ulab_streq(item, "packages")) {
                    s->setup.create_packages = 1;
                } else if (ulab_streq(item, "subscribers")) {
                    s->setup.create_subscribers = 1;
                } else if (ulab_streq(item, "sims")) s->setup.create_sims = 1;
                else goto unknown;
            } else goto unknown;
            continue;
        }

        if (sec == SEC_PROVIDER) {
            if (ind == 2 && ulab_streq(key, "type")) {
                if (ulab_copy(s->provider.type,
                    sizeof(s->provider.type), val)) goto bad;
            } else goto unknown;
            continue;
        }

        if (sec == SEC_RUNTIME) {
            if (ind == 2 && ulab_streq(key, "start")) {
                s->runtime.start_nodes = parse_inline_list(val, "nodes");
                s->runtime.start_ues = parse_inline_list(val, "ues");
            } else if (ind == 2 && ulab_streq(key, "wait")) {
                s->runtime.wait_nodes_ready = parse_inline_list(val,
                    "nodes_ready");
                s->runtime.wait_ues_attached = parse_inline_list(val,
                    "ues_attached");
            } else goto unknown;
            continue;
        }
        if (sec == SEC_PROFILES || sec == SEC_PROFILE_ONE ||
            sec == SEC_PROFILE_BUCKET) {
            if (ind == 2) {
                if (s->profile_count >= ULAB_MAX_BUCKETS) goto many;
                prof = &s->profiles[s->profile_count++];
                memset(prof, 0, sizeof(*prof));
                if (ulab_copy(prof->name, sizeof(prof->name), key)) goto bad;
                sec = SEC_PROFILE_ONE;
            } else if (ind == 4 && prof != NULL) {
                if (prof->bucket_count >= ULAB_MAX_BUCKETS) goto many;
                bucket = &prof->buckets[prof->bucket_count++];
                memset(bucket, 0, sizeof(*bucket));
                if (ulab_copy(bucket->name, sizeof(bucket->name), key)) {
                    goto bad;
                }
                sec = SEC_PROFILE_BUCKET;
            } else if (ind == 6 && bucket != NULL) {
                if (ulab_streq(key, "percent")) {
                    if (ulab_parse_u32(val, &bucket->percent)) goto bad;
                } else if (ulab_streq(key, "amount_mb")) {
                    if (ulab_parse_u64(val, &bucket->amount_mb)) goto bad;
                } else goto unknown;
            } else goto unknown;
            continue;
        }
        if (sec == SEC_PHASES) {
            if (ind == 2 && ulab_starts(p, "- ")) {
                if (s->phase_count >= ULAB_MAX_PHASES) goto many;
                phase = &s->phases[s->phase_count++];
                memset(phase, 0, sizeof(*phase));
                if (parse_item_value(p, &key, &val) ||
                    !ulab_streq(key, "name")) goto bad;
                if (ulab_copy(phase->name, sizeof(phase->name), val)) goto bad;
            } else if (ind == 4 && ulab_streq(key, "events")) {
                sec = SEC_PHASE_EVENTS;
            } else if (ind == 4 && ulab_streq(key, "checks")) {
                sec = SEC_PHASE_CHECKS;
            } else goto unknown;
            continue;
        }
        if (sec == SEC_PHASE_EVENTS) {
            if (ind == 6 && ulab_starts(p, "- ")) {
                if (parse_item_value(p, &key, &val) ||
                    !ulab_streq(key, "type")) goto bad;
                event = new_event(phase, val, err);
                if (event == NULL) goto fail;
            } else if (ind == 8 && event != NULL &&
                       ulab_streq(key, "expect") && val[0] == '\0') {
                sec = SEC_EVENT_EXPECT;
            } else if (ind == 8 && event != NULL) {
                if (apply_event_field(event, key, val) != ULAB_OK) goto unknown;
            } else if (ind == 4 && ulab_streq(key, "checks")) {
                sec = SEC_PHASE_CHECKS;
            } else goto unknown;
            continue;
        }
        if (sec == SEC_EVENT_EXPECT) {
            if (ind == 10 && event != NULL) {
                if (apply_event_expect_field(event, key, val) != ULAB_OK) {
                    goto unknown;
                }
            } else if (ind == 8 && event != NULL) {
                sec = SEC_PHASE_EVENTS;
                if (apply_event_field(event, key, val) != ULAB_OK) goto unknown;
            } else if (ind == 6 && ulab_starts(p, "- ")) {
                sec = SEC_PHASE_EVENTS;
                if (parse_item_value(p, &key, &val) ||
                    !ulab_streq(key, "type")) goto bad;
                event = new_event(phase, val, err);
                if (event == NULL) goto fail;
            } else if (ind == 4 && ulab_streq(key, "checks")) {
                sec = SEC_PHASE_CHECKS;
            } else goto unknown;
            continue;
        }
        if (sec == SEC_PHASE_CHECKS || sec == SEC_FINAL_CHECKS) {
            check_spec_t *arr = sec == SEC_FINAL_CHECKS ?
                s->final_checks : phase->checks;
            size_t *cnt = sec == SEC_FINAL_CHECKS ?
                &s->final_check_count : &phase->check_count;

            if (sec == SEC_PHASE_CHECKS && ind == 2 &&
                ulab_starts(p, "- ")) {
                if (s->phase_count >= ULAB_MAX_PHASES) goto many;
                phase = &s->phases[s->phase_count++];
                memset(phase, 0, sizeof(*phase));
                if (parse_item_value(p, &key, &val) ||
                    !ulab_streq(key, "name")) goto bad;
                if (ulab_copy(phase->name, sizeof(phase->name), val)) goto bad;
                sec = SEC_PHASES;
                continue;
            }
            if ((ind == 2 || ind == 6) && ulab_starts(p, "- ")) {
                if (parse_item_value(p, &key, &val) ||
                    !ulab_streq(key, "type")) goto bad;
                check = new_check(arr, cnt, val, err);
                if (check == NULL) goto fail;
            } else if ((ind == 4 || ind == 8) && check != NULL) {
                if (apply_check_field(check, key, val) != ULAB_OK) goto unknown;
            } else if (ind == 2 && sec != SEC_FINAL_CHECKS) {
                sec = SEC_PHASES;
            } else goto unknown;
            continue;
        }
    }

    fclose(fp);
    return ULAB_OK;

bad:
    snprintf(err->msg, sizeof(err->msg), "line %d: invalid value", lineno);
    goto fail;
unknown:
    snprintf(err->msg, sizeof(err->msg), "line %d: unknown field", lineno);
    goto fail;
many:
    snprintf(err->msg, sizeof(err->msg), "line %d: too many entries", lineno);
fail:
    fclose(fp);
    return ULAB_ERR;
}

void scenario_list_events(void) {
    int i;

    for (i = EVT_TRAFFIC; i <= EVT_CHECK; i++) {
        printf("%s\n", scenario_event_name((event_type_t)i));
    }
}

void scenario_list_checks(void) {
    int i;

    for (i = CHECK_BACKEND_COUNT; i <= CHECK_BALANCE_NON_NEGATIVE; i++) {
        printf("%s\n", scenario_check_name((check_type_t)i));
    }
}
