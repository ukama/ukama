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
        (!s->setup.create_packages || !s->setup.create_subscribers ||
         !s->setup.create_sims)) {
        return fail(err,
                    "UE scenarios require packages, subscribers, and sims");
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
        for (j = 0; j < phase->event_count; j++) {
            const event_spec_t *event;

            event = &phase->events[j];
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

            if (event->type == EVT_SET_PACKAGE_ACTIVE) {
                if (event->package_ref[0] == '\0' || !event->has_active) {
                    return fail(err,
                                "set_package_active requires package and "
                                "active");
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
                 check->type == CHECK_KPI_TREND) &&
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
        }
    }

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

    return ULAB_OK;
}
