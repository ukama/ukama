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

#define GEN_MAX_TEXT 4096

typedef struct {
    char model[ULAB_MAX_REF];
    char mode[ULAB_MAX_REF];
    char out_dir[ULAB_MAX_PATH];
    char models_dir[ULAB_MAX_PATH];
    char templates_dir[ULAB_MAX_PATH];
} gen_opts_t;

typedef struct {
    const char *entity;
    const char *mode;
    const char *templ;
    char name[ULAB_MAX_NAME];
    char seed[ULAB_MAX_REF];
    char priority[ULAB_MAX_REF];
    char status[ULAB_MAX_REF];
    char phase[ULAB_MAX_NAME];
    char events[GEN_MAX_TEXT];
    char checks[GEN_MAX_TEXT];
    char final_checks[GEN_MAX_TEXT];
} gen_case_t;

static const char *entities[] = {
    "org", "network", "site", "node", "sim", "subscriber", "package"
};

static const char *modes[] = {
    "smoke", "transition", "negative", "pairwise", "full"
};

static const char *templates_all[] = {
    "state-transition",
    "blocked-transition",
    "lifecycle-cleanup",
    "permission-check",
    "retry-idempotency",
    "partial-failure",
    "wrong-org-network",
    "empty-state",
    "boundary-values",
    "backend-failure",
    "runtime-effect",
    "read-model-check"
};

static void usage(void) {
    printf("usage:\n");
    printf("  ukama-lab generate --model <name|all> --mode <name|all> [options]\n");
    printf("options:\n");
    printf("  --out <dir>        output directory; default: scenarios/generated\n");
    printf("  --models <dir>     model directory; default: models\n");
    printf("  --templates <dir>  template directory; default: templates/generated\n");
}

static int is_in(const char *v, const char **arr, size_t n) {
    size_t i;

    if (ulab_streq(v, "all")) {
        return 1;
    }
    for (i = 0; i < n; i++) {
        if (ulab_streq(v, arr[i])) {
            return 1;
        }
    }
    return 0;
}

static int parse_opts(int argc, char **argv, gen_opts_t *opts) {
    int i;

    memset(opts, 0, sizeof(*opts));
    ulab_copy(opts->model, sizeof(opts->model), "all");
    ulab_copy(opts->mode, sizeof(opts->mode), "smoke");
    ulab_copy(opts->out_dir, sizeof(opts->out_dir), "scenarios/generated");
    ulab_copy(opts->models_dir, sizeof(opts->models_dir), "models");
    ulab_copy(opts->templates_dir, sizeof(opts->templates_dir),
              "templates/generated");

    for (i = 0; i < argc; i++) {
        if (ulab_streq(argv[i], "--model") && i + 1 < argc) {
            if (ulab_copy(opts->model, sizeof(opts->model), argv[++i])) {
                return ULAB_EUSAGE;
            }
        } else if (ulab_streq(argv[i], "--mode") && i + 1 < argc) {
            if (ulab_copy(opts->mode, sizeof(opts->mode), argv[++i])) {
                return ULAB_EUSAGE;
            }
        } else if (ulab_streq(argv[i], "--out") && i + 1 < argc) {
            if (ulab_copy(opts->out_dir, sizeof(opts->out_dir), argv[++i])) {
                return ULAB_EUSAGE;
            }
        } else if (ulab_streq(argv[i], "--models") && i + 1 < argc) {
            if (ulab_copy(opts->models_dir, sizeof(opts->models_dir),
                          argv[++i])) {
                return ULAB_EUSAGE;
            }
        } else if (ulab_streq(argv[i], "--templates") && i + 1 < argc) {
            if (ulab_copy(opts->templates_dir, sizeof(opts->templates_dir),
                          argv[++i])) {
                return ULAB_EUSAGE;
            }
        } else if (ulab_streq(argv[i], "--help")) {
            usage();
            return ULAB_EUSAGE;
        } else {
            fprintf(stderr, "unknown generate option: %s\n", argv[i]);
            return ULAB_EUSAGE;
        }
    }

    if (!is_in(opts->model, entities, sizeof(entities) / sizeof(entities[0]))) {
        fprintf(stderr, "unknown model: %s\n", opts->model);
        return ULAB_EUSAGE;
    }
    if (!is_in(opts->mode, modes, sizeof(modes) / sizeof(modes[0]))) {
        fprintf(stderr, "unknown mode: %s\n", opts->mode);
        return ULAB_EUSAGE;
    }

    return ULAB_OK;
}

static int path_join3(char *out, size_t n, const char *a, const char *b,
                      const char *c) {
    int rc;

    rc = snprintf(out, n, "%s/%s/%s", a, b, c);
    return rc < 0 || (size_t)rc >= n ? ULAB_ERR : ULAB_OK;
}

static int path_join2(char *out, size_t n, const char *a, const char *b) {
    int rc;

    rc = snprintf(out, n, "%s/%s", a, b);
    return rc < 0 || (size_t)rc >= n ? ULAB_ERR : ULAB_OK;
}

static int load_model(const gen_opts_t *opts, const char *entity) {
    char path[ULAB_MAX_PATH];
    char want[ULAB_MAX_LINE];
    char line[ULAB_MAX_LINE];
    FILE *fp;

    if (snprintf(want, sizeof(want), "entity: %s", entity) >=
        (int)sizeof(want)) {
        return ULAB_ERR;
    }
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
        char *p = ulab_trim(line);
        if (ulab_starts(p, want)) {
            fclose(fp);
            return ULAB_OK;
        }
    }
    fclose(fp);
    fprintf(stderr, "model entity mismatch: %s\n", path);
    return ULAB_ERR;
}

static int template_for_mode(const char *mode, const char *templ) {
    if (ulab_streq(mode, "full")) return 1;
    if (ulab_streq(mode, "smoke")) {
        return ulab_streq(templ, "read-model-check") ||
               ulab_streq(templ, "runtime-effect");
    }
    if (ulab_streq(mode, "transition")) {
        return ulab_streq(templ, "state-transition") ||
               ulab_streq(templ, "lifecycle-cleanup");
    }
    if (ulab_streq(mode, "negative")) {
        return ulab_streq(templ, "blocked-transition") ||
               ulab_streq(templ, "partial-failure") ||
               ulab_streq(templ, "wrong-org-network") ||
               ulab_streq(templ, "backend-failure");
    }
    if (ulab_streq(mode, "pairwise")) {
        return ulab_streq(templ, "boundary-values") ||
               ulab_streq(templ, "read-model-check") ||
               ulab_streq(templ, "runtime-effect") ||
               ulab_streq(templ, "permission-check") ||
               ulab_streq(templ, "retry-idempotency") ||
               ulab_streq(templ, "empty-state");
    }
    return 0;
}

static void set_common(gen_case_t *c, const char *entity, const char *mode,
                       const char *templ) {
    uint32_t seed;

    memset(c, 0, sizeof(*c));
    c->entity = entity;
    c->mode = mode;
    c->templ = templ;
    seed = ulab_hash32(entity, 6000);
    seed = ulab_hash32(mode, seed);
    seed = ulab_hash32(templ, seed);
    snprintf(c->name, sizeof(c->name), "generated-%s-%s-%s", entity, mode,
             templ);
    snprintf(c->seed, sizeof(c->seed), "%u", 6000u + (seed % 200000u));
    snprintf(c->priority, sizeof(c->priority), "%s",
             ulab_streq(mode, "smoke") ? "p0" : "p2");
    snprintf(c->status, sizeof(c->status), "active");
    snprintf(c->phase, sizeof(c->phase), "%s", templ);
}

static void set_read_model(gen_case_t *c) {
    snprintf(c->events, sizeof(c->events),
             "      - type: check\n");
    snprintf(c->checks, sizeof(c->checks),
             "      - type: backend_count\n"
             "        target: networks\n"
             "        expected: from_world\n"
             "      - type: backend_count\n"
             "        target: sites\n"
             "        expected: from_world\n"
             "      - type: backend_count\n"
             "        target: nodes\n"
             "        expected: from_world\n"
             "      - type: backend_count\n"
             "        target: sims\n"
             "        expected: from_world\n");
    snprintf(c->final_checks, sizeof(c->final_checks),
             "  - type: dashboard_loads\n"
             "    networks: all\n");
}

static void set_runtime(gen_case_t *c) {
    snprintf(c->events, sizeof(c->events),
             "      - type: traffic\n"
             "        ues: all\n"
             "        amount_mb: 1\n");
    snprintf(c->checks, sizeof(c->checks),
             "      - type: traffic_allowed\n"
             "        ues: all\n"
             "        amount_mb: 1\n");
    snprintf(c->final_checks, sizeof(c->final_checks),
             "  - type: balance_non_negative\n"
             "    ues: all\n");
}

static void set_state_transition(gen_case_t *c) {
    if (ulab_streq(c->entity, "sim") || ulab_streq(c->entity, "subscriber")) {
        snprintf(c->events, sizeof(c->events),
                 "      - type: set_sim_status\n"
                 "        ues: all\n"
                 "        status: inactive\n"
                 "      - type: set_sim_status\n"
                 "        ues: all\n"
                 "        status: active\n");
        snprintf(c->checks, sizeof(c->checks),
                 "      - type: backend_count\n"
                 "        target: sims\n"
                 "        expected: from_world\n");
    } else if (ulab_streq(c->entity, "package")) {
        snprintf(c->events, sizeof(c->events),
                 "      - type: remove_package_from_sim\n"
                 "        ues: all\n"
                 "      - type: add_package_to_sim\n"
                 "        ues: all\n"
                 "        package: daily_1gb\n");
        snprintf(c->checks, sizeof(c->checks),
                 "      - type: backend_count\n"
                 "        target: sims\n"
                 "        expected: from_world\n");
    } else if (ulab_streq(c->entity, "node")) {
        snprintf(c->events, sizeof(c->events),
                 "      - type: wait_nodes_ready\n"
                 "        nodes: all\n");
        snprintf(c->checks, sizeof(c->checks),
                 "      - type: node_ready\n"
                 "        nodes: all\n");
    } else {
        set_read_model(c);
    }
    if (c->final_checks[0] == '\0') {
        snprintf(c->final_checks, sizeof(c->final_checks),
                 "  - type: balance_non_negative\n"
                 "    ues: all\n");
    }
}

static void set_blocked_transition(gen_case_t *c) {
    snprintf(c->events, sizeof(c->events),
             "      - type: create_ues\n"
             "        count_per_site: 1\n"
             "        expect:\n"
             "          result: failure\n"
             "          error_contains: \"not enabled\"\n");
    snprintf(c->checks, sizeof(c->checks),
             "      - type: backend_count\n"
             "        target: sims\n"
             "        expected: from_world\n");
    snprintf(c->final_checks, sizeof(c->final_checks),
             "  - type: balance_non_negative\n"
             "    ues: all\n");
}

static void set_lifecycle_cleanup(gen_case_t *c) {
    snprintf(c->events, sizeof(c->events),
             "      - type: remove_package_from_sim\n"
             "        ues: all\n"
             "      - type: add_package_to_sim\n"
             "        ues: all\n"
             "        package: daily_1gb\n");
    snprintf(c->checks, sizeof(c->checks),
             "      - type: backend_count\n"
             "        target: sims\n"
             "        expected: from_world\n");
    snprintf(c->final_checks, sizeof(c->final_checks),
             "  - type: balance_non_negative\n"
             "    ues: all\n");
}

static void set_add_missing_package(gen_case_t *c, const char *contains) {
    snprintf(c->events, sizeof(c->events),
             "      - type: add_package_to_sim\n"
             "        ues: all\n"
             "        expect:\n"
             "          result: failure\n"
             "          error_contains: \"%s\"\n",
             contains);
    snprintf(c->checks, sizeof(c->checks),
             "      - type: backend_count\n"
             "        target: sims\n"
             "        expected: from_world\n");
    snprintf(c->final_checks, sizeof(c->final_checks),
             "  - type: balance_non_negative\n"
             "    ues: all\n");
}

static void set_add_unknown_package(gen_case_t *c) {
    snprintf(c->events, sizeof(c->events),
             "      - type: add_package_to_sim\n"
             "        ues: all\n"
             "        package: wrong_network_package\n"
             "        expect:\n"
             "          result: failure\n"
             "          error_contains: \"unknown package\"\n");
    snprintf(c->checks, sizeof(c->checks),
             "      - type: backend_count\n"
             "        target: sims\n"
             "        expected: from_world\n");
    snprintf(c->final_checks, sizeof(c->final_checks),
             "  - type: balance_non_negative\n"
             "    ues: all\n");
}

static void set_backend_failure(gen_case_t *c) {
    snprintf(c->events, sizeof(c->events),
             "      - type: set_sim_status\n"
             "        ues: all\n"
             "        expect:\n"
             "          result: failure\n"
             "          error_contains: \"missing status\"\n");
    snprintf(c->checks, sizeof(c->checks),
             "      - type: backend_count\n"
             "        target: sims\n"
             "        expected: from_world\n");
    snprintf(c->final_checks, sizeof(c->final_checks),
             "  - type: balance_non_negative\n"
             "    ues: all\n");
}

static void build_case(gen_case_t *c, const char *entity, const char *mode,
                       const char *templ) {
    set_common(c, entity, mode, templ);

    if (ulab_streq(templ, "state-transition")) set_state_transition(c);
    else if (ulab_streq(templ, "blocked-transition")) set_blocked_transition(c);
    else if (ulab_streq(templ, "lifecycle-cleanup")) set_lifecycle_cleanup(c);
    else if (ulab_streq(templ, "permission-check")) set_blocked_transition(c);
    else if (ulab_streq(templ, "retry-idempotency")) set_lifecycle_cleanup(c);
    else if (ulab_streq(templ, "partial-failure")) {
        set_add_missing_package(c, "missing package");
    } else if (ulab_streq(templ, "wrong-org-network")) set_add_unknown_package(c);
    else if (ulab_streq(templ, "empty-state")) set_read_model(c);
    else if (ulab_streq(templ, "boundary-values")) set_runtime(c);
    else if (ulab_streq(templ, "backend-failure")) set_backend_failure(c);
    else if (ulab_streq(templ, "runtime-effect")) set_runtime(c);
    else set_read_model(c);
}

static const char *token_value(const gen_case_t *c, const char *tok,
                               size_t len) {
    if (len == 4 && strncmp(tok, "NAME", len) == 0) return c->name;
    if (len == 4 && strncmp(tok, "SEED", len) == 0) return c->seed;
    if (len == 6 && strncmp(tok, "ENTITY", len) == 0) return c->entity;
    if (len == 4 && strncmp(tok, "MODE", len) == 0) return c->mode;
    if (len == 8 && strncmp(tok, "TEMPLATE", len) == 0) return c->templ;
    if (len == 8 && strncmp(tok, "PRIORITY", len) == 0) return c->priority;
    if (len == 6 && strncmp(tok, "STATUS", len) == 0) return c->status;
    if (len == 5 && strncmp(tok, "PHASE", len) == 0) return c->phase;
    if (len == 6 && strncmp(tok, "EVENTS", len) == 0) return c->events;
    if (len == 6 && strncmp(tok, "CHECKS", len) == 0) return c->checks;
    if (len == 12 && strncmp(tok, "FINAL_CHECKS", len) == 0) {
        return c->final_checks;
    }
    return NULL;
}

static void emit_line(FILE *out, const char *line, const gen_case_t *c) {
    const char *p = line;

    while (*p) {
        if (*p == '@') {
            const char *e = strchr(p + 1, '@');
            const char *v;
            if (e != NULL) {
                v = token_value(c, p + 1, (size_t)(e - p - 1));
                if (v != NULL) {
                    fputs(v, out);
                    p = e + 1;
                    continue;
                }
            }
        }
        fputc(*p++, out);
    }
}

static int render_case(const gen_opts_t *opts, const gen_case_t *c) {
    char tpath[ULAB_MAX_PATH];
    char odir[ULAB_MAX_PATH];
    char opath[ULAB_MAX_PATH];
    char file[ULAB_MAX_NAME];
    char tmpl_file[ULAB_MAX_NAME];
    FILE *in;
    FILE *out;
    char line[ULAB_MAX_LINE];

    if (snprintf(tmpl_file, sizeof(tmpl_file), "%s.yaml.tmpl", c->templ) >=
        (int)sizeof(tmpl_file)) {
        return ULAB_ERR;
    }
    if (path_join2(tpath, sizeof(tpath), opts->templates_dir, tmpl_file)) {
        return ULAB_ERR;
    }
    if (path_join3(odir, sizeof(odir), opts->out_dir, c->entity, c->mode)) {
        return ULAB_ERR;
    }
    if (ulab_mkdir_p(odir)) {
        return ULAB_ERR;
    }
    if (snprintf(file, sizeof(file), "%s.yaml", c->name) >=
        (int)sizeof(file)) {
        return ULAB_ERR;
    }
    if (path_join2(opath, sizeof(opath), odir, file)) {
        return ULAB_ERR;
    }

    in = fopen(tpath, "r");
    if (in == NULL) {
        fprintf(stderr, "missing template: %s\n", tpath);
        return ULAB_ERR;
    }
    out = fopen(opath, "w");
    if (out == NULL) {
        fclose(in);
        fprintf(stderr, "unable to write: %s\n", opath);
        return ULAB_ERR;
    }

    while (fgets(line, sizeof(line), in) != NULL) {
        emit_line(out, line, c);
    }

    fclose(out);
    fclose(in);
    printf("generated %s\n", opath);
    return ULAB_OK;
}

static int generate_one(const gen_opts_t *opts, const char *entity,
                        const char *mode) {
    size_t i;

    if (load_model(opts, entity)) {
        return ULAB_ERR;
    }
    for (i = 0; i < sizeof(templates_all) / sizeof(templates_all[0]); i++) {
        gen_case_t c;
        if (!template_for_mode(mode, templates_all[i])) {
            continue;
        }
        build_case(&c, entity, mode, templates_all[i]);
        if (render_case(opts, &c)) {
            return ULAB_ERR;
        }
    }
    return ULAB_OK;
}

int generator_run(int argc, char **argv) {
    gen_opts_t opts;
    size_t i;
    size_t j;

    if (parse_opts(argc, argv, &opts) != ULAB_OK) {
        usage();
        return ULAB_EUSAGE;
    }

    for (i = 0; i < sizeof(entities) / sizeof(entities[0]); i++) {
        if (!ulab_streq(opts.model, "all") && !ulab_streq(opts.model,
            entities[i])) {
            continue;
        }
        for (j = 0; j < sizeof(modes) / sizeof(modes[0]); j++) {
            if (!ulab_streq(opts.mode, "all") && !ulab_streq(opts.mode,
                modes[j])) {
                continue;
            }
            if (generate_one(&opts, entities[i], modes[j])) {
                return ULAB_ERR;
            }
        }
    }

    return ULAB_OK;
}
