/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <unistd.h>

#include "hdlc.h"
#include "ret_mode.h"
#include "ret_model.h"
#include "ret_pty.h"
#include "retap.h"
#include "retap_codes.h"
#include "usys_log.h"
#include "version.h"
#include "xid.h"

#define RET_MODE_RAW_MAX                   HDLC_MAX_FRAME
#define RET_MODE_INFO_MAX                  HDLC_MAX_INFO

static void log_hex(const char *prefix, const uint8_t *buf, size_t len)
{
    char line[512];
    size_t off = 0;
    size_t i;

    if (prefix == NULL || buf == NULL) {
        return;
    }

    for (i = 0; i < len && off + 4 < sizeof(line); i++) {
        off += snprintf(line + off, sizeof(line) - off, "%02X%s", buf[i],
                        (i + 1 == len) ? "" : " ");
    }

    usys_log_debug("%s len=%zu %s%s", prefix, len, line,
                   (i < len) ? " ..." : "");
}

static bool build_info_field(uint8_t *buf, size_t size, size_t *len)
{
    const char *fields[] = { "RET1T1", "EMU-0001", "A1", VERSION };
    size_t off = 0;
    size_t i;
    size_t n;

    if (buf == NULL || len == NULL) {
        return false;
    }

    for (i = 0; i < sizeof(fields) / sizeof(fields[0]); i++) {
        n = strlen(fields[i]);
        if (n > 255 || off + 1 + n > size) {
            return false;
        }
        buf[off++] = (uint8_t)n;
        memcpy(&buf[off], fields[i], n);
        off += n;
    }

    *len = off;
    return true;
}

static bool ret_fail(uint8_t procedure,
                     uint8_t reason,
                     uint8_t *out,
                     size_t outSize,
                     size_t *outLen)
{
    return retap_encode_fail_response(procedure,
                                      reason,
                                      NULL,
                                      0,
                                      out,
                                      outSize,
                                      outLen);
}

static bool ret_ok(uint8_t procedure,
                   const uint8_t *data,
                   size_t dataLen,
                   uint8_t *out,
                   size_t outSize,
                   size_t *outLen)
{
    return retap_encode_ok_response(procedure,
                                    data,
                                    dataLen,
                                    out,
                                    outSize,
                                    outLen);
}

static bool handle_retap(RetModel *model,
                         const uint8_t *info,
                         size_t infoLen,
                         uint8_t *response,
                         size_t responseSize,
                         size_t *responseLen)
{
    RetapRequest request;
    uint8_t data[RET_MODE_INFO_MAX];
    size_t dataLen = 0;
    int16_t target;
    uint16_t raw;

    if (model == NULL || response == NULL || responseLen == NULL) {
        return false;
    }

    if (!retap_decode_request(info, infoLen, &request)) {
        usys_log_debug("ret-emu: malformed RETAP request len=%zu", infoLen);
        return ret_fail(RETAP_PROC_GET_ERROR_STATUS,
                        RETAP_RC_DATA_ERROR,
                        response,
                        responseSize,
                        responseLen);
    }

    usys_log_debug("ret-emu: RETAP request procedure=0x%02X data_len=%zu",
                   request.procedure,
                   request.dataLen);

    if (model->busy && request.procedure != RETAP_PROC_GET_ERROR_STATUS) {
        return ret_fail(request.procedure,
                        RETAP_RC_BUSY,
                        response,
                        responseSize,
                        responseLen);
    }

    switch (request.procedure) {
    case RETAP_PROC_GET_INFORMATION:
        if (request.dataLen != 0) {
            return ret_fail(request.procedure,
                            RETAP_RC_DATA_ERROR,
                            response,
                            responseSize,
                            responseLen);
        }
        if (!build_info_field(data, sizeof(data), &dataLen)) {
            return false;
        }
        return ret_ok(request.procedure, data, dataLen, response, responseSize, responseLen);

    case RETAP_PROC_GET_ERROR_STATUS:
        if (request.dataLen != 0) {
            return ret_fail(request.procedure,
                            RETAP_RC_DATA_ERROR,
                            response,
                            responseSize,
                            responseLen);
        }
        return ret_ok(request.procedure,
                      model->activeErrors,
                      model->activeErrorCount,
                      response,
                      responseSize,
                      responseLen);

    case RETAP_PROC_CLEAR_ACTIVE_ALARMS:
        if (request.dataLen != 0) {
            return ret_fail(request.procedure,
                            RETAP_RC_DATA_ERROR,
                            response,
                            responseSize,
                            responseLen);
        }
        ret_model_clear_errors(model);
        if (!model->configured) {
            ret_model_set_error(model, RETAP_RC_NOT_SCALED);
        }
        if (!model->calibrated) {
            ret_model_set_error(model, RETAP_RC_NOT_CALIBRATED);
        }
        return ret_ok(request.procedure, NULL, 0, response, responseSize, responseLen);

    case RETAP_PROC_ALARM_SUBSCRIBE:
        if (request.dataLen != 0) {
            return ret_fail(request.procedure,
                            RETAP_RC_DATA_ERROR,
                            response,
                            responseSize,
                            responseLen);
        }
        model->alarmSubscribed = true;
        return ret_ok(request.procedure, NULL, 0, response, responseSize, responseLen);

    case RETAP_PROC_SEND_CONFIG_DATA:
        if (request.dataLen == 0 || request.dataLen > RETAP_CONFIG_SEGMENT_MAX) {
            return ret_fail(request.procedure,
                            RETAP_RC_DATA_ERROR,
                            response,
                            responseSize,
                            responseLen);
        }
        model->configured = true;
        ret_model_clear_error(model, RETAP_RC_NOT_SCALED);
        return ret_ok(request.procedure, NULL, 0, response, responseSize, responseLen);

    case RETAP_PROC_CALIBRATE:
        if (request.dataLen != 0) {
            return ret_fail(request.procedure,
                            RETAP_RC_DATA_ERROR,
                            response,
                            responseSize,
                            responseLen);
        }
        if (model->requiresConfig && !model->configured) {
            return ret_fail(request.procedure,
                            RETAP_RC_NOT_SCALED,
                            response,
                            responseSize,
                            responseLen);
        }
        model->calibrated = true;
        ret_model_clear_error(model, RETAP_RC_NOT_CALIBRATED);
        return ret_ok(request.procedure, NULL, 0, response, responseSize, responseLen);

    case RETAP_PROC_GET_TILT:
        if (request.dataLen != 0) {
            return ret_fail(request.procedure,
                            RETAP_RC_DATA_ERROR,
                            response,
                            responseSize,
                            responseLen);
        }
        if (!model->calibrated) {
            return ret_fail(request.procedure,
                            RETAP_RC_NOT_CALIBRATED,
                            response,
                            responseSize,
                            responseLen);
        }
        raw = (uint16_t)model->tiltTenths;
        data[0] = (uint8_t)(raw & 0xFF);
        data[1] = (uint8_t)((raw >> 8) & 0xFF);
        return ret_ok(request.procedure, data, 2, response, responseSize, responseLen);

    case RETAP_PROC_SET_TILT:
        if (request.dataLen != 2) {
            return ret_fail(request.procedure,
                            RETAP_RC_DATA_ERROR,
                            response,
                            responseSize,
                            responseLen);
        }
        if (!model->calibrated) {
            return ret_fail(request.procedure,
                            RETAP_RC_NOT_CALIBRATED,
                            response,
                            responseSize,
                            responseLen);
        }
        raw = (uint16_t)request.data[0] | ((uint16_t)request.data[1] << 8);
        target = (int16_t)raw;
        if (target < model->minTiltTenths || target > model->maxTiltTenths) {
            return ret_fail(request.procedure,
                            RETAP_RC_OUT_OF_RANGE,
                            response,
                            responseSize,
                            responseLen);
        }
        model->tiltTenths = target;
        return ret_ok(request.procedure, NULL, 0, response, responseSize, responseLen);

    case RETAP_PROC_RESET_SOFTWARE:
        ret_model_reset_l2(model);
        model->configured = !model->requiresConfig;
        model->calibrated = !model->requiresConfig;
        ret_model_clear_errors(model);
        if (!model->configured) {
            ret_model_set_error(model, RETAP_RC_NOT_SCALED);
        }
        if (!model->calibrated) {
            ret_model_set_error(model, RETAP_RC_NOT_CALIBRATED);
        }
        return ret_ok(request.procedure, NULL, 0, response, responseSize, responseLen);

    default:
        return ret_fail(request.procedure,
                        RETAP_RC_UNKNOWN_PROCEDURE,
                        response,
                        responseSize,
                        responseLen);
    }
}

static bool send_hdlc_response(int fd,
                               uint8_t address,
                               uint8_t control,
                               const uint8_t *info,
                               size_t infoLen)
{
    uint8_t raw[RET_MODE_RAW_MAX];
    size_t rawLen = 0;

    if (!hdlc_encode_addr_info(address,
                               control,
                               info,
                               infoLen,
                               raw,
                               sizeof(raw),
                               &rawLen)) {
        usys_log_debug("ret-emu: failed to encode HDLC response");
        return false;
    }

    log_hex("ret-emu TX", raw, rawLen);
    return ret_pty_write_all(fd, raw, rawLen);
}

static bool send_xid_device_response(int fd, const RetModel *model)
{
    uint8_t info[RET_MODE_INFO_MAX];
    size_t infoLen = 0;

    if (!xid_build_device_response_info(model->uniqueId,
                                        model->uniqueIdLen,
                                        model->address,
                                        model->deviceType,
                                        true,
                                        model->vendorCode,
                                        info,
                                        sizeof(info),
                                        &infoLen)) {
        return false;
    }

    return send_hdlc_response(fd,
                              model->address,
                              hdlc_xid_ctrl(true),
                              info,
                              infoLen);
}

static bool send_one_octet_xid_response(int fd,
                                        const RetModel *model,
                                        uint8_t pi,
                                        uint8_t value)
{
    uint8_t info[RET_MODE_INFO_MAX];
    size_t infoLen = 0;

    if (!xid_build_one_octet_info(pi, value, info, sizeof(info), &infoLen)) {
        return false;
    }

    return send_hdlc_response(fd,
                              model->address,
                              hdlc_xid_ctrl(true),
                              info,
                              infoLen);
}

static bool handle_xid_frame(int fd, RetModel *model, const HdlcFrame *rx)
{
    XidAddressingParams params;

    if (model == NULL || rx == NULL) {
        return false;
    }

    if (!xid_parse_addressing_info(rx->info, rx->infoLen, &params)) {
        usys_log_debug("ret-emu: ignoring malformed XID info_len=%zu", rx->infoLen);
        return true;
    }

    if (params.hasMask) {
        if (model->state != RET_L2_NO_ADDRESS ||
            rx->address != RET_EMU_ADDR_BROADCAST ||
            !params.hasUniqueId || params.uniqueIdLen != params.maskLen) {
            usys_log_debug("ret-emu: ignoring invalid scan state=%s addr=0x%02X",
                           ret_l2_state_name(model->state), rx->address);
            return true;
        }

        if (!xid_unique_id_mask_match(model->uniqueId,
                                      model->uniqueIdLen,
                                      params.uniqueId,
                                      params.uniqueIdLen,
                                      params.mask,
                                      params.maskLen)) {
            usys_log_debug("ret-emu: XID scan did not match unique-id");
            return true;
        }

        usys_log_info("ret-emu: XID scan match unique_id_len=%zu", model->uniqueIdLen);
        return send_xid_device_response(fd, model);
    }

    if (params.hasAddress) {
        if (model->state != RET_L2_NO_ADDRESS || rx->address != RET_EMU_ADDR_BROADCAST) {
            usys_log_debug("ret-emu: ignoring address assignment state=%s addr=0x%02X",
                           ret_l2_state_name(model->state), rx->address);
            return true;
        }

        if (params.address == RET_EMU_ADDR_DEFAULT ||
            params.address == RET_EMU_ADDR_BROADCAST) {
            usys_log_debug("ret-emu: invalid assigned address=0x%02X", params.address);
            return true;
        }

        if (!xid_assignment_matches(&params,
                                    model->uniqueId,
                                    model->uniqueIdLen,
                                    model->deviceType,
                                    model->vendorCode)) {
            usys_log_debug("ret-emu: address assignment did not match device");
            return true;
        }

        model->address = params.address;
        model->state = RET_L2_ADDRESS_ASSIGNED;
        model->primaryNsExpected = 0;
        model->secondaryNs = 0;

        usys_log_info("ret-emu: address assigned addr=0x%02X state=%s",
                      model->address,
                      ret_l2_state_name(model->state));

        return send_xid_device_response(fd, model);
    }

    if (rx->address != model->address || model->state == RET_L2_NO_ADDRESS) {
        usys_log_debug("ret-emu: ignoring addressed XID state=%s addr=0x%02X expected=0x%02X",
                       ret_l2_state_name(model->state), rx->address, model->address);
        return true;
    }

    if (params.has3gppRelease) {
        if (params.release != RET_EMU_3GPP_RELEASE_ID) {
            usys_log_debug("ret-emu: unsupported 3GPP release=0x%02X", params.release);
            return true;
        }
        model->has3gppRelease = true;
        model->negotiated3gppRelease = params.release;
        return send_one_octet_xid_response(fd,
                                           model,
                                           AISG_XID_PI_3GPP_RELEASE,
                                           params.release);
    }

    if (params.hasAisgVersion) {
        if (params.aisgVersion != RET_EMU_AISG_VERSION) {
            usys_log_debug("ret-emu: unsupported AISG version=0x%02X", params.aisgVersion);
            return true;
        }
        model->hasAisgVersion = true;
        model->negotiatedAisgVersion = params.aisgVersion;
        return send_one_octet_xid_response(fd,
                                           model,
                                           AISG_XID_PI_AISG_VERSION,
                                           params.aisgVersion);
    }

    usys_log_debug("ret-emu: unsupported XID parameter set ignored");
    return true;
}

static bool handle_i_frame(int fd, RetModel *model, const HdlcFrame *rx)
{
    uint8_t responseInfo[RETAP_MAX_ENCODED];
    size_t responseInfoLen = 0;
    uint8_t rxNs;
    uint8_t rxNr;
    uint8_t ctrl;

    if (model == NULL || rx == NULL) {
        return false;
    }

    if (model->state != RET_L2_CONNECTED || rx->address != model->address) {
        usys_log_debug("ret-emu: I-frame before connected ignored state=%s addr=0x%02X expected=0x%02X",
                       ret_l2_state_name(model->state), rx->address, model->address);
        return true;
    }

    rxNs = hdlc_ns(rx->control);
    rxNr = hdlc_nr(rx->control);

    if (rxNs != model->primaryNsExpected) {
        usys_log_debug("ret-emu: I-frame sequence error ns=%u expected=%u",
                       rxNs,
                       model->primaryNsExpected);
        return send_hdlc_response(fd,
                                  model->address,
                                  hdlc_frmr_ctrl(true),
                                  NULL,
                                  0);
    }

    if (rxNr != model->secondaryNs) {
        usys_log_debug("ret-emu: I-frame primary ack nr=%u expected=%u",
                       rxNr,
                       model->secondaryNs);
        return send_hdlc_response(fd,
                                  model->address,
                                  hdlc_frmr_ctrl(true),
                                  NULL,
                                  0);
    }

    if (!handle_retap(model,
                      rx->info,
                      rx->infoLen,
                      responseInfo,
                      sizeof(responseInfo),
                      &responseInfoLen)) {
        return false;
    }

    ctrl = hdlc_i_ctrl(model->secondaryNs,
                       (uint8_t)((rxNs + 1) & 0x07),
                       true);

    if (!send_hdlc_response(fd,
                            model->address,
                            ctrl,
                            responseInfo,
                            responseInfoLen)) {
        return false;
    }

    model->primaryNsExpected = (uint8_t)((rxNs + 1) & 0x07);
    model->secondaryNs = (uint8_t)((model->secondaryNs + 1) & 0x07);

    return true;
}

static bool handle_frame(int fd, RetModel *model, const HdlcFrame *rx)
{
    if (model == NULL || rx == NULL) {
        return false;
    }

    usys_log_debug("ret-emu: RX addr=0x%02X ctrl=0x%02X info_len=%zu state=%s",
                   rx->address,
                   rx->control,
                   rx->infoLen,
                   ret_l2_state_name(model->state));

    if (hdlc_is_xid(rx->control)) {
        return handle_xid_frame(fd, model, rx);
    }

    if (hdlc_is_snrm(rx->control)) {
        if (model->state != RET_L2_ADDRESS_ASSIGNED || rx->address != model->address) {
            usys_log_debug("ret-emu: SNRM ignored state=%s addr=0x%02X expected=0x%02X",
                           ret_l2_state_name(model->state), rx->address, model->address);
            return true;
        }

        model->state = RET_L2_CONNECTED;
        model->primaryNsExpected = 0;
        model->secondaryNs = 0;

        usys_log_info("ret-emu: SNRM accepted state=CONNECTED addr=0x%02X",
                      model->address);
        return send_hdlc_response(fd,
                                  model->address,
                                  hdlc_ua_ctrl(true),
                                  NULL,
                                  0);
    }

    if (hdlc_is_disc(rx->control)) {
        if (rx->address == model->address) {
            send_hdlc_response(fd,
                               model->address,
                               hdlc_ua_ctrl(true),
                               NULL,
                               0);
            ret_model_reset_l2(model);
        }
        return true;
    }

    if (hdlc_is_i_frame(rx->control)) {
        return handle_i_frame(fd, model, rx);
    }

    if (rx->address == model->address && model->state != RET_L2_NO_ADDRESS) {
        return send_hdlc_response(fd,
                                  model->address,
                                  hdlc_frmr_ctrl(true),
                                  NULL,
                                  0);
    }

    return true;
}

bool ret_mode_run(const EmuConfig *config, volatile sig_atomic_t *running)
{
    RetModel model;
    int fd = -1;
    char slaveName[128];
    uint8_t raw[RET_MODE_RAW_MAX];
    size_t rawLen;
    HdlcFrame frame;

    if (config == NULL || running == NULL) {
        return false;
    }

    ret_model_init(&model,
                   config->retVendorCode,
                   config->retSerial,
                   config->retRequiresConfig,
                   (int16_t)config->retInitialTiltTenths,
                   (int16_t)config->retMinTiltTenths,
                   (int16_t)config->retMaxTiltTenths);

    if (!ret_pty_open(config->retPtyPath, &fd, slaveName, sizeof(slaveName))) {
        return false;
    }

    usys_log_info("ret-emu: mode=ret vendor=%s serial=%s uid_len=%zu requires_config=%d tilt=%d range=[%d,%d]",
                  model.vendorCodeStr,
                  model.serial,
                  model.uniqueIdLen,
                  model.requiresConfig ? 1 : 0,
                  model.tiltTenths,
                  model.minTiltTenths,
                  model.maxTiltTenths);

    while (*running) {
        rawLen = 0;
        memset(raw, 0, sizeof(raw));
        memset(&frame, 0, sizeof(frame));

        if (!ret_pty_read_hdlc_frame(fd,
                                     running,
                                     raw,
                                     sizeof(raw),
                                     &rawLen)) {
            if (*running) {
                usys_log_warn("ret-emu: failed to read HDLC frame");
            }
            continue;
        }

        log_hex("ret-emu RX", raw, rawLen);

        if (!hdlc_decode_frame(raw, rawLen, &frame)) {
            usys_log_debug("ret-emu: invalid HDLC frame ignored");
            continue;
        }

        if (!handle_frame(fd, &model, &frame)) {
            usys_log_warn("ret-emu: failed to process frame");
        }
    }

    ret_pty_close(fd, config->retPtyPath);
    usys_log_info("ret-emu: stopped");

    return true;
}
