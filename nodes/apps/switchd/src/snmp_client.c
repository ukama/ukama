/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <stdbool.h>
#include <unistd.h>
#include <errno.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <sys/select.h>
#include "snmp_client.h"
#include "types.h"
#include "log.h"

#define ASN_INTEGER         0x02
#define ASN_OCTET_STR       0x04
#define ASN_NULL            0x05
#define ASN_OBJECT_ID       0x06
#define ASN_SEQUENCE        0x30
#define ASN_IPADDRESS       0x40
#define ASN_COUNTER32       0x41
#define ASN_GAUGE32         0x42
#define ASN_TIMETICKS       0x43
#define ASN_OPAQUE          0x44
#define ASN_COUNTER64       0x46

#define SNMP_GET            0xA0
#define SNMP_GETNEXT        0xA1
#define SNMP_RESPONSE       0xA2
#define SNMP_SET            0xA3

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

static int ber_get_len(const uint8_t *buf, size_t bufLen, size_t *len, size_t *used) {
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

static int ber_put_tlv(uint8_t type, const uint8_t *val, size_t vlen, uint8_t *out, size_t outLen, size_t *written) {
    size_t lused;
    if (outLen < 1) return -1;
    out[0] = type;
    lused = ber_put_len(out + 1, vlen);
    if (1 + lused + vlen > outLen) return -1;
    if (vlen > 0 && val != NULL) memcpy(out + 1 + lused, val, vlen);
    *written = 1 + lused + vlen;
    return 0;
}

static int ber_encode_int(int64_t value, uint8_t *buf, size_t bufLen, size_t *written) {
    uint8_t tmp[9];
    size_t i = sizeof(tmp);
    int64_t v = value;
    bool neg = (value < 0);
    do {
        tmp[--i] = (uint8_t)(v & 0xFF);
        v >>= 8;
    } while ((neg && v != -1) || (!neg && v != 0));
    if (!neg && (tmp[i] & 0x80)) tmp[--i] = 0;
    if (neg && !(tmp[i] & 0x80)) tmp[--i] = 0xFF;
    return ber_put_tlv(ASN_INTEGER, tmp + i, sizeof(tmp) - i, buf, bufLen, written);
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

static int ber_encode_string(uint8_t type, const void *data, size_t len, uint8_t *buf, size_t bufLen, size_t *written) {
    return ber_put_tlv(type, (const uint8_t *)data, len, buf, bufLen, written);
}

static int ber_encode_oid(const uint32_t *oid, size_t oidLen, uint8_t *buf, size_t bufLen, size_t *written) {
    uint8_t tmp[256];
    size_t n = 0;
    size_t i;
    uint32_t v;
    uint8_t stack[8];
    size_t s;
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

static int ber_decode_oid(const uint8_t *buf, size_t len, uint32_t *oid, size_t *oidLen) {
    size_t i;
    size_t n = 0;
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

static int decode_value(uint8_t type, const uint8_t *buf, size_t len, SnmpValue *out) {
    size_t olen;
    int64_t iv;
    if (out == NULL) return -1;
    memset(out, 0, sizeof(*out));
    switch (type) {
        case ASN_INTEGER:
            if (ber_decode_int(buf, len, &iv) != 0) return -1;
            out->type = SNMP_VALUE_INT;
            out->intValue = iv;
            return 0;
        case ASN_COUNTER32:
        case ASN_GAUGE32:
        case ASN_TIMETICKS:
        case ASN_COUNTER64:
            if (ber_decode_int(buf, len, &iv) != 0) return -1;
            out->type = SNMP_VALUE_UINT;
            out->uintValue = (uint64_t)iv;
            return 0;
        case ASN_OCTET_STR:
        case ASN_IPADDRESS:
        case ASN_OPAQUE:
            out->type = SNMP_VALUE_STRING;
            if (len >= sizeof(out->stringValue)) len = sizeof(out->stringValue) - 1;
            memcpy(out->stringValue, buf, len);
            out->stringValue[len] = '\0';
            return 0;
        case ASN_OBJECT_ID:
            olen = sizeof(out->oid) / sizeof(out->oid[0]);
            if (ber_decode_oid(buf, len, out->oid, &olen) != 0) return -1;
            out->type = SNMP_VALUE_OID;
            out->oidLen = olen;
            return 0;
        case ASN_NULL:
            out->type = SNMP_VALUE_NONE;
            return 0;
        default:
            return -1;
    }
}

static int build_request(int pduType, const char *community, int32_t reqId,
                         const uint32_t *oid, size_t oidLen,
                         const SnmpValue *setVal,
                         uint8_t *out, size_t outLen, size_t *written) {
    uint8_t oidTlv[256];
    uint8_t valTlv[300];
    uint8_t vbContent[600];
    uint8_t vbSeq[700];
    uint8_t vblSeq[800];
    uint8_t pduContent[900];
    uint8_t pdu[1000];
    uint8_t msgContent[1100];
    size_t nOid, nVal, nVbSeq, nVbl, nReq, nErr, nErrIdx;
    size_t nPdu, nVer, nCom, nMsg;
    uint8_t *ptr;

    if (ber_encode_oid(oid, oidLen, oidTlv, sizeof(oidTlv), &nOid) != 0) {
        return -1;
    }

    if (setVal == NULL || setVal->type == SNMP_VALUE_NONE) {
        if (ber_put_tlv(ASN_NULL, NULL, 0, valTlv, sizeof(valTlv),
                        &nVal) != 0) {
            return -1;
        }
    } else if (setVal->type == SNMP_VALUE_INT) {
        if (ber_encode_int(setVal->intValue, valTlv, sizeof(valTlv),
                           &nVal) != 0) {
            return -1;
        }
    } else if (setVal->type == SNMP_VALUE_STRING) {
        if (ber_encode_string(ASN_OCTET_STR, setVal->stringValue,
                              strlen(setVal->stringValue),
                              valTlv, sizeof(valTlv), &nVal) != 0) {
            return -1;
        }
    } else {
        return -1;
    }

    ptr = vbContent;
    memcpy(ptr, oidTlv, nOid);
    ptr += nOid;
    memcpy(ptr, valTlv, nVal);
    ptr += nVal;

    if (ber_put_tlv(ASN_SEQUENCE,
                    vbContent,
                    (size_t)(ptr - vbContent),
                    vbSeq,
                    sizeof(vbSeq),
                    &nVbSeq) != 0) {
        return -1;
    }

    /*
     * Do not wrap in-place. The old code used vblSeq as both source and
     * destination here, corrupting SNMP packets before the emulator could
     * parse them.
     */
    if (ber_put_tlv(ASN_SEQUENCE,
                    vbSeq,
                    nVbSeq,
                    vblSeq,
                    sizeof(vblSeq),
                    &nVbl) != 0) {
        return -1;
    }

    ptr = pduContent;
    if (ber_encode_int(reqId, ptr, sizeof(pduContent), &nReq) != 0) {
        return -1;
    }
    ptr += nReq;

    if (ber_encode_int(0, ptr,
                       sizeof(pduContent) - (size_t)(ptr - pduContent),
                       &nErr) != 0) {
        return -1;
    }
    ptr += nErr;

    if (ber_encode_int(0, ptr,
                       sizeof(pduContent) - (size_t)(ptr - pduContent),
                       &nErrIdx) != 0) {
        return -1;
    }
    ptr += nErrIdx;

    memcpy(ptr, vblSeq, nVbl);
    ptr += nVbl;

    if (ber_put_tlv((uint8_t)pduType,
                    pduContent,
                    (size_t)(ptr - pduContent),
                    pdu,
                    sizeof(pdu),
                    &nPdu) != 0) {
        return -1;
    }

    ptr = msgContent;
    if (ber_encode_int(1, ptr, sizeof(msgContent), &nVer) != 0) {
        return -1;
    }
    ptr += nVer;

    if (ber_encode_string(ASN_OCTET_STR,
                          community,
                          strlen(community),
                          ptr,
                          sizeof(msgContent) - (size_t)(ptr - msgContent),
                          &nCom) != 0) {
        return -1;
    }
    ptr += nCom;

    memcpy(ptr, pdu, nPdu);
    ptr += nPdu;

    if (ber_put_tlv(ASN_SEQUENCE,
                    msgContent,
                    (size_t)(ptr - msgContent),
                    out,
                    outLen,
                    &nMsg) != 0) {
        return -1;
    }

    *written = nMsg;
    return 0;
}

static int parse_response(const uint8_t *buf, size_t len, int32_t expectReqId, SnmpVarBind *out) {
    size_t off = 0, l = 0, used = 0;
    uint8_t type;
    int64_t iv;
    uint32_t oid[64];
    size_t oidLen;

    if (len < 2 || buf[off++] != ASN_SEQUENCE) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used;
    if (off + l > len) return -1;

    if (buf[off++] != ASN_INTEGER) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used + l;

    if (buf[off++] != ASN_OCTET_STR) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used + l;

    type = buf[off++];
    if (type != SNMP_RESPONSE) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used;

    if (buf[off++] != ASN_INTEGER) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used;
    if (ber_decode_int(buf + off, l, &iv) != 0 || iv != expectReqId) return -1;
    off += l;

    if (buf[off++] != ASN_INTEGER) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used;
    if (ber_decode_int(buf + off, l, &iv) != 0 || iv != 0) return -1;
    off += l;

    if (buf[off++] != ASN_INTEGER) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used + l;

    if (buf[off++] != ASN_SEQUENCE) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used;
    if (buf[off++] != ASN_SEQUENCE) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used;

    if (buf[off++] != ASN_OBJECT_ID) return -1;
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used;
    oidLen = sizeof(oid) / sizeof(oid[0]);
    if (ber_decode_oid(buf + off, l, oid, &oidLen) != 0) return -1;
    off += l;

    if (out != NULL) {
        memset(out, 0, sizeof(*out));
        memcpy(out->oid, oid, oidLen * sizeof(uint32_t));
        out->oidLen = oidLen;
    }

    type = buf[off++];
    if (ber_get_len(buf + off, len - off, &l, &used) != 0) return -1;
    off += used;
    if (out != NULL) return decode_value(type, buf + off, l, &out->value);
    return 0;
}

int snmp_session_init(SnmpSession *s, const char *host, int port,
                      const char *community, int timeoutMs, int retries) {
    if (s == NULL || host == NULL || community == NULL) return SWITCHD_ERR_INVAL;
    memset(s, 0, sizeof(*s));
    snprintf(s->host, sizeof(s->host), "%s", host);
    snprintf(s->community, sizeof(s->community), "%s", community);
    s->port = port;
    s->timeoutMs = timeoutMs;
    s->retries = retries;
    return SWITCHD_OK;
}

static void snmp_oid_debug(const uint32_t *oid, size_t oidLen,
                           char *buf, size_t bufLen) {
    size_t i;
    size_t off;

    if (!buf || bufLen == 0) {
        return;
    }

    off = 0;
    buf[0] = '\0';

    for (i = 0; i < oidLen && off < bufLen; i++) {
        off += (size_t)snprintf(buf + off,
                                (off < bufLen) ? bufLen - off : 0,
                                "%s%u",
                                (i == 0) ? "" : ".",
                                oid[i]);
    }
}

static int snmp_request(SnmpSession *s, int pduType, const uint32_t *oid, size_t oidLen,
                        const SnmpValue *setVal, SnmpVarBind *out) {
    int sockFd = -1;
    struct sockaddr_in addr;
    uint8_t req[1200];
    uint8_t resp[1500];
    size_t reqLen;
    ssize_t n;
    int attempt;
    int timeoutMs;
    int32_t reqId;
    fd_set rfds;
    struct timeval tv;
    char oidStr[256];

    if (!s || !oid || oidLen == 0) {
        return SWITCHD_ERR_INVAL;
    }

    snmp_oid_debug(oid, oidLen, oidStr, sizeof(oidStr));
    timeoutMs = (s->timeoutMs > 0) ? s->timeoutMs : 500;

    reqId = (int32_t)(time(NULL) ^ getpid() ^ rand());
    if (build_request(pduType,
                      s->community,
                      reqId,
                      oid,
                      oidLen,
                      setVal,
                      req,
                      sizeof(req),
                      &reqLen) != 0) {
        log_error("snmp: build request failed pdu=%d oid=%s",
                  pduType, oidStr);
        return SWITCHD_ERR_PROTOCOL;
    }

    sockFd = socket(AF_INET, SOCK_DGRAM, 0);
    if (sockFd < 0) {
        log_error("snmp: socket failed host=%s port=%d errno=%d",
                  s->host, s->port, errno);
        return SWITCHD_ERR_IO;
    }

    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_port = htons((uint16_t)s->port);

    if (inet_pton(AF_INET, s->host, &addr.sin_addr) != 1) {
        log_error("snmp: invalid host '%s'", s->host);
        close(sockFd);
        return SWITCHD_ERR_INVAL;
    }

    for (attempt = 0; attempt <= s->retries; attempt++) {
        if (sendto(sockFd,
                   req,
                   reqLen,
                   0,
                   (struct sockaddr *)&addr,
                   sizeof(addr)) < 0) {
            log_error("snmp: sendto failed host=%s port=%d oid=%s errno=%d",
                      s->host, s->port, oidStr, errno);
            close(sockFd);
            return SWITCHD_ERR_IO;
        }

        FD_ZERO(&rfds);
        FD_SET(sockFd, &rfds);

        tv.tv_sec = timeoutMs / 1000;
        tv.tv_usec = (timeoutMs % 1000) * 1000;

        n = select(sockFd + 1, &rfds, NULL, NULL, &tv);
        if (n < 0) {
            if (errno == EINTR) {
                continue;
            }

            log_error("snmp: select failed host=%s port=%d oid=%s errno=%d",
                      s->host, s->port, oidStr, errno);
            close(sockFd);
            return SWITCHD_ERR_IO;
        }

        if (n == 0) {
            log_error("snmp: timeout host=%s port=%d oid=%s attempt=%d/%d",
                      s->host,
                      s->port,
                      oidStr,
                      attempt + 1,
                      s->retries + 1);
            continue;
        }

        n = recvfrom(sockFd, resp, sizeof(resp), 0, NULL, NULL);
        if (n <= 0) {
            log_error("snmp: recv failed host=%s port=%d oid=%s errno=%d",
                      s->host, s->port, oidStr, errno);
            continue;
        }

        close(sockFd);

        if (parse_response(resp, (size_t)n, reqId, out) == 0) {
            log_debug("snmp: ok pdu=%d oid=%s bytes=%zd",
                      pduType, oidStr, n);
            return SWITCHD_OK;
        }

        log_error("snmp: parse response failed pdu=%d oid=%s bytes=%zd",
                  pduType, oidStr, n);
        return SWITCHD_ERR_SNMP;
    }

    close(sockFd);
    return SWITCHD_ERR_TIMEOUT;
}

int snmp_get(SnmpSession *s, const uint32_t *oid, size_t oidLen, SnmpVarBind *out) {
    return snmp_request(s, SNMP_GET, oid, oidLen, NULL, out);
}

int snmp_get_next(SnmpSession *s, const uint32_t *oid, size_t oidLen, SnmpVarBind *out) {
    return snmp_request(s, SNMP_GETNEXT, oid, oidLen, NULL, out);
}

int snmp_set_integer(SnmpSession *s, const uint32_t *oid, size_t oidLen, int32_t value) {
    SnmpValue v;
    memset(&v, 0, sizeof(v));
    v.type = SNMP_VALUE_INT;
    v.intValue = value;
    return snmp_request(s, SNMP_SET, oid, oidLen, &v, NULL);
}

int snmp_set_string(SnmpSession *s, const uint32_t *oid, size_t oidLen, const char *value) {
    SnmpValue v;
    memset(&v, 0, sizeof(v));
    v.type = SNMP_VALUE_STRING;
    snprintf(v.stringValue, sizeof(v.stringValue), "%s", value);
    return snmp_request(s, SNMP_SET, oid, oidLen, &v, NULL);
}

int snmp_walk(SnmpSession *s, const uint32_t *baseOid, size_t baseLen,
              int (*cb)(const SnmpVarBind *vb, void *arg), void *arg) {
    SnmpVarBind vb;
    uint32_t cur[64];
    size_t curLen = baseLen;
    int ret;
    memcpy(cur, baseOid, baseLen * sizeof(uint32_t));
    for (;;) {
        ret = snmp_get_next(s, cur, curLen, &vb);
        if (ret != SWITCHD_OK) return ret;
        if (!snmp_oid_has_prefix(vb.oid, vb.oidLen, baseOid, baseLen)) return SWITCHD_OK;
        if (cb(&vb, arg) != 0) return SWITCHD_OK;
        memcpy(cur, vb.oid, vb.oidLen * sizeof(uint32_t));
        curLen = vb.oidLen;
    }
}

int snmp_oid_from_string(const char *s, uint32_t *oid, size_t *oidLen) {
    char tmp[256];
    char *tok;
    char *save;
    size_t n = 0;
    if (s == NULL || oid == NULL || oidLen == NULL) return -1;
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

int snmp_oid_has_prefix(const uint32_t *oid, size_t oidLen, const uint32_t *prefix, size_t prefixLen) {
    size_t i;
    if (oidLen < prefixLen) return 0;
    for (i = 0; i < prefixLen; i++) if (oid[i] != prefix[i]) return 0;
    return 1;
}

char *snmp_oid_to_string(const uint32_t *oid, size_t oidLen, char *buf, size_t bufLen) {
    size_t i;
    size_t off = 0;
    for (i = 0; i < oidLen; i++) {
        off += (size_t)snprintf(buf + off, (off < bufLen) ? bufLen - off : 0,
                                "%s%u", (i == 0) ? "" : ".", oid[i]);
        if (off >= bufLen) break;
    }
    return buf;
}
