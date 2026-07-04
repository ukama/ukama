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
#include <sys/stat.h>
#include <unistd.h>

#include "generator.h"
#include "scenario.h"
#include "util.h"

#define GEN_MAX_ACTIONS       64
#define GEN_MAX_LIST          32
#define GEN_MAX_ITEMS         128
#define GEN_MAX_PROFILES      16
#define GEN_MAX_FAMILIES      32
#define GEN_MAX_PACKAGES_OUT  10

typedef struct {
    char model[ULAB_MAX_REF];
    char out_dir[ULAB_MAX_PATH];
    char models_dir[ULAB_MAX_PATH];
} gen_opts_t;

typedef struct {
    char name[ULAB_MAX_REF];
    char event[ULAB_MAX_REF];
    char status[ULAB_MAX_REF];
    char package_ref[ULAB_MAX_REF];
    char selector[ULAB_MAX_REF];
    char from[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t from_count;
    char to[ULAB_MAX_REF];
    char blocked_from[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t blocked_count;
    char guards[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t guard_count;
    char checks[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t check_count;
    char runtime[ULAB_MAX_REF];
    char priority[ULAB_MAX_REF];
} action_rule_t;

typedef struct {
    char entity[ULAB_MAX_REF];
    char states[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t state_count;
    action_rule_t actions[GEN_MAX_ACTIONS];
    size_t action_count;
} model_def_t;

typedef struct {
    char name[ULAB_MAX_REF];
    char mode[ULAB_MAX_REF];
    char expect[ULAB_MAX_REF];
    uint32_t networks;
    uint32_t sites_per_network;
    uint32_t ues_per_site;
    uint32_t packages;
    uint32_t sim_pool;
} matrix_item_t;

typedef struct {
    char name[ULAB_MAX_REF];
    char kind[ULAB_MAX_REF];
    matrix_item_t items[GEN_MAX_ITEMS];
    size_t item_count;
} matrix_profile_t;

typedef struct {
    char name[ULAB_MAX_REF];
    char priority[ULAB_MAX_REF];
    char flows[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t flow_count;
    char topologies[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t topology_count;
    char scales[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t scale_count;
    char runtime[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t runtime_count;
    char failures[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t failure_count;
    char software[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t software_count;
    char verification[GEN_MAX_LIST][ULAB_MAX_REF];
    size_t verification_count;
} matrix_family_t;

typedef struct {
    matrix_profile_t profiles[GEN_MAX_PROFILES];
    size_t profile_count;
    matrix_family_t families[GEN_MAX_FAMILIES];
    size_t family_count;
} matrix_def_t;

static const char *default_entities[] = {
    "sim", "package", "subscriber", "node"
};

static const char *default_profiles[] = {
    "topology", "scale", "runtime", "failure", "software", "verification"
};

static const char *default_families[] = {
    "smoke", "usage", "sim_pool", "package", "subscriber", "lifecycle",
    "node_ops", "site_ops", "software_update", "failure", "dashboard"
};

static void usage(void) {
    printf("usage:\n");
    printf("  ukama-lab generate --model <name|all> [options]\n");
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
    if (used >= n) {
        return ULAB_ERR;
    }
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

    for (i = 0; s[i] != '\0'; i++) {
        if (s[i] == ' ' || s[i] == '/' || s[i] == ':' || s[i] == ',') {
            s[i] = '_';
        }
    }
}

static int build_path3(char *out, size_t n,
                       const char *a, const char *b, const char *c) {
    out[0] = '\0';
    if (append_str(out, n, a)) return ULAB_ERR;
    if (append_str(out, n, "/")) return ULAB_ERR;
    if (append_str(out, n, b)) return ULAB_ERR;
    if (append_str(out, n, "/")) return ULAB_ERR;
    if (append_str(out, n, c)) return ULAB_ERR;
    return ULAB_OK;
}

static int build_yaml_path(char *out, size_t n,
                           const char *dir, const char *name) {
    out[0] = '\0';
    if (append_str(out, n, dir)) return ULAB_ERR;
    if (append_str(out, n, "/")) return ULAB_ERR;
    if (append_str(out, n, name)) return ULAB_ERR;
    if (append_str(out, n, ".yaml")) return ULAB_ERR;
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
        } else if (ulab_streq(argv[i], "--mode") && i + 1 < argc) {
            i++;
        } else if (ulab_streq(argv[i], "--out") && i + 1 < argc) {
            if (ulab_copy(opts->out_dir, sizeof(opts->out_dir), argv[++i])) {
                return ULAB_EUSAGE;
            }
        } else if (ulab_streq(argv[i], "--models") && i + 1 < argc) {
            if (ulab_copy(opts->models_dir, sizeof(opts->models_dir), argv[++i])) {
                return ULAB_EUSAGE;
            }
        } else if (ulab_streq(argv[i], "--templates") && i + 1 < argc) {
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

static int parse_list(char arr[GEN_MAX_LIST][ULAB_MAX_REF], size_t *count,
                      char *val) {
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

static action_rule_t *new_action(model_def_t *m, const char *name) {
    action_rule_t *a;

    if (m->action_count >= GEN_MAX_ACTIONS) {
        return NULL;
    }
    a = &m->actions[m->action_count++];
    memset(a, 0, sizeof(*a));
    if (ulab_copy(a->name, sizeof(a->name), name)) {
        return NULL;
    }
    clean_token(a->name);
    ulab_copy(a->event, sizeof(a->event), "check");
    ulab_copy(a->priority, sizeof(a->priority), "p2");
    return a;
}

static int parse_action_field(action_rule_t *a, const char *key, char *val) {
    if (ulab_streq(key, "event")) return ulab_copy(a->event, sizeof(a->event), val);
    if (ulab_streq(key, "status")) return ulab_copy(a->status, sizeof(a->status), val);
    if (ulab_streq(key, "package")) return ulab_copy(a->package_ref, sizeof(a->package_ref), val);
    if (ulab_streq(key, "selector")) return ulab_copy(a->selector, sizeof(a->selector), val);
    if (ulab_streq(key, "from")) return parse_list(a->from, &a->from_count, val);
    if (ulab_streq(key, "to")) return ulab_copy(a->to, sizeof(a->to), val);
    if (ulab_streq(key, "blocked_from")) return parse_list(a->blocked_from, &a->blocked_count, val);
    if (ulab_streq(key, "guards")) return parse_list(a->guards, &a->guard_count, val);
    if (ulab_streq(key, "checks")) return parse_list(a->checks, &a->check_count, val);
    if (ulab_streq(key, "runtime")) return ulab_copy(a->runtime, sizeof(a->runtime), val);
    if (ulab_streq(key, "priority")) return ulab_copy(a->priority, sizeof(a->priority), val);
    return ULAB_ERR;
}

static int model_yaml_path(char *out, size_t n,
                           const gen_opts_t *opts, const char *entity) {
    out[0] = '\0';
    if (append_str(out, n, opts->models_dir)) return ULAB_ERR;
    if (append_str(out, n, "/")) return ULAB_ERR;
    if (append_str(out, n, entity)) return ULAB_ERR;
    if (append_str(out, n, ".yaml")) return ULAB_ERR;
    return ULAB_OK;
}

static int load_model(const gen_opts_t *opts, const char *entity,
                      model_def_t *model) {
    char path[ULAB_MAX_PATH];
    char line[ULAB_MAX_LINE];
    FILE *fp;
    int in_actions;
    action_rule_t *cur;

    memset(model, 0, sizeof(*model));
    in_actions = 0;
    cur = NULL;

    if (model_yaml_path(path, sizeof(path), opts, entity)) {
        return ULAB_ERR;
    }

    fp = fopen(path, "r");
    if (fp == NULL) {
        fprintf(stderr, "missing model: %s\n", path);
        return ULAB_ERR;
    }

    while (fgets(line, sizeof(line), fp) != NULL) {
        char *key;
        char *val;
        char *p;

        p = ulab_trim(line);
        if (*p == '\0' || *p == '#') {
            continue;
        }
        if (ulab_starts(p, "- name:")) {
            val = ulab_trim(strchr(p, ':') + 1);
            cur = new_action(model, val);
            if (cur == NULL) {
                fclose(fp);
                return ULAB_ERR;
            }
            continue;
        }
        if (strchr(p, ':') == NULL) {
            continue;
        }
        key = p;
        val = strchr(p, ':');
        *val++ = '\0';
        key = ulab_trim(key);
        val = ulab_trim(val);

        if (ulab_streq(key, "entity")) {
            if (ulab_copy(model->entity, sizeof(model->entity), val)) {
                fclose(fp);
                return ULAB_ERR;
            }
            clean_token(model->entity);
            in_actions = 0;
        } else if (ulab_streq(key, "states")) {
            if (val[0] != '\0') {
                if (parse_list(model->states, &model->state_count, val)) {
                    fclose(fp);
                    return ULAB_ERR;
                }
            }
            in_actions = 0;
        } else if (ulab_starts(p, "- ") && !in_actions) {
            if (add_word(model->states, &model->state_count, ulab_trim(p + 2))) {
                fclose(fp);
                return ULAB_ERR;
            }
        } else if (ulab_streq(key, "actions")) {
            in_actions = 1;
        } else if (in_actions && cur != NULL) {
            if (parse_action_field(cur, key, val)) {
                fprintf(stderr, "unknown action field in %s: %s\n", path, key);
                fclose(fp);
                return ULAB_ERR;
            }
        }
    }

    fclose(fp);

    if (!ulab_streq(model->entity, entity)) {
        fprintf(stderr, "model entity mismatch: %s\n", path);
        return ULAB_ERR;
    }
    if (model->action_count == 0) {
        fprintf(stderr, "model has no actions: %s\n", path);
        return ULAB_ERR;
    }

    return ULAB_OK;
}

static int has_check(const action_rule_t *a, const char *name) {
    size_t i;

    for (i = 0; i < a->check_count; i++) {
        if (ulab_streq(a->checks[i], name)) {
            return 1;
        }
    }
    return 0;
}

static const char *action_status(const action_rule_t *a) {
    if (a->status[0] != '\0') {
        return a->status;
    }
    if (a->to[0] != '\0') {
        if (ulab_streq(a->to, "inactive") || ulab_streq(a->to, "deactivated") ||
            ulab_streq(a->to, "released") || ulab_streq(a->to, "suspended") ||
            ulab_streq(a->to, "retired") || ulab_streq(a->to, "archived")) {
            return "inactive";
        }
    }
    return "active";
}

static void write_common_packages(FILE *f, uint32_t package_count) {
    uint32_t i;
    uint32_t n;

    n = package_count;
    if (n == 0) n = 1;
    if (n > GEN_MAX_PACKAGES_OUT) n = GEN_MAX_PACKAGES_OUT;

    fprintf(f, "packages:\n");
    for (i = 0; i < n; i++) {
        fprintf(f,
                "  - ref: pkg_%03u\n"
                "    name: Package %03u\n"
                "    data_mb: %u\n"
                "    duration_days: %u\n"
                "    amount: %u.00\n"
                "    assign_percent: %u\n",
                i + 1, i + 1, 1024u * (i + 1), i + 1, i + 1,
                i == 0 ? 100u : 0u);
    }
    fprintf(f, "\n");
}

static void write_common_setup(FILE *f) {
    fprintf(f,
            "setup:\n"
            "  create_via_bff:\n"
            "    - networks\n"
            "    - sites\n"
            "    - nodes\n"
            "    - node_site_links\n"
            "    - packages\n"
            "    - subscribers\n"
            "    - sims\n\n"
            "runtime:\n"
            "  start: [nodes, ues]\n"
            "  wait: [nodes_ready, ues_attached]\n\n");
}

static void write_base_world(FILE *f, uint32_t networks,
                             uint32_t sites_per_network,
                             uint32_t ues_per_site,
                             uint32_t package_count) {
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
            networks ? networks : 1u,
            sites_per_network ? sites_per_network : 1u,
            ues_per_site ? ues_per_site : 1u);
    write_common_packages(f, package_count);
    write_common_setup(f);
}

static void write_action_event(FILE *f, const model_def_t *m,
                               const action_rule_t *a, int blocked) {
    (void)m;

    if (blocked) {
        fprintf(f,
                "      - type: set_sim_status\n"
                "        ues: all\n"
                "        expect:\n"
                "          result: failure\n"
                "          error_contains: \"missing status\"\n");
        return;
    }

    if (ulab_streq(a->event, "set_sim_status")) {
        fprintf(f,
                "      - type: set_sim_status\n"
                "        ues: all\n"
                "        status: %s\n", action_status(a));
    } else if (ulab_streq(a->event, "add_package_to_sim")) {
        fprintf(f,
                "      - type: add_package_to_sim\n"
                "        ues: all\n"
                "        package: %s\n",
                a->package_ref[0] ? a->package_ref : "pkg_002");
    } else if (ulab_streq(a->event, "remove_package_from_sim")) {
        fprintf(f,
                "      - type: remove_package_from_sim\n"
                "        ues: all\n");
    } else if (ulab_streq(a->event, "restart_nodes")) {
        fprintf(f,
                "      - type: restart_nodes\n"
                "        type_selector: %s\n"
                "        count_per_network: 1\n",
                a->selector[0] ? a->selector : "tower");
    } else {
        fprintf(f, "      - type: check\n");
    }
}

static void write_basic_checks(FILE *f, const char *target) {
    fprintf(f,
            "    checks:\n"
            "      - type: backend_count\n"
            "        target: %s\n"
            "        expected: from_world\n",
            target != NULL && target[0] ? target : "sims");
}

static void write_action_checks(FILE *f, const model_def_t *m,
                                const action_rule_t *a, int blocked) {
    write_basic_checks(f, ulab_streq(m->entity, "node") ? "nodes" : "sims");

    if (!blocked && has_check(a, "list_contains")) {
        fprintf(f,
                "      - type: list_contains\n"
                "        view: sims\n"
                "        ref: ue-000001\n");
    }
    if (!blocked && has_check(a, "status_equals") &&
        (ulab_streq(m->entity, "sim") || ulab_streq(m->entity, "subscriber"))) {
        fprintf(f,
                "      - type: status_equals\n"
                "        entity: sim\n"
                "        ref: ue-000001\n"
                "        status: %s\n", action_status(a));
    }
    if (!blocked && ulab_streq(a->runtime, "traffic_allowed")) {
        fprintf(f,
                "      - type: traffic_allowed\n"
                "        ues: all\n"
                "        amount_mb: 1\n");
    } else if (!blocked && ulab_streq(a->runtime, "traffic_blocked")) {
        fprintf(f,
                "      - type: traffic_blocked\n"
                "        ues: all\n"
                "        amount_mb: 1\n");
    }
}

static void build_entity_name(char *out, size_t n,
                              const char *entity, const char *action,
                              const char *case_name,
                              const char *from_state) {
    out[0] = '\0';
    append_str(out, n, entity);
    append_str(out, n, "-");
    append_str(out, n, action);
    append_str(out, n, "-");
    append_str(out, n, case_name);
    append_str(out, n, "-");
    append_str(out, n, from_state && from_state[0] ? from_state : "default");
    clean_token(out);
}

static int write_entity_scenario(const gen_opts_t *opts, const model_def_t *m,
                                 const action_rule_t *a, const char *case_name,
                                 const char *from_state, int blocked,
                                 FILE *index) {
    char dir[ULAB_MAX_PATH];
    char path[ULAB_MAX_PATH];
    char name[ULAB_MAX_NAME];
    uint32_t seed;
    FILE *f;

    if (build_path3(dir, sizeof(dir), opts->out_dir, m->entity, "")) {
        return ULAB_ERR;
    }
    if (dir[0] != '\0' && dir[strlen(dir) - 1] == '/') {
        dir[strlen(dir) - 1] = '\0';
    }
    if (ulab_mkdir_p(dir)) {
        return ULAB_ERR;
    }

    build_entity_name(name, sizeof(name), m->entity, a->name, case_name,
                      from_state);
    if (build_yaml_path(path, sizeof(path), dir, name)) {
        return ULAB_ERR;
    }

    seed = ulab_hash32(m->entity, blocked ? 11000 : 9000);
    seed = ulab_hash32(a->name, seed);
    seed = ulab_hash32(from_state ? from_state : "default", seed);

    f = fopen(path, "w");
    if (f == NULL) {
        fprintf(stderr, "unable to write: %s\n", path);
        return ULAB_ERR;
    }

    fprintf(f,
            "version: 1\n"
            "name: %s\n"
            "seed: %u\n"
            "suite: generated\n"
            "priority: %s\n"
            "tags: [generated, entity, %s, %s, %s%s]\n"
            "status: active\n"
            "generated: true\n"
            "entity: %s\n"
            "action: %s\n\n",
            name, seed, a->priority[0] ? a->priority : "p2",
            m->entity, a->name, case_name, blocked ? ", negative" : "",
            m->entity, a->name);

    fprintf(f,
            "# model_from: %s\n"
            "# model_to: %s\n"
            "# model_expected: %s\n\n",
            from_state && from_state[0] ? from_state : "default",
            a->to[0] ? a->to : "unchanged",
            blocked ? "failure" : "success");

    write_base_world(f, 1, 1, 1, 2);
    fprintf(f,
            "phases:\n"
            "  - name: action\n"
            "    events:\n");
    write_action_event(f, m, a, blocked);
    write_action_checks(f, m, a, blocked);

    fprintf(f,
            "\nfinal_checks:\n"
            "  - type: balance_non_negative\n"
            "    ues: all\n");

    fclose(f);

    fprintf(index,
            "  - file: %s/%s.yaml\n"
            "    entity: %s\n"
            "    action: %s\n"
            "    case: %s\n"
            "    from: %s\n"
            "    expected: %s\n"
            "    priority: %s\n",
            m->entity, name, m->entity, a->name, case_name,
            from_state && from_state[0] ? from_state : "default",
            blocked ? "failure" : "success",
            a->priority[0] ? a->priority : "p2");

    return ULAB_OK;
}

static int generate_action(const gen_opts_t *opts, const model_def_t *m,
                           const action_rule_t *a, FILE *index) {
    size_t i;

    if (a->from_count == 0) {
        if (write_entity_scenario(opts, m, a, "success", "default", 0, index)) {
            return ULAB_ERR;
        }
    } else {
        for (i = 0; i < a->from_count; i++) {
            if (write_entity_scenario(opts, m, a, "success", a->from[i], 0, index)) {
                return ULAB_ERR;
            }
        }
    }

    for (i = 0; i < a->blocked_count; i++) {
        if (write_entity_scenario(opts, m, a, "blocked", a->blocked_from[i], 1, index)) {
            return ULAB_ERR;
        }
    }

    for (i = 0; i < a->guard_count; i++) {
        if (write_entity_scenario(opts, m, a, "guard", a->guards[i], 1, index)) {
            return ULAB_ERR;
        }
    }

    return ULAB_OK;
}

static int matrix_path(char *out, size_t n, const gen_opts_t *opts,
                       const char *subdir, const char *name) {
    out[0] = '\0';
    if (append_str(out, n, opts->models_dir)) return ULAB_ERR;
    if (append_str(out, n, "/")) return ULAB_ERR;
    if (append_str(out, n, subdir)) return ULAB_ERR;
    if (append_str(out, n, "/")) return ULAB_ERR;
    if (append_str(out, n, name)) return ULAB_ERR;
    if (append_str(out, n, ".yaml")) return ULAB_ERR;
    return ULAB_OK;
}

static matrix_profile_t *new_profile(matrix_def_t *matrix,
                                     const char *name) {
    matrix_profile_t *p;

    if (matrix->profile_count >= GEN_MAX_PROFILES) {
        return NULL;
    }
    p = &matrix->profiles[matrix->profile_count++];
    memset(p, 0, sizeof(*p));
    if (ulab_copy(p->name, sizeof(p->name), name)) {
        return NULL;
    }
    return p;
}

static matrix_family_t *new_family(matrix_def_t *matrix,
                                   const char *name) {
    matrix_family_t *f;

    if (matrix->family_count >= GEN_MAX_FAMILIES) {
        return NULL;
    }
    f = &matrix->families[matrix->family_count++];
    memset(f, 0, sizeof(*f));
    if (ulab_copy(f->name, sizeof(f->name), name)) {
        return NULL;
    }
    clean_token(f->name);
    ulab_copy(f->priority, sizeof(f->priority), "p2");
    return f;
}

static matrix_item_t *profile_add_item(matrix_profile_t *profile,
                                       const char *name) {
    matrix_item_t *item;

    if (profile->item_count >= GEN_MAX_ITEMS) {
        return NULL;
    }
    item = &profile->items[profile->item_count++];
    memset(item, 0, sizeof(*item));
    if (ulab_copy(item->name, sizeof(item->name), name)) {
        return NULL;
    }
    clean_token(item->name);
    return item;
}

static int parse_profile_field(matrix_item_t *item,
                               const char *key,
                               char *val) {
    if (ulab_streq(key, "mode") || ulab_streq(key, "condition")) {
        return ulab_copy(item->mode, sizeof(item->mode), val);
    }
    if (ulab_streq(key, "expect")) {
        return ulab_copy(item->expect, sizeof(item->expect), val);
    }
    if (ulab_streq(key, "networks")) {
        return ulab_parse_u32(val, &item->networks);
    }
    if (ulab_streq(key, "sites_per_network")) {
        return ulab_parse_u32(val, &item->sites_per_network);
    }
    if (ulab_streq(key, "ues_per_site")) {
        return ulab_parse_u32(val, &item->ues_per_site);
    }
    if (ulab_streq(key, "packages")) {
        return ulab_parse_u32(val, &item->packages);
    }
    if (ulab_streq(key, "sim_pool")) {
        return ulab_parse_u32(val, &item->sim_pool);
    }
    return ULAB_OK;
}

static int load_profile(const gen_opts_t *opts, const char *name,
                        matrix_def_t *matrix) {
    char path[ULAB_MAX_PATH];
    char line[ULAB_MAX_LINE];
    FILE *fp;
    matrix_profile_t *profile;
    matrix_item_t *cur;

    if (matrix_path(path, sizeof(path), opts, "profiles", name)) {
        return ULAB_ERR;
    }
    fp = fopen(path, "r");
    if (fp == NULL) {
        fprintf(stderr, "missing matrix profile: %s\n", path);
        return ULAB_ERR;
    }
    profile = new_profile(matrix, name);
    if (profile == NULL) {
        fclose(fp);
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
            val = ulab_trim(strchr(p, ':') + 1);
            cur = profile_add_item(profile, val);
            if (cur == NULL) {
                fclose(fp);
                return ULAB_ERR;
            }
            continue;
        }
        if (strchr(p, ':') == NULL) {
            continue;
        }
        key = p;
        val = strchr(p, ':');
        *val++ = '\0';
        key = ulab_trim(key);
        val = ulab_trim(val);

        if (ulab_streq(key, "profile_type")) {
            if (ulab_copy(profile->kind, sizeof(profile->kind), val)) {
                fclose(fp);
                return ULAB_ERR;
            }
            clean_token(profile->kind);
        } else if (cur != NULL) {
            if (parse_profile_field(cur, key, val)) {
                fclose(fp);
                return ULAB_ERR;
            }
        }
    }

    fclose(fp);

    if (profile->kind[0] == '\0' || profile->item_count == 0) {
        fprintf(stderr, "invalid matrix profile: %s\n", path);
        return ULAB_ERR;
    }
    return ULAB_OK;
}

static int parse_family_field(matrix_family_t *family,
                              const char *key,
                              char *val) {
    if (ulab_streq(key, "priority")) {
        return ulab_copy(family->priority, sizeof(family->priority), val);
    }
    if (ulab_streq(key, "flows")) {
        return parse_list(family->flows, &family->flow_count, val);
    }
    if (ulab_streq(key, "topologies")) {
        return parse_list(family->topologies, &family->topology_count, val);
    }
    if (ulab_streq(key, "scales")) {
        return parse_list(family->scales, &family->scale_count, val);
    }
    if (ulab_streq(key, "runtime")) {
        return parse_list(family->runtime, &family->runtime_count, val);
    }
    if (ulab_streq(key, "failures")) {
        return parse_list(family->failures, &family->failure_count, val);
    }
    if (ulab_streq(key, "software")) {
        return parse_list(family->software, &family->software_count, val);
    }
    if (ulab_streq(key, "verification")) {
        return parse_list(family->verification, &family->verification_count, val);
    }
    return ULAB_OK;
}

static int load_family(const gen_opts_t *opts, const char *name,
                       matrix_def_t *matrix) {
    char path[ULAB_MAX_PATH];
    char line[ULAB_MAX_LINE];
    FILE *fp;
    matrix_family_t *family;

    if (matrix_path(path, sizeof(path), opts, "families", name)) {
        return ULAB_ERR;
    }
    fp = fopen(path, "r");
    if (fp == NULL) {
        fprintf(stderr, "missing scenario family: %s\n", path);
        return ULAB_ERR;
    }
    family = NULL;

    while (fgets(line, sizeof(line), fp) != NULL) {
        char *p;
        char *key;
        char *val;

        p = ulab_trim(line);
        if (*p == '\0' || *p == '#') {
            continue;
        }
        if (strchr(p, ':') == NULL) {
            continue;
        }
        key = p;
        val = strchr(p, ':');
        *val++ = '\0';
        key = ulab_trim(key);
        val = ulab_trim(val);

        if (ulab_streq(key, "family")) {
            family = new_family(matrix, val);
            if (family == NULL) {
                fclose(fp);
                return ULAB_ERR;
            }
        } else if (family != NULL) {
            if (parse_family_field(family, key, val)) {
                fclose(fp);
                return ULAB_ERR;
            }
        }
    }

    fclose(fp);

    if (family == NULL || family->flow_count == 0 ||
        family->topology_count == 0 || family->scale_count == 0 ||
        family->runtime_count == 0 || family->failure_count == 0 ||
        family->verification_count == 0) {
        fprintf(stderr, "invalid scenario family: %s\n", path);
        return ULAB_ERR;
    }
    return ULAB_OK;
}

static const matrix_profile_t *find_profile(const matrix_def_t *matrix,
                                            const char *kind) {
    size_t i;

    for (i = 0; i < matrix->profile_count; i++) {
        if (ulab_streq(matrix->profiles[i].name, kind) ||
            ulab_streq(matrix->profiles[i].kind, kind)) {
            return &matrix->profiles[i];
        }
    }
    return NULL;
}

static const matrix_item_t *profile_item(const matrix_def_t *matrix,
                                         const char *kind,
                                         const char *name) {
    const matrix_profile_t *p;
    size_t i;

    p = find_profile(matrix, kind);
    if (p == NULL) {
        return NULL;
    }
    for (i = 0; i < p->item_count; i++) {
        if (ulab_streq(p->items[i].name, name)) {
            return &p->items[i];
        }
    }
    return NULL;
}

static int load_matrix_models(const gen_opts_t *opts, matrix_def_t *matrix) {
    size_t i;

    memset(matrix, 0, sizeof(*matrix));
    for (i = 0; i < sizeof(default_profiles) / sizeof(default_profiles[0]); i++) {
        if (load_profile(opts, default_profiles[i], matrix)) {
            return ULAB_ERR;
        }
    }
    for (i = 0; i < sizeof(default_families) / sizeof(default_families[0]); i++) {
        if (load_family(opts, default_families[i], matrix)) {
            return ULAB_ERR;
        }
    }
    return ULAB_OK;
}

static int is_runtime_supported(const char *runtime) {
    return ulab_streq(runtime, "normal") || ulab_streq(runtime, "node_restart");
}

static int is_verification_supported(const char *verification) {
    return ulab_streq(verification, "bff_mutation") ||
           ulab_streq(verification, "read_model") ||
           ulab_streq(verification, "dashboard") ||
           ulab_streq(verification, "runtime");
}

static int is_software_supported(const char *software) {
    return ulab_streq(software, "no_update");
}

static int is_family_flow_supported(const char *family, const char *flow) {
    if (ulab_streq(family, "software_update")) return 0;
    if (ulab_streq(family, "site_ops")) {
        return ulab_streq(flow, "restart_site");
    }
    if (ulab_streq(family, "lifecycle")) {
        return ulab_streq(flow, "sim_suspend") ||
               ulab_streq(flow, "sim_deactivate");
    }
    return 1;
}

static int should_mark_wip(const char *family, const char *flow,
                           const char *runtime, const char *failure,
                           const char *software, const char *verification) {
    if (!is_family_flow_supported(family, flow)) return 1;
    if (!is_runtime_supported(runtime)) return 1;
    if (!ulab_streq(failure, "none")) return 1;
    if (!is_software_supported(software)) return 1;
    if (!is_verification_supported(verification)) return 1;
    return 0;
}

static void build_matrix_name(char *out, size_t n,
                              const char *family, const char *flow,
                              const char *topology, const char *scale,
                              const char *runtime, const char *failure,
                              const char *software, const char *verification) {
    out[0] = '\0';
    append_str(out, n, family);
    append_str(out, n, "-");
    append_str(out, n, flow);
    append_str(out, n, "-");
    append_str(out, n, topology);
    append_str(out, n, "-");
    append_str(out, n, scale);
    append_str(out, n, "-");
    append_str(out, n, runtime);
    append_str(out, n, "-");
    append_str(out, n, failure);
    if (!ulab_streq(software, "no_update")) {
        append_str(out, n, "-");
        append_str(out, n, software);
    }
    append_str(out, n, "-");
    append_str(out, n, verification);
    clean_token(out);
}

static const char *matrix_entity_for_family(const char *family) {
    if (ulab_streq(family, "node_ops") || ulab_streq(family, "software_update")) return "node";
    if (ulab_streq(family, "package")) return "package";
    if (ulab_streq(family, "subscriber")) return "subscriber";
    if (ulab_streq(family, "site_ops")) return "site";
    if (ulab_streq(family, "dashboard")) return "dashboard";
    return "sim";
}

static void write_matrix_event(FILE *f, const char *family, const char *flow,
                               const char *runtime, int wip) {
    if (wip) {
        fprintf(f, "      - type: check\n");
        return;
    }

    if (ulab_streq(runtime, "node_restart") ||
        (ulab_streq(family, "node_ops") && ulab_streq(flow, "restart")) ||
        (ulab_streq(family, "site_ops") && ulab_streq(flow, "restart_site"))) {
        fprintf(f,
                "      - type: restart_nodes\n"
                "        type_selector: tower\n"
                "        count_per_network: 1\n");
    } else if (ulab_streq(family, "usage") ||
               ulab_streq(flow, "ue_traffic") ||
               ulab_streq(flow, "traffic")) {
        fprintf(f,
                "      - type: traffic\n"
                "        ues: all\n"
                "        amount_mb: 1\n");
    } else if (ulab_streq(flow, "profile_traffic")) {
        fprintf(f,
                "      - type: traffic_by_profile\n"
                "        profile: mixed\n");
    } else if (ulab_streq(flow, "add_to_sim") || ulab_streq(flow, "topup")) {
        fprintf(f,
                "      - type: add_package_to_sim\n"
                "        ues: all\n"
                "        package: pkg_002\n");
    } else if (ulab_streq(flow, "remove_from_sim")) {
        fprintf(f,
                "      - type: remove_package_from_sim\n"
                "        ues: all\n");
    } else if (ulab_streq(flow, "sim_suspend") || ulab_streq(flow, "sim_deactivate")) {
        fprintf(f,
                "      - type: set_sim_status\n"
                "        ues: all\n"
                "        status: inactive\n");
    } else {
        fprintf(f, "      - type: check\n");
    }
}

static void write_matrix_checks(FILE *f, const char *family,
                                const char *verification,
                                const char *runtime, int wip) {
    const char *target;

    target = (ulab_streq(family, "node_ops") ||
              ulab_streq(family, "software_update")) ? "nodes" : "sims";

    fprintf(f,
            "    checks:\n"
            "      - type: backend_count\n"
            "        target: %s\n"
            "        expected: from_world\n", target);

    if (!wip && ulab_streq(verification, "read_model")) {
        fprintf(f,
                "      - type: list_contains\n"
                "        view: sims\n"
                "        ref: ue-000001\n");
    }
    if (!wip && ulab_streq(verification, "dashboard")) {
        fprintf(f,
                "      - type: dashboard_loads\n"
                "        view: network_overview\n");
    }
    if (!wip && ulab_streq(verification, "runtime") &&
        (ulab_streq(runtime, "normal") || ulab_streq(runtime, "node_restart"))) {
        fprintf(f,
                "      - type: traffic_allowed\n"
                "        ues: all\n"
                "        amount_mb: 1\n");
    }
}

static void write_profile_section(FILE *f, const char *flow) {
    if (!ulab_streq(flow, "profile_traffic")) {
        return;
    }
    fprintf(f,
            "profiles:\n"
            "  - name: mixed\n"
            "    buckets:\n"
            "      - name: light\n"
            "        percent: 70\n"
            "        amount_mb: 1\n"
            "      - name: heavy\n"
            "        percent: 30\n"
            "        amount_mb: 3\n\n");
}

static int write_matrix_scenario(const gen_opts_t *opts,
                                 const matrix_family_t *family,
                                 const matrix_item_t *topology,
                                 const matrix_item_t *scale,
                                 const char *flow,
                                 const char *runtime,
                                 const char *failure,
                                 const char *software,
                                 const char *verification,
                                 FILE *index,
                                 size_t *count_out) {
    char dir[ULAB_MAX_PATH];
    char path[ULAB_MAX_PATH];
    char name[ULAB_MAX_NAME];
    uint32_t seed;
    uint32_t networks;
    uint32_t sites;
    uint32_t ues;
    uint32_t packages;
    int wip;
    FILE *f;

    if (build_path3(dir, sizeof(dir), opts->out_dir, family->name, "")) {
        return ULAB_ERR;
    }
    if (dir[0] != '\0' && dir[strlen(dir) - 1] == '/') {
        dir[strlen(dir) - 1] = '\0';
    }
    if (ulab_mkdir_p(dir)) {
        return ULAB_ERR;
    }

    build_matrix_name(name, sizeof(name), family->name, flow,
                      topology->name, scale->name, runtime, failure,
                      software, verification);
    if (build_yaml_path(path, sizeof(path), dir, name)) {
        return ULAB_ERR;
    }

    networks = topology->networks ? topology->networks : 1;
    sites = topology->sites_per_network ? topology->sites_per_network : 1;
    ues = scale->ues_per_site ? scale->ues_per_site : 1;
    packages = scale->packages ? scale->packages : 1;
    wip = should_mark_wip(family->name, flow, runtime, failure,
                          software, verification);

    seed = ulab_hash32(family->name, 41000);
    seed = ulab_hash32(flow, seed);
    seed = ulab_hash32(topology->name, seed);
    seed = ulab_hash32(scale->name, seed);
    seed = ulab_hash32(runtime, seed);
    seed = ulab_hash32(failure, seed);
    seed = ulab_hash32(software, seed);
    seed = ulab_hash32(verification, seed);

    f = fopen(path, "w");
    if (f == NULL) {
        fprintf(stderr, "unable to write: %s\n", path);
        return ULAB_ERR;
    }

    fprintf(f,
            "version: 1\n"
            "name: %s\n"
            "seed: %u\n"
            "suite: generated\n"
            "priority: %s\n"
            "tags: [generated, matrix, %s, %s, %s, %s, %s, %s, %s, %s%s]\n"
            "status: %s\n"
            "generated: true\n"
            "entity: %s\n"
            "action: %s\n\n",
            name, seed, family->priority[0] ? family->priority : "p2",
            family->name, flow, topology->name, scale->name, runtime,
            failure, software, verification, wip ? ", future" : "",
            wip ? "wip" : "active",
            matrix_entity_for_family(family->name), flow);

    fprintf(f,
            "# family: %s\n"
            "# flow: %s\n"
            "# topology: %s\n"
            "# scale: %s\n"
            "# runtime_condition: %s\n"
            "# failure_condition: %s\n"
            "# software_condition: %s\n"
            "# verification: %s\n"
            "# generated_status: %s\n\n",
            family->name, flow, topology->name, scale->name, runtime,
            failure, software, verification, wip ? "wip" : "active");

    write_base_world(f, networks, sites, ues, packages);
    write_profile_section(f, flow);

    fprintf(f,
            "phases:\n"
            "  - name: matrix_action\n"
            "    events:\n");
    write_matrix_event(f, family->name, flow, runtime, wip);
    write_matrix_checks(f, family->name, verification, runtime, wip);

    fprintf(f,
            "\nfinal_checks:\n"
            "  - type: balance_non_negative\n"
            "    ues: all\n");

    fclose(f);

    fprintf(index,
            "  - file: %s/%s.yaml\n"
            "    family: %s\n"
            "    flow: %s\n"
            "    topology: %s\n"
            "    scale: %s\n"
            "    runtime: %s\n"
            "    failure: %s\n"
            "    software: %s\n"
            "    verification: %s\n"
            "    status: %s\n"
            "    priority: %s\n",
            family->name, name, family->name, flow, topology->name,
            scale->name, runtime, failure, software, verification,
            wip ? "wip" : "active",
            family->priority[0] ? family->priority : "p2");

    (*count_out)++;
    return ULAB_OK;
}

static int generate_family_matrix(const gen_opts_t *opts,
                                  const matrix_def_t *matrix,
                                  const matrix_family_t *family,
                                  FILE *index,
                                  size_t *count_out) {
    size_t fi;
    size_t ti;
    size_t si;
    size_t ri;
    size_t fai;
    size_t swi;
    size_t vi;

    for (fi = 0; fi < family->flow_count; fi++) {
        for (ti = 0; ti < family->topology_count; ti++) {
            const matrix_item_t *topology;

            topology = profile_item(matrix, "topology", family->topologies[ti]);
            if (topology == NULL) return ULAB_ERR;

            for (si = 0; si < family->scale_count; si++) {
                const matrix_item_t *scale;

                scale = profile_item(matrix, "scale", family->scales[si]);
                if (scale == NULL) return ULAB_ERR;

                for (ri = 0; ri < family->runtime_count; ri++) {
                    for (fai = 0; fai < family->failure_count; fai++) {
                        for (swi = 0; swi < family->software_count ||
                             (swi == 0 && family->software_count == 0); swi++) {
                            const char *software;

                            software = family->software_count ?
                                       family->software[swi] : "no_update";
                            for (vi = 0; vi < family->verification_count; vi++) {
                                if (write_matrix_scenario(opts, family,
                                                          topology, scale,
                                                          family->flows[fi],
                                                          family->runtime[ri],
                                                          family->failures[fai],
                                                          software,
                                                          family->verification[vi],
                                                          index, count_out)) {
                                    return ULAB_ERR;
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    return ULAB_OK;
}

static void write_matrix_index(FILE *index, const matrix_def_t *matrix) {
    size_t i;

    fprintf(index, "matrix:\n");
    fprintf(index, "  profiles:\n");
    for (i = 0; i < matrix->profile_count; i++) {
        fprintf(index, "    - name: %s\n", matrix->profiles[i].name);
        fprintf(index, "      kind: %s\n", matrix->profiles[i].kind);
        fprintf(index, "      count: %zu\n", matrix->profiles[i].item_count);
    }
    fprintf(index, "  families:\n");
    for (i = 0; i < matrix->family_count; i++) {
        fprintf(index, "    - name: %s\n", matrix->families[i].name);
        fprintf(index, "      priority: %s\n", matrix->families[i].priority);
        fprintf(index, "      flows: %zu\n", matrix->families[i].flow_count);
    }
}

static int generate_model(const gen_opts_t *opts, const char *entity,
                          FILE *index) {
    model_def_t model;
    size_t i;

    if (load_model(opts, entity, &model)) {
        return ULAB_ERR;
    }

    for (i = 0; i < model.action_count; i++) {
        if (generate_action(opts, &model, &model.actions[i], index)) {
            return ULAB_ERR;
        }
    }

    return ULAB_OK;
}

int generator_run(int argc, char **argv) {
    gen_opts_t opts;
    matrix_def_t matrix;
    char index_path[ULAB_MAX_PATH];
    FILE *index;
    size_t i;
    size_t matrix_count;

    if (parse_opts(argc, argv, &opts) != ULAB_OK) {
        usage();
        return ULAB_EUSAGE;
    }

    if (load_matrix_models(&opts, &matrix)) {
        return ULAB_ERR;
    }
    printf("loaded matrix profiles=%zu families=%zu\n", matrix.profile_count,
           matrix.family_count);

    if (ulab_mkdir_p(opts.out_dir)) {
        fprintf(stderr, "unable to create output dir: %s\n", opts.out_dir);
        return ULAB_ERR;
    }

    index_path[0] = '\0';
    if (append_str(index_path, sizeof(index_path), opts.out_dir)) return ULAB_ERR;
    if (append_str(index_path, sizeof(index_path), "/index.yaml")) return ULAB_ERR;

    index = fopen(index_path, "w");
    if (index == NULL) {
        fprintf(stderr, "unable to write: %s\n", index_path);
        return ULAB_ERR;
    }

    write_matrix_index(index, &matrix);
    fprintf(index, "generated_scenarios:\n");

    if (!ulab_streq(opts.model, "all")) {
        if (generate_model(&opts, opts.model, index)) {
            fclose(index);
            return ULAB_ERR;
        }
    } else {
        for (i = 0; i < sizeof(default_entities) / sizeof(default_entities[0]); i++) {
            if (generate_model(&opts, default_entities[i], index)) {
                fclose(index);
                return ULAB_ERR;
            }
        }
    }

    matrix_count = 0;
    for (i = 0; i < matrix.family_count; i++) {
        if (generate_family_matrix(&opts, &matrix, &matrix.families[i],
                                   index, &matrix_count)) {
            fclose(index);
            return ULAB_ERR;
        }
    }

    fclose(index);
    printf("generated matrix scenarios=%zu\n", matrix_count);
    printf("generated index %s\n", index_path);
    return ULAB_OK;
}
