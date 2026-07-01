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

#include "runner.h"
#include "check.h"
#include "event.h"
#include "log.h"
#include "util.h"
#include "generator.h"

static void usage(void) {
    printf("ukama-lab %s\n", ULAB_VERSION);
    printf("usage:\n");
    printf("  ukama-lab validate <scenario.yaml> [options]\n");
    printf("  ukama-lab generate --model <name|all> --mode <name|all> [options]\n");
    printf("  ukama-lab list-checks\n");
    printf("  ukama-lab list-events\n");
    printf("  ukama-lab version\n");
    printf("options:\n");
    printf("  --repo <dir>          ukama repo root path (MUST)\n");
    printf("  --seed <n>            override scenario seed\n");
    printf("  --bff <url>           BFF GraphQL endpoint\n");
    printf("  --out <dir>           output directory\n");
    printf("  --run-id <id>         fixed run id; existing created.json is cleaned first\n");
    printf("  --scripts <dir>       runtime script directory\n");
    printf("  --sim-type <type>     SIM pool type; default: ukama_data\n");
    printf("  --warehouse-url <url> warehouse API URL for run SIM provisioning\n");
    printf("  --factory-url <url>   sim factory API URL for run SIM export\n");
    printf("  --asr-url <url>       optional ukama-agent ASR API URL for post-allocation check\n");
    printf("generate options:\n");
    printf("  --model <name|all>    org/network/site/node/sim/subscriber/package\n");
    printf("  --mode <name|all>     smoke/transition/negative/pairwise/full\n");
    printf("  --models <dir>        model directory; default: models\n");
    printf("  --templates <dir>     template directory; default: templates/generated\n");
    printf("  --quiet               summary only\n");
    printf("  --verbose             debug logs\n");
}

static void opts_init(runner_opts_t *o) {
    memset(o, 0, sizeof(*o));
    ulab_copy(o->bff_url, sizeof(o->bff_url),
              ulab_getenv_default("UKAMA_LAB_BFF",
                                  "http://localhost:8080/graphql"));
    ulab_copy(o->out_dir, sizeof(o->out_dir),
              ulab_getenv_default("UKAMA_LAB_OUT", "runs"));
    ulab_copy(o->script_dir, sizeof(o->script_dir),
              ulab_getenv_default("UKAMA_LAB_SCRIPTS", "scripts"));
    ulab_copy(o->sim_type, sizeof(o->sim_type),
              ulab_getenv_default("UKAMA_LAB_SIM_TYPE", "ukama_data"));
    ulab_copy(o->warehouse_url, sizeof(o->warehouse_url),
              ulab_getenv_default("UKAMA_LAB_WAREHOUSE_URL",
                                  "http://warehouse-ukama.udev.ukama.com"));
    ulab_copy(o->factory_url, sizeof(o->factory_url),
              ulab_getenv_default("UKAMA_LAB_FACTORY_URL",
                                  "http://factory-ukama.udev.ukama.com"));
    ulab_copy(o->asr_url, sizeof(o->asr_url),
              ulab_getenv_default("UKAMA_LAB_ASR_URL", ""));
}

static int parse_opts(int argc, char **argv, int start, runner_opts_t *o) {
    int i;

    for (i = start; i < argc; i++) {
        if (ulab_streq(argv[i], "--seed") && i + 1 < argc) {
            if (ulab_parse_u32(argv[++i], &o->seed_override)) {
                return ULAB_EUSAGE;
            }
            o->has_seed_override = 1;
        } else if (ulab_streq(argv[i], "--bff") && i + 1 < argc) {
            ulab_copy(o->bff_url, sizeof(o->bff_url), argv[++i]);
        } else if (ulab_streq(argv[i], "--out") && i + 1 < argc) {
            ulab_copy(o->out_dir, sizeof(o->out_dir), argv[++i]);
        } else if (ulab_streq(argv[i], "--run-id") && i + 1 < argc) {
            ulab_copy(o->run_id, sizeof(o->run_id), argv[++i]);
        } else if (ulab_streq(argv[i], "--scripts") && i + 1 < argc) {
            ulab_copy(o->script_dir, sizeof(o->script_dir), argv[++i]);
        } else if (ulab_streq(argv[i], "--repo") && i + 1 < argc) {
            ulab_copy(o->repo, sizeof(o->repo), argv[++i]);
        } else if (ulab_streq(argv[i], "--sim-type") && i + 1 < argc) {
            ulab_copy(o->sim_type, sizeof(o->sim_type), argv[++i]);
        } else if (ulab_streq(argv[i], "--warehouse-url") && i + 1 < argc) {
            ulab_copy(o->warehouse_url, sizeof(o->warehouse_url), argv[++i]);
        } else if (ulab_streq(argv[i], "--factory-url") && i + 1 < argc) {
            ulab_copy(o->factory_url, sizeof(o->factory_url), argv[++i]);
        } else if (ulab_streq(argv[i], "--asr-url") && i + 1 < argc) {
            ulab_copy(o->asr_url, sizeof(o->asr_url), argv[++i]);
        } else if (ulab_streq(argv[i], "--quiet")) {
            o->quiet = 1;
        } else if (ulab_streq(argv[i], "--verbose")) {
            o->verbose = 1;
        } else {
            fprintf(stderr, "unknown option: %s\n", argv[i]);
            return ULAB_EUSAGE;
        }
    }

    return ULAB_OK;
}

int main(int argc, char **argv) {

    runner_opts_t opts;
    int rc;

    if (argc < 2) {
        usage();
        return ULAB_EUSAGE;
    }

    if (ulab_streq(argv[1], "version")) {
        printf("ukama-lab %s scenario-v%d\n", ULAB_VERSION,
               ULAB_SCHEMA_VER);
        return ULAB_OK;
    }

    if (ulab_streq(argv[1], "generate")) {
        return generator_run(argc - 2, argv + 2);
    }

    if (ulab_streq(argv[1], "list-checks")) {
        check_list_supported();
        return ULAB_OK;
    }

    if (ulab_streq(argv[1], "list-events")) {
        event_list_supported();
        return ULAB_OK;
    }

    if (argc < 3) {
        usage();
        return ULAB_EUSAGE;
    }

    opts_init(&opts);
    ulab_copy(opts.scenario_path, sizeof(opts.scenario_path), argv[2]);
    rc = parse_opts(argc, argv, 3, &opts);
    if (rc != ULAB_OK) {
        usage();
        return rc;
    }

    /* repo path is must else we wont know how to build virtual node/ue */
    if (opts.repo[0] == '\0' || strstr(opts.repo, "ukama") == NULL) {
        printf("Missing --repo. Ukama repo root is MUST\n");
        usage();
        return ULAB_EUSAGE;
    }

    ulab_log_set_quiet(opts.quiet);
    ulab_log_set_verbose(opts.verbose);

    if (ulab_streq(argv[1], "validate")) {
        return runner_validate(&opts);
    }

    usage();
    return ULAB_EUSAGE;
}
