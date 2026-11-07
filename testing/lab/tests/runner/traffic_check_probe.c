/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "check.h"
#include "hardware_observe.h"
#include "runtime_hardware.h"
#include "util.h"

int check_runtime(check_ctx_t *, const check_spec_t *, check_result_t *,
                   ulab_error_t *);

static const char *pdns[2];
static int traffic_result;

static int verify_pdns(const selector_result_t *ues, const char *expected,
                        ulab_error_t *err) {
    size_t i;

    for (i = 0; i < ues->count; i++) {
        if (!ulab_streq(pdns[ues->idx[i]], expected)) {
            snprintf(err->msg, sizeof(err->msg), "PDN is %s, expected %s",
                     pdns[ues->idx[i]], expected);
            return ULAB_ERR;
        }
    }

    return ULAB_OK;
}

int runtime_is_hardware(const runtime_t *rt) {
    (void)rt;
    return 0;
}

int hardware_check_runtime(check_ctx_t *ctx, const check_spec_t *check,
                            check_result_t *result, ulab_error_t *err) {
    (void)ctx;
    (void)check;
    (void)result;
    (void)err;
    abort();
}

int runtime_generate_traffic(runtime_t *rt, const world_t *world,
                             const selector_result_t *ues, uint64_t amount,
                             ulab_error_t *err) {
    (void)rt;
    (void)world;
    (void)ues;
    (void)amount;

    if (traffic_result != ULAB_OK) {
        snprintf(err->msg, sizeof(err->msg), "traffic failed");
    }

    return traffic_result;
}

int runtime_verify_ue_pdns(runtime_t *rt, const world_t *world,
                            const selector_result_t *ues,
                            const check_spec_t *check, ulab_error_t *err) {
    (void)rt;
    (void)world;

    assert(check->timeout_seconds == 0);
    return verify_pdns(ues, check->expected, err);
}

int runtime_verify_ue_sessions(runtime_t *rt, const world_t *world,
                                const selector_result_t *ues,
                                ulab_error_t *err) {
    size_t i;

    (void)rt;
    (void)world;

    for (i = 0; i < ues->count; i++) {
        if (!ulab_streq(pdns[ues->idx[i]], "allowed") &&
            !ulab_streq(pdns[ues->idx[i]], "blocked")) {
            snprintf(err->msg, sizeof(err->msg), "no confirmed live PDN");
            return ULAB_ERR;
        }
    }

    return ULAB_OK;
}

int runtime_verify_ue_policy_blocks(runtime_t *rt, const world_t *world,
                                     const selector_result_t *ues,
                                     ulab_error_t *err) {
    (void)rt;
    (void)world;
    return verify_pdns(ues, "blocked", err);
}

int main(int argc, char **argv) {
    ue_t ues[2];
    world_t world;
    runtime_t runtime;
    check_ctx_t ctx;
    check_spec_t check;
    check_result_t result;
    ulab_error_t err;
    int rc;
    int i;

    if (argc < 3 || argc > 5) {
        return 2;
    }

    memset(ues, 0, sizeof(ues));
    memset(&world, 0, sizeof(world));
    memset(&runtime, 0, sizeof(runtime));
    memset(&ctx, 0, sizeof(ctx));
    memset(&check, 0, sizeof(check));
    memset(&result, 0, sizeof(result));
    memset(&err, 0, sizeof(err));

    world.ues      = ues;
    world.ue_count = argc - 3;
    ctx.world     = &world;
    ctx.runtime   = &runtime;

    check.ues.kind = SEL_ALL;
    if (ulab_streq(argv[1], "unavailable")) {
        check.type = CHECK_TRAFFIC_UNAVAILABLE;
    } else if (ulab_streq(argv[1], "blocked")) {
        check.type = CHECK_TRAFFIC_BLOCKED;
    } else if (ulab_streq(argv[1], "allowed")) {
        check.type = CHECK_TRAFFIC_ALLOWED;
    } else {
        return 2;
    }

    traffic_result = ulab_streq(argv[2], "ok") ? ULAB_OK : ULAB_ERR;
    for (i = 3; i < argc; i++) {
        pdns[i - 3] = argv[i];
    }
    setenv("ULAB_CONTROL_WAIT_SEC", "0", 1);

    rc = check_runtime(&ctx, &check, &result, &err);
    printf("rc=%d passed=%d %s %s\n", rc, result.passed,
           result.detail, err.msg);
    return rc != ULAB_OK || !result.passed ? 1 : 0;
}
