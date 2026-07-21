/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <math.h>
#include <stdio.h>
#include <stdlib.h>

#include "check.h"
#include "selector.h"
#include "util.h"

static int payment_status_matches(const char *actual, const char *expected) {
    if (expected == NULL || expected[0] == '\0') {
        return 1;
    }
    if (ulab_streq(expected, "settled")) {
        return ulab_streq(actual, "completed") ||
            ulab_streq(actual, "success");
    }
    return ulab_streq(actual, expected);
}

static int payment_matches(const bff_payment_t *payment,
                           const ue_t *ue,
                           const subscriber_t *subscriber,
                           const check_spec_t *check) {
    double actual_amount;
    double tolerance;

    if (ue->last_payment_id[0] != '\0' &&
        !ulab_streq(payment->id, ue->last_payment_id)) {
        return 0;
    }
    if (ue->last_payment_id[0] == '\0' && subscriber != NULL &&
        subscriber->email[0] != '\0' && payment->payer_email[0] != '\0' &&
        !ulab_streq(payment->payer_email, subscriber->email)) {
        return 0;
    }
    if (!payment_status_matches(payment->status, check->status)) {
        return 0;
    }
    if (check->currency[0] != '\0' &&
        !ulab_streq(payment->currency, check->currency)) {
        return 0;
    }
    if (!check->has_expected_value) {
        return 1;
    }

    actual_amount = strtod(payment->amount, NULL);
    tolerance = check->tolerance_value > 0 ?
        check->tolerance_value : 0.000001;
    return fabs(actual_amount - check->expected_value) <= tolerance;
}

int check_payment(check_ctx_t *ctx, const check_spec_t *check,
                  check_result_t *res, ulab_error_t *err) {
    selector_result_t ues;
    size_t i;
    size_t ok;
    size_t observed_total;

    if (selector_resolve_ues(ctx->world, &check->ues, &ues, err)) {
        return ULAB_ERR;
    }

    ok = 0;
    observed_total = 0;
    for (i = 0; i < ues.count; i++) {
        ue_t *ue;
        package_t *pkg;
        subscriber_t *subscriber;
        bff_payment_t payments[ULAB_MAX_BFF_PAYMENTS];
        size_t payment_count;
        size_t j;
        size_t matching;

        ue = &ctx->world->ues[ues.idx[i]];
        pkg = check->package_ref[0] ?
            world_package_for_network(ctx->world, check->package_ref,
                                      ue->network_ref) :
            world_package_by_ref(ctx->world, ue->package_ref);
        subscriber = world_subscriber_by_ref(ctx->world,
                                             ue->subscriber_ref);
        if (pkg == NULL) {
            continue;
        }
        if (bff_get_package_payments(ctx->bff, pkg, payments,
                                     ULAB_MAX_BFF_PAYMENTS,
                                     &payment_count, err)) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }

        matching = 0;
        for (j = 0; j < payment_count; j++) {
            if (check->type == CHECK_PAYMENT_COUNT) {
                if (payment_status_matches(payments[j].status,
                                           check->status)) {
                    matching++;
                }
            } else if (payment_matches(&payments[j], ue, subscriber, check)) {
                matching++;
            }
        }
        observed_total += matching;
        if ((check->type == CHECK_PAYMENT_COUNT &&
             matching == check->expected_count) ||
            (check->type == CHECK_PAYMENT_EQUALS && matching > 0)) {
            ok++;
        }
    }

    res->passed = ok == ues.count;
    snprintf(res->detail, sizeof(res->detail),
             "%s=%zu/%zu observed=%zu expected_count=%u",
             scenario_check_name(check->type), ok, ues.count,
             observed_total, check->expected_count);
    selector_result_free(&ues);
    return ULAB_OK;
}
