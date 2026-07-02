/*
 * AISG protocol golden tests.
 *
 * Scope: pure HDLC / XID / RETAP helpers used by ctrl and aisg-emu --mode ret.
 * These tests are deliberately independent of the Ukama platform build tree.
 */

#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>
#include <string.h>

#include "hdlc.h"
#include "retap.h"
#include "retap_ops.h"
#include "xid.h"

#define CHECK(expr) do {                                                       \
    if (!(expr)) {                                                             \
        fprintf(stderr, "FAIL %s:%d: %s\n", __FILE__, __LINE__, #expr);       \
        return false;                                                          \
    }                                                                          \
} while (0)

static bool bytes_eq(const uint8_t *a, const uint8_t *b, size_t len)
{
    return memcmp(a, b, len) == 0;
}

static bool test_hdlc_roundtrip_and_escaping(void)
{
    HdlcFrame tx;
    HdlcFrame rx;
    uint8_t raw[HDLC_MAX_FRAME];
    size_t rawLen = 0;
    bool sawEsc7e = false;
    bool sawEsc7d = false;
    size_t i;

    memset(&tx, 0, sizeof(tx));
    memset(&rx, 0, sizeof(rx));

    tx.address = 0x01;
    tx.control = hdlc_i_ctrl(3, 2, true);
    tx.info[0] = 0x34;
    tx.info[1] = 0x7E;
    tx.info[2] = 0x7D;
    tx.info[3] = 0x00;
    tx.infoLen = 4;

    CHECK(hdlc_encode_frame(&tx, raw, sizeof(raw), &rawLen));
    CHECK(rawLen >= 8);
    CHECK(raw[0] == HDLC_FLAG);
    CHECK(raw[rawLen - 1] == HDLC_FLAG);

    for (i = 0; i + 1 < rawLen; i++) {
        if (raw[i] == HDLC_ESCAPE && raw[i + 1] == (uint8_t)(0x7E ^ HDLC_ESCAPE_XOR)) {
            sawEsc7e = true;
        }
        if (raw[i] == HDLC_ESCAPE && raw[i + 1] == (uint8_t)(0x7D ^ HDLC_ESCAPE_XOR)) {
            sawEsc7d = true;
        }
    }

    CHECK(sawEsc7e);
    CHECK(sawEsc7d);
    CHECK(hdlc_decode_frame(raw, rawLen, &rx));
    CHECK(rx.address == tx.address);
    CHECK(rx.control == tx.control);
    CHECK(rx.infoLen == tx.infoLen);
    CHECK(bytes_eq(rx.info, tx.info, tx.infoLen));
    CHECK(hdlc_is_i_frame(rx.control));
    CHECK(hdlc_ns(rx.control) == 3);
    CHECK(hdlc_nr(rx.control) == 2);
    CHECK(hdlc_poll_final(rx.control));

    raw[rawLen - 3] ^= 0x01; /* corrupt FCS/data area */
    CHECK(!hdlc_decode_frame(raw, rawLen, &rx));

    return true;
}

static bool test_xid_scan_and_assignment(void)
{
    uint8_t info[256];
    size_t infoLen = 0;
    XidAddressingParams params;
    const uint8_t uid[] = { 'U', 'K', 'A', 'M', 'A', '0', '0', '1' };
    uint16_t vendor = ((uint16_t)'U' << 8) | (uint8_t)'K';

    CHECK(xid_build_scan_info(info, sizeof(info), &infoLen));
    CHECK(xid_parse_addressing_info(info, infoLen, &params));
    CHECK(params.hasUniqueId);
    CHECK(params.uniqueIdLen == AISG_XID_SCAN_ID_LEN);
    CHECK(params.hasMask);
    CHECK(params.maskLen == AISG_XID_SCAN_ID_LEN);
    CHECK(!params.hasAddress);
    CHECK(xid_unique_id_mask_match(uid,
                                   sizeof(uid),
                                   params.uniqueId,
                                   params.uniqueIdLen,
                                   params.mask,
                                   params.maskLen));

    CHECK(xid_build_assign_info(uid,
                                sizeof(uid),
                                0x01,
                                0x01,
                                true,
                                vendor,
                                info,
                                sizeof(info),
                                &infoLen));
    CHECK(xid_parse_addressing_info(info, infoLen, &params));
    CHECK(params.hasUniqueId);
    CHECK(params.hasAddress);
    CHECK(params.address == 0x01);
    CHECK(params.hasDeviceType);
    CHECK(params.deviceType == 0x01);
    CHECK(params.hasVendorCode);
    CHECK(params.vendorCode == vendor);
    CHECK(!params.hasMask);
    CHECK(xid_assignment_matches(&params, uid, sizeof(uid), 0x01, vendor));
    CHECK(!xid_assignment_matches(&params, uid, sizeof(uid), 0x11, vendor));

    /* Address assignment containing PI=3/bit-mask must not match. */
    CHECK(xid_build_scan_info(info, sizeof(info), &infoLen));
    CHECK(xid_parse_addressing_info(info, infoLen, &params));
    CHECK(!xid_assignment_matches(&params, uid, sizeof(uid), 0x01, vendor));

    return true;
}

static bool test_retap_golden_packets(void)
{
    RetapRequest req;
    RetapResponse resp;
    uint8_t buf[RETAP_MAX_ENCODED];
    uint8_t expected[8];
    size_t len = 0;
    int16_t tilt;
    uint8_t payload[2];

    CHECK(retap_build_get_information(&req));
    CHECK(retap_encode_request(&req, buf, sizeof(buf), &len));
    expected[0] = RETAP_PROC_GET_INFORMATION;
    expected[1] = 0x00;
    expected[2] = 0x00;
    CHECK(len == 3);
    CHECK(bytes_eq(buf, expected, len));

    CHECK(retap_build_get_tilt(&req));
    CHECK(retap_encode_request(&req, buf, sizeof(buf), &len));
    expected[0] = RETAP_PROC_GET_TILT;
    expected[1] = 0x00;
    expected[2] = 0x00;
    CHECK(len == 3);
    CHECK(bytes_eq(buf, expected, len));

    CHECK(retap_build_set_tilt(&req, 32));
    CHECK(retap_encode_request(&req, buf, sizeof(buf), &len));
    expected[0] = RETAP_PROC_SET_TILT;
    expected[1] = 0x02;
    expected[2] = 0x00;
    expected[3] = 0x20;
    expected[4] = 0x00;
    CHECK(len == 5);
    CHECK(bytes_eq(buf, expected, len));

    /* GetTilt response: proc, len=3, OK, low, high. */
    buf[0] = RETAP_PROC_GET_TILT;
    buf[1] = 0x03;
    buf[2] = 0x00;
    buf[3] = RETAP_RETURN_OK;
    buf[4] = 0x20;
    buf[5] = 0x00;
    CHECK(retap_decode_response(buf, 6, &resp));
    CHECK(retap_response_is_ok(&resp));
    CHECK(resp.procedure == RETAP_PROC_GET_TILT);
    CHECK(resp.dataLen == 2);
    CHECK(retap_parse_get_tilt(&resp, &tilt));
    CHECK(tilt == 32);

    /* Failure response: proc, len=2, FAIL, reason. */
    buf[0] = RETAP_PROC_SET_TILT;
    buf[1] = 0x02;
    buf[2] = 0x00;
    buf[3] = RETAP_RETURN_FAIL;
    buf[4] = RETAP_RC_NOT_CALIBRATED;
    CHECK(retap_decode_response(buf, 5, &resp));
    CHECK(retap_response_is_fail(&resp));
    CHECK(resp.failureReason == RETAP_RC_NOT_CALIBRATED);
    CHECK(retap_failure_to_ctrl_code(resp.failureReason) == CtrlCodeNotCalibrated);

    /* Secondary-side OK response builder. */
    payload[0] = 0x20;
    payload[1] = 0x00;
    CHECK(retap_encode_ok_response(RETAP_PROC_GET_TILT,
                                   payload,
                                   sizeof(payload),
                                   buf,
                                   sizeof(buf),
                                   &len));
    CHECK(len == 6);
    CHECK(buf[0] == RETAP_PROC_GET_TILT);
    CHECK(buf[1] == 0x03);
    CHECK(buf[2] == 0x00);
    CHECK(buf[3] == RETAP_RETURN_OK);
    CHECK(buf[4] == 0x20);
    CHECK(buf[5] == 0x00);

    CHECK(retap_decode_request(expected, 5, &req));
    CHECK(req.procedure == RETAP_PROC_SET_TILT);
    CHECK(req.dataLen == 2);
    CHECK(req.data[0] == 0x20 && req.data[1] == 0x00);

    CHECK(!retap_decode_request(expected, 1, &req));
    CHECK(!retap_decode_response(expected, 1, &resp));

    return true;
}

static bool test_retap_config_limits(void)
{
    RetapRequest req;
    uint8_t data[RETAP_CONFIG_SEGMENT_MAX + 1];

    memset(data, 0xA5, sizeof(data));
    CHECK(retap_build_send_configuration_data(&req,
                                              data,
                                              RETAP_CONFIG_SEGMENT_MAX));
    CHECK(req.procedure == RETAP_PROC_SEND_CONFIG_DATA);
    CHECK(req.dataLen == RETAP_CONFIG_SEGMENT_MAX);
    CHECK(!retap_build_send_configuration_data(&req,
                                               data,
                                               RETAP_CONFIG_SEGMENT_MAX + 1));
    return true;
}

int main(void)
{
    struct {
        const char *name;
        bool (*fn)(void);
    } tests[] = {
        { "hdlc_roundtrip_and_escaping", test_hdlc_roundtrip_and_escaping },
        { "xid_scan_and_assignment",     test_xid_scan_and_assignment },
        { "retap_golden_packets",        test_retap_golden_packets },
        { "retap_config_limits",         test_retap_config_limits },
    };
    size_t i;

    for (i = 0; i < sizeof(tests) / sizeof(tests[0]); i++) {
        if (!tests[i].fn()) {
            fprintf(stderr, "not ok - %s\n", tests[i].name);
            return 1;
        }
        printf("ok - %s\n", tests[i].name);
    }

    printf("AISG protocol golden tests passed\n");
    return 0;
}
