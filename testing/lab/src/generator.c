/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#include "generator.h"
#include "scenario.h"
#include "util.h"

#define GEN_MAX_ACTIONS 64

typedef struct {
    char model[ULAB_MAX_REF];
    char out_dir[ULAB_MAX_PATH];
    char models_dir[ULAB_MAX_PATH];
} gen_opts_t;

typedef struct {
    char entity[ULAB_MAX_REF];
    char actions[GEN_MAX_ACTIONS][ULAB_MAX_REF];
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
            /* accepted for compatibility; phase-6 stores suite/priority in YAML */
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
            /* accepted for compatibility; generator is model-driven now */
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

static int add_action(model_def_t *m, const char *action) {
    if (action == NULL || action[0] == '\0') {
        return ULAB_OK;
    }
    if (m->action_count >= GEN_MAX_ACTIONS) {
        return ULAB_ERR;
    }
    return ulab_copy(m->actions[m->action_count++], ULAB_MAX_REF, action);
}

static int parse_actions_inline(model_def_t *m, char *val) {
    char *p;
    char *tok;

    p = ulab_trim(val);
    if (*p == '[') p++;
    tok = strtok(p, ",]");
    while (tok != NULL) {
        char *a = ulab_trim(tok);
        if (add_action(m, a)) {
            return ULAB_ERR;
        }
        tok = strtok(NULL, ",]");
    }
    return ULAB_OK;
}

static int load_model(const gen_opts_t *opts, const char *entity,
                      model_def_t *model) {
    char path[ULAB_MAX_PATH];
    char line[ULAB_MAX_LINE];
    FILE *fp;
    int in_actions;

    memset(model, 0, sizeof(*model));
    in_actions = 0;

    if (snprintf(path, sizeof(path), "%s/%s.yaml", opts->models_dir,
                 entity) >= (int)sizeof(path)) {
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
        if (ulab_starts(p, "- ") && in_actions) {
            if (add_action(model, ulab_trim(p + 2))) {
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
        } else if (ulab_streq(key, "actions")) {
            in_actions = 1;
            if (val[0] != '\0') {
                if (parse_actions_inline(model, val)) {
                    fclose(fp);
                    return ULAB_ERR;
                }
            }
        } else {
            in_actions = 0;
        }
    }

    fclose(fp);

    if (!ulab_streq(model->entity, entity)) {
        fprintf(stderr, "model entity mismatch: %s\n", path);
        return ULAB_ERR;
    }
    if (model->action_count == 0) {
        add_action(model, "read");
    }

    return ULAB_OK;
}

static const char *event_for_action(const char *entity, const char *action) {
    if (ulab_streq(entity, "node") && ulab_streq(action, "restart")) {
        return "restart_nodes";
    }
    if ((ulab_streq(entity, "package") && ulab_streq(action, "add_to_sim")) ||
        (ulab_streq(entity, "sim") && ulab_streq(action, "add_package"))) {
        return "add_package_to_sim";
    }
    if ((ulab_streq(entity, "package") && ulab_streq(action, "remove_from_sim")) ||
        (ulab_streq(entity, "sim") && ulab_streq(action, "remove_package"))) {
        return "remove_package_from_sim";
    }
    if (ulab_streq(entity, "sim") &&
        (ulab_streq(action, "set_active") || ulab_streq(action, "set_inactive") ||
         ulab_streq(action, "release") || ulab_streq(action, "deactivate"))) {
        return "set_sim_status";
    }
    return "check";
}

static const char *status_for_action(const char *action) {
    if (ulab_streq(action, "set_inactive") || ulab_streq(action, "release") ||
        ulab_streq(action, "deactivate") || ulab_streq(action, "archive") ||
        ulab_streq(action, "suspend") || ulab_streq(action, "retire")) {
        return "inactive";
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

static int write_success_scenario(const gen_opts_t *opts,
                                  const char *entity,
                                  const char *action,
                                  FILE *index) {
    char dir[ULAB_MAX_PATH];
    char path[ULAB_MAX_PATH];
    char name[ULAB_MAX_NAME];
    const char *event;
    uint32_t seed;
    FILE *f;

    event = event_for_action(entity, action);
    seed = ulab_hash32(entity, 9000);
    seed = ulab_hash32(action, seed);

    if (path_join2(dir, sizeof(dir), opts->out_dir, entity)) return ULAB_ERR;
    if (ulab_mkdir_p(dir)) return ULAB_ERR;
    snprintf(name, sizeof(name), "%s-%s-success", entity, action);
    snprintf(path, sizeof(path), "%s/%s.yaml", dir, name);

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
            "priority: p1\n"
            "tags: [generated, %s, %s, success]\n"
            "status: active\n"
            "generated: true\n"
            "entity: %s\n"
            "action: %s\n\n",
            name, seed, entity, action, entity, action);
    write_common_world(f);

    fprintf(f, "phases:\n  - name: action\n    events:\n");
    if (ulab_streq(event, "add_package_to_sim")) {
        fprintf(f, "      - type: add_package_to_sim\n        ues: all\n        package: weekly_10gb\n");
    } else if (ulab_streq(event, "remove_package_from_sim")) {
        fprintf(f, "      - type: remove_package_from_sim\n        ues: all\n");
    } else if (ulab_streq(event, "set_sim_status")) {
        fprintf(f, "      - type: set_sim_status\n        ues: all\n        status: %s\n", status_for_action(action));
    } else if (ulab_streq(event, "restart_nodes")) {
        fprintf(f, "      - type: restart_nodes\n        type_selector: tower\n        count_per_network: 1\n");
    } else {
        fprintf(f, "      - type: check\n");
    }

    fprintf(f,
            "    checks:\n"
            "      - type: backend_count\n"
            "        target: sims\n"
            "        expected: from_world\n"
            "      - type: list_contains\n"
            "        view: sims\n"
            "        ref: ue-000001\n");
    if (ulab_streq(entity, "sim")) {
        fprintf(f,
                "      - type: status_equals\n"
                "        entity: sim\n"
                "        ref: ue-000001\n"
                "        status: %s\n", status_for_action(action));
    }
    fprintf(f,
            "\nfinal_checks:\n"
            "  - type: balance_non_negative\n"
            "    ues: all\n");

    fclose(f);
    if (index != NULL) {
        fprintf(index, "  - file: %s/%s.yaml\n    entity: %s\n    action: %s\n    case: success\n    suite: generated\n    priority: p1\n",
                entity, name, entity, action);
    }
    printf("generated %s\n", path);
    return ULAB_OK;
}

static int write_blocked_scenario(const gen_opts_t *opts,
                                  const char *entity,
                                  const char *action,
                                  FILE *index) {
    char dir[ULAB_MAX_PATH];
    char path[ULAB_MAX_PATH];
    char name[ULAB_MAX_NAME];
    uint32_t seed;
    FILE *f;

    seed = ulab_hash32(entity, 10000);
    seed = ulab_hash32(action, seed);

    if (path_join2(dir, sizeof(dir), opts->out_dir, entity)) return ULAB_ERR;
    if (ulab_mkdir_p(dir)) return ULAB_ERR;
    snprintf(name, sizeof(name), "%s-%s-blocked", entity, action);
    snprintf(path, sizeof(path), "%s/%s.yaml", dir, name);

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
            "priority: p2\n"
            "tags: [generated, %s, %s, blocked, negative]\n"
            "status: active\n"
            "generated: true\n"
            "entity: %s\n"
            "action: %s\n\n",
            name, seed, entity, action, entity, action);
    write_common_world(f);

    fprintf(f,
            "phases:\n"
            "  - name: blocked_action\n"
            "    events:\n"
            "      - type: set_sim_status\n"
            "        ues: all\n"
            "        expect:\n"
            "          result: failure\n"
            "          error_contains: \"missing status\"\n"
            "    checks:\n"
            "      - type: backend_count\n"
            "        target: sims\n"
            "        expected: from_world\n\n"
            "final_checks:\n"
            "  - type: balance_non_negative\n"
            "    ues: all\n");

    fclose(f);
    if (index != NULL) {
        fprintf(index, "  - file: %s/%s.yaml\n    entity: %s\n    action: %s\n    case: blocked\n    suite: generated\n    priority: p2\n",
                entity, name, entity, action);
    }
    printf("generated %s\n", path);
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
        if (write_success_scenario(opts, entity, model.actions[i], index)) {
            return ULAB_ERR;
        }
        if (write_blocked_scenario(opts, entity, model.actions[i], index)) {
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
