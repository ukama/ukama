/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "selector.h"
#include "util.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <stdint.h>

static int add_idx(selector_result_t *r, size_t idx) {
    size_t *n = realloc(r->idx, (r->count + 1) * sizeof(size_t));
    if (n == NULL) return ULAB_ERR;
    r->idx = n;
    r->idx[r->count++] = idx;
    return ULAB_OK;
}

static void init(selector_result_t *r) {
    r->idx = NULL;
    r->count = 0;
}

static int has_idx(const selector_result_t *r, size_t idx) {
    size_t i;

    for (i = 0; i < r->count; i++) {
        if (r->idx[i] == idx) {
            return 1;
        }
    }
    return 0;
}

void selector_result_free(selector_result_t *r) {
    if (r == NULL) return;
    free(r->idx);
    r->idx = NULL;
    r->count = 0;
}

int selector_resolve_ues(const world_t *w, const selector_t *sel,
                         selector_result_t *out, ulab_error_t *err) {
    size_t i;
    init(out);

    if (sel->kind == SEL_NONE || sel->kind == SEL_ALL) {
        for (i = 0; i < w->ue_count; i++) if (add_idx(out, i)) return ULAB_ERR;
        return ULAB_OK;
    }

    if (sel->kind == SEL_SAMPLE_PER_SITE) {
        size_t s;
        uint32_t pick;

        for (s = 0; s < w->site_count; s++) {
            for (pick = 0; pick < sel->count; pick++) {
                size_t best = w->ue_count;
                uint32_t best_score = UINT32_MAX;

                for (i = 0; i < w->ue_count; i++) {
                    uint32_t score;

                    if (!ulab_streq(w->ues[i].site_ref, w->sites[s].ref)) {
                        continue;
                    }
                    if (has_idx(out, i)) {
                        continue;
                    }

                    score = ulab_hash32(w->ues[i].ref, w->seed ^ pick);
                    if (score < best_score) {
                        best_score = score;
                        best = i;
                    }
                }

                if (best == w->ue_count) {
                    break;
                }
                if (add_idx(out, best)) {
                    return ULAB_ERR;
                }
            }
        }
        return ULAB_OK;
    }

    snprintf(err->msg, sizeof(err->msg), "unsupported UE selector");
    return ULAB_ERR;
}

int selector_resolve_nodes(const world_t *w, const selector_t *sel,
                           selector_result_t *out, ulab_error_t *err) {
    size_t i;
    init(out);
    if (sel->kind == SEL_NONE || sel->kind == SEL_ALL) {
        for (i = 0; i < w->node_count; i++) {
            if (add_idx(out, i)) return ULAB_ERR;
        }
        return ULAB_OK;
    }
    if (sel->kind == SEL_NODE_TYPE ||
        sel->kind == SEL_NODE_TYPE_COUNT_PER_NETWORK) {
        for (i = 0; i < w->node_count; i++) {
            if (ulab_streq(w->nodes[i].type, sel->value)) {
                if (add_idx(out, i)) return ULAB_ERR;
                if (sel->kind == SEL_NODE_TYPE_COUNT_PER_NETWORK &&
                    out->count >= sel->count) break;
            }
        }
        return ULAB_OK;
    }
    snprintf(err->msg, sizeof(err->msg), "unsupported node selector");
    return ULAB_ERR;
}

int selector_resolve_sites(const world_t *w, const selector_t *sel,
                           selector_result_t *out, ulab_error_t *err) {
    size_t i;
    (void)err;
    init(out);
    if (sel->kind == SEL_NONE || sel->kind == SEL_ALL) {
        for (i = 0; i < w->site_count; i++) {
            if (add_idx(out, i)) return ULAB_ERR;
        }
        return ULAB_OK;
    }
    return ULAB_ERR;
}

int selector_resolve_networks(const world_t *w, const selector_t *sel,
                              selector_result_t *out, ulab_error_t *err) {
    size_t i;
    (void)err;
    init(out);
    if (sel->kind == SEL_NONE || sel->kind == SEL_ALL) {
        for (i = 0; i < w->network_count; i++) {
            if (add_idx(out, i)) return ULAB_ERR;
        }
        return ULAB_OK;
    }
    return ULAB_ERR;
}
