/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <jansson.h>
#include <string.h>

#include "rlogd.h"
#include "usys_log.h"

extern ThreadData *gData;

static json_t *legacy_record(const char *log) {
    json_t *record;

    record = json_object();
    if (!record) {
        return NULL;
    }

    json_object_set_new(record, "schema", json_string("ukama.log.v1"));
    json_object_set_new(record, "level", json_string("info"));
    json_object_set_new(record, "app", json_string("legacy"));
    json_object_set_new(record, "component", json_string("legacy"));
    json_object_set_new(record, "event", json_string("legacy_log"));
    json_object_set_new(record, "msg", json_string(log ? log : ""));
    json_object_set_new(record, "structured", json_false());

    return record;
}

void process_logs(void *nodeID, const char *log) {
    json_error_t error;
    json_t *record;

    (void)nodeID;

    if (!gData || !gData->store || !log) {
        return;
    }

    record = json_loads(log, 0, &error);
    if (!record || !json_is_object(record)) {
        if (record) {
            json_decref(record);
        }
        record = legacy_record(log);
    }

    if (!record) {
        return;
    }

    if (log_store_append(gData->store, record, NULL) != 0) {
        usys_log_error("Unable to append record to canonical log store");
    }

    json_decref(record);
}
