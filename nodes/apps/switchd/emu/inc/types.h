/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef TYPES_H
#define TYPES_H

#include <pthread.h>
#include <stdint.h>
#include <stddef.h>
#include <time.h>

#include "switchemu.h"

typedef struct {
    int httpPort;
    int snmpPort;
    int tftpPort;
    int notifyPort;
    int logLevel;
    char bindAddr[64];
    char stateFile[EMU_MAX_PATH];
    char scenario[64];
    char notifyHost[64];
    char notifyPath[128];
    int notifyEnabled;
} EmuConfig;

typedef struct {
    uint32_t id;
    char name[32];
    char media[16];
    int present;
    int adminUp;
    int linkUp;
    uint32_t speedMbps;
    int fullDuplex;
    uint64_t rxBytes;
    uint64_t txBytes;
    uint64_t rxPackets;
    uint64_t txPackets;
    uint64_t rxErrors;
    uint64_t txErrors;
    int poeSupported;
    int poeAdminEnabled;
    int poeOperStatus;
    int poePowerMw;
    int poeCurrentMa;
    int poeVoltageMv;
    int poeClassId;
    int faultPoe;
    int faultLink;
    time_t updatedAt;
} EmuPortState;

typedef enum {
    FW_IDLE = 0,
    FW_STAGED,
    FW_APPLYING,
    FW_REBOOTING,
    FW_DONE,
    FW_FAILED
} EmuFirmwareState;

typedef struct {
    EmuFirmwareState state;
    char stagedFilename[128];
    char stagedPath[EMU_MAX_PATH];
    char stagedVersion[64];
    int executeStatus;
    int applyShouldFail;
    int rebootDelaySec;
    int applyDelaySec;
    time_t stateSince;
} EmuFirmware;

typedef struct {
    char manufacturer[64];
    char serial[64];
    char hardwareVersion[64];
    char softwareVersion[64];
    int reachable;
    int rebooting;
    int systemTempC;
    int ambientTempC;
    int inputVoltageMv;
    int systemCurrentMa;
    int systemPowerMw;
    int poeBudgetMw;
    int poeUsedMw;
    int alarmLinkFailure;
    int alarmPoeFailure;
    time_t updatedAt;
} EmuSwitchInfo;

typedef struct {
    int snmpDelayMs;
    int unreachable;
    int flapPortId;
    int flapPeriodSec;
    int tftpFail;
    int snmpSetFail;
} EmuFaults;

typedef struct {
    int active;
    int code;
    char source[32];
    char severity[16];
    char message[128];
    time_t raisedAt;
} EmuAlarm;

typedef struct {
    pthread_mutex_t lock;
    int running;
    int terminate;
    int httpFd;
    int snmpFd;
    int tftpFd;
    pthread_t httpThread;
    pthread_t snmpThread;
    pthread_t tftpThread;
    pthread_t engineThread;
    EmuConfig cfg;
    EmuSwitchInfo info;
    EmuPortState ports[EMU_MAX_PORTS];
    size_t portCount;
    EmuFirmware firmware;
    EmuFaults faults;
    EmuAlarm alarms[EMU_MAX_ALARMS];
    size_t alarmCount;
    char activeScenario[64];
} EmuModel;

#endif /* TYPES_H */
