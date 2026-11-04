/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <assert.h>
#include <stdlib.h>

/* Exercise the production checks without linking unrelated runtime actions. */
#include "../../src/check_analytics.c"
#include "../../src/check_usage.c"

static json_t *fixture;

static int transport(void *ctx, const char *operation, const char *query,
                      const char *variables, json_t **response,
                      ulab_error_t *err) {
    (void)ctx;
    (void)operation;
    (void)query;
    (void)variables;

    if (fixture == NULL) {
        snprintf(err->msg, sizeof(err->msg), "transport failed");
        return ULAB_ERR;
    }
    *response = json_incref(fixture);
    return ULAB_OK;
}

static void response(const char *key, const char *scope, double value) {
    json_t *rows;
    json_t *row;

    json_decref(fixture);
    rows = json_array();
    if (key != NULL) {
        row = json_pack("{s:s,s:f,s:s,s:[{s:s,s:s}]}",
                        "kpi", key, "value", value, "unit", "MB", "scope",
                        "key", "network_id", "value", scope);
        assert(row != NULL);
        assert(json_array_append_new(rows, row) == 0);
    }
    fixture = json_pack("{s:{s:{s:o}}}", "data", "getKpiValues", "values", rows);
    assert(fixture != NULL);
}

static void validate_files(int argc, char **argv) {
    scenario_t *scenario;
    ulab_error_t err = {0};
    int i;

    scenario = calloc(1, sizeof(*scenario));
    assert(scenario != NULL);
    for (i = 1; i < argc; i++) {
        memset(scenario, 0, sizeof(*scenario));
        if (scenario_load(argv[i], scenario, &err) ||
            scenario_validate(scenario, &err)) {
            fprintf(stderr, "%s: %s\n", argv[i], err.msg);
            abort();
        }
    }
    free(scenario);
}

int main(int argc, char **argv) {
    network_t network = {.ref = "net-001", .bff_id = "network-1"};
    ue_t ue = {.ref = "ue-001", .network_ref = "net-001"};
    world_t world = {.networks = &network, .network_count = 1,
                     .ues = &ue, .ue_count = 1};
    model_ue_t model_ue_value = {.ue_ref = "ue-001"};
    model_t model = {.ues = &model_ue_value, .ue_count = 1};
    runtime_t runtime = {0};
    bff_client_t client = {.transport = transport};
    check_ctx_t ctx = {.world = &world, .model = &model,
                       .runtime = &runtime, .bff = &client};
    check_spec_t check = {.type = CHECK_KPI_VALUE,
                          .key = "SITE_UPTIME", .span = "daily", .op = "AVG",
                          .scope_key = "network_id", .networks = {.kind = SEL_ALL},
                          .has_expected_value = 1, .expected_value = 0,
                          .timeout_seconds = 1, .poll_seconds = 1};
    check_result_t result = {0};
    ulab_error_t err = {0};

    validate_files(argc, argv);
    response("SITE_UPTIME", "network-1", 0);
    assert(check_kpi(&ctx, &check, &result, &err) == ULAB_OK && result.passed);
    ulab_copy(check.expected, sizeof(check.expected), "previous");
    check.has_expected_value = 0;
    ulab_copy(check.comparator, sizeof(check.comparator), "greater_than");
    response("SITE_UPTIME", "network-1", 25);
    assert(check_kpi(&ctx, &check, &result, &err) == ULAB_OK && result.passed);
    assert(runtime.last_kpi.value == 25);

    /* An unchanged value cannot satisfy recovery, nor can an invalid percent. */
    assert(check_kpi(&ctx, &check, &result, &err) == ULAB_OK && !result.passed);
    response("SITE_UPTIME", "network-1", 101);
    assert(check_kpi(&ctx, &check, &result, &err) == ULAB_OK && !result.passed);
    assert(runtime.last_kpi.value == 25);
    ulab_copy(check.comparator, sizeof(check.comparator), "less_than");
    response("SITE_UPTIME", "network-1", 20);
    assert(check_kpi(&ctx, &check, &result, &err) == ULAB_OK && result.passed);
    ulab_copy(check.comparator, sizeof(check.comparator), "greater_than");
    response("SITE_UPTIME", "network-1", 28);
    assert(check_kpi(&ctx, &check, &result, &err) == ULAB_OK && result.passed);
    ulab_copy(check.scope_value, sizeof(check.scope_value), "different-network");
    assert(check_kpi(&ctx, &check, &result, &err) == ULAB_ERR);
    check.scope_value[0] = '\0';
    runtime.last_kpi.valid = 0;
    assert(check_kpi(&ctx, &check, &result, &err) == ULAB_ERR);

    memset(&check, 0, sizeof(check));
    ulab_copy(check.target, sizeof(check.target), "network");
    check.networks.kind = SEL_ALL;
    check.timeout_seconds = 1;
    check.poll_seconds = 1;
    response(NULL, NULL, 0);
    assert(check_usage_aggregate(&ctx, &check, &result, &err) == ULAB_OK && result.passed);
    response("DATA_USAGE", "network-1", 0);
    assert(check_usage_aggregate(&ctx, &check, &result, &err) == ULAB_OK && result.passed);
    response("DATA_USAGE", "other-network", 0);
    assert(check_usage_aggregate(&ctx, &check, &result, &err) == ULAB_OK && !result.passed);
    response("DATA_USAGE", "network-1", 1);
    assert(check_usage_aggregate(&ctx, &check, &result, &err) == ULAB_OK && !result.passed);
    response(NULL, NULL, 0);
    model_ue_value.used_mb = 1;
    assert(check_usage_aggregate(&ctx, &check, &result, &err) == ULAB_OK && !result.passed);
    model_ue_value.used_mb = 0;
    ulab_copy(check.key, sizeof(check.key), "REVENUE");
    assert(check_usage_aggregate(&ctx, &check, &result, &err) == ULAB_OK && !result.passed);
    check.key[0] = '\0';
    json_decref(fixture);
    fixture = json_object();
    assert(check_usage_aggregate(&ctx, &check, &result, &err) == ULAB_ERR);
    json_decref(fixture);
    fixture = NULL;
    assert(check_usage_aggregate(&ctx, &check, &result, &err) == ULAB_ERR);

    puts("PASS P0: measured uptime, scoped baselines, empty-network usage and errors");
    return 0;
}
