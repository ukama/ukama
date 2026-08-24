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
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <time.h>
#include <unistd.h>

static unsigned int env_seconds(const char *name,
                                unsigned int fallback,
                                unsigned int maximum) {
    const char *value;
    char *end;
    unsigned long seconds;

    value = getenv(name);
    if (value == NULL || value[0] == '\0') {
        return fallback;
    }

    end = NULL;
    seconds = strtoul(value, &end, 10);
    if (end == value || *end != '\0' || seconds > maximum) {
        return fallback;
    }

    return (unsigned int)seconds;
}

static unsigned int control_wait_seconds(void) {
    return env_seconds("ULAB_CONTROL_WAIT_SEC", 30u, 300u);
}

static unsigned int software_wait_seconds(void) {
    return env_seconds("ULAB_SOFTWARE_UPDATE_TIMEOUT_SEC", 180u, 1800u);
}

static unsigned int software_poll_seconds(void) {
    unsigned int seconds;

    seconds = env_seconds("ULAB_SOFTWARE_UPDATE_POLL_SEC", 5u, 60u);
    return seconds == 0 ? 1u : seconds;
}

static int app_version_matches(const bff_app_state_t *app,
                               const char *expected) {
    if (app == NULL || expected == NULL || expected[0] == '\0') {
        return 0;
    }

    return ulab_streq(app->version, expected);
}

static int app_is_running(const bff_app_state_t *app) {
    if (app == NULL) {
        return 0;
    }

    return strcasecmp(app->status, "running") == 0 ||
           strcasecmp(app->status, "active") == 0;
}

static int check_release_unavailable(check_ctx_t *ctx,
                                     const check_spec_t *check,
                                     check_result_t *res,
                                     ulab_error_t *err) {
    bff_release_t release;
    int found;

    if (check->app[0] == '\0' || check->expected[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "release_unavailable requires app and version");
        return ULAB_ERR;
    }

    memset(&release, 0, sizeof(release));
    found = 0;
    if (bff_get_release(ctx->bff, check->app, "app", check->expected,
                        &release, &found, err)) {
        return ULAB_ERR;
    }

    res->passed = !found || !release.available;
    snprintf(res->detail, sizeof(res->detail),
             "app=%s version=%s found=%s available=%s",
             check->app, check->expected,
             found ? "true" : "false",
             found && release.available ? "true" : "false");
    return ULAB_OK;
}

static int check_node_version(check_ctx_t *ctx,
                              const check_spec_t *check,
                              check_result_t *res,
                              ulab_error_t *err) {
    selector_result_t sel;
    bff_app_state_t app_state;
    ulab_error_t query_err;
    const char *expected;
    unsigned int timeout;
    unsigned int poll;
    unsigned int elapsed;
    size_t i;
    size_t matched;
    char last_node[ULAB_MAX_ID];
    char last_version[ULAB_MAX_REF];
    char last_tag[ULAB_MAX_REF];
    char last_status[ULAB_MAX_REF];
    char last_error[ULAB_MAX_ERR];

    memset(&sel, 0, sizeof(sel));
    memset(last_node, 0, sizeof(last_node));
    memset(last_version, 0, sizeof(last_version));
    memset(last_tag, 0, sizeof(last_tag));
    memset(last_status, 0, sizeof(last_status));
    memset(last_error, 0, sizeof(last_error));

    expected = check->expected[0] ? check->expected : check->status;
    if (check->app[0] == '\0' || expected == NULL || expected[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "node_version_equals requires app and version/tag");
        return ULAB_ERR;
    }

    if (selector_resolve_nodes(ctx->world, &check->nodes, &sel, err)) {
        return ULAB_ERR;
    }

    if (sel.count == 0) {
        res->passed = 0;
        snprintf(res->detail, sizeof(res->detail),
                 "app=%s expected=%s selected_nodes=0",
                 check->app, expected);
        selector_result_free(&sel);
        return ULAB_OK;
    }

    timeout = software_wait_seconds();
    poll = software_poll_seconds();
    elapsed = 0;

    for (;;) {
        matched = 0;

        for (i = 0; i < sel.count; i++) {
            node_t *node;

            node = &ctx->world->nodes[sel.idx[i]];
            memset(&app_state, 0, sizeof(app_state));
            memset(&query_err, 0, sizeof(query_err));

            ulab_copy(last_node, sizeof(last_node), node->bff_id);
            if (bff_get_node_app(ctx->bff, node, check->app,
                                 &app_state, &query_err)) {
                ulab_copy(last_error, sizeof(last_error), query_err.msg);
                continue;
            }

            last_error[0] = '\0';
            ulab_copy(last_version, sizeof(last_version), app_state.version);
            ulab_copy(last_tag, sizeof(last_tag), app_state.tag);
            ulab_copy(last_status, sizeof(last_status), app_state.status);

            if (app_version_matches(&app_state, expected)) {
                matched++;
            }
        }

        if (matched == sel.count) {
            res->passed = 1;
            snprintf(res->detail, sizeof(res->detail),
                     "app=%s expected=%s matched=%zu/%zu elapsed=%us",
                     check->app, expected, matched, sel.count, elapsed);
            selector_result_free(&sel);
            return ULAB_OK;
        }

        if (elapsed >= timeout) {
            break;
        }

        sleep(poll);
        if (timeout - elapsed < poll) {
            elapsed = timeout;
        } else {
            elapsed += poll;
        }
    }

    res->passed = 0;
    if (last_error[0] != '\0') {
        snprintf(res->detail, sizeof(res->detail),
                 "app=%.128s expected=%.96s matched=%zu/%zu "
                 "node=%.192s timeout=%us error=%.384s",
                 check->app, expected, matched, sel.count, last_node,
                 timeout, last_error);
    } else {
        snprintf(res->detail, sizeof(res->detail),
                 "app=%.128s expected=%.96s matched=%zu/%zu "
                 "node=%.192s version=%.96s tag=%.96s status=%.64s "
                 "timeout=%us",
                 check->app, expected, matched, sel.count, last_node,
                 last_version, last_tag, last_status, timeout);
    }

    selector_result_free(&sel);
    return ULAB_OK;
}

static int check_node_health(check_ctx_t *ctx,
                             const check_spec_t *check,
                             check_result_t *res,
                             ulab_error_t *err) {
    selector_result_t sel;
    bff_node_status_t node_state;
    bff_app_state_t app_state;
    ulab_error_t query_err;
    unsigned int timeout;
    unsigned int poll;
    unsigned int elapsed;
    size_t i;
    size_t healthy;
    char last_node[ULAB_MAX_ID];
    char last_connectivity[ULAB_MAX_REF];
    char last_app_status[ULAB_MAX_REF];
    char last_error[ULAB_MAX_ERR];

    memset(&sel, 0, sizeof(sel));
    memset(last_node, 0, sizeof(last_node));
    memset(last_connectivity, 0, sizeof(last_connectivity));
    memset(last_app_status, 0, sizeof(last_app_status));
    memset(last_error, 0, sizeof(last_error));

    if (selector_resolve_nodes(ctx->world, &check->nodes, &sel, err)) {
        return ULAB_ERR;
    }

    if (sel.count == 0) {
        res->passed = 0;
        snprintf(res->detail, sizeof(res->detail), "selected_nodes=0");
        selector_result_free(&sel);
        return ULAB_OK;
    }

    timeout = software_wait_seconds();
    poll = software_poll_seconds();
    elapsed = 0;

    for (;;) {
        healthy = 0;

        for (i = 0; i < sel.count; i++) {
            node_t *node;
            int node_ok;
            int app_ok;

            node = &ctx->world->nodes[sel.idx[i]];
            memset(&node_state, 0, sizeof(node_state));
            memset(&query_err, 0, sizeof(query_err));
            ulab_copy(last_node, sizeof(last_node), node->bff_id);

            if (bff_get_node_status(ctx->bff, node, &node_state,
                                   &query_err)) {
                ulab_copy(last_error, sizeof(last_error), query_err.msg);
                continue;
            }

            last_error[0] = '\0';
            ulab_copy(last_connectivity, sizeof(last_connectivity),
                      node_state.connectivity);
            node_ok = strcasecmp(node_state.connectivity, "online") == 0;
            app_ok = 1;

            if (check->app[0] != '\0') {
                memset(&app_state, 0, sizeof(app_state));
                memset(&query_err, 0, sizeof(query_err));
                if (bff_get_node_app(ctx->bff, node, check->app,
                                     &app_state, &query_err)) {
                    ulab_copy(last_error, sizeof(last_error), query_err.msg);
                    app_ok = 0;
                } else {
                    ulab_copy(last_app_status, sizeof(last_app_status),
                              app_state.status);
                    app_ok = app_is_running(&app_state);
                }
            }

            if (node_ok && app_ok) {
                healthy++;
            }
        }

        if (healthy == sel.count) {
            res->passed = 1;
            snprintf(res->detail, sizeof(res->detail),
                     "nodes_online=%zu/%zu%s%s elapsed=%us",
                     healthy, sel.count,
                     check->app[0] ? " app=" : "",
                     check->app[0] ? check->app : "", elapsed);
            selector_result_free(&sel);
            return ULAB_OK;
        }

        if (elapsed >= timeout) {
            break;
        }

        sleep(poll);
        if (timeout - elapsed < poll) {
            elapsed = timeout;
        } else {
            elapsed += poll;
        }
    }

    res->passed = 0;
    if (last_error[0] != '\0') {
        snprintf(res->detail, sizeof(res->detail),
                 "healthy=%zu/%zu node=%.192s timeout=%us "
                 "error=%.512s",
                 healthy, sel.count, last_node, timeout, last_error);
    } else {
        snprintf(res->detail, sizeof(res->detail),
                 "healthy=%zu/%zu node=%.192s connectivity=%.64s "
                 "app=%.128s app_status=%.64s timeout=%us",
                 healthy, sel.count, last_node, last_connectivity,
                 check->app[0] ? check->app : "n/a",
                 check->app[0] ? last_app_status : "n/a", timeout);
    }

    selector_result_free(&sel);
    return ULAB_OK;
}

int check_runtime(check_ctx_t *ctx, const check_spec_t *check,
                  check_result_t *res, ulab_error_t *err) {
    selector_result_t sel;
    size_t i;
    size_t ok = 0;

    if (check->type == CHECK_RELEASE_UNAVAILABLE) {
        return check_release_unavailable(ctx, check, res, err);
    }

    if (check->type == CHECK_NODE_VERSION_EQUALS) {
        return check_node_version(ctx, check, res, err);
    }

    if (check->type == CHECK_NODE_HEALTH_OK) {
        return check_node_health(ctx, check, res, err);
    }

    if (check->type == CHECK_UE_ATTACHED) {
        if (selector_resolve_ues(ctx->world, &check->ues, &sel, err)) {
            return ULAB_ERR;
        }
        for (i = 0; i < sel.count; i++) {
            if (ctx->world->ues[sel.idx[i]].attached) ok++;
        }
        res->passed = ok == sel.count;
        snprintf(res->detail, sizeof(res->detail), "attached=%zu/%zu", ok,
                 sel.count);
        selector_result_free(&sel);
        return ULAB_OK;
    }

    if (check->type == CHECK_TRAFFIC_ALLOWED ||
        check->type == CHECK_TRAFFIC_BLOCKED ||
        check->type == CHECK_TRAFFIC_UNAVAILABLE) {
        ulab_error_t tmp;
        ulab_error_t session_err;
        uint64_t amount;
        unsigned int timeout;
        unsigned int attempt;
        time_t deadline;
        int rc;
        int expected;
        int n;

        amount = check->expected_used_mb ? check->expected_used_mb : 1;
        timeout = control_wait_seconds();
        deadline = time(NULL) + (time_t)timeout;
        rc = ULAB_ERR;
        expected = 0;
        attempt = 0;
        memset(&tmp, 0, sizeof(tmp));
        memset(&session_err, 0, sizeof(session_err));

        if (selector_resolve_ues(ctx->world, &check->ues, &sel, err)) {
            return ULAB_ERR;
        }

        if (check->type == CHECK_TRAFFIC_BLOCKED &&
            runtime_verify_ue_sessions(ctx->runtime, ctx->world, &sel,
                                       &session_err) != ULAB_OK) {
            res->passed = 0;
            snprintf(res->detail, sizeof(res->detail),
                     "traffic block precondition failed: %.768s",
                     session_err.msg);
            selector_result_free(&sel);
            return ULAB_OK;
        }

        do {
            attempt++;
            memset(&tmp, 0, sizeof(tmp));
            rc = runtime_generate_traffic(ctx->runtime, ctx->world, &sel,
                                          amount, &tmp);
            if (check->type == CHECK_TRAFFIC_ALLOWED) {
                expected = rc == ULAB_OK;
            } else {
                expected = rc != ULAB_OK;
            }

            if (expected || time(NULL) >= deadline) {
                break;
            }
            sleep(1);
        } while (1);

        if (expected && check->type == CHECK_TRAFFIC_BLOCKED) {
            memset(&session_err, 0, sizeof(session_err));
            if (runtime_verify_ue_policy_blocks(ctx->runtime, ctx->world,
                                                &sel, &session_err) != ULAB_OK) {
                expected = 0;
            }
        }

        res->passed = expected;
        n = snprintf(res->detail, sizeof(res->detail),
                     "ues=%zu amount_mb=%llu runtime_rc=%d attempts=%u",
                     sel.count, (unsigned long long)amount, rc,
                     attempt);
        if (n > 0 && (size_t)n < sizeof(res->detail) && tmp.msg[0]) {
            n += snprintf(res->detail + n,
                          sizeof(res->detail) - (size_t)n,
                          " traffic_error=%.256s", tmp.msg);
        }
        if (!expected && check->type == CHECK_TRAFFIC_BLOCKED &&
            n > 0 && (size_t)n < sizeof(res->detail) && session_err.msg[0]) {
            snprintf(res->detail + n,
                     sizeof(res->detail) - (size_t)n,
                     " session_error=%.256s", session_err.msg);
        }

        selector_result_free(&sel);
        return ULAB_OK;
    }

    if (selector_resolve_nodes(ctx->world, &check->nodes, &sel, err)) {
        return ULAB_ERR;
    }
    if (check->type == CHECK_NODE_READY) {
        res->passed = 1;
        snprintf(res->detail, sizeof(res->detail),
                 "runtime nodes=%zu", sel.count);
        selector_result_free(&sel);
        return ULAB_OK;
    }
    for (i = 0; i < sel.count; i++) {
        bff_node_status_t st = {0};
        if (bff_get_node_status(ctx->bff, &ctx->world->nodes[sel.idx[i]],
            &st, err)) {
            selector_result_free(&sel);
            return ULAB_ERR;
        }
        if (check->expected[0] == '\0' || ulab_streq(st.state,
            check->expected)) {
            ok++;
        }
    }
    res->passed = ok == sel.count;
    snprintf(res->detail, sizeof(res->detail), "node_state=%zu/%zu", ok,
             sel.count);
    selector_result_free(&sel);
    return ULAB_OK;
}
