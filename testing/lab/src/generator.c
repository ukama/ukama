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

#include "generator.h"
#include "scenario.h"
#include "util.h"

#define GEN_MAX_FAMILIES  32
#define GEN_MAX_CASES     512
#define GEN_MAX_PROFILES  16
#define GEN_MAX_LIST      24
#define GEN_MAX_TAGS      384

/*
 * Phase-6 generator rule:
 *   - do NOT create blind Cartesian products.
 *   - read explicit family cases from models/families YAML files.
 *   - use simple/medium/large topology+scale profiles to materialize world.
 *   - generate one meaningful scenario per named case.
 */

typedef struct {
    char model[ULAB_MAX_REF];
    char out_dir[ULAB_MAX_PATH];
    char models_dir[ULAB_MAX_PATH];
} gen_opts_t;

typedef struct {
    char name[ULAB_MAX_REF];
    uint32_t networks;
    uint32_t sites_per_network;
    uint32_t ues_per_site;
    uint32_t packages;
    uint64_t traffic_mb_per_ue;
} topo_profile_t;

typedef struct {
    char name[ULAB_MAX_REF];
    char priority[ULAB_MAX_REF];
    char tags[GEN_MAX_TAGS];
    char setup[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t setup_count;
    char events[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t event_count;
    char checks[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t check_count;
    char topology[ULAB_MAX_REF];
} gen_case_t;

typedef struct {
    char name[ULAB_MAX_REF];
    char priority[ULAB_MAX_REF];
    gen_case_t cases[GEN_MAX_CASES];
    size_t case_count;
} gen_family_t;

typedef struct {
    topo_profile_t profiles[GEN_MAX_PROFILES];
    size_t profile_count;
    gen_family_t families[GEN_MAX_FAMILIES];
    size_t family_count;
} gen_catalog_t;

static const char *default_families[] = {
    "smoke", "backend", "usage", "sim_pool", "package", "subscriber",
    "lifecycle", "node_ops", "site_ops", "software_update", "failure",
    "scale"
};

static void usage(void) {
    printf("usage:\n");
    printf("  ukama-lab generate --model <all|family> [options]\n");
    printf("options:\n");
    printf("  --out <dir>        output directory; default: scenarios/generated\n");
    printf("  --models <dir>     model directory; default: models\n");
}

static int append_str(char *dst, size_t n, const char *src) {
    size_t used;
    size_t i;

    if (dst == NULL || n == 0 || src == NULL) {
        return ULAB_ERR;
    }
    used = strlen(dst);
    for (i = 0; src[i] != '\0'; i++) {
        if (used + 1 >= n) {
            dst[n - 1] = '\0';
            return ULAB_ERR;
        }
        dst[used++] = src[i];
    }
    dst[used] = '\0';
    return ULAB_OK;
}

static void clean_token(char *s) {
    size_t i;

    for (i = 0; s != NULL && s[i] != '\0'; i++) {
        if (s[i] == ' ' || s[i] == '/' || s[i] == ':' ||
            s[i] == ',' || s[i] == '[' || s[i] == ']') {
            s[i] = '_';
        }
    }
}

static int path_join(char *out, size_t n, const char *a, const char *b) {
    out[0] = '\0';
    if (append_str(out, n, a)) return ULAB_ERR;
    if (append_str(out, n, "/")) return ULAB_ERR;
    if (append_str(out, n, b)) return ULAB_ERR;
    return ULAB_OK;
}

static int path_join3(char *out, size_t n, const char *a, const char *b,
                      const char *c) {
    out[0] = '\0';
    if (append_str(out, n, a)) return ULAB_ERR;
    if (append_str(out, n, "/")) return ULAB_ERR;
    if (append_str(out, n, b)) return ULAB_ERR;
    if (append_str(out, n, "/")) return ULAB_ERR;
    if (append_str(out, n, c)) return ULAB_ERR;
    return ULAB_OK;
}

static int parse_opts(int argc, char **argv, gen_opts_t *opts) {
    int i;

    memset(opts, 0, sizeof(*opts));
    ulab_copy(opts->model, sizeof(opts->model), "all");
    ulab_copy(opts->out_dir, sizeof(opts->out_dir), "scenarios/generated");
    ulab_copy(opts->models_dir, sizeof(opts->models_dir), "models");

    for (i = 0; i < argc; i++) {
        if (ulab_streq(argv[i], "--model") && i + 1 < argc) {
            if (ulab_copy(opts->model, sizeof(opts->model), argv[++i])) {
                return ULAB_EUSAGE;
            }
        } else if (ulab_streq(argv[i], "--out") && i + 1 < argc) {
            if (ulab_copy(opts->out_dir, sizeof(opts->out_dir), argv[++i])) {
                return ULAB_EUSAGE;
            }
        } else if (ulab_streq(argv[i], "--models") && i + 1 < argc) {
            if (ulab_copy(opts->models_dir, sizeof(opts->models_dir), argv[++i])) {
                return ULAB_EUSAGE;
            }
        } else if ((ulab_streq(argv[i], "--mode") ||
                    ulab_streq(argv[i], "--templates")) && i + 1 < argc) {
            /* accepted for backwards-compatible scripts; ignored */
            i++;
        } else if (ulab_streq(argv[i], "--help")) {
            usage();
            return ULAB_EUSAGE;
        } else {
            fprintf(stderr, "unknown generate option: %s\n", argv[i]);
            return ULAB_EUSAGE;
        }
    }

    return ULAB_OK;
}

static int split_kv(char *line, char **key, char **val) {
    char *c;

    c = strchr(line, ':');
    if (c == NULL) {
        return ULAB_ERR;
    }
    *c = '\0';
    *key = ulab_trim(line);
    *val = ulab_trim(c + 1);
    return ULAB_OK;
}

static int add_word(char arr[GEN_MAX_LIST][ULAB_MAX_REF], size_t *count,
                    const char *word) {
    if (word == NULL || word[0] == '\0') {
        return ULAB_OK;
    }
    if (*count >= GEN_MAX_LIST) {
        return ULAB_ERR;
    }
    if (ulab_copy(arr[*count], ULAB_MAX_REF, word)) {
        return ULAB_ERR;
    }
    clean_token(arr[*count]);
    (*count)++;
    return ULAB_OK;
}

static int parse_inline_list(char arr[GEN_MAX_LIST][ULAB_MAX_REF],
                             size_t *count, char *val) {
    char *p;
    char *tok;

    *count = 0;
    p = ulab_trim(val);
    if (*p == '[') {
        p++;
    }
    tok = strtok(p, ",]");
    while (tok != NULL) {
        char *w = ulab_trim(tok);
        if (w[0] != '\0' && add_word(arr, count, w)) {
            return ULAB_ERR;
        }
        tok = strtok(NULL, ",]");
    }
    return ULAB_OK;
}

static const topo_profile_t *find_profile(const gen_catalog_t *cat,
                                          const char *name) {
    size_t i;

    for (i = 0; i < cat->profile_count; i++) {
        if (ulab_streq(cat->profiles[i].name, name)) {
            return &cat->profiles[i];
        }
    }
    return NULL;
}

static int load_topology_profiles(const gen_opts_t *opts,
                                  gen_catalog_t *cat) {
    char path[ULAB_MAX_PATH];
    char line[ULAB_MAX_LINE];
    FILE *fp;
    topo_profile_t *cur;

    if (path_join3(path, sizeof(path), opts->models_dir, "profiles",
                   "topology.yaml")) {
        return ULAB_ERR;
    }
    fp = fopen(path, "r");
    if (fp == NULL) {
        fprintf(stderr, "missing topology profile: %s\n", path);
        return ULAB_ERR;
    }

    cur = NULL;
    while (fgets(line, sizeof(line), fp) != NULL) {
        char *p;
        char *key;
        char *val;

        p = ulab_trim(line);
        if (*p == '\0' || *p == '#') {
            continue;
        }
        if (ulab_starts(p, "- name:")) {
            if (cat->profile_count >= GEN_MAX_PROFILES) {
                fclose(fp);
                return ULAB_ERR;
            }
            cur = &cat->profiles[cat->profile_count++];
            memset(cur, 0, sizeof(*cur));
            val = ulab_trim(strchr(p, ':') + 1);
            if (ulab_copy(cur->name, sizeof(cur->name), val)) {
                fclose(fp);
                return ULAB_ERR;
            }
            clean_token(cur->name);
            continue;
        }
        if (cur == NULL || split_kv(p, &key, &val)) {
            continue;
        }
        if (ulab_streq(key, "networks")) {
            if (ulab_parse_u32(val, &cur->networks)) goto bad;
        } else if (ulab_streq(key, "sites_per_network")) {
            if (ulab_parse_u32(val, &cur->sites_per_network)) goto bad;
        }
    }

    fclose(fp);
    return cat->profile_count > 0 ? ULAB_OK : ULAB_ERR;

bad:
    fclose(fp);
    return ULAB_ERR;
}

static int load_scale_profiles(const gen_opts_t *opts,
                               gen_catalog_t *cat) {
    char path[ULAB_MAX_PATH];
    char line[ULAB_MAX_LINE];
    FILE *fp;
    topo_profile_t *cur;

    if (path_join3(path, sizeof(path), opts->models_dir, "profiles",
                   "scale.yaml")) {
        return ULAB_ERR;
    }
    fp = fopen(path, "r");
    if (fp == NULL) {
        fprintf(stderr, "missing scale profile: %s\n", path);
        return ULAB_ERR;
    }

    cur = NULL;
    while (fgets(line, sizeof(line), fp) != NULL) {
        char *p;
        char *key;
        char *val;
        size_t i;

        p = ulab_trim(line);
        if (*p == '\0' || *p == '#') {
            continue;
        }
        if (ulab_starts(p, "- name:")) {
            val = ulab_trim(strchr(p, ':') + 1);
            cur = NULL;
            for (i = 0; i < cat->profile_count; i++) {
                if (ulab_streq(cat->profiles[i].name, val)) {
                    cur = &cat->profiles[i];
                    break;
                }
            }
            continue;
        }
        if (cur == NULL || split_kv(p, &key, &val)) {
            continue;
        }
        if (ulab_streq(key, "ues_per_site")) {
            if (ulab_parse_u32(val, &cur->ues_per_site)) goto bad;
        } else if (ulab_streq(key, "packages")) {
            if (ulab_parse_u32(val, &cur->packages)) goto bad;
        } else if (ulab_streq(key, "traffic_mb_per_ue")) {
            if (ulab_parse_u64(val, &cur->traffic_mb_per_ue)) goto bad;
        }
    }

    fclose(fp);
    return ULAB_OK;

bad:
    fclose(fp);
    return ULAB_ERR;
}

static gen_family_t *add_family(gen_catalog_t *cat, const char *name) {
    gen_family_t *f;

    if (cat->family_count >= GEN_MAX_FAMILIES) {
        return NULL;
    }
    f = &cat->families[cat->family_count++];
    memset(f, 0, sizeof(*f));
    if (ulab_copy(f->name, sizeof(f->name), name)) {
        return NULL;
    }
    clean_token(f->name);
    ulab_copy(f->priority, sizeof(f->priority), "p2");
    return f;
}

static gen_case_t *add_case(gen_family_t *family, const char *name) {
    gen_case_t *c;

    if (family->case_count >= GEN_MAX_CASES) {
        return NULL;
    }
    c = &family->cases[family->case_count++];
    memset(c, 0, sizeof(*c));
    if (ulab_copy(c->name, sizeof(c->name), name)) {
        return NULL;
    }
    clean_token(c->name);
    ulab_copy(c->topology, sizeof(c->topology), "simple");
    ulab_copy(c->priority, sizeof(c->priority),
              family->priority[0] ? family->priority : "p2");
    return c;
}

static int load_family_file(const gen_opts_t *opts, const char *name,
                            gen_catalog_t *cat) {
    char path[ULAB_MAX_PATH];
    char line[ULAB_MAX_LINE];
    FILE *fp;
    gen_family_t *family;
    gen_case_t *cur;
    int in_cases;

    if (path_join3(path, sizeof(path), opts->models_dir, "families", name)) {
        return ULAB_ERR;
    }
    if (append_str(path, sizeof(path), ".yaml")) {
        return ULAB_ERR;
    }
    fp = fopen(path, "r");
    if (fp == NULL) {
        fprintf(stderr, "missing family model: %s\n", path);
        return ULAB_ERR;
    }

    family = NULL;
    cur = NULL;
    in_cases = 0;

    while (fgets(line, sizeof(line), fp) != NULL) {
        char *p;
        char *key;
        char *val;

        p = ulab_trim(line);
        if (*p == '\0' || *p == '#') {
            continue;
        }
        if (ulab_streq(p, "cases:")) {
            in_cases = 1;
            continue;
        }
        if (split_kv(p, &key, &val)) {
            continue;
        }

        if (!in_cases) {
            if (ulab_streq(key, "family")) {
                family = add_family(cat, val);
                if (family == NULL) goto bad;
            } else if (family != NULL && ulab_streq(key, "priority")) {
                if (ulab_copy(family->priority, sizeof(family->priority), val)) {
                    goto bad;
                }
            }
            continue;
        }

        if (family == NULL) {
            goto bad;
        }
        if (ulab_streq(key, "- name")) {
            cur = add_case(family, val);
            if (cur == NULL) goto bad;
            continue;
        }
        if (cur == NULL) {
            continue;
        }
        if (ulab_streq(key, "topology")) {
            if (ulab_copy(cur->topology, sizeof(cur->topology), val)) goto bad;
            clean_token(cur->topology);
        } else if (ulab_streq(key, "priority")) {
            if (ulab_copy(cur->priority, sizeof(cur->priority), val)) goto bad;
        } else if (ulab_streq(key, "tags")) {
            if (ulab_copy(cur->tags, sizeof(cur->tags), val)) goto bad;
        } else if (ulab_streq(key, "setup")) {
            if (parse_inline_list(cur->setup, &cur->setup_count, val)) goto bad;
        } else if (ulab_streq(key, "events")) {
            if (parse_inline_list(cur->events, &cur->event_count, val)) goto bad;
        } else if (ulab_streq(key, "checks")) {
            if (parse_inline_list(cur->checks, &cur->check_count, val)) goto bad;
        }
    }

    fclose(fp);
    if (family == NULL || family->case_count == 0) {
        fprintf(stderr, "invalid or empty family model: %s\n", path);
        return ULAB_ERR;
    }
    return ULAB_OK;

bad:
    fclose(fp);
    fprintf(stderr, "invalid family model: %s\n", path);
    return ULAB_ERR;
}

static int load_catalog(const gen_opts_t *opts, gen_catalog_t *cat) {
    size_t i;

    memset(cat, 0, sizeof(*cat));
    if (load_topology_profiles(opts, cat)) {
        return ULAB_ERR;
    }
    if (load_scale_profiles(opts, cat)) {
        return ULAB_ERR;
    }

    if (!ulab_streq(opts->model, "all")) {
        return load_family_file(opts, opts->model, cat);
    }

    for (i = 0; i < sizeof(default_families) / sizeof(default_families[0]); i++) {
        if (load_family_file(opts, default_families[i], cat)) {
            return ULAB_ERR;
        }
    }
    return ULAB_OK;
}

static int list_has(const char arr[GEN_MAX_LIST][ULAB_MAX_REF], size_t count,
                    const char *word) {
    size_t i;

    for (i = 0; i < count; i++) {
        if (ulab_streq(arr[i], word)) {
            return 1;
        }
    }
    return 0;
}

static int case_has_tag(const gen_case_t *c, const char *word) {
    return c->tags[0] != '\0' && strstr(c->tags, word) != NULL;
}

static int event_supported(const char *event) {
    return ulab_streq(event, "check") ||
           ulab_streq(event, "traffic") ||
           ulab_streq(event, "traffic_1gb") ||
           ulab_streq(event, "traffic_2gb_per_ue") ||
           ulab_streq(event, "traffic_5gb_per_ue") ||
           ulab_streq(event, "traffic_by_profile") ||
           ulab_streq(event, "restart_nodes") ||
           ulab_streq(event, "restart_tower") ||
           ulab_streq(event, "restart_site") ||
           ulab_streq(event, "toggle_service_off") ||
           ulab_streq(event, "toggle_service_on") ||
           ulab_streq(event, "toggle_radio_off") ||
           ulab_streq(event, "toggle_radio_on") ||
           ulab_streq(event, "mark_node_offline") ||
           ulab_streq(event, "restore_node") ||
           ulab_streq(event, "software_update_tower") ||
           ulab_streq(event, "software_update_amplifier") ||
           ulab_streq(event, "software_update_controller") ||
           ulab_streq(event, "software_update_same_version") ||
           ulab_streq(event, "add_package_to_sim") ||
           ulab_streq(event, "topup") ||
           ulab_streq(event, "remove_package_from_sim") ||
           ulab_streq(event, "set_sim_inactive") ||
           ulab_streq(event, "set_sim_active");
}

static int check_supported(const char *check) {
    return ulab_streq(check, "backend_count") ||
           ulab_streq(check, "backend_count_networks") ||
           ulab_streq(check, "backend_count_sites") ||
           ulab_streq(check, "backend_count_nodes") ||
           ulab_streq(check, "backend_count_sims") ||
           ulab_streq(check, "backend_count_subscribers") ||
           ulab_streq(check, "backend_count_packages") ||
           ulab_streq(check, "backend_package_count") ||
           ulab_streq(check, "ue_attached") ||
           ulab_streq(check, "list_contains") ||
           ulab_streq(check, "usage_per_sim") ||
           ulab_streq(check, "balance_non_negative") ||
           ulab_streq(check, "traffic_allowed") ||
           ulab_streq(check, "traffic_blocked") ||
           ulab_streq(check, "node_ready") ||
           ulab_streq(check, "nodes_ready") ||
           ulab_streq(check, "node_health_ok") ||
           ulab_streq(check, "node_version_equals") ||
           ulab_streq(check, "history_preserved") ||
           ulab_streq(check, "relationship_exists") ||
           ulab_streq(check, "relationship_ended") ||
           ulab_streq(check, "audit_event_exists") ||
           ulab_streq(check, "status_equals") ||
           ulab_streq(check, "sim_assigned") ||
           ulab_streq(check, "sim_not_in_unassigned_pool") ||
           ulab_streq(check, "subscriber_has_sim") ||
           ulab_streq(check, "package_active");
}

static int setup_supported(const char *setup) {
    return ulab_streq(setup, "network") ||
           ulab_streq(setup, "site") ||
           ulab_streq(setup, "sites") ||
           ulab_streq(setup, "nodes") ||
           ulab_streq(setup, "subscriber") ||
           ulab_streq(setup, "subscribers") ||
           ulab_streq(setup, "sim") ||
           ulab_streq(setup, "sims") ||
           ulab_streq(setup, "package") ||
           ulab_streq(setup, "packages") ||
           ulab_streq(setup, "active_package") ||
           ulab_streq(setup, "active_subscriber") ||
           ulab_streq(setup, "sim_pool_one") ||
           ulab_streq(setup, "sim_csv_one") ||
           ulab_streq(setup, "known_iccid") ||
           ulab_streq(setup, "service_off") ||
           ulab_streq(setup, "radio_off") ||
           ulab_streq(setup, "node_offline") ||
           ulab_streq(setup, "restart_required");
}

static int case_is_wip(const gen_family_t *family, const gen_case_t *c) {
    size_t i;

    (void)family;
    for (i = 0; i < c->event_count; i++) {
        if (!event_supported(c->events[i])) {
            return 1;
        }
    }
    for (i = 0; i < c->check_count; i++) {
        if (!check_supported(c->checks[i])) {
            return 1;
        }
    }
    for (i = 0; i < c->setup_count; i++) {
        if (!setup_supported(c->setup[i])) {
            return 1;
        }
    }
    if (case_has_tag(c, "negative") || case_has_tag(c, "partial") ||
        case_has_tag(c, "retry") || case_has_tag(c, "duplicate") ||
        case_has_tag(c, "wrong_network") || case_has_tag(c, "rollback")) {
        return 1;
    }
    return 0;
}

static uint64_t event_amount_mb(const char *event, const topo_profile_t *p) {
    if (ulab_streq(event, "traffic_1gb")) return 1024;
    if (ulab_streq(event, "traffic_2gb_per_ue")) return 2048;
    if (ulab_streq(event, "traffic_5gb_per_ue")) return 5120;
    return p->traffic_mb_per_ue ? p->traffic_mb_per_ue : 1;
}

static int case_needs_ues(const gen_family_t *family, const gen_case_t *c) {
    size_t i;

    if (ulab_streq(family->name, "usage") ||
        ulab_streq(family->name, "node_ops") ||
        ulab_streq(family->name, "site_ops") ||
        ulab_streq(family->name, "lifecycle") ||
        ulab_streq(family->name, "scale") ||
        case_has_tag(c, "runtime")) {
        return 1;
    }
    for (i = 0; i < c->event_count; i++) {
        if (strstr(c->events[i], "traffic") != NULL ||
            strstr(c->events[i], "restart") != NULL ||
            strstr(c->events[i], "service") != NULL ||
            strstr(c->events[i], "radio") != NULL ||
            strstr(c->events[i], "offline") != NULL ||
            strstr(c->events[i], "restore") != NULL) {
            return 1;
        }
    }
    return 0;
}

static int case_has_event(const gen_case_t *c, const char *event) {
    return list_has(c->events, c->event_count, event);
}

static int case_has_setup(const gen_case_t *c, const char *setup) {
    return list_has(c->setup, c->setup_count, setup);
}

static uint32_t scenario_package_count(const gen_case_t *c,
                                       const topo_profile_t *p) {
    uint32_t count;

    count = p->packages ? p->packages : 1;
    if ((case_has_event(c, "add_package_to_sim") ||
         case_has_event(c, "topup")) && count < 2) {
        count = 2;
    }
    if (count > ULAB_MAX_PACKAGES) {
        count = ULAB_MAX_PACKAGES;
    }
    return count;
}

static void write_packages(FILE *f, uint32_t count, uint64_t data_mb) {
    uint32_t i;

    if (count == 0) {
        count = 1;
    }
    if (data_mb == 0) {
        data_mb = 1024;
    }

    fprintf(f, "packages:\n");
    for (i = 0; i < count; i++) {
        fprintf(f,
                "  - ref: pkg_%03u\n"
                "    name: Data Package %u\n"
                "    data_mb: %llu\n"
                "    duration_days: 1\n"
                "    amount: %.2f\n"
                "    assign_percent: %u\n",
                i + 1, i + 1, (unsigned long long)data_mb,
                (double)(i + 1), i == 0 ? 100u : 0u);
    }
    fprintf(f, "\n");
}

static void write_profile_section(FILE *f, const gen_case_t *c) {
    if (!case_has_event(c, "traffic_by_profile")) {
        return;
    }
    fprintf(f,
            "profiles:\n"
            "  mixed_usage:\n"
            "    light:\n"
            "      percent: 70\n"
            "      amount_mb: 512\n"
            "    heavy:\n"
            "      percent: 30\n"
            "      amount_mb: 2048\n\n");
}

static void write_setup_preamble(FILE *f, const gen_case_t *c) {
    if (case_has_setup(c, "service_off")) {
        fprintf(f,
                "      - type: toggle_service\n"
                "        state: off\n");
    }
    if (case_has_setup(c, "radio_off")) {
        fprintf(f,
                "      - type: toggle_radio\n"
                "        state: off\n");
    }
    if (case_has_setup(c, "node_offline")) {
        fprintf(f, "      - type: mark_node_offline\n");
    }
}

static const char *software_node_type(const char *event) {
    if (ulab_streq(event, "software_update_amplifier")) {
        return "amplifier";
    }
    if (ulab_streq(event, "software_update_controller")) {
        return "controller";
    }
    return "tower";
}

static int software_update_success_case(const gen_family_t *family,
                                        const gen_case_t *c) {
    const char *event;

    if (family == NULL || c == NULL || c->event_count == 0 ||
        !ulab_streq(family->name, "software_update")) {
        return 0;
    }

    event = c->events[0];
    return ulab_streq(event, "software_update_tower") ||
           ulab_streq(event, "software_update_amplifier") ||
           ulab_streq(event, "software_update_controller");
}

static void write_software_description(FILE *f, const gen_case_t *c) {
    const char *node_type;

    node_type = software_node_type(c->events[0]);
    fprintf(f,
            "description: Verify the example app is healthy at the "
            "expected current version on the %s node, update it from "
            "the exact tar.gz version published in Hub, and confirm "
            "the running target version and app health after the "
            "update.\n",
            node_type);
}

static void write_software_preflight(FILE *f, const gen_case_t *c) {
    const char *node_type;

    node_type = software_node_type(c->events[0]);
    fprintf(f,
            "  - name: verify_app_before_update\n"
            "    checks:\n"
            "      - type: node_health_ok\n"
            "        type_selector: %s\n"
            "        count_per_network: 1\n"
            "        app: example\n"
            "      - type: node_version_equals\n"
            "        type_selector: %s\n"
            "        count_per_network: 1\n"
            "        app: example\n"
            "        version: ${ULAB_SOFTWARE_CURRENT_VERSION}\n\n",
            node_type, node_type);
}

static int write_one_event(FILE *f, const char *event,
                           const topo_profile_t *profile) {
    if (ulab_streq(event, "traffic") || ulab_streq(event, "traffic_1gb") ||
        ulab_streq(event, "traffic_2gb_per_ue") ||
        ulab_streq(event, "traffic_5gb_per_ue")) {
        fprintf(f,
                "      - type: traffic\n"
                "        ues: all\n"
                "        amount_mb: %llu\n",
                (unsigned long long)event_amount_mb(event, profile));
        return ULAB_OK;
    }
    if (ulab_streq(event, "traffic_by_profile")) {
        fprintf(f,
                "      - type: traffic_by_profile\n"
                "        profile: mixed_usage\n");
        return ULAB_OK;
    }
    if (ulab_streq(event, "restart_nodes") ||
        ulab_streq(event, "restart_tower")) {
        fprintf(f,
                "      - type: wait_node_connectivity\n"
                "        type_selector: tower\n"
                "        count_per_network: 1\n"
                "        connectivity: Online\n"
                "        seconds: 120\n"
                "      - type: restart_nodes\n"
                "        type_selector: tower\n"
                "        count_per_network: 1\n"
                "      - type: wait_node_connectivity\n"
                "        type_selector: tower\n"
                "        count_per_network: 1\n"
                "        connectivity: Offline\n"
                "        seconds: 120\n"
                "      - type: wait_node_connectivity\n"
                "        type_selector: tower\n"
                "        count_per_network: 1\n"
                "        connectivity: Online\n"
                "        seconds: 300\n"
                "      - type: wait_nodes_ready\n"
                "        type_selector: tower\n"
                "        count_per_network: 1\n"
                "      - type: wait_ues_attached\n"
                "        ues: all\n");
        return ULAB_OK;
    }
    if (ulab_streq(event, "restart_site")) {
        fprintf(f,
                "      - type: wait_node_connectivity\n"
                "        nodes: all\n"
                "        connectivity: Online\n"
                "        seconds: 120\n"
                "      - type: restart_site\n"
                "        nodes: all\n"
                "      - type: wait_node_connectivity\n"
                "        nodes: all\n"
                "        connectivity: Offline\n"
                "        seconds: 120\n"
                "      - type: wait_node_connectivity\n"
                "        nodes: all\n"
                "        connectivity: Online\n"
                "        seconds: 300\n"
                "      - type: wait_nodes_ready\n"
                "        nodes: all\n"
                "      - type: wait_ues_attached\n"
                "        ues: all\n");
        return ULAB_OK;
    }
    if (ulab_streq(event, "toggle_service_off")) {
        fprintf(f,
                "      - type: toggle_service\n"
                "        state: off\n");
        return ULAB_OK;
    }
    if (ulab_streq(event, "toggle_service_on")) {
        fprintf(f,
                "      - type: toggle_service\n"
                "        state: on\n");
        return ULAB_OK;
    }
    if (ulab_streq(event, "toggle_radio_off")) {
        fprintf(f,
                "      - type: toggle_radio\n"
                "        state: off\n");
        return ULAB_OK;
    }
    if (ulab_streq(event, "toggle_radio_on")) {
        fprintf(f,
                "      - type: toggle_radio\n"
                "        state: on\n");
        return ULAB_OK;
    }
    if (ulab_streq(event, "mark_node_offline")) {
        fprintf(f, "      - type: mark_node_offline\n");
        return ULAB_OK;
    }
    if (ulab_streq(event, "restore_node")) {
        fprintf(f, "      - type: restore_node\n");
        return ULAB_OK;
    }
    if (ulab_streq(event, "software_update_tower") ||
        ulab_streq(event, "software_update_amplifier") ||
        ulab_streq(event, "software_update_controller")) {
        fprintf(f,
                "      - type: software_update\n"
                "        type_selector: %s\n"
                "        count_per_network: 1\n"
                "        app: example\n"
                "        tag: ${ULAB_SOFTWARE_TARGET_VERSION}\n",
                software_node_type(event));
        return ULAB_OK;
    }
    if (ulab_streq(event, "software_update_same_version")) {
        fprintf(f,
                "      - type: software_update\n"
                "        type_selector: tower\n"
                "        count_per_network: 1\n"
                "        app: example\n"
                "        tag: current\n");
        return ULAB_OK;
    }
    if (ulab_streq(event, "add_package_to_sim") ||
        ulab_streq(event, "topup")) {
        fprintf(f,
                "      - type: add_package_to_sim\n"
                "        ues: all\n"
                "        package: pkg_002\n");
        return ULAB_OK;
    }
    if (ulab_streq(event, "remove_package_from_sim")) {
        fprintf(f,
                "      - type: remove_package_from_sim\n"
                "        ues: all\n");
        return ULAB_OK;
    }
    if (ulab_streq(event, "set_sim_inactive")) {
        fprintf(f,
                "      - type: set_sim_status\n"
                "        ues: all\n"
                "        status: inactive\n");
        return ULAB_OK;
    }
    if (ulab_streq(event, "set_sim_active")) {
        fprintf(f,
                "      - type: set_sim_status\n"
                "        ues: all\n"
                "        status: active\n");
        return ULAB_OK;
    }
    if (ulab_streq(event, "check")) {
        fprintf(f, "      - type: check\n");
        return ULAB_OK;
    }

    fprintf(f, "      - type: check\n");
    return ULAB_ERR;
}

static void write_check_backend_count(FILE *f, const char *target) {
    fprintf(f,
            "      - type: backend_count\n"
            "        target: %s\n"
            "        expected: from_world\n", target);
}

static void write_one_check(FILE *f, const char *check, const char *family,
                            const gen_case_t *c) {
    if (ulab_streq(check, "backend_count_networks")) {
        write_check_backend_count(f, "networks");
    } else if (ulab_streq(check, "backend_count_sites")) {
        write_check_backend_count(f, "sites");
    } else if (ulab_streq(check, "backend_count_nodes") ||
               ulab_streq(check, "nodes_ready")) {
        write_check_backend_count(f, "nodes");
    } else if (ulab_streq(check, "backend_count_sims") ||
               ulab_streq(check, "backend_count")) {
        write_check_backend_count(f, "sims");
    } else if (ulab_streq(check, "backend_count_subscribers")) {
        write_check_backend_count(f, "subscribers");
    } else if (ulab_streq(check, "backend_count_packages") ||
               ulab_streq(check, "backend_package_count")) {
        write_check_backend_count(f, "packages");
    } else if (ulab_streq(check, "list_contains") ||
               ulab_streq(check, "sim_assigned") ||
               ulab_streq(check, "sim_not_in_unassigned_pool") ||
               ulab_streq(check, "subscriber_has_sim")) {
        fprintf(f,
                "      - type: list_contains\n"
                "        view: sims\n"
                "        ref: ue-000001\n");
    } else if (ulab_streq(check, "usage_per_sim")) {
        fprintf(f,
                "      - type: usage_per_sim\n"
                "        ues: all\n"
                "        expected: from_model\n"
                "        tolerance_percent: 2\n");
    } else if (ulab_streq(check, "balance_non_negative")) {
        fprintf(f,
                "      - type: balance_non_negative\n"
                "        ues: all\n");
    } else if (ulab_streq(check, "ue_attached")) {
        fprintf(f,
                "      - type: ue_attached\n"
                "        ues: all\n");
    } else if (ulab_streq(check, "traffic_allowed")) {
        fprintf(f,
                "      - type: traffic_allowed\n"
                "        ues: all\n"
                "        amount_mb: 1\n");
    } else if (ulab_streq(check, "traffic_blocked")) {
        fprintf(f,
                "      - type: traffic_blocked\n"
                "        ues: all\n"
                "        amount_mb: 1\n");
    } else if (ulab_streq(check, "node_ready")) {
        fprintf(f,
                "      - type: node_ready\n"
                "        nodes: all\n");
    } else if (ulab_streq(check, "node_health_ok")) {
        if (ulab_streq(family, "software_update") &&
            c != NULL && c->event_count > 0) {
            fprintf(f,
                    "      - type: node_health_ok\n"
                    "        type_selector: %s\n"
                    "        count_per_network: 1\n"
                    "        app: example\n",
                    software_node_type(c->events[0]));
        } else {
            fprintf(f, "      - type: node_health_ok\n");
        }
    } else if (ulab_streq(check, "node_version_equals")) {
        fprintf(f,
                "      - type: node_version_equals\n"
                "        type_selector: %s\n"
                "        count_per_network: 1\n"
                "        app: example\n"
                "        version: %s\n",
                c != NULL && c->event_count > 0 ?
                software_node_type(c->events[0]) : "tower",
                ulab_streq(family, "software_update") ?
                "${ULAB_SOFTWARE_TARGET_VERSION}" : "v2.0.0");
    } else if (ulab_streq(check, "history_preserved")) {
        fprintf(f,
                "      - type: history_preserved\n"
                "        entity: %s\n", family);
    } else if (ulab_streq(check, "relationship_exists")) {
        fprintf(f,
                "      - type: relationship_exists\n"
                "        entity: %s\n", family);
    } else if (ulab_streq(check, "relationship_ended")) {
        fprintf(f,
                "      - type: relationship_ended\n"
                "        entity: %s\n", family);
    } else if (ulab_streq(check, "audit_event_exists")) {
        fprintf(f,
                "      - type: audit_event_exists\n"
                "        entity: %s\n", family);
    } else if (ulab_streq(check, "status_equals")) {
        fprintf(f,
                "      - type: status_equals\n"
                "        entity: sim\n"
                "        ref: ue-000001\n"
                "        status: active\n");
    } else if (ulab_streq(check, "package_active")) {
        fprintf(f,
                "      - type: package_active\n"
                "        ues: all\n");
    } else {
        write_check_backend_count(f, "sims");
    }
}

static void write_world(FILE *f, const topo_profile_t *p,
                        int software_only) {
    uint32_t networks;
    uint32_t sites;
    uint32_t ues;

    networks = p->networks ? p->networks : 1;
    sites = p->sites_per_network ? p->sites_per_network : 1;
    ues = software_only ? 0 :
          (p->ues_per_site ? p->ues_per_site : 1);

    fprintf(f,
            "provider:\n"
            "  type: virtual\n\n"
            "world:\n"
            "  networks: %u\n"
            "  sites_per_network: %u\n"
            "  nodes_per_site:\n"
            "    tower: 1\n"
            "    amplifier: 1\n"
            "    controller: 1\n"
            "  ues_per_site: %u\n\n",
            networks, sites, ues);
}

static void write_setup(FILE *f, int software_only) {
    fprintf(f,
            "setup:\n"
            "  create_via_bff:\n"
            "    - networks\n"
            "    - sites\n"
            "    - nodes\n"
            "    - node_site_links\n");
    if (!software_only) {
        fprintf(f,
                "    - packages\n"
                "    - subscribers\n"
                "    - sims\n");
    }
    fprintf(f, "\n");
}

static void write_runtime(FILE *f, int needs_ues) {
    fprintf(f,
            "runtime:\n"
            "  start: %s\n"
            "  wait: %s\n\n",
            needs_ues ? "[nodes, ues]" : "[nodes]",
            needs_ues ? "[nodes_ready, ues_attached]" : "[nodes_ready]");
}


static void normalized_tags(const gen_case_t *c, char *out, size_t out_len) {
    const char *p;
    size_t len;

    if (out_len == 0) {
        return;
    }
    out[0] = '\0';
    if (c == NULL || c->tags[0] == '\0') {
        return;
    }

    p = c->tags;
    while (*p == ' ' || *p == '[') {
        p++;
    }
    len = strlen(p);
    while (len > 0 && (p[len - 1] == ' ' || p[len - 1] == ']')) {
        len--;
    }
    if (len >= out_len) {
        len = out_len - 1;
    }
    memcpy(out, p, len);
    out[len] = '\0';
}

static int write_case_scenario(const gen_opts_t *opts,
                               const gen_family_t *family,
                               const gen_case_t *c,
                               const topo_profile_t *p,
                               FILE *index,
                               size_t *active_count,
                               size_t *wip_count) {
    char dir[ULAB_MAX_PATH];
    char path[ULAB_MAX_PATH];
    char tags[GEN_MAX_TAGS];
    uint32_t seed;
    uint32_t packages;
    int wip;
    int needs_ues;
    int software_only;
    size_t i;
    FILE *f;

    if (path_join(dir, sizeof(dir), opts->out_dir, family->name)) {
        return ULAB_ERR;
    }
    if (ulab_mkdir_p(dir)) {
        return ULAB_ERR;
    }
    if (path_join(path, sizeof(path), dir, c->name)) {
        return ULAB_ERR;
    }
    if (append_str(path, sizeof(path), ".yaml")) {
        return ULAB_ERR;
    }

    normalized_tags(c, tags, sizeof(tags));
    wip = case_is_wip(family, c);
    software_only = software_update_success_case(family, c);
    needs_ues = software_only ? 0 : case_needs_ues(family, c);
    packages = scenario_package_count(c, p);
    seed = ulab_hash32(family->name, 62000);
    seed = ulab_hash32(c->name, seed);

    f = fopen(path, "w");
    if (f == NULL) {
        fprintf(stderr, "unable to write scenario: %s\n", path);
        return ULAB_ERR;
    }

    fprintf(f,
            "version: 1\n"
            "name: %s\n",
            c->name);
    if (software_only) {
        write_software_description(f, c);
    }
    fprintf(f,
            "seed: %u\n"
            "suite: generated\n"
            "priority: %s\n"
            "tags: [generated, %s%s%s]\n"
            "status: %s\n"
            "generated: true\n"
            "entity: %s\n"
            "action: %s\n\n"
            "# family: %s\n"
            "# case: %s\n"
            "# topology: %s\n"
            "# generated_status: %s\n\n",
            seed, c->priority[0] ? c->priority : family->priority,
            family->name, tags[0] ? ", " : "", tags,
            wip ? "wip" : "active",
            family->name, c->event_count ? c->events[0] : "check",
            family->name, c->name, p->name, wip ? "wip" : "active");

    write_world(f, p, software_only);
    if (!software_only) {
        write_packages(f, packages, p->traffic_mb_per_ue);
    }
    write_setup(f, software_only);
    write_runtime(f, needs_ues);
    write_profile_section(f, c);

    fprintf(f, "phases:\n");
    if (!wip && software_only) {
        write_software_preflight(f, c);
    }
    fprintf(f,
            "  - name: action\n"
            "    events:\n");
    write_setup_preamble(f, c);
    if (c->event_count == 0 || wip) {
        fprintf(f, "      - type: check\n");
    } else {
        for (i = 0; i < c->event_count; i++) {
            write_one_event(f, c->events[i], p);
        }
    }

    fprintf(f, "    checks:\n");
    if (c->check_count == 0 || wip) {
        write_check_backend_count(f, "sims");
    } else {
        for (i = 0; i < c->check_count; i++) {
            write_one_check(f, c->checks[i], family->name, c);
        }
    }

    if (!software_only) {
        fprintf(f,
                "\nfinal_checks:\n"
                "  - type: balance_non_negative\n"
                "    ues: all\n");
    }

    fclose(f);

    fprintf(index,
            "  - file: %s/%s.yaml\n"
            "    family: %s\n"
            "    case: %s\n"
            "    topology: %s\n"
            "    status: %s\n"
            "    priority: %s\n",
            family->name, c->name, family->name, c->name, p->name,
            wip ? "wip" : "active",
            c->priority[0] ? c->priority : family->priority);

    if (wip) {
        (*wip_count)++;
    } else {
        (*active_count)++;
    }
    return ULAB_OK;
}

static int generate_family(const gen_opts_t *opts,
                           const gen_catalog_t *cat,
                           const gen_family_t *family,
                           FILE *index,
                           size_t *active_count,
                           size_t *wip_count) {
    size_t i;

    for (i = 0; i < family->case_count; i++) {
        const topo_profile_t *p;
        p = find_profile(cat, family->cases[i].topology);
        if (p == NULL) {
            fprintf(stderr, "unknown topology %.128s in case %.128s\n",
                    family->cases[i].topology, family->cases[i].name);
            return ULAB_ERR;
        }
        if (write_case_scenario(opts, family, &family->cases[i], p, index,
                                active_count, wip_count)) {
            return ULAB_ERR;
        }
    }
    return ULAB_OK;
}

static void write_index_header(FILE *index, const gen_catalog_t *cat) {
    size_t i;

    fprintf(index, "catalog:\n");
    fprintf(index, "  topologies:\n");
    for (i = 0; i < cat->profile_count; i++) {
        fprintf(index,
                "    - name: %s\n"
                "      networks: %u\n"
                "      sites_per_network: %u\n"
                "      ues_per_site: %u\n"
                "      packages: %u\n"
                "      traffic_mb_per_ue: %llu\n",
                cat->profiles[i].name,
                cat->profiles[i].networks,
                cat->profiles[i].sites_per_network,
                cat->profiles[i].ues_per_site,
                cat->profiles[i].packages,
                (unsigned long long)cat->profiles[i].traffic_mb_per_ue);
    }
    fprintf(index, "generated_scenarios:\n");
}

int generator_run(int argc, char **argv) {
    gen_opts_t opts;
    gen_catalog_t *cat;
    char index_path[ULAB_MAX_PATH];
    FILE *index;
    size_t i;
    size_t active_count;
    size_t wip_count;

    if (parse_opts(argc, argv, &opts) != ULAB_OK) {
        usage();
        return ULAB_EUSAGE;
    }

    cat = calloc(1, sizeof(*cat));
    if (cat == NULL) {
        return ULAB_ERR;
    }

    if (load_catalog(&opts, cat)) {
        free(cat);
        return ULAB_ERR;
    }

    if (ulab_mkdir_p(opts.out_dir)) {
        free(cat);
        fprintf(stderr, "unable to create output dir: %s\n", opts.out_dir);
        return ULAB_ERR;
    }
    if (path_join(index_path, sizeof(index_path), opts.out_dir, "index.yaml")) {
        free(cat);
        return ULAB_ERR;
    }
    index = fopen(index_path, "w");
    if (index == NULL) {
        free(cat);
        fprintf(stderr, "unable to write: %s\n", index_path);
        return ULAB_ERR;
    }

    write_index_header(index, cat);
    active_count = 0;
    wip_count = 0;
    for (i = 0; i < cat->family_count; i++) {
        if (generate_family(&opts, cat, &cat->families[i], index,
                            &active_count, &wip_count)) {
            fclose(index);
            free(cat);
            return ULAB_ERR;
        }
    }

    fclose(index);
    printf("generated cases active=%zu wip=%zu index=%s\n",
           active_count, wip_count, index_path);
    free(cat);
    return ULAB_OK;
}
