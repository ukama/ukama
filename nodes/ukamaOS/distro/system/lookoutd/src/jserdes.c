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

static char *json_dup_string_or_null(JsonObj *json, const char *key) {

    JsonObj *entry = NULL;
    const char *value = NULL;

    if (json == NULL || key == NULL) {
        return NULL;
    }

    entry = json_object_get(json, key);
    if (!json_is_string(entry)) {
        return NULL;
    }

    value = json_string_value(entry);
    if (value == NULL) {
        return NULL;
    }

    return strdup(value);
}

static bool json_get_bool_default(JsonObj *json,
                                  const char *key,
                                  bool defVal) {

    JsonObj *entry = NULL;

    if (json == NULL || key == NULL) {
        return defVal;
    }

    entry = json_object_get(json, key);
    if (!json_is_boolean(entry)) {
        return defVal;
    }

    return json_boolean_value(entry) ? true : false;
}

static int json_get_int_default(JsonObj *json,
                                const char *key,
                                int defVal) {

    JsonObj *entry = NULL;

    if (json == NULL || key == NULL) {
        return defVal;
    }

    entry = json_object_get(json, key);
    if (!json_is_integer(entry)) {
        return defVal;
    }

    return (int)json_integer_value(entry);
}

bool json_deserialize_starter_status(StarterStatusData *starter,
                                     JsonObj *json) {

    JsonObj *jStarter = NULL;

    if (starter == NULL || json == NULL) {
        return USYS_FALSE;
    }

    memset(starter, 0, sizeof(StarterStatusData));

    jStarter = json_object_get(json, "starterd");
    if (!json_is_object(jStarter)) {
        starter->available = USYS_FALSE;
        return USYS_FALSE;
    }

    starter->available = USYS_TRUE;
    starter->state = json_dup_string_or_null(jStarter, "state");
    starter->updateInProgress =
        json_get_bool_default(jStarter, "updateInProgress", false);
    starter->switchRequested =
        json_get_bool_default(jStarter, "switchRequested", false);
    starter->terminateRequested =
        json_get_bool_default(jStarter, "terminateRequested", false);
    starter->exitCode =
        json_get_int_default(jStarter, "exitCode", 0);

    return USYS_TRUE;
}

static const char *node_type_to_string(LookoutNodeType nodeType) {

    switch (nodeType) {
    case LOOKOUT_NODE_TOWER:
        return "tower";

    case LOOKOUT_NODE_AMPLIFIER:
        return "amplifier";

    case LOOKOUT_NODE_CONTROL:
        return "control";

    default:
        return "unknown";
    }
}

static void json_add_capabilities(JsonObj *root, LookoutNodeType nodeType) {

    JsonObj *array = NULL;

    if (root == NULL) {
        return;
    }

    array = json_array();
    if (array == NULL) {
        return;
    }

    switch (nodeType) {
    case LOOKOUT_NODE_TOWER:
        json_array_append_new(array, json_string("power"));
        json_array_append_new(array, json_string("cellular"));
        json_array_append_new(array, json_string("radio"));
        json_array_append_new(array, json_string("gps"));
        json_array_append_new(array, json_string("backhaul"));
        break;

    case LOOKOUT_NODE_AMPLIFIER:
        json_array_append_new(array, json_string("power"));
        json_array_append_new(array, json_string("radio"));
        json_array_append_new(array, json_string("fem"));
        break;

    case LOOKOUT_NODE_CONTROL:
        json_array_append_new(array, json_string("power"));
        json_array_append_new(array, json_string("switch"));
        json_array_append_new(array, json_string("controller"));
        json_array_append_new(array, json_string("backhaul"));
        break;

    default:
        break;
    }

    json_object_set_new(root, "capabilities", array);
}

static long get_uptime_sec(void) {

    FILE *fp = NULL;
    double uptime = 0.0;

    fp = fopen("/proc/uptime", "r");
    if (fp == NULL) {
        return 0;
    }

    if (fscanf(fp, "%lf", &uptime) != 1) {
        fclose(fp);
        return 0;
    }

    fclose(fp);
    return (long)uptime;
}

static void json_add_app_resources(JsonObj *jApp, CappRuntime *runtime) {

    JsonObj *resources = NULL;

    if (jApp == NULL || runtime == NULL) {
        return;
    }

    resources = json_object();
    if (resources == NULL) {
        return;
    }

    json_object_set_new(resources,
                        "cpuPercent",
                        json_real(runtime->cpu));
    json_object_set_new(resources,
                        "memoryRssKb",
                        json_integer(runtime->memory));
    json_object_set_new(resources,
                        "diskReadBytes",
                        json_integer(0));
    json_object_set_new(resources,
                        "diskWriteBytes",
                        json_integer(runtime->disk));

    json_object_set_new(jApp, "resources", resources);
}

static void json_add_apps(JsonObj *root, CappList *list) {

    JsonObj *apps = NULL;
    JsonObj *jApp = NULL;
    CappList *ptr = NULL;
    Capp *capp = NULL;

    if (root == NULL) {
        return;
    }

    apps = json_array();
    if (apps == NULL) {
        return;
    }

    for (ptr = list; ptr; ptr = ptr->next) {
        capp = ptr->capp;
        if (capp == NULL || capp->runtime == NULL) {
            continue;
        }

        jApp = json_object();
        if (jApp == NULL) {
            continue;
        }

        json_object_set_new(jApp,
                            "space",
                            json_string(capp->space ? capp->space : ""));
        json_object_set_new(jApp,
                            "name",
                            json_string(capp->name ? capp->name : ""));
        json_object_set_new(jApp,
                            "tag",
                            json_string(capp->tag ? capp->tag : ""));
        json_object_set_new(jApp,
                            "version",
                            json_string(capp->tag ? capp->tag : ""));
        json_object_set_new(jApp,
                            "state",
                            json_string(capp->runtime->status ?
                                        capp->runtime->status : "unknown"));
        json_object_set_new(jApp,
                            "pid",
                            json_integer(capp->runtime->pid));

        json_add_app_resources(jApp, capp->runtime);
        json_array_append_new(apps, jApp);
    }

    json_object_set_new(root, "apps", apps);
}

static void json_add_starter(JsonObj *system, StarterStatusData *starter) {

    JsonObj *jStarter = NULL;

    if (system == NULL || starter == NULL) {
        return;
    }

    jStarter = json_object();
    if (jStarter == NULL) {
        return;
    }

    json_object_set_new(jStarter,
                        "available",
                        json_boolean(starter->available));
    json_object_set_new(jStarter,
                        "state",
                        json_string(starter->state ?
                                    starter->state : "unknown"));
    json_object_set_new(jStarter,
                        "updateInProgress",
                        json_boolean(starter->updateInProgress));
    json_object_set_new(jStarter,
                        "switchRequested",
                        json_boolean(starter->switchRequested));
    json_object_set_new(jStarter,
                        "terminateRequested",
                        json_boolean(starter->terminateRequested));
    json_object_set_new(jStarter,
                        "exitCode",
                        json_integer(starter->exitCode));

    json_object_set_new(system, "starter", jStarter);
}

static void json_add_power(JsonObj *system, PowerStatusData *power) {

    JsonObj *jPower = NULL;

    if (system == NULL || power == NULL) {
        return;
    }

    jPower = json_object();
    if (jPower == NULL) {
        return;
    }

    json_object_set_new(jPower,
                        "available",
                        json_boolean(power->available));

    if (power->available == USYS_TRUE) {
        json_object_set_new(jPower, "ok", json_boolean(power->ok));
        json_object_set_new(jPower,
                            "board",
                            json_string(power->board ? power->board : ""));
        json_object_set_new(jPower,
                            "reason",
                            json_string(power->reason ? power->reason : ""));
        json_object_set_new(jPower,
                            "totalWatts",
                            json_real(power->totalWatts));
        json_object_set_new(jPower,
                            "temperatureC",
                            json_real(power->temperatureC));
    } else {
        json_object_set_new(jPower,
                            "error",
                            json_string("power_service_unreachable"));
    }

    json_object_set_new(system, "power", jPower);
}

static void json_add_system(JsonObj *root, LookoutStatusData *status) {

    JsonObj *system = NULL;

    if (root == NULL || status == NULL) {
        return;
    }

    system = json_object();
    if (system == NULL) {
        return;
    }

    json_object_set_new(system,
                        "uptimeSec",
                        json_integer(get_uptime_sec()));

    json_add_starter(system, &status->starter);
    json_add_power(system, &status->power);

    json_object_set_new(root, "system", system);
}

static void json_add_cellular(JsonObj *interfaces,
                              CellularStatusData *cellular) {

    JsonObj *jCellular = NULL;

    if (interfaces == NULL || cellular == NULL ||
        cellular->available != USYS_TRUE) {
        return;
    }

    jCellular = json_object();
    if (jCellular == NULL) {
        return;
    }

    json_object_set_new(jCellular,
                        "available",
                        cellular->service ? json_true() : json_false());
    if (cellular->service) {
        json_object_set_new(jCellular,
                            "service",
                            json_string(cellular->service));
    }
    if (cellular->error) {
        json_object_set_new(jCellular,
                            "error",
                            json_string(cellular->error));
    }

    json_object_set_new(interfaces, "cellular", jCellular);
}

static void json_add_radio(JsonObj *interfaces,
                           RadioStatusData *radio) {

    JsonObj *jRadio = NULL;

    if (interfaces == NULL || radio == NULL ||
        radio->available != USYS_TRUE) {
        return;
    }

    jRadio = json_object();
    if (jRadio == NULL) {
        return;
    }

    json_object_set_new(jRadio, "available", json_true());
    json_object_set_new(jRadio,
                        "state",
                        json_string(radio->state ? radio->state : "unknown"));

    json_object_set_new(interfaces, "radio", jRadio);
}

static void json_add_gps(JsonObj *interfaces,
                         GPSClientData *gps) {

    JsonObj *jGps = NULL;

    if (interfaces == NULL || gps == NULL) {
        return;
    }

    if (gps->available != USYS_TRUE) {
        return;
    }

    jGps = json_object();
    if (jGps == NULL) {
        return;
    }

    json_object_set_new(jGps, "available", json_true());
    json_object_set_new(jGps, "lock", json_boolean(gps->lock));
    json_object_set_new(jGps,
                        "coordinates",
                        json_string(gps->coordinates ?
                                    gps->coordinates :
                                    LOOKOUT_GPS_COORD_NA));
    json_object_set_new(jGps,
                        "time",
                        json_string(gps->time ?
                                    gps->time :
                                    LOOKOUT_GPS_TIME_NA));

    json_object_set_new(interfaces, "gps", jGps);
}

static JsonObj *copy_string_field(JsonObj *src, const char *key) {

    JsonObj *entry = NULL;

    if (src == NULL || key == NULL) {
        return NULL;
    }

    entry = json_object_get(src, key);
    if (!json_is_string(entry)) {
        return NULL;
    }

    return json_string(json_string_value(entry));
}

static JsonObj *copy_bool_field(JsonObj *src, const char *key) {

    JsonObj *entry = NULL;

    if (src == NULL || key == NULL) {
        return NULL;
    }

    entry = json_object_get(src, key);
    if (!json_is_boolean(entry)) {
        return NULL;
    }

    return json_boolean(json_boolean_value(entry));
}

static JsonObj *copy_int_field(JsonObj *src, const char *key) {

    JsonObj *entry = NULL;

    if (src == NULL || key == NULL) {
        return NULL;
    }

    entry = json_object_get(src, key);
    if (!json_is_integer(entry)) {
        return NULL;
    }

    return json_integer(json_integer_value(entry));
}

static JsonObj *copy_num_field(JsonObj *src, const char *key) {

    JsonObj *entry = NULL;

    if (src == NULL || key == NULL) {
        return NULL;
    }

    entry = json_object_get(src, key);
    if (!json_is_number(entry)) {
        return NULL;
    }

    return json_real(json_number_value(entry));
}

static void set_if(JsonObj *dst,
                   const char *key,
                   JsonObj *value) {

    if (dst == NULL || key == NULL || value == NULL) {
        return;
    }

    json_object_set_new(dst, key, value);
}

static void json_add_switch_ports(JsonObj *jSwitch, JsonObj *portsRaw) {

    JsonObj *ports = NULL;
    JsonObj *srcArray = NULL;
    JsonObj *src = NULL;
    JsonObj *dst = NULL;

    int i = 0;
    int count = 0;

    if (jSwitch == NULL || portsRaw == NULL) {
        return;
    }

    srcArray = portsRaw;
    if (json_is_object(portsRaw)) {
        srcArray = json_object_get(portsRaw, "ports");
    }

    if (!json_is_array(srcArray)) {
        return;
    }

    ports = json_array();
    if (ports == NULL) {
        return;
    }

    count = json_array_size(srcArray);
    for (i = 0; i < count; i++) {
        src = json_array_get(srcArray, i);
        if (!json_is_object(src)) {
            continue;
        }

        dst = json_object();
        if (dst == NULL) {
            continue;
        }

        set_if(dst, "id", copy_int_field(src, "id"));
        set_if(dst, "name", copy_string_field(src, "name"));
        set_if(dst, "present", copy_bool_field(src, "present"));
        set_if(dst, "adminState", copy_string_field(src, "adminState"));
        set_if(dst, "linkState", copy_string_field(src, "linkState"));
        set_if(dst, "poeState", copy_string_field(src, "poeState"));
        set_if(dst, "poeOperational",
               copy_bool_field(src, "poeOperational"));
        set_if(dst, "speedBps", copy_int_field(src, "speedBps"));
        set_if(dst, "powerWatts", copy_num_field(src, "powerWatts"));
        set_if(dst, "fault", copy_string_field(src, "fault"));

        json_array_append_new(ports, dst);
    }

    json_object_set_new(jSwitch, "ports", ports);
}

static void json_add_switch(JsonObj *interfaces,
                            SwitchStatusData *sw) {

    JsonObj *jSwitch = NULL;
    JsonObj *jSwitchd = NULL;
    JsonObj *jInfo = NULL;
    JsonObj *policy = NULL;

    if (interfaces == NULL || sw == NULL || sw->available != USYS_TRUE) {
        return;
    }

    jSwitch = json_object();
    if (jSwitch == NULL) {
        return;
    }

    if (sw->status == NULL) {
        json_object_set_new(jSwitch, "available", json_false());
        json_object_set_new(jSwitch,
                            "error",
                            json_string("switch_service_unreachable"));
        json_object_set_new(interfaces, "switch", jSwitch);
        return;
    }

    json_object_set_new(jSwitch, "available", json_true());

    jSwitchd = json_object_get(sw->status, "switchd");
    jInfo = json_object_get(sw->status, "switch");

    if (json_is_object(jSwitchd)) {
        set_if(jSwitch, "state", copy_string_field(jSwitchd, "state"));
        set_if(jSwitch,
               "reachable",
               copy_bool_field(jSwitchd, "reachable"));
    }

    if (json_is_object(jInfo)) {
        set_if(jSwitch, "model", copy_string_field(jInfo, "model"));
        set_if(jSwitch,
               "softwareVersion",
               copy_string_field(jInfo, "softwareVersion"));
        set_if(jSwitch, "portCount", copy_int_field(jInfo, "portCount"));
    }

    policy = json_object();
    if (policy) {
        if (sw->policy) {
            set_if(policy, "state", copy_string_field(sw->policy, "state"));
            set_if(policy, "hash", copy_string_field(sw->policy, "hash"));
            set_if(policy, "source", copy_string_field(sw->policy, "source"));
            set_if(policy, "error", copy_string_field(sw->policy, "error"));
        }

        json_object_set_new(jSwitch, "policy", policy);
    }

    json_add_switch_ports(jSwitch, sw->ports);

    json_object_set_new(interfaces, "switch", jSwitch);
}

static void json_add_controller(JsonObj *interfaces,
                                ControllerStatusData *controller) {

    JsonObj *jController = NULL;
    JsonObj *src = NULL;
    JsonObj *solar = NULL;
    JsonObj *battery = NULL;
    JsonObj *load = NULL;
    JsonObj *ctrl = NULL;

    if (interfaces == NULL || controller == NULL ||
        controller->available != USYS_TRUE) {
        return;
    }

    jController = json_object();
    if (jController == NULL) {
        return;
    }

    src = controller->status;
    if (src == NULL) {
        json_object_set_new(jController, "available", json_false());
        json_object_set_new(jController,
                            "error",
                            json_string("controller_service_unreachable"));
        json_object_set_new(interfaces, "controller", jController);
        return;
    }

    json_object_set_new(jController, "available", json_true());

    set_if(jController, "commOk", copy_bool_field(src, "commOk"));
    set_if(jController, "chargeState", copy_string_field(src, "chargeState"));
    set_if(jController, "errorCode", copy_int_field(src, "errorCode"));
    set_if(jController, "error", copy_string_field(src, "error"));
    set_if(jController,
           "activeAlarmCount",
           copy_int_field(src, "activeAlarmCount"));

    solar = json_object_get(src, "solar");
    if (json_is_object(solar)) {
        json_object_set(jController, "solar", solar);
    }

    battery = json_object_get(src, "battery");
    if (json_is_object(battery)) {
        json_object_set(jController, "battery", battery);
    }

    ctrl = json_object_get(src, "controller");
    if (json_is_object(ctrl)) {
        load = json_object();
        if (load) {
            set_if(load, "outputOn",
                   copy_bool_field(ctrl, "loadOutputOn"));
            set_if(load, "currentA",
                   copy_num_field(ctrl, "loadCurrentA"));
            json_object_set_new(jController, "load", load);
        }
    }

    json_object_set_new(interfaces, "controller", jController);
}

static void json_add_backhaul(JsonObj *interfaces,
                              BackhaulStatusData *backhaul) {

    JsonObj *jBackhaul = NULL;
    JsonObj *src = NULL;

    if (interfaces == NULL || backhaul == NULL ||
        backhaul->available != USYS_TRUE) {
        return;
    }

    jBackhaul = json_object();
    if (jBackhaul == NULL) {
        return;
    }

    src = backhaul->status;
    if (src == NULL) {
        json_object_set_new(jBackhaul, "available", json_false());
        json_object_set_new(jBackhaul,
                            "error",
                            json_string("backhaul_service_unreachable"));
        json_object_set_new(interfaces, "backhaul", jBackhaul);
        return;
    }

    json_object_set_new(jBackhaul, "available", json_true());
    set_if(jBackhaul, "state", copy_string_field(src, "backhaulState"));
    set_if(jBackhaul, "linkGuess", copy_string_field(src, "linkGuess"));
    set_if(jBackhaul, "confidence", copy_num_field(src, "confidence"));

    json_object_set_new(interfaces, "backhaul", jBackhaul);
}

static void json_add_fem(JsonObj *interfaces,
                         FemStatusData *fem) {

    JsonObj *jFem = NULL;
    JsonObj *fems = NULL;

    if (interfaces == NULL || fem == NULL ||
        fem->available != USYS_TRUE) {
        return;
    }

    jFem = json_object();
    if (jFem == NULL) {
        return;
    }

    if (fem->status == NULL) {
        json_object_set_new(jFem, "available", json_false());
        json_object_set_new(jFem,
                            "error",
                            json_string("fem_service_unreachable"));
        json_object_set_new(interfaces, "fem", jFem);
        return;
    }

    json_object_set_new(jFem, "available", json_true());

    fems = json_object_get(fem->status, "fems");
    if (json_is_array(fems)) {
        json_object_set(jFem, "fems", fems);
    }

    json_object_set_new(interfaces, "fem", jFem);
}

static void json_add_interfaces(JsonObj *root,
                                LookoutStatusData *status) {

    JsonObj *interfaces = NULL;

    if (root == NULL || status == NULL) {
        return;
    }

    interfaces = json_object();
    if (interfaces == NULL) {
        return;
    }

    json_add_cellular(interfaces, &status->cellular);
    json_add_radio(interfaces, &status->radio);
    json_add_gps(interfaces, &status->gps);
    json_add_switch(interfaces, &status->sw);
    json_add_controller(interfaces, &status->controller);
    json_add_backhaul(interfaces, &status->backhaul);
    json_add_fem(interfaces, &status->fem);

    json_object_set_new(root, "interfaces", interfaces);
}

bool json_serialize_health_report(JsonObj **json,
                                  Config *config,
                                  CappList *list,
                                  LookoutStatusData *status) {

    char buffer[MAX_BUFFER] = {0};
    time_t currTime;

    if (json == NULL || config == NULL || status == NULL) {
        return USYS_FALSE;
    }

    *json = json_object();
    if (*json == NULL) {
        return USYS_FALSE;
    }

    json_object_set_new(*json,
                        "schemaVersion",
                        json_string(LOOKOUT_SCHEMA_VERSION));
    json_object_set_new(*json,
                        "nodeId",
                        json_string(config->nodeID ? config->nodeID : ""));
    json_object_set_new(*json,
                        "nodeType",
                        json_string(node_type_to_string(config->nodeType)));

    time(&currTime);
    snprintf(buffer, sizeof(buffer), "%ld", (long)currTime);
    json_object_set_new(*json, "reportedAt", json_string(buffer));

    json_add_capabilities(*json, config->nodeType);
    json_add_system(*json, status);
    json_add_interfaces(*json, status);
    json_add_apps(*json, list);
    json_object_set_new(*json, "events", json_array());

    return USYS_TRUE;
}

void json_free(JsonObj** json) {

    if (*json) {
        json_decref(*json);
        *json = NULL;
    }
}
