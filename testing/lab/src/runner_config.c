/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc. */

#include <jansson.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>

#include "runner.h"
#include "util.h"

int runner_load_scenario(const runner_opts_t *opts, scenario_t *scenario,
                          hardware_config_t *hardware, ulab_error_t *err) {
    memset(hardware, 0, sizeof(*hardware));
    if (scenario_load(opts->scenario_path, scenario, err)) return ULAB_ESCENARIO;
    if (opts->has_seed_override) scenario->seed = opts->seed_override;
    if (opts->provider[0] != '\0')
        ulab_copy(scenario->provider.type, sizeof(scenario->provider.type),
                   opts->provider);
    if (scenario_validate(scenario, err)) return ULAB_ESCENARIO;
    if (ulab_streq(scenario->provider.type, "hardware")) {
        if (opts->hardware_config[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "hardware provider requires --hardware <rack.json>");
            return ULAB_ESCENARIO;
        }
        if (hardware_config_load(opts->hardware_config, hardware, err) ||
            hardware_check_capacity(scenario, hardware, err) ||
            hardware_execution_supported(scenario, err))
            return ULAB_ESCENARIO;
    } else if (opts->hardware_config[0] != '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "--hardware requires provider.type: hardware or --provider hardware");
        return ULAB_ESCENARIO;
    }
    return ULAB_OK;
}

static json_t *preflight_json(const scenario_t *scenario, const world_t *world,
                               const hardware_config_t *hardware) {
    json_t *root;
    json_t *nodes;
    json_t *ues;
    size_t i;
    size_t physical_ues;
    int is_hardware;

    is_hardware = ulab_streq(scenario->provider.type, "hardware");
    nodes = json_array();
    ues = json_array();
    root = NULL;
    physical_ues = 0;
    if (nodes == NULL || ues == NULL) goto done;
    if (is_hardware) {
        for (i = 0; i < world->node_count; i++) {
            const node_t *node;

            node = &world->nodes[i];
            if (json_array_append_new(nodes, json_pack("{s:s,s:s,s:s}",
                "ref", node->ref, "id", node->id, "site_ref", node->site_ref)))
                goto done;
        }
        for (i = 0; i < world->ue_count; i++) {
            const ue_t *ue;

            ue = &world->ues[i];
            if (ue->hardware_slot == 0) continue;
            physical_ues++;
            if (json_array_append_new(ues, json_pack("{s:s,s:s,s:i,s:s,s:s,s:i}",
                "ref", ue->ref, "site_ref", ue->site_ref,
                "slot", (int)ue->hardware_slot, "iccid", ue->iccid,
                "imsi", ue->imsi, "iperf_port",
                (int)hardware->iperf_ports[ue->hardware_slot - 1]))) goto done;
        }
    }
    root = json_pack("{s:s,s:s,s:b,s:{s:I,s:I,s:I,s:I},s:O,s:O}",
        "scenario", scenario->name, "provider", scenario->provider.type,
        "execution_available", 1,
        "world", "networks", (json_int_t)world->network_count,
        "sites", (json_int_t)world->site_count,
        "nodes", (json_int_t)world->node_count,
        "ues", (json_int_t)world->ue_count,
        "physical_nodes", nodes, "physical_ues", ues);
    if (root != NULL && is_hardware) {
        if (json_object_set_new(root, "configured_slots", json_integer(hardware->ue_count)) ||
            json_object_set_new(root, "backend_only_sims",
                                json_integer(world->ue_count - physical_ues)) ||
            json_object_set_new(root, "media_host", json_string(hardware->media_host))) {
            json_decref(root);
            root = NULL;
        }
    }
done:
    json_decref(nodes);
    json_decref(ues);
    return root;
}

int runner_preflight(const runner_opts_t *opts) {
    scenario_t *scenario;
    hardware_config_t hardware;
    world_t world;
    ulab_error_t err;
    struct stat st;
    json_t *output;
    int rc;

    memset(&world, 0, sizeof(world));
    memset(&err, 0, sizeof(err));
    if (stat(opts->scenario_path, &st) == 0 && S_ISDIR(st.st_mode)) {
        fprintf(stderr, "preflight requires one scenario file, not a directory\n");
        return ULAB_EUSAGE;
    }
    scenario = calloc(1, sizeof(*scenario));
    if (scenario == NULL) return ULAB_EINTERNAL;
    rc = runner_load_scenario(opts, scenario, &hardware, &err);
    if (rc != ULAB_OK) goto done;
    if (scenario_execution_supported(scenario, &err)) {
        rc = ULAB_ESCENARIO;
        goto done;
    }
    if (world_generate(scenario, "preflight", &world, &err)) {
        rc = ULAB_EINTERNAL;
        goto done;
    }
    if (ulab_streq(scenario->provider.type, "hardware") &&
        hardware_bind_world(&hardware, &world, &err)) {
        rc = ULAB_ESCENARIO;
        goto done;
    }
    output = preflight_json(scenario, &world, &hardware);
    if (output == NULL || json_dumpf(output, stdout, JSON_INDENT(2)) != 0) {
        snprintf(err.msg, sizeof(err.msg), "cannot write preflight result");
        rc = ULAB_EINTERNAL;
    } else {
        putchar('\n');
    }
    json_decref(output);
done:
    if (rc != ULAB_OK) fprintf(stderr, "%s\n", err.msg);
    world_free(&world);
    free(scenario);
    return rc;
}
