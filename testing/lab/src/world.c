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

#include "world.h"
#include "util.h"

static int alloc_world(const scenario_t *s, world_t *w) {

    size_t sites     = s->world.networks * s->world.sites_per_network;
    size_t nodes_per = s->world.tower_per_site +
        s->world.amplifier_per_site +
        s->world.controller_per_site;
    size_t nodes    = sites * nodes_per;
    size_t ues      = sites * s->world.ues_per_site;
    size_t packages = s->package_count * s->world.networks;

    memset(w, 0, sizeof(*w));
    w->networks    = calloc(s->world.networks, sizeof(network_t));
    w->sites       = calloc(sites, sizeof(site_t));
    w->nodes       = calloc(nodes, sizeof(node_t));
    w->subscribers = calloc(ues, sizeof(subscriber_t));
    w->ues         = calloc(ues, sizeof(ue_t));
    w->packages    = calloc(packages, sizeof(package_t));
    if (!w->networks ||
        !w->sites ||
        !w->nodes ||
        !w->subscribers ||
        !w->ues ||
        !w->packages) {
        return ULAB_ERR;
    }

    w->network_count    = s->world.networks;
    w->site_count       = sites;
    w->node_count       = nodes;
    w->subscriber_count = ues;
    w->ue_count         = ues;
    w->package_count    = packages;

    return ULAB_OK;
}

static uint32_t seeded_u32(const scenario_t *s,
                           const char *scope,
                           size_t a,
                           size_t b,
                           size_t c) {
    char key[ULAB_MAX_LINE];

    snprintf(key, sizeof(key), "%s:%zu:%zu:%zu", scope, a, b, c);
    return ulab_hash32(key, s->seed);
}


static void ue_ip_for_site_index(size_t site_ue_idx,
                                 char *out,
                                 size_t out_len) {

    size_t host;
    size_t third;
    size_t fourth;

    /*
     * UE CIDR is 192.168.8.0/22. 192.168.8.1 is used by EPC TUN,
     * so the first UE starts at 192.168.8.2. IPs are reused per site
     * because each tower node owns its own isolated UE subnet.
     */
    host = site_ue_idx + 2;
    third = ULAB_UE_SUBNET_BASE_C + (host / 256);
    fourth = host % 256;

    snprintf(out, out_len, "%u.%u.%zu.%zu",
             ULAB_UE_SUBNET_BASE_A,
             ULAB_UE_SUBNET_BASE_B,
             third,
             fourth);
}

static const char *pick_package(const scenario_t *s, size_t ue_idx) {

    uint32_t roll;
    uint32_t acc = 0;
    size_t i;

    roll = seeded_u32(s, "package", ue_idx, 0, 0) % 100u;

    for (i = 0; i < s->package_count; i++) {
        acc += s->packages[i].assign_percent;
        if (roll < acc) {
            return s->packages[i].ref;
        }
    }

    return s->packages[s->package_count - 1].ref;
}

static void make_package_ref(char *out, size_t out_len,
                             const char *base_ref,
                             const char *network_ref) {

    snprintf(out, out_len, "%.96s__%.24s", base_ref, network_ref);
}

static void make_package_name(char *out, size_t out_len,
                              const char *base_name,
                              const network_t *network) {

    /*
     * Package names are globally unique in the new backend.
     * Scope each scenario package template to the generated network id.
     * For --subscriber, network->bff_id is filled from --network-id later,
     * but network->id still contains the unique lab run id.
     */
    snprintf(out, out_len, "%.120s %.120s", base_name, network->id);
}

static void add_package_for_network(world_t *w,
                                    const scenario_t *s,
                                    size_t *idx,
                                    const network_t *network,
                                    size_t package_spec_idx) {

    const package_spec_t *spec = &s->packages[package_spec_idx];
    package_t *p = &w->packages[(*idx)++];

    make_package_ref(p->ref, sizeof(p->ref), spec->ref, network->ref);
    ulab_copy(p->base_ref, sizeof(p->base_ref), spec->ref);
    ulab_copy(p->network_ref, sizeof(p->network_ref), network->ref);
    make_package_name(p->name, sizeof(p->name), spec->name, network);
    p->data_mb       = spec->data_mb;
    p->duration_days = spec->duration_days;
    p->amount        = spec->amount;
}

static void add_node(world_t *w, size_t *idx, const char *type,
                     const site_t *site, size_t n) {
    node_t *node = &w->nodes[(*idx)++];

    snprintf(node->ref, sizeof(node->ref), "%.32s-%.64s-%03zu", type,
             site->ref, n);
    snprintf(node->id, sizeof(node->id), "%.240s-%.120s", w->run_id,
             node->ref);
    snprintf(node->name, sizeof(node->name), "%.255s", node->id);
    snprintf(node->type, sizeof(node->type), "%s", type);
    snprintf(node->site_ref, sizeof(node->site_ref), "%s", site->ref);
    snprintf(node->network_ref, sizeof(node->network_ref), "%s",
             site->network_ref);
}

int world_generate(const scenario_t *s,
                   const char *run_id,
                   world_t *w,
                   ulab_error_t *err) {
    size_t i;
    size_t j;
    size_t k;
    size_t site_idx = 0;
    size_t node_idx    = 0;
    size_t package_idx = 0;
    size_t ue_idx      = 0;

    if (alloc_world(s, w) != ULAB_OK) {
        snprintf(err->msg, sizeof(err->msg), "world allocation failed");
        return ULAB_ERR;
    }
    ulab_copy(w->run_id, sizeof(w->run_id), run_id);
    w->seed = s->seed;

    for (i = 0; i < w->network_count; i++) {
        network_t *net = &w->networks[i];
        snprintf(net->ref,  sizeof(net->ref),  "net-%03zu", i + 1);
        snprintf(net->id,   sizeof(net->id),   "%.240s-%.120s", run_id,
                 net->ref);
        snprintf(net->name, sizeof(net->name), "%.255s", net->id);

        for (j = 0; j < s->package_count; j++) {
            add_package_for_network(w, s, &package_idx, net, j);
        }

        for (j = 0; j < s->world.sites_per_network; j++) {
            site_t *site = &w->sites[site_idx++];
            snprintf(site->ref, sizeof(site->ref), "site-%03zu", site_idx);
            snprintf(site->id,  sizeof(site->id), "%.240s-%.120s", run_id,
                     site->ref);
            snprintf(site->name, sizeof(site->name), "%.255s", site->id);
            snprintf(site->network_ref, sizeof(site->network_ref), "%s",
                     net->ref);

            for (k = 0; k < s->world.tower_per_site; k++) {
                add_node(w, &node_idx, ULAB_NODE_TOWER, site, k + 1);
            }

            for (k = 0; k < s->world.amplifier_per_site; k++) {
                add_node(w, &node_idx, ULAB_NODE_AMPLIFIER, site, k + 1);
            }

            for (k = 0; k < s->world.controller_per_site; k++) {
                add_node(w, &node_idx, ULAB_NODE_CONTROLLER, site, k + 1);
            }

            for (k = 0; k < s->world.ues_per_site; k++) {
                subscriber_t *sub = &w->subscribers[ue_idx];
                ue_t *ue = &w->ues[ue_idx];
                size_t num = ue_idx + 1;
                uint32_t phone = seeded_u32(s, "phone", ue_idx, i, j) %
                    10000000u;

                snprintf(sub->ref, sizeof(sub->ref), "sub-%06zu", num);
                snprintf(sub->id, sizeof(sub->id), "%.240s-%.120s", run_id,
                         sub->ref);
                snprintf(sub->name, sizeof(sub->name), "Lab User %zu",
                         num);
                snprintf(sub->email, sizeof(sub->email),
                         "%.180s-%06zu@ukama.test", run_id, num);
                snprintf(sub->phone, sizeof(sub->phone), "+1555%07u", phone);
                snprintf(sub->network_ref, sizeof(sub->network_ref), "%s",
                         net->ref);
                snprintf(sub->site_ref, sizeof(sub->site_ref), "%s",
                         site->ref);

                snprintf(ue->ref, sizeof(ue->ref), "ue-%06zu", num);
                snprintf(ue->id, sizeof(ue->id), "%.240s-%.120s",
                         run_id, ue->ref);
                snprintf(ue->iccid, sizeof(ue->iccid),
                         "890100%013zu", num);
                snprintf(ue->imsi, sizeof(ue->imsi), "001010%09zu", num);
                snprintf(ue->subscriber_ref, sizeof(ue->subscriber_ref),
                         "%s", sub->ref);
                snprintf(ue->network_ref, sizeof(ue->network_ref), "%s",
                         net->ref);
                snprintf(ue->site_ref, sizeof(ue->site_ref), "%s",
                         site->ref);
                make_package_ref(ue->package_ref, sizeof(ue->package_ref),
                                 pick_package(s, ue_idx), net->ref);
                ue_ip_for_site_index(k, ue->ip, sizeof(ue->ip));
                ue_idx++;
            }
        }
    }
    return ULAB_OK;
}

void world_free(world_t *w) {
    if (w == NULL) return;
    free(w->networks);
    free(w->sites);
    free(w->nodes);
    free(w->subscribers);
    free(w->ues);
    free(w->packages);
    memset(w, 0, sizeof(*w));
}

network_t *world_network_by_ref(world_t *w, const char *ref) {
    size_t i;
    for (i = 0; i < w->network_count; i++) {
        if (ulab_streq(w->networks[i].ref, ref)) return &w->networks[i];
    }
    return NULL;
}

site_t *world_site_by_ref(world_t *w, const char *ref) {
    size_t i;
    for (i = 0; i < w->site_count; i++) {
        if (ulab_streq(w->sites[i].ref, ref)) return &w->sites[i];
    }
    return NULL;
}

node_t *world_node_by_ref(world_t *w, const char *ref) {
    size_t i;
    for (i = 0; i < w->node_count; i++) {
        if (ulab_streq(w->nodes[i].ref, ref)) return &w->nodes[i];
    }
    return NULL;
}

ue_t *world_ue_by_ref(world_t *w, const char *ref) {
    size_t i;
    for (i = 0; i < w->ue_count; i++) {
        if (ulab_streq(w->ues[i].ref, ref)) return &w->ues[i];
    }
    return NULL;
}

package_t *world_package_by_ref(world_t *w, const char *ref) {
    size_t i;
    for (i = 0; i < w->package_count; i++) {
        if (ulab_streq(w->packages[i].ref, ref)) return &w->packages[i];
    }
    return NULL;
}

package_t *world_package_for_network(world_t *w,
                                     const char *package_ref,
                                     const char *network_ref) {
    size_t i;

    if (package_ref == NULL || network_ref == NULL) {
        return NULL;
    }

    for (i = 0; i < w->package_count; i++) {
        package_t *p = &w->packages[i];

        if (!ulab_streq(p->network_ref, network_ref)) {
            continue;
        }
        if (ulab_streq(p->ref, package_ref) ||
            ulab_streq(p->base_ref, package_ref)) {
            return p;
        }
    }

    return NULL;
}

int world_add_ues(world_t *w, const scenario_t *s, const char *phase,
                  const selector_t *sites, uint32_t count_per_site,
                  const char *package_ref, ulab_error_t *err) {
    (void)w;
    (void)s;
    (void)phase;
    (void)sites;
    (void)count_per_site;
    (void)package_ref;

    snprintf(err->msg, sizeof(err->msg),
             "create_ues is reserved for v1.1 in this build");
    return ULAB_ERR;
}

int world_write_json(const world_t *w, const char *path) {

    FILE *f = fopen(path, "w");
    if (!f) return ULAB_ERR;

    fprintf(f, "{\n");
    fprintf(f, "  \"run_id\": \"%s\",\n",   w->run_id);
    fprintf(f, "  \"seed\": %u,\n",          w->seed);
    fprintf(f, "  \"networks\": %zu,\n",    w->network_count);
    fprintf(f, "  \"sites\": %zu,\n",       w->site_count);
    fprintf(f, "  \"nodes\": %zu,\n",       w->node_count);
    fprintf(f, "  \"subscribers\": %zu,\n", w->subscriber_count);
    fprintf(f, "  \"ues\": %zu,\n",         w->ue_count);
    fprintf(f, "  \"packages\": %zu\n",     w->package_count);
    fprintf(f, "}\n");
    fclose(f);

    return ULAB_OK;
}

void world_print(const world_t *w) {

    if (w == NULL) {
        return;
    }

    printf("World:\n");
    printf("  seed:        %u\n", w->seed);
    printf("  networks:    %zu\n", w->network_count);
    printf("  sites:       %zu\n", w->site_count);
    printf("  nodes:       %zu\n", w->node_count);
    printf("  subscribers: %zu\n", w->subscriber_count);
    printf("  ues:         %zu\n", w->ue_count);
    printf("  packages:    %zu\n", w->package_count);
}
