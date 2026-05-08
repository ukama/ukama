/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <time.h>

#include "jserdes.h"
#include "lookout.h"
#include "errorcode.h"
#include "json_types.h"
#include "web_client.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"

/* defined in resources.c */
extern char *get_radio_status(void);

void json_log(json_t *json) {

    char *str = NULL;

    str = json_dumps(json, 0);
    if (str) {
        log_debug("json str: %s", str);
        free(str);
    }
}

static bool get_json_entry(json_t *json,
                           char *key,
                           json_type type,
                           char **strValue,
                           int *intValue,
                           double *doubleValue) {

    json_t *jEntry = NULL;

    if (json == NULL || key == NULL) return USYS_FALSE;

    jEntry = json_object_get(json, key);
    if (jEntry == NULL) {
        log_error("Missing %s key in json", key);
        return USYS_FALSE;
    }

    switch(type) {
    case (JSON_STRING):
        *strValue = strdup(json_string_value(jEntry));
        break;

    case (JSON_INTEGER):
        *intValue = json_integer_value(jEntry);
        break;

    case (JSON_REAL):
        *doubleValue = json_real_value(jEntry);
        break;

    default:
        log_error("Invalid type for json key-value: %d", type);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

bool json_deserialize_node_id(char **nodeID, JsonObj *json) {

    JsonObj *jNodeInfo = NULL;

    if (json == NULL) return USYS_FALSE;

    jNodeInfo = json_object_get(json, JTAG_NODE_INFO);
    if (jNodeInfo == NULL) {
        log_error("Missing mandatory %s from JSON", JTAG_NODE_INFO);
        return USYS_FALSE;
    }

    if (get_json_entry(jNodeInfo,
                       JTAG_UUID,
                       JSON_STRING,
                       nodeID,
                       NULL,
                       NULL) == USYS_FALSE) {
        log_error("Error deserializing node info");
        json_log(json);
        *nodeID = NULL;

        return USYS_FALSE;
    }

    return USYS_TRUE;
}

static const char *starter_state_to_ukama_status(const char *state) {

    if (state == NULL) {
        return "Unknown";
    }

    if (strcmp(state, "running") == 0) {
        return "Active";
    }

    if (strcmp(state, "starting") == 0 ||
        strcmp(state, "pending") == 0 ||
        strcmp(state, "switching") == 0 ||
        strcmp(state, "installing") == 0) {
        return "Pending";
    }

    if (strcmp(state, "stopped") == 0 ||
        strcmp(state, "exited") == 0 ||
        strcmp(state, "failed") == 0 ||
        strcmp(state, "fatal") == 0) {
        return "Failure";
    }

    return "Unknown";
}

static void add_starter_app(CappList **cappList,
                            const char *spaceName,
                            JsonObj *jApp) {

    JsonObj *jName = NULL;
    JsonObj *jTag = NULL;
    JsonObj *jState = NULL;
    JsonObj *jPid = NULL;

    const char *name = NULL;
    const char *tag = NULL;
    const char *state = NULL;
    const char *status = NULL;
    int pid = 0;

    if (cappList == NULL || spaceName == NULL || jApp == NULL) {
        return;
    }

    jName  = json_object_get(jApp, JTAG_NAME);
    jTag   = json_object_get(jApp, JTAG_TAG);
    jState = json_object_get(jApp, "state");
    jPid   = json_object_get(jApp, JTAG_PID);

    if (!json_is_string(jName)) {
        return;
    }

    name = json_string_value(jName);

    if (json_is_string(jTag)) {
        tag = json_string_value(jTag);
    } else {
        tag = "latest";
    }

    if (json_is_string(jState)) {
        state = json_string_value(jState);
    } else {
        state = "unknown";
    }

    if (json_is_integer(jPid)) {
        pid = (int)json_integer_value(jPid);
    }

    status = starter_state_to_ukama_status(state);

    add_capp_to_list(cappList,
                     spaceName,
                     name,
                     tag,
                     status,
                     pid);
}

/*
 * Latest starter.d status:
 *
 * {
 *   "spaces": [
 *     {
 *       "name": "boot",
 *       "apps": [
 *         {
 *           "name": "noded",
 *           "tag": "old_arch-7bd935d06",
 *           "state": "running",
 *           "pid": 11
 *         }
 *       ]
 *     }
 *   ],
 *   "starterd": {
 *     "exitCode": 0,
 *     "switchRequested": false,
 *     "updateInProgress": false
 *   }
 * }
 */
bool json_deserialize_capps(CappList **cappList, JsonObj *json) {

    int i = 0;
    int j = 0;
    int spaceCount = 0;
    int appCount = 0;
    int totalApps = 0;

    JsonObj *jSpaces = NULL;
    JsonObj *jSpace = NULL;
    JsonObj *jSpaceName = NULL;
    JsonObj *jApps = NULL;
    JsonObj *jApp = NULL;

    const char *spaceName = NULL;

    if (json == NULL || cappList == NULL) {
        return USYS_FALSE;
    }

    jSpaces = json_object_get(json, "spaces");
    if (!json_is_array(jSpaces)) {
        usys_log_error("starter.d status missing spaces[]");
        return USYS_FALSE;
    }

    spaceCount = json_array_size(jSpaces);

    for (i = 0; i < spaceCount; i++) {
        jSpace = json_array_get(jSpaces, i);
        if (jSpace == NULL) continue;

        jSpaceName = json_object_get(jSpace, JTAG_NAME);
        jApps = json_object_get(jSpace, "apps");

        if (!json_is_string(jSpaceName) || !json_is_array(jApps)) {
            continue;
        }

        spaceName = json_string_value(jSpaceName);
        appCount = json_array_size(jApps);

        for (j = 0; j < appCount; j++) {
            jApp = json_array_get(jApps, j);
            add_starter_app(cappList, spaceName, jApp);
            totalApps++;
        }
    }

    usys_log_debug("Received %d capps from starter.d", totalApps);

    return USYS_TRUE;
}

static void json_add_resources_to_capp_report(JsonObj **json,
                                              CappRuntime *runtime) {

    JsonObj *jArray = NULL;
    JsonObj *jMemory = NULL;
    JsonObj *jDisk = NULL;
    JsonObj *jCpu = NULL;

    char buffer[MAX_BUFFER] = {0};

    json_object_set_new(*json, JTAG_RESOURCES, json_array());
    jArray = json_object_get(*json, JTAG_RESOURCES);
    if (jArray == NULL) return;

    jMemory = json_object();
    jDisk   = json_object();
    jCpu    = json_object();

    sprintf(buffer, "%d", runtime->memory);
    json_object_set_new(jMemory,
                        JTAG_NAME,
                        json_string("memory"));
    json_object_set_new(jMemory,
                        JTAG_VALUE,
                        json_string(buffer));

    sprintf(buffer, "%d", runtime->disk);
    json_object_set_new(jDisk,
                        JTAG_NAME,
                        json_string("disk"));
    json_object_set_new(jDisk,
                        JTAG_VALUE,
                        json_string(buffer));

    sprintf(buffer, "%f", runtime->cpu);
    json_object_set_new(jCpu,
                        JTAG_NAME,
                        json_string("cpu"));
    json_object_set_new(jCpu,
                        JTAG_VALUE,
                        json_string(buffer));

    json_array_append_new(jArray, jMemory);
    json_array_append_new(jArray, jDisk);
    json_array_append_new(jArray, jCpu);
}

static void json_add_system_info_to_report(JsonObj **json,
                                           bool radioAvailable) {

    JsonObj *jArray = NULL;
    JsonObj *jRadio = NULL;

    jRadio = json_object();

    json_object_set_new(*json, JTAG_SYSTEM, json_array());
    jArray = json_object_get(*json, JTAG_SYSTEM);
    if (jArray == NULL) return;

    json_object_set_new(jRadio,
                        JTAG_NAME,
                        json_string("radio"));

    if (radioAvailable == USYS_TRUE) {
        json_object_set_new(jRadio,
                            JTAG_VALUE,
                            json_string(get_radio_status()));
    } else {
        json_object_set_new(jRadio,
                            JTAG_VALUE,
                            json_string(LOOKOUT_STATUS_NA));
    }

    json_array_append_new(jArray, jRadio);
}

static void json_add_gps_info_to_report(JsonObj **json, GPSClientData *gps) {

    JsonObj *jArray = NULL;
    JsonObj *jEntry = NULL;

    const char *lockStr = LOOKOUT_STATUS_NA;
    const char *coord = LOOKOUT_GPS_COORD_NA;
    const char *timeStr = LOOKOUT_GPS_TIME_NA;

    if (json == NULL || *json == NULL || gps == NULL) return;

    jArray = json_object_get(*json, JTAG_SYSTEM);
    if (jArray == NULL || !json_is_array(jArray)) return;

    if (gps->available == USYS_TRUE) {
        lockStr = gps->gpsLock ? "true" : "false";

        if (gps->coordinates && gps->coordinates[0] != '\0') {
            coord = gps->coordinates;
        }

        if (gps->gpsTime && gps->gpsTime[0] != '\0') {
            timeStr = gps->gpsTime;
        }
    }

    jEntry = json_object();
    json_object_set_new(jEntry, JTAG_NAME, json_string("coordinates"));
    json_object_set_new(jEntry, JTAG_VALUE, json_string(coord));
    json_array_append_new(jArray, jEntry);

    jEntry = json_object();
    json_object_set_new(jEntry, JTAG_NAME, json_string("gpsLock"));
    json_object_set_new(jEntry, JTAG_VALUE, json_string(lockStr));
    json_array_append_new(jArray, jEntry);

    jEntry = json_object();
    json_object_set_new(jEntry, JTAG_NAME, json_string("gpsTime"));
    json_object_set_new(jEntry, JTAG_VALUE, json_string(timeStr));
    json_array_append_new(jArray, jEntry);
}

static void json_add_system_entry(JsonObj *jArray,
                                  const char *name,
                                  const char *value) {

    JsonObj *jEntry = NULL;

    if (jArray == NULL || !json_is_array(jArray) ||
        name == NULL || value == NULL) {
        return;
    }

    jEntry = json_object();
    if (jEntry == NULL) {
        return;
    }

    json_object_set_new(jEntry, JTAG_NAME, json_string(name));
    json_object_set_new(jEntry, JTAG_VALUE, json_string(value));
    json_array_append_new(jArray, jEntry);
}

static void json_add_switch_policy_info_to_report(JsonObj **json,
                                                  SwitchPolicyStatusData *sp) {

    JsonObj *jArray = NULL;

    const char *switchValue = LOOKOUT_STATUS_NA;
    const char *siteID      = LOOKOUT_STATUS_NA;
    const char *state       = LOOKOUT_STATUS_NA;
    const char *hash        = LOOKOUT_STATUS_NA;
    const char *source      = LOOKOUT_STATUS_NA;
    const char *error       = "";

    if (json == NULL || *json == NULL || sp == NULL) return;

    if (sp->available != USYS_TRUE) {
        return;
    }

    jArray = json_object_get(*json, JTAG_SYSTEM);
    if (jArray == NULL || !json_is_array(jArray)) return;

    if (sp->switchAvailable == USYS_TRUE) {
        switchValue = "available";

        if (sp->siteID && sp->siteID[0] != '\0') {
            siteID = sp->siteID;
        }

        if (sp->policyState && sp->policyState[0] != '\0') {
            state = sp->policyState;
        }

        if (sp->policyHash && sp->policyHash[0] != '\0') {
            hash = sp->policyHash;
        }

        if (sp->policySource && sp->policySource[0] != '\0') {
            source = sp->policySource;
        }

        if (sp->policyError && sp->policyError[0] != '\0') {
            error = sp->policyError;
        }
    } else {
        switchValue = "unavailable";
        state = "unknown";
        hash = LOOKOUT_STATUS_NA;
        source = LOOKOUT_STATUS_NA;
        error = "switchd_unreachable";
    }

    json_add_system_entry(jArray, "switch", switchValue);
    json_add_system_entry(jArray, "switchPolicySiteID", siteID);
    json_add_system_entry(jArray, "switchPolicy", state);
    json_add_system_entry(jArray, "switchPolicyHash", hash);
    json_add_system_entry(jArray, "switchPolicySource", source);
    json_add_system_entry(jArray, "switchPolicyError", error);
}

bool json_serialize_health_report(JsonObj **json,
                                  char *nodeID,
                                  CappList *list,
                                  GPSClientData *gps,
                                  SwitchPolicyStatusData *switchPolicy,
                                  bool radioAvailable) {

    JsonObj *jArray = NULL;
    JsonObj *jCapp = NULL;
    CappList *ptr = NULL;
    Capp *capp = NULL;

    char buffer[MAX_BUFFER] = {0};
    time_t currTime;

    *json = json_object();
    if (*json == NULL) return USYS_FALSE;

    json_object_set_new(*json,
                        JTAG_NODE_ID,
                        json_string(nodeID ? nodeID : ""));

    time(&currTime);
    sprintf(buffer, "%ld", (long)currTime);
    json_object_set_new(*json,
                        JTAG_TIMESTAMP,
                        json_string(buffer));

    json_object_set_new(*json, JTAG_CAPPS, json_array());
    jArray = json_object_get(*json, JTAG_CAPPS);
    if (jArray == NULL) {
        json_decref(*json);
        *json = NULL;
        return USYS_FALSE;
    }

    for (ptr = list; ptr; ptr = ptr->next) {
        capp = ptr->capp;
        if (capp == NULL || capp->runtime == NULL) {
            continue;
        }

        jCapp = json_object();
        if (jCapp == NULL) {
            json_decref(*json);
            *json = NULL;
            return USYS_FALSE;
        }

        json_object_set_new(jCapp,
                            JTAG_SPACE,
                            json_string(capp->space ? capp->space : ""));
        json_object_set_new(jCapp,
                            JTAG_NAME,
                            json_string(capp->name ? capp->name : ""));
        json_object_set_new(jCapp,
                            JTAG_TAG,
                            json_string(capp->tag ? capp->tag : ""));
        json_object_set_new(jCapp,
                            JTAG_STATUS,
                            json_string(capp->runtime->status ?
                                        capp->runtime->status : "Unknown"));

        json_add_resources_to_capp_report(&jCapp, capp->runtime);
        json_array_append_new(jArray, jCapp);
    }

    json_add_system_info_to_report(json, radioAvailable);
    json_add_gps_info_to_report(json, gps);
    json_add_switch_policy_info_to_report(json, switchPolicy);

    return USYS_TRUE;
}

void json_free(JsonObj** json) {

    if (*json) {
        json_decref(*json);
        *json = NULL;
    }
}
