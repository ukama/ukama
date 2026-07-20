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
#include <unistd.h>

#include "aisg_v2.h"
#include "hdlc.h"
#include "retap.h"
#include "retap_ops.h"
#include "serial.h"
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

static bool test_serial_fill_and_shared_flags(void)
{
    HdlcFrame tx1;
    HdlcFrame tx2;
    HdlcFrame rx;
    SerialPort port;
    uint8_t raw1[HDLC_MAX_FRAME];
    uint8_t raw2[HDLC_MAX_FRAME];
    uint8_t stream[HDLC_MAX_FRAME * 2];
    uint8_t frame[HDLC_MAX_FRAME];
    size_t raw1Len = 0;
    size_t raw2Len = 0;
    size_t streamLen = 0;
    size_t frameLen = 0;
    int fds[2];

    memset(&tx1, 0, sizeof(tx1));
    memset(&tx2, 0, sizeof(tx2));
    memset(&rx, 0, sizeof(rx));
    memset(&port, 0, sizeof(port));

    tx1.address = 0x01;
    tx1.control = hdlc_rr_ctrl(0, true);
    tx2.address = 0x01;
    tx2.control = hdlc_i_ctrl(0, 0, true);
    tx2.info[0] = 0x34;
    tx2.infoLen = 1;

    CHECK(hdlc_encode_frame(&tx1, raw1, sizeof(raw1), &raw1Len));
    CHECK(hdlc_encode_frame(&tx2, raw2, sizeof(raw2), &raw2Len));
    CHECK(pipe(fds) == 0);
    port.fd = fds[0];

    /* Noise and any number of leading/fill flags must be ignored. */
    stream[streamLen++] = 0x55;
    stream[streamLen++] = HDLC_FLAG;
    stream[streamLen++] = HDLC_FLAG;
    memcpy(&stream[streamLen], raw1, raw1Len);
    streamLen += raw1Len;
    CHECK(write(fds[1], stream, streamLen) == (ssize_t)streamLen);
    CHECK(serial_read_frame(&port,
                            frame,
                            sizeof(frame),
                            &frameLen,
                            100));
    CHECK(hdlc_decode_frame(frame, frameLen, &rx));
    CHECK(rx.address == tx1.address);
    CHECK(rx.control == tx1.control);

    /* The previous closing flag is allowed to open the next frame. */
    CHECK(write(fds[1], &raw2[1], raw2Len - 1) == (ssize_t)(raw2Len - 1));
    CHECK(serial_read_frame(&port,
                            frame,
                            sizeof(frame),
                            &frameLen,
                            100));
    CHECK(hdlc_decode_frame(frame, frameLen, &rx));
    CHECK(rx.address == tx2.address);
    CHECK(rx.control == tx2.control);
    CHECK(rx.infoLen == 1 && rx.info[0] == 0x34);

    close(fds[0]);
    close(fds[1]);
    return true;
}

static bool test_xid_scan_and_assignment(void)
{
    HdlcFrame scan;
    uint8_t raw[HDLC_MAX_FRAME];
    uint8_t info[256];
    size_t rawLen = 0;
    size_t infoLen = 0;
    XidAddressingParams params;
    const uint8_t uid[] = { 'U', 'K', 'A', 'M', 'A', '0', '0', '1' };
    uint16_t vendor = ((uint16_t)'U' << 8) | (uint8_t)'K';

    CHECK(xid_build_scan_info(info, sizeof(info), &infoLen));
    CHECK(infoLen == 45);
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

    memset(&scan, 0, sizeof(scan));
    scan.address = AISG_ADDR_BROADCAST;
    scan.control = hdlc_xid_ctrl(true);
    memcpy(scan.info, info, infoLen);
    scan.infoLen = infoLen;
    CHECK(hdlc_encode_frame(&scan, raw, sizeof(raw), &rawLen));
    CHECK(rawLen == 51);
    CHECK(raw[0] == 0x7E && raw[1] == 0xFF && raw[2] == 0xBF);
    CHECK(raw[3] == 0x81 && raw[4] == 0xF0 && raw[5] == 0x2A);
    CHECK(raw[6] == 0x01 && raw[7] == 0x13);
    CHECK(raw[27] == 0x03 && raw[28] == 0x13);
    CHECK(raw[48] == 0x7F && raw[49] == 0x0F && raw[50] == 0x7E);
    CHECK(AISG_HDLC_DEFAULT_FRAME_MAX == 78);
    CHECK(AISG_HDLC_DEFAULT_INFO_MAX == 74);

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

static bool test_real_ret_scan_response_without_pi2(void)
{
    static const uint8_t raw[] = {
        0x7E, 0x00, 0xBF, 0x81, 0xF0, 0x1C, 0x01, 0x13,
        0x54, 0x43, 0x30, 0x30, 0x34, 0x42, 0x4C, 0x32,
        0x33, 0x33, 0x37, 0x59, 0x31, 0x30, 0x30, 0x30,
        0x39, 0x30, 0x31, 0x06, 0x02, 0x54, 0x43, 0x04,
        0x01, 0x01, 0x00, 0x30, 0x7E
    };
    static const uint8_t expectedUid[] = {
        0x54, 0x43, 0x30, 0x30, 0x34, 0x42, 0x4C, 0x32,
        0x33, 0x33, 0x37, 0x59, 0x31, 0x30, 0x30, 0x30,
        0x39, 0x30, 0x31
    };
    HdlcFrame frame;
    XidAddressingParams params;

    memset(&frame, 0, sizeof(frame));
    memset(&params, 0, sizeof(params));

    CHECK(hdlc_decode_frame(raw, sizeof(raw), &frame));
    CHECK(frame.address == AISG_ADDR_DEFAULT);
    CHECK(hdlc_is_xid(frame.control));
    CHECK(xid_parse_addressing_info(frame.info, frame.infoLen, &params));
    CHECK(params.hasUniqueId);
    CHECK(params.uniqueIdLen == sizeof(expectedUid));
    CHECK(bytes_eq(params.uniqueId, expectedUid, sizeof(expectedUid)));
    CHECK(!params.hasAddress);
    CHECK(params.hasVendorCode && params.vendorCode == 0x5443);
    CHECK(params.hasDeviceType &&
          params.deviceType == AISG_DEVICE_TYPE_SINGLE_RET);

    return true;
}

static bool test_real_ret_identical_release_xid(void)
{
    /* Exact command/response seen on the real RET for Release 6. */
    static const uint8_t captured[] = {
        0x7E, 0x01, 0xBF, 0x81, 0xF0, 0x03,
        0x05, 0x01, 0x06, 0xDE, 0xB5, 0x7E
    };
    HdlcFrame frame;
    XidAddressingParams params;
    uint8_t encoded[HDLC_MAX_FRAME];
    size_t infoLen = 0;
    size_t encodedLen = 0;

    memset(&frame, 0, sizeof(frame));
    memset(&params, 0, sizeof(params));

    CHECK(hdlc_decode_frame(captured, sizeof(captured), &frame));
    CHECK(frame.address == AISG_ADDR_ASSIGNED);
    CHECK(hdlc_is_xid(frame.control));
    CHECK(xid_parse_addressing_info(frame.info, frame.infoLen, &params));
    CHECK(params.has3gppRelease);
    CHECK(params.release == AISG_3GPP_RELEASE_ID);

    /* Accepting the offered value produces a byte-identical XID response. */
    memset(&frame, 0, sizeof(frame));
    frame.address = AISG_ADDR_ASSIGNED;
    frame.control = hdlc_xid_ctrl(true);
    CHECK(xid_build_one_octet_info(AISG_XID_PI_3GPP_RELEASE,
                                   AISG_3GPP_RELEASE_ID,
                                   frame.info,
                                   sizeof(frame.info),
                                   &infoLen));
    frame.infoLen = infoLen;
    CHECK(hdlc_encode_frame(&frame,
                            encoded,
                            sizeof(encoded),
                            &encodedLen));
    CHECK(encodedLen == sizeof(captured));
    CHECK(bytes_eq(encoded, captured, sizeof(captured)));

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
        { "serial_fill_and_shared_flags", test_serial_fill_and_shared_flags },
        { "xid_scan_and_assignment",     test_xid_scan_and_assignment },
        { "real_ret_scan_without_pi2",   test_real_ret_scan_response_without_pi2 },
        { "real_ret_identical_release_xid", test_real_ret_identical_release_xid },
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
