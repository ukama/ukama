/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>

#include "runtime.h"
#include "util.h"

int main(int argc, char **argv) {
    runtime_t rt = {0};
    world_t world = {0};
    site_t site = {0};
    ulab_error_t err = {0};
    size_t i;
    int rc;

    if (argc < 3 || ulab_copy(rt.run_dir, sizeof(rt.run_dir), argv[1])) {
        return ULAB_EUSAGE;
    }

    ulab_copy(site.ref, sizeof(site.ref), "site-001");
    ulab_copy(site.network_ref, sizeof(site.network_ref), "net-001");
    world.sites = &site;
    world.site_count = 1;
    world.node_count = (size_t)argc - 2;
    world.nodes = calloc(world.node_count, sizeof(*world.nodes));
    if (world.nodes == NULL) {
        return ULAB_ERR;
    }

    for (i = 0; i < world.node_count; i++) {
        if (ulab_copy(world.nodes[i].id, sizeof(world.nodes[i].id), argv[i + 2])) {
            free(world.nodes);
            return ULAB_EUSAGE;
        }
        ulab_copy(world.nodes[i].site_ref, sizeof(world.nodes[i].site_ref), site.ref);
        ulab_copy(world.nodes[i].type, sizeof(world.nodes[i].type), ULAB_NODE_TOWER);
    }

    rc = runtime_load_workload_sites(&rt, &world, &err);
    if (rc != ULAB_OK) {
        fprintf(stderr, "%s\n", err.msg);
    }
    free(world.nodes);
    return rc;
}
