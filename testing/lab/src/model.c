/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "model.h"
#include "util.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int model_init(model_t *m, const world_t *w) {
    return model_sync_world(m, w);
}

void model_free(model_t *m) {
    if (m == NULL) return;
    free(m->ues);
    memset(m, 0, sizeof(*m));
}

int model_sync_world(model_t *m, const world_t *w) {
    size_t old = m->ue_count;
    size_t i;

    if (w->ue_count <= old) {
        return ULAB_OK;
    }
    m->ues = realloc(m->ues, w->ue_count * sizeof(model_ue_t));
    if (m->ues == NULL) {
        return ULAB_ERR;
    }
    for (i = old; i < w->ue_count; i++) {
        const ue_t *ue = &w->ues[i];
        package_t *pkg = world_package_by_ref((world_t *)w, ue->package_ref);
        model_ue_t *mu = &m->ues[i];

        memset(mu, 0, sizeof(*mu));
        snprintf(mu->ue_ref, sizeof(mu->ue_ref), "%s", ue->ref);
        snprintf(mu->sim_id, sizeof(mu->sim_id), "%s", ue->bff_id);
        snprintf(mu->package_ref, sizeof(mu->package_ref), "%s",
                 ue->package_ref);
        mu->package_mb = pkg ? pkg->data_mb : 0;
        mu->remaining_mb = mu->package_mb;
    }
    m->ue_count = w->ue_count;
    return ULAB_OK;
}

model_ue_t *model_ue(model_t *m, const char *ue_ref) {
    size_t i;
    for (i = 0; i < m->ue_count; i++) {
        if (ulab_streq(m->ues[i].ue_ref, ue_ref)) return &m->ues[i];
    }
    return NULL;
}

const model_ue_t *model_ue_const(const model_t *m, const char *ue_ref) {
    size_t i;
    for (i = 0; i < m->ue_count; i++) {
        if (ulab_streq(m->ues[i].ue_ref, ue_ref)) return &m->ues[i];
    }
    return NULL;
}

int model_add_usage(model_t *m, const char *ue_ref, uint64_t amount_mb) {
    model_ue_t *ue = model_ue(m, ue_ref);
    if (ue == NULL) return ULAB_ERR;
    ue->used_mb += amount_mb;
    if (amount_mb >= ue->remaining_mb) {
        ue->remaining_mb = 0;
    } else {
        ue->remaining_mb -= amount_mb;
    }
    return ULAB_OK;
}

uint64_t model_site_usage(const model_t *m, const world_t *w,
                          const char *site_ref) {
    uint64_t total = 0;
    size_t i;
    for (i = 0; i < w->ue_count; i++) {
        if (ulab_streq(w->ues[i].site_ref, site_ref)) {
            const model_ue_t *mu = model_ue_const(m, w->ues[i].ref);
            if (mu) total += mu->used_mb;
        }
    }
    return total;
}

uint64_t model_network_usage(const model_t *m, const world_t *w,
                             const char *network_ref) {
    uint64_t total = 0;
    size_t i;
    for (i = 0; i < w->ue_count; i++) {
        if (ulab_streq(w->ues[i].network_ref, network_ref)) {
            const model_ue_t *mu = model_ue_const(m, w->ues[i].ref);
            if (mu) total += mu->used_mb;
        }
    }
    return total;
}

int model_balance_non_negative(const model_t *m) {
    size_t i;
    for (i = 0; i < m->ue_count; i++) {
        if (m->ues[i].used_mb > m->ues[i].package_mb &&
            m->ues[i].remaining_mb != 0) {
            return 0;
        }
    }
    return 1;
}

int model_write_json(const model_t *m, const char *path) {
    FILE *f = fopen(path, "w");
    size_t i;
    if (!f) return ULAB_ERR;
    fprintf(f, "{\n  \"ues\": [\n");
    for (i = 0; i < m->ue_count; i++) {
        fprintf(f, "    {\"ue\":\"%s\",\"used_mb\":%llu,",
                m->ues[i].ue_ref,
                (unsigned long long)m->ues[i].used_mb);
        fprintf(f, "\"remaining_mb\":%llu}%s\n",
                (unsigned long long)m->ues[i].remaining_mb,
                i + 1 == m->ue_count ? "" : ",");
    }
    fprintf(f, "  ]\n}\n");
    fclose(f);
    return ULAB_OK;
}
