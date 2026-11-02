/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SERVICE_CLIENT_H_
#define SERVICE_CLIENT_H_

#include <stdbool.h>
#include <jansson.h>

#define HTTP_TRANSPORT_ERROR 0
#define HTTP_SCHEME          "http://"
#define SERVICE_PATH_LEN     16

#define ENV_NETWORK_URL "UKAMA_NETWORK_URL"
#define ENV_PCRF_URL    "UKAMA_PCRF_URL"
#define ENV_EPC_URL     "UKAMA_EPC_URL"
#define ENV_EPC_KIND    "UKAMA_EPC_KIND"

#define DEF_NETWORK_URL "http://127.0.0.1:18091/v1"
#define DEF_PCRF_URL    "http://127.0.0.1:18030/v1"
#define DEF_EPCEMU_URL  "http://127.0.0.1:18028/v1"

#define EPC_KIND_REAL       "epc"
#define EPC_KIND_EMULATED   "epcemu"
#define SERVICE_STATE_ON    "on"
#define SERVICE_STATE_OFF   "off"
#define SERVICE_API_SERVICE "/service"
#define SERVICE_API_STATUS  "/status"

#define JTAG_STATE           "state"
#define JTAG_SERVICE         "service"
#define JTAG_READY           "ready"
#define JTAG_PCRF            "pcrf"
#define JTAG_PDN_COUNT       "pdn_count"
#define JTAG_CLEANUP_PENDING "cleanup_pending"
#define JTAG_CONFIRMED       "confirmed"

#define SERVICE_MEMBER_COUNT 3
#define SERVICE_URL_LEN      256

typedef enum {
    SERVICE_MEMBER_NETWORK = 0,
    SERVICE_MEMBER_PCRF,
    SERVICE_MEMBER_EPC
} ServiceMember;

typedef struct {
    char url[SERVICE_URL_LEN];
    bool emulated;
} ServicePeer;

typedef struct {
    bool reachable;
    bool ready;
    bool on;
    bool confirmed;
    int pdnCount;
} ServiceObservation;

int service_http(const char *method, const char *url, json_t *body,
                 json_t **reply);
bool service_peers_init(ServicePeer peers[SERVICE_MEMBER_COUNT]);
bool service_peer_set(const ServicePeer *peer, bool enabled);
ServiceObservation service_peer_read(const ServicePeer *peer);

#endif /* SERVICE_CLIENT_H_ */
