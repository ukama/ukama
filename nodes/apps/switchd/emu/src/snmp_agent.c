/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <arpa/inet.h>
#include <errno.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <time.h>
#include <unistd.h>

#include "oid_map.h"
#include "snmp_agent.h"

#define ASN_INTEGER     0x02
#define ASN_OCTET_STR   0x04
#define ASN_NULL        0x05
#define ASN_OBJECT_ID   0x06
#define ASN_SEQUENCE    0x30

#define SNMP_GET        0xA0
#define SNMP_GETNEXT    0xA1
#define SNMP_RESPONSE   0xA2
#define SNMP_SET        0xA3

#define SNMP_ERR_NONE        0
#define SNMP_ERR_NOSUCHNAME  2

static size_t ber_put_len(uint8_t *buf, size_t len) {
    if (len < 128) {
        buf[0] = (uint8_t)len;
        return 1;
    }
    if (len < 256) {
        buf[0] = 0x81;
        buf[1] = (uint8_t)len;
        return 2;
    }
    buf[0] = 0x82;
    buf[1] = (uint8_t)((len >> 8) & 0xFF);
    buf[2] = (uint8_t)(len & 0xFF);
    return 3;
}

static int ber_get_len(const uint8_t *buf, size_t bufLen,
                       size_t *len, size_t *used) {
    if (bufLen < 1) return -1;
    if ((buf[0] & 0x80) == 0) {
        *len = buf[0];
        *used = 1;
        return 0;
    }
    if (buf[0] == 0x81 && bufLen >= 2) {
        *len = buf[1];
        *used = 2;
        return 0;
    }
    if (buf[0] == 0x82 && bufLen >= 3) {
        *len = ((size_t)buf[1] << 8) | buf[2];
        *used = 3;
        return 0;
    }
    return -1;
}

static int ber_put_tlv(uint8_t type, const uint8_t *val, size_t vlen,
                       uint8_t *out, size_t outLen, size_t *written) {
    size_t lused;
    if (outLen < 1) return -1;
    out[0] = type;
    lused = ber_put_len(out + 1, vlen);
    if (1 + lused + vlen > outLen) return -1;
    if (vlen > 0 && val != NULL)
        memcpy(out + 1 + lused, val, vlen);
    *written = 1 + lused + vlen;
    return 0;
}

static int ber_encode_int(int64_t value, uint8_t *buf, size_t bufLen,
                          size_t *written) {
    uint8_t tmp[9];
    size_t i = sizeof(tmp);
    int64_t v = value;
    int neg = (value < 0);
    do {
        tmp[--i] = (uint8_t)(v & 0xFF);
        v >>= 8;
    } while ((neg && v != -1) || (!neg && v != 0));
    if (!neg && (tmp[i] & 0x80)) tmp[--i] = 0;
    if (neg && !(tmp[i] & 0x80)) tmp[--i] = 0xFF;
    return ber_put_tlv(ASN_INTEGER, tmp + i, sizeof(tmp) - i,
                       buf, bufLen, written);
}

static int ber_decode_int(const uint8_t *buf, size_t len, int64_t *value) {
    size_t i;
    int64_t v = 0;
    if (len == 0 || len > 8) return -1;
    v = (buf[0] & 0x80) ? -1 : 0;
    for (i = 0; i < len; i++) v = (v << 8) | buf[i];
    *value = v;
    return 0;
}

static int ber_encode_string(uint8_t type, const void *data, size_t len,
                             uint8_t *buf, size_t bufLen, size_t *written) {
    return ber_put_tlv(type, (const uint8_t *)data, len,
                       buf, bufLen, written);
}

static int ber_encode_oid(const uint32_t *oid, size_t oidLen,
                          uint8_t *buf, size_t bufLen, size_t *written) {
    uint8_t tmp[256];
    size_t n = 0, i, s;
    uint32_t v;
    uint8_t stack[8];
    if (oidLen < 2) return -1;
    tmp[n++] = (uint8_t)(oid[0] * 40 + oid[1]);
    for (i = 2; i < oidLen; i++) {
        v = oid[i];
        s = 0;
        do {
            stack[s++] = (uint8_t)(v & 0x7F);
            v >>= 7;
        } while (v > 0 && s < sizeof(stack));
        while (s > 1) tmp[n++] = stack[--s] | 0x80;
        tmp[n++] = stack[0];
    }
    return ber_put_tlv(ASN_OBJECT_ID, tmp, n, buf, bufLen, written);
}

static int ber_decode_oid(const uint8_t *buf, size_t len,
                          uint32_t *oid, size_t *oidLen) {
    size_t i, n = 0;
    uint32_t v = 0;
    if (len == 0 || *oidLen < 2) return -1;
    oid[n++] = buf[0] / 40;
    oid[n++] = buf[0] % 40;
    for (i = 1; i < len; i++) {
        v = (v << 7) | (buf[i] & 0x7F);
        if ((buf[i] & 0x80) == 0) {
            if (n >= *oidLen) return -1;
            oid[n++] = v;
            v = 0;
        }
    }
    *oidLen = n;
    return 0;
}

static void oid_to_str(const uint32_t *oid, size_t len,
                       char *buf, size_t bufLen) {
    size_t off = 0;
    for (size_t i = 0; i < len && off < bufLen; i++) {
        off += (size_t)snprintf(buf + off,
                                (off < bufLen) ? bufLen - off : 0,
                                "%s%u", (i == 0) ? "" : ".", oid[i]);
    }
}

static int str_to_oid(const char *s, uint32_t *oid, size_t *oidLen) {
    char tmp[256];
    char *tok, *save;
    size_t n = 0;
    snprintf(tmp, sizeof(tmp), "%s", s);
    tok = strtok_r(tmp, ".", &save);
    while (tok != NULL) {
        if (n >= *oidLen) return -1;
        oid[n++] = (uint32_t)strtoul(tok, NULL, 10);
        tok = strtok_r(NULL, ".", &save);
    }
    *oidLen = n;
    return 0;
}

typedef struct {
    int32_t  reqId;
    uint8_t  pduType;
    uint32_t oid[64];
    size_t   oidLen;
    int64_t  setValue;
    int      hasSetValue;
} SnmpRequest;

static int parse_request(const uint8_t *buf, size_t len, SnmpRequest *req) {
    size_t off = 0, l = 0, used = 0;
    int64_t iv;

    memset(req, 0, sizeof(*req));

    if (off >= len || buf[off++] != ASN_SEQUENCE) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used;

    if (off >= len || buf[off++] != ASN_INTEGER) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used + l;

    if (off >= len || buf[off++] != ASN_OCTET_STR) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used + l;

    if (off >= len) return -1;
    req->pduType = buf[off++];
    if (req->pduType != SNMP_GET && req->pduType != SNMP_GETNEXT &&
        req->pduType != SNMP_SET)
        return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used;

    if (off >= len || buf[off++] != ASN_INTEGER) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used;
    if (ber_decode_int(buf + off, l, &iv) != 0) return -1;
    req->reqId = (int32_t)iv;
    off += l;

    if (off >= len || buf[off++] != ASN_INTEGER) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used + l;

    if (off >= len || buf[off++] != ASN_INTEGER) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used + l;

    if (off >= len || buf[off++] != ASN_SEQUENCE) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used;

    if (off >= len || buf[off++] != ASN_SEQUENCE) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used;

    if (off >= len || buf[off++] != ASN_OBJECT_ID) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used;
    req->oidLen = sizeof(req->oid) / sizeof(req->oid[0]);
    if (ber_decode_oid(buf + off, l, req->oid, &req->oidLen) != 0) return -1;
    off += l;

    if (off < len) {
        uint8_t vtype = buf[off++];
        if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
        off += used;
        if (vtype == ASN_INTEGER && l > 0) {
            if (ber_decode_int(buf + off, l, &iv) == 0) {
                req->setValue = iv;
                req->hasSetValue = 1;
            }
        }
    }

    return 0;
}

static int build_response(int32_t reqId, int errStatus,
                          const uint32_t *oid, size_t oidLen,
                          int isString, int intVal,
                          const char *strVal,
                          uint8_t *out, size_t outLen, size_t *written) {
    uint8_t oidTlv[256], valTlv[300];
    uint8_t vbContent[600], vbSeq[700], vblSeq[800];
    uint8_t pduContent[900], pdu[1000];
    uint8_t msgContent[1100];
    size_t nOid, nVal, nVb, nVbSeq, nVbl;
    size_t nReq, nErr, nErrIdx, nPdu;
    size_t nVer, nCom;
    uint8_t *p;

    if (ber_encode_oid(oid, oidLen, oidTlv, sizeof(oidTlv), &nOid) != 0)
        return -1;

    if (errStatus != 0) {
        if (ber_put_tlv(ASN_NULL, NULL, 0, valTlv, sizeof(valTlv), &nVal) != 0)
            return -1;
    } else if (isString) {
        size_t slen = strVal ? strlen(strVal) : 0;
        if (ber_encode_string(ASN_OCTET_STR, strVal, slen,
                              valTlv, sizeof(valTlv), &nVal) != 0)
            return -1;
    } else {
        if (ber_encode_int(intVal, valTlv, sizeof(valTlv), &nVal) != 0)
            return -1;
    }

    p = vbContent;
    memcpy(p, oidTlv, nOid); p += nOid;
    memcpy(p, valTlv, nVal); p += nVal;
    nVb = (size_t)(p - vbContent);

    if (ber_put_tlv(ASN_SEQUENCE, vbContent, nVb,
                    vbSeq, sizeof(vbSeq), &nVbSeq) != 0)
        return -1;

    if (ber_put_tlv(ASN_SEQUENCE, vbSeq, nVbSeq,
                    vblSeq, sizeof(vblSeq), &nVbl) != 0)
        return -1;

    p = pduContent;
    if (ber_encode_int(reqId, p,
                       sizeof(pduContent), &nReq) != 0)
        return -1;
    p += nReq;
    if (ber_encode_int(errStatus, p,
                       sizeof(pduContent) - (size_t)(p - pduContent), &nErr) != 0)
        return -1;
    p += nErr;
    if (ber_encode_int(0, p,
                       sizeof(pduContent) - (size_t)(p - pduContent), &nErrIdx) != 0)
        return -1;
    p += nErrIdx;
    memcpy(p, vblSeq, nVbl);
    p += nVbl;

    if (ber_put_tlv(SNMP_RESPONSE, pduContent, (size_t)(p - pduContent),
                    pdu, sizeof(pdu), &nPdu) != 0)
        return -1;

    p = msgContent;
    if (ber_encode_int(1, p, sizeof(msgContent), &nVer) != 0) return -1;
    p += nVer;
    if (ber_encode_string(ASN_OCTET_STR, "public", 6,
                          p, sizeof(msgContent) - (size_t)(p - msgContent),
                          &nCom) != 0)
        return -1;
    p += nCom;
    memcpy(p, pdu, nPdu);
    p += nPdu;

    if (ber_put_tlv(ASN_SEQUENCE, msgContent, (size_t)(p - msgContent),
                    out, outLen, written) != 0)
        return -1;

    return 0;
}

static void *snmp_main(void *arg) {
    EmuModel *model = (EmuModel *)arg;
    uint8_t buf[EMU_SNMP_BUF];

    while (model->running) {
        struct sockaddr_in cli;
        socklen_t slen = sizeof(cli);

        ssize_t bytes = recvfrom(model->snmpFd, buf, sizeof(buf), 0,
                                 (struct sockaddr *)&cli, &slen);
        if (bytes < 0) {
            if (errno == EINTR) continue;
            if (!model->running) break;
            continue;
        }

        if (model->faults.unreachable || !model->info.reachable) {
            continue;
        }

        if (model->faults.snmpDelayMs > 0) {
            struct timespec ts;
            ts.tv_sec  = model->faults.snmpDelayMs / 1000;
            ts.tv_nsec = (long)(model->faults.snmpDelayMs % 1000) * 1000000L;
            nanosleep(&ts, NULL);
        }

        SnmpRequest req;
        if (parse_request(buf, (size_t)bytes, &req) != 0) {
            continue;
        }

        char oidStr[256];
        oid_to_str(req.oid, req.oidLen, oidStr, sizeof(oidStr));

        uint8_t resp[1500];
        size_t respLen = 0;
        int intVal = 0;
        char strVal[256] = {0};
        int isString = 0;
        int errStatus = SNMP_ERR_NONE;

        uint32_t respOid[64];
        size_t respOidLen = req.oidLen;
        memcpy(respOid, req.oid, req.oidLen * sizeof(uint32_t));

        pthread_mutex_lock(&model->lock);

        if (req.pduType == SNMP_GET) {
            if (oid_get_int(model, oidStr, &intVal) == STATUS_OK) {
                isString = 0;
            } else if (oid_get_string(model, oidStr, strVal,
                                       sizeof(strVal)) == STATUS_OK) {
                isString = 1;
            } else {
                errStatus = SNMP_ERR_NOSUCHNAME;
            }

        } else if (req.pduType == SNMP_GETNEXT) {
            char nextStr[256];
            if (oid_get_next(model, oidStr, nextStr,
                             sizeof(nextStr)) == STATUS_OK) {
                respOidLen = sizeof(respOid) / sizeof(respOid[0]);
                str_to_oid(nextStr, respOid, &respOidLen);

                if (oid_get_int(model, nextStr, &intVal) == STATUS_OK) {
                    isString = 0;
                } else if (oid_get_string(model, nextStr, strVal,
                                           sizeof(strVal)) == STATUS_OK) {
                    isString = 1;
                } else {
                    errStatus = SNMP_ERR_NOSUCHNAME;
                }
            } else {
                errStatus = SNMP_ERR_NOSUCHNAME;
            }

        } else if (req.pduType == SNMP_SET) {
            if (req.hasSetValue) {
                if (oid_set_int(model, oidStr,
                                (int)req.setValue) != STATUS_OK) {
                    errStatus = SNMP_ERR_NOSUCHNAME;
                }
                intVal = (int)req.setValue;
            } else {
                errStatus = SNMP_ERR_NOSUCHNAME;
            }
        }

        pthread_mutex_unlock(&model->lock);

        if (build_response(req.reqId, errStatus,
                           respOid, respOidLen,
                           isString, intVal, strVal,
                           resp, sizeof(resp), &respLen) == 0) {
            sendto(model->snmpFd, resp, respLen, 0,
                   (struct sockaddr *)&cli, slen);
        }
    }

    return NULL;
}

int snmp_agent_start(EmuModel *model) {
    struct sockaddr_in addr;

    model->snmpFd = socket(AF_INET, SOCK_DGRAM, 0);
    if (model->snmpFd < 0) {
        return STATUS_NOK;
    }

    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_port   = htons((uint16_t)model->cfg.snmpPort);
    addr.sin_addr.s_addr = inet_addr(model->cfg.bindAddr);

    if (bind(model->snmpFd, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
        return STATUS_NOK;
    }

    return pthread_create(&model->snmpThread, NULL, snmp_main, model);
}

void snmp_agent_stop(EmuModel *model) {
    if (model->snmpFd >= 0) {
        close(model->snmpFd);
        model->snmpFd = -1;
    }

    if (model->snmpThread != 0U) {
        pthread_join(model->snmpThread, NULL);
        model->snmpThread = 0;
    }
}
