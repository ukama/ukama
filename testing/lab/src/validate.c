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

static int fail(ulab_error_t *err, const char *msg) {
    snprintf(err->msg, sizeof(err->msg), "%s", msg);
    return ULAB_ERR;
}

int scenario_validate(const scenario_t *s, ulab_error_t *err) {
    uint32_t pct = 0;
    size_t i;

    if (s->version != ULAB_SCHEMA_VER) {
        return fail(err, "unsupported scenario version");
    }

    if (s->name[0] == '\0') {
        return fail(err, "missing scenario name");
    }

    if (s->world.networks == 0 || s->world.sites_per_network == 0) {
        return fail(err, "world must include networks and sites");
    }

    if (s->world.ues_per_site == 0) {
        return fail(err, "world.ues_per_site must be > 0");
    }

    if (s->world.tower_per_site + s->world.amplifier_per_site +
        s->world.controller_per_site == 0) {
        return fail(err, "world.nodes_per_site must include nodes");
    }

    if (s->package_count == 0) {
        return fail(err, "at least one package is required");
    }

    for (i = 0; i < s->package_count; i++) {
        const package_spec_t *p = &s->packages[i];
        if (p->ref[0] == '\0' || p->name[0] == '\0') {
            return fail(err, "package ref/name is required");
        }

        if (p->data_mb == 0 || p->duration_days == 0) {
            return fail(err, "package data_mb/duration_days invalid");
        }

        pct += p->assign_percent;
    }

    if (pct != 100) {
        return fail(err, "package assign_percent values must add to 100");
    }

    if (!s->setup.create_networks || !s->setup.create_sites ||
        !s->setup.create_nodes || !s->setup.create_packages ||
        !s->setup.create_subscribers || !s->setup.create_sims) {
        return fail(err, "setup.create_via_bff missing required entries");
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

    return ULAB_OK;
}
