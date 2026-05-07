/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SWITCHD_TYPES_H
#define SWITCHD_TYPES_H

#include <pthread.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <sys/types.h>
#include <time.h>

#define SWITCHD_MAX_PORTS          32
#define SWITCHD_NAME_LEN           64
#define SWITCHD_STR_LEN            128
#define SWITCHD_LONG_STR_LEN       256
#define SWITCHD_OP_DETAIL_LEN      256
#define SWITCHD_ALARM_RESOURCE_LEN 64
#define SWITCHD_ALARM_TEXT_LEN     160
#define SWITCHD_STAGE_NAME_LEN     128
#define SWITCHD_STAGE_PATH_LEN     256
#define SWITCHD_SHA256_LEN         80

typedef enum {
    SWITCHD_OK = 0,
    SWITCHD_ERR_NOMEM,
    SWITCHD_ERR_INVAL,
    SWITCHD_ERR_IO,
    SWITCHD_ERR_TIMEOUT,
    SWITCHD_ERR_BUSY,
    SWITCHD_ERR_NOTFOUND,
    SWITCHD_ERR_UNSUPPORTED,
    SWITCHD_ERR_SNMP,
    SWITCHD_ERR_PROTOCOL,
    SWITCHD_ERR_STATE,
    SWITCHD_ERR_AUTH,
    SWITCHD_ERR_INTERNAL
} SwitchdError;

typedef enum {
    SWITCHD_STATE_INIT = 0,
    SWITCHD_STATE_READY,
    SWITCHD_STATE_BUSY,
    SWITCHD_STATE_DEGRADED,
    SWITCHD_STATE_UPDATING,
    SWITCHD_STATE_RECOVERING,
    SWITCHD_STATE_ERROR,
    SWITCHD_STATE_TERMINATING
} SwitchdState;

typedef enum {
    SWITCHD_OP_NONE = 0,
    SWITCHD_OP_PORT_ADMIN_SET,
    SWITCHD_OP_PORT_POE_SET,
    SWITCHD_OP_PORT_POE_CYCLE,
    SWITCHD_OP_SWITCH_REBOOT,
    SWITCHD_OP_FW_STAGE,
    SWITCHD_OP_FW_APPLY
} SwitchdOperationType;

typedef enum {
    SWITCHD_OP_STATE_IDLE = 0,
    SWITCHD_OP_STATE_RUNNING,
    SWITCHD_OP_STATE_DONE,
    SWITCHD_OP_STATE_FAILED
} SwitchdOperationState;

typedef enum {
    SWITCHD_FW_IDLE = 0,
    SWITCHD_FW_STAGED,
    SWITCHD_FW_APPLYING,
    SWITCHD_FW_REBOOTING,
    SWITCHD_FW_RECONNECTING,
    SWITCHD_FW_VERIFYING,
    SWITCHD_FW_DONE,
    SWITCHD_FW_FAILED
} SwitchdFirmwareState;

typedef enum {
    SWITCHD_ALARM_SWITCH_UNREACHABLE = 1,
    SWITCHD_ALARM_SWITCH_RECOVERED,
    SWITCHD_ALARM_PORT_LINK_DOWN,
    SWITCHD_ALARM_PORT_POE_OFF,
    SWITCHD_ALARM_PORT_POE_FAULT,
    SWITCHD_ALARM_HIGH_SYSTEM_TEMP,
    SWITCHD_ALARM_HIGH_AMBIENT_TEMP,
    SWITCHD_ALARM_FIRMWARE_FAILED,
    SWITCHD_ALARM_FIRMWARE_DONE
} SwitchdAlarmCode;

typedef enum {
    SWITCHD_ALARM_SEV_INFO = 0,
    SWITCHD_ALARM_SEV_WARNING,
    SWITCHD_ALARM_SEV_CRITICAL
} SwitchdAlarmSeverity;

typedef struct {
    bool supportsPortAdmin;
    bool supportsPoeControl;
    bool supportsPoeCycle;
    bool supportsPortCounters;
    bool supportsPowerMetrics;
    bool supportsSystemMetrics;
    bool supportsFirmwareUpdate;
    bool supportsSaveConfig;
    uint32_t maxPorts;
} SwitchCapabilities;

typedef struct {
    char vendor[SWITCHD_NAME_LEN];
    char model[SWITCHD_NAME_LEN];
    char serial[SWITCHD_NAME_LEN];
    char hardwareVersion[SWITCHD_NAME_LEN];
    char softwareVersion[SWITCHD_NAME_LEN];
    char managementAddress[SWITCHD_NAME_LEN];
    bool reachable;
    uint32_t portCount;
    int pollFailures;
    time_t updatedAt;
} SwitchInfo;

typedef struct {
    double poeTotalPowerWatts;
    double poeMaxPowerWatts;
    double systemTemperatureC;
    double ambientTemperatureC;
    double systemPowerWatts;
    double inputVoltage;
    double systemCurrentAmps;
    bool inputLinkFailureAlarm;
    bool inputPoeFailureAlarm;
    time_t updatedAt;
} SwitchKpis;

typedef struct {
    uint32_t id;
    char name[SWITCHD_NAME_LEN];
    char media[16];
    bool present;
    bool adminUp;
    bool linkUp;
    bool poeSupported;
    bool poeEnabled;
    bool poeOperational;
    int poeClass;
    double powerWatts;
    double voltage;
    double currentAmps;
    uint64_t speedBps;
    uint64_t rxBytes;
    uint64_t txBytes;
    uint64_t rxPackets;
    uint64_t txPackets;
    uint64_t rxErrors;
    uint64_t txErrors;
    uint64_t rxDrops;
    uint64_t txDrops;
    char fault[SWITCHD_STR_LEN];
    time_t updatedAt;
} SwitchPortState;

typedef struct {
    uint64_t id;
    SwitchdOperationType type;
    SwitchdOperationState state;
    uint32_t portId;
    int progress;
    SwitchdError error;
    char detail[SWITCHD_OP_DETAIL_LEN];
    time_t startedAt;
    time_t endedAt;
} SwitchOperation;

typedef struct {
    SwitchdAlarmCode code;
    SwitchdAlarmSeverity severity;
    char resource[SWITCHD_ALARM_RESOURCE_LEN];
    char text[SWITCHD_ALARM_TEXT_LEN];
    bool active;
    bool latched;
    time_t firstSeen;
    time_t lastSeen;
    time_t lastSent;
} SwitchAlarm;

typedef struct {
    char path[SWITCHD_STAGE_PATH_LEN];
    char version[SWITCHD_NAME_LEN];
    char sha256[SWITCHD_SHA256_LEN];
    char tftpFilename[SWITCHD_STAGE_NAME_LEN];
    off_t size;
    SwitchdFirmwareState state;
    int executeStatus;
    time_t stagedAt;
    time_t updatedAt;
    char detail[SWITCHD_OP_DETAIL_LEN];
} SwitchFirmware;



typedef enum {
    SWITCH_POLICY_STATE_MISSING = 0,
    SWITCH_POLICY_STATE_LOADED,
    SWITCH_POLICY_STATE_INVALID
} SwitchPolicyState;

typedef enum {
    SWITCH_PORT_POLICY_UNKNOWN = 0,
    SWITCH_PORT_POLICY_PROTECTED,
    SWITCH_PORT_POLICY_FREE,
    SWITCH_PORT_POLICY_NEVER_OFF_REMOTE,
    SWITCH_PORT_POLICY_DISABLED
} SwitchPortPolicyType;

typedef enum {
    SWITCH_POLICY_ACTION_ADMIN_UP = 0,
    SWITCH_POLICY_ACTION_ADMIN_DOWN,
    SWITCH_POLICY_ACTION_POE_ON,
    SWITCH_POLICY_ACTION_POE_OFF,
    SWITCH_POLICY_ACTION_POE_CYCLE
} SwitchPolicyAction;

typedef struct {
    uint32_t port;
    char role[SWITCHD_NAME_LEN];
    char nodeId[SWITCHD_NAME_LEN];
    char klass[SWITCHD_NAME_LEN];
    SwitchPortPolicyType policy;
    bool present;
} SwitchPortPolicy;

typedef struct {
    SwitchPolicyState state;
    char siteId[SWITCHD_NAME_LEN];
    char source[SWITCHD_NAME_LEN];
    char updatedAt[SWITCHD_NAME_LEN];
    char path[SWITCHD_STAGE_PATH_LEN];
    char error[SWITCHD_OP_DETAIL_LEN];
    SwitchPortPolicy ports[SWITCHD_MAX_PORTS];
    time_t loadedAt;
} SwitchPolicy;

typedef struct {
    char driverName[SWITCHD_NAME_LEN];
    char httpHost[64];
    int httpPort;
    char urlPrefix[32];

    char snmpHost[64];
    int snmpPort;
    char snmpCommunity[64];
    int snmpVersion;
    int snmpTimeoutMs;
    int snmpRetries;

    int pollStatusSec;
    int pollKpisSec;
    int pollInfoSec;
    int alarmScanSec;

    int commandTimeoutMs;
    int firmwareReconnectSec;
    int firmwareVerifySec;
    int poeCycleMs;

    char notifyUrl[SWITCHD_LONG_STR_LEN];
    int notifyTimeoutMs;

    char tftpBindIp[64];
    int tftpPort;
    char tftpRoot[SWITCHD_STAGE_PATH_LEN];

    bool strictLinkAlarms;
    bool saveAfterWrite;
    char policyPath[SWITCHD_STAGE_PATH_LEN];
} SwitchdConfig;

struct SwitchDriver;
typedef struct SwitchDriver SwitchDriver;

typedef struct {
    SwitchdConfig config;
    SwitchdState state;
    bool terminate;

    pthread_mutex_t stateMutex;
    pthread_mutex_t opMutex;
    pthread_mutex_t alarmMutex;
    pthread_mutex_t driverMutex;

    pthread_t pollerThread;
    pthread_t alarmThread;
    bool pollerRunning;
    bool alarmRunning;

    SwitchCapabilities caps;
    SwitchInfo info;
    SwitchKpis kpis;
    SwitchPortState ports[SWITCHD_MAX_PORTS];
    uint32_t portCount;

    SwitchOperation op;
    uint64_t nextOpId;
    SwitchFirmware fw;

    SwitchAlarm alarms[SWITCHD_MAX_PORTS + 8];
    size_t alarmCount;

    SwitchPolicy policy;

    SwitchDriver *driver;
} SwitchdContext;

#endif
