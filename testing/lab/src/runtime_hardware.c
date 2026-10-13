/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc. */
#include <inttypes.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#include "runtime_hardware.h"
#include "rack_client.h"
#include "log.h"
#include "util.h"

struct hardware_runtime {
    hardware_config_t config;
    rack_client_t *rack;
    int power_touched;
    int media_touched;
    int media_running;
    unsigned int slots_touched;
};

static uint64_t now_ms(void) {
    struct timespec t;
    clock_gettime(CLOCK_MONOTONIC, &t);
    return (uint64_t)t.tv_sec * 1000u + (uint64_t)t.tv_nsec / 1000000u;
}

static int fail(ulab_error_t *err, const char *text) {
    snprintf(err->msg, sizeof(err->msg), "hardware: %s", text);
    return ULAB_ERR;
}

int runtime_is_hardware(const runtime_t *rt) {
    return rt != NULL && ulab_streq(rt->provider, "hardware");
}

int hardware_runtime_healthy(runtime_t *rt, ulab_error_t *err) {
    if (!runtime_is_hardware(rt))
        return ULAB_OK;
    if (rt->hardware == NULL)
        return fail(err, "runtime is not open");
    /* Backend-only scenarios have no rack resources to acquire. */
    if (rt->hardware->config.site_count == 0)
        return ULAB_OK;
    return rack_client_healthy(rt->hardware->rack, err);
}

static int command(runtime_t *rt, const char *text, char id[65], json_t **result,
                   ulab_error_t *err) {
    *result = NULL;
    if (rt->hardware == NULL)
        return fail(err, "runtime is not open");
    return rack_client_call(rt->hardware->rack, text, id, result, err);
}

static int status(runtime_t *rt, json_t **out, ulab_error_t *err) {
    char id[65];
    json_int_t age;
    if (command(rt, "ue status", id, out, err))
        return ULAB_ERR;
    age = json_integer_value(json_object_get(*out, "age_ms"));
    if (!json_is_true(json_object_get(*out, "reachable")) ||
        !json_is_true(json_object_get(*out, "lease_ready")) ||
        !json_is_integer(json_object_get(*out, "age_ms")) || age < 0 || age > 5000 ||
        !json_is_array(json_object_get(json_object_get(*out, "gadget"), "slots"))) {
        json_decref(*out);
        *out = NULL;
        return fail(err, "UE gadget status unavailable, stale or lease not ready");
    }
    return ULAB_OK;
}

static json_t *slot_value(json_t *root, uint32_t slot) {
    json_t *slots;
    json_t *entry;
    size_t i;
    slots = json_object_get(json_object_get(root, "gadget"), "slots");
    json_array_foreach(slots, i, entry) {
        if (json_integer_value(json_object_get(entry, "slot")) == slot)
            return entry;
    }
    return NULL;
}

int hardware_runtime_open(runtime_t *rt, const hardware_config_t *config, const world_t *world,
                          ulab_error_t *err) {
    struct hardware_runtime *h;
    json_t *reply = NULL;
    uint64_t deadline;
    size_t i;

    h = calloc(1, sizeof(*h));
    if (h == NULL)
        return ULAB_EINTERNAL;
    rt->hardware = h;
    h->config = *config;
    if (world->site_count == 0) {
        h->config.site_count = 0;
        return ULAB_OK;
    }
    if (rack_client_open(&h->rack, config, rt->run_dir, err))
        return ULAB_ERR;
    for (i = 0; i < world->ue_count; i++) {
        if (world->ues[i].hardware_slot)
            break;
    }
    if (i == world->ue_count)
        return ULAB_OK;
    /* Session acquisition lets the controller acquire its existing UE lease.
     * This is protocol readiness only; no modem attach or backend setup probe. */
    deadline = now_ms() + config->command_timeout_seconds * 1000u;
    while (status(rt, &reply, err)) {
        if (now_ms() >= deadline || rack_client_healthy(h->rack, err))
            return ULAB_ERR;
        sleep(1);
    }
    for (i = 0; i < world->ue_count; i++) {
        if (world->ues[i].hardware_slot && slot_value(reply, world->ues[i].hardware_slot) == NULL) {
            json_decref(reply);
            return fail(err, "a configured UE slot is absent from the gadget");
        }
    }
    json_decref(reply);
    ulab_error_clear(err);
    return ULAB_OK;
}

int hardware_runtime_power(runtime_t *rt, int enabled, ulab_error_t *err) {
    json_t *reply = NULL;
    char id[65];
    int rc;
    if (rt->hardware == NULL)
        return fail(err, "runtime is not open");
    if (enabled)
        rt->hardware->power_touched = 1;
    rc = command(rt, enabled ? "dut on" : "dut off", id, &reply, err);
    if (!rc && (!json_is_boolean(json_object_get(reply, "enabled")) ||
                json_is_true(json_object_get(reply, "enabled")) != enabled))
        rc = fail(err, "DUT power acknowledgement does not match the request");
    json_decref(reply);
    if (!rc && !enabled)
        rt->hardware->power_touched = 0;
    return rc;
}

int hardware_runtime_media(runtime_t *rt, int enabled, ulab_error_t *err) {
    struct hardware_runtime *h = rt->hardware;
    json_t *reply = NULL;
    json_t *ports;
    const char *host;
    char id[65];
    int rc;
    size_t i;
    if (h == NULL)
        return fail(err, "runtime is not open");
    if (enabled && h->media_running)
        return hardware_runtime_healthy(rt, err);
    if (enabled)
        h->media_touched = 1;
    rc = command(rt, enabled ? "media start" : "media stop", id, &reply, err);
    if (!rc && (!json_is_boolean(json_object_get(reply, "running")) ||
                json_is_true(json_object_get(reply, "running")) != enabled))
        rc = fail(err, "media acknowledgement does not match the request");
    if (!rc && enabled) {
        host = json_string_value(json_object_get(reply, "traffic_address"));
        ports = json_object_get(reply, "iperf_ports");
        if (host == NULL || strcmp(host, h->config.media_host) != 0 ||
            json_array_size(ports) != ULAB_HARDWARE_MAX_UES)
            rc = fail(err, "media endpoint differs from the rack inventory");
        for (i = 0; !rc && i < ULAB_HARDWARE_MAX_UES; i++) {
            if (json_integer_value(json_array_get(ports, i)) != h->config.iperf_ports[i])
                rc = fail(err, "media ports differ from the rack inventory");
        }
    }
    json_decref(reply);
    if (!rc)
        h->media_running = enabled;
    if (!rc && !enabled)
        h->media_touched = 0;
    return rc;
}

static int selected_slots(const world_t *world, const selector_result_t *ues, unsigned int *mask,
                          ulab_error_t *err) {
    size_t i;
    *mask = 0;
    for (i = 0; i < ues->count; i++) {
        uint32_t slot;
        if (ues->idx[i] >= world->ue_count)
            return fail(err, "invalid UE selection");
        slot = world->ues[ues->idx[i]].hardware_slot;
        if (slot == 0 || slot > ULAB_HARDWARE_MAX_UES)
            return fail(err, "selected UE is a backend-only SIM or has no gadget slot");
        if (*mask & (1u << (slot - 1)))
            return fail(err, "duplicate UE slot selection");
        *mask |= 1u << (slot - 1);
    }
    return ULAB_OK;
}

static int run_ues(runtime_t *rt, const world_t *world, const selector_result_t *ues,
                   const char *action, uint64_t amount_mb, uint64_t deadline, ulab_error_t *err) {
    struct hardware_runtime *h = rt->hardware;
    char ids[ULAB_HARDWARE_MAX_UES][65] = {{0}};
    unsigned int pending;
    unsigned int mask;
    unsigned int slot;
    int traffic;
    int failed = 0;

    if (h == NULL)
        return fail(err, "runtime is not open");
    if (selected_slots(world, ues, &mask, err))
        return ULAB_ERR;
    traffic = ulab_streq(action, "uplink");
    if (traffic && (amount_mb == 0 || amount_mb > UINT32_MAX))
        return fail(err, "traffic amount_mb must be in 1..4294967295");
    if (!traffic && !ulab_streq(action, "start") && !ulab_streq(action, "connect") &&
        !ulab_streq(action, "disconnect") && !ulab_streq(action, "shutdown"))
        return fail(err, "unsupported UE action");
    pending = mask;
    for (slot = 1; slot <= ULAB_HARDWARE_MAX_UES; slot++) {
        char text[128];
        json_t *reply = NULL;
        int rc;
        if (!(mask & (1u << (slot - 1))))
            continue;
        if (now_ms() >= deadline)
            return fail(err, "UE operation timed out");
        if (traffic)
            snprintf(text, sizeof(text), "ue uplink %u --mb %" PRIu64, slot, amount_mb);
        else
            snprintf(text, sizeof(text), "ue %s %u", action, slot);
        h->slots_touched |= 1u << (slot - 1);
        rc = command(rt, text, ids[slot - 1], &reply, err);
        if (!rc &&
            (!json_is_true(json_object_get(reply, "queued")) ||
             json_integer_value(json_object_get(reply, "slot")) != slot ||
             !ulab_streq(json_string_value(json_object_get(reply, "request_id")), ids[slot - 1])))
            rc = fail(err, "invalid UE queue acknowledgement");
        json_decref(reply);
        if (rc)
            return ULAB_ERR;
    }
    while (pending) {
        json_t *reply = NULL;
        if (now_ms() >= deadline)
            return fail(err, "UE operation timed out");
        if (status(rt, &reply, err))
            return ULAB_ERR;
        for (slot = 1; slot <= ULAB_HARDWARE_MAX_UES; slot++) {
            json_t *entry;
            json_t *operation;
            json_t *submission;
            json_t *list;
            const char *state;
            size_t i;
            if (!(pending & (1u << (slot - 1))))
                continue;
            list = json_object_get(reply, "submissions");
            json_array_foreach(list, i, submission) {
                if (json_integer_value(json_object_get(submission, "slot")) == slot &&
                    ulab_streq(json_string_value(json_object_get(submission, "request_id")),
                               ids[slot - 1])) {
                    state = json_string_value(json_object_get(submission, "status"));
                    if (ulab_streq(state, "rejected") || ulab_streq(state, "unknown")) {
                        json_decref(reply);
                        return fail(err, "UE gadget rejected the request");
                    }
                }
            }
            entry = slot_value(reply, slot);
            if (entry == NULL) {
                json_decref(reply);
                return fail(err, "UE slot disappeared");
            }
            operation = json_object_get(entry, "operation");
            if (!ulab_streq(json_string_value(json_object_get(operation, "request_id")),
                            ids[slot - 1]))
                continue;
            if (!ulab_streq(json_string_value(json_object_get(operation, "action")), action)) {
                json_decref(reply);
                return fail(err, "UE operation action does not match the request");
            }
            state = json_string_value(json_object_get(operation, "status"));
            if (ulab_streq(state, "queued") || ulab_streq(state, "running"))
                continue;
            if (!ulab_streq(state, "completed")) {
                failed = 1;
            } else if (traffic) {
                json_t *report = json_object_get(json_object_get(entry, "state"), "test");
                /* Observe the command's result, without re-counting bytes or
                 * introducing new transfer/accounting assertions. */
                if (!ulab_streq(json_string_value(json_object_get(report, "request_id")),
                                ids[slot - 1]) ||
                    !ulab_streq(json_string_value(json_object_get(report, "status")), "passed"))
                    failed = 1;
            }
            pending &= ~(1u << (slot - 1));
        }
        json_decref(reply);
        if (pending && now_ms() >= deadline)
            return fail(err, "UE operation timed out");
        if (pending)
            sleep(1);
    }
    return failed ? fail(err, "UE operation or traffic command failed") : ULAB_OK;
}

int hardware_runtime_ues(runtime_t *rt, const world_t *world, const selector_result_t *ues,
                         const char *action, uint64_t amount_mb, ulab_error_t *err) {
    if (rt->hardware == NULL)
        return fail(err, "runtime is not open");
    return run_ues(rt, world, ues, action, amount_mb,
                   now_ms() + (uint64_t)rt->hardware->config.operation_timeout_seconds * 1000u,
                   err);
}

static int observe(runtime_t *rt, const world_t *world, const selector_result_t *ues, int attached,
                   ulab_error_t *err) {
    json_t *reply = NULL;
    unsigned int mask;
    size_t i;
    int rc = ULAB_OK;
    if (selected_slots(world, ues, &mask, err) || status(rt, &reply, err))
        return ULAB_ERR;
    for (i = 0; i < ues->count; i++) {
        json_t *entry;
        json_t *connected;
        json_t *age;
        entry = slot_value(reply, world->ues[ues->idx[i]].hardware_slot);
        age = json_object_get(entry, "state_age_ms");
        if (!json_is_integer(age) || json_integer_value(age) < 0 ||
            json_integer_value(age) > 5000) {
            rc = fail(err, "UE slot state is unavailable or stale");
            break;
        }
        connected =
            json_object_get(json_object_get(json_object_get(entry, "state"), "modem"), "connected");
        if (!json_is_boolean(connected) || json_is_true(connected) != attached) {
            rc = fail(err, attached ? "UE data session is not connected"
                                    : "UE data session is not detached");
            break;
        }
    }
    json_decref(reply);
    return rc;
}

int hardware_runtime_sessions(runtime_t *rt, const world_t *world, const selector_result_t *ues,
                              ulab_error_t *err) {
    return observe(rt, world, ues, 1, err);
}

int hardware_runtime_wait(runtime_t *rt, world_t *world, const selector_result_t *ues, int attached,
                          unsigned int timeout, ulab_error_t *err) {
    uint64_t deadline = now_ms() + (uint64_t)timeout * 1000u;
    size_t i;
    if (attached && run_ues(rt, world, ues, "connect", 0, deadline, err))
        return ULAB_ERR;
    while (observe(rt, world, ues, attached, err)) {
        if (now_ms() >= deadline || hardware_runtime_healthy(rt, err))
            return ULAB_ERR;
        sleep(1);
    }
    for (i = 0; i < ues->count; i++)
        world->ues[ues->idx[i]].attached = attached;
    ulab_error_clear(err);
    return ULAB_OK;
}

int hardware_runtime_cleanup_ues(runtime_t *rt, const world_t *world, const selector_result_t *ues,
                                 ulab_error_t *err) {
    selector_result_t selected;
    size_t indices[ULAB_HARDWARE_MAX_UES];
    size_t i;
    int rc;
    if (rt->hardware == NULL || rt->hardware->rack == NULL)
        return ULAB_OK;
    selected.idx = indices;
    selected.count = 0;
    for (i = 0; i < world->ue_count; i++) {
        uint32_t slot = world->ues[i].hardware_slot;
        size_t j;
        if (!slot || slot > ULAB_HARDWARE_MAX_UES ||
            !(rt->hardware->slots_touched & (1u << (slot - 1))))
            continue;
        if (ues != NULL) {
            for (j = 0; j < ues->count && ues->idx[j] != i; j++) {
            }
            if (j == ues->count)
                continue;
        }
        if (selected.count >= ULAB_HARDWARE_MAX_UES)
            return fail(err, "duplicate cleanup slot");
        indices[selected.count++] = i;
    }
    if (selected.count == 0)
        return ULAB_OK;
    rc = hardware_runtime_ues(rt, world, &selected, "disconnect", 0, err);
    if (rc == ULAB_OK) {
        for (i = 0; i < selected.count; i++) {
            world->ues[indices[i]].attached = 0;
        }
    }
    return rc;
}

int hardware_runtime_cleanup(runtime_t *rt, ulab_error_t *err) {
    struct hardware_runtime *h = rt->hardware;
    ulab_error_t current = {{0}};
    int rc = ULAB_OK;
    if (h == NULL)
        return ULAB_OK;
    if (h->media_touched && hardware_runtime_media(rt, 0, &current)) {
        *err = current;
        rc = ULAB_ERR;
    }
    if (h->power_touched && hardware_runtime_power(rt, 0, &current)) {
        *err = current;
        rc = ULAB_ERR;
    }
    if (rack_client_close(&h->rack, &current)) {
        *err = current;
        rc = ULAB_ERR;
    }
    return rc;
}

void hardware_runtime_close(runtime_t *rt) {
    ulab_error_t err = {{0}};
    if (rt->hardware == NULL)
        return;
    /* Normal cleanup already stopped run-owned resources. If cleanup was
     * explicitly disabled, only release the lease, preserving that policy. */
    if (rack_client_close(&rt->hardware->rack, &err))
        ulab_log_warn("%s", err.msg);
    free(rt->hardware);
    rt->hardware = NULL;
}
