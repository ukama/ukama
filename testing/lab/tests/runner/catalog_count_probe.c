/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>

/* Exercise the production count checker without linking unrelated checks. */
#include "../../src/check_console.c"

static int fixture_transport(void *ctx, const char *operation,
                              const char *query, const char *variables,
                              json_t **response, ulab_error_t *err) {
    json_t *fixture = ctx;

    (void)query;
    printf("REQUEST %s %s\n", operation, variables);
    if (json_is_true(json_object_get(fixture, "transport_error"))) {
        snprintf(err->msg, sizeof(err->msg), "fixture transport failed");
        return ULAB_ERR;
    }
    *response = json_incref(fixture);
    return ULAB_OK;
}

static int run_count(check_ctx_t *ctx, const check_spec_t *check,
                      int all_lists, size_t *checked) {
    check_result_t result = {0};
    ulab_error_t err = {0};
    int rc;

    if (check->type != CHECK_LIST_COUNT_EQUALS ||
        (!all_lists && !ulab_streq(check->target, "packages") &&
         !ulab_streq(check->target, "plans"))) {
        return 0;
    }
    (*checked)++;
    rc = check_list_count(ctx, check, &result, &err);
    printf("CHECK rc=%d passed=%d %s %s\n",
           rc, result.passed, result.detail, err.msg);
    return rc != ULAB_OK || !result.passed;
}

int main(int argc, char **argv) {
    scenario_t *scenario;
    world_t world = {0};
    bff_client_t client = {0};
    check_ctx_t ctx = {0};
    ulab_error_t err = {0};
    json_error_t json_err;
    json_t *fixture;
    size_t i;
    size_t p;
    size_t checked = 0;
    int failed = 0;
    int all_lists;

    if (argc != 4) return ULAB_EUSAGE;
    all_lists = ulab_streq(argv[3], "all");
    fixture = json_load_file(argv[2], 0, &json_err);
    if (fixture == NULL) {
        fprintf(stderr, "%s\n", json_err.text);
        return ULAB_ERR;
    }
    scenario = calloc(1, sizeof(*scenario));
    if (scenario == NULL) return ULAB_ERR;
    if (scenario_load(argv[1], scenario, &err) ||
        scenario_validate(scenario, &err) ||
        world_generate(scenario, "catalog-test", &world, &err)) {
        fprintf(stderr, "%s\n", err.msg);
        failed = 1;
        goto done;
    }
    for (i = 0; i < world.network_count; i++) {
        snprintf(world.networks[i].bff_id,
                 sizeof(world.networks[i].bff_id), "network-%zu", i);
    }
    for (i = 0; i < world.package_count; i++) {
        if (!ulab_streq(argv[3], "uncreated")) {
            snprintf(world.packages[i].bff_id,
                     sizeof(world.packages[i].bff_id), "owned-%zu", i);
        }
    }
    client.transport = fixture_transport;
    client.transport_ctx = fixture;
    ctx.scenario = scenario;
    ctx.world = &world;
    ctx.bff = &client;
    for (p = 0; p < scenario->phase_count; p++) {
        for (i = 0; i < scenario->phases[p].check_count; i++) {
            failed |= run_count(&ctx, &scenario->phases[p].checks[i],
                                 all_lists, &checked);
        }
    }
    for (i = 0; i < scenario->final_check_count; i++) {
        failed |= run_count(&ctx, &scenario->final_checks[i],
                             all_lists, &checked);
    }
    if (checked == 0) failed = 1;
done:
    world_free(&world);
    free(scenario);
    json_decref(fixture);
    return failed;
}
