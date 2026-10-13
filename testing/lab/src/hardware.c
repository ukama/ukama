/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc. */

#include <arpa/inet.h>
#include <ctype.h>
#include <inttypes.h>
#include <jansson.h>
#include <stdio.h>
#include <string.h>

#include "hardware.h"
#include "util.h"

static int fail(ulab_error_t *err, const char *message) {
    snprintf(err->msg, sizeof(err->msg), "hardware: %s", message);
    return ULAB_ERR;
}

static int keys(json_t *object, const char *const *allowed,
                 ulab_error_t *err) {
    const char *key;
    json_t *value;
    size_t i;

    if (!json_is_object(object)) return fail(err, "expected a JSON object");
    json_object_foreach(object, key, value) {
        (void)value;
        for (i = 0; allowed[i] != NULL; i++) {
            if (ulab_streq(key, allowed[i])) break;
        }
        if (allowed[i] == NULL) {
            snprintf(err->msg, sizeof(err->msg),
                     "hardware: unknown field %.128s", key);
            return ULAB_ERR;
        }
    }
    return ULAB_OK;
}

static int string_field(json_t *object, const char *key, char *out,
                         size_t length, ulab_error_t *err) {
    json_t *value;
    const char *text;
    size_t i;

    value = json_object_get(object, key);
    text = json_string_value(value);
    if (text == NULL || text[0] == '\0' ||
        json_string_length(value) >= length ||
        strlen(text) != json_string_length(value)) goto invalid;
    for (i = 0; text[i] != '\0'; i++) {
        if (iscntrl((unsigned char)text[i])) goto invalid;
    }
    return ulab_copy(out, length, text);

invalid:
    snprintf(err->msg, sizeof(err->msg),
             "hardware: %s must be a nonempty string shorter than %zu bytes",
             key, length);
    return ULAB_ERR;
}

static int integer_field(json_t *value, uint32_t min, uint32_t max,
                          uint32_t *out, const char *name,
                          ulab_error_t *err) {
    json_int_t number;

    number = json_integer_value(value);
    if (!json_is_integer(value) || number < min || number > max) {
        snprintf(err->msg, sizeof(err->msg),
                 "hardware: %s must be an integer in %u..%u", name, min, max);
        return ULAB_ERR;
    }
    *out = (uint32_t)number;
    return ULAB_OK;
}

static int token(const char *text) {
    size_t i;

    if (text[0] == '\0' || text[0] == '-') return 0;
    for (i = 0; text[i] != '\0'; i++) {
        if (!isalnum((unsigned char)text[i]) &&
            strchr("-_.:", text[i]) == NULL) return 0;
    }
    return 1;
}

static int digits(const char *text, size_t max) {
    size_t i;

    if (text[0] == '\0' || strlen(text) > max) return 0;
    for (i = 0; text[i] != '\0'; i++) {
        if (text[i] < '0' || text[i] > '9') return 0;
    }
    return 1;
}

static const hardware_site_t *find_site(const hardware_config_t *config,
                                        const char *ref) {
    size_t i;

    for (i = 0; i < config->site_count; i++) {
        if (ulab_streq(config->sites[i].ref, ref)) return &config->sites[i];
    }
    return NULL;
}

static const hardware_ue_t *find_ue(const hardware_config_t *config,
                                    const char *ref) {
    size_t i;

    for (i = 0; i < config->ue_count; i++) {
        if (ulab_streq(config->ues[i].ref, ref)) return &config->ues[i];
    }
    return NULL;
}

static int load_rack(json_t *root, hardware_config_t *config,
                      ulab_error_t *err) {
    static const char *const fields[] = {
        "host", "user", "ssh_port", "identity_file", "known_hosts_file",
        "command_timeout_seconds", "operation_timeout_seconds", "ssh_path", NULL
    };
    json_t *rack;
    json_t *port;

    rack = json_object_get(root, "rack");
    if (keys(rack, fields, err) ||
        string_field(rack, "host", config->host, sizeof(config->host), err) ||
        string_field(rack, "user", config->user, sizeof(config->user), err) ||
        string_field(rack, "identity_file", config->identity_file,
                     sizeof(config->identity_file), err) ||
        string_field(rack, "known_hosts_file", config->known_hosts_file,
                     sizeof(config->known_hosts_file), err)) return ULAB_ERR;
    if (!token(config->host) || !token(config->user))
        return fail(err, "rack host/user must be plain names, not commands or URLs");
    if (config->identity_file[0] != '/' || config->known_hosts_file[0] != '/')
        return fail(err, "SSH identity_file and known_hosts_file must be absolute paths");
    if (strpbrk(config->known_hosts_file, "\"\\") != NULL)
        return fail(err, "known_hosts_file cannot contain quotes or backslashes");
    ulab_copy(config->ssh_path, sizeof(config->ssh_path), "/usr/bin/ssh");
    if (json_object_get(rack, "ssh_path") != NULL &&
        string_field(rack, "ssh_path", config->ssh_path, sizeof(config->ssh_path), err)) return ULAB_ERR;
    if (config->ssh_path[0] != '/') return fail(err, "ssh_path must be absolute");
    config->command_timeout_seconds = 20;
    config->operation_timeout_seconds = 1020;
    port = json_object_get(rack, "command_timeout_seconds");
    if (port != NULL && integer_field(port, 1, 120, &config->command_timeout_seconds,
                                      "command_timeout_seconds", err)) return ULAB_ERR;
    port = json_object_get(rack, "operation_timeout_seconds");
    if (port != NULL && integer_field(port, 1, 3720, &config->operation_timeout_seconds,
                                      "operation_timeout_seconds", err)) return ULAB_ERR;
    config->ssh_port = 22;
    port = json_object_get(rack, "ssh_port");
    if (port != NULL && integer_field(port, 1, 65535, &config->ssh_port,
                                      "ssh_port", err)) return ULAB_ERR;
    return ULAB_OK;
}

static int load_media(json_t *root, hardware_config_t *config,
                       ulab_error_t *err) {
    static const char *const fields[] = {"host", "iperf_ports", NULL};
    json_t *media;
    json_t *ports;
    struct in_addr address;
    size_t i;

    media = json_object_get(root, "media");
    if (keys(media, fields, err) ||
        string_field(media, "host", config->media_host,
                     sizeof(config->media_host), err)) return ULAB_ERR;
    if (inet_pton(AF_INET, config->media_host, &address) != 1)
        return fail(err, "media.host must be the media server's IPv4 traffic address");
    ports = json_object_get(media, "iperf_ports");
    if (!json_is_array(ports) || json_array_size(ports) != ULAB_HARDWARE_MAX_UES)
        return fail(err, "media.iperf_ports requires one port for each of slots 1 and 2");
    for (i = 0; i < ULAB_HARDWARE_MAX_UES; i++) {
        if (integer_field(json_array_get(ports, i), 1, 65535,
                           &config->iperf_ports[i], "iperf port", err))
            return ULAB_ERR;
    }
    if (config->iperf_ports[0] == config->iperf_ports[1])
        return fail(err, "iperf ports must be distinct");
    return ULAB_OK;
}

static int load_sites(json_t *root, hardware_config_t *config,
                       ulab_error_t *err) {
    static const char *const fields[] = {
        "ref", "tnode_id", "anode_id", "cnode_id", NULL
    };
    json_t *array;
    json_t *entry;
    size_t i;
    size_t j;
    size_t n;
    size_t m;

    array = json_object_get(root, "sites");
    if (!json_is_array(array) || json_array_size(array) > ULAB_HARDWARE_MAX_SITES)
        return fail(err, "sites must be an array with at most 16 entries");
    config->site_count = json_array_size(array);
    json_array_foreach(array, i, entry) {
        hardware_site_t *site;
        const char *ids[3];

        site = &config->sites[i];
        if (keys(entry, fields, err) ||
            string_field(entry, "ref", site->ref, sizeof(site->ref), err) ||
            string_field(entry, "tnode_id", site->tnode_id,
                         sizeof(site->tnode_id), err) ||
            string_field(entry, "anode_id", site->anode_id,
                         sizeof(site->anode_id), err) ||
            string_field(entry, "cnode_id", site->cnode_id,
                         sizeof(site->cnode_id), err)) return ULAB_ERR;
        ids[0] = site->tnode_id;
        ids[1] = site->anode_id;
        ids[2] = site->cnode_id;
        if (!token(site->ref)) return fail(err, "invalid site ref");
        for (n = 0; n < 3; n++) {
            if (!token(ids[n])) return fail(err, "invalid physical node ID");
            for (m = 0; m < n; m++) {
                if (ulab_streq(ids[n], ids[m]))
                    return fail(err, "duplicate physical node ID");
            }
        }
        for (j = 0; j < i; j++) {
            const hardware_site_t *prior;

            prior = &config->sites[j];
            if (ulab_streq(site->ref, prior->ref))
                return fail(err, "duplicate site ref");
            for (n = 0; n < 3; n++) {
                if (ulab_streq(ids[n], prior->tnode_id) ||
                    ulab_streq(ids[n], prior->anode_id) ||
                    ulab_streq(ids[n], prior->cnode_id))
                    return fail(err, "duplicate physical node ID");
            }
        }
    }
    return ULAB_OK;
}

static int load_ues(json_t *root, hardware_config_t *config,
                     ulab_error_t *err) {
    static const char *const fields[] = {
        "ref", "site_ref", "slot", "iccid", "imsi", NULL
    };
    json_t *array;
    json_t *entry;
    size_t i;
    size_t j;

    array = json_object_get(root, "ues");
    if (!json_is_array(array) || json_array_size(array) > ULAB_HARDWARE_MAX_UES)
        return fail(err, "ues must be an array with at most 2 entries");
    config->ue_count = json_array_size(array);
    json_array_foreach(array, i, entry) {
        hardware_ue_t *ue;

        ue = &config->ues[i];
        if (keys(entry, fields, err) ||
            string_field(entry, "ref", ue->ref, sizeof(ue->ref), err) ||
            string_field(entry, "site_ref", ue->site_ref, sizeof(ue->site_ref), err) ||
            string_field(entry, "iccid", ue->iccid, sizeof(ue->iccid), err) ||
            string_field(entry, "imsi", ue->imsi, sizeof(ue->imsi), err) ||
            integer_field(json_object_get(entry, "slot"), 1,
                           ULAB_HARDWARE_MAX_UES, &ue->slot, "UE slot", err))
            return ULAB_ERR;
        if (!token(ue->ref) || find_site(config, ue->site_ref) == NULL)
            return fail(err, "invalid UE ref or unknown UE site_ref");
        if (!digits(ue->iccid, 22) || !digits(ue->imsi, 15))
            return fail(err, "ICCID/IMSI must be decimal strings (max 22/15 digits)");
        for (j = 0; j < i; j++) {
            const hardware_ue_t *prior;

            prior = &config->ues[j];
            if (ue->slot == prior->slot || ulab_streq(ue->ref, prior->ref) ||
                ulab_streq(ue->iccid, prior->iccid) ||
                ulab_streq(ue->imsi, prior->imsi))
                return fail(err, "duplicate UE ref, slot, ICCID or IMSI");
        }
    }
    return ULAB_OK;
}

int hardware_config_load(const char *path, hardware_config_t *config,
                          ulab_error_t *err) {
    static const char *const fields[] = {
        "version", "rack", "media", "sites", "ues", NULL
    };
    json_t *root;
    json_error_t error;
    uint32_t version;
    int rc;

    memset(config, 0, sizeof(*config));
    root = json_load_file(path, JSON_REJECT_DUPLICATES, &error);
    if (root == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "hardware config %.256s:%d: %.512s", path, error.line, error.text);
        return ULAB_ERR;
    }
    rc = keys(root, fields, err) ||
        integer_field(json_object_get(root, "version"), 1, 1, &version,
                       "version", err) ||
        load_rack(root, config, err) || load_media(root, config, err) ||
        load_sites(root, config, err) || load_ues(root, config, err);
    json_decref(root);
    if (rc) memset(config, 0, sizeof(*config));
    return rc ? ULAB_ERR : ULAB_OK;
}

int hardware_check_capacity(const scenario_t *scenario,
                             const hardware_config_t *config,
                             ulab_error_t *err) {
    uint64_t sites;
    uint64_t ues;

    sites = (uint64_t)scenario->world.networks * scenario->world.sites_per_network;
    if (sites > config->site_count) {
        snprintf(err->msg, sizeof(err->msg),
                 "hardware: scenario requires %" PRIu64 " sites; configured %zu",
                 sites, config->site_count);
        return ULAB_ERR;
    }
    /* Check sites first so the next multiplication cannot overflow. */
    ues = sites * scenario->world.ues_per_site;
    if (ues > config->ue_count) {
        snprintf(err->msg, sizeof(err->msg),
                 "hardware: scenario requires %" PRIu64 " UEs; configured %zu slots",
                 ues, config->ue_count);
        return ULAB_ERR;
    }
    if (sites > 0 && (scenario->world.tower_per_site > 1 ||
        scenario->world.amplifier_per_site > 1 ||
        scenario->world.controller_per_site > 1))
        return fail(err, "each physical site supports at most one node of each type");
    return ULAB_OK;
}

int hardware_bind_world(const hardware_config_t *config, world_t *world,
                         ulab_error_t *err) {
    size_t i;

    /* Validate every binding before changing any identity. Extra inventory is OK. */
    for (i = 0; i < world->site_count; i++) {
        if (find_site(config, world->sites[i].ref) == NULL) {
            snprintf(err->msg, sizeof(err->msg), "hardware: missing site mapping for %s",
                     world->sites[i].ref);
            return ULAB_ERR;
        }
    }
    for (i = 0; i < world->ue_count; i++) {
        const ue_t *ue;
        const hardware_ue_t *binding;

        ue = &world->ues[i];
        if (ue->site_ref[0] == '\0') continue; /* Backend-only SIM fixture. */
        binding = find_ue(config, ue->ref);
        if (binding == NULL || !ulab_streq(binding->site_ref, ue->site_ref)) {
            snprintf(err->msg, sizeof(err->msg),
                     "hardware: missing UE mapping for %s at %s", ue->ref, ue->site_ref);
            return ULAB_ERR;
        }
    }
    for (i = 0; i < world->node_count; i++) {
        const node_t *node;

        node = &world->nodes[i];
        if (find_site(config, node->site_ref) == NULL ||
            (!ulab_streq(node->type, ULAB_NODE_TOWER) &&
             !ulab_streq(node->type, ULAB_NODE_AMPLIFIER) &&
             !ulab_streq(node->type, ULAB_NODE_CONTROLLER)))
            return fail(err, "node has no physical site/type mapping");
    }
    for (i = 0; i < world->site_count; i++) {
        site_t *site;
        const hardware_site_t *binding;

        site = &world->sites[i];
        binding = find_site(config, site->ref);
        ulab_copy(site->tnode_id, sizeof(site->tnode_id), binding->tnode_id);
        ulab_copy(site->anode_id, sizeof(site->anode_id), binding->anode_id);
        ulab_copy(site->cnode_id, sizeof(site->cnode_id), binding->cnode_id);
    }
    for (i = 0; i < world->node_count; i++) {
        node_t *node;
        const hardware_site_t *binding;
        const char *id;

        node = &world->nodes[i];
        binding = find_site(config, node->site_ref);
        id = ulab_streq(node->type, ULAB_NODE_TOWER) ? binding->tnode_id :
            ulab_streq(node->type, ULAB_NODE_AMPLIFIER) ? binding->anode_id :
                                                        binding->cnode_id;
        ulab_copy(node->id, sizeof(node->id), id);
        ulab_copy(node->bff_id, sizeof(node->bff_id), id);
        node->physical = 1;
    }
    for (i = 0; i < world->ue_count; i++) {
        ue_t *ue;
        const hardware_ue_t *binding;

        ue = &world->ues[i];
        if (ue->site_ref[0] == '\0') continue;
        binding = find_ue(config, ue->ref);
        ulab_copy(ue->iccid, sizeof(ue->iccid), binding->iccid);
        ulab_copy(ue->imsi, sizeof(ue->imsi), binding->imsi);
        ue->hardware_slot = binding->slot;
        ue->ip[0] = '\0'; /* Assigned by the real EPC, never by EPCEMU. */
    }
    return ULAB_OK;
}

static int supported_checks(const check_spec_t *checks, size_t count, ulab_error_t *err) {
    size_t i;
    for (i = 0; i < count; i++) {
        if (checks[i].type == CHECK_TRAFFIC_BLOCKED)
            return fail(err, "traffic_blocked requires real EPC policy evidence; hardware observation is not available");
        if ((checks[i].type == CHECK_TRAFFIC_ALLOWED || checks[i].type == CHECK_TRAFFIC_UNAVAILABLE) &&
            checks[i].expected_used_mb > UINT32_MAX)
            return fail(err, "traffic check amount_mb exceeds the UE-gadget limit");
    }
    return ULAB_OK;
}

int hardware_execution_supported(const scenario_t *scenario, ulab_error_t *err) {
    size_t p;
    size_t i;
    for (p = 0; p <= scenario->phase_count; p++) {
        const check_spec_t *checks;
        size_t count;
        checks = p == scenario->phase_count ? scenario->final_checks : scenario->phases[p].checks;
        count = p == scenario->phase_count ? scenario->final_check_count : scenario->phases[p].check_count;
        if (supported_checks(checks, count, err)) return ULAB_ERR;
        if (p == scenario->phase_count) break;
        for (i = 0; i < scenario->phases[p].event_count; i++) {
            const event_spec_t *event = &scenario->phases[p].events[i];
            if (supported_checks(event->checks, event->check_count, err)) return ULAB_ERR;
            switch (event->type) {
            case EVT_DISCONNECT_NODES:
            case EVT_RECONNECT_NODES:
            case EVT_MARK_NODE_OFFLINE:
            case EVT_RESTORE_NODE:
            case EVT_FAILURE_CONTROL:
                snprintf(err->msg, sizeof(err->msg), "hardware: event %s is not supported by this runtime",
                         scenario_event_name(event->type));
                return ULAB_ERR;
            case EVT_CONFIGURE_SITES:
                if (ulab_streq(event->status, "hold")) return fail(err, "configure_sites hold needs virtual node fault injection");
                break;
            case EVT_TRAFFIC:
                if (event->amount_mb > UINT32_MAX) return fail(err, "traffic amount_mb exceeds the UE-gadget limit");
                break;
            default: break;
            }
        }
    }
    for (p = 0; p < scenario->profile_count; p++) {
        for (i = 0; i < scenario->profiles[p].bucket_count; i++) {
            if (scenario->profiles[p].buckets[i].amount_mb > UINT32_MAX)
                return fail(err, "profile amount_mb exceeds the UE-gadget limit");
        }
    }
    return ULAB_OK;
}
