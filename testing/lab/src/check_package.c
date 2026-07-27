/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "check.h"
#include "selector.h"
#include "util.h"
#include <stdio.h>
#include <time.h>
#include <unistd.h>

int check_package(check_ctx_t *ctx, const check_spec_t *check,
                  check_result_t *res, ulab_error_t *err) {
    selector_result_t ues;
    size_t i;
    size_t ok = 0;

    if (check->type == CHECK_BALANCE_NON_NEGATIVE) {
        res->passed = model_balance_non_negative(ctx->model);
        snprintf(res->detail, sizeof(res->detail),
                 "model balance non-negative");
        return ULAB_OK;
    }
    if (check->type == CHECK_PACKAGE_REMAINING) {
        res->skipped = check->required ? 0 : 1;
        res->passed = 0;
        snprintf(res->detail, sizeof(res->detail),
                 "BFF remaining balance query not available yet%s",
                 check->required ? " (required)" : "");
        return ULAB_OK;
    }

    if (check->type == CHECK_PACKAGE_STATE ||
        check->type == CHECK_PACKAGE_ASSIGNMENT_COUNT) {
        time_t deadline;
        unsigned int poll;

        if (selector_resolve_ues(ctx->world, &check->ues, &ues, err)) {
            return ULAB_ERR;
        }
        deadline = time(NULL) + (time_t)(check->timeout_seconds ?
                                        check->timeout_seconds : 300u);
        poll = check->poll_seconds ? check->poll_seconds : 5u;
        do {
            ok = 0;
            for (i = 0; i < ues.count; i++) {
                ue_t *ue;
                package_t *pkg;
                bff_sim_package_t assignments[ULAB_MAX_BFF_SIM_PACKAGES];
                size_t assignment_count;
                size_t j;
                size_t matching;
                int active;

                ue = &ctx->world->ues[ues.idx[i]];
                pkg = check->package_ref[0] ?
                    world_package_for_network(ctx->world, check->package_ref,
                                              ue->network_ref) : NULL;
                if (check->package_ref[0] && pkg == NULL) {
                    continue;
                }
                if (bff_get_sim_packages(ctx->bff, ue, assignments,
                                         ULAB_MAX_BFF_SIM_PACKAGES,
                                         &assignment_count, err)) {
                    selector_result_free(&ues);
                    return ULAB_ERR;
                }

                matching = 0;
                active = 0;
                for (j = 0; j < assignment_count; j++) {
                    if (pkg != NULL &&
                        !ulab_streq(assignments[j].package_id,
                                    pkg->bff_id)) {
                        continue;
                    }
                    matching++;
                    if (assignments[j].active) {
                        active = 1;
                    }
                }

                if (check->type == CHECK_PACKAGE_ASSIGNMENT_COUNT &&
                    matching == check->expected_count) {
                    ok++;
                } else if (check->type == CHECK_PACKAGE_STATE &&
                           ((ulab_streq(check->expected, "active") && active) ||
                            (ulab_streq(check->expected, "queued") &&
                             matching > 0 && !active) ||
                            (ulab_streq(check->expected, "inactive") &&
                             matching > 0 && !active) ||
                            (ulab_streq(check->expected, "absent") &&
                             matching == 0))) {
                    ok++;
                }
            }
            if (ok != ues.count && time(NULL) < deadline) {
                sleep(poll > 60u ? 60u : poll);
            }
        } while (ok != ues.count && time(NULL) < deadline);

        res->passed = ok == ues.count;
        if (check->type == CHECK_PACKAGE_STATE) {
            snprintf(res->detail, sizeof(res->detail),
                     "%s=%zu/%zu expected=%s",
                     scenario_check_name(check->type), ok, ues.count,
                     check->expected);
        } else {
            snprintf(res->detail, sizeof(res->detail),
                     "%s=%zu/%zu expected=%u",
                     scenario_check_name(check->type), ok, ues.count,
                     check->expected_count);
        }
        selector_result_free(&ues);
        return ULAB_OK;
    }
    if (selector_resolve_ues(ctx->world, &check->ues, &ues, err)) {
        return ULAB_ERR;
    }
    for (i = 0; i < ues.count; i++) {
        ue_t *ue = &ctx->world->ues[ues.idx[i]];
        package_t *pkg = check->package_ref[0] ?
            world_package_for_network(ctx->world, check->package_ref,
                                      ue->network_ref) :
            world_package_by_ref(ctx->world, ue->package_ref);
        int active = 0;
        if (pkg == NULL) continue;
        if (bff_get_packages_for_sim(ctx->bff, ue, pkg->bff_id, &active,
            err)) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        if (active) ok++;
    }
    res->passed = ok == ues.count;
    snprintf(res->detail, sizeof(res->detail), "package_active=%zu/%zu", ok,
             ues.count);
    selector_result_free(&ues);
    return ULAB_OK;
}
