/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_JSON_TYPES_H_
#define INC_JSON_TYPES_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "notify_macros.h"
#include "usys_types.h"

#include "jansson.h"

/* JSON types used by parser */
typedef  json_t  JsonObj;
typedef  json_error_t JsonErrObj;

#define MAX_JSON_SCHEMA     5

/* MFG Data Json Tags */
#define JTAG_HEADER                     "header"
#define JTAG_VERSION                    "version"
#define JTAG_MAJOR_VERSION              "major"
#define JTAG_MINOR_VERSION              "minor"
#define JTAG_IDX_TABLE_OFFSET           "indexTableOffset"
#define JTAG_IDX_TUPLE_SIZE             "indexTupleSize"
#define JTAG_IDX_TUPLE_MAX_COUNT        "indexTupleMaxCount"
#define JTAG_IDX_CURR_TUPLE             "indexCurrentTuple"
#define JTAG_MODULE_CAPABILITY          "moduleCapability"
#define JTAG_MODULE_CAP_AUTONOMOUS      "AUTONOMOUS"
#define JTAG_MODULE_CAP_DEPENDENT       "DEPENDENT"
#define JTAG_MODULE_MODE                "moduleMode"
#define JTAG_MODULE_MODE_MASTER         "MASTER"
#define JTAG_MODULE_MODE_SLAVE          "SLAVE"
#define JTAG_MODULE_DEV_OWNER           "moduleDeviceOwner"
#define JTAG_DEV_OWNER                  "OWNER"
#define JTAG_DEV_LEDER                  "LENDER"
#define JTAG_INDEX_TABLE                "indexTable"
#define JTAG_FIELD_ID                   "fieldId"
#define JTAG_PAYLOAD_OFFSET             "payloadOffset"
#define JTAG_PAYLOAD_SIZE               "payloadSize"
#define JTAG_PAYLOAD_CRC                "payloadCRC"
#define JTAG_PAYLOAD_VERSION            "payloadVersion"
#define JTAG_STATE                      "state"
#define JTAG_STATE_ENABLED              "ENABLED"
#define JTAG_STATE_DISABLED             "DISABLED"
#define JTAG_VALID                      "valid"
#define JTAG_NODE_INFO                  "nodeInfo"
#define JTAG_UUID                       "UUID"
#define JTAG_NAME                       "name"
#define JTAG_TYPE                       "type"
#define JTAG_PART_NUMBER                "partNumber"
#define JTAG_SKEW                       "skew"
#define JTAG_MAC                        "mac"
#define JTAG_SW_VERISION                "swVersion"
#define JTAG_PROD_SW_VERSION            "prodSwVersion"
#define JTAG_ASM_DATE                   "assemblyDate"
#define JTAG_OEM_NAME                   "oemName"
#define JTAG_MODULE_COUNT               "moduleCount"
#define JTAG_NODE_CONFIG                "nodeConfig"
#define JTAG_INVT_SYSFS_FILE            "invtSysFsFile"
#define JTAG_INVT_DEV_INFO              "invtDeviceInfo"
#define JTAG_BUS                        "bus"
#define JTAG_ADDRESS                    "address"
#define JTAG_MODULE_INFO                "moduleInfo"
#define JTAG_HW_VERSION                 "hwVersion"
#define JTAG_MFG_DATE                   "manufacturingDate"
#define JTAG_MFG_NAME                   "manufacturerName"
#define JTAG_DEVICE_COUNT               "deviceCount"
#define JTAG_MODULE_CONFIG              "moduleConfig"
#define JTAG_DESCRIPTION                "description"
#define JTAG_CLASS                      "class"
#define JTAG_DEV_SYSFS_FILE             "devSysFsFile"
#define JTAG_DEV_HW_ATTRS               "devHwAttrs"
#define JTAG_FACTORY_CONFIG             "factoryConfig"
#define JTAG_USER_CONFIG                "userConfig"
#define JTAG_FACTORY_CALIB              "factCalibaration"
#define JTAG_USER_CALIB                 "userCalibration"
#define JTAG_BOOTSTRAP_CERTS            "bootstrapCerts"
#define JTAG_CLOUD_CERTS                "cloudCerts"
#define JTAG_GPIO_DIRECTION             "direction"
#define JTAG_GPIO_NUMBER                "number"
#define JTAG_GPIO_DIRECTION_IN          "INPUT"
#define JTAG_GPIO_DIRECTION_OUT         "OUTPUT"
#define JTAG_UART                       "uartNumber"
#define JTAG_CHIP_SELECT                "chipSelect"
#define JTAG_DEVICE                     "device"
#define JTAG_VERSION                    "version"
#define JTAG_PROPERTY_TABLE             "propertyTable"
#define JTAG_ID                         "id"
#define JTAG_DATA_TYPE                  "dataType"
#define JTAG_PERMISSION                 "perm"
#define JTAG_AVAILABILITY               "available"
#define JTAG_PROPERTY_NAME              "propName"
#define JTAG_PROPERTY_DESC              "propDesc"
#define JTAG_PROPERTY_TYPE              "propType"
#define JTAG_UNITS                      "units"
#define JTAG_SYS_FS_FILE                "sysFsFile"
#define JTAG_DEPENDENT                  "dependent"
#define JTAG_CURR_PROP_ID               "currentValPropertyId"
#define JTAG_LIMIT_PROP_ID              "limitValPropertyId"
#define JTAG_ALERT_COND                 "alertCondition"
#define JTAG_VALUE                      "value"

#define JTAG_ERROR                      "error"
#define JTAG_ERROR_CODE                 "code"
#define JTAG_ERROR_CSTRING              "string"

#define JTAG_API_LIST                   "api"
#define JTAG_METHOD                     "method"
#define JTAG_URL_EP                     "endPoint"

#define JTAG_NOTIFY                     "notify"
#define JTAG_SERVICE_NAME               "serviceName"
#define JTAG_NOTIFICATION_TYPE          "notificationType"
#define JTAG_NODE_ID                    "nodeID"
#define JTAG_NODE_TYPE                  "nodeType"
#define JTAG_NOTIF_SEVERITY             "severity"
#define JTAG_NOTIF_DETAILS              "details"

#define JTAG_EPOCH_TIME                 "time"

typedef enum {
    TYPE_NULL = 0,
    TYPE_CHAR = 1,
    TYPE_BOOL,
    TYPE_UINT8,
    TYPE_UINT16,
    TYPE_UINT32,
    TYPE_INT8,
    TYPE_INT16,
    TYPE_INT32,
    TYPE_INT,
    TYPE_FLOAT,
    TYPE_DOUBLE,
    TYPE_ENUM,
    TYPE_STRING,
    TYPE_MAX
} DataType;

#ifdef __cplusplus
}
#endif

#endif /* INC_JSON_TYPES_H_ */
