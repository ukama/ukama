/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "scenario.h"
#include "world.h"
#include "util.h"

static int fail(ulab_error_t *err, const char *msg) {
    snprintf(err->msg, sizeof(err->msg), "%s", msg);
    return ULAB_ERR;
}

int scenario_validate(const scenario_t *s, ulab_error_t *err) {
    uint32_t pct = 0;
    size_t assertion_count = 0;
    size_t i;
    size_t j;

    if (s->version != ULAB_SCHEMA_VER) {
        return fail(err, "unsupported scenario version");
    }

    if (s->name[0] == '\0') {
        return fail(err, "missing scenario name");
    }

    if (!ulab_streq(s->status, "active") &&
        !ulab_streq(s->status, "wip") &&
        !ulab_streq(s->status, "skip") &&
        !ulab_streq(s->status, "xfail")) {
        return fail(err, "scenario status must be active/wip/skip/xfail");
    }

    if (!ulab_streq(s->provider.type, "virtual")) {
        return fail(err, "unsupported provider type");
    }

    if (s->world.networks == 0 || s->world.sites_per_network == 0) {
        return fail(err, "world must include networks and sites");
    }

    if (s->world.tower_per_site == 0) {
        return fail(err, "world.nodes_per_site must include tower node");
    }

    if (s->world.ues_per_site > 0 && s->world.ues_per_site >
        s->world.tower_per_site * ULAB_MAX_UES_PER_TOWER) {
        return fail(err, "world.ues_per_site exceeds 500 UEs per tower");
    }

    if (s->world.sims_per_subscriber > s->world.ues_per_site &&
        s->world.ues_per_site > 0) {
        return fail(err,
                    "world.sims_per_subscriber cannot exceed ues_per_site");
    }

    if (s->world.tower_per_site + s->world.amplifier_per_site +
        s->world.controller_per_site == 0) {
        return fail(err, "world.nodes_per_site must include nodes");
    }

    if (s->world.ues_per_site > 0 && s->package_count == 0) {
        return fail(err, "at least one package is required");
    }

    for (i = 0; i < s->package_count; i++) {
        const package_spec_t *p = &s->packages[i];
        int network_exists;

        if (p->ref[0] == '\0' || p->name[0] == '\0') {
            return fail(err, "package ref/name is required");
        }

        if (p->data_mb == 0 ||
            (p->duration_days == 0 && p->duration_minutes == 0)) {
            return fail(err,
                        "package data_mb and duration_days or "
                        "duration_minutes are required");
        }

        if (p->duration_days != 0 && p->duration_minutes != 0) {
            return fail(err,
                        "package duration_days and duration_minutes are "
                        "mutually exclusive");
        }

        if (!ulab_streq(p->scope, "network") &&
            !ulab_streq(p->scope, "organization")) {
            return fail(err,
                        "package scope must be network or organization");
        }
        if (ulab_streq(p->scope, "organization") &&
            p->network_ref[0] != '\0') {
            return fail(err,
                        "organization package cannot specify network");
        }
        if (p->network_ref[0] != '\0' &&
            !ulab_starts(p->network_ref, "net-")) {
            return fail(err, "package network must use a net-NNN ref");
        }
        network_exists = p->network_ref[0] == '\0';
        for (j = 0; !network_exists && j < s->world.networks; j++) {
            char expected_ref[ULAB_MAX_REF];

            snprintf(expected_ref, sizeof(expected_ref), "net-%03zu", j + 1);
            network_exists = ulab_streq(p->network_ref, expected_ref);
        }
        if (!network_exists) {
            return fail(err, "package network does not exist in world");
        }

        pct += p->assign_percent;
    }

    if (s->package_count > 0 && pct != 100) {
        return fail(err, "package assign_percent values must add to 100");
    }

    if (!s->setup.create_networks || !s->setup.create_sites ||
        !s->setup.create_nodes) {
        return fail(err, "setup.create_via_bff missing required entries");
    }

    if (s->package_count > 0 && !s->setup.create_packages) {
        return fail(err, "setup.create_via_bff missing packages");
    }

    if (s->world.ues_per_site > 0 &&
        (!s->setup.create_packages || !s->setup.create_subscribers)) {
        return fail(err,
                    "UE scenarios require packages and subscribers");
    }

    if (s->world.ues_per_site > 0 && !s->setup.create_sims &&
        (s->runtime.start_ues || s->runtime.wait_ues_attached)) {
        return fail(err,
                    "runtime UE scenarios require setup sims");
    }

    if (s->world.ues_per_site == 0 &&
        (s->runtime.start_ues || s->runtime.wait_ues_attached)) {
        return fail(err, "runtime UE actions require world.ues_per_site > 0");
    }

    for (i = 0; i < s->profile_count; i++) {
        uint32_t pp = 0;
        size_t j;
        for (j = 0; j < s->profiles[i].bucket_count; j++) {
            pp += s->profiles[i].buckets[j].percent;
        }
        if (pp != 100) {
            return fail(err, "profile percent values must add to 100");
        }
    }

    if (s->phase_count == 0) {
        return fail(err, "at least one phase is required");
    }

    for (i = 0; i < s->phase_count; i++) {
        const phase_spec_t *phase;

        phase = &s->phases[i];
        assertion_count += phase->check_count;
        for (j = 0; j < phase->event_count; j++) {
            const event_spec_t *event;

            event = &phase->events[j];
            if (event->type == EVT_WAIT_UE_SESSIONS ||
                event->type == EVT_FINALIZE_UE_SESSIONS) {
                if (event->ues.kind == SEL_NONE) {
                    return fail(err,
                                "UE session event requires UE selector");
                }
                if (event->amount_mb > 300) {
                    return fail(err,
                                "UE session event exceeds 300 seconds");
                }
            }

            if (event->type == EVT_WAIT_NODE_CONNECTIVITY) {
                if (event->nodes.kind == SEL_NONE) {
                    return fail(err,
                                "wait_node_connectivity requires "
                                "node selector");
                }
                if (event->status[0] == '\0') {
                    return fail(err,
                                "wait_node_connectivity requires "
                                "connectivity");
                }
                if (!ulab_streq(event->status, "Online") &&
                    !ulab_streq(event->status, "Offline")) {
                    return fail(err,
                                "wait_node_connectivity connectivity "
                                "must be Online or Offline");
                }
                if (event->amount_mb > 900) {
                    return fail(err,
                                "wait_node_connectivity exceeds 900 seconds");
                }
            }

            if (event->type == EVT_PROMOTE_RELEASE) {
                if (event->app[0] == '\0' || event->tag[0] == '\0') {
                    return fail(err,
                                "promote_release event requires app and tag");
                }
                continue;
            }

            if (event->type == EVT_PURCHASE_PACKAGE) {
                if (event->ues.kind == SEL_NONE ||
                    event->package_ref[0] == '\0') {
                    return fail(err,
                                "purchase_package requires ues and package");
                }
                continue;
            }

            if (event->type == EVT_PURCHASE_PACKAGES_PARALLEL) {
                if (event->ues.kind == SEL_NONE ||
                    event->package_ref[0] == '\0' ||
                    event->other_package_ref[0] == '\0') {
                    return fail(err,
                                "purchase_packages_parallel requires ues, "
                                "package and other_package");
                }
                if (ulab_streq(event->package_ref,
                               event->other_package_ref)) {
                    return fail(err,
                                "parallel purchase packages must differ");
                }
                continue;
            }

            if (event->type == EVT_ALLOCATE_SIM) {
                if (event->ues.kind == SEL_NONE ||
                    event->package_ref[0] == '\0') {
                    return fail(err,
                                "allocate_sim requires ues and package");
                }
                continue;
            }

            if (event->type == EVT_CREATE_INVALID_PACKAGE) {
                if (event->package_ref[0] == '\0' ||
                    event->variant[0] == '\0') {
                    return fail(err,
                                "create_invalid_package requires package "
                                "and variant");
                }
                if (!ulab_streq(event->variant, "allowance") &&
                    !ulab_streq(event->variant, "duration") &&
                    !ulab_streq(event->variant, "price") &&
                    !ulab_streq(event->variant, "currency")) {
                    return fail(err,
                                "unsupported invalid package variant");
                }
                continue;
            }

            if (event->type == EVT_WAIT_PACKAGE_BOUNDARY) {
                if (event->ues.kind == SEL_NONE ||
                    event->package_ref[0] == '\0') {
                    return fail(err,
                                "wait_package_boundary requires ues and "
                                "package");
                }
                continue;
            }

            if (event->type == EVT_SET_PACKAGE_ACTIVE) {
                if (event->package_ref[0] == '\0' || !event->has_active) {
                    return fail(err,
                                "set_package_active requires package and "
                                "active");
                }
                continue;
            }

            if (event->type == EVT_WAIT_SIM_STATUS) {
                if (event->ues.kind == SEL_NONE ||
                    event->status[0] == '\0') {
                    return fail(err,
                                "wait_sim_status requires ues and status");
                }
                if (event->amount_mb > 900) {
                    return fail(err,
                                "wait_sim_status exceeds 900 seconds");
                }
                continue;
            }

            if (event->type == EVT_RESTART_SITE) {
                if (event->sites.kind == SEL_NONE &&
                    event->nodes.kind == SEL_NONE) {
                    return fail(err,
                                "restart_site requires site selector");
                }
                continue;
            }

            if (event->type == EVT_TOGGLE_INTERNET_SWITCH) {
                if (event->sites.kind == SEL_NONE &&
                    event->nodes.kind == SEL_NONE) {
                    return fail(err,
                                "toggle_internet_switch requires site "
                                "selector");
                }
                if (!event->has_port || event->port == 0) {
                    return fail(err,
                                "toggle_internet_switch requires port > 0");
                }
                if (event->status[0] == '\0') {
                    return fail(err,
                                "toggle_internet_switch requires state");
                }
                continue;
            }

            if (event->type == EVT_FAILURE_CONTROL) {
                if ((!ulab_streq(event->target, "payment") &&
                     !ulab_streq(event->target, "software")) ||
                    (!ulab_streq(event->status, "on") &&
                     !ulab_streq(event->status, "off") &&
                     !ulab_streq(event->status, "enabled") &&
                     !ulab_streq(event->status, "disabled") &&
                     !ulab_streq(event->status, "true") &&
                     !ulab_streq(event->status, "false"))) {
                    return fail(err,
                                "failure_control requires target "
                                "payment|software and state on|off");
                }
                continue;
            }

            if (event->type == EVT_DISCONNECT_NODES ||
                event->type == EVT_RECONNECT_NODES) {
                if (event->nodes.kind == SEL_NONE) {
                    return fail(err,
                                "node network event requires node selector");
                }
                continue;
            }

            if (event->type != EVT_SOFTWARE_UPDATE) {
                continue;
            }

            if (event->app[0] == '\0' || event->tag[0] == '\0') {
                return fail(err,
                            "software_update event requires app and tag");
            }

            if (event->nodes.kind == SEL_NONE) {
                return fail(err,
                            "software_update event requires node selector");
            }
        }

        for (j = 0; j < phase->check_count; j++) {
            const check_spec_t *check;

            check = &phase->checks[j];
            if (check->type == CHECK_NODE_VERSION_EQUALS &&
                (check->app[0] == '\0' || check->expected[0] == '\0')) {
                return fail(err,
                            "node_version_equals requires app and version/tag");
            }
            if (check->type == CHECK_RELEASE_UNAVAILABLE &&
                (check->app[0] == '\0' || check->expected[0] == '\0')) {
                return fail(err,
                            "release_unavailable requires app and version");
            }
            if ((check->type == CHECK_KPI_VALUE ||
                 check->type == CHECK_KPI_TREND ||
                 check->type == CHECK_KPI_CONTRACT ||
                 check->type == CHECK_KPI_ROLLUP_CONSISTENCY) &&
                check->key[0] == '\0') {
                return fail(err, "KPI checks require key");
            }
            if (check->type == CHECK_KPI_VALUE &&
                !check->has_expected_value) {
                return fail(err, "kpi_value requires expected_value");
            }
            if (check->type == CHECK_PERFORMANCE_REPORT_CELL &&
                (check->report[0] == '\0' || check->column[0] == '\0' ||
                 check->package_ref[0] == '\0' ||
                 !check->has_expected_value)) {
                return fail(err,
                            "performance_report_cell requires report, "
                            "column, package and expected_value");
            }
            if (check->type == CHECK_PERFORMANCE_REPORT_ROW &&
                (check->report[0] == '\0' ||
                 check->package_ref[0] == '\0')) {
                return fail(err,
                            "performance_report_row requires report and "
                            "package");
            }
            if ((check->type == CHECK_REVENUE_SUMMARY ||
                 check->type == CHECK_PACKAGE_DASHBOARD_METRIC ||
                 check->type == CHECK_NETWORK_OVERVIEW_METRIC) &&
                (check->column[0] == '\0' ||
                 !check->has_expected_value)) {
                return fail(err,
                            "console metric check requires column and "
                            "expected_value");
            }
            if (check->type == CHECK_SUBSCRIBER_BILLING_SUMMARY &&
                check->ues.kind == SEL_NONE) {
                return fail(err,
                            "subscriber_billing_summary requires ues");
            }
            if (check->type == CHECK_PAYMENT_ENTITLEMENT_RECONCILES &&
                (check->ues.kind == SEL_NONE ||
                 check->package_ref[0] == '\0')) {
                return fail(err,
                            "payment_entitlement_reconciles requires ues "
                            "and package");
            }
            if (check->type == CHECK_CONSOLE_INVENTORY_RECONCILES &&
                check->target[0] == '\0') {
                return fail(err,
                            "console_inventory_reconciles requires target");
            }
            if (check->type == CHECK_USAGE_AGGREGATE &&
                check->target[0] == '\0') {
                return fail(err, "usage_aggregate requires target");
            }
            if ((check->type == CHECK_PACKAGE_CATALOG_EQUALS ||
                 check->type == CHECK_PACKAGE_VISIBLE ||
                 check->type == CHECK_PACKAGE_HIDDEN ||
                 check->type == CHECK_PACKAGE_BUSINESS_METRICS ||
                 check->type == CHECK_PACKAGE_ASSIGNMENT_CHAIN) &&
                check->package_ref[0] == '\0') {
                return fail(err, "data-package check requires package");
            }
            if (check->type == CHECK_PACKAGE_ASSIGNMENT_CHAIN &&
                check->ues.kind == SEL_NONE) {
                return fail(err,
                            "package_assignment_chain requires ues");
            }
            if (check->type == CHECK_PACKAGE_ASSIGNMENT_CHAIN &&
                check->expected[0] != '\0' &&
                !ulab_streq(check->expected, "all") &&
                !ulab_streq(check->expected, "contiguous")) {
                return fail(err,
                            "package_assignment_chain expected must be "
                            "all or contiguous");
            }
            if (check->type == CHECK_PACKAGE_BUSINESS_METRICS &&
                !check->has_expected_value &&
                !check->has_expected_count) {
                return fail(err,
                            "package_business_metrics requires expected "
                            "value or count");
            }
            if (check->type == CHECK_PACKAGE_NAME_AVAILABLE &&
                (check->package_ref[0] == '\0' ||
                 check->variant[0] == '\0')) {
                return fail(err,
                            "package_name_available requires package and "
                            "variant");
            }
            if (check->type == CHECK_PACKAGE_NAME_AVAILABLE &&
                !ulab_streq(check->variant, "allowance") &&
                !ulab_streq(check->variant, "duration") &&
                !ulab_streq(check->variant, "price") &&
                !ulab_streq(check->variant, "currency")) {
                return fail(err,
                            "unsupported package name variant");
            }
            if (check->type == CHECK_SIM_UNALLOCATED &&
                check->ues.kind == SEL_NONE) {
                return fail(err, "sim_unallocated requires ues");
            }
        }
    }

    assertion_count += s->final_check_count;

    for (i = 0; i < s->final_check_count; i++) {
        const check_spec_t *check;

        check = &s->final_checks[i];
        if (check->type == CHECK_NODE_VERSION_EQUALS &&
            (check->app[0] == '\0' || check->expected[0] == '\0')) {
            return fail(err,
                        "node_version_equals requires app and version/tag");
        }
        if (check->type == CHECK_RELEASE_UNAVAILABLE &&
            (check->app[0] == '\0' || check->expected[0] == '\0')) {
            return fail(err,
                        "release_unavailable requires app and version");
        }
    }


    if (ulab_streq(s->suite, "p0") && ulab_streq(s->status, "active") &&
        assertion_count == 0) {
        return fail(err, "active p0 scenario requires at least one check");
    }

    return ULAB_OK;
}
