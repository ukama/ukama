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

static void usage(void) {
    printf("ukama-lab %s\n", ULAB_VERSION);
    printf("usage:\n");
    printf("  ukama-lab validate <scenario.yaml> [options]\n");
    printf("  ukama-lab dry-run <scenario.yaml> [options]\n");
    printf("  ukama-lab list-checks\n");
    printf("  ukama-lab list-events\n");
    printf("  ukama-lab version\n");
    printf("options:\n");
    printf("  --repo <dir>     ukama repo root path (MUST)\n");
    printf("  --seed <n>       override scenario seed\n");
    printf("  --bff <url>      BFF GraphQL endpoint\n");
    printf("  --out <dir>      output directory\n");
    printf("  --scripts <dir>  runtime script directory\n");
    printf("  --setup-only     create BFF world and skip runtime\n");
    printf("  --subscriber     create package/subscriber/SIM only\n");
    printf("  --network-id <id> existing BFF network id for --subscriber\n");
    printf("  --sim-csv <file> upload SIM CSV before allocation\n");
    printf("  --sim-type <type> SIM pool type; default: ukama_data\n");
    printf("  --cleanup        force cleanup after setup-only/subscriber run\n");
    printf("  --keep           skip cleanup and keep runtime/resources\n");
    printf("  --keep-on-failure keep runtime/resources when validate fails\n");
    printf("  --print-world    dry-run: print generated world sample\n");
    printf("  --quiet          summary only\n");
    printf("  --verbose        debug logs\n");
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
    ulab_copy(o->subscriber_network_id,
              sizeof(o->subscriber_network_id),
              ulab_getenv_default("UKAMA_LAB_NETWORK_ID", ""));
    ulab_copy(o->sim_csv_path, sizeof(o->sim_csv_path),
              ulab_getenv_default("UKAMA_LAB_SIM_CSV", ""));
    ulab_copy(o->sim_type, sizeof(o->sim_type),
              ulab_getenv_default("UKAMA_LAB_SIM_TYPE", "ukama_data"));
    o->cleanup = 0;
    o->keep = 0;
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
        } else if (ulab_streq(argv[i], "--scripts") && i + 1 < argc) {
            ulab_copy(o->script_dir, sizeof(o->script_dir), argv[++i]);
        } else if (ulab_streq(argv[i], "--repo") && i + 1 < argc) {
            ulab_copy(o->repo, sizeof(o->repo), argv[++i]);
        } else if (ulab_streq(argv[i], "--setup-only")) {
            o->setup_only = 1;
        } else if (ulab_streq(argv[i], "--subscriber")) {
            o->subscriber_only = 1;
            o->setup_only = 1;
        } else if (ulab_streq(argv[i], "--network-id") && i + 1 < argc) {
            ulab_copy(o->subscriber_network_id,
                      sizeof(o->subscriber_network_id), argv[++i]);
        } else if (ulab_streq(argv[i], "--sim-csv") && i + 1 < argc) {
            ulab_copy(o->sim_csv_path, sizeof(o->sim_csv_path), argv[++i]);
        } else if (ulab_streq(argv[i], "--sim-type") && i + 1 < argc) {
            ulab_copy(o->sim_type, sizeof(o->sim_type), argv[++i]);
        } else if (ulab_streq(argv[i], "--print-world")) {
            o->print_world = 1;
        } else if (ulab_streq(argv[i], "--print-plan")) {
            o->print_plan = 1;
        } else if (ulab_streq(argv[i], "--cleanup")) {
            o->cleanup = 1;
            o->keep = 0;
        } else if (ulab_streq(argv[i], "--keep")) {
            o->keep = 1;
            o->cleanup = 0;
        } else if (ulab_streq(argv[i], "--keep-on-failure")) {
            o->keep_on_failure = 1;
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
    if (!opts.subscriber_only &&
        (opts.repo[0] == '\0' || strstr(opts.repo, "ukama") == NULL)) {
        printf("Missing --repo. Ukama repo root is MUST\n");
        usage();
        return ULAB_EUSAGE;
    }

    ulab_log_set_quiet(opts.quiet);
    ulab_log_set_verbose(opts.verbose);

    if (ulab_streq(argv[1], "validate")) {
        return runner_validate(&opts);
    }

    if (ulab_streq(argv[1], "dry-run")) {
        return runner_dry_run(&opts);
    }

    usage();
    return ULAB_EUSAGE;
}
