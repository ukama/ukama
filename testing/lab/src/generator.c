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

#define GEN_MAX_ACTIONS 64
#define GEN_MAX_LIST    16

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

static const char *default_entities[] = {
    "sim", "package", "subscriber", "node"
};

static void usage(void) {
    printf("usage:\n");
    printf("  ukama-lab generate --model <name|all> [options]\n");
    printf("options:\n");
    printf("  --out <dir>        output directory; default: scenarios/generated\n");
    printf("  --models <dir>     model directory; default: models\n");
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

static int path_join2(char *out, size_t n, const char *a, const char *b) {
    int rc;

    rc = snprintf(out, n, "%s/%s", a, b);
    return rc < 0 || (size_t)rc >= n ? ULAB_ERR : ULAB_OK;
}

static int path_yaml(char *out, size_t n, const char *dir, const char *name) {
    int rc;

    rc = snprintf(out, n, "%s/%s.yaml", dir, name);
    return rc < 0 || (size_t)rc >= n ? ULAB_ERR : ULAB_OK;
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

    if (path_yaml(path, sizeof(path), opts->models_dir, entity)) {
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
            in_actions = 0;
            cur = NULL;
        } else if (ulab_streq(key, "states")) {
            if (parse_list(model->states, &model->state_count, val)) {
                fclose(fp);
                return ULAB_ERR;
            }
            in_actions = 0;
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

static void write_common_world(FILE *f) {
    fprintf(f,
            "provider:\n"
            "  type: virtual\n\n"
            "world:\n"
            "  networks: 1\n"
            "  sites_per_network: 1\n"
            "  nodes_per_site:\n"
            "    tower: 1\n"
            "    amplifier: 1\n"
            "    controller: 1\n"
            "  ues_per_site: 1\n\n"
            "packages:\n"
            "  - ref: daily_1gb\n"
            "    name: Daily 1GB\n"
            "    data_mb: 1024\n"
            "    duration_days: 1\n"
            "    amount: 1.00\n"
            "    assign_percent: 100\n"
            "  - ref: weekly_10gb\n"
            "    name: Weekly 10GB\n"
            "    data_mb: 10240\n"
            "    duration_days: 7\n"
            "    amount: 7.00\n"
            "    assign_percent: 0\n\n"
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
                a->package_ref[0] ? a->package_ref : "weekly_10gb");
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

static void write_checks(FILE *f, const model_def_t *m,
                         const action_rule_t *a, int blocked) {
    fprintf(f,
            "    checks:\n"
            "      - type: backend_count\n");
    if (ulab_streq(m->entity, "node")) {
        fprintf(f, "        target: nodes\n");
    } else {
        fprintf(f, "        target: sims\n");
    }
    fprintf(f, "        expected: from_world\n");

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

static int write_scenario(const gen_opts_t *opts, const model_def_t *m,
                          const action_rule_t *a, const char *case_name,
                          const char *from_state, int blocked, FILE *index) {
    char dir[ULAB_MAX_PATH];
    char path[ULAB_MAX_PATH];
    char name[ULAB_MAX_NAME];
    uint32_t seed;
    FILE *f;

    if (path_join2(dir, sizeof(dir), opts->out_dir, m->entity)) {
        return ULAB_ERR;
    }
    if (ulab_mkdir_p(dir)) {
        return ULAB_ERR;
    }

    snprintf(name, sizeof(name), "%s-%s-%s-%s", m->entity, a->name,
             case_name, from_state && from_state[0] ? from_state : "default");
    if (path_yaml(path, sizeof(path), dir, name)) {
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
            "tags: [generated, %s, %s, %s%s]\n"
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

    write_common_world(f);
    fprintf(f,
            "phases:\n"
            "  - name: action\n"
            "    events:\n");
    write_action_event(f, m, a, blocked);
    write_checks(f, m, a, blocked);

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

    printf("generated %s\n", path);
    return ULAB_OK;
}

static int generate_action(const gen_opts_t *opts, const model_def_t *m,
                           const action_rule_t *a, FILE *index) {
    size_t i;

    if (a->from_count == 0) {
        if (write_scenario(opts, m, a, "success", "default", 0, index)) {
            return ULAB_ERR;
        }
    } else {
        for (i = 0; i < a->from_count; i++) {
            if (write_scenario(opts, m, a, "success", a->from[i], 0, index)) {
                return ULAB_ERR;
            }
        }
    }

    for (i = 0; i < a->blocked_count; i++) {
        if (write_scenario(opts, m, a, "blocked", a->blocked_from[i], 1, index)) {
            return ULAB_ERR;
        }
    }

    for (i = 0; i < a->guard_count; i++) {
        if (write_scenario(opts, m, a, "guard", a->guards[i], 1, index)) {
            return ULAB_ERR;
        }
    }

    return ULAB_OK;
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
    char index_path[ULAB_MAX_PATH];
    FILE *index;
    size_t i;

    if (parse_opts(argc, argv, &opts) != ULAB_OK) {
        usage();
        return ULAB_EUSAGE;
    }

    if (ulab_mkdir_p(opts.out_dir)) {
        fprintf(stderr, "unable to create output dir: %s\n", opts.out_dir);
        return ULAB_ERR;
    }

    if (path_join2(index_path, sizeof(index_path), opts.out_dir, "index.yaml")) {
        return ULAB_ERR;
    }
    index = fopen(index_path, "w");
    if (index == NULL) {
        fprintf(stderr, "unable to write: %s\n", index_path);
        return ULAB_ERR;
    }
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

    fclose(index);
    printf("generated index %s\n", index_path);
    return ULAB_OK;
}
