/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "scenario.h"
#include "util.h"
#include <ctype.h>
#include <math.h>
#include <stdio.h>
#include <string.h>

static int fail(ulab_error_t *err, const char *message) {
    snprintf(err->msg, sizeof(err->msg), "%s", message);
    return ULAB_ERR;
}

int scenario_is_web_event(event_type_t type) {
    return type >= EVT_WEB_OPEN && type <= EVT_WEB_COMMERCE;
}

int scenario_is_web_check(check_type_t type) {
    return type >= CHECK_WEB_KPI_EQUALS && type <= CHECK_WEB_COMMERCE_EQUALS;
}

int scenario_has_webapp(const scenario_t *s) {
    size_t i;
    size_t j;
    if (s->webapp.present || s->setup.has_webapp_entities) return 1;
    for (i = 0; i < s->phase_count; i++) {
        for (j = 0; j < s->phases[i].event_count; j++)
            if (scenario_is_web_event(s->phases[i].events[j].type)) return 1;
        for (j = 0; j < s->phases[i].check_count; j++)
            if (scenario_is_web_check(s->phases[i].checks[j].type)) return 1;
    }
    for (i = 0; i < s->final_check_count; i++)
        if (scenario_is_web_check(s->final_checks[i].type)) return 1;
    return 0;
}

int scenario_execution_supported(const scenario_t *s, ulab_error_t *err) {
    size_t p;
    size_t i;
    if (s->version != ULAB_WEBAPP_SCHEMA_VER) return ULAB_OK;
    if (ulab_streq(s->status, "xfail"))
        return fail(err, "webapp xfail execution is unsupported; infrastructure failures must remain failures");
    if (ulab_streq(s->status, "wip") || ulab_streq(s->status, "skip")) return ULAB_OK;
    for (p = 0; p < s->phase_count; p++) {
        for (i = 0; i < s->phases[p].event_count; i++) {
            const event_spec_t *event;
            event = &s->phases[p].events[i];
            if (ulab_streq(event->view, "welcome") || ulab_streq(event->view, "unauthorized"))
                return fail(err, "welcome/unauthorized browser handlers are not implemented");
        }
    }
    return ULAB_OK;
}

/* These are semantic browser views, not BFF query/view names. */
static int view_valid(const char *view) {
    static const char *const views[] = {
        "business_home", "business_revenue", "business_customers",
        "business_packages", "business_data_plans", "business_members",
        "business_sim_pool", "business_support", "business_settings",
        "network_home", "network_sites", "network_site_detail",
        "network_nodes", "network_node_detail", "network_customers",
        "network_node_pool", "network_sim_pool", "network_support",
        "network_settings", "customer_customers", "customer_data_plans",
        "customer_settings", "welcome", "unauthorized"
    };
    size_t i;
    for (i = 0; i < sizeof(views) / sizeof(views[0]); i++)
        if (ulab_streq(view, views[i])) return 1;
    return 0;
}

static int view_needs_network(const char *view) {
    return !ulab_streq(view, "welcome") && !ulab_streq(view, "unauthorized") &&
        !ulab_ends(view, "_settings") && !ulab_ends(view, "_members") &&
        !ulab_ends(view, "_sim_pool") && !ulab_ends(view, "_node_pool") &&
        !ulab_streq(view, "business_data_plans");
}

static int numbered_ref(const char *ref, const char *prefix, uint32_t max,
                        uint32_t *number) {
    char canonical[ULAB_MAX_REF];
    uint32_t value;
    if (!ulab_starts(ref, prefix) ||
        ulab_parse_u32(ref + strlen(prefix), &value) || value == 0 || value > max)
        return 0;
    snprintf(canonical, sizeof(canonical), "%s%03u", prefix, value);
    if (!ulab_streq(ref, canonical)) return 0;
    if (number != NULL) *number = value;
    return 1;
}

static int one_network(const scenario_t *s, const selector_t *selector) {
    return selector->kind == SEL_REF &&
        numbered_ref(selector->value, "net-", s->world.networks, NULL);
}

static int detail_reference(const scenario_t *s, const event_spec_t *event) {
    uint32_t network;
    uint32_t site;
    uint32_t total;
    char expected[ULAB_MAX_REF];
    const char *type;
    const char *types[] = {"tower", "amplifier", "controller"};
    size_t i;

    if (!numbered_ref(event->networks.value, "net-", s->world.networks,
                      &network)) return 0;
    total = s->world.networks * s->world.sites_per_network;
    if (ulab_streq(event->view, "network_site_detail")) {
        if (event->sites.kind != SEL_REF || event->nodes.kind != SEL_NONE ||
            !numbered_ref(event->sites.value, "site-", total, &site)) return 0;
    } else {
        if (event->nodes.kind != SEL_REF || event->sites.kind != SEL_NONE)
            return 0;
        type = NULL;
        for (i = 0; i < sizeof(types) / sizeof(types[0]); i++) {
            snprintf(expected, sizeof(expected), "%s-site-", types[i]);
            if (ulab_starts(event->nodes.value, expected)) {
                type = types[i];
                break;
            }
        }
        if (type == NULL) return 0;
        /* Current virtual site topology has one node of each type. */
        for (site = 1; site <= total; site++) {
            snprintf(expected, sizeof(expected), "%s-site-%03u-001", type, site);
            if (ulab_streq(event->nodes.value, expected)) break;
        }
        if (site > total) return 0;
    }
    return s->world.sites_per_network > 0 &&
        (site - 1) / s->world.sites_per_network + 1 == network;
}

static int operation_event(const event_spec_t *e, ulab_error_t *err) {
    int site = ulab_streq(e->view, "network_site_detail");
    int node = ulab_streq(e->view, "network_node_detail");
    int fill = ulab_streq(e->target, "fill_confirmation");
    int toggle = ulab_streq(e->target, "set_radio") || ulab_streq(e->target, "set_service");
    int update = ulab_streq(e->target, "update_software") || ulab_streq(e->target, "retry_update");
    int common = ulab_streq(e->target, "open_restart") || ulab_streq(e->target, "confirm_restart") || ulab_streq(e->target, "cancel_dialog");
    if (!(site || node) || !(common || (site && (fill || toggle ||
        ulab_streq(e->target, "open_site_actions") || ulab_streq(e->target, "close_site_actions"))) ||
        (node && (update || ulab_streq(e->target, "open_software")))))
        return fail(err, "web_action requires a supported operation on its matching detail view");
    if (fill) {
        if ((!!(e->web_fields & (1u << 6)) + !!(e->web_fields & (1u << 7))) != 1 ||
            ((e->web_fields & (1u << 7)) && !ulab_streq(e->variant, "site_name")))
            return fail(err, "fill_confirmation requires value or value_from: site_name");
    } else if (toggle) {
        if (!ulab_streq(e->status, "on") && !ulab_streq(e->status, "off"))
            return fail(err, "radio/service value must be on or off");
        if (e->variant[0]) return fail(err, "value_from is only valid for fill_confirmation");
    } else if (e->web_fields & ((1u << 6) | (1u << 7)))
        return fail(err, "value/value_from are not valid for this action");
    if (update) {
        if (!e->app[0] || !e->tag[0]) return fail(err, "software action requires app and expected target tag");
    } else if (e->web_fields & ((1u << 8) | (1u << 9)))
        return fail(err, "app/tag are only valid for software update actions");
    return ULAB_OK;
}

static int commerce_ue(const scenario_t *s, const selector_t *sel) {
    uint32_t total = s->world.networks * (s->world.sites_per_network * s->world.ues_per_site + s->world.sims_per_network);
    uint32_t n;
    char expected[ULAB_MAX_REF];
    if (sel->kind != SEL_REF || !ulab_starts(sel->value, "ue-") || ulab_parse_u32(sel->value + 3, &n) || !n || n > total) return 0;
    snprintf(expected, sizeof(expected), "ue-%06u", n);
    return ulab_streq(expected, sel->value);
}
static const package_spec_t *commerce_package(const scenario_t *s, const char *ref) {
    size_t i;
    uint32_t n;
    char expected[ULAB_MAX_REF];
    for (i = 0; i < s->package_count; i++) {
        const package_spec_t *p = &s->packages[i];
        if (ulab_streq(p->scope, "organization")) {
            snprintf(expected, sizeof(expected), "%.96s__org", p->ref);
            if (ulab_streq(ref, expected)) return p;
        } else for (n = 1; n <= s->world.networks; n++) {
            char net[ULAB_MAX_REF];
            snprintf(net, sizeof(net), "net-%03u", n);
            if (p->network_ref[0] && !ulab_streq(p->network_ref, net)) continue;
            snprintf(expected, sizeof(expected), "%.96s__%.24s", p->ref, net);
            if (ulab_streq(ref, expected)) return p;
        }
    }
    return NULL;
}
static int commerce_event(const scenario_t *s, const event_spec_t *e, ulab_error_t *err) {
    int plan = ulab_streq(e->target, "create_plan") || ulab_streq(e->target, "edit_plan");
    int combined = ulab_streq(e->target, "allocate_sim") || ulab_streq(e->target, "top_up") || ulab_streq(e->target, "cancel_top_up") || ulab_streq(e->target, "open_receipt");
    int customer = ulab_streq(e->target, "create_customer") || ulab_streq(e->target, "open_customer") || ulab_streq(e->target, "close_customer") || ulab_streq(e->target, "activate_sim") || ulab_streq(e->target, "deactivate_sim");
    int close = ulab_streq(e->target, "close_dialog");
    const package_spec_t *p = e->package_ref[0] ? commerce_package(s, e->package_ref) : NULL;
    if (!(plan || combined || customer || close) || !one_network(s, &e->networks))
        return fail(err, "web_commerce requires a supported action and one world network");
    if (!ulab_streq(e->view, plan ? "business_data_plans" : close ? e->view : "customer_customers") ||
        (close && !ulab_streq(e->view, "business_data_plans") && !ulab_streq(e->view, "customer_customers")))
        return fail(err, "commerce action is on the wrong view");
    if (!!e->package_ref[0] != !!(plan || combined) || ((plan || combined) && !p) ||
        ((customer || combined) ? !commerce_ue(s, &e->ues) : e->ues.kind != SEL_NONE))
        return fail(err, "commerce action requires matching package/ues references");
    if (ulab_streq(e->target, "create_plan")) {
        if (!ulab_streq(e->variant, "MB") && !ulab_streq(e->variant, "GB")) return fail(err, "create_plan requires unit: MB|GB");
        if (ulab_streq(e->variant, "GB") && p->data_mb % 1024) return fail(err, "GB plans require data_mb divisible by 1024; use MB otherwise");
    } else if (e->variant[0]) return fail(err, "unit is only valid for create_plan");
    return ULAB_OK;
}

static int browser_event(const scenario_t *s, const event_spec_t *event,
                         ulab_error_t *err) {
    if (event->expect_result[0] || event->error_contains[0])
        return fail(err, "web actions cannot mask execution failures; "
                    "assert the visible validation/error state instead");
    if (event->timeout_seconds == 0 || event->timeout_seconds > 900 ||
        event->timeout_seconds > s->webapp.scenario_timeout_seconds)
        return fail(err, "web action timeout must be 1..900 and fit scenario timeout");
    if (event->type == EVT_WEB_COMMERCE) return commerce_event(s, event, err);
    if (event->type == EVT_WEB_RELOAD) return ULAB_OK;
    if (event->type == EVT_WEB_TAB) {
        if (!ulab_streq(event->profile, "primary") && !ulab_streq(event->profile, "secondary"))
            return fail(err, "web_tab requires primary or secondary");
        return ULAB_OK;
    }
    if (event->type == EVT_WEB_ACTION && operation_event(event, err)) return ULAB_ERR;
    if (event->type == EVT_WEB_SELECT_NETWORK) {
        if (!one_network(s, &event->networks))
            return fail(err, "web_select_network requires one existing net-NNN reference");
        return ULAB_OK;
    }
    if (!view_valid(event->view)) return fail(err, "web_open has an unknown view");
    if ((view_needs_network(event->view) || event->networks.kind != SEL_NONE) &&
        !one_network(s, &event->networks))
        return fail(err, "web_open requires one existing net-NNN reference for this view");
    if (ulab_streq(event->view, "network_site_detail") ||
        ulab_streq(event->view, "network_node_detail")) {
        if (!detail_reference(s, event))
            return fail(err, "web_open detail reference is missing or outside its network");
    } else if (event->sites.kind != SEL_NONE || event->nodes.kind != SEL_NONE) {
        return fail(err, "sites/nodes selectors are only valid for matching detail views");
    }
    return ULAB_OK;
}

static int node_selector(const scenario_t *s, const selector_t *selector) {
    char expected[ULAB_MAX_REF];
    const char *types[] = {"tower", "amplifier", "controller"};
    uint32_t total = s->world.networks * s->world.sites_per_network;
    uint32_t site;
    size_t i;
    if (total == 0) return 0;
    if (selector->kind == SEL_ALL) return 1;
    if (selector->kind != SEL_REF) return 0;
    for (site = 1; site <= total; site++) {
        for (i = 0; i < sizeof(types) / sizeof(types[0]); i++) {
            snprintf(expected, sizeof(expected), "%s-site-%03u-001", types[i], site);
            if (ulab_streq(selector->value, expected)) return 1;
        }
    }
    return 0;
}

static int requirement_valid(const char *id) {
    const unsigned char *p;
    if (!ulab_starts(id, "WEB-") || strlen(id) < 5) return 0;
    for (p = (const unsigned char *)id; *p; p++)
        if (!isupper(*p) && !isdigit(*p) && *p != '-') return 0;
    return 1;
}

static int commerce_check_scope(const check_spec_t *c, ulab_error_t *err) {
    static const char *const plans[] = {"Plan terms", "Plan price", "Plan scope", "Validity", "Price", "Data volume", "Unit"};
    static const char *const customer[] = {"ICCID", "SIM status", "Phone", "Active plan", "Cycle usage", "Total usage", "Receipt total", "Receipt payment ID", "Receipt method", "Receipt status", "Receipt empty"};
    static const char *const scoped[] = {"Package count", "Package status", "Package dates", "Package days", "Receipt plan"};
    size_t i;
    int valid = 0;
    if (ulab_streq(c->view, "business_data_plans")) {
        for (i = 0; i < sizeof(plans) / sizeof(plans[0]); i++) if (ulab_streq(c->label, plans[i])) valid = 1;
        if (!c->package_ref[0] || c->ues.kind != SEL_NONE) valid = 0;
    } else if (ulab_streq(c->view, "business_packages")) {
        valid = c->package_ref[0] && c->ues.kind == SEL_NONE &&
            (ulab_streq(c->label, "Performance price") || ulab_streq(c->label, "Performance sold") ||
             ulab_streq(c->label, "Performance revenue") || ulab_streq(c->label, "Performance share"));
    } else if (ulab_streq(c->view, "business_sim_pool")) {
        valid = ulab_streq(c->label, "Pool status") && c->ues.kind == SEL_REF && !c->package_ref[0];
    } else if (ulab_streq(c->view, "customer_customers")) {
        for (i = 0; i < sizeof(customer) / sizeof(customer[0]); i++) if (ulab_streq(c->label, customer[i])) valid = 1;
        for (i = 0; i < sizeof(scoped) / sizeof(scoped[0]); i++) if (ulab_streq(c->label, scoped[i]) && c->package_ref[0]) valid = 1;
        if (c->ues.kind != SEL_REF) valid = 0;
    }
    if ((ulab_streq(c->key, "iccid") && !ulab_streq(c->label, "ICCID")) ||
        (ulab_streq(c->key, "payment_id") && !ulab_streq(c->label, "Receipt payment ID")) ||
        (ulab_streq(c->key, "plan_name") && !ulab_streq(c->label, "Active plan") && !ulab_streq(c->label, "Receipt plan"))) valid = 0;
    return valid ? ULAB_OK : fail(err, "commerce label, view and resource selectors do not match");
}

static int browser_check(const scenario_t *s, const check_spec_t *check,
                         ulab_error_t *err) {
    if (!scenario_is_web_check(check->type))
        return fail(err, "webapp acceptance checks must use web_*; BFF checks belong to v1");
    if (!view_valid(check->view) || !check->label[0])
        return fail(err, "web checks require a known view and a visible label");
    if (!requirement_valid(check->requirement))
        return fail(err, "web checks require a WEB-* requirement identifier");
    if (!check->immediate)
        return fail(err, "web checks must run immediately at their declared phase");
    if (check->timeout_seconds == 0 || check->timeout_seconds > 900 ||
        check->timeout_seconds > s->webapp.scenario_timeout_seconds)
        return fail(err, "web check timeout must be 1..900 and fit scenario timeout");
    if (check->type == CHECK_WEB_COMMERCE_EQUALS) {
        if (commerce_check_scope(check, err)) return ULAB_ERR;
        if (!ulab_streq(check->view, "business_data_plans") && !ulab_streq(check->view, "customer_customers") && !ulab_streq(check->view, "business_sim_pool") && !ulab_streq(check->view, "business_packages"))
            return fail(err, "unsupported commerce check view");
        if (check->package_ref[0] && !commerce_package(s, check->package_ref)) return fail(err, "unknown commerce package reference");
        if (check->ues.kind != SEL_NONE && !commerce_ue(s, &check->ues)) return fail(err, "unknown commerce UE reference");
        if ((!!(check->web_fields & (1u << 4)) + !!check->key[0]) != 1)
            return fail(err, "commerce check requires expected or expected_property exclusively");
        if (check->key[0] && !ulab_streq(check->key, "iccid") && !ulab_streq(check->key, "payment_id") && !ulab_streq(check->key, "plan_name"))
            return fail(err, "commerce expected_property must be iccid/payment_id/plan_name");
        if ((ulab_streq(check->key, "plan_name") && !check->package_ref[0]) ||
            ((ulab_streq(check->key, "iccid") || ulab_streq(check->key, "payment_id")) && check->ues.kind == SEL_NONE))
            return fail(err, "commerce expected_property requires its resource selector");
        return ULAB_OK;
    }
    if (check->web_fields & ((1u << 7) | (1u << 8))) {
        selector_t sel;
        memset(&sel, 0, sizeof(sel)); sel.kind = SEL_REF;
        ulab_copy(sel.value, sizeof(sel.value), check->ref);
        if (check->type != CHECK_WEB_FIELD_EQUALS || (check->web_fields & (1u << 4)) ||
            !node_selector(s, &sel) || (!ulab_streq(check->key, "id") && !ulab_streq(check->key, "model") && !ulab_streq(check->key, "site_name")))
            return fail(err, "expected_ref/property requires an existing node, id/model/site_name, and no literal expected");
    }
    if ((check->web_fields & (1u << 11)) && (!ulab_streq(check->variant, "contains") ||
        check->type != CHECK_WEB_FIELD_EQUALS || !ulab_ends(check->label, "reason") ||
        !check->expected[0] || check->ref[0]))
        return fail(err, "match: contains is limited to nonempty literal reason checks");
    if ((!check->app[0] && (ulab_streq(check->label, "Software status") ||
        ulab_streq(check->label, "Current version") || ulab_streq(check->label, "Target version") ||
        ulab_streq(check->label, "Update Now") || ulab_streq(check->label, "Retry update"))))
        return fail(err, "software fields/actions require app");
    if ((check->web_fields & (1u << 10)) && (!check->app[0]))
        return fail(err, "app cannot be empty");
    if (check->app[0] && (!ulab_streq(check->view, "network_node_detail") ||
        !(ulab_streq(check->label, "Software status") || ulab_streq(check->label, "Current version") ||
          ulab_streq(check->label, "Target version") || ulab_streq(check->label, "Update Now") ||
          ulab_streq(check->label, "Retry update"))))
        return fail(err, "app scopes only software fields/actions on node detail");
    if (check->nodes.kind != SEL_NONE && (check->type != CHECK_WEB_FIELD_EQUALS ||
        !ulab_streq(check->view, "network_nodes") || check->nodes.kind != SEL_REF || !node_selector(s, &check->nodes)))
        return fail(err, "web field nodes selector requires one node on network_nodes");
    if ((check->type == CHECK_WEB_KPI_EQUALS || check->type == CHECK_WEB_FIELD_EQUALS) &&
        !(check->web_fields & (1u << 4)) && !check->ref[0])
        return fail(err, "web text checks require expected (including an explicit empty string)");
    if (check->type == CHECK_WEB_TABLE_COUNT_EQUALS && !check->has_expected_count)
        return fail(err, "web_table_count_equals requires expected_count");
    if (check->type == CHECK_WEB_ACTION_AVAILABLE && !check->has_expected_value)
        return fail(err, "web_action_available requires available: true|false");
    return ULAB_OK;
}

int scenario_webapp_validate(const scenario_t *s, ulab_error_t *err) {
    size_t i;
    size_t j;
    size_t k;
    size_t checks = s->final_check_count;
    const webapp_spec_t *w = &s->webapp;
    const char *host;
    const unsigned char *p;
    const event_spec_t *event;
    const phase_spec_t *phase;

    if (!s->name[0] || !ulab_streq(s->suite, "webapp"))
        return fail(err, "version 2 requires a name and suite: webapp");
    if (!ulab_streq(s->status, "active") && !ulab_streq(s->status, "wip") &&
        !ulab_streq(s->status, "skip") && !ulab_streq(s->status, "xfail"))
        return fail(err, "scenario status must be active/wip/skip/xfail");
    if (!ulab_streq(s->priority, "p0") && !ulab_streq(s->priority, "p1") &&
        !ulab_streq(s->priority, "p2"))
        return fail(err, "webapp priority must be p0/p1/p2");
    if (!w->present) return fail(err, "version 2 requires a webapp block");
    host = ulab_starts(w->base_url, "http://") ? w->base_url + 7 :
        (ulab_starts(w->base_url, "https://") ? w->base_url + 8 : NULL);
    if (host == NULL || !*host || *host == '/' || *host == ':' ||
        strpbrk(host, "@?#") != NULL)
        return fail(err, "webapp.base_url requires an http(s) URL without credentials/query/fragment");
    for (p = (const unsigned char *)host; *p; p++)
        if (isspace(*p) || iscntrl(*p)) return fail(err, "invalid webapp.base_url");
    if (!ulab_streq(w->browser, "chromium") && !ulab_streq(w->browser, "firefox") &&
        !ulab_streq(w->browser, "webkit"))
        return fail(err, "webapp.browser must be chromium/firefox/webkit");
    if (!w->auth_state[0]) return fail(err, "webapp.auth_state is required");
    if (w->scenario_timeout_seconds == 0 || w->scenario_timeout_seconds > 86400 ||
        w->action_timeout_seconds == 0 || w->action_timeout_seconds > 900 ||
        w->check_timeout_seconds == 0 || w->check_timeout_seconds > 900 ||
        w->action_timeout_seconds > w->scenario_timeout_seconds ||
        w->check_timeout_seconds > w->scenario_timeout_seconds)
        return fail(err, "invalid webapp deadlines: steps 1..900s, scenario 1..86400s");
    if (!ulab_streq(s->provider.type, "virtual"))
        return fail(err, "webapp foundation currently supports provider: virtual");
    if (s->setup.create_networks || s->setup.create_sites || s->setup.create_nodes ||
        s->setup.create_node_site_links || s->setup.create_packages ||
        s->setup.create_subscribers || s->setup.create_sims)
        return fail(err, "webapp scenarios must provision operator resources via create_via_webapp");
    if (s->profile_count || s->runtime.start_ues || s->runtime.wait_ues_attached || s->world.sims_per_subscriber > 1)
        return fail(err, "webapp commerce uses one SIM per customer and explicit start_ues events after UI allocation");
    if ((s->world.ues_per_site && !s->world.sites_per_network) ||
        ((s->world.ues_per_site || s->world.sims_per_network) && !s->world.networks) ||
        s->world.ues_per_site > 100 || s->world.sims_per_network > 100 ||
        (uint64_t)s->world.networks * ((uint64_t)s->world.sites_per_network * s->world.ues_per_site + s->world.sims_per_network) > 100)
        return fail(err, "invalid webapp SIM/UE topology (limit 100 per site/network)");
    for (i = 0; i < s->package_count; i++) {
        const package_spec_t *p = &s->packages[i];
        uint32_t days = p->duration_minutes ? p->duration_minutes / 1440 : p->duration_days;
        if (!p->ref[0] || !p->name[0] || !p->data_mb || p->data_mb > 1048576 ||
            !isfinite(p->amount) || p->amount <= 0 || p->amount > 1000000 ||
            !p->currency[0] || !p->country[0] || !p->active ||
            (p->duration_days && p->duration_minutes) || p->duration_minutes % 1440 ||
            (days != 1 && days != 7 && days != 30))
            return fail(err, "webapp plans require positive amount/data and exactly 1/7/30 days or 1440/10080/43200 minutes");
        if (!ulab_streq(p->scope, "organization") && !ulab_streq(p->scope, "network"))
            return fail(err, "plan scope must be organization or network");
        if (!s->world.networks || (p->network_ref[0] &&
            (ulab_streq(p->scope, "organization") || !numbered_ref(p->network_ref, "net-", s->world.networks, NULL))))
            return fail(err, "plan scope requires existing world networks");
        for (j = 0; j < i; j++) if (ulab_streq(p->ref, s->packages[j].ref)) return fail(err, "duplicate package reference");
    }
    if (s->world.networks > 100 || s->world.sites_per_network > 100)
        return fail(err, "webapp foundation world limit is 100 networks and 100 sites per network");
    if (!!s->world.networks != !!(s->setup.webapp_entities & WEB_SETUP_NETWORKS))
        return fail(err, "world networks and create_via_webapp networks must agree");
    if (!!s->world.sites_per_network != !!(s->setup.webapp_entities & WEB_SETUP_SITES))
        return fail(err, "world sites and create_via_webapp sites must agree");
    if (s->world.sites_per_network > 0) {
        if (!s->world.networks || s->world.tower_per_site != 1 ||
            s->world.amplifier_per_site != 1 || s->world.controller_per_site != 1 ||
            !s->runtime.start_nodes || !s->runtime.wait_nodes_ready)
            return fail(err, "webapp sites require a network, one tower/amplifier/controller, "
                        "runtime start nodes and wait nodes_ready");
    } else if (s->world.tower_per_site || s->world.amplifier_per_site ||
               s->world.controller_per_site || s->runtime.start_nodes ||
               s->runtime.wait_nodes_ready) {
        return fail(err, "node runtime requires world sites");
    }
    if (!s->phase_count) return fail(err, "at least one phase is required");
    for (i = 0; i < s->phase_count; i++) {
        phase = &s->phases[i];
        if (!phase->name[0] || (!phase->event_count && !phase->check_count))
            return fail(err, "webapp phases require a name and at least one event/check");
        for (k = 0; k < i; k++)
            if (ulab_streq(phase->name, s->phases[k].name))
                return fail(err, "webapp phase names must be unique");
        checks += phase->check_count;
        for (j = 0; j < phase->event_count; j++) {
            event = &phase->events[j];
            if (scenario_is_web_event(event->type)) {
                if (browser_event(s, event, err)) return ULAB_ERR;
            } else if (event->type == EVT_START_UES || event->type == EVT_TRAFFIC) {
                if (event->timeout_seconds > 900 || event->timeout_seconds > s->webapp.scenario_timeout_seconds ||
                    ((event->web_fields & 1u) && !event->timeout_seconds) || !s->world.ues_per_site || s->world.sims_per_network ||
                    (event->ues.kind != SEL_ALL && !commerce_ue(s, &event->ues)) ||
                    (event->type == EVT_TRAFFIC && !event->amount_mb) || event->expect_result[0] || event->error_contains[0])
                    return fail(err, "UE runtime events require site-anchored UEs, positive traffic and no ignored failures");
            } else if (event->type == EVT_DISCONNECT_NODES ||
                       event->type == EVT_RECONNECT_NODES) {
                if (!node_selector(s, &event->nodes))
                    return fail(err, "node fault requires an existing node reference or all");
                if (event->expect_result[0] || event->error_contains[0])
                    return fail(err, "webapp fault injection failures cannot be ignored");
            } else {
                return fail(err, "event is outside the webapp foundation contract; "
                            "operator actions must use browser handlers");
            }
        }
        for (j = 0; j < phase->check_count; j++)
            if (browser_check(s, &phase->checks[j], err)) return ULAB_ERR;
    }
    for (i = 0; i < s->final_check_count; i++)
        if (browser_check(s, &s->final_checks[i], err)) return ULAB_ERR;
    if (!checks) return fail(err, "webapp scenarios require at least one browser check");
    return ULAB_OK;
}
