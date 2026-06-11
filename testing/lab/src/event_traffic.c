/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "event.h"
#include "selector.h"
#include "util.h"
#include <string.h>

static const profile_spec_t *find_profile(const scenario_t *s,
                                          const char *name) {
    size_t i;
    for (i = 0; i < s->profile_count; i++) {
        if (ulab_streq(s->profiles[i].name, name)) return &s->profiles[i];
    }
    return NULL;
}

static uint64_t profile_amount(const profile_spec_t *p, size_t idx) {
    uint32_t bucket = (uint32_t)(idx % 100);
    uint32_t acc = 0;
    size_t i;

    for (i = 0; i < p->bucket_count; i++) {
        acc += p->buckets[i].percent;
        if (bucket < acc) return p->buckets[i].amount_mb;
    }
    return p->buckets[p->bucket_count - 1].amount_mb;
}

int event_traffic(event_ctx_t *ctx, const event_spec_t *event,
                  ulab_error_t *err) {
    selector_result_t ues;
    size_t i;

    if (selector_resolve_ues(ctx->world, &event->ues, &ues, err)) {
        return ULAB_ERR;
    }
    if (event->type == EVT_TRAFFIC) {
        if (runtime_generate_traffic(ctx->runtime, ctx->world, &ues,
            event->amount_mb, err)) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        for (i = 0; i < ues.count; i++) {
            model_add_usage(ctx->model, ctx->world->ues[ues.idx[i]].ref,
                            event->amount_mb);
        }
    } else {
        const profile_spec_t *p = find_profile(ctx->scenario, event->profile);
        if (p == NULL) {
            snprintf(err->msg, sizeof(err->msg), "unknown profile: %s",
                     event->profile);
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        for (i = 0; i < ues.count; i++) {
            selector_result_t one = {0};
            uint64_t mb = profile_amount(p, i);
            one.idx = &ues.idx[i];
            one.count = 1;
            if (runtime_generate_traffic(ctx->runtime, ctx->world, &one,
                mb, err)) {
                selector_result_free(&ues);
                return ULAB_ERR;
            }
            model_add_usage(ctx->model, ctx->world->ues[ues.idx[i]].ref, mb);
        }
    }
    selector_result_free(&ues);
    return ULAB_OK;
}
