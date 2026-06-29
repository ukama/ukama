/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "scenario.h"
#include "log.h"
#include "util.h"

typedef enum {
    SEC_NONE = 0,
    SEC_WORLD,
    SEC_NODES_PER_SITE,
    SEC_PACKAGES,
    SEC_SETUP,
    SEC_SETUP_LIST,
    SEC_PROVIDER,
    SEC_RUNTIME,
    SEC_PROFILES,
    SEC_PROFILE_ONE,
    SEC_PROFILE_BUCKET,
    SEC_PHASES,
    SEC_PHASE_EVENTS,
    SEC_PHASE_CHECKS,
    SEC_FINAL_CHECKS
} parse_sec_t;

void scenario_init(scenario_t *s) {
    memset(s, 0, sizeof(*s));
    snprintf(s->suite, sizeof(s->suite), "default");
    snprintf(s->priority, sizeof(s->priority), "p2");
    snprintf(s->status, sizeof(s->status), "active");
    snprintf(s->provider.type, sizeof(s->provider.type), "virtual");
}

const char *scenario_event_name(event_type_t type) {
    switch (type) {
    case EVT_TRAFFIC: return "traffic";
    case EVT_TRAFFIC_BY_PROFILE: return "traffic_by_profile";
    case EVT_CREATE_UES: return "create_ues";
    case EVT_START_UES: return "start_ues";
    case EVT_WAIT_UES_ATTACHED: return "wait_ues_attached";
    case EVT_RESTART_NODES: return "restart_nodes";
    case EVT_WAIT_NODES_READY: return "wait_nodes_ready";
    case EVT_CHECK: return "check";
    default: return "unknown";
    }
}

const char *scenario_check_name(check_type_t type) {
    switch (type) {
    case CHECK_MODEL_COUNT: return "model_count";
    case CHECK_BFF_COUNT: return "bff_count";
    case CHECK_NODE_READY: return "node_ready";
    case CHECK_UE_ATTACHED: return "ue_attached";
    case CHECK_USAGE_PER_SIM: return "usage_per_sim";
    case CHECK_USAGE_SAMPLE: return "usage_sample";
    case CHECK_PACKAGE_ACTIVE: return "package_active";
    case CHECK_PACKAGE_REMAINING: return "package_remaining";
    case CHECK_NODE_STATE: return "node_state";
    case CHECK_DASHBOARD_LOADS: return "dashboard_loads";
    case CHECK_BALANCE_NON_NEGATIVE: return "balance_non_negative";
    default: return "unknown";
    }
}

int scenario_event_from_name(const char *name, event_type_t *out) {
    if (ulab_streq(name, "traffic")) *out = EVT_TRAFFIC;
    else if (ulab_streq(name, "traffic_by_profile")) {
        *out = EVT_TRAFFIC_BY_PROFILE;
    } else if (ulab_streq(name, "create_ues")) *out = EVT_CREATE_UES;
    else if (ulab_streq(name, "start_ues")) *out = EVT_START_UES;
    else if (ulab_streq(name, "wait_ues_attached")) {
        *out = EVT_WAIT_UES_ATTACHED;
    } else if (ulab_streq(name, "restart_nodes")) {
        *out = EVT_RESTART_NODES;
    } else if (ulab_streq(name, "wait_nodes_ready")) {
        *out = EVT_WAIT_NODES_READY;
    } else if (ulab_streq(name, "check")) *out = EVT_CHECK;
    else return ULAB_ERR;
    return ULAB_OK;
}

int scenario_check_from_name(const char *name, check_type_t *out) {
    if (ulab_streq(name, "count") || ulab_streq(name, "model_count")) {
        *out = CHECK_MODEL_COUNT;
    } else if (ulab_streq(name, "bff_count")) *out = CHECK_BFF_COUNT;
    else if (ulab_streq(name, "node_ready")) *out = CHECK_NODE_READY;
    else if (ulab_streq(name, "ue_attached")) *out = CHECK_UE_ATTACHED;
    else if (ulab_streq(name, "usage_per_sim")) {
        *out = CHECK_USAGE_PER_SIM;
    } else if (ulab_streq(name, "usage_sample")) {
        *out = CHECK_USAGE_SAMPLE;
    } else if (ulab_streq(name, "package_active")) {
        *out = CHECK_PACKAGE_ACTIVE;
    } else if (ulab_streq(name, "package_remaining")) {
        *out = CHECK_PACKAGE_REMAINING;
    } else if (ulab_streq(name, "node_state")) *out = CHECK_NODE_STATE;
    else if (ulab_streq(name, "dashboard_loads")) {
        *out = CHECK_DASHBOARD_LOADS;
    } else if (ulab_streq(name, "balance_non_negative")) {
        *out = CHECK_BALANCE_NON_NEGATIVE;
    } else return ULAB_ERR;
    return ULAB_OK;
}

static int indent_of(const char *s) {
    int n = 0;

    while (*s == ' ') {
        n++;
        s++;
    }
    return n;
}

static void strip_comment(char *s) {
    int quote = 0;

    while (*s) {
        if (*s == '"' || *s == '\'') {
            quote = !quote;
        }
        if (!quote && *s == '#') {
            *s = '\0';
            return;
        }
        s++;
    }
}

static int split_kv(char *s, char **key, char **val) {
    char *p;

    p = strchr(s, ':');
    if (p == NULL) {
        return ULAB_ERR;
    }
    *p = '\0';
    *key = ulab_trim(s);
    *val = ulab_trim(p + 1);
    return ULAB_OK;
}

static int parse_selector_value(selector_t *sel, const char *key,
                                const char *val) {
    memset(sel, 0, sizeof(*sel));
    if (ulab_streq(key, "ues") || ulab_streq(key, "nodes") ||
        ulab_streq(key, "sites") || ulab_streq(key, "networks")) {
        if (ulab_streq(val, "all")) {
            sel->kind = SEL_ALL;
            return ULAB_OK;
        }
    }
    return ULAB_ERR;
}

static int parse_inline_list(const char *v, const char *item) {
    char buf[ULAB_MAX_LINE];
    char *p;
    char *tok;

    if (!ulab_starts(v, "[") || !ulab_ends(v, "]")) {
        return 0;
    }
    if (ulab_copy(buf, sizeof(buf), v + 1) != ULAB_OK) {
        return 0;
    }
    buf[strlen(buf) - 1] = '\0';
    tok = strtok_r(buf, ",", &p);
    while (tok != NULL) {
        tok = ulab_trim(tok);
        if (ulab_streq(tok, item)) {
            return 1;
        }
        tok = strtok_r(NULL, ",", &p);
    }
    return 0;
}

static event_spec_t *new_event(phase_spec_t *p, const char *type,
                               ulab_error_t *err) {
    event_spec_t *e;

    if (p->event_count >= ULAB_MAX_EVENTS) {
        snprintf(err->msg, sizeof(err->msg), "too many events");
        return NULL;
    }
    e = &p->events[p->event_count++];
    memset(e, 0, sizeof(*e));
    if (scenario_event_from_name(type, &e->type) != ULAB_OK) {
        snprintf(err->msg, sizeof(err->msg), "unknown event type: %s", type);
        return NULL;
    }
    return e;
}

static check_spec_t *new_check(check_spec_t *arr, size_t *cnt,
                               const char *type, ulab_error_t *err) {
    check_spec_t *c;

    if (*cnt >= ULAB_MAX_CHECKS) {
        snprintf(err->msg, sizeof(err->msg), "too many checks");
        return NULL;
    }
    c = &arr[(*cnt)++];
    memset(c, 0, sizeof(*c));
    c->tolerance_percent = 2;
    if (scenario_check_from_name(type, &c->type) != ULAB_OK) {
        snprintf(err->msg, sizeof(err->msg), "unknown check type: %s", type);
        return NULL;
    }
    return c;
}

static int apply_check_field(check_spec_t *c, const char *key,
                             const char *val) {
    if (ulab_streq(key, "target")) return ulab_copy(c->target,
        sizeof(c->target), val);
    if (ulab_streq(key, "expected")) return ulab_copy(c->expected,
        sizeof(c->expected), val);
    if (ulab_streq(key, "package")) return ulab_copy(c->package_ref,
        sizeof(c->package_ref), val);
    if (ulab_streq(key, "expected_used_mb")) {
        return ulab_parse_u64(val, &c->expected_used_mb);
    }
    if (ulab_streq(key, "expected_remaining_mb")) {
        return ulab_parse_u64(val, &c->expected_remaining_mb);
    }
    if (ulab_streq(key, "tolerance_percent")) {
        return ulab_parse_u32(val, &c->tolerance_percent);
    }
    if (ulab_streq(key, "required")) {
        c->required = ulab_streq(val, "true") || ulab_streq(val, "1");
        return ULAB_OK;
    }
    if (parse_selector_value(&c->ues, key, val) == ULAB_OK &&
        ulab_streq(key, "ues")) return ULAB_OK;
    if (parse_selector_value(&c->nodes, key, val) == ULAB_OK &&
        ulab_streq(key, "nodes")) return ULAB_OK;
    if (parse_selector_value(&c->sites, key, val) == ULAB_OK &&
        ulab_streq(key, "sites")) return ULAB_OK;
    if (parse_selector_value(&c->networks, key, val) == ULAB_OK &&
        ulab_streq(key, "networks")) return ULAB_OK;
    if (ulab_streq(key, "sample_per_site")) {
        c->ues.kind = SEL_SAMPLE_PER_SITE;
        return ulab_parse_u32(val, &c->ues.count);
    }
    if (ulab_streq(key, "created_in_phase")) {
        c->ues.kind = SEL_CREATED_IN_PHASE;
        return ulab_copy(c->ues.value, sizeof(c->ues.value), val);
    }
    return ULAB_ERR;
}

static int apply_event_field(event_spec_t *e, const char *key,
                             const char *val) {
    if (ulab_streq(key, "name")) return ulab_copy(e->name,
        sizeof(e->name), val);
    if (ulab_streq(key, "amount_mb")) {
        return ulab_parse_u64(val, &e->amount_mb);
    }
    if (ulab_streq(key, "profile")) return ulab_copy(e->profile,
        sizeof(e->profile), val);
    if (ulab_streq(key, "count_per_site")) {
        return ulab_parse_u32(val, &e->count_per_site);
    }
    if (ulab_streq(key, "package")) return ulab_copy(e->package_ref,
        sizeof(e->package_ref), val);
    if (parse_selector_value(&e->ues, key, val) == ULAB_OK &&
        ulab_streq(key, "ues")) return ULAB_OK;
    if (parse_selector_value(&e->sites, key, val) == ULAB_OK &&
        ulab_streq(key, "sites")) return ULAB_OK;
    if (parse_selector_value(&e->nodes, key, val) == ULAB_OK &&
        ulab_streq(key, "nodes")) return ULAB_OK;
    if (ulab_streq(key, "created_in_phase")) {
        e->ues.kind = SEL_CREATED_IN_PHASE;
        return ulab_copy(e->ues.value, sizeof(e->ues.value), val);
    }
    if (ulab_streq(key, "affected_by_phase")) {
        e->nodes.kind = SEL_AFFECTED_BY_PHASE;
        return ulab_copy(e->nodes.value, sizeof(e->nodes.value), val);
    }
    if (ulab_streq(key, "type_selector")) {
        e->nodes.kind = SEL_NODE_TYPE;
        return ulab_copy(e->nodes.value, sizeof(e->nodes.value), val);
    }
    if (ulab_streq(key, "count_per_network")) {
        e->nodes.kind = SEL_NODE_TYPE_COUNT_PER_NETWORK;
        return ulab_parse_u32(val, &e->nodes.count);
    }
    return ULAB_ERR;
}

static int parse_item_value(char *line, char **key, char **val) {
    char *p = ulab_trim(line);

    if (!ulab_starts(p, "- ")) {
        return ULAB_ERR;
    }
    p += 2;
    return split_kv(p, key, val);
}

int scenario_load(const char *path, scenario_t *s, ulab_error_t *err) {
    FILE *fp;
    char line[ULAB_MAX_LINE];
    parse_sec_t sec = SEC_NONE;
    phase_spec_t *phase = NULL;
    event_spec_t *event = NULL;
    check_spec_t *check = NULL;
    package_spec_t *pkg = NULL;
    profile_spec_t *prof = NULL;
    profile_bucket_t *bucket = NULL;
    int lineno = 0;

    scenario_init(s);
    fp = fopen(path, "r");
    if (fp == NULL) {
        snprintf(err->msg, sizeof(err->msg), "unable to open %s", path);
        return ULAB_ERR;
    }

    while (fgets(line, sizeof(line), fp) != NULL) {
        char *key;
        char *val;
        char *p;
        int ind;

        lineno++;
        strip_comment(line);
        p = ulab_trim(line);
        if (*p == '\0') {
            continue;
        }
        ind = indent_of(line);
        if (!ulab_starts(p, "- ") &&
            split_kv(p, &key, &val) != ULAB_OK) {
            snprintf(err->msg, sizeof(err->msg), "line %d: bad syntax", lineno);
            fclose(fp);
            return ULAB_ERR;
        }

        if (ind == 0 && !ulab_starts(p, "- ")) {
            if (ulab_streq(key, "version")) {
                if (ulab_parse_u32(val, &s->version) != ULAB_OK) goto bad;
            } else if (ulab_streq(key, "name")) {
                if (ulab_copy(s->name, sizeof(s->name), val) != ULAB_OK) {
                    goto bad;
                }
            } else if (ulab_streq(key, "seed")) {
                if (ulab_parse_u32(val, &s->seed) != ULAB_OK) goto bad;
            } else if (ulab_streq(key, "suite")) {
                if (ulab_copy(s->suite, sizeof(s->suite), val)) goto bad;
            } else if (ulab_streq(key, "priority")) {
                if (ulab_copy(s->priority, sizeof(s->priority), val)) goto bad;
            } else if (ulab_streq(key, "tags")) {
                if (ulab_copy(s->tags, sizeof(s->tags), val)) goto bad;
            } else if (ulab_streq(key, "status")) {
                if (ulab_copy(s->status, sizeof(s->status), val)) goto bad;
            } else if (ulab_streq(key, "world")) sec = SEC_WORLD;
            else if (ulab_streq(key, "packages")) sec = SEC_PACKAGES;
            else if (ulab_streq(key, "setup")) sec = SEC_SETUP;
            else if (ulab_streq(key, "provider")) sec = SEC_PROVIDER;
            else if (ulab_streq(key, "runtime")) sec = SEC_RUNTIME;
            else if (ulab_streq(key, "profiles")) sec = SEC_PROFILES;
            else if (ulab_streq(key, "phases")) sec = SEC_PHASES;
            else if (ulab_streq(key, "final_checks")) sec = SEC_FINAL_CHECKS;
            else goto unknown;
            continue;
        }

        if (sec == SEC_WORLD) {
            if (ind == 2 && ulab_streq(key, "networks")) {
                if (ulab_parse_u32(val, &s->world.networks) != ULAB_OK) {
                    goto bad;
                }
            } else if (ind == 2 && ulab_streq(key, "sites_per_network")) {
                if (ulab_parse_u32(val, &s->world.sites_per_network)) goto bad;
            } else if (ind == 2 && ulab_streq(key, "ues_per_site")) {
                if (ulab_parse_u32(val, &s->world.ues_per_site)) goto bad;
            } else if (ind == 2 && ulab_streq(key, "nodes_per_site")) {
                sec = SEC_NODES_PER_SITE;
            } else goto unknown;
            continue;
        }
        if (sec == SEC_NODES_PER_SITE) {
            if (ind == 4 && ulab_streq(key, "tower")) {
                if (ulab_parse_u32(val, &s->world.tower_per_site)) goto bad;
            } else if (ind == 4 && ulab_streq(key, "amplifier")) {
                if (ulab_parse_u32(val, &s->world.amplifier_per_site)) goto bad;
            } else if (ind == 4 && ulab_streq(key, "controller")) {
                if (ulab_parse_u32(val, &s->world.controller_per_site)) goto bad;
            } else if (ind == 2 && ulab_streq(key, "ues_per_site")) {
                if (ulab_parse_u32(val, &s->world.ues_per_site)) goto bad;
                sec = SEC_WORLD;
            } else if (ind == 2 && ulab_streq(key, "networks")) {
                if (ulab_parse_u32(val, &s->world.networks)) goto bad;
                sec = SEC_WORLD;
            } else if (ind == 2 && ulab_streq(key, "sites_per_network")) {
                if (ulab_parse_u32(val, &s->world.sites_per_network)) goto bad;
                sec = SEC_WORLD;
            } else goto unknown;
            continue;
        }
        if (sec == SEC_PACKAGES) {
            if (ind == 2 && ulab_starts(p, "- ")) {
                if (s->package_count >= ULAB_MAX_PACKAGES) goto many;
                pkg = &s->packages[s->package_count++];
                memset(pkg, 0, sizeof(*pkg));
                if (parse_item_value(p, &key, &val) ||
                    !ulab_streq(key, "ref")) goto bad;
                if (ulab_copy(pkg->ref, sizeof(pkg->ref), val)) goto bad;
            } else if (ind == 4 && pkg != NULL) {
                if (ulab_streq(key, "name")) {
                    if (ulab_copy(pkg->name, sizeof(pkg->name), val)) goto bad;
                } else if (ulab_streq(key, "data_mb")) {
                    if (ulab_parse_u64(val, &pkg->data_mb)) goto bad;
                } else if (ulab_streq(key, "duration_days")) {
                    if (ulab_parse_u32(val, &pkg->duration_days)) goto bad;
                } else if (ulab_streq(key, "amount")) {
                    if (ulab_parse_double(val, &pkg->amount)) goto bad;
                } else if (ulab_streq(key, "assign_percent")) {
                    if (ulab_parse_u32(val, &pkg->assign_percent)) goto bad;
                } else goto unknown;
            } else goto unknown;
            continue;
        }
        if (sec == SEC_SETUP) {
            if (ind == 2 && ulab_streq(key, "create_via_bff")) {
                sec = SEC_SETUP_LIST;
            } else goto unknown;
            continue;
        }
        if (sec == SEC_SETUP_LIST) {
            if (ind == 4 && ulab_starts(p, "- ")) {
                char *item = ulab_trim(p + 2);
                if (ulab_streq(item, "networks")) s->setup.create_networks = 1;
                else if (ulab_streq(item, "sites")) s->setup.create_sites = 1;
                else if (ulab_streq(item, "nodes")) s->setup.create_nodes = 1;
                else if (ulab_streq(item, "node_site_links")) {
                    s->setup.create_node_site_links = 1;
                } else if (ulab_streq(item, "packages")) {
                    s->setup.create_packages = 1;
                } else if (ulab_streq(item, "subscribers")) {
                    s->setup.create_subscribers = 1;
                } else if (ulab_streq(item, "sims")) s->setup.create_sims = 1;
                else goto unknown;
            } else goto unknown;
            continue;
        }

        if (sec == SEC_PROVIDER) {
            if (ind == 2 && ulab_streq(key, "type")) {
                if (ulab_copy(s->provider.type,
                    sizeof(s->provider.type), val)) goto bad;
            } else goto unknown;
            continue;
        }

        if (sec == SEC_RUNTIME) {
            if (ind == 2 && ulab_streq(key, "start")) {
                s->runtime.start_nodes = parse_inline_list(val, "nodes");
                s->runtime.start_ues = parse_inline_list(val, "ues");
            } else if (ind == 2 && ulab_streq(key, "wait")) {
                s->runtime.wait_nodes_ready = parse_inline_list(val,
                    "nodes_ready");
                s->runtime.wait_ues_attached = parse_inline_list(val,
                    "ues_attached");
            } else goto unknown;
            continue;
        }
        if (sec == SEC_PROFILES || sec == SEC_PROFILE_ONE ||
            sec == SEC_PROFILE_BUCKET) {
            if (ind == 2) {
                if (s->profile_count >= ULAB_MAX_BUCKETS) goto many;
                prof = &s->profiles[s->profile_count++];
                memset(prof, 0, sizeof(*prof));
                if (ulab_copy(prof->name, sizeof(prof->name), key)) goto bad;
                sec = SEC_PROFILE_ONE;
            } else if (ind == 4 && prof != NULL) {
                if (prof->bucket_count >= ULAB_MAX_BUCKETS) goto many;
                bucket = &prof->buckets[prof->bucket_count++];
                memset(bucket, 0, sizeof(*bucket));
                if (ulab_copy(bucket->name, sizeof(bucket->name), key)) {
                    goto bad;
                }
                sec = SEC_PROFILE_BUCKET;
            } else if (ind == 6 && bucket != NULL) {
                if (ulab_streq(key, "percent")) {
                    if (ulab_parse_u32(val, &bucket->percent)) goto bad;
                } else if (ulab_streq(key, "amount_mb")) {
                    if (ulab_parse_u64(val, &bucket->amount_mb)) goto bad;
                } else goto unknown;
            } else goto unknown;
            continue;
        }
        if (sec == SEC_PHASES) {
            if (ind == 2 && ulab_starts(p, "- ")) {
                if (s->phase_count >= ULAB_MAX_PHASES) goto many;
                phase = &s->phases[s->phase_count++];
                memset(phase, 0, sizeof(*phase));
                if (parse_item_value(p, &key, &val) ||
                    !ulab_streq(key, "name")) goto bad;
                if (ulab_copy(phase->name, sizeof(phase->name), val)) goto bad;
            } else if (ind == 4 && ulab_streq(key, "events")) {
                sec = SEC_PHASE_EVENTS;
            } else if (ind == 4 && ulab_streq(key, "checks")) {
                sec = SEC_PHASE_CHECKS;
            } else goto unknown;
            continue;
        }
        if (sec == SEC_PHASE_EVENTS) {
            if (ind == 6 && ulab_starts(p, "- ")) {
                if (parse_item_value(p, &key, &val) ||
                    !ulab_streq(key, "type")) goto bad;
                event = new_event(phase, val, err);
                if (event == NULL) goto fail;
            } else if (ind == 8 && event != NULL) {
                if (apply_event_field(event, key, val) != ULAB_OK) goto unknown;
            } else if (ind == 4 && ulab_streq(key, "checks")) {
                sec = SEC_PHASE_CHECKS;
            } else goto unknown;
            continue;
        }
        if (sec == SEC_PHASE_CHECKS || sec == SEC_FINAL_CHECKS) {
            check_spec_t *arr = sec == SEC_FINAL_CHECKS ?
                s->final_checks : phase->checks;
            size_t *cnt = sec == SEC_FINAL_CHECKS ?
                &s->final_check_count : &phase->check_count;

            if (sec == SEC_PHASE_CHECKS && ind == 2 &&
                ulab_starts(p, "- ")) {
                if (s->phase_count >= ULAB_MAX_PHASES) goto many;
                phase = &s->phases[s->phase_count++];
                memset(phase, 0, sizeof(*phase));
                if (parse_item_value(p, &key, &val) ||
                    !ulab_streq(key, "name")) goto bad;
                if (ulab_copy(phase->name, sizeof(phase->name), val)) goto bad;
                sec = SEC_PHASES;
                continue;
            }
            if ((ind == 2 || ind == 6) && ulab_starts(p, "- ")) {
                if (parse_item_value(p, &key, &val) ||
                    !ulab_streq(key, "type")) goto bad;
                check = new_check(arr, cnt, val, err);
                if (check == NULL) goto fail;
            } else if ((ind == 4 || ind == 8) && check != NULL) {
                if (apply_check_field(check, key, val) != ULAB_OK) goto unknown;
            } else if (ind == 2 && sec != SEC_FINAL_CHECKS) {
                sec = SEC_PHASES;
            } else goto unknown;
            continue;
        }
    }

    fclose(fp);
    return ULAB_OK;

bad:
    snprintf(err->msg, sizeof(err->msg), "line %d: invalid value", lineno);
    goto fail;
unknown:
    snprintf(err->msg, sizeof(err->msg), "line %d: unknown field", lineno);
    goto fail;
many:
    snprintf(err->msg, sizeof(err->msg), "line %d: too many entries", lineno);
fail:
    fclose(fp);
    return ULAB_ERR;
}

void scenario_list_events(void) {
    int i;

    for (i = EVT_TRAFFIC; i <= EVT_CHECK; i++) {
        printf("%s\n", scenario_event_name((event_type_t)i));
    }
}

void scenario_list_checks(void) {
    int i;

    for (i = CHECK_MODEL_COUNT; i <= CHECK_BALANCE_NON_NEGATIVE; i++) {
        printf("%s\n", scenario_check_name((check_type_t)i));
    }
}
